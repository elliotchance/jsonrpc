package jsonrpc_test

import (
	"errors"
	"github.com/elliotchance/jsonrpc"
	"math/rand"
	"reflect"
	"regexp"
	"testing"
	"github.com/stretchr/testify/assert"
	"fmt"
)

func TestErrorMessageForCode(t *testing.T) {
	tests := map[string]struct {
		code    int
		message string
	}{
		"0 => Unknown error":            {0, "Unknown error"},
		"-1 (invalid) => Unknown error": {0, "Unknown error"},
		"-32700 => Parse error":         {jsonrpc.ParseError, "Parse error"},
		"-32600 => Invalid request":     {jsonrpc.InvalidRequest, "Invalid request"},
		"-32601 => Method not found":    {jsonrpc.MethodNotFound, "Method not found"},
		"-32602 => Invalid params":      {jsonrpc.InvalidParams, "Invalid params"},
		"-32603 => Internal error":      {jsonrpc.InternalError, "Internal error"},
		"-32000 => Server error 1":      {jsonrpc.ServerError, "Server error"},
		"-32000 => Server error 2":      {-32000 - rand.Intn(98), "Server error"},
		"-32098 => Server error":        {-32098, "Server error"},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			if jsonrpc.ErrorMessageForCode(test.code) != test.message {
				t.Errorf("TestErrorMessageForCode: %v", test.message)
			}
		})
	}
}

func TestGenerateRequestId(t *testing.T) {
	// Try a bunch of times to make sure it's unique.
	tries := 10
	values := map[string]bool{}

	for i := 0; i < tries; i += 1 {
		id := jsonrpc.GenerateRequestId()
		if !regexp.MustCompile("[0-9a-f]{16}").MatchString(id) {
			t.Errorf("TestGenerateRequestId: %v", id)
		}

		values[id] = true
	}

	if len(values) != tries {
		t.Errorf("TestGenerateRequestId: %v != %v", len(values), tries)
	}
}

func TestNewRequestResponder(t *testing.T) {
	request := jsonrpc.NewRequestResponder("2.0", "1", "method", []int{1, 2})

	if request.Version() != "2.0" {
		t.Errorf("TestNewRequestResponder: Version: %v != %v", request.Version(), "2.0")
	}
	if request.Id() != "1" {
		t.Errorf("TestNewRequestResponder: Id: %v != %v", request.Id(), "1")
	}
	if request.Method() != "method" {
		t.Errorf("TestNewRequestResponder: Method: %v != %v", request.Method(), "method")
	}
	if !reflect.DeepEqual(request.Params(), []int{1, 2}) {
		t.Errorf("TestNewRequestResponder: Params: %v != %v", request.Params(), []int{1, 2})
	}
}

func TestRequest_String(t *testing.T) {
	request := jsonrpc.NewRequestResponder("2.0", "1", "method", []int{1, 2})

	expected := `{"jsonrpc":"2.0","method":"method","params":[1,2],"id":"1"}`
	if request.String() != expected {
		t.Errorf("TestRequest_String: %v != %v", request.String(), expected)
	}
}

func TestResponse_String(t *testing.T) {
	t.Run("will render success response as JSON", func(t *testing.T) {
		response := jsonrpc.NewSuccessResponse("1", []int{1, 2})

		expected := `{"jsonrpc":"2.0","id":"1","result":[1,2]}`
		if response.String() != expected {
			t.Errorf("TestResponse_String: %v != %v", response.String(), expected)
		}
	})

	t.Run("will render error response as JSON", func(t *testing.T) {
		response := jsonrpc.NewErrorResponse(1, jsonrpc.InvalidRequest, "Oops!")

		expected := `{"jsonrpc":"2.0","id":1,"error":{"code":-32600,"message":"Oops!"}}`
		if response.String() != expected {
			t.Errorf("TestResponse_String: %v != %v", response.String(), expected)
		}
	})

	t.Run("will render server error response as JSON", func(t *testing.T) {
		response := jsonrpc.NewServerErrorResponse(2, errors.New("bad stuff happened"))

		expected := `{"jsonrpc":"2.0","id":2,"error":{"code":-32000,"message":"bad stuff happened"}}`
		if response.String() != expected {
			t.Errorf("TestResponse_String: %v != %v", response.String(), expected)
		}
	})
}

