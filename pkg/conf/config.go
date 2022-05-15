package conf

import "time"

type Config struct {
	IP            string
	NetServerType string
	SerialType    string
}

type AggregateConfig struct {
	DepthAggregatorMillsecs time.Duration
}
