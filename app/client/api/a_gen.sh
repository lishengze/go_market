#!/bin/bash

script_path="$(cd "$(dirname "$0")" && pwd)"
project_path="$script_path"/../../..
echo "$script_path"
echo "$project_path"

goctl  api go -api client.api -dir .  --home="$project_path"/deploy/goctl/1.3.4 -style=goZero

