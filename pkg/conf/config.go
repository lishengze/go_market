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

type AggregateConfigAtom struct {
	AggregateFreq time.Duration
	PublishLevel  int
	IsPublish     bool
}
type AggregateConfig struct {
	DepthAggregatorConfigMap map[string]AggregateConfigAtom
}

func (a *AggregateConfig) String() string {
	result := ""

	for symbol, aggregate_atom := range a.DepthAggregatorConfigMap {
		result += symbol + " " + fmt.Sprintf(", %+v\n", aggregate_atom)
	}

	return result
}

type HedgeConfig struct {
	FeeKind  int
	FeeValue float64
}

type RiskCtrlConfig struct {
	HedgeConfigMap map[string]*HedgeConfig

	PricePrecison  uint32
	VolumePrecison uint32

	PriceBiasValue float64
	PriceBiasKind  int

	VolumeBiasValue float64
	VolumeBiasKind  int

	PriceMinumChange float64
}

func (r *RiskCtrlConfig) String() string {

	hedge_info := ""
	for symbol, hedge_config := range r.HedgeConfigMap {
		hedge_info += fmt.Sprintf("%s: %+v\n", symbol, *hedge_config)
	}

	return fmt.Sprintf("HedgeConfigMap: %s \nPricePrecison: %v, VolumePrecison: %v \nPriceBiasValue: %v, PriceBiasKind: %v \nVolumeBiasValue: %v, VolumeBiasKind: %v\nPriceMinumChange:%v \n",
		hedge_info,
		r.PricePrecison, r.VolumePrecison,
		r.PriceBiasValue, r.PriceBiasKind,
		r.VolumeBiasValue, r.VolumeBiasKind,
		r.PriceMinumChange)

	// return fmt.Sprintf("HedgeConfigMap: %s\nPricePrecison: %v",
	// 	hedge_info, r.PricePrecison)

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
	fmt.Println(c.Nacos.IpAddr, ": ", c.Nacos.Port)
	fmt.Printf("%+v\n", c)
}
