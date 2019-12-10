package server

import (
	"encoding/json"
	"testing"
)

func TestJsonRpcServer_AddHandler(t *testing.T) {

	s := NewServer()

	s.AddHandler(test_PassHandler{}, "Handle_")

	if len(s.Methods) != 1 {
		t.Errorf("Methods were not extracted from the first handler")
	}

	s.AddHandler(test_PassHandler{}, "Handle_")

	if len(s.Methods) != 1 {
		t.Errorf("Methods were extracted from the first handler twice")
	}

	s.AddHandler(test_ErrorHandler{}, "Handle_")

	if len(s.Methods) != 2 {
		t.Errorf("Methods were not extracted from the second handler")
	}
}

func TestJsonRpcServer_GetMethod(t *testing.T) {
	s := NewServer()

	s.AddHandler(test_PassHandler{}, "Handle_")

	if len(s.Methods) != 1 {
		t.Errorf("Methods were not extracted from the first handler")
	}

	m, ok := s.GetMethod("pass")

	if !ok {
		t.Errorf("Method was not found")
	} else {
		if m.Name != "pass" {
			t.Errorf("wrong method was extracted")
		}
	}
}

func TestJsonRpcServer_ProcessRawRequest(t *testing.T) {

	s := NewServer()
	s.AddHandler(test_EmptyHandler{}, "Handle_")
	rc := EmptyRequestContext()
	rc.RawRequest = []byte(`{"badjson":`)
	err := s.ProcessRawRequest(&rc)
	if err == nil {
		t.Errorf("Bad json parsing did not generate an error")
	}

	testData := []struct {
		Handler Handler
		In      string
		Out     string
	}{
		{
			test_EmptyHandler{},
			`{"jsonrpc":"2.0","id":"test","method":"empty","params":{"name":"lol"}}`,
			`{"jsonrpc":"2.0","id":"test","result":{}}`,
		},
		{
			test_ConstHandler{},
			`{"jsonrpc":"2.0","id":"test","method":"const","params":{"name":"lol"}}`,
			`{"jsonrpc":"2.0","id":"test","result":{"value":"test"}}`,
		},
		{
			test_PassHandler{},
			`{"jsonrpc":"2.0","id":"test","method":"pass","params":{"name":"lol"}}`,
			`{"jsonrpc":"2.0","id":"test","result":{"value":"lol"}}`,
		},
		{
			test_ErrorHandler{},
			`{"jsonrpc":"2.0","id":"test","method":"error","params":{"name":"lol"}}`,
			`{"jsonrpc":"2.0","id":"test","error":{"code":666,"message":"error"}}`,
		},
		{
			test_PassHandler{},
			`{"jsonrpc":"2.0","id":"test","method":"lol","params":{"name":"lol"}}`,
			`{"jsonrpc":"2.0","id":"test","error":{"code":-32601,"message":"Method not found"}}`,
		},
	}

	for k, data := range testData {
		s := NewServer()
		s.AddHandler(data.Handler, "Handle_")

		rc := EmptyRequestContext()
		rc.RawRequest = []byte(data.In)
		err := rc.ParseRawRequest()
		if err != nil {
			t.Errorf("%d Request parse failed", k)
		}
		_ = s.ProcessRawRequest(&rc)

		if string(rc.RawRequest) != data.In {
			t.Errorf("%d Request context was not passed through properly", k)
		}
		if string(rc.RawResponse) != data.Out {
			t.Errorf("%d Request was not processed properly", k)
		}
	}
}

func TestJsonRpcServer_ProcessRawBatch(t *testing.T) {
	testData := []struct {
		Handler Handler
		In      []string
		Out     string
	}{
		{
			test_EmptyHandler{},
			[]string{},
			`[]`,
		},
		{
			test_EmptyHandler{},
			[]string{`{"jsonrpc":"2.0","id":"test","method":"empty","params":{"name":"lol"}}`},
			`[{"jsonrpc":"2.0","id":"test","result":{}}]`,
		},
		{
			test_EmptyHandler{},
			[]string{
				`{"jsonrpc":"2.0","id":"test","method":"empty","params":{"name":"lol"}}`,
				`{"badjson`,
			},
			`[{"jsonrpc":"2.0","id":"test","result":{}},{"jsonrpc":"2.0","error":{"code":-32700,"message":"Parse error"}}]`,
		},
		{
			test_PassHandler{},
			[]string{
				`{"jsonrpc":"2.0","id":"t1","method":"pass","params":{"name":"lol"}}`,
				`{"jsonrpc":"2.0","id":"t2","method":"nope","params":{"name":"kek"}}`,
			},
			`[{"jsonrpc":"2.0","id":"t1","result":{"value":"lol"}},{"jsonrpc":"2.0","id":"t2","error":{"code":-32601,"message":"Method not found"}}]`,
		},
	}

	for k, data := range testData {
		s := NewServer()
		s.AddHandler(data.Handler, "Handle_")

		rc := EmptyRequestContext()

		batch := make([]json.RawMessage, len(data.In))
		for i, v := range data.In {
			batch[i] = []byte(v)
		}

		err := s.ProcessRawBatch(batch, &rc)
		if err != nil {
			t.Errorf("%d Batch processing generated an error", k)
			t.Errorf(err.Error())
		} else if string(rc.RawResponse) != data.Out {
			t.Errorf("%d Request was not processed properly", k)
			t.Errorf(string(rc.RawResponse))
		}
	}
}

