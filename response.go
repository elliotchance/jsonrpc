package jsonrpc

import (
	"encoding/json"
	"fmt"
)

const (
	// This is not declared in the JSON-RPC spec but it can be used with
	// ErrorCode() to test that there was no error.
	Success = 0

	// Invalid JSON was received by the server or an error occurred on the
	// server while parsing the JSON.
	ParseError = -32700

	// The JSON sent is not a valid Request object. This could be because
	// they have specific invalid types or that the version is not "2.0".
	InvalidRequest = -32600

	// The method does not exist or is not available.
	MethodNotFound = -32601

	// Invalid method parameter(s).
	InvalidParams = -32602

	// Internal JSON-RPC error. This means an error with something to do
	// with passing the request to the handler. If an error happens in the
	// handler itself it would be a ServerError.
	InternalError = -32603

	// Reserved for implementation-defined server-errors. This is the
	// maximum value permitted, other values for the same error can range to
	// -32099. These error codes would be understood by the receiver.
	ServerError = -32000

	// The lower bound for the server error range. You would probably not
	// use this constant directly unless you had a special reason to, use
	// jsonRpcServerError instead.
	ServerErrorMin = -32099
)

// Provides immutable information about a response. A response will either be a
// success or failure. This can be tested with:
//
//     if response.ErrorCode() == jsonrpc.Success {
//         // response.Result()
//     } else {
//         // response.ErrorCode()
//         // response.ErrorMessage()
//     }
//
type Response interface {
	fmt.Stringer
	Version() string
	Id() interface{}
	Result() interface{}
	ErrorCode() int
	ErrorMessage() string
}

type Responses []Response

// A JSON-RPC error is made up of a code and a message. It is acceptable for the
// message to be empty - the server will replace it with the generic message
// returned from ErrorMessageForCode().
type errorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// A JSON-RPC response object.
type response struct {
	ResponseVersion string         `json:"jsonrpc"`
	ResponseId      interface{}    `json:"id"`
	ResponseResult  interface{}    `json:"result,omitempty"`
	ResponseError   *errorResponse `json:"error,omitempty"`
}

func (response *response) Version() string {
	return response.ResponseVersion
}

func (response *response) Id() interface{} {
	return response.ResponseId
}

func (response *response) Result() interface{} {
	return response.ResponseResult
}

func (response *response) ErrorCode() int {
	if response.ResponseError == nil {
		return Success
	}

	return response.ResponseError.Code
}

func (response *response) ErrorMessage() string {
	if response.ResponseError == nil {
		return ""
	}

	return response.ResponseError.Message
}

// The string representation of a response will be the JSON encoded value. This
// JSON is expected to be a perfectly valid JSON-RPC response.
func (response *response) String() string {
	b, err := json.Marshal(response)
	if err != nil {
		// I don't know what would cause this situation. There is nothing we can
		// do except return an empty string (which would not occur in any
		// successful situation).
		return ""
	}

	return string(b)
}

// Create a response containing a successful response.
//
// Permitted values for id are a string or integer. Other types would be against
// the JSON-RPC specification and at your own risk. As long as you are passing
// through the original id from the request this shouldn't give you any issues.
//
// The result can be of any type, but the server will expect that the client can
// handle this appropriately.
func NewSuccessResponse(id interface{}, result interface{}) Response {
	return &response{
		ResponseVersion: "2.0",
		ResponseId:      id,
		ResponseResult:  result,
	}
}

// Create a response containing an error.
//
// Permitted values for id are a string or integer. Other types would be against
// the JSON-RPC specification and at your own risk. As long as you are passing
// through the original id from the request this shouldn't give you any issues.
//
// code is one of the integer constants, such as InvalidRequest. You may also
// use a range of integers to represent different types of ServerError but
// client will have to understand these codes.
//
// The message should be a human-readable description of the error and should
// not contain sensitive details (such as passwords). You may provide an empty
// string for message to use the message from ErrorMessageForCode() instead.
func NewErrorResponse(id interface{}, code int, message string) Response {
	if message == "" {
		message = ErrorMessageForCode(code)
	}

	return &response{
		ResponseVersion: "2.0",
		ResponseId:      id,
		ResponseError: &errorResponse{
			Code:    code,
			Message: message,
		},
	}
}

// A convenience method for converting a standard error into a ServerError.
//
// It is assumed to be a generic ServerError since that covers any general
// errors.
//
// If the parameters you receive are not valid or in a format that is understood
// (since they could be an array or a map) you should use:
//  ServerErrorResponse{Code:InvalidParams, Message:"Missing foo"}
func NewServerErrorResponse(id interface{}, err error) Response {
	return NewErrorResponse(id, ServerError, err.Error())
}

// Get the generic error message for the error code.
func ErrorMessageForCode(code int) string {
	switch code {
	case ParseError:
		return "Parse error"

	case InvalidRequest:
		return "Invalid request"

	case MethodNotFound:
		return "Method not found"

	case InvalidParams:
		return "Invalid params"

	case InternalError:
		return "Internal error"
	}

	if code >= ServerErrorMin && code <= ServerError {
		return "Server error"
	}

	return "Unknown error"
}

func (responses Responses) String() string {
	b, err := json.Marshal(responses)
	if err != nil {
		// I don't know what would cause this situation. I really don't
		// want to panic, so just return a different string instead.
		return ""
	}

	return string(b)
}

func NewResponsesFromJSON(data []byte) (Responses, error) {
	if data[0] == '[' {
		rawResponses := []*response{}
		err := json.Unmarshal(data, &rawResponses)
		if err != nil {
			return nil, err
		}

		responses := make([]Response, len(rawResponses))
		for i := range rawResponses {
			responses[i] = rawResponses[i]
		}

		return responses, err
	}

	response := new(response)
	err := json.Unmarshal(data, response)
	if err != nil {
		return nil, err
	}

	return Responses{response}, err
}
