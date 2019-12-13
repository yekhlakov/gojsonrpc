package server

import (
	"github.com/yekhlakov/gojsonrpc/common"
)

// InvokeMethod the method with given request and response
// A lot of reflection magic is going on here
func InvokeMethod(rc *RequestContext, m JsonRpcMethod) (err error) {

	// Prepare an empty Response
	rc.MakeEmptyResponse()

	// Bind params for the method
	boundParams, err := m.BindParams(rc.JsonRpcRequest.Params)
	if err != nil {
		rc.MakeErrorResponse(common.InvalidParamsError)
		return
	}

	// Call the method and get back the results which is an array of Values
	results := m.Method.Func.Call(boundParams)

	// Extract the Error, and if it is not nil, put it into the Response and return
	rawError, err := m.ExtractError(results)
	if err != nil {
		rc.MakeErrorResponse(common.InternalError)
		return
	} else if rawError != nil {
		rc.JsonRpcResponse.Error = rawError
		rc.JsonRpcResponse.Result = nil

		return
	}

	// Extract the Result, encode it into raw data and put into the Response
	rawResult, err := m.ExtractResult(results)
	if err != nil {
		rc.MakeErrorResponse(common.InternalError)
	} else {
		rc.JsonRpcResponse.Error = nil
		rc.JsonRpcResponse.Result = rawResult
	}

	return
}
