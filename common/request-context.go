package common

import (
	"encoding/json"
	"fmt"
)

// Create an empty Request Context
func EmptyRequestContext() RequestContext {
	return RequestContext{
		JsonRpcRequest:  Request{},
		JsonRpcResponse: Response{},
		RawResponse:     nil,
		RawRequest:      nil,
		Logger:          nil,
		Data:            map[string]interface{}{},
	}
}

func (rc *RequestContext) MakeEmptyResponse() {
	rc.JsonRpcResponse = rc.JsonRpcRequest.MakeResponse(nil, nil)
}

// Create a JSON-RPC Response with a given Error and put it into the Context
func (rc *RequestContext) MakeErrorResponse(e Error) {
	rc.JsonRpcResponse = rc.JsonRpcRequest.MakeErrorResponse(e)
}

// Create a JSON-RPC request from Raw Request in this Context
func (rc *RequestContext) ParseRawRequest() (err error) {
	if err = json.Unmarshal(rc.RawRequest, &rc.JsonRpcRequest); err != nil {
		rc.MakeErrorResponse(ParseError)
	} else if rc.JsonRpcRequest.JsonRPC != "2.0" || rc.JsonRpcRequest.Method == "" {
		rc.MakeErrorResponse(InvalidRequestError)
		err = fmt.Errorf("invalid request")
	}

	return
}

// Create Raw Request from JSON-RPC Request in this Context
func (rc *RequestContext) RebuildRawRequest() (err error) {
	rc.RawRequest, err = json.Marshal(rc.JsonRpcRequest)
	if err != nil {
		rc.Logger.Println("failed to build raw request", err.Error())
	}
	return
}

// Create Raw Response from JSON-RPC Response in this Context
func (rc *RequestContext) RebuildRawResponse() (err error) {
	rc.RawResponse, err = json.Marshal(rc.JsonRpcResponse)
	if err != nil {
		rc.Logger.Println("failed to build raw response", err.Error())
		rc.MakeErrorResponse(InternalError)
		// Try again (possibly failing once more)
		rc.RawResponse, err = json.Marshal(rc.JsonRpcResponse)
	}

	return
}

// Create JSON-RPC Response from Raw Response in this Context
func (rc *RequestContext) ParseRawResponse() (err error) {
	return json.Unmarshal(rc.RawResponse, &rc.JsonRpcResponse)
}

// Apply a processing pipeline to the context
func (rc *RequestContext) ApplyPipeline(stages *[]Stage) (ok bool) {
	ok = true

	for _, stage := range *stages {
		if ok = stage(rc); !ok {
			return
		}
	}

	return
}
