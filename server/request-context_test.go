package server

import (
    "encoding/json"
    "testing"

    "github.com/yekhlakov/gojsonrpc/common"
)

func TestEmptyRequestContext(t *testing.T) {
    rc := EmptyRequestContext()

    if rc.JsonRpcRequest.Method != "" {
        t.Errorf("request context JsonRpcRequest method is not nil")
    }

    if rc.JsonRpcRequest.Params != nil {
        t.Errorf("request context JsonRpcRequest params is not nil")
    }

    if rc.RawRequest != nil {
        t.Errorf("request context raw response is not nil")
    }

    if rc.RawResponse != nil {
        t.Errorf("request context raw response is not nil")
    }

    if rc.Logger != nil {
        t.Errorf("request context logger is not nil")
    }

    if len(rc.Data) != 0 {
        t.Errorf("request context data is not empty")
    }
}

func TestRequestContext_MakeEmptyResponse(t *testing.T) {

    rc := EmptyRequestContext()
    rc.JsonRpcRequest = common.Request{
        JsonRPC: "666",
        Id:      "777",
        Method:  "888",
        Params:  nil,
    }

    rc.MakeEmptyResponse()

    if rc.JsonRpcResponse.JsonRPC != "2.0" {
        t.Errorf("request context JsonRpcResponse json-rpc version is not 2.0")
    }

    if rc.JsonRpcResponse.Id != rc.JsonRpcRequest.Id {
        t.Errorf("request context JsonRpcResponse Id is not equal to JsonRpcRequest Id")
    }

    if rc.JsonRpcResponse.Result != nil {
        t.Errorf("request context JsonRpcResponse is not empty")
    }

    if rc.JsonRpcResponse.Error != nil {
        t.Errorf("request context JsonRpcResponse Error is not empty")
    }
}

func TestRequestContext_MakeErrorResponse(t *testing.T) {
    rc := EmptyRequestContext()
    rc.JsonRpcRequest = common.Request{
        JsonRPC: "666",
        Id:      "777",
        Method:  "888",
        Params:  nil,
    }

    rc.MakeErrorResponse(common.MethodNotFoundError)

    if rc.JsonRpcResponse.JsonRPC != "2.0" {
        t.Errorf("request context JsonRpcResponse json-rpc version is not 2.0")
    }

    if rc.JsonRpcResponse.Id != rc.JsonRpcRequest.Id {
        t.Errorf("request context JsonRpcResponse Id is not equal to JsonRpcRequest Id")
    }

    if rc.JsonRpcResponse.Result != nil {
        t.Errorf("request context JsonRpcResponse is not empty")
    }

    if rc.JsonRpcResponse.Error == nil {
        t.Errorf("request context JsonRpcResponse Error is empty")
    }

    e := common.Error{}
    err := json.Unmarshal(rc.JsonRpcResponse.Error, &e)
    if err != nil {
        t.Errorf("Error was not passed through properly")
        t.Errorf(err.Error())
    } else {
        if e.Code != common.MethodNotFoundError.Code {
            t.Error("Wrong error returned")
        }
    }
}

type TestRequestContext_Params_Struct struct {
    Xxx json.Number `json:"xxx"`
}

func TestRequestContext_ParseRawRequest_Error(t *testing.T) {

    testData := []struct {
        In  string
        Out string
        Err bool
    }{
        {
            `{"qwe"}`,
            `{"code":-32700,"message":"Parse error"}`,
            true,
        },
        {
            `{}`,
            `{"code":-32600,"message":"Invalid request"}`,
            true,
        },
        {
            `{"jsonrpc":"1.0","method":"test","params":{"xxx":666}}`,
            `{"code":-32600,"message":"Invalid request"}`,
            true,
        },
        {
            `{"jsonrpc":"2.0","params":{"xxx":666}}`,
            `{"code":-32600,"message":"Invalid request"}`,
            true,
        },
        {
            `{"jsonrpc":"2.0","method":"test","params":{"xxx":666}}`,
            ``,
            false,
        },
    }

    for k, data := range testData {
        rc := EmptyRequestContext()
        rc.RawRequest = []byte(data.In)
        err := rc.ParseRawRequest()
        if err != nil && !data.Err {
            t.Errorf("%d Error was produced when it shouldn't", k)
        }
        if err == nil && data.Err {
            t.Errorf("%d Error was not produced when it should", k)
        }

        if string(rc.JsonRpcResponse.Error) != data.Out {
            t.Errorf("%d json-rpc error was not generated properly", k)
            t.Errorf(string(rc.JsonRpcResponse.Error))
        }
    }

}

func TestRequestContext_ParseRawRequest(t *testing.T) {

    rc := EmptyRequestContext()
    rc.RawRequest = []byte(`{"jsonrpc":"2.0","method":"test","params":{"xxx":666}}`)
    err := rc.ParseRawRequest()
    if err != nil {
        t.Errorf("Could not unmarshal raw request")
        t.Error(err)
    } else {
        if rc.JsonRpcRequest.JsonRPC != "2.0" {
            t.Errorf("Raw request jsonrpc was not unmarshaled properly")
        }
        if rc.JsonRpcRequest.Method != "test" {
            t.Errorf("Raw request method was not unmarshaled properly")
        }
        if string(rc.JsonRpcRequest.Params) != `{"xxx":666}` {
            t.Errorf("Raw request was not unmarshaled properly")
        }

        var params TestRequestContext_Params_Struct

        err := json.Unmarshal(rc.JsonRpcRequest.Params, &params)
        if err != nil {
            t.Errorf("Could not unmarshal request params")
            t.Error(err)
        } else {
            if params.Xxx != "666" {
                t.Errorf("request params were not unmarshaled properly")
                t.Error(params.Xxx)
            }
        }
    }
}

func TestRequestContext_RebuildRawResponse(t *testing.T) {
    rc := RequestContext{
        JsonRpcResponse: common.Response{
            JsonRPC: "2.0",
            Id:      "test",
            Error:   nil,
            Result:  []byte(`{"xxx":555}`),
        },
    }

    err := rc.RebuildRawResponse()
    if err != nil {
        t.Errorf("Could not unmarshal json rpc response")
        t.Errorf(err.Error())
    } else {
        if string(rc.RawResponse) != `{"jsonrpc":"2.0","id":"test","result":{"xxx":555}}` {
            t.Errorf("Response was not unmarshaled properly")
            t.Errorf(string(rc.RawResponse))
        }
    }
}
