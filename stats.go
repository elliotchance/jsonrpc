package jsonrpc

// StatReporter provides statistics for the JSON-RPC server.
//
// You can see examples for the each of the statistics against different message
// types by looking at the specTests.
type StatReporter interface {
	// The total payloads have been received by the server. The number of
	// payloads is not the number of requests received or processed. All
	// individual payloads including success, malformed, invalid, error, batch
	// or notification are considered a single payload.
	TotalPayloads() int

	// The number of requests processed by the server. That is, the number of
	// requests that ended up calling a handler. Malformed and invalid requests
	// are not considered requests. Batch requests will only count requests that
	// call a handler. Other jobs in the batch will not be counted towards the
	// total requests.
	TotalRequests() int
}

func (server *SimpleServer) TotalPayloads() int {
	return server.totalPayloads
}

func (server *SimpleServer) TotalRequests() int {
	return server.totalRequests
}
