A simple JSON-RPC server for Go.

# Installation

```bash
go get -u github.com/elliotchance/jsonrpc
```

# Handler

A handler is a function that uses the following definition:

```go
func sum(request jsonrpc.RequestResponder) jsonrpc.Response {
	total := 0.0
	for _, x := range request.Params().([]interface{}) {
		total += x.(float64)
	}

	return request.NewSuccessResponse(total)
}
```

A handler must return `request.NewSuccessResponse` or
`request.NewErrorResponse`.

# Server

Creating a new server and attaching handlers:

```go
server := jsonrpc.NewSimpleServer()
server.SetHandler("sum", sum)
```

# Requests

The safest and easiest way to handle request is to pass the JSON bytes directly
to the `Handle` method of the server:

```go
responses := server.Handle(
	[]byte(`{"jsonrpc": "2.0", "method": "sum", "params": [1,2,4], "id": 1}`))
// [{"jsonrpc": "2.0", "result": 7, "id": 1}]
```

The JSON bytes could contain a single request or an array of requests (as
described in JSON-RPC 2.0). The number of responses returned may be zero or more
depending on if the requests are notifications.

There is no guaranteed order on the responses. You should use `Id()` to pair
responses with the appropriate request.

## Stateful Requests

Stateful requests allow you to pass extra state to the handler that only exist
for that single request.

State can be passed in by using `HandleWithState` with an extra parameter:

```go
responses := server.HandleWithState(
	[]byte(`{"jsonrpc": "2.0", "method": "add", "params": [13], "id": 1}`),
    jsonrpc.State{"offset": 25.0},
)
```

The handler can access state through `State(key)`:

```go
func add(request jsonrpc.RequestResponder) jsonrpc.Response {
	total := request.State("offset").(float64) + request.Params()[0].(float64)
		
	return request.NewSuccessResponse()
}
```

If the state `key` does not exist then `nil` is returned.
