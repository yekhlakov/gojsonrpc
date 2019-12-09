package common

import (
    "encoding/json"
    "fmt"
    "math/rand"
    "testing"
)

func TestRequest_MakeErrorResponse(t *testing.T) {
    testData:= []struct {
        Id string
        Method string
        Error Error
        Expect json.RawMessage
    }{
        {
            fmt.Sprint(rand.Uint64()),
            fmt.Sprint(rand.Uint64()),
            InvalidParamsError,
            []byte(fmt.Sprintf(`{"code":%s,"message":"%s"}`, InvalidParamsError.Code, InvalidParamsError.Message)),
        },
        {
            fmt.Sprint(rand.Uint64()),
            fmt.Sprint(rand.Uint64()),
            MethodNotFoundError,
            []byte(fmt.Sprintf(`{"code":%s,"message":"%s"}`, MethodNotFoundError.Code, MethodNotFoundError.Message)),
        },
    }

    for k, input := range testData {
        request := Request{
            JsonRPC: "2.0",
            Id:      input.Id,
            Method:  input.Method,
            Params:  nil,
        }

        response := request.MakeErrorResponse(input.Error)

        if response.Id != request.Id {
            t.Errorf("%d Id was not passed through properly", k)
        }

        if response.Result != nil {
            t.Errorf("%d Result should be empty", k)
        }

        if response.Error == nil {
            t.Errorf("%d Error should not be empty", k)
        } else {
            if string(response.Error) != string(input.Expect) {
                t.Errorf("%d Error was not passed through properly", k)
            }
        }
    }
}

func TestRequest_MakeResponse(t *testing.T) {

    testData := []struct {
        Id     string
        Method string
        Result json.RawMessage
        Error  json.RawMessage
    }{
        {
            fmt.Sprint(rand.Uint64()),
            fmt.Sprint(rand.Uint64()),
            nil,
            nil,
        },
        {
            fmt.Sprint(rand.Uint64()),
            fmt.Sprint(rand.Uint64()),
            []byte(fmt.Sprint(rand.Uint64())),
            nil,
        },
        {
            fmt.Sprint(rand.Uint64()),
            fmt.Sprint(rand.Uint64()),
            nil,
            []byte(fmt.Sprint(rand.Uint64())),
        },
        {
            fmt.Sprint(rand.Uint64()),
            fmt.Sprint(rand.Uint64()),
            []byte(fmt.Sprint(rand.Uint64())),
            []byte(fmt.Sprint(rand.Uint64())),
        },    }

    for k, input := range testData {
        request := Request{
            JsonRPC: "2.0",
            Id:      input.Id,
            Method:  input.Method,
            Params:  nil,
        }

        response := request.MakeResponse(input.Result, input.Error)

        if response.Id != request.Id {
            t.Errorf("%d Id was not passed through properly", k)
        }

        if input.Error == nil && string(response.Result) != string(input.Result) {
            t.Errorf("%d Result was not passed through properly", k)
        }

        if input.Error != nil && string(response.Error) != string(input.Error) {
            t.Errorf("%d Error was not passed through properly", k)
        }
    }

}
