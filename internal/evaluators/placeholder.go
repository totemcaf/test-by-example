package evaluators

import (
	"fmt"
	"reflect"

	"bitbucket.org/altscore/test-by-example.git/internal/contexts"
)

type Placeholder interface {
	Value(c contexts.RunningContext) string
}

type literalPlaceholder struct {
	value string
}

func (l literalPlaceholder) Value(_ contexts.RunningContext) string {
	return l.value
}

type variablePlaceholder struct {
	varName string
}

func (v variablePlaceholder) Value(c contexts.RunningContext) string {
	if value := c.Get(v.varName); value == nil {
		return "<nil>"
	} else {
		// Remove the interface
		return fmt.Sprintf("%v", reflect.ValueOf(value).Interface())
	}
}
