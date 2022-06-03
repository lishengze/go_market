package config

import (
	"fmt"
	"sync"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/zrpc"
)

var g_single_native_config *Config
var lock = &sync.Mutex{}

// Mysql:
//   Addr: "bcts:bcts@tcp(127.0.0.1:3306)/market"
//   max_open_conns: 16
//   max_idle_conns: 8
//   conn_max_life_time: 300
type MysqlConfig struct {
	Addr               string
	max_open_conns     int
	max_idle_conns     int
	conn_max_life_time int
}

type Config struct {
	zrpc.RpcServerConf

	IP            string
	NetServerType string
	SerialType    string

	Nacos     NacosConfig
	LogConfig logx.LogConf
	Mysql     MysqlConfig
}

func (n *Config) ParseFile(file_name string) error {
	conf.MustLoad(file_name, n)
	return nil
}

func NATIVE_CONFIG_INIT(file_name string) error {
	native_config := NATIVE_CONFIG()
	return native_config.ParseFile(file_name)
}

func NATIVE_CONFIG() *Config {
	if g_single_native_config == nil {
		lock.Lock()
		defer lock.Unlock()

		if g_single_native_config == nil {
			g_single_native_config = new(Config)
			fmt.Println("Init Single")
		} else {
			fmt.Println("Second Judge")
		}
	} else {
		// fmt.Println("Single already created!")
	}

	return g_single_native_config
}

// func TESTCONFIG() *RiskCtlTestConfig {
// 	native_config := NATIVE_CONFIG()
// 	return &native_config.RiskTestConfig
// }

func TestConf() {
	var c Config
	configFile := "client.yaml"
	conf.MustLoad(configFile, &c)
	fmt.Println(c.Nacos.IpAddr, ": ", c.Nacos.Port)
	fmt.Printf("%+v\n", c)
}
