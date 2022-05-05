package main

import (
	"exterior-interactor/pkg/nacosadapter"
	"fmt"
	"github.com/zeromicro/go-zero/core/conf"
	"log"
	"time"
)

type C struct {
	Symbols []string
}

func main()  {
	serverConfig := nacosadapter.ServerConfig{
		IpAddr: "36.255.220.139",
		Port:   8848,
	}
	clientConfig := nacosadapter.ClientConfig{
		NamespaceId:         "exterior-interactor", //namespace id
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "./log",
		CacheDir:            "./cache",
		LogLevel:            "debug",
	}
	config := nacosadapter.Config{
		ServerConfig: &serverConfig,
		ClientConfig: &clientConfig,
	}
	client := nacosadapter.NewClient(&config)

	currencyContent, err := client.GetConfigContent("MPU", "MPU_GROUP")
	if err!=nil{
		log.Fatalln(err)
	}
	fmt.Println(1111111, currencyContent)
	c:=C{}
	err=conf.LoadConfigFromYamlBytes([]byte(currencyContent),&c)
	if err!=nil{
		log.Fatalln(err)
	}
	fmt.Println("******",c)



	go client.ListenConfig("MPU", "MPU_GROUP", ConfigChange)

	time.Sleep(100 * time.Second)
}


func ConfigChange(namespace, group, dataId, data string) {
	fmt.Println("66666666,config changed,,,, group:" + group + ", dataId:" + dataId + ", content:" + data + ", namespace:" + namespace)
}
