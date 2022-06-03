package nacosAdapter

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest"
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
	err := client.publishConfig("currency", "parameter-management", currencyStr)
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
	err = client.publishConfig("hedging-platform", "parameter-management", hedgingStr)
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

	currencyContent, _ := client.getConfigContent("currency", "parameter-management")
	fmt.Println(1111111, currencyContent)
	var currencyMap map[string]interface{}
	contentBytes := []byte(currencyContent)
	err := json.Unmarshal(contentBytes, &currencyMap)
	if err != nil {
		fmt.Println(111222, err.Error())
	}
	fmt.Println(222211111, currencyMap["ChineseName"], currencyMap)

	hedgingContent, _ := client.getConfigContent("hedging-platform", "parameter-management")
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

	go client.listenConfig("currency", "parameter-management", ConfigChange)

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
	err := client.publishConfig("currency", "parameter-management", currencyStr)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(1111, err)

	time.Sleep(100 * time.Second)
}

type BctsConfig struct {
	rest.RestConf
	JwtAuth struct {
		AccessSecret string
		AccessExpire int64
	}
	Nacos Config
}

//common/nacosAdapter/nacos_test.go
// vs
//app/admin/cmd/api/etc/admin.yaml
//var configFile = flag.String("f", "etc/admin.yaml", "the config file")
var configFile = flag.String("f", "./../../app/admin/cmd/api/etc/admin.yaml", "the config file")

func TestClient_GetAll(t *testing.T) {
	flag.Parse()
	var c BctsConfig
	conf.MustLoad(*configFile, &c)

	nacosClient := NewClient(&c.Nacos)

	currencies, err := nacosClient.GetCurrency("")
	if err != nil {
		logx.Error(err)
	} else {
		for _, currency := range currencies {
			logx.Infof("%+v", currency)
			break
		}
	}

	logx.Info("GetHedging")
	hedgingList, err := nacosClient.GetHedging("")
	if err != nil {
		logx.Error(err)
	} else {
		//logx.Infof("%+v", hedgingList)
		for _, hedging := range hedgingList {
			logx.Infof("%+v", hedging)
			break
		}
	}

	logx.Info("GetSymbol")
	symbols, err := nacosClient.GetSymbol("")
	if err != nil {
		logx.Error(err)
	} else {
		//logx.Infof("%+v", symbols)
		for _, symbol := range symbols {
			logx.Infof("%+v", symbol)
			break
		}
	}
	tradeSwitch, err := nacosClient.GetTradeSwitch()
	if err != nil {
		logx.Error(err)
	} else {
		logx.Infof("%+v", tradeSwitch)
	}

	groupTradeParam, err := nacosClient.GetGroupTradeParam("WXBroker-VIP2")
	if err != nil {
		logx.Error(err)
	} else {
		//logx.Infof("%+v", groupTradeParam)
		for _, v := range groupTradeParam {
			logx.Infof("%+v", v)
			break
		}
	}
}

func TestClient_ListenAll(t *testing.T) {
	flag.Parse()
	var c BctsConfig
	conf.MustLoad(*configFile, &c)

	nacosClient := NewClient(&c.Nacos)

	nacosClient.ListenCurrency(func(namespace, group, dataId, data string) {
		logx.Info(namespace)
		logx.Info(group)
		logx.Info(dataId)
		logx.Info(data)
	})
	time.Sleep(time.Hour)
}

func TestClient_GetChanged(t *testing.T) {
	flag.Parse()
	var c BctsConfig
	conf.MustLoad(*configFile, &c)

	nacosClient := NewClient(&c.Nacos)
	for i := 0; i < 50; i++ {
		content, err := nacosClient.getConfigContent(Test_PARAMS, BCTS_GROUP)
		if err != nil {
			logx.Error(err)
		} else {
			logx.Infof("%+v", content)
		}

		time.Sleep(time.Second * 1)
	}
}
