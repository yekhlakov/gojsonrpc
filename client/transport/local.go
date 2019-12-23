package transport

import (
	"github.com/yekhlakov/gojsonrpc/common"
	"github.com/yekhlakov/gojsonrpc/server"
)

type Local struct {
	Logged
	Server               *server.JsonRpcServer
	PreProcessingStages  []common.Stage
	PostProcessingStages []common.Stage
}

func (t *Local) PerformRequest(rc *common.RequestContext) error {
	rc.ApplyPipeline(&t.PreProcessingStages)
	err := t.Server.ProcessRawInput(rc)
	rc.ApplyPipeline(&t.PostProcessingStages)
	return err
}

func (t *Local) AddPreProcessingStage(stage common.Stage) {
	t.PreProcessingStages = append(t.PreProcessingStages, stage)
}

func (t *Local) AddPostProcessingStage(stage common.Stage) {
	t.PostProcessingStages = append(t.PostProcessingStages, stage)
}
