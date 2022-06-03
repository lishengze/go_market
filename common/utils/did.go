package utils

import (
	"github.com/rpcxio/did/snowflake"
)

var (
	IDDid = &Did{}
)

type Did struct {
	Node *snowflake.Node
}

type Snowflake struct {
	Node     int64 `json:"node"`
	Epoch    int64 `json:"epoch"`
	NodeBits uint8 `json:"nodeBits"`
}

func NewDid(snowflakeConfig *Snowflake) *Did {
	node, err := NewSnowflakeNode(snowflakeConfig)
	if err != nil {
		panic("NewDid panic")
	}
	return &Did{Node: node}
}

func NewSnowflakeNode(c *Snowflake) (node *snowflake.Node, err error) {
	return snowflake.NewNode(c.Node, c.Epoch, c.NodeBits, 22-uint8(c.NodeBits))
}

func (d *Did) Generate() int64 {
	return d.Node.Generate()
}

func (d *Did) GetNode() *snowflake.Node {
	return d.Node
}
