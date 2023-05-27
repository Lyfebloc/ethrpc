package ethrpc

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/Lyfebloc/logger"
)

const (
	MethodCall = "call"

	MethodAggregate = "aggregate"

	MethodTryAggregate = "tryAggregate"

	MethodGetCurrentBlockTimestamp = "getCurrentBlockTimestamp"
)

type (
	// RequestMiddleware type is for request middleware, called before a request is sent
	RequestMiddleware func(*Client, *Request) error

	// ResponseMiddleware type is for response middleware, called after a response has been received
	ResponseMiddleware func(*Client, *Response) error
)

type Client struct {
	ethClient         *ethclient.Client
	multiCallContract common.Address
	beforeRequest     []RequestMiddleware
	afterResponse     []ResponseMiddleware
}

func (c *Client) SetMulticallContract(multiCallContract common.Address) *Client {
	c.multiCallContract = multiCallContract

	return c
}

func (c *Client) R() *Request {
	r := &Request{
		client: c,
	}

	return r
}

func (c *Client) NewRequest() *Request {
	return c.R()
}

func (c *Client) execute(req *Request) (*Response, error) {
	var err error

	// Apply Request middlewares
	for _, f := range c.beforeRequest {
		if err = f(c, req); err != nil {
			return nil, err
		}
	}

	resp, err := c.ethClient.CallContract(req.Context(), req.RawCallMsg, nil)
	if err != nil {
		logger.Errorf("failed to call multicall, err: %v", err)
		return nil, err
	}

	response := &Response{
		Request:     req,
		RawResponse: resp,
	}

	// Apply Response middleware
	for _, f := range c.afterResponse {
		if err = f(c, response); err != nil {
			break
		}
	}

	return response, err
}

func createClient(ec *ethclient.Client) *Client {
	c := &Client{
		ethClient: ec,
	}

	// default before request middlewares
	c.beforeRequest = []RequestMiddleware{
		parseRequestCallParam,
	}

	// default after response middlewares
	c.afterResponse = []ResponseMiddleware{
		parseResponse,
	}

	return c
}
