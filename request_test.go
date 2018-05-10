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

	t.Run("Empty", func(t *testing.T) {
		r, err := jsonrpc.NewRequestFromJSON([]byte(``))

		assert.EqualError(t, err, "Empty input")
		assert.Nil(t, r)
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

func TestNewRequestsFromJSON(t *testing.T) {
	t.Run("Single", func(t *testing.T) {
		request := jsonrpc.NewRequestResponder("2.0", 123, "foo", "bar")
		r, err := jsonrpc.NewRequestsFromJSON(request.Bytes())

		assert.NoError(t, err)
		assert.Len(t, r, 1)

		assert.Equal(t, r[0].Version(), "2.0")
		assert.Equal(t, r[0].Id(), 123.0)
		assert.Equal(t, r[0].Method(), "foo")
		assert.Equal(t, r[0].Params(), "bar")
	})

	t.Run("Empty", func(t *testing.T) {
		r, err := jsonrpc.NewRequestsFromJSON([]byte(``))

		assert.EqualError(t, err, "Empty input")
		assert.Nil(t, r)
	})

	t.Run("Malformed", func(t *testing.T) {
		r, err := jsonrpc.NewRequestsFromJSON([]byte(`{`))

		assert.EqualError(t, err, "Parse error")
		assert.Nil(t, r)
	})

	t.Run("Batch", func(t *testing.T) {
		request1 := jsonrpc.NewRequestResponder("2.0", 123, "foo", "bar")
		request2 := jsonrpc.NewRequestResponder("2.0", 456, "baz", "qux")
		r, err := jsonrpc.NewRequestsFromJSON([]byte(
			"[" + request1.String() + "," + request2.String() + "]"))

		assert.NoError(t, err)
		assert.Len(t, r, 2)

		assert.Equal(t, r[0].Version(), "2.0")
		assert.Equal(t, r[0].Id(), 123.0)
		assert.Equal(t, r[0].Method(), "foo")
		assert.Equal(t, r[0].Params(), "bar")

		assert.Equal(t, r[1].Version(), "2.0")
		assert.Equal(t, r[1].Id(), 456.0)
		assert.Equal(t, r[1].Method(), "baz")
		assert.Equal(t, r[1].Params(), "qux")
	})

	t.Run("BadVersionType", func(t *testing.T) {
		request := `{"jsonrpc":2, "id":123, "method":"foo"}`
		r, err := jsonrpc.NewRequestsFromJSON([]byte(request))

		assert.EqualError(t, err, "Version (jsonrpc) must be a string.")
		assert.Nil(t, r)
	})

	t.Run("BadVersionType", func(t *testing.T) {
		request := `{"jsonrpc":"2.0", "id":123, "method":null}`
		r, err := jsonrpc.NewRequestsFromJSON([]byte(request))

		assert.EqualError(t, err, "Method must be a string.")
		assert.Nil(t, r)
	})
}
