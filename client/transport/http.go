package transport

import (
	"github.com/yekhlakov/gojsonrpc/common"
)

type Http struct {
	Logged
	PreProcessingStages  []common.Stage
	PostProcessingStages []common.Stage
}

func (t *Http) PerformRequest(rc *common.RequestContext) error {
	rc.ApplyPipeline(&t.PreProcessingStages)
	// TODO: http processing

	rc.ApplyPipeline(&t.PostProcessingStages)
	return nil
}

func (t *Http) AddPreProcessingStage(stage common.Stage) {
	t.PreProcessingStages = append(t.PreProcessingStages, stage)
}

func (t *Http) AddPostProcessingStage(stage common.Stage) {
	t.PostProcessingStages = append(t.PostProcessingStages, stage)
}
