package jsonrpc_test

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/elliotchance/jsonrpc"
)

func TestResponse_Bytes(t *testing.T) {
	response := jsonrpc.NewSuccessResponse("foo", "bar")

	assert.Equal(t,
		"{\"jsonrpc\":\"2.0\",\"id\":\"foo\",\"result\":\"bar\"}",
		string(response.Bytes()))
}

func TestResponses_Bytes(t *testing.T) {
	responses := jsonrpc.Responses{jsonrpc.NewSuccessResponse("foo", "bar")}

	assert.Equal(t,
		"[{\"jsonrpc\":\"2.0\",\"id\":\"foo\",\"result\":\"bar\"}]",
		string(responses.Bytes()))
}
