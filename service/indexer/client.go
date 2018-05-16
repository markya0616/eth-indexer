// Copyright 2018 AMIS Technologies
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package indexer

import (
	"context"
	"fmt"
	"math/big"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

//go:generate mockery -name EthClient

type EthClient interface {
	BlockByNumber(ctx context.Context, number *big.Int) (*types.Block, error)
	BlockByHash(ctx context.Context, hash common.Hash) (*types.Block, error)
	TransactionReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error)
	TransactionByHash(ctx context.Context, hash common.Hash) (tx *types.Transaction, isPending bool, err error)
	SubscribeNewHead(ctx context.Context, ch chan<- *types.Header) (ethereum.Subscription, error)
	DumpBlock(ctx context.Context, blockNr int64) (*state.Dump, error)
	BalanceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (*big.Int, error)
	CallContract(ctx context.Context, msg ethereum.CallMsg, blockNumber *big.Int) ([]byte, error)
	ModifiedAccountStatesByNumber(ctx context.Context, num uint64) (*state.DirtyDump, error)
	CodeAt(ctx context.Context, account common.Address, blockNumber *big.Int) ([]byte, error)
	Close()
}

type client struct {
	*ethclient.Client
	rpc *rpc.Client
}

func NewClient(url string) (EthClient, error) {
	rpcClient, err := rpc.Dial(url)
	if err != nil {
		return nil, err
	}
	return &client{
		Client: ethclient.NewClient(rpcClient),
		rpc:    rpcClient,
	}, nil
}

func (c *client) DumpBlock(ctx context.Context, blockNr int64) (*state.Dump, error) {
	r := &state.Dump{}
	err := c.rpc.CallContext(ctx, r, "debug_dumpBlock", fmt.Sprintf("0x%x", blockNr))
	return r, err
}

func (c *client) ModifiedAccountStatesByNumber(ctx context.Context, num uint64) (*state.DirtyDump, error) {
	r := &state.DirtyDump{}
	err := c.rpc.CallContext(ctx, r, "debug_getModifiedAccountStatesByNumber", num)
	return r, err
}

func (c *client) Close() {
	c.rpc.Close()
}
