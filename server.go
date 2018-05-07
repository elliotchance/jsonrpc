package jsonrpc

import (
	"encoding/json"
)

// A handler is a function that is able to respond to a server request.
type RequestHandler func(RequestResponder) Response

type Server interface {
	SetHandler(methodName string, handler RequestHandler)
	HandleRequest(request RequestResponder) Response
	Handle(jsonRequest []byte) Responses
	HandleWithState(jsonRequest []byte, state State) Responses
	GetHandler(methodName string) RequestHandler
}

type SimpleServer struct {
	requestHandlers map[string]RequestHandler
}

// SetHandler will register (or replace) a handler for a method.
func (server *SimpleServer) SetHandler(methodName string, handler RequestHandler) {
	server.requestHandlers[methodName] = handler
}

func (server *SimpleServer) GetHandler(methodName string) RequestHandler {
	return server.requestHandlers[methodName]
}

// Requests can be handled two ways, but creating and passing a request
// directly:
//
//     request := jsonrpc.NewRequestResponder("1", "sayHello", map[string]string{"name": "Bob"})
//     response := server.HandleRequest(request)
//     fmt.Printf("%s", response.Result())
//
//     // Hello, Bob
//
// The first argument to NewRequest is the ID. This can be any string, integer
// or nil. If the ID is nil the request is called a "notification" and you will
// not receive a result from the server. In any other case the ID has no effect
// on how the request is processed. However, clients rely on this ID to be able
// to route and log results correctly back to where they came from.
//
// It is recommended that you always use a unique value for this. There is the
// provided function:
//
//     GenerateRequestId()
//
// The second method is to pass the raw request:
//
//     rawRequest := `{"jsonrpc": "2.0", "method": "sayHello", "params": {"name": "Bob"}, "id": 1}`
//     responses := server.Handle([]byte(rawRequest))
//     fmt.Printf("%s", responses[0].Result())
//
//     // Hello, Bob
//
// Handle() returns an array of Response interfaces to allow batch processing.
// The "Batch Requests" second explains this in more detail.
func (server *SimpleServer) HandleRequest(request RequestResponder) (response Response) {
	// Always recover from a panic and send it back as an error.
	defer func(id interface{}) {
		if r := recover(); r != nil {
			response = request.NewErrorResponse(ServerError, "")
		}
	}(request.Id())

	// We only support 2.0 right now.
	if request.Version() != "2.0" {
		return request.NewErrorResponse(InvalidRequest,
			"Version is not 2.0.")
	}

	handler := server.requestHandlers[request.Method()]
	if handler == nil {
		return request.NewErrorResponse(MethodNotFound, "")
	}

	return handler(request)
}

func (server *SimpleServer) handleSingle(jsonRequest []byte, isPartOfBatch bool, state State) Response {
	var requestMap map[string]interface{}
	err := json.Unmarshal(jsonRequest, &requestMap)
	if err != nil {
		errCode := ParseError

		// The JSON-RPC spec says that for a batch request, any
		// individual requests that would normally throw a ParseError
		// here should be treated as InvalidRequest instead.
		if isPartOfBatch {
			errCode = InvalidRequest
		}

		// It is unlikely that we will have an "id" but we might as well
		// try.
		return NewErrorResponse(requestMap["id"], errCode, "")
	}

	// Catch some type errors before creating the real request.
	if _, ok := requestMap["jsonrpc"].(string); !ok {
		return NewErrorResponse(requestMap["id"],
			InvalidRequest, "Version (jsonrpc) must be a string.")
	}
	if _, ok := requestMap["method"].(string); !ok {
		return NewErrorResponse(requestMap["id"],
			InvalidRequest, "Method must be a string.")
	}

	return server.HandleRequest(NewRequestResponderWithState(
		requestMap["jsonrpc"].(string),
		requestMap["id"],
		requestMap["method"].(string),
		requestMap["params"],
		state,
	))
}

func appendResponseIfNeeded(responses *Responses, response Response) {
	// Notifications do not receive results
	if response.Id() != nil || response.ErrorCode() != Success {
		*responses = append(*responses, response)
	}
}

// Batch Requests:
//
// Batch requests allow multiple requests to be handled as a single group. A
// batch request is simply an array of requests, which each one being processed
// independently and all results sent back at the same time.
//
//     rawRequest := `[
//       {"jsonrpc": "2.0", "method": "sayHello", "params": {"name": "Bob"}, "id": 1},
//       {"jsonrpc": "2.0", "method": "sayHello", "params": {"name": "John"}},
//       {"jsonrpc": "2.0", "method": "sayHello", "params": {"name": "Jane"}, "id": 2}
//     ]`
//     responses := server.Handle([]byte(rawRequest))
//     fmt.Printf("%s", responses[0].Result())
//     fmt.Printf("%s", responses[1].Result())
//
//     // Hello, Bob
//     // Hello, Jane
//
// You will get a Response for every non-notification (every request with a
// non-nil ID). The order of the responses is not predictable against the order
// of the requests. You should use the response IDs to correlate results in a
// batch result.
//
// It is also important to note that the order in which the requests are
// processed (whether single requests or batch) in a are non-deterministic and
// should be considered to be run all at the same time.
func (server *SimpleServer) HandleWithState(jsonRequest []byte, state State) Responses {
	responses := make(Responses, 0)

	// Check for a batch request.
	var batchRequest []interface{}
	err := json.Unmarshal(jsonRequest, &batchRequest)
	if err == nil {
		// It is a batch request, make sure it is not empty. Normally I
		// wouldn't care and happily return an empty array of results
		// back but the JSON-RPC spec says this is an invalid request.
		if len(batchRequest) == 0 {
			return Responses{NewErrorResponse(nil, InvalidRequest,
				"Batch is empty.")}
		}

		// Validate each of the requests because some of them may be
		// good and some invalid.
		for _, probableRequest := range batchRequest {
			// We have to mashal each request back to JSON, then
			// treat each one as an independent request.
			rawMessage, err := json.Marshal(probableRequest)
			if err != nil {
				// This condition should not be possible since
				// we have already unmarshalled this object
				// once. Still, better to be safe than sorry.
				response := NewErrorResponse(nil, ParseError,
					err.Error())
				appendResponseIfNeeded(&responses, response)
				continue
			}

			response := server.handleSingle(rawMessage, true, state)
			appendResponseIfNeeded(&responses, response)
		}
	} else {
		response := server.handleSingle(jsonRequest, false, state)
		appendResponseIfNeeded(&responses, response)
	}

	return responses
}

func (server *SimpleServer) Handle(jsonRequest []byte) Responses {
	return server.HandleWithState(jsonRequest, State{})
}

// Example:
//
//     func sayHello(request jsonrpc.RequestResponder) jsonrpc.Response {
//         return request.NewSuccessResponse("Hello, " + request.Param("name").(string))
//     }
//
//     server := jsonrpc.NewSimpleServer()
//     server.SetHandler("sayHello", sayHello)
func NewSimpleServer() *SimpleServer {
	server := new(SimpleServer)
	server.requestHandlers = make(map[string]RequestHandler)

	return server
}
