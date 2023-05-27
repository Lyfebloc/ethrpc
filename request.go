package ethrpc

import (
	"context"
	"encoding/binary"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
)

type Call struct {
	ABI       abi.ABI
	UnpackABI []abi.ABI
	Target    string
	Method    string
	Params    []interface{}
	Output    []interface{}
}

func (c *Call) SetOutput(output []interface{}) *Call {
	c.Output = output

	return c
}

// autofillUnpackABI fills the call's UnpackABI in case it's not set
func (c *Call) autofillUnpackABI() {
	if c.UnpackABI == nil {
		c.UnpackABI = []abi.ABI{c.ABI}

		return
	}

	if len(c.UnpackABI) == 0 {
		c.UnpackABI = append(c.UnpackABI, c.ABI)

		return
	}
}

type Request struct {
	client         *Client
	Method         string
	RequireSuccess bool
	Calls          []*Call
	ctx            context.Context
	RawCallMsg     ethereum.CallMsg
}

// Context method returns the Context if it's already set in request
// otherwise it creates new one using `context.Background()`.
func (r *Request) Context() context.Context {
	if r.ctx == nil {
		return context.Background()
	}
	return r.ctx
}

// SetContext method sets the context.Context for current Request. It allows
// to interrupt the request execution if ctx.Done() channel is closed.
// See https://blog.golang.org/context article and the "context" package
// documentation.
func (r *Request) SetContext(ctx context.Context) *Request {
	r.ctx = ctx
	return r
}

// AddCall adds a call to the request
// it will autofill the UnpackABI in case it's not set
func (r *Request) AddCall(c *Call, output []interface{}) *Request {
	c.autofillUnpackABI()
	c.SetOutput(output)
	r.Calls = append(r.Calls, c)

	return r
}

func (r *Request) SetRequireSuccess(requireSuccess bool) *Request {
	r.RequireSuccess = requireSuccess

	return r
}

func (r *Request) Execute(method string) (*Response, error) {
	r.Method = method

	return r.client.execute(r)
}

func (r *Request) Call() (*Response, error) {
	return r.Execute(MethodCall)
}

func (r *Request) Aggregate() (*Response, error) {
	return r.Execute(MethodAggregate)
}

func (r *Request) TryAggregate() (*Response, error) {
	return r.Execute(MethodTryAggregate)
}

func (r *Request) GetCurrentBlockTimestamp() (uint64, error) {
	res, err := r.Execute(MethodGetCurrentBlockTimestamp)
	if err != nil {
		return 0, err
	}

	// Block timestamp is a uint256 number so returned byte array should have length equals 256/8.
	if len(res.RawResponse) != 256/8 {
		return 0, ErrUnexpectedResponse
	}

	// Golang does not have uint256 but uint64 is enough to store a timestamp value.
	// Data is transferred in network byte order therefore just need to parse last 8 bytes in array.
	blockTimestamp := binary.BigEndian.Uint64(res.RawResponse[len(res.RawResponse)-8:])

	return blockTimestamp, nil
}
