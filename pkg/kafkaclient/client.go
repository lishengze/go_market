package kafkaclient

import (
	"github.com/segmentio/kafka-go"
)

type KafkaClient struct {
	*kafka.Conn
}

func a() {
	//kafka.Dial()
	//kafka.DialLeader()
}
