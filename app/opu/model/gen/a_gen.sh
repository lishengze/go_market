#!/bin/bash

script_path="$(cd "$(dirname "$0")" && pwd)"
project_path="$script_path"/../../../..
#echo "$script_path"
#echo "$project_path"
#cd "$project_path" && pwd && cd -

cd "$script_path" && rm -rf ../*model.go && rm -rf ../vars.go

goctl model mysql ddl -c -src ./create_table.sql -dir ../ -home "$project_path"/pkg/goctl/template/1.2.4-cli && git add ./..
