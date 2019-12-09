package server

import (
    "encoding/json"
    "fmt"

    "github.com/yekhlakov/gojsonrpc/common"
)

// Create an empty Request Context
func EmptyRequestContext() RequestContext {
    return RequestContext{
        JsonRpcRequest:  common.Request{},
        JsonRpcResponse: common.Response{},
        RawResponse:     nil,
        RawRequest:      nil,
        Logger:          nil,
        Data:            map[string]interface{}{},
    }
}

func (rc * RequestContext) MakeEmptyResponse() {
    rc.JsonRpcResponse = rc.JsonRpcRequest.MakeResponse(nil, nil)
}

// Create a JSON-RPC Response with a given Error and put it into the Context
func (rc * RequestContext) MakeErrorResponse(e common.Error) {
    rc.JsonRpcResponse = rc.JsonRpcRequest.MakeErrorResponse(e)
}

// Create a JSON-RPC request from Raw Request in this Context
func (rc * RequestContext) ParseRawRequest() (err error) {
    if err = json.Unmarshal(rc.RawRequest, &rc.JsonRpcRequest); err != nil {
        rc.MakeErrorResponse(common.ParseError)
    } else if rc.JsonRpcRequest.JsonRPC != "2.0" || rc.JsonRpcRequest.Method == "" {
        rc.MakeErrorResponse(common.InvalidRequestError)
        err = fmt.Errorf("invalid request")
    }

    return
}

// Create Raw Response from JSON-RPC Response in this Context
func (rc * RequestContext) RebuildRawResponse() (err error) {
    rc.RawResponse, err = json.Marshal(rc.JsonRpcResponse)
    if err != nil {
        rc.Logger.Println("failed to build raw response", err.Error())
        rc.MakeErrorResponse(common.InternalError)
        // Try again (possibly failing once more)
        rc.RawResponse, err = json.Marshal(rc.JsonRpcResponse)
    }

    return
}
