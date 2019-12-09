package common

import (
    "encoding/json"
)

// General JSON-RPC request
type Request struct {
    JsonRPC string          `json:"jsonrpc"`
    Id      string          `json:"id"`
    Method  string          `json:"method"`
    Params  json.RawMessage `json:"params"`
}

// General JSON-RPC response
type Response struct {
    JsonRPC string          `json:"jsonrpc"`
    Id      string          `json:"id,omitempty"`
    Error   json.RawMessage `json:"error,omitempty"`
    Result  json.RawMessage `json:"result,omitempty"`
}

// General JSON-RPC error payload
type ErrorData struct {
    Errors []Error `json:"errors,omitempty"`
}

// General JSON-RPC response error
type Error struct {
    Code    json.Number     `json:"code"`
    Message string          `json:"message"`
    Data    json.RawMessage `json:"data,omitempty"`
}
