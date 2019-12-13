package server

import (
	"testing"

	"github.com/yekhlakov/gojsonrpc/common"
)

// Testing method invocation
func TestInvokeMethod(t *testing.T) {
	testData := []struct {
		Name    string
		Handler Handler
		In      string
		Out     string
	}{
		{
			"empty",
			test_EmptyHandler{},
			`{"jsonrpc":"2.0","id":"test","method":"whatever","params":{"name":"lol"}}`,
			`{"jsonrpc":"2.0","id":"test","result":{}}`,
		},
		{
			"const",
			test_ConstHandler{},
			`{"jsonrpc":"2.0","id":"test","method":"const","params":{"name":"lol"}}`,
			`{"jsonrpc":"2.0","id":"test","result":{"value":"test"}}`,
		},
		{
			"pass",
			test_PassHandler{},
			`{"jsonrpc":"2.0","id":"test","method":"pass","params":{"name":"lol"}}`,
			`{"jsonrpc":"2.0","id":"test","result":{"value":"lol"}}`,
		},
		{
			"error",
			test_ErrorHandler{},
			`{"jsonrpc":"2.0","id":"test","method":"pass","params":{"name":"lol"}}`,
			`{"jsonrpc":"2.0","id":"test","error":{"code":666,"message":"error"}}`,
		},
		{
			"invalid params",
			test_PassHandler{},
			`{"jsonrpc":"2.0","id":"test","method":"pass","params":{"name":["lol"]}}`,
			`{"jsonrpc":"2.0","id":"test","error":{"code":-32602,"message":"Invalid params"}}`,
		},
		{
			"bad method result unmarshaling",
			test_WrongHandler{},
			`{"jsonrpc":"2.0","id":"test","method":"wrong","params":{"name":"lol"}}`,
			`{"jsonrpc":"2.0","id":"test","error":{"code":-32603,"message":"Internal error"}}`,
		},
	}

	for k, data := range testData {
		m := ExtractMethods(data.Handler, "Handle_")

		rc := &RequestContext{
			JsonRpcRequest:  common.Request{},
			JsonRpcResponse: common.Response{},
			RawRequest:      []byte(data.In),
			RawResponse:     nil,
			Logger:          nil,
			Data:            nil,
		}
		_ = rc.ParseRawRequest()

		_ = InvokeMethod(rc, m[0])
		_ = rc.RebuildRawResponse()

		if string(rc.RawResponse) != data.Out {
			t.Errorf("%d %s : Method invocation returned wrong results", k, data.Name)
			t.Error(string(rc.RawResponse))
		}
	}
}
