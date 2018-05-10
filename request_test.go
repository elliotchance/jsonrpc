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

func TestNewRequestFromJSON(t *testing.T) {
	t.Run("Single", func(t *testing.T) {
		request := jsonrpc.NewRequestResponder("2.0", 123, "foo", "bar")
		r, err := jsonrpc.NewRequestFromJSON(request.Bytes())

		assert.NoError(t, err)
		assert.Equal(t, r.Version(), "2.0")
		assert.Equal(t, r.Id(), 123.0)
		assert.Equal(t, r.Method(), "foo")
		assert.Equal(t, r.Params(), "bar")
	})

	t.Run("Malformed", func(t *testing.T) {
		r, err := jsonrpc.NewRequestFromJSON([]byte(`{`))

		assert.EqualError(t, err, "Parse error")
		assert.Nil(t, r)
	})

	t.Run("Batch", func(t *testing.T) {
		request := jsonrpc.NewRequestResponder("2.0", 123, "foo", "bar")
		r, err :=
			jsonrpc.NewRequestFromJSON([]byte("[" + request.String() + "]"))

		assert.EqualError(t, err, "Parse error")
		assert.Nil(t, r)
	})

	t.Run("BadVersionType", func(t *testing.T) {
		request := `{"jsonrpc":2, "id":123, "method":"foo"}`
		r, err := jsonrpc.NewRequestFromJSON([]byte(request))

		assert.EqualError(t, err, "Version (jsonrpc) must be a string.")
		assert.Nil(t, r)
	})

	t.Run("BadVersionType", func(t *testing.T) {
		request := `{"jsonrpc":"2.0", "id":123, "method":null}`
		r, err := jsonrpc.NewRequestFromJSON([]byte(request))

		assert.EqualError(t, err, "Method must be a string.")
		assert.Nil(t, r)
	})
}
