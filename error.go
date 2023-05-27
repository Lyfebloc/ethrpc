package ethrpc

import (
	"errors"
	"fmt"
)

var (
	ErrMethodNotSupported = errors.New("method not supported")
	ErrWrongCallParam     = errors.New("wrong call param")
	ErrUnexpectedResponse = errors.New("unexpected response")
)

type UnPackMulticallError struct {
	OriginalErr error
}

func NewUnPackMulticallError(originalErr error) error {
	return &UnPackMulticallError{
		OriginalErr: originalErr,
	}
}

func (e UnPackMulticallError) Error() string {
	return fmt.Sprintf("Unpack Multicall Error: %v", e.OriginalErr)
}
