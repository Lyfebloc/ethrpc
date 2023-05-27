package ethrpc

import _ "embed"

//go:embed abis/Multicall.json
var multicallABIJson []byte

//go:embed abis/DmmPool.json
var poolABIJson []byte
