package server

import (
	"encoding/json"

	"github.com/yekhlakov/gojsonrpc/common"
)

// An JsonRpcServer does actual processing

// Create a new JsonRpcServer
func NewServer() *JsonRpcServer {
	return &JsonRpcServer{
		Methods: make(map[string]JsonRpcMethod),
	}
}

// Add a handler (that is effectively a collection of methods)
func (e *JsonRpcServer) AddHandler(handler Handler, methodNamePrefix string) {
	methods := ExtractMethods(handler, methodNamePrefix)

	for _, method := range methods {
		e.Methods[method.Name] = method
	}
}

// Get a method from the server
func (e *JsonRpcServer) GetMethod(name string) (method JsonRpcMethod, ok bool) {
	method, ok = e.Methods[name]
	return
}

// Get RAW request (probably a batch), return RAW response
func (e *JsonRpcServer) ProcessRawInput(context *common.RequestContext) (err error) {

	for _, b := range context.RawRequest {
		// skip initial whitespace
		if b == 9 || b == 32 || b == 10 || b == 13 {
			continue
		}

		if b == '[' {
			// Batch

			batch := struct {
				Batch []json.RawMessage `json:"batch,omitempty"`
			}{}

			rawBatch := []byte(`{"batch":`)
			rawBatch = append(rawBatch, context.RawRequest...)
			rawBatch = append(rawBatch, []byte(`}`)...)

			if err = json.Unmarshal(rawBatch, &batch); err != nil {
				context.MakeErrorResponse(common.ParseError)
				break
			}

			if len(batch.Batch) == 0 {
				context.MakeErrorResponse(common.InvalidRequestError)
				break
			}

			return e.ProcessRawBatch(batch.Batch, context)
		} else if b == '{' {
			if err = json.Unmarshal(context.RawRequest, &context.JsonRpcRequest); err != nil {
				context.MakeErrorResponse(common.ParseError)
				break
			}

			return e.ProcessRawRequest(context)
		} else {
			context.MakeErrorResponse(common.InvalidRequestError)
			break
		}

	}

	_ = context.RebuildRawResponse()

	return err
}

// Get a list of RAW requests of the batch, process each request, return RAW batch response
func (e *JsonRpcServer) ProcessRawBatch(batch []json.RawMessage, context *common.RequestContext) (err error) {

	results := make([]json.RawMessage, len(batch))

	for i, rawRequest := range batch {
		localContext := *context
		localContext.RawRequest = rawRequest
		_ = e.ProcessRawRequest(&localContext)
		results[i] = localContext.RawResponse
	}

	context.RawResponse, err = json.Marshal(results)

	return
}

// Process RAW request, return RAW result
func (e *JsonRpcServer) ProcessRawRequest(context *common.RequestContext) (err error) {

	// Get Json-Rpc request from byte array
	err = context.ParseRawRequest()
	if err != nil {
		_ = context.RebuildRawResponse()
		return
	}

	// Get method from the server
	if method, ok := e.GetMethod(context.JsonRpcRequest.Method); ok {
		// Apply pre-processing pipeline
		context.ApplyPipeline(&e.PreProcessingStages)

		// InvokeMethod the method
		err = InvokeMethod(context, method)

		// Apply post-processing pipeline
		context.ApplyPipeline(&e.PostProcessingStages)

	} else {
		context.MakeErrorResponse(common.MethodNotFoundError)
		err = nil
	}

	// Rebuild raw response
	_ = context.RebuildRawResponse()

	return
}
