package jsonrpc

// StatReporter provides statistics for the JSON-RPC server.
type StatReporter interface {
	// The total payloads have been received by the server. The number of
	// payloads is not the number of requests received or processed. All
	// individual payloads including success, malformed, invalid, error, batch
	// or notification are considered a single payload.
	TotalPayloads() int

	// The number of requests received by the server, regardless of the
	// validity, result or if anything was sent back to the client. A batch
	// request is treated as multiple requests. An empty batch request is
	// considered to be zero requests.
	TotalRequests() int
}

func (server *SimpleServer) TotalPayloads() int {
	return server.totalPayloads
}

func (server *SimpleServer) TotalRequests() int {
	return server.totalRequests
}
