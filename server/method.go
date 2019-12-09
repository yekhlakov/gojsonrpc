package server

import (
    "encoding/json"
    "fmt"
    "reflect"
    "strings"

    "github.com/yekhlakov/gojsonrpc/common"
)

// Extract methods from the handler using the method name prefix
func ExtractMethods(handler Handler, methodNamePrefix string) (r []JsonRpcMethod) {
    t := reflect.TypeOf(handler)

    for i := 0; i < t.NumMethod(); i++ {
        m := t.Method(i)
        if !strings.HasPrefix(m.Name, methodNamePrefix) {
            continue
        }

        // Should have exactly 2 input parameters: *receiver and *params
        // TODO: add support for Array-Params requests
        if m.Type.NumIn() != 2 {
            continue
        }

        // Should return exactly 3 parameters: result, jsonrpc error, and go error
        if m.Type.NumOut() != 3 {
            continue
        }

        // Type-checks out params
        if !m.Type.Out(1).ConvertibleTo(reflect.TypeOf(common.Error{})) {
            continue
        }
        if m.Type.Out(2).Name() != "error" {
            continue
        }

        description := JsonRpcMethod{
            Receiver:   handler,
            Name:       strings.TrimPrefix(m.Name, methodNamePrefix),
            Method:     m,
            ParamsType: m.Type.In(1),
            ResultType: m.Type.Out(0),
        }

        r = append(r, description)
    }
    return r
}

// Extract the Params from a Json-Rpc Request and convert them into types required by the Method Handler
func (m JsonRpcMethod) BindParams(params json.RawMessage) ([]reflect.Value, error) {
    // Try to create an object of the type that is required by method's params
    // TODO: add support for Array-Params requests
    paramsType := m.ParamsType
    paramsObject := reflect.New(paramsType)
    err := json.Unmarshal(params, paramsObject.Interface())
    if err != nil {
        return []reflect.Value{}, err
    }

    return []reflect.Value{reflect.ValueOf(m.Receiver), paramsObject.Elem()}, nil
}

// Extract an Error from Values returned from a Method invocation
func (m JsonRpcMethod) ExtractError(results []reflect.Value) (r []byte, err error) {

    if !results[1].Type().ConvertibleTo(reflect.TypeOf(common.Error{})) {
        return nil, fmt.Errorf("not an error")
    }

    ePtr := results[1].Interface().(common.Error)

    if ePtr.Code == "" && ePtr.Message == "" && ePtr.Data == nil {
        return nil, nil
    }

    r, err = json.Marshal(ePtr)
    if err != nil {
        r = nil
    }

    return
}

// Extract a Result from Values returned from a Method invocation
func (m JsonRpcMethod) ExtractResult(results []reflect.Value) (r []byte, err error) {

    if !results[0].Type().ConvertibleTo(m.ResultType) {
        return nil, fmt.Errorf("not a valid result")
    }

    r, err = json.Marshal(results[0].Convert(m.ResultType).Interface())

    if err != nil {
        r = nil
    }

    return
}

// InvokeMethod the method with given request and response
// A lot of reflection magic is going on here
func (rc *RequestContext) InvokeMethod(m JsonRpcMethod) (err error) {

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
