package jsonrpc

import (
	"time"
	"sync/atomic"
)

// StatReporter provides statistics for the JSON-RPC server.
//
// You can see examples for the each of the statistics against different message
// types by looking at the specTests.
type StatReporter interface {
	// TotalPayloads is total payloads have been received by the server. The
	// number of payloads is not the number of requests received or processed.
	// All individual payloads including success, malformed, invalid, error,
	// batch or notification are considered a single payload.
	TotalPayloads() uint64

	// TotalRequests is the number of requests processed by the server. That is,
	// the number of requests that ended up calling a handler. Malformed and
	// invalid requests are not considered requests. Batch requests will only
	// count requests that call a handler. Other jobs in the batch will not be
	// counted towards the total requests.
	TotalRequests() uint64

	// TotalSuccessResponses returns the number of successful responses sent
	// back. A notification does not send back a result so it will not increment
	// this counter.
	TotalSuccessResponses() uint64

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
	TotalErrorResponses() uint64

	// TotalSuccessNotifications returns the number of notifications sent to the
	// server that returned success from the handler. Malformed or invalid
	// requests are not included in this count.
	TotalNotificationSuccesses() uint64

	// TotalNotificationErrors returns the number of notifications sent to the
	// server that did not return a success from the handler. Malformed or
	// invalid requests are not included in this count.
	TotalNotificationErrors() uint64

	// Uptime returns the duration that the server has been running for.
	Uptime() time.Duration

	// CurrentActiveRequests returns the number of requests that are inflight.
	// This does not include requests that are queued.
	CurrentActiveRequests() uint64
}

func (server *SimpleServer) TotalPayloads() uint64 {
	return atomic.LoadUint64(&server.totalPayloads)
}

func (server *SimpleServer) TotalRequests() uint64 {
	return atomic.LoadUint64(&server.totalRequests)
}

func (server *SimpleServer) TotalSuccessResponses() uint64 {
	return atomic.LoadUint64(&server.totalSuccessResponses)
}

func (server *SimpleServer) TotalErrorResponses() uint64 {
	return atomic.LoadUint64(&server.totalErrorResponses)
}

func (server *SimpleServer) TotalNotificationSuccesses() uint64 {
	return atomic.LoadUint64(&server.totalSuccessNotifications)
}

func (server *SimpleServer) TotalNotificationErrors() uint64 {
	return atomic.LoadUint64(&server.totalErrorNotifications)
}

func (server *SimpleServer) Uptime() time.Duration {
	return time.Now().Sub(server.startTime)
}

func (server *SimpleServer) CurrentActiveRequests() uint64 {
	return atomic.LoadUint64(&server.currentActiveRequests)
}