func TestJsonRpcServer_ProcessRawInput(t *testing.T) {
	testData := []struct {
		Name    string
		Handler Handler
		In      string
		Out     string
	}{
		{
			"Empty Batch",
			test_EmptyHandler{},
			`[]`,
			`{"jsonrpc":"2.0","error":{"code":-32600,"message":"Invalid request"}}`,
		},
		{
			"Leading whitespace",
			test_EmptyHandler{},
			"    \t\r\n  []",
			`{"jsonrpc":"2.0","error":{"code":-32600,"message":"Invalid request"}}`, // todo: remake
		},
		{
			"Bad input",
			test_EmptyHandler{},
			"666",
			`{"jsonrpc":"2.0","error":{"code":-32600,"message":"Invalid request"}}`,
		},
		{
			"Bad json",
			test_EmptyHandler{},
			"[...]",
			`{"jsonrpc":"2.0","error":{"code":-32700,"message":"Parse error"}}`,
		},
		{
			"Ok batch",
			test_EmptyHandler{},
			`[{"jsonrpc":"2.0","id":"test","method":"empty","params":{"name":"lol"}}]`,
			`[{"jsonrpc":"2.0","id":"test","result":{}}]`,
		},
		{
			"Invalid request in a batch",
			test_EmptyHandler{},
			`[{"jsonrpc":"2.0","id":"test","method":"empty","params":{"name":"lol"}},{"ololo":"trololo"}]`,
			`[{"jsonrpc":"2.0","id":"test","result":{}},{"jsonrpc":"2.0","error":{"code":-32600,"message":"Invalid request"}}]`,
		},
		{
			"Bad request",
			test_EmptyHandler{},
			"{.}",
			`{"jsonrpc":"2.0","error":{"code":-32700,"message":"Parse error"}}`,
		},
		{
			"Ok request",
			test_EmptyHandler{},
			`{"jsonrpc":"2.0","id":"test","method":"empty","params":{"name":"lol"}}`,
			`{"jsonrpc":"2.0","id":"test","result":{}}`,
		},
	}

	for k, data := range testData {
		s := NewServer()
		s.AddHandler(data.Handler, "Handle_")

		rc := EmptyRequestContext()
		rc.RawRequest = []byte(data.In)
		_ = s.ProcessRawInput(&rc)

		if string(rc.RawResponse) != data.Out {
			t.Errorf("%d '%s': input was not processed properly", k, data.Name)
			t.Errorf("expected %s", data.Out)
			t.Errorf("received %s", string(rc.RawResponse))
		}
	}
}

func TestRequestContext_ApplyPipeline(t *testing.T) {
	rc := EmptyRequestContext()

	stages := []Stage{
		func(context *RequestContext) bool {
			context.Data["lol"] = "kek"
			return true
		},

		// This should be called last
		func(context *RequestContext) bool {
			context.Data["lol"] = "cheburek"
			return false
		},

		// This should not be called
		func(context *RequestContext) bool {
			context.Data["lol"] = "azaza"
			return true
		},
	}

	rc.applyPipeline(&stages)

	v, ok := rc.Data["lol"]
	if !ok {
		t.Errorf("pipeline not applied")
	}

	if v != "cheburek" {
		t.Errorf("pipeline was not applied correctly")
	}
}

func TestRequestContext_ApplyPipeline2(t *testing.T) {
	rc := EmptyRequestContext()

	stages := []Stage{
		func(context *RequestContext) bool {
			context.Data["lol"] = "kek"
			return true
		},

		func(context *RequestContext) bool {
			context.Data["lol"] = "cheburek"
			return true
		},

		// This should be called too
		func(context *RequestContext) bool {
			context.Data["lol"] = "azaza"
			return true
		},
	}

	rc.applyPipeline(&stages)

	v, ok := rc.Data["lol"]
	if !ok {
		t.Errorf("pipeline not applied")
	}

	if v != "azaza" {
		t.Errorf("pipeline was not applied correctly")
	}
}
