package core

import (
	"exterior-interactor/app/mpu/rpc/internal/config"
	"exterior-interactor/app/mpu/rpc/mpupb"
	"exterior-interactor/pkg/exchangeapi/exchanges/binance"
	"exterior-interactor/pkg/exchangeapi/exchanges/ftx"
	"exterior-interactor/pkg/exchangeapi/exmodel"
	"exterior-interactor/pkg/exchangeapi/extools"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/gogo/protobuf/proto"
	nacosclients "github.com/nacos-group/nacos-sdk-go/clients"
	nacosconfigclient "github.com/nacos-group/nacos-sdk-go/clients/config_client"
	nacosconstant "github.com/nacos-group/nacos-sdk-go/common/constant"
	nacosconf "github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type (
	Mpu interface {
		//Sub(symbols ...exmodel.StdSymbol)
	}

	mpu struct {
		exchange string
		extools.SymbolManager
		extools.DepthManager
		extools.MarketTradeManager
		extools.KlineGenerator

		kafkaSyncProducer sarama.SyncProducer
		nacosClient       nacosconfigclient.IConfigClient
	}

	MpuConfig struct {
		Symbols []string
	}
)

func newNacosClient(c config.Config, onConfigUpdate func(c MpuConfig)) nacosconfigclient.IConfigClient {
	var serverConfigs []nacosconstant.ServerConfig
	for _, item := range c.NacosConf.Server {
		serverConfigs = append(serverConfigs, nacosconstant.ServerConfig{
			Scheme:      "",
			ContextPath: "",
			IpAddr:      item.Host,
			Port:        item.Port,
		})
	}

	clientConfig := &nacosconstant.ClientConfig{
		NamespaceId:         c.NacosConf.Client.NamespaceId,
		TimeoutMs:           c.NacosConf.Client.TimeoutMs,
		LogDir:              c.NacosConf.Client.LogDir,
		CacheDir:            c.NacosConf.Client.CacheDir,
		LogLevel:            c.NacosConf.Client.LogLevel,
		NotLoadCacheAtStart: true,
	}

	nacosClient, err := nacosclients.NewConfigClient(nacosconf.NacosClientParam{
		ClientConfig:  clientConfig,
		ServerConfigs: serverConfigs,
	})

	if err != nil {
		panic(fmt.Sprintf("create nacos client err:%v", err))
	}

	nacosConfig := nacosconf.ConfigParam{
		DataId:  c.NacosConf.DataId,
		Group:   c.NacosConf.Group,
		Content: "",
		DatumId: "",
		Type:    nacosconf.YAML,
		OnChange: func(namespace, group, dataId, data string) {
			var c MpuConfig
			err := conf.LoadConfigFromYamlBytes([]byte(data), &c)
			if err != nil {
				logx.Errorf("load config from nacos err:%v, data:%s", err, data)
				return
			}
			logx.Infof("detect config update, config:%v", c)
			onConfigUpdate(c)
		},
	}

	err = nacosClient.ListenConfig(nacosConfig)

	if err != nil {
		panic(fmt.Sprintf("listen nacos config err:%v", err))
	}

	return nacosClient
}

func newKafkaSyncProducer(c config.Config) sarama.SyncProducer {
	kafkaConfig := sarama.NewConfig()
	kafkaConfig.Producer.RequiredAcks = sarama.WaitForAll          // 发送完数据需要leader和follow都确认
	kafkaConfig.Producer.Partitioner = sarama.NewRandomPartitioner // 新选出一个partition
	kafkaConfig.Producer.Return.Successes = true                   // 成功交付的消息将在success channel返回

	client, err := sarama.NewSyncProducer([]string{c.KafkaConf.Address}, kafkaConfig)
	if err != nil {
		panic(fmt.Sprintf("create kafka client err:%v", err))
	}

	return client
}

