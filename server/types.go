package server

import (
	"log"
	"reflect"

	"github.com/yekhlakov/gojsonrpc/common"
)

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
	PreProcessingStages  []common.Stage
	Methods              map[string]JsonRpcMethod
	PostProcessingStages []common.Stage
	Logger               *log.Logger
}
