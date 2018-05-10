package jsonrpc_test

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/elliotchance/jsonrpc"
)

func TestRequest_Bytes(t *testing.T) {
	request := jsonrpc.NewRequestResponder("2.0", 123, "foo", "bar")

	assert.Equal(t,
		"{\"jsonrpc\":\"2.0\",\"method\":\"foo\",\"params\":\"bar\",\"id\":123}",
		string(request.Bytes()))
}
