package config

import (
	"fmt"
	"sync"
	"time"

	"github.com/zeromicro/go-zero/core/conf"
)

var g_single_native_config *NativeConfig
var lock = &sync.Mutex{}

type Config struct {
	IP            string
	NetServerType string
	SerialType    string
	NacosIP       string
}

type RiskCtlTestConfig struct {
	FeeRiskctrlOpen       bool
	BiasRiskctrlOpen      bool
	WatermarkRiskctrlOpen bool
	PricesionRiskctrlOpen bool
}

type NacosConfig struct {
	IpAddr              string
	Port                int32
	NamespaceId         string
	TimeoutMs           int
	NotLoadCacheAtStart bool
	LogDir              string
	CacheDir            string
	RotateTime          string
	MaxAge              int32
	LogLevel            string
}
type NativeConfig struct {
	IP            string
	NetServerType string
	SerialType    string

	Nacos          NacosConfig
	RiskTestConfig RiskCtlTestConfig
}

func (n *NativeConfig) ParseFile(file_name string) error {
	conf.MustLoad(file_name, n)
	return nil
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

func (t *RiskCtlTestConfig) Init() {
	t.FeeRiskctrlOpen = true
	t.BiasRiskctrlOpen = true
	t.WatermarkRiskctrlOpen = true
	t.PricesionRiskctrlOpen = true

}

// func TESTCONFIG_INIT(file_name string) error {
// 	test_config := TESTCONFIG()

// 	test_config.Init()

// 	return nil
// }

func NATIVE_CONFIG_INIT(file_name string) error {
	native_config := NATIVE_CONFIG()
	return native_config.ParseFile(file_name)
}

func NATIVE_CONFIG() *NativeConfig {
	if g_single_native_config == nil {
		lock.Lock()
		defer lock.Unlock()

		if g_single_native_config == nil {
			g_single_native_config = new(NativeConfig)
			fmt.Println("Init Single")
		} else {
			fmt.Println("Second Judge")
		}
	} else {
		// fmt.Println("Single already created!")
	}

	return g_single_native_config
}

func TESTCONFIG() *RiskCtlTestConfig {
	native_config := NATIVE_CONFIG()
	return &native_config.RiskTestConfig
}

func TestConf() {
	var c NativeConfig
	configFile := "client.yaml"
	conf.MustLoad(configFile, &c)
	fmt.Printf("%+v\n", c)
}
