package ethrpc

import (
	"bytes"

	"github.com/ethereum/go-ethereum/accounts/abi"
)

var (
	dmmPoolABI   abi.ABI
	multicallABI abi.ABI
)

func init() {
	builder := []struct {
		ABI  *abi.ABI
		data []byte
	}{
		{&dmmPoolABI, poolABIJson},
		{&multicallABI, multicallABIJson},
	}

	for _, b := range builder {
		var err error
		*b.ABI, err = abi.JSON(bytes.NewReader(b.data))
		if err != nil {
			panic(err)
		}
	}
}
