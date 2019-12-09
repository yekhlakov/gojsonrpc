package common

import (
    "encoding/json"
)

var ParseError = Error{
    Code:    "-32700",
    Message: "Parse error",
}

var InvalidRequestError = Error{
    Code:    "-32600",
    Message: "Invalid request",
}

var MethodNotFoundError = Error{
    Code:    "-32601",
    Message: "Method not found",
}

var InvalidParamsError = Error{
    Code:    "-32602",
    Message: "Invalid params",
}

var InternalError = Error{
    Code:    "-32603",
    Message: "Internal error",
}

// Create a Response to the given Request and put the given Error in it
func (rq Request) MakeErrorResponse(e Error) Response {
    errorData, err := json.Marshal(e)
    if err != nil {
        errorData, _ = json.Marshal(nil)
    }
    return rq.MakeResponse(nil, errorData)
}

// Create a Response for the request
func (rq Request) MakeResponse(resultData json.RawMessage, errorData json.RawMessage) Response {
    return Response{
        JsonRPC: "2.0",
        Id:      rq.Id,
        Error:   errorData,
        Result:  resultData,
    }
}
