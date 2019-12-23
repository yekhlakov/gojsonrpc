package transport

import (
	"encoding/json"
	"fmt"

	"github.com/yekhlakov/gojsonrpc/common"
	"github.com/yekhlakov/gojsonrpc/server"
)

// ERROR transport
// Consumes all, returns requested errors
type Error struct {
	Logged
	Server       *server.JsonRpcServer
	ErrorMessage string
	RawResponse  json.RawMessage
	JsonRpcError common.Error
}

// No request handling, always error out
func (t *Error) PerformRequest(rc *common.RequestContext) error {
	if t.ErrorMessage != "" {
		return fmt.Errorf(t.ErrorMessage)
	}

	if t.RawResponse != nil {
		rc.RawResponse = t.RawResponse
		return nil
	}

	rc.MakeErrorResponse(t.JsonRpcError)
	return nil
}

// No pre-processing stages
func (t *Error) AddPreProcessingStage(stage common.Stage) {
}

// No post-processing stages
func (t *Error) AddPostProcessingStage(stage common.Stage) {
}