func NewMpu(c config.Config) Mpu {
	m := &mpu{
		exchange:           c.Exchange,
		SymbolManager:      nil,
		DepthManager:       nil,
		MarketTradeManager: nil,
		KlineGenerator:     nil,
		nacosClient:        nil,
		kafkaSyncProducer:  newKafkaSyncProducer(c),
	}

	m.nacosClient = newNacosClient(c, m.onConfigUpdate) // 放在最后

	switch exmodel.Exchange(c.Exchange) {
	case exmodel.BINANCE:
		api := binance.NewNativeApiWithProxy(exmodel.EmptyAccountConfig, c.Proxy)
		m.SymbolManager = binance.NewSymbolManager(api)
		m.DepthManager = binance.NewDepthManager(m.SymbolManager, api)
		m.MarketTradeManager = binance.NewMarketTradeManager(m.SymbolManager, api)
		m.KlineGenerator = extools.NewKlineGenerator()
	case exmodel.FTX:
		api := ftx.NewNativeApiWithProxy(exmodel.EmptyAccountConfig, c.Proxy)
		m.SymbolManager = ftx.NewSymbolManager(api)
		m.DepthManager = ftx.NewDepthManager(m.SymbolManager, api)
		m.MarketTradeManager = ftx.NewMarketTradeManager(m.SymbolManager, api)
		m.KlineGenerator = extools.NewKlineGenerator()
	default:
		panic(fmt.Sprintf("mpu not support exchange:%s", c.Exchange))
	}

	// 启动时加载一次配置
	content, err := m.nacosClient.GetConfig(nacosconf.ConfigParam{
		DataId: c.NacosConf.DataId,
		Group:  c.NacosConf.Group,
		Type:   nacosconf.YAML,
	})
	if err != nil {
		panic(fmt.Sprintf("get nacos config err:%v", err))
	}

	var mpuConfig MpuConfig
	err = conf.LoadConfigFromYamlBytes([]byte(content), &mpuConfig)
	if err != nil {
		panic(fmt.Sprintf("convert nacos config err:%v", err))
	}
	logx.Infof("MPU load config: %v", c)
	m.onConfigUpdate(mpuConfig)

	m.run()

	return m
}

func (o *mpu) onConfigUpdate(c MpuConfig) {
	for _, s := range c.Symbols {
		stdSymbol := exmodel.StdSymbol(s)
		o.MarketTradeManager.Sub(stdSymbol)
		o.DepthManager.Sub(stdSymbol)
	}
}

func (o *mpu) run() {
	go o.dispatchKline()
	go o.dispatchDepth()
	go o.dispatchTrade()
}

func (o *mpu) dispatchDepth() {
	depthCh := o.DepthManager.OutputCh()
	for {
		select {
		case depth := <-depthCh:
			bytes, err := o.convertDepth(depth)
			if err != nil {
				continue
			}

			var (
				kafkaTopic = fmt.Sprintf("DEPTH.%s.%s", depth.Symbol.StdSymbol.String(), depth.Exchange.String())
				chLen      = len(depthCh)
				msgs       = []*sarama.ProducerMessage{{
					Topic: kafkaTopic,
					Value: sarama.ByteEncoder(bytes),
				}}
			)

			for i := 0; i < chLen; i++ {
				depth = <-depthCh
				bytes, err := o.convertDepth(depth)
				if err != nil {
					continue
				}
				msgs = append(msgs, &sarama.ProducerMessage{
					Topic: kafkaTopic,
					Value: sarama.ByteEncoder(bytes),
				})
			}

			//start := time.Now()
			err = o.kafkaSyncProducer.SendMessages(msgs)
			//fmt.Printf("send %d messages,cost:%v \n", len(msgs), time.Now().Sub(start))

			if err != nil {
				logx.Errorf("kafkaConn write depth err:%s, data:%+v", err, *depth)
			}
		}
	}
}

