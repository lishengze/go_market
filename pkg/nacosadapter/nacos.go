package nacosadapter

import (
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

type ServerConfig struct {
	IpAddr string
	Port   uint64
}

type ClientConfig struct {
	NamespaceId         string
	TimeoutMs           uint64
	NotLoadCacheAtStart bool
	LogDir              string
	CacheDir            string
	LogLevel            string
}

type Config struct {
	*ServerConfig
	*ClientConfig
}

type Client struct {
	iClient config_client.IConfigClient
}

func NewClient(c *Config) *Client {
	sc := []constant.ServerConfig{
		{
			IpAddr: c.IpAddr,
			Port:   c.Port,
		},
	}

	cc := constant.ClientConfig{
		NamespaceId:         c.NamespaceId,
		TimeoutMs:           c.TimeoutMs,
		NotLoadCacheAtStart: c.NotLoadCacheAtStart,
		LogDir:              c.LogDir,
		CacheDir:            c.CacheDir,
		LogLevel:            c.LogLevel,
	}

	client, err := clients.NewConfigClient(
		vo.NacosClientParam{
			ClientConfig:  &cc,
			ServerConfigs: sc,
		},
	)

	if err != nil {
		panic(err)
	}
	return &Client{
		iClient: client,
	}
}

func (c *Client) GetConfigContent(dataId string, group string) (string, error) {
	content, err := c.iClient.GetConfig(vo.ConfigParam{
		DataId: dataId,
		Group:  group,
	})
	if err != nil {
		return "", err
	}
	return content, err
}

func (c *Client) PublishConfig(dataId string, group string, content string) error {
	_, err := c.iClient.PublishConfig(vo.ConfigParam{
		DataId:  dataId,
		Group:   group,
		Content: content,
	})
	return err
}

func (c *Client) DeleteConfig(dataId string, group string) error {
	_, err := c.iClient.DeleteConfig(vo.ConfigParam{
		DataId: dataId,
		Group:  group,
	})
	return err
}

func (c *Client) ListenConfig(dataId string, group string, f func(namespace, group, dataId, data string)) error {
	err := c.iClient.ListenConfig(vo.ConfigParam{
		DataId:   dataId,
		Group:    group,
		OnChange: f,
	})
	return err
}
