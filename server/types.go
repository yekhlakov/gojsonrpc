package server

import (
    "encoding/json"
    "log"
    "reflect"

    "github.com/yekhlakov/gojsonrpc/common"
)

type RequestContext struct {
    JsonRpcRequest  common.Request
    JsonRpcResponse common.Response
    RawRequest      json.RawMessage
    RawResponse     json.RawMessage
    Logger          *log.Logger
    Data            map[string]interface{}
}

// Generalized processing stage
// It takes a RequestContext and returns a modified context + maybe an error
// If the second return is false, the request won't be processed further
type Stage func(context *RequestContext) bool

// Generalized Handler type
// a Handler is an object with methods that have names like 'namePrefix_jsonRpcMethodName'
// these methods will get invoked by endpoints
type Handler interface{}

// A struct for keeping JSON-RPC method descriptions
type JsonRpcMethod struct {
    Receiver   Handler
    Name       string
    Method     reflect.Method
    ParamsType reflect.Type
    ResultType reflect.Type
}

// A Server for actual handling of requests
type JsonRpcServer struct {
    PreProcessingStages  []Stage
    Methods              map[string]JsonRpcMethod
    PostProcessingStages []Stage
    Logger               *log.Logger
}
