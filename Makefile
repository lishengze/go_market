GO     ?= go
GOFMT  ?= gofmt
GOPATH:=$(shell go env GOPATH)
VERSION=$(shell git describe --tags --always)
PROTO_FILES=$(shell find app/*/api -name *.proto)
SERVICES=$(shell ls -d app/*/ | cut -f2 -d'/')
GOCTL_TEMPLATE_HOME=./deploy/goctl/1.3.4
#SERVICES = $(foreach p, $(RPC_SERVICES),$(p))


.PHONY: init
# init env
init:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install github.com/zeromicro/go-zero/tools/goctl@latest
	go install github.com/zeromicro/goctl-go-compact@latest
	go mod tidy

.PHONY: all
# 生成全部的rpc 和 api的服务
all: init rpc

.PHONY: rpc
# generate rpc server
rpc:
	@@echo Build rpc with `$(GO) version` and `goctl -version`
	for pf in $$(find app/*/api -name *.proto); do \
		service=$$(ls -d $$pf | cut -f2 -d'/'); \
		echo Build $$service; \
		goctl rpc protoc --proto_path=. \
         	       --go_out=./app/$$service/api \
         	       --go-grpc_out=./app/$$service/api \
         	       --zrpc_out=./app/$$service \
         	       --home=$(GOCTL_TEMPLATE_HOME)/rpc \
         	       --style=go_zero \
        		  $$pf; \
	done

# show help
help:
	@echo ''
	@echo 'Usage:'
	@echo ' make [target]'
	@echo ''
	@echo 'Targets:'
	@awk '/^[a-zA-Z\-\_0-9]+:/ { \
	helpMessage = match(lastLine, /^# (.*)/); \
		if (helpMessage) { \
			helpCommand = substr($$1, 0, index($$1, ":")-1); \
			helpMessage = substr(lastLine, RSTART + 2, RLENGTH); \
			printf "\033[36m%-22s\033[0m %s\n", helpCommand,helpMessage; \
		} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST)

.DEFAULT_GOAL := help