func (o *mpu) dispatchTrade() {
	tradeCh := o.MarketTradeManager.OutputCh()
	for {
		select {
		case trade := <-tradeCh:
			o.KlineGenerator.InputMarketTrade(trade) // 首先推送到到 kline

			bytes, err := o.convertStreamMarketTrade(trade)
			if err != nil {
				continue
			}

			var (
				kafkaTopic = fmt.Sprintf("TRADE.%s.%s", trade.Symbol.StdSymbol.String(), trade.Exchange.String())
				chLen      = len(tradeCh)
				msgs       = []*sarama.ProducerMessage{{
					Topic: kafkaTopic,
					Value: sarama.ByteEncoder(bytes),
				}}
			)

			for i := 0; i < chLen; i++ {
				trade = <-tradeCh
				o.KlineGenerator.InputMarketTrade(trade) // 首先推送到到 kline
				bytes, err := o.convertStreamMarketTrade(trade)
				if err != nil {
					continue
				}
				msgs = append(msgs, &sarama.ProducerMessage{
					Topic: kafkaTopic,
					Value: sarama.ByteEncoder(bytes),
				})
			}


			//start := time.Now()
			err = o.kafkaSyncProducer.SendMessages(msgs)
			//fmt.Printf("send %d messages,cost:%v \n", len(msgs), time.Now().Sub(start))

			if err != nil {
				logx.Errorf("kafkaConn write trade err:%s, data:%+v", err, *trade)
			}
		}
	}
}

func (o *mpu) dispatchKline() {
	klineCh := o.KlineGenerator.OutputCh()
	for {
		select {
		case kline := <-klineCh:
			bytes, err := o.convertKline(kline)
			if err != nil {
				continue
			}

			var (
				kafkaTopic = fmt.Sprintf("KLINE.%s.%s", kline.Symbol.StdSymbol.String(), kline.Exchange.String())
				chLen      = len(klineCh)
				msgs       = []*sarama.ProducerMessage{{
					Topic: kafkaTopic,
					Value: sarama.ByteEncoder(bytes),
				}}
			)

			for i := 0; i < chLen; i++ {
				kline = <-klineCh
				bytes, err := o.convertKline(kline)
				if err != nil {
					continue
				}
				msgs = append(msgs, &sarama.ProducerMessage{
					Topic: kafkaTopic,
					Value: sarama.ByteEncoder(bytes),
				})
			}

			//start := time.Now()
			err = o.kafkaSyncProducer.SendMessages(msgs)
			//fmt.Printf("send %d messages,cost:%v \n", len(msgs), time.Now().Sub(start))

			if err != nil {
				logx.Errorf("kafkaConn write kline err:%s, data:%+v", err, *kline)
			}
		}
	}
}

func (o *mpu) convertStreamMarketTrade(in *exmodel.StreamMarketTrade) ([]byte, error) {
	data := &mpupb.Trade{
		Timestamp: timestamppb.New(in.Time),
		Exchange:  in.Exchange.String(),
		Symbol:    in.Symbol.StdSymbol.String(),
		Price:     in.Price,
		Volume:    in.Volume,
	}

	msg, err := proto.Marshal(data)
	if err != nil {
		logx.Errorf("convertStreamMarketTrade err:%v, data:%v ", err, data)
	}
	return msg, err
}

func (o *mpu) convertKline(in *exmodel.Kline) ([]byte, error) {
	data := &mpupb.Kline{
		Timestamp:  timestamppb.New(in.Time),
		Exchange:   in.Exchange.String(),
		Symbol:     in.Symbol.StdSymbol.String(),
		Open:       in.Open,
		High:       in.High,
		Low:        in.Low,
		Close:      in.Close,
		Volume:     in.Volume,
		Value:      in.Value,
		Resolution: in.Resolution,
	}

	msg, err := proto.Marshal(data)
	if err != nil {
		logx.Errorf("convertKline err:%v, data:%v ", err, data)
	}
	return msg, err
}

func (o *mpu) convertDepth(in *exmodel.StreamDepth) ([]byte, error) {
	var (
		asks, bids []*mpupb.PriceVolume
	)

	for _, r := range in.Asks {
		asks = append(asks, &mpupb.PriceVolume{
			Price:  r[0],
			Volume: r[1],
		})
	}

	for _, r := range in.Bids {
		bids = append(bids, &mpupb.PriceVolume{
			Price:  r[0],
			Volume: r[1],
		})
	}

	data := &mpupb.Depth{
		Timestamp: timestamppb.New(in.Time),
		Exchange:  in.Exchange.String(),
		Symbol:    in.Symbol.StdSymbol.String(),
		Asks:      asks,
		Bids:      bids,
	}

	msg, err := proto.Marshal(data)
	if err != nil {
		logx.Errorf("convertDepth err:%v, data:%v ", err, data)
	}
	return msg, err
}
