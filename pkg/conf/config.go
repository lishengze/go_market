package conf

import (
	"fmt"
	"sync"
	"time"
)

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

type TestConfig struct {
	FEE_RISKCTRL_OPEN       bool
	BIAS_RISKCTRL_OPEN      bool
	WATERMARK_RISKCTRL_OPEN bool
	PRICESION_RISKCTRL_OPEN bool
}

func (t *TestConfig) Init() {
	t.FEE_RISKCTRL_OPEN = true
	t.BIAS_RISKCTRL_OPEN = true
	t.WATERMARK_RISKCTRL_OPEN = true
	t.PRICESION_RISKCTRL_OPEN = true
}

var g_single_testconfig *TestConfig
var lock = &sync.Mutex{}

func TESTCONFIG_INIT(file_name string) error {
	test_config := TESTCONFIG()

	test_config.Init()

	return nil
}

func TESTCONFIG() *TestConfig {
	if g_single_testconfig == nil {
		lock.Lock()
		defer lock.Unlock()

		if g_single_testconfig == nil {
			g_single_testconfig = new(TestConfig)
			fmt.Println("Init Single")
		} else {
			fmt.Println("Second Judge")
		}
	} else {
		fmt.Println("Single already created!")
	}

	return g_single_testconfig
}
