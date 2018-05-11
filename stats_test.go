package jsonrpc_test

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/elliotchance/jsonrpc"
)

func TestSimpleServer_TotalPayloads(t *testing.T) {
	server := newTestServer()
	previousValue := 0
	assert.Equal(t, previousValue, server.TotalPayloads())

	t.Run("Handle", func(t *testing.T) {
		for testName, test := range specTests {
			server.Handle([]byte(test.j))

			assert.Equal(t, test.statsPayloads,server.TotalPayloads() -
				previousValue, "%s: %s", testName, test.j)
			previousValue = server.TotalPayloads()
		}
	})

	t.Run("HandleWithState", func(t *testing.T) {
		for testName, test := range specTests {
			server.HandleWithState([]byte(test.j), jsonrpc.State{})

			assert.Equal(t, test.statsPayloads, server.TotalPayloads() -
				previousValue, "%s: %s", testName, test.j)
			previousValue = server.TotalPayloads()
		}
	})

	t.Run("HandleRequest", func(t *testing.T) {
		for testName, test := range specTests {
			request, err := jsonrpc.NewRequestFromJSON([]byte(test.j))

			// We are only testing single requests here so ignore the ones that
			// are multi or invalid.
			if err != nil {
				continue
			}

			server.HandleRequest(request)

			assert.Equal(t, test.statsPayloads, server.TotalPayloads() -
				previousValue, "%s: %s", testName, test.j)
			previousValue = server.TotalPayloads()
		}
	})
}

func TestSimpleServer_TotalRequests(t *testing.T) {
	server := newTestServer()
	previousValue := 0
	assert.Equal(t, previousValue, server.TotalRequests())

	t.Run("Handle", func(t *testing.T) {
		for testName, test := range specTests {
			server.Handle([]byte(test.j))

			assert.Equal(t, test.statsRequests, server.TotalRequests() -
				previousValue, "%s: %s", testName, test.j)
			previousValue = server.TotalRequests()
		}
	})

	t.Run("HandleWithState", func(t *testing.T) {
		for testName, test := range specTests {
			server.HandleWithState([]byte(test.j), jsonrpc.State{})

			assert.Equal(t, test.statsRequests, server.TotalRequests() -
				previousValue, "%s: %s", testName, test.j)
			previousValue = server.TotalRequests()
		}
	})

	t.Run("HandleRequest", func(t *testing.T) {
		for testName, test := range specTests {
			request, err := jsonrpc.NewRequestFromJSON([]byte(test.j))

			// We are only testing single requests here so ignore the ones that
			// are multi or invalid.
			if err != nil {
				continue
			}

			server.HandleRequest(request)

			assert.Equal(t, test.statsRequests, server.TotalRequests() -
				previousValue, "%s: %s", testName, test.j)
			previousValue = server.TotalRequests()
		}
	})
}