func TestSimpleServer_SetHandler(t *testing.T) {
	t.Run("handler can be replaced", func(t *testing.T) {
		server := jsonrpc.NewSimpleServer()

		server.SetHandler("subtract", subtract)
		server.SetHandler("subtract", sum)

		actual := reflect.ValueOf(server.GetHandler("subtract")).Pointer()
		expected := reflect.ValueOf(sum).Pointer()

		if actual != expected {
			t.Errorf("TestSimpleServer_SetHandler: %v != %v", actual, expected)
		}
	})

	t.Run("missing handler is nil", func(t *testing.T) {
		server := jsonrpc.NewSimpleServer()

		actual := server.GetHandler("subtract")
		if actual != nil {
			t.Errorf("TestSimpleServer_SetHandler: %v != nil", actual)
		}
	})
}

func TestJSONRPCSpecification(t *testing.T) {
	// All of these examples were provided from the official spec at:
	// http://www.jsonrpc.org/specification#examples
	tests := map[string]struct {
		input           string
		expectedResults jsonrpc.Responses
	}{
		"rpc call with positional parameters 1": {
			`{"jsonrpc": "2.0", "method": "subtract", "params": [42, 23], "id": 1}`,
			// `{"jsonrpc": "2.0", "result": 19, "id": 1}`,
			jsonrpc.Responses{
				jsonrpc.NewSuccessResponse(float64(1), float64(19)),
			},
		},
		"rpc call with positional parameters 2": {
			`{"jsonrpc": "2.0", "method": "subtract", "params": [23, 42], "id": 2}`,
			// `{"jsonrpc": "2.0", "result": -19, "id": 2}`,
			jsonrpc.Responses{
				jsonrpc.NewSuccessResponse(float64(2), float64(-19)),
			},
		},
		"rpc call with named parameters 1": {
			`{"jsonrpc": "2.0", "method": "subtract", "params": {"subtrahend": 23, "minuend": 42}, "id": 3}`,
			// `{"jsonrpc": "2.0", "result": 19, "id": 3}`,
			jsonrpc.Responses{
				jsonrpc.NewSuccessResponse(float64(3), float64(19)),
			},
		},
		"rpc call with named parameters 2": {
			`{"jsonrpc": "2.0", "method": "subtract", "params": {"minuend": 42, "subtrahend": 23}, "id": 4}`,
			// `{"jsonrpc": "2.0", "result": 19, "id": 4}`,
			jsonrpc.Responses{
				jsonrpc.NewSuccessResponse(float64(4), float64(19)),
			},
		},
		"a notification 1": {
			`{"jsonrpc": "2.0", "method": "subtract", "params": [1,2,3,4,5]}`,
			// ``,
			jsonrpc.Responses{},
		},
		"a notification 2": {
			`{"jsonrpc": "2.0", "method": "subtract"}`,
			// ``,
			jsonrpc.Responses{},
		},
		"rpc call of non-existent method": {
			`{"jsonrpc": "2.0", "method": "foobar", "id": 1}`,
			// `{"jsonrpc": "2.0", "error": {"code": -32601, "message": "Method not found"}, "id": "1"}`,
			jsonrpc.Responses{
				jsonrpc.NewErrorResponse(float64(1), jsonrpc.MethodNotFound, ""),
			},
		},
		"rpc call with invalid JSON": {
			`{"jsonrpc": "2.0", "method": "foobar, "params": "bar", "baz]`,
			// `{"jsonrpc": "2.0", "error": {"code": -32700, "message": "Parse error"}, "id": null}`,
			jsonrpc.Responses{
				jsonrpc.NewErrorResponse(nil, jsonrpc.ParseError, ""),
			},
		},
		"rpc call with invalid Request object": {
			`{"jsonrpc": "2.0", "method": 1, "params": "bar"}`,
			// `{"jsonrpc": "2.0", "error": {"code": -32600, "message": "Invalid request"}, "id": null}`,
			jsonrpc.Responses{
				jsonrpc.NewErrorResponse(nil, jsonrpc.InvalidRequest, "Method must be a string."),
			},
		},
		"rpc call Batch, invalid JSON": {
			`[
				{"jsonrpc": "2.0", "method": "sum", "params": [1,2,4], "id": "1"},
				{"jsonrpc": "2.0", "method"
			]`,
			// `{"jsonrpc": "2.0", "error": {"code": -32700, "message": "Parse error"}, "id": null}`,
			jsonrpc.Responses{
				jsonrpc.NewErrorResponse(nil, jsonrpc.ParseError, ""),
			},
		},
		"rpc call with an empty Array": {
			`[]`,
			// `{"jsonrpc": "2.0", "error": {"code": -32600, "message": "Invalid request"}, "id": null}`,
			jsonrpc.Responses{
				jsonrpc.NewErrorResponse(nil, jsonrpc.InvalidRequest, "Batch is empty."),
			},
		},
		"rpc call with an invalid Batch (but not empty)": {
			`[1]`,
			// `[{"jsonrpc": "2.0", "error": {"code": -32600, "message": "Invalid request"}, "id": null}]`,
			jsonrpc.Responses{
				jsonrpc.NewErrorResponse(nil, jsonrpc.InvalidRequest, ""),
			},
		},
		"rpc call with invalid Batch": {
			`[1,2,3]`,
			// `[
			// 	{"jsonrpc": "2.0", "error": {"code": -32600, "message": "Invalid request"}, "id": null},
			// 	{"jsonrpc": "2.0", "error": {"code": -32600, "message": "Invalid request"}, "id": null},
			// 	{"jsonrpc": "2.0", "error": {"code": -32600, "message": "Invalid request"}, "id": null}
			// ]`,
			jsonrpc.Responses{
				jsonrpc.NewErrorResponse(nil, jsonrpc.InvalidRequest, ""),
				jsonrpc.NewErrorResponse(nil, jsonrpc.InvalidRequest, ""),
				jsonrpc.NewErrorResponse(nil, jsonrpc.InvalidRequest, ""),
			},
		},
		"rpc call Batch": {
			`[
				{"jsonrpc": "2.0", "method": "sum", "params": [1,2,4], "id": 1},
				{"jsonrpc": "2.0", "method": "notify_hello", "params": [7]},
				{"jsonrpc": "2.0", "method": "subtract", "params": [42,23], "id": 2},
				{"foo": "boo"},
				{"jsonrpc": "2.0", "method": "foo.get", "params": {"name": "myself"}, "id": 5},
				{"jsonrpc": "2.0", "method": "get_data", "id": 9}
			]`,
			// `[
			// 	{"jsonrpc": "2.0", "result": 7, "id": 1},
			// 	{"jsonrpc": "2.0", "result": 19, "id": 2},
			// 	{"jsonrpc": "2.0", "error": {"code": -32600, "message": "Invalid request"}, "id": null},
			// 	{"jsonrpc": "2.0", "error": {"code": -32601, "message": "Method not found"}, "id": 5},
			// 	{"jsonrpc": "2.0", "result": ["hello", 5], "id": 9}
			// ]`,
			jsonrpc.Responses{
				jsonrpc.NewSuccessResponse(float64(1), float64(7)),
				jsonrpc.NewSuccessResponse(float64(2), float64(19)),
				jsonrpc.NewErrorResponse(nil, jsonrpc.InvalidRequest, "Version (jsonrpc) must be a string."),
				jsonrpc.NewErrorResponse(float64(5), jsonrpc.MethodNotFound, ""),
				jsonrpc.NewSuccessResponse(float64(9), []interface{}{"hello", float64(5)}),
			},
		},
		"rpc call Batch (all notifications)": {
			`[
				{"jsonrpc": "2.0", "method": "sum", "params": [1,2,4]},
				{"jsonrpc": "2.0", "method": "notify_hello", "params": [7]}
			]`,
			// ``,
			jsonrpc.Responses{},
		},

		// The tests below are extras for other edge cases not covered above.
		"wrong version": {
			`{"jsonrpc": "2", "method": "subtract", "params": [42, 23], "id": 2}`,
			// `{"jsonrpc": "2.0", "error": {"code": -32600, "message": "Invalid request"}, "id": 2}`,
			jsonrpc.Responses{
				jsonrpc.NewErrorResponse(float64(2), jsonrpc.InvalidRequest, "Version is not 2.0."),
			},
		},
		"bad version": {
			`{"jsonrpc": true, "method": "subtract", "params": [42, 23], "id": 2}`,
			// `{"jsonrpc": "2.0", "error": {"code": -32600, "message": "Invalid request"}, "id": 2}`,
			jsonrpc.Responses{
				jsonrpc.NewErrorResponse(float64(2), jsonrpc.InvalidRequest, "Version (jsonrpc) must be a string."),
			},
		},
		"missing version": {
			`{"method": "subtract", "params": [42, 23], "id": 2}`,
			// `{"jsonrpc": "2.0", "error": {"code": -32600, "message": "Invalid request"}, "id": 2}`,
			jsonrpc.Responses{
				jsonrpc.NewErrorResponse(float64(2), jsonrpc.InvalidRequest, "Version (jsonrpc) must be a string."),
			},
		},

		// The server much always recover from a panic(). We do not
		// return the error because it ay contain sensitive information.
		// Instead a generic Internal error will do.
		"recover from panic": {
			`{"jsonrpc": "2.0", "method": "panic", "id": 2}`,
			// `{"jsonrpc": "2.0", "error": {"code": -32000, "message": "Server error"}, "id": 2}`,
			jsonrpc.Responses{
				jsonrpc.NewErrorResponse(float64(2), jsonrpc.ServerError, ""),
			},
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			server := newTestServer()
			responses := server.Handle([]byte(test.input))

			if !reflect.DeepEqual(responses, test.expectedResults) {
				t.Errorf("TestJSONRPCSpecification: %v != %v", responses, test.expectedResults)
			}
		})
	}
}

