package conf

import "time"

type Config struct {
	IP            string
	NetServerType string
	SerialType    string
	NacosIP       string
}

type AggregateConfig struct {
	DepthAggregatorMillsecsMap map[string]time.Duration
}

type HedgeConfig struct {
	FeeKind  int
	FeeValue float64
}

type RiskCtrlConfig struct {
	HedgeConfigMap map[string]HedgeConfig

	PricePrecison  uint32
	VolumePrecison uint32

	PriceBiasValue float64
	PriceBiasKind  int

	VolumeBiasValue float64
	VolumeBiasKind  int

	PriceMinumChange float64
}
