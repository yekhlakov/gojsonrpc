package client

import (
	"log"

	"github.com/yekhlakov/gojsonrpc/common"
)

type Transport interface {
	PerformRequest(rc *common.RequestContext) error
	AddPreProcessingStage(stage common.Stage)
	AddPostProcessingStage(stage common.Stage)
	SetLogger(l *log.Logger) error
}

type Client struct {
	T      Transport
	logger *log.Logger
}

// Dummy interface for method signature containers
type Handler interface {
}
