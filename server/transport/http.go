package transport

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/yekhlakov/gojsonrpc/common"
	"github.com/yekhlakov/gojsonrpc/server"
)

// Http Request context
// This extends the common Json-Rpc Request Context
type HttpRequestContext struct {
	HttpRequest *http.Request
	common.RequestContext
	HttpResponse http.ResponseWriter
}

// A Stage for processing an Http Request before or after the Json-Rpc processing
type HttpStage func(context *HttpRequestContext) bool

// This is the actual HTTP transport
type HttpTransport struct {
	Mux              *http.ServeMux
	PreServerStages  []HttpStage
	Endpoints        map[string]*server.JsonRpcServer
	PostServerStages []HttpStage
	logger           *log.Logger
}

// Create a new HTTP transport listening on a given hostName:port
func NewHttpTransport(hostName string) HttpTransport {

	transport := HttpTransport{
		Mux:              http.NewServeMux(),
		PreServerStages:  []HttpStage{},
		Endpoints:        map[string]*server.JsonRpcServer{},
		PostServerStages: []HttpStage{},
		logger:           log.New(ioutil.Discard, "", 0),
	}

	go func() {
		err := http.ListenAndServe(hostName, transport.Mux)
		fmt.Println(err)
	}()

	return transport
}

// Set the logger for the Http Transport
// This logger will be also set for all (current and future) endpoints of this HTTP transport
func (t *HttpTransport) SetLogger(logger *log.Logger) error {
	if logger == nil {
		return fmt.Errorf("nil logger not allowed")
	}
	t.logger = logger

	for k := range t.Endpoints {
		t.Endpoints[k].Logger = logger
	}

	return nil
}

// Process a given pipeline of HttpStage's
// that is sequentially apply each stage to the context until an error is returned or no more stages left
func (hrc *HttpRequestContext) applyPipeline(stages *[]HttpStage) (ok bool) {
	ok = true

	for _, stage := range *stages {
		if ok = stage(hrc); !ok {
			return
		}
	}

	return
}

// Get an endpoint by url
func (t *HttpTransport) GetEndpoint(url string) *server.JsonRpcServer {
	if e, ok := t.Endpoints[url]; ok {
		return e
	}

	return nil
}

// Add an existing JSON-RPC server to a transport at given endpoint URL
func (t *HttpTransport) AddEndpoint(url string, s *server.JsonRpcServer) (*server.JsonRpcServer, error) {

	if s == nil {
		return nil, fmt.Errorf("nil server not allowed")
	}

	// Check if this url is already registered
	if t.GetEndpoint(url) != nil {
		return nil, fmt.Errorf("the url is already registered")
	}

	t.Endpoints[url] = s
	s.Logger = t.logger

	// Register a handler function on the transport for the newly added endpoint
	t.Mux.HandleFunc(url, func(w http.ResponseWriter, r *http.Request) {

		// Create Context
		context := HttpRequestContext{
			HttpRequest:    r,
			HttpResponse:   w,
			RequestContext: common.EmptyRequestContext(),
		}

		var err error
		context.RawRequest, err = ioutil.ReadAll(r.Body)
		if err != nil {
			context.RequestContext.MakeErrorResponse(common.InvalidRequestError)
		} else {
			_ = t.ProcessRequest(s, &context)
		}

		if err != nil && context.RawResponse == nil {
			// Some error should have been set by a pre-processing stage so just regenerate the response
			_ = context.RebuildRawResponse()
		}

		// Write the response
		if _, err = w.Write(context.RawResponse); err != nil {
			// Looks like we can't write to output, so no error will ever be returned
			s.Logger.Println("http response write error", err.Error())
		}
	})

	return s, nil
}

// Process the request, return the result that should be ready to write out
func (t *HttpTransport) ProcessRequest(s *server.JsonRpcServer, hrc *HttpRequestContext) (ok bool) {
	if ok = hrc.applyPipeline(&t.PreServerStages); !ok {
		return
	}

	_ = s.ProcessRawInput(&hrc.RequestContext)

	_ = hrc.applyPipeline(&t.PostServerStages)

	return true
}

func (t *HttpTransport) AddPreServerStage(stage HttpStage) {
	t.PreServerStages = append(t.PreServerStages, stage)
}

func (t *HttpTransport) AddPostServerStage(stage HttpStage) {
	t.PostServerStages = append(t.PostServerStages, stage)
}
