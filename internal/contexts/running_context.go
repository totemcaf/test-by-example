package contexts

import (
	"github.com/totemcaf/test-by-example.git/internal/model"
	"go.uber.org/zap"
)

// RunningContext contains the defined and captured values for a test run
// This implementation is based on code from: https://gist.github.com/hvoecking/10772475
type RunningContext interface {
	Set(name string, expression model.AnyValue)
	Get(name string) interface{}
}

type runningContext struct {
	entries map[string]model.AnyValue
	logger  *zap.SugaredLogger
}

func NewRunningContext(logger *zap.SugaredLogger) RunningContext {
	return &runningContext{
		make(map[string]model.AnyValue, 0),
		logger,
	}
}

// Set sets the value of a variable
func (c runningContext) Set(name string, expression model.AnyValue) {
	c.logger.Debugf("Setting %s to %s\n", name, expression)
	c.entries[name] = expression
}

// Get returns the current value of a variable
func (c runningContext) Get(name string) interface{} {
	return c.entries[name]
}
