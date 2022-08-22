package evaluators

import (
	"github.com/totemcaf/test-by-example.git/internal/contexts"
	"github.com/totemcaf/test-by-example.git/pkg/jsonx"
)

type jsonXEvaluator struct {
	parser  jsonx.Parser
	context contexts.RunningContext
}

func NewJsonXEvaluator(context contexts.RunningContext) *jsonXEvaluator {
	return &jsonXEvaluator{
		context: context,
		parser:  jsonx.NewParser(),
	}
}

func (j jsonXEvaluator) Evaluate(obj interface{}) interface{} {
	return j.evaluate(obj)
}

func (j jsonXEvaluator) evaluate(obj interface{}) jsonx.JsonX {
	return j.parser.Parse(obj).Eval(j.context)
}

func (j jsonXEvaluator) EvaluateStr(expression string) string {
	return j.evaluate(expression).String()
}

func (c evaluator) name() {

}
