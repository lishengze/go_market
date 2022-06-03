package http

import (
	"context"
	"crypto/tls"
	"fmt"
	nethttp "net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"
)

type mockRoundTripper struct{}

func (rt *mockRoundTripper) RoundTrip(req *nethttp.Request) (resp *nethttp.Response, err error) {
	return
}

func TestWithTransport(t *testing.T) {
	ov := &mockRoundTripper{}
	o := WithTransport(ov)
	co := &clientOptions{}
	o(co)
	if !reflect.DeepEqual(co.transport, ov) {
		t.Errorf("expected transport to be %v, got %v", ov, co.transport)
	}
}

func TestWithTimeout(t *testing.T) {
	ov := 1 * time.Second
	o := WithTimeout(ov)
	co := &clientOptions{}
	o(co)
	if !reflect.DeepEqual(co.timeout, ov) {
		t.Errorf("expected timeout to be %v, got %v", ov, co.timeout)
	}
}

func TestWithTLSConfig(t *testing.T) {
	ov := &tls.Config{}
	o := WithTLSConfig(ov)
	co := &clientOptions{}
	o(co)
	if !reflect.DeepEqual(co.tlsConf, ov) {
		t.Errorf("expected tls config to be %v, got %v", ov, co.tlsConf)
	}
}

func TestWithUserAgent(t *testing.T) {
	ov := "go-zero"
	o := WithUserAgent(ov)
	co := &clientOptions{}
	o(co)
	if !reflect.DeepEqual(co.userAgent, ov) {
		t.Errorf("expected user agent to be %v, got %v", ov, co.userAgent)
	}
}

func TestWithEndpoint(t *testing.T) {
	ov := "some-endpoint"
	o := WithEndpoint(ov)
	co := &clientOptions{}
	o(co)
	if !reflect.DeepEqual(co.endpoint, ov) {
		t.Errorf("expected endpoint to be %v, got %v", ov, co.endpoint)
	}
}

func TestWithRequetEncoder(t *testing.T) {
	o := &clientOptions{}
	v := func(r *nethttp.Request) *nethttp.Request {
		return nil
	}
	WithRequetEncoder(v)(o)
	if o.encoder == nil {
		t.Errorf("expected encoder to be not nil")
	}
}

func TestWithInterceptor(t *testing.T) {
	o := &clientOptions{}
	v := []Interceptor{
		func(r *nethttp.Request) (*nethttp.Request, ResponseHandler) {
			return r, func(resp *nethttp.Response, err error) {}
		},
	}

	WithInterceptor(v...)(o)
	r := new(nethttp.Request)
	for _, opt := range o.encoder {
		r = opt(r)
	}

	if o.encoder == nil || o.responseHandlers == nil {
		t.Errorf("expected interceptor to be %+v, got encoder %+v and responseHandler %+v", v, o.encoder, o.responseHandlers)
	}
}

func TestNewClient(t *testing.T) {
	ts := httptest.NewServer(nethttp.HandlerFunc(func(w nethttp.ResponseWriter, r *nethttp.Request) {
		fmt.Fprintln(w, "Hello, client")
	}))
	defer ts.Close()
	cli, err := NewClient(context.Background(),
		ts.URL,
		WithTimeout(60*time.Second),
		WithTransport(&nethttp.Transport{
			MaxIdleConnsPerHost: 50,
			MaxIdleConns:        50,
			TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
		}),
	)
	if err != nil {
		t.Error(err)
	}

	resp, err := cli.Do(nethttp.MethodPost, nil)
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(resp.StatusCode, nethttp.StatusOK) {
		t.Errorf("expected status code to be %v,got status code %v", nethttp.StatusOK, resp.StatusCode)
	}
}
