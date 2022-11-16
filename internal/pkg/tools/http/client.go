package http

import (
	"bytes"
	"context"
	"crypto/tls"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/tools/errors"
	"github.com/rs/xid"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

type Client interface {
	PerformRequest(ctx context.Context, req Request, res Response) error
}

type Request interface {
	URL() string
	Method() string
}

type RequestWithHeaders interface {
	Request
	Headers() http.Header
}

type RequestWithBody interface {
	Request
	Body() ([]byte, error)
}

type RequestWithRequestID interface {
	Request
	RequestID() string
}

type Response interface {
	ReadFrom(*http.Response) error
}

type Config struct {
	ConnectTimeout   time.Duration `mapstructure:"connect_timeout"`
	ReadWriteTimeout time.Duration `mapstructure:"read_write_timeout"`

	KeepAlive time.Duration `mapstructure:"keep_alive"`

	MaxIdleConns        int `mapstructure:"max_idle_conns"`
	MaxIdleConnsPerHost int `mapstructure:"max_idle_conns_per_host"`

	InsecureSkipVerify bool `mapstructure:"insecure_skip_verify"`
}

const (
	headerRequestID = "X-Request-ID"
)

type simpleHTTPClient struct {
	Client http.Client
}

func NewSimpleHTTPClient(
	cfg Config,
) (Client, error) {
	dialer := net.Dialer{
		Timeout: cfg.ConnectTimeout,
	}

	tlsConfig := tls.Config{
		MinVersion: tls.VersionTLS12,
	}

	transport := http.Transport{
		DialContext:           dialer.DialContext,
		MaxIdleConns:          cfg.MaxIdleConns,
		MaxIdleConnsPerHost:   cfg.MaxIdleConnsPerHost,
		IdleConnTimeout:       cfg.KeepAlive,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		TLSClientConfig:       &tlsConfig,
	}

	if cfg.InsecureSkipVerify {
		tlsConfig.InsecureSkipVerify = true
	}

	httpClient := http.Client{Transport: &transport, Timeout: cfg.ReadWriteTimeout}

	return &simpleHTTPClient{
		Client: httpClient,
	}, nil
}

func (c *simpleHTTPClient) PerformRequest(ctx context.Context, req Request, res Response) (err error) {
	httpReq, err := buildHTTPRequest(ctx, req)
	if err != nil {
		return errors.Wrap(err, "failed to convert request model to http request")
	}

	if httpReq.Body != nil {
		var reqBody []byte
		reqBody, err = ioutil.ReadAll(httpReq.Body)
		if err != nil {
			return errors.Wrap(err, "failed to read request body")
		}
		httpReq.Body.Close()
		httpReq.Body = ioutil.NopCloser(bytes.NewReader(reqBody))
	}

	httpResp, err := c.Client.Do(httpReq)
	if err != nil {
		return errors.Wrap(err, "failed to do request")
	}

	respBody, err := ioutil.ReadAll(httpResp.Body)
	if err != nil {
		return errors.Wrap(err, "failed to read response body")
	}
	httpResp.Body.Close()

	httpResp.Body = ioutil.NopCloser(bytes.NewReader(respBody))

	defer func() { httpResp.Body.Close() }()

	return res.ReadFrom(httpResp)
}

func buildHTTPRequest(ctx context.Context, req Request) (*http.Request, error) {
	var reqBody io.Reader

	if reqWithBody, ok := req.(RequestWithBody); ok {
		reqBodyRaw, err := reqWithBody.Body()
		if err != nil {
			return nil, errors.Wrap(err, "failed to fetch request body")
		}

		reqBody = bytes.NewReader(reqBodyRaw)
	}

	httpReq, err := http.NewRequestWithContext(ctx, req.Method(), req.URL(), reqBody)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create request with context")
	}

	if reqWithRequestID, ok := req.(RequestWithRequestID); ok {
		reqID := reqWithRequestID.RequestID()
		if reqID == "" {
			reqID = xid.New().String()
		}
		httpReq.Header.Set(headerRequestID, reqID)
	}

	if reqWithHeaders, ok := req.(RequestWithHeaders); ok {
		reqHeaders := reqWithHeaders.Headers()
		for k := range reqHeaders {
			httpReq.Header.Set(k, reqHeaders.Get(k))
		}
	}

	return httpReq, nil
}
