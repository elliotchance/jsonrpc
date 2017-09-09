A simple JSON-RPC server for Go.

# Installation

```bash
go get -u github.com/elliotchance/jsonrpc
```

# Example

Create the handler function:

```go
func sum(request jsonrpc.RequestResponder) jsonrpc.Response {
	total := 0.0
	for _, x := range request.Params().([]interface{}) {
		total += x.(float64)
	}

	return request.NewSuccessResponse(total)
}
```

Create the server:

```go
server := jsonrpc.NewSimpleServer()
server.SetHandler("sum", sum)
```

Making requests:

```go
responses := server.Handle([]byte(`{"jsonrpc": "2.0", "method": "sum", "params": [1,2,4], "id": 1}`))
// [{"jsonrpc": "2.0", "result": 7, "id": 1}]
```
