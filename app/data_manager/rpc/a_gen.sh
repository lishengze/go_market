#!/bin/bash
#goctl api go -api client.api -dir . --home ~/.goctl/1.2.5/api

script_path="$(cd "$(dirname "$0")" && pwd)"
project_path="$script_path"/../../../..
echo "$script_path"
echo "$project_path"

goctl rpc protoc marketData.proto --go_out=./types --go-grpc_out=./types --zrpc_out=.   --home="$project_path"/deploy/goctl/1.3.4 --style=goZero
