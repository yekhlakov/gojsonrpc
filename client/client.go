package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"

	"github.com/yekhlakov/gojsonrpc/client/transport"
	"github.com/yekhlakov/gojsonrpc/common"
)

// Create a new empty Client
func New() *Client {
	return &Client{
		T:      &transport.Discard{},
		logger: log.New(ioutil.Discard, "", 0),
	}
}

// Set client transport
func (c *Client) SetTransport(t Transport) error {
	if t == nil {
		return fmt.Errorf("nil transport not allowed")
	}

	c.T = t
	return nil
}

// Set client logger. The same logger will be set for client's transport
func (c *Client) SetLogger(l *log.Logger) error {
	if l == nil {
		return fmt.Errorf("nil logger not allowed")
	}
	c.logger = l

	return c.T.SetLogger(l)
}

// Do the request
func (c *Client) PerformRequest(rc *common.RequestContext) error {
	if err := rc.RebuildRawRequest(); err != nil {
		return err
	}
	if err := c.T.PerformRequest(rc); err != nil {
		return err
	}
	if err := rc.ParseRawResponse(); err != nil {
		return err
	}

	return nil
}

// Create a new Json-Rpc Request
func NewJsonRpcRequest(method string, params interface{}) (common.Request, error) {
	r := common.Request{
		JsonRPC: "2.0",
		Id:      fmt.Sprintf("%8.8x%8.8x", rand.Uint64(), rand.Uint64()),
		Method:  method,
		Params:  nil,
	}

	p, err := json.Marshal(params)
	if err == nil {
		r.Params = p
	}

	return r, err
}

// Create a new Request Context for the given Client
func (c *Client) NewRequestContext(method string, params interface{}) (rc common.RequestContext, err error) {
	rc = common.EmptyRequestContext()
	rc.Logger = c.logger
	rc.JsonRpcRequest, err = NewJsonRpcRequest(method, params)
	return
}

// The main entry point for request processing
func (c *Client) Request(method string, params interface{}) (response common.Response, err error) {
	rc, err := c.NewRequestContext(method, params)
	if err != nil {
		return
	}

	err = c.PerformRequest(&rc)
	if err != nil {
		return
	}

	return rc.JsonRpcResponse, nil
}
