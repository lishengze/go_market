package http

import (
	"context"
	"crypto/tls"
	"net/http"
	"time"

	"github.com/zeromicro/go-zero/rest/httpc"
)

type (
	ClientOption func(*clientOptions)

	// Client is an HTTP transport client.
	clientOptions struct {
		ctx              context.Context
		tlsConf          *tls.Config
		timeout          time.Duration
		endpoint         string
		userAgent        string
		encoder          []httpc.Option
		transport        http.RoundTripper
		responseHandlers []ResponseHandler
	}
)

// WithTransport with client transport.
func WithTransport(trans http.RoundTripper) ClientOption {
	return func(o *clientOptions) {
		o.transport = trans
	}
}

// WithTimeout with client request timeout.
func WithTimeout(d time.Duration) ClientOption {
	return func(o *clientOptions) {
		o.timeout = d
	}
}

// WithUserAgent with client user agent.
func WithUserAgent(ua string) ClientOption {
	return func(o *clientOptions) {
		o.userAgent = ua
	}
}

// WithRequetEncoder with client request encoder.
func WithRequetEncoder(fn ...httpc.Option) ClientOption {
	return func(o *clientOptions) {
		o.encoder = fn
	}
}

// WithEndpoint with client addr.
func WithEndpoint(endpoint string) ClientOption {
	return func(o *clientOptions) {
		o.endpoint = endpoint
	}
}

// WithTLSConfig with tls config.
func WithTLSConfig(c *tls.Config) ClientOption {
	return func(o *clientOptions) {
		o.tlsConf = c
	}
}

// WithInterceptor with http Interceptor.
func WithInterceptor(interceptors ...Interceptor) ClientOption {
	return func(o *clientOptions) {
		var interceptorEncoders []httpc.Option
		for _, interceptor := range interceptors {
			interceptorEncoders = append(interceptorEncoders, func(r *http.Request) *http.Request {
				r, h := interceptor(r)
				o.responseHandlers = append(o.responseHandlers, h)
				return r
			})
		}

		o.encoder = append(o.encoder, interceptorEncoders...)
	}
}

type Client struct {
	opts     clientOptions
	endpoint string
	sv       httpc.Service
}

func NewClient(ctx context.Context, endpoint string, opts ...ClientOption) (Client, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	options := clientOptions{
		ctx:       ctx,
		timeout:   2000 * time.Millisecond,
		transport: http.DefaultTransport,
	}
	for _, o := range opts {
		o(&options)
	}
	if options.tlsConf != nil {
		if tr, ok := options.transport.(*http.Transport); ok {
			tr.TLSClientConfig = options.tlsConf
		}
	}

	sv := httpc.NewServiceWithClient(endpoint, &http.Client{
		Timeout:   options.timeout,
		Transport: options.transport,
	}, options.encoder...)

	return Client{
		opts:     options,
		endpoint: endpoint,
		sv:       sv,
	}, nil
}

// Do sends an HTTP request with the given arguments and returns an HTTP response.
func (c Client) Do(method string, data interface{}) (*http.Response, error) {
	return c.do(c.opts.ctx, method, data)
}

// DoRequest sends an HTTP request to the service.
func (c Client) DoRequest(r *http.Request) (*http.Response, error) {
	return c.doRequest(c.opts.ctx, r)
}

func (c Client) do(ctx context.Context, method string, data interface{}) (*http.Response, error) {
	resp, err := c.sv.Do(ctx, method, c.endpoint, data)
	for i := len(c.opts.responseHandlers) - 1; i >= 0; i-- {
		c.opts.responseHandlers[i](resp, err)
	}
	return resp, err
}

func (c Client) doRequest(ctx context.Context, r *http.Request) (*http.Response, error) {
	resp, err := c.sv.DoRequest(r)
	for i := len(c.opts.responseHandlers) - 1; i >= 0; i-- {
		c.opts.responseHandlers[i](resp, err)
	}
	return resp, err
}
