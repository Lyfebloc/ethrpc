package ethrpc

import (
	"github.com/ethereum/go-ethereum/ethclient"
)

// New method creates a new RPC client.
func New(url string) *Client {
	ec, err := ethclient.Dial(url)
	if err != nil {
		panic(err)
	}

	return createClient(ec)
}

// NewWithClient method creates a new RPC client with given `ethclient.Client`.
func NewWithClient(ec *ethclient.Client) *Client {
	return createClient(ec)
}
