package httptools

import (
	"crypto/tls"
	"net/http"
	"net/url"
	"time"
)

type HttpClient struct {
	*http.Client
	Proxy string
}

func NewHttpClient() *HttpClient {
	tr := &http.Transport{
		MaxIdleConnsPerHost: 200,
		MaxIdleConns:        200,
		TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
	}

	return &HttpClient{
		Client: &http.Client{
			Transport:     tr,
			CheckRedirect: nil,
			Jar:           nil,
			Timeout:       time.Second * 15,
		},
		Proxy: "",
	}
}

func NewHttpClientWithProxy(proxy string) (*HttpClient, error) {
	tr := &http.Transport{
		MaxIdleConnsPerHost: 200,
		MaxIdleConns:        200,
		TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
	}

	if proxy != "" {
		uProxy, err := url.Parse(proxy)
		if err != nil {
			return nil, err
		}
		tr.Proxy = http.ProxyURL(uProxy)
	}

	return &HttpClient{
		Client: &http.Client{
			Transport:     tr,
			CheckRedirect: nil,
			Jar:           nil,
			Timeout:       time.Second * 15,
		},
		Proxy: proxy,
	}, nil
}
