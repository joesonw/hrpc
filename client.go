package hrpc

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"google.golang.org/protobuf/proto"

	"github.com/joesonw/hrpc/status"
)

type Client struct {
	baseURL string

	httpClient *http.Client
	codec      Codec
	middleware ClientMiddleware
}

func NewClient(baseURL string, options ...ClientOption) *Client {
	o := &clientOptions{
		codec:      JSONCodec{},
		httpClient: http.DefaultClient,
	}

	for _, apply := range options {
		apply(o)
	}

	return &Client{
		baseURL: strings.TrimRight(baseURL, "/"),

		httpClient: o.httpClient,
		codec:      o.codec,
		middleware: o.middleware,
	}
}

func (c *Client) Invoke(ctx context.Context, method string, req, res proto.Message, opts ...CallOption) error {
	u := c.baseURL + method
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, u, nil)
	if err != nil {
		return fmt.Errorf("unable to create http.Request: %w", err)
	}

	b, err := c.codec.Marshal(req)
	if err != nil {
		return fmt.Errorf("unable to marshal request: %w", err)
	}
	httpReq.Body = io.NopCloser(bytes.NewBuffer(b))

	var httpRes *http.Response
	if c.middleware != nil {
		httpRes, err = c.middleware(httpReq, c.httpClient, method, req, res, func(r *http.Request, client *http.Client, req, reply proto.Message, opts ...CallOption) (*http.Response, error) {
			return client.Do(httpReq)
		}, opts...)
	} else {
		httpRes, err = c.httpClient.Do(httpReq)
	}

	if err != nil {
		return fmt.Errorf("unalbe to do request: %w", err)
	}
	b, err = io.ReadAll(httpRes.Body)
	if err != nil {
		return fmt.Errorf("unable to read response: %w", err)
	}
	_ = httpRes.Body.Close()

	if httpRes.StatusCode >= 400 {
		return status.Wrap(httpRes.StatusCode, errors.New(string(b)))
	}

	if err := c.codec.Unmarshal(b, res); err != nil {
		return fmt.Errorf("unable to unmarshal response: %w", err)
	}

	return nil
}