func newTestServer() *jsonrpc.SimpleServer {
	server := jsonrpc.NewSimpleServer()

	server.SetHandler("subtract", subtract)
	server.SetHandler("sum", sum)
	server.SetHandler("notify_hello", notifyHello)
	server.SetHandler("get_data", getData)
	server.SetHandler("panic", forcePanic)
	server.SetHandler("handlerWithState", handlerWithState)

	return server
}

//noinspection GoUnusedParameter
func subtract(request jsonrpc.RequestResponder) jsonrpc.Response {
	switch p := request.Params().(type) {
	case []interface{}:
		return request.NewSuccessResponse(p[0].(float64) - p[1].(float64))
	case map[string]interface{}:
		return request.NewSuccessResponse(p["minuend"].(float64) - p["subtrahend"].(float64))
	}

	return request.NewSuccessResponse(nil)
}

//noinspection GoUnusedParameter
func sum(request jsonrpc.RequestResponder) jsonrpc.Response {
	total := 0.0
	for _, x := range request.Params().([]interface{}) {
		total += x.(float64)
	}

	return request.NewSuccessResponse(total)
}

//noinspection GoUnusedParameter
func notifyHello(request jsonrpc.RequestResponder) jsonrpc.Response {
	return request.NewSuccessResponse(nil)
}

