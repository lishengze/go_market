package config

import (
	"flag"
	"fmt"
	"os"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/zrpc"
)

// Mysql:
//   Addr: "bcts:bcts@tcp(127.0.0.1:3306)/market"
//   max_open_conns: 16
//   max_idle_conns: 8
//   conn_max_life_time: 300
type MysqlConfig struct {
	Addr string
	// max_open_conns     int
	// max_idle_conns     int
	// conn_max_life_time int
}

type DingConfig struct {
	Secret string
	Token  string
}

type KafkaConfig struct {
	IP string `json:",optional"`
}

type CommConfig struct {
	KafkaConfig
	NetServerType string `json:",optional"`
	SerialType    string `json:",optional"`
}

type Config struct {
	zrpc.RpcServerConf

	Comm      CommConfig
	Nacos     NacosConfig
	LogConfig logx.LogConf
	Mysql     MysqlConfig
}

func (n *Config) ParseFile(file_name string) error {
	conf.MustLoad(file_name, n)
	return nil
}

// func NATIVE_CONFIG_INIT(file_name string) error {
// 	native_config := NATIVE_CONFIG()
// 	return native_config.ParseFile(file_name)
// }

// var g_single_native_config *Config
// var lock = &sync.Mutex{}

// func NATIVE_CONFIG() *Config {
// 	if g_single_native_config == nil {
// 		lock.Lock()
// 		defer lock.Unlock()

// 		if g_single_native_config == nil {
// 			g_single_native_config = new(Config)
// 			fmt.Println("Init Single")
// 		} else {
// 			fmt.Println("Second Judge")
// 		}
// 	} else {
// 		// fmt.Println("Single already created!")
// 	}

// 	return g_single_native_config
// }

// func TESTCONFIG() *RiskCtlTestConfig {
// 	native_config := NATIVE_CONFIG()
// 	return &native_config.RiskTestConfig
// }

func TestConf() {
	var c Config
	flag.Parse()

	fmt.Printf("Args: %+v \n", os.Args)
	env := os.Args[1]

	fmt.Printf("env: %+v \n", env)
	var configFile = flag.String("f", "etc/"+env+"/marketData.yaml", "the config file")

	conf.MustLoad(*configFile, &c)
	// fmt.Println(c.Nacos.IpAddr, ": ", c.Nacos.Port)
	fmt.Printf("%+v\n", c)
}
