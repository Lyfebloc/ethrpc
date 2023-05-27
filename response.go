package ethrpc

type Response struct {
	Request     *Request
	RawResponse []byte
	// Result is an array that contains response result for all calls in the request
	Result []bool
}
