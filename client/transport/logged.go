package transport

import (
	"fmt"
	"log"
)

// Auxiliary type holding a logger (for transports)
type Logged struct {
	logger *log.Logger
}

func (t *Logged) SetLogger(l *log.Logger) error {
	if l == nil {
		return fmt.Errorf("nil logger not allowed")
	}

	t.logger = l
	return nil
}