//noinspection GoUnusedParameter
func getData(request jsonrpc.RequestResponder) jsonrpc.Response {
	return request.NewSuccessResponse([]interface{}{"hello", 5.0})
}

//noinspection GoUnusedParameter
func forcePanic(request jsonrpc.RequestResponder) jsonrpc.Response {
	panic("uh-oh!")

	return request.NewSuccessResponse(nil)
}

func handlerWithState(request jsonrpc.RequestResponder) jsonrpc.Response {
	return request.NewSuccessResponse(request.State("foo"))
}

func TestStatefulRequestMissingKey(t *testing.T) {
	server := newTestServer()
	r := `{"jsonrpc": "2.0", "method": "handlerWithState", "params": [42, 23], "id": 1}`
	responses := server.Handle([]byte(r))

	assert.Len(t, responses, 1)
	assert.Nil(t, responses[0].Result())
}

func TestStatefulRequestWithKey(t *testing.T) {
	server := newTestServer()
	r := `{"jsonrpc": "2.0", "method": "handlerWithState", "params": [42, 23], "id": 1}`
	state := jsonrpc.State{
		"foo": "bar",
	}
	responses := server.HandleWithState([]byte(r), state)

	assert.Len(t, responses, 1)
	assert.Equal(t, "bar", responses[0].Result())
}

func TestSimpleServerIsAServer(t *testing.T) {
	server := newTestServer()
	assert.Implements(t, (*jsonrpc.Server)(nil), server)
}

func TestRequestResponderIsAStringer(t *testing.T) {
	request := jsonrpc.NewRequestResponder("2.0", 123, "foo", nil)
	assert.Implements(t, (*fmt.Stringer)(nil), request)
}

func TestResponseIsAStringer(t *testing.T) {
	response := jsonrpc.NewSuccessResponse(123, "foo")
	assert.Implements(t, (*fmt.Stringer)(nil), response)
}

func TestResponseStringIsJSON(t *testing.T) {
	response := jsonrpc.NewSuccessResponse(123, "foo")
	assert.Equal(t, "{\"jsonrpc\":\"2.0\",\"id\":123,\"result\":\"foo\"}", response.String())
}

func TestResponsesIsAStringer(t *testing.T) {
	responses := jsonrpc.Responses{jsonrpc.NewSuccessResponse(123, "foo")}
	assert.Implements(t, (*fmt.Stringer)(nil), responses)
}
