package jsonrpc

// StatReporter provides statistics for the JSON-RPC server.
//
// You can see examples for the each of the statistics against different message
// types by looking at the specTests.
type StatReporter interface {
	// TotalPayloads is total payloads have been received by the server. The
	// number of payloads is not the number of requests received or processed.
	// All individual payloads including success, malformed, invalid, error,
	// batch or notification are considered a single payload.
	TotalPayloads() int

	// TotalRequests is the number of requests processed by the server. That is,
	// the number of requests that ended up calling a handler. Malformed and
	// invalid requests are not considered requests. Batch requests will only
	// count requests that call a handler. Other jobs in the batch will not be
	// counted towards the total requests.
	TotalRequests() int

	// TotalSuccessResponses returns the number of successful responses sent
	// back. A notification does not send back a result so it will not increment
	// this counter.
	TotalSuccessResponses() int

	// TotalErrorResponses returns the number of individual unsuccessful
	// responses sent back. A notification does not send back a result so it
	// will not increment this counter.
	//
	// This will also include requests that fail to even make it to the handler
	// and have to send back and error like a Parse error or invalid JSON-RPC
	// version.
	//
	// A batch request may contain zero or more failures if the JSON is not
	// malformed. However, a batch containing many jobs that is malformed JSON
	// (so the individual request cannot be parsed) will result in a single
	// Parse error sent back which will only count as one error response.
	TotalErrorResponses() int
}

func (server *SimpleServer) TotalPayloads() int {
	return server.totalPayloads
}

func (server *SimpleServer) TotalRequests() int {
	return server.totalRequests
}

func (server *SimpleServer) TotalSuccessResponses() int {
	return server.totalSuccessResponses
}

func (server *SimpleServer) TotalErrorResponses() int {
	return server.totalErrorResponses
}
