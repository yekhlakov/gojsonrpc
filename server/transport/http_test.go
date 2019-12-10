package transport

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"testing"
	"time"

	"github.com/yekhlakov/gojsonrpc/server"
)

func TestHttpTransport_SetLogger(t *testing.T) {
	transport := NewHttpTransport("localhost:56666")

	e, _ := transport.AddEndpoint("/lol", &server.JsonRpcServer{})

	if transport.SetLogger(nil) == nil {
		t.Errorf("Nil logger was accepted")
	}

	l := log.New(ioutil.Discard, "lol", 0)

	if transport.SetLogger(l) != nil {
		t.Errorf("Logger was not set")
	} else {
		if transport.logger != l {
			t.Errorf("Logger was set wrongly")
		}
		if e.Logger != l {
			t.Errorf("Logger was not set for endpoints")
		}
	}
}

func TestHttpTransport_AddEndpoint(t *testing.T) {
	transport := NewHttpTransport("localhost:56666")

	_, e := transport.AddEndpoint("/lol", nil)
	if e == nil {
		t.Errorf("Empty endpoint was accepted")
	}
}

func TestHttpTransport_AddEndpoint2(t *testing.T) {
	transport := NewHttpTransport("localhost:56666")
	server1 := server.JsonRpcServer{}

	s, e := transport.AddEndpoint("/lol", &server1)
	if e != nil {
		t.Errorf("Could not add endpoint")
	} else if s != &server1 {
		t.Errorf("Endpoint was added wrongly")
	}
}

func TestHttpTransport_AddEndpoint3(t *testing.T) {
	transport := NewHttpTransport("localhost:56666")
	server1 := server.JsonRpcServer{}

	_, _ = transport.AddEndpoint("/lol", &server1)
	_, e := transport.AddEndpoint("/lol", &server1)
	if e == nil {
		t.Errorf("Duplicate endpoint was accepted")
	}
}

func TestHttpTransport_GetEndpoint(t *testing.T) {
	transport := NewHttpTransport("localhost:56666")

	server1 := server.JsonRpcServer{}
	server2 := server.JsonRpcServer{}

	transport.Endpoints["lol"] = &server1
	transport.Endpoints["kek"] = &server2

	if transport.GetEndpoint("lol") != &server1 {
		t.Errorf("Endpoint was not returned")
	}

	if transport.GetEndpoint("kek") != &server2 {
		t.Errorf("Endpoint was not returned")
	}

	if transport.GetEndpoint("cheburek") != nil {
		t.Errorf("Endpoint was returned when it shouldn't be")
	}
}

func TestHttpRequestContext_ApplyPipeline(t *testing.T) {
	context := HttpRequestContext{
		RequestContext: server.EmptyRequestContext(),
	}

	stages := []HttpStage{
		func(context *HttpRequestContext) bool {
			context.Data["lol"] = "kek"
			return true
		},

		// This should be called last
		func(context *HttpRequestContext) bool {
			context.Data["lol"] = "cheburek"
			return false
		},

		// This should not be called
		func(context *HttpRequestContext) bool {
			context.Data["lol"] = "azaza"
			return true
		},
	}

	context.applyPipeline(&stages)

	v, ok := context.Data["lol"]
	if !ok {
		t.Errorf("pipeline not applied")
	}

	if v != "cheburek" {
		t.Errorf("pipeline was not applied correctly")
	}
}

func TestHttpRequestContext_ApplyPipeline2(t *testing.T) {
	context := HttpRequestContext{
		RequestContext: server.EmptyRequestContext(),
	}

	stages := []HttpStage{
		func(context *HttpRequestContext) bool {
			context.Data["lol"] = "kek"
			return true
		},

		func(context *HttpRequestContext) bool {
			context.Data["lol"] = "cheburek"
			return true
		},

		// This should be called
		func(context *HttpRequestContext) bool {
			context.Data["lol"] = "azaza"
			return true
		},
	}

	context.applyPipeline(&stages)

	v, ok := context.Data["lol"]
	if !ok {
		t.Errorf("pipeline not applied")
	}

	if v != "azaza" {
		t.Errorf("pipeline was not applied correctly")
	}
}

func TestHttpTransport_ProcessRequest(t *testing.T) {
	transport := NewHttpTransport("localhost:56666")
	server1 := server.JsonRpcServer{}
	_, _ = transport.AddEndpoint("/lol", &server1)
	context := HttpRequestContext{
		RequestContext: server.EmptyRequestContext(),
	}
	context.RawRequest = []byte(`{}`)

	if !transport.ProcessRequest(&server1, &context) {
		t.Errorf("Error processing request")
	}

	transport.AddPreServerStage(func(context *HttpRequestContext) bool {
		return false
	})

	if transport.ProcessRequest(&server1, &context) {
		t.Errorf("Error was not generated")
	}
}

func TestHttpTransport(t *testing.T) {
	transport := NewHttpTransport("localhost:56666")
	server1 := server.JsonRpcServer{}
	_, _ = transport.AddEndpoint("/lol", &server1)

	time.Sleep(100 * time.Millisecond)

	r, err := http.Post(
		"http://localhost:56666/lol",
		"application/json",
		bytes.NewReader([]byte(`{"jsonrpc":"2.0","method":"kek"}`)),
	)

	if err == nil {
		o, e := ioutil.ReadAll(r.Body)

		if e != nil {
			t.Errorf("Error reading response body")
		} else if string(o) != `{"jsonrpc":"2.0","error":{"code":-32601,"message":"Method not found"}}` {
			t.Errorf("Wrong response received")
			t.Errorf(string(o))
		}
	} else {
		t.Errorf("Got http post error %s", err.Error())
	}
}
