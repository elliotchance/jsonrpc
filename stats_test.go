package jsonrpc_test

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/elliotchance/jsonrpc"
	"time"
)

func TestSimpleServer_TotalPayloads(t *testing.T) {
	server := newTestServer()
	previousValue := uint64(0)
	assert.Equal(t, previousValue, server.TotalPayloads())

	t.Run("Handle", func(t *testing.T) {
		for testName, test := range specTests {
			server.Handle([]byte(test.j))

			assert.Equal(t, test.statsPayloads, server.TotalPayloads()-
				previousValue, "%s: %s", testName, test.j)
			previousValue = server.TotalPayloads()
		}
	})

	t.Run("HandleWithState", func(t *testing.T) {
		for testName, test := range specTests {
			server.HandleWithState([]byte(test.j), jsonrpc.State{})

			assert.Equal(t, test.statsPayloads, server.TotalPayloads()-
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

			assert.Equal(t, test.statsPayloads, server.TotalPayloads()-
				previousValue, "%s: %s", testName, test.j)
			previousValue = server.TotalPayloads()
		}
	})
}

func TestSimpleServer_TotalRequests(t *testing.T) {
	server := newTestServer()
	previousValue := uint64(0)
	assert.Equal(t, previousValue, server.TotalRequests())

	t.Run("Handle", func(t *testing.T) {
		for testName, test := range specTests {
			server.Handle([]byte(test.j))

			assert.Equal(t, test.statsRequests, server.TotalRequests()-
				previousValue, "%s: %s", testName, test.j)
			previousValue = server.TotalRequests()
		}
	})

	t.Run("HandleWithState", func(t *testing.T) {
		for testName, test := range specTests {
			server.HandleWithState([]byte(test.j), jsonrpc.State{})

			assert.Equal(t, test.statsRequests, server.TotalRequests()-
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

			assert.Equal(t, test.statsRequests, server.TotalRequests()-
				previousValue, "%s: %s", testName, test.j)
			previousValue = server.TotalRequests()
		}
	})
}

func TestSimpleServer_TotalSuccessResponses(t *testing.T) {
	server := newTestServer()
	previousValue := uint64(0)
	assert.Equal(t, previousValue, server.TotalSuccessResponses())

	t.Run("Handle", func(t *testing.T) {
		for testName, test := range specTests {
			server.Handle([]byte(test.j))

			assert.Equal(t, test.statsSuccess, server.TotalSuccessResponses()-
				previousValue, "%s: %s", testName, test.j)
			previousValue = server.TotalSuccessResponses()
		}
	})

	t.Run("HandleWithState", func(t *testing.T) {
		for testName, test := range specTests {
			server.HandleWithState([]byte(test.j), jsonrpc.State{})

			assert.Equal(t, test.statsSuccess, server.TotalSuccessResponses()-
				previousValue, "%s: %s", testName, test.j)
			previousValue = server.TotalSuccessResponses()
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

			assert.Equal(t, test.statsSuccess, server.TotalSuccessResponses()-
				previousValue, "%s: %s", testName, test.j)
			previousValue = server.TotalSuccessResponses()
		}
	})
}

func TestSimpleServer_TotalErrorResponses(t *testing.T) {
	server := newTestServer()
	previousValue := uint64(0)
	assert.Equal(t, previousValue, server.TotalErrorResponses())

	t.Run("Handle", func(t *testing.T) {
		for testName, test := range specTests {
			server.Handle([]byte(test.j))

			assert.Equal(t, test.statsError, server.TotalErrorResponses()-
				previousValue, "%s: %s", testName, test.j)
			previousValue = server.TotalErrorResponses()
		}
	})

	t.Run("HandleWithState", func(t *testing.T) {
		for testName, test := range specTests {
			server.HandleWithState([]byte(test.j), jsonrpc.State{})

			assert.Equal(t, test.statsError, server.TotalErrorResponses()-
				previousValue, "%s: %s", testName, test.j)
			previousValue = server.TotalErrorResponses()
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

			assert.Equal(t, test.statsError, server.TotalErrorResponses()-
				previousValue, "%s: %s", testName, test.j)
			previousValue = server.TotalErrorResponses()
		}
	})
}

func TestSimpleServer_TotalNotificationSuccesses(t *testing.T) {
	server := newTestServer()
	previousValue := uint64(0)
	assert.Equal(t, previousValue, server.TotalNotificationSuccesses())

	t.Run("Handle", func(t *testing.T) {
		for testName, test := range specTests {
			server.Handle([]byte(test.j))

			assert.Equal(t, test.statsSuccessNotifications,
				server.TotalNotificationSuccesses()-previousValue,
				"%s: %s", testName, test.j)
			previousValue = server.TotalNotificationSuccesses()
		}
	})

	t.Run("HandleWithState", func(t *testing.T) {
		for testName, test := range specTests {
			server.HandleWithState([]byte(test.j), jsonrpc.State{})

			assert.Equal(t, test.statsSuccessNotifications,
				server.TotalNotificationSuccesses()-previousValue,
				"%s: %s", testName, test.j)
			previousValue = server.TotalNotificationSuccesses()
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

			assert.Equal(t, test.statsSuccessNotifications,
				server.TotalNotificationSuccesses()-previousValue,
				"%s: %s", testName, test.j)
			previousValue = server.TotalNotificationSuccesses()
		}
	})
}

func TestSimpleServer_TotalNotificationErrors(t *testing.T) {
	server := newTestServer()
	previousValue := uint64(0)
	assert.Equal(t, previousValue, server.TotalNotificationErrors())

	t.Run("Handle", func(t *testing.T) {
		for testName, test := range specTests {
			server.Handle([]byte(test.j))

			assert.Equal(t, test.statsErrorNotifications,
				server.TotalNotificationErrors()-previousValue,
				"%s: %s", testName, test.j)
			previousValue = server.TotalNotificationErrors()
		}
	})

	t.Run("HandleWithState", func(t *testing.T) {
		for testName, test := range specTests {
			server.HandleWithState([]byte(test.j), jsonrpc.State{})

			assert.Equal(t, test.statsErrorNotifications,
				server.TotalNotificationErrors()-previousValue,
				"%s: %s", testName, test.j)
			previousValue = server.TotalNotificationErrors()
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

			assert.Equal(t, test.statsErrorNotifications,
				server.TotalNotificationErrors()-previousValue,
				"%s: %s", testName, test.j)
			previousValue = server.TotalNotificationErrors()
		}
	})
}

func TestSimpleServer_Uptime(t *testing.T) {
	server := newTestServer()

	firstUptime := server.Uptime()
	assert.True(t, firstUptime > 0)

	t.Run("AfterMillisecond", func(t *testing.T) {
		time.Sleep(time.Millisecond)

		assert.True(t, server.Uptime()-firstUptime > 0)
		assert.True(t, server.Uptime()-firstUptime < 10*time.Millisecond)
	})
}

func TestSimpleServer_CurrentActiveRequests(t *testing.T) {
	server := newTestServer()

	assert.Equal(t, uint64(0), server.CurrentActiveRequests())

	t.Run("Handle", func(t *testing.T) {
		for _, test := range specTests {
			server.Handle([]byte(test.j))

			assert.Equal(t, uint64(0), server.CurrentActiveRequests())
		}
	})

	t.Run("HandleWithState", func(t *testing.T) {
		for _, test := range specTests {
			server.HandleWithState([]byte(test.j), jsonrpc.State{})

			assert.Equal(t, uint64(0), server.CurrentActiveRequests())
		}
	})

	t.Run("HandleRequest", func(t *testing.T) {
		for _, test := range specTests {
			request, err := jsonrpc.NewRequestFromJSON([]byte(test.j))

			// We are only testing single requests here so ignore the ones that
			// are multi or invalid.
			if err != nil {
				continue
			}

			server.HandleRequest(request)

			assert.Equal(t, uint64(0), server.CurrentActiveRequests())
		}
	})

	t.Run("DuringRequest", func(t *testing.T) {
		assert.Equal(t, uint64(0), server.CurrentActiveRequests())

		done := make(chan bool)
		go func() {
			server.Handle([]byte(`{"jsonrpc":"2.0","method":"hangUntilChannel"}`))
			done <- true
		}()

		<-hangStarted
		assert.Equal(t, uint64(1), server.CurrentActiveRequests())

		waitForChannel <- true
		<-done
		assert.Equal(t, uint64(0), server.CurrentActiveRequests())
	})
}
