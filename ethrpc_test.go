package ethrpc

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/suite"
)

type RPCTestSuite struct {
	suite.Suite

	client *Client
}

func (ts *RPCTestSuite) SetupTest() {
	// Setup RPC server
	rpcClient := New("https://eth.llamarpc.com")
	rpcClient.SetMulticallContract(common.HexToAddress("0x5ba1e12693dc8f9c48aad8770482f4739beed696"))

	ts.client = rpcClient
}

func (ts *RPCTestSuite) TestTryAggregate() {
	type TradeInfo struct {
		Reserve0       *big.Int
		Reserve1       *big.Int
		VReserve0      *big.Int
		VReserve1      *big.Int
		FeeInPrecision *big.Int
	}

	pools := []string{
		"0x88e6A0c2dDD26FEEb64F039a2c41296FcB3f5640",
		"0xCBCdF9626bC03E24f779434178A73a0B4bad62eD",
		"0x8ad599c3A0ff1De082011EFDDc58f1908eb6e6D8",
	}

	reserves := make([]TradeInfo, len(pools))
	req := ts.client.NewRequest()

	for i, p := range pools {
		req.AddCall(&Call{
			ABI:    dmmPoolABI,
			Target: p,
			Method: "getTradeInfo",
			Params: nil,
		}, []interface{}{&reserves[i]})
	}

	res, err := req.TryAggregate()

	fmt.Printf("%+v", reserves)

	ts.Require().NoError(err)
	ts.Require().Len(res.Result, len(req.Calls))
}

func TestRPCTestSuite(t *testing.T) {
	suite.Run(t, new(RPCTestSuite))
}
