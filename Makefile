# This Makefile is meant to be used by people that do not usually work
# with Go source code. If you know what GOPATH is then you probably
# don't need to bother with make.

PHONY += all docker clean

include main.mk

TARGETS := $(sort $(notdir $(wildcard ./cmd/*)))
PHONY += $(TARGETS)

all: $(TARGETS)

.SECONDEXPANSION:
$(TARGETS): $(addprefix $(GOBIN)/,$$@)

$(GOBIN):
	@mkdir -p $@

$(GOBIN)/%: $(GOBIN) FORCE
	@go build -v -o $@ ./cmd/$(notdir $@)
	@echo "Done building."
	@echo "Run \"$(subst $(CURDIR),.,$@)\" to launch $(notdir $@)."

PROTOC_INCLUDES := \
		-I$(CURDIR)/vendor/github.com/gogo/protobuf/types \
		-I$(CURDIR)/vendor/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
		-I$(GOPATH)/src

GRPC_PROTOS := \
	service/pb/*.proto

service-grpc: FORCE
	@protoc $(PROTOC_INCLUDES) \
		--gofast_out=plugins=grpc,\
Mgoogle/protobuf/any.proto=github.com/gogo/protobuf/types,\
Mgoogle/protobuf/duration.proto=github.com/gogo/protobuf/types,\
Mgoogle/protobuf/struct.proto=github.com/gogo/protobuf/types,\
Mgoogle/protobuf/timestamp.proto=github.com/gogo/protobuf/types,\
Mgoogle/protobuf/wrappers.proto=github.com/gogo/protobuf/types:$(GOPATH)/src \
		$(addprefix $(CURDIR)/,$(GRPC_PROTOS))

	@protoc $(PROTOC_INCLUDES) \
		--grpc-gateway_out=logtostderr=true:$(GOPATH)/src $(addprefix $(CURDIR)/,$(GRPC_PROTOS))

migration-%:
	@$(MAKE) -f migration/Makefile $@

coverage.txt:
	@touch $@

test: coverage.txt FORCE
	@for d in `go list ./... | grep -v vendor | grep -v mock`; do		\
		go test -v -coverprofile=profile.out -covermode=atomic $$d;	\
		if [ $$? -eq 0 ]; then						\
			echo "\033[32mPASS\033[0m:\t$$d";			\
			if [ -f profile.out ]; then				\
				cat profile.out >> coverage.txt;		\
				rm profile.out;					\
			fi							\
		else								\
			echo "\033[31mFAIL\033[0m:\t$$d";			\
			exit -1;						\
		fi								\
	done;

contracts: FORCE
	$(shell solc contracts/erc20.sol --bin --abi --optimize --overwrite --output-dir contracts)
	$(shell abigen --type ERC20Token --abi contracts/ERC20Token.abi -bin contracts/ERC20Token.bin -out contracts/erc20_token.go --pkg contracts)
	$(shell abigen --type MithrilToken --abi contracts/MithrilToken.abi -bin contracts/MithrilToken.bin -out contracts/mithril_token.go --pkg contracts)

%-docker:
	@docker build -f ./cmd/$(subst -docker,,$@)/Dockerfile -t $(DOCKER_IMAGE)-$(subst -docker,,$@):$(REV) .
	@docker tag $(DOCKER_IMAGE)-$(subst -docker,,$@):$(REV) $(DOCKER_IMAGE)-$(subst -docker,,$@):latest

%-docker.push:
	@docker push $(DOCKER_IMAGE)-$(subst -docker.push,,$@):$(REV)
	@docker push $(DOCKER_IMAGE)-$(subst -docker.push,,$@):latest

clean:
	rm -fr $(GOBIN)/*

PHONY: help
help:
	@echo  'Generic targets:'
	@echo  '  service                       - Build indexer service'
	@echo  ''
	@echo  'Execute "make" or "make all" to build all targets marked with [*] '
	@echo  'For further info see the ./README.md file'

.PHONY: $(PHONY)

.PHONY: FORCE
FORCE:
