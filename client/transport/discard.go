package transport

import (
	"github.com/yekhlakov/gojsonrpc/common"
	"github.com/yekhlakov/gojsonrpc/server"
)

// Black hole transport
// Consumes all, returns nothing
type Discard struct {
	Logged
	Server *server.JsonRpcServer
}

// No request handling
func (t *Discard) PerformRequest(rc *common.RequestContext) error {
	return nil
}

// No pre-processing stages
func (t *Discard) AddPreProcessingStage(stage common.Stage) {
}

// No post-processing stages
func (t *Discard) AddPostProcessingStage(stage common.Stage) {
}
