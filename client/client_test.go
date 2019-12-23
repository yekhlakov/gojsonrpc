package client

import (
	"encoding/json"
	"log"
	"testing"

	"github.com/yekhlakov/gojsonrpc/client/transport"
	"github.com/yekhlakov/gojsonrpc/common"
	"github.com/yekhlakov/gojsonrpc/server"
)

func TestNew(t *testing.T) {
	c := New()

	if _, ok := c.T.(*transport.Discard); !ok {
		t.Errorf("wrong default transport")
	}

	if c.logger == nil {
		t.Errorf("nil default logger")
	}
}

type test_Transport struct {
	transport.Logged
}

// No request handling
func (t *test_Transport) PerformRequest(rc *common.RequestContext) error {
	return nil
}

// No pre-processing stages
func (t *test_Transport) AddPreProcessingStage(stage common.Stage) {
}

// No post-processing stages
func (t *test_Transport) AddPostProcessingStage(stage common.Stage) {
}

func TestClient_SetTransport(t *testing.T) {
	c := New()

	e := c.SetTransport(nil)
	if e == nil {
		t.Errorf("nil transport was attached")
	}

	tp := test_Transport{}

	e = c.SetTransport(&tp)
	if e != nil {
		t.Errorf("could not attach a transport")
	} else if c.T != &tp {
		t.Errorf("transport was not attached")
	}
}

func TestClient_SetLogger(t *testing.T) {
	c := New()

	e := c.SetLogger(nil)
	if e == nil {
		t.Errorf("nil logger was attached")
	}

	tp := test_Transport{}
	_ = c.SetTransport(&tp)

	lg := log.Logger{}

	e = c.SetLogger(&lg)

	if e != nil {
		t.Errorf("could not attach a logger")
	} else if c.logger != &lg {
		t.Errorf("logger was not attached")
	}
	// TODO: test logger propagation
}

func TestClient_PerformRequest_Error1(t *testing.T) {
	c := New()

	rc, _ := c.NewRequestContext("test", struct{}{})
	rc.JsonRpcRequest.Params = []byte("LOL") // bad json

	err := c.PerformRequest(&rc)
	if err == nil {
		t.Errorf("bad request was processed")
	}
}

func TestClient_PerformRequest_Error2(t *testing.T) {
	rc := common.EmptyRequestContext()
	rc.RawRequest = []byte(`{"jsonrpc":"2.0","id":"test","method":"pass"}`)

	c := New()
	_ = c.SetTransport(&transport.Error{
		ErrorMessage: "KEK",
	})

	err := c.PerformRequest(&rc)

	if err == nil {
		t.Errorf("error was not produced")
	} else if err.Error() != "KEK" {
		t.Errorf("wrong error was produced: %s", err.Error())
	}
}

func TestClient_PerformRequest_Error3(t *testing.T) {
	rc := common.EmptyRequestContext()
	rc.RawRequest = []byte(`{"jsonrpc":"2.0","id":"test","method":"pass"}`)

	c := New()
	_ = c.SetTransport(&transport.Error{
		RawResponse: []byte("KEK"), // bad json
	})

	err := c.PerformRequest(&rc)

	if err == nil {
		t.Errorf("error was not produced")
	}
}

type test_Client struct {
	C *Client
}

// Method params struct
type test_ClientPassParams struct {
	Name string `json:"name"`
}

// Method response struct
type test_ClientPassResult struct {
	Value string `json:"value"`
}

// Method signature
// Methods receives a json-marshal-able struct holding method params
// Methods returns:
//  - a struct holding the result,
//  - a JSON-RPC error received from the server,
//  - and an error probably generated locally
func (c test_Client) Pass(params test_ClientPassParams) (
	result test_ClientPassResult,
	error common.Error,
	err error,
) {

	response, err := c.C.Request("pass", params)
	if err != nil {
		return
	}

	if response.Error != nil {
		err = json.Unmarshal(response.Error, &error)
		return
	}

	err = json.Unmarshal(response.Result, &result)
	return
}

func TestClient_Request(t *testing.T) {
	jsonRpcClient := New()

	tr := transport.Local{
		Server: &server.JsonRpcServer{},
	}

	_ = jsonRpcClient.SetTransport(&tr)

	testClient := test_Client{jsonRpcClient}

	_, e, err := testClient.Pass(test_ClientPassParams{"qwer"})

	if err != nil {
		t.Errorf("got error while processing request")
	} else if e.Code != common.MethodNotFoundError.Code {
		t.Errorf("got wrong error from the server %s, %s", e.Code, e.Message)
	}
}
