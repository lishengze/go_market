package nacosadapter

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

func TestClient_PublishConfig(t *testing.T) {
	serverConfig := ServerConfig{
		IpAddr: "127.0.0.1",
		Port:   8848,
	}
	clientConfig := ClientConfig{
		NamespaceId:         "bcts-test", //namespace id
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "../../../logs/nacos/log",
		CacheDir:            "../../../logs/nacos/cache",
		LogLevel:            "debug",
	}
	config := Config{
		ServerConfig: &serverConfig,
		ClientConfig: &clientConfig,
	}
	client := NewClient(&config)

	currencyMap := map[string]interface{}{
		"currencyId":     "USDT",
		"ChineseName":    "泰达币",
		"EnglishName":    "Tether USD",
		"kind":           "公链数字货币",
		"minUnit":        1,
		"switch":         false,
		"minWithdraw":    10,
		"maxWithdraw":    1000,
		"maxDayWithdraw": 10000,
		"threshold":      100000,
		"feeKind":        1,
		"fee":            2,
		"minFee":         1,
	}
	currencyBytes, _ := json.Marshal(currencyMap)
	currencyStr := string(currencyBytes)
	err := client.PublishConfig("currency", "parameter-management", currencyStr)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(1111, err)

	hedgingMap := map[string]interface{}{
		"platformId":     "123",
		"instrument":     "BTC_USDT",
		"minUnit":        1,
		"minChangePrice": 0.1,
		"maxMargin":      5,
		"maxOrder":       100,
		"buyPriceLimit":  10000,
		"sellPriceLimit": 5000,
		"maxDayWithdraw": 10000,
		"maxMatchLevel":  5,
	}
	hedgingBytes, _ := json.Marshal(hedgingMap)
	hedgingStr := string(hedgingBytes)
	err = client.PublishConfig("hedging-platform", "parameter-management", hedgingStr)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(22222, err)
}

func TestClient_GetConfigContent(t *testing.T) {
	serverConfig := ServerConfig{
		IpAddr: "127.0.0.1",
		Port:   8848,
	}
	clientConfig := ClientConfig{
		NamespaceId:         "bcts-test", //namespace id
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "../../../logs/nacos/log",
		CacheDir:            "../../../logs/nacos/cache",
		LogLevel:            "debug",
	}
	config := Config{
		ServerConfig: &serverConfig,
		ClientConfig: &clientConfig,
	}
	client := NewClient(&config)

	currencyContent, _ := client.GetConfigContent("currency", "parameter-management")
	fmt.Println(1111111, currencyContent)
	var currencyMap map[string]interface{}
	contentBytes := []byte(currencyContent)
	err := json.Unmarshal(contentBytes, &currencyMap)
	if err != nil {
		fmt.Println(111222, err.Error())
	}
	fmt.Println(222211111, currencyMap["ChineseName"], currencyMap)

	hedgingContent, _ := client.GetConfigContent("hedging-platform", "parameter-management")
	fmt.Println(3333333, hedgingContent)
	var hedgingMap map[string]interface{}
	hedgingBytes := []byte(hedgingContent)
	err = json.Unmarshal(hedgingBytes, &hedgingMap)
	if err != nil {
		fmt.Println(33331111, err.Error())
	}
	fmt.Println(33333332222, hedgingMap["instrument"], hedgingMap)

}

func ConfigChange(namespace, group, dataId, data string) {
	fmt.Println("66666666,config changed,,,, group:" + group + ", dataId:" + dataId + ", content:" + data + ", namespace:" + namespace)
}

func TestClient_ListenConfig(t *testing.T) {
	serverConfig := ServerConfig{
		IpAddr: "127.0.0.1",
		Port:   8848,
	}
	clientConfig := ClientConfig{
		NamespaceId:         "bcts-test", //namespace id
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "../../../logs/nacos/log",
		CacheDir:            "../../../logs/nacos/cache",
		LogLevel:            "debug",
	}
	config := Config{
		ServerConfig: &serverConfig,
		ClientConfig: &clientConfig,
	}
	client := NewClient(&config)

	go client.ListenConfig("currency", "parameter-management", ConfigChange)

	currencyMap := map[string]interface{}{
		"currencyId":     "USDT2",
		"ChineseName":    "泰达币2",
		"EnglishName":    "Tether USD2",
		"kind":           "公链数字货币2",
		"minUnit":        1,
		"switch":         false,
		"minWithdraw":    10,
		"maxWithdraw":    1000,
		"maxDayWithdraw": 10000,
		"threshold":      100000,
		"feeKind":        1,
		"fee":            2,
		"minFee":         1,
	}
	currencyBytes, _ := json.Marshal(currencyMap)
	currencyStr := string(currencyBytes)
	err := client.PublishConfig("currency", "parameter-management", currencyStr)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(1111, err)

	time.Sleep(100 * time.Second)
}
