#!/bin/bash

script_path="$(cd "$(dirname "$0")" && pwd)"
project_path="$script_path"/../../..
echo "$script_path"
echo "$project_path"

goctl model mysql ddl -c -src ddl/create_table.sql -dir . --home="$project_path"/deploy/goctl/1.3.4 --style=goZero
