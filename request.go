package jsonrpc

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/rand"
	"strconv"
)

// Provides immutable information about a request.
type Request interface {
	fmt.Stringer
	Version() string
	Method() string
	Params() interface{}
	Id() interface{}
}

// Allows a request to produce responses. These are convenience functions so
// that the request ID (an potentially version) are set correctly in the
// response.
//
// You can use the similarly named functions if this interface is not available.
type Responder interface {
	NewSuccessResponse(result interface{}) Response
	NewErrorResponse(code int, message string) Response
	NewServerErrorResponse(err error) Response
}

// Provides immutable information about a request and has the ability to
// generate response objects from the request.
type RequestResponder interface {
	Request
	Responder
}

// A JSON-RPC request object.
type request struct {
	RequestVersion string      `json:"jsonrpc"`
	RequestMethod  string      `json:"method"`
	RequestParams  interface{} `json:"params,omitempty"`
	RequestId      interface{} `json:"id"`
}

func (request *request) Version() string {
	return request.RequestVersion
}

func (request *request) Method() string {
	return request.RequestMethod
}

func (request *request) Params() interface{} {
	return request.RequestParams
}

func (request *request) Id() interface{} {
	return request.RequestId
}

func (request *request) NewSuccessResponse(result interface{}) Response {
	return NewSuccessResponse(request.Id(), result)
}

func (request *request) NewErrorResponse(code int, message string) Response {
	return NewErrorResponse(request.Id(), code, message)
}

func (request *request) NewServerErrorResponse(err error) Response {
	return NewServerErrorResponse(request.Id(), err)
}

// The string representation of a request will be the JSON encoded value. This
// JSON is expected to be a perfectly valid JSON-RPC request.
func (request *request) String() string {
	b, err := json.Marshal(request)
	if err != nil {
		// I don't know what would cause this situation. I really don't
		// want to panic, so just return a different string instead.
		return "<Request>"
	}

	return string(b)
}

// Create a JSON-RPC request that is also able to produce responses.
//
// If the id is nil it will be considered a notification and no response will be
// send back.
//
// If params is nil then it will not be included, other acceptable types are an
// array or map for ordered and named-parameters respectively.
func NewRequestResponder(version string, id interface{}, method string, params interface{}) RequestResponder {
	return &request{
		RequestVersion: version,
		RequestId:      id,
		RequestMethod:  method,
		RequestParams:  params,
	}
}

// Use this with NewRequestResponder() to generate a random IDs for your
// requests. The generated ID will be a 32 digit hexadecimal value, the same you
// would see from an MD5.
//
// This can be important for logging to match up requests and response. It also
// allows tests to make sure the id is being maintained correctly.
func GenerateRequestId() string {
	hash := md5.Sum([]byte(strconv.Itoa(rand.Int())))
	return hex.EncodeToString(hash[:])
}
