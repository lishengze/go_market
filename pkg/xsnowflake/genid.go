package xsnowflake

import (
	"github.com/bwmarrin/snowflake"
)

var node *snowflake.Node

func init() {
	n, err := snowflake.NewNode(1)
	if err != nil {
		panic(err)
	}
	node = n
}

func GenerateId() snowflake.ID {
	return node.Generate()
}
