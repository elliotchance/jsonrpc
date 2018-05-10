package jsonrpc_test

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/elliotchance/jsonrpc"
)

func TestSimpleServer_TotalPayloads(t *testing.T) {
	t.Run("Handle", func(t *testing.T) {
		server := newTestServer()
		total := 0
		for _, test := range specTests {
			server.Handle([]byte(test.j))
			total += test.statsPayloads
		}

		assert.Equal(t, total, server.TotalPayloads())
	})

	t.Run("HandleWithState", func(t *testing.T) {
		server := newTestServer()
		total := 0
		for _, test := range specTests {
			server.HandleWithState([]byte(test.j), jsonrpc.State{})
			total += test.statsPayloads
		}

		assert.Equal(t, total, server.TotalPayloads())
	})

	t.Run("HandleRequest", func(t *testing.T) {
		server := newTestServer()
		total := 0
		for _, test := range specTests {
			request, err := jsonrpc.NewRequestFromJSON([]byte(test.j))

			// We are only testing single requests here so ignore the ones that
			// are multi or invalid.
			if err != nil {
				continue
			}

			server.HandleRequest(request)
			total += test.statsPayloads
		}

		assert.Equal(t, total, server.TotalPayloads())
	})
}
