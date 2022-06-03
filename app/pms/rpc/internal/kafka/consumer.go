package kafka

import (
	"context"
	"fmt"
	"log"
	"market_server/app/pms/rpc/internal/kafka/logic"
	"market_server/app/pms/rpc/internal/kafka/types"
	"market_server/app/pms/rpc/internal/svc"
	"os"
	"strings"

	"github.com/Shopify/sarama"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/threading"
	"google.golang.org/protobuf/proto"
)

type Router struct {
	Handler func(ctx context.Context, svcCtx *svc.ServiceContext) logic.LogicHandle
}

type Kafka struct {
	Topics         []string `toml:"topic"`
	Broker         string   `toml:"broker"`
	Partition      int32    `toml:"partition"`
	Replication    int16    `toml:"replication"`
	Group          string   `toml:"group"`
	Version        string   `toml:"version"`
	Routers        Router
	serviceContext *svc.ServiceContext
}

func NewKafka(service *svc.ServiceContext) *Kafka {
	return &Kafka{
		Version:        "2.8.0",
		Group:          "1",
		Broker:         "152.32.254.76:9117",
		Topics:         []string{"TRADEx-ADA_USDT", "TRADEx-AR_USDT", "TRADEx-ATLAS_USD", "TRADEx-ATOM_USDT", "TRADEx-AUDIO_USD", "TRADEx-AUDIO_USDT", "TRADEx-AVAX_USD", "TRADEx-AVAX_USDT", "TRADEx-AXS_USD", "TRADEx-AXS_USDT", "TRADEx-BNB_USD", "TRADEx-BNB_USDT", "TRADEx-BSV_USDT", "TRADEx-BTC_USD", "TRADEx-BTC_USDT", "TRADEx-BTT_USD", "TRADEx-BTT_USDT", "TRADEx-C98_USD", "TRADEx-C98_USDT", "TRADEx-CTK_USDT", "TRADEx-DAI_USD", "TRADEx-DAI_USDT", "TRADEx-DOT_USD", "TRADEx-DOT_USDT", "TRADEx-EOS_USDT", "TRADEx-ETH_USD", "TRADEx-ETH_USDT", "TRADEx-FSN_USDT", "TRADEx-FTT_USD", "TRADEx-FTT_USDT", "TRADEx-GLM_USDT", "TRADEx-GNX_USDT", "TRADEx-HT_USD", "TRADEx-HT_USDT", "TRADEx-ICP_USDT", "TRADEx-KSM_USDT", "TRADEx-LINK_USD", "TRADEx-LINK_USDT", "TRADEx-LTC_USD", "TRADEx-LTC_USDT", "TRADEx-MANA_USD", "TRADEx-MANA_USDT", "TRADEx-MATIC_USD", "TRADEx-MATIC_USDT", "TRADEx-MNGO_USD", "TRADEx-NMR_USDT", "TRADEx-OKB_USD", "TRADEx-PAX_USDT", "TRADEx-PHA_USDT", "TRADEx-QTUM_USDT", "TRADEx-RAY_USD", "TRADEx-RAY_USDT", "TRADEx-SLP_USD", "TRADEx-SLP_USDT", "TRADEx-SOL_USD", "TRADEx-SOL_USDT", "TRADEx-SRM_USD", "TRADEx-SRM_USDT", "TRADEx-SXP_USD", "TRADEx-SXP_USDT", "TRADEx-TRX_USD", "TRADEx-TRX_USDT", "TRADEx-USDT_USD", "TRADEx-VET_USDT", "TRADEx-XRP_USD", "TRADEx-XRP_USDT", "TRADEx-USDK_USDT", "TRADEx-USDC_USDT", "TRADEx-TUSD_USDT", "TRADEx-PST_USDT", "TRADEx-BUSD_USDT", "TRADEx-BTC_HUSD"},
		serviceContext: service,
	}
}

func (k *Kafka) AddRouter(router Router) {
	k.Routers = router
}

func (k *Kafka) RunKafka() {
	for _, topic := range k.Topics {
		singleTopic := func() {
			sarama.Logger = log.New(os.Stdout, fmt.Sprintf("[%s]", topic), log.LstdFlags)
			consumer, err := sarama.NewConsumer(strings.Split(k.Broker, ","), nil)
			if err != nil {
				logx.Error("fail to start consumer, err:%v", err)
				return
			}
			pc, err := consumer.ConsumePartition(topic, 0, sarama.OffsetNewest)
			if err != nil {
				logx.Error("failed to start consumer for partition %d,err:%v", 0, err)
				return
			}
			defer func() {
				pc.AsyncClose()
			}()
			priceValue := types.TradeData{}
			for msg := range pc.Messages() {
				err = proto.Unmarshal([]byte(msg.Value), &priceValue)
				if err != nil {
					logx.Error("ConsumeClaim error(%v)", err)
					return
				}
				err = k.Routers.Handler(context.Background(), k.serviceContext).Handle(&priceValue)
				if err != nil {
					logx.Error(err)
					return
				}
			}
		}
		threading.GoSafe(singleTopic)
	}
}

func (k *Kafka) SingleTopic(topic string) {
	sarama.Logger = log.New(os.Stdout, fmt.Sprintf("[%s]", topic), log.LstdFlags)
	consumer, err := sarama.NewConsumer(strings.Split(k.Broker, ","), nil)
	if err != nil {
		logx.Error("fail to start consumer, err:%v", err)
		return
	}

	pc, err := consumer.ConsumePartition(topic, 0, sarama.OffsetNewest)
	if err != nil {
		logx.Error("failed to start consumer for partition %d,err:%v", 0, err)
		return
	}
	defer func() {
		pc.AsyncClose()
	}()
	// 异步从每个分区消费信息
	priceValue := types.TradeData{}
	for msg := range pc.Messages() {
		err = proto.Unmarshal([]byte(msg.Value), &priceValue)
		if err != nil {
			logx.Error("ConsumeClaim error(%v)", err)
			return
		}
		err = k.Routers.Handler(context.Background(), k.serviceContext).Handle(&priceValue)
		if err != nil {
			logx.Error(err)
			return
		}
	}
}
