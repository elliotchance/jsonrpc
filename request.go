package jsonrpc

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/rand"
	"strconv"
	"github.com/pkg/errors"
)

// Provides immutable information about a request.
type Request interface {
	Version() string
	Method() string
	Params() interface{}
	Id() interface{}
	State(key string) interface{}

	// Serialization
	fmt.Stringer
	Bytes() []byte
}

// State can be optionally provided with Handle requests to pass extra state to
// the handler for that individual request.
type State map[string]interface{}

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
	requestState   State
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

func (request *request) State(key string) interface{} {
	return request.requestState[key]
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
	return string(request.Bytes())
}

// Create a JSON-RPC request that is also able to produce responses.
//
// If the id is nil it will be considered a notification and no response will be
// send back.
//
// If params is nil then it will not be included, other acceptable types are an
// array or map for ordered and named-parameters respectively.
func NewRequestResponderWithState(version string, id interface{}, method string, params interface{}, state State) RequestResponder {
	return &request{
		RequestVersion: version,
		RequestId:      id,
		RequestMethod:  method,
		RequestParams:  params,
		requestState:   state,
	}
}

func NewRequestResponder(version string, id interface{}, method string, params interface{}) RequestResponder {
	return NewRequestResponderWithState(version, id, method, params, State{})
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

// The bytes representation of a request will be the JSON encoded value. This
// JSON is expected to be a perfectly valid JSON-RPC request.
func (request *request) Bytes() []byte {
	b, err := json.Marshal(request)
	if err != nil {
		// I don't know what would cause this situation. There is nothing we can
		// do except return an empty string (which would not occur in any
		// successful situation).
		return nil
	}

	return b
}

func newRequestResponderFromJSON(jsonRequest []byte, isPartOfBatch bool, state State) (RequestResponder, interface{}, int, string) {
	var requestMap map[string]interface{}
	err := json.Unmarshal(jsonRequest, &requestMap)
	if err != nil {
		errCode := ParseError

		// The JSON-RPC spec says that for a batch request, any individual
		// requests that would normally throw a ParseError here should be
		// treated as InvalidRequest instead.
		if isPartOfBatch {
			errCode = InvalidRequest
		}

		// It is unlikely that we will have an "id" but we might as well try.
		return nil, requestMap["id"], errCode, ErrorMessageForCode(errCode)
	}

	// Catch some type errors before creating the real request.
	if _, ok := requestMap["jsonrpc"].(string); !ok {
		return nil, requestMap["id"],
			InvalidRequest, "Version (jsonrpc) must be a string."
	}
	if _, ok := requestMap["method"].(string); !ok {
		return nil, requestMap["id"], InvalidRequest, "Method must be a string."
	}

	return NewRequestResponderWithState(
		requestMap["jsonrpc"].(string),
		requestMap["id"],
		requestMap["method"].(string),
		requestMap["params"],
		state,
	), requestMap["id"], Success, ""
}

func NewRequestFromJSON(data []byte) (Request, error) {
	r, _, _, errMessage := newRequestResponderFromJSON(data, false, nil)
	if errMessage != "" {
		return nil, errors.New(errMessage)
	}

	return r, nil
}
