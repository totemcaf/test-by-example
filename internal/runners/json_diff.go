package runners

import (
	"bitbucket.org/altscore/test-by-example.git/internal/contexts"
	"bitbucket.org/altscore/test-by-example.git/internal/model"
)

type Diff struct {
}

type JsonDiffer struct {
	contexts.RunningContext
}

// TODO check if it should be removed

func NewJsonDiffer(context contexts.RunningContext) *JsonDiffer {
	return &JsonDiffer{context}
}

func (d JsonDiffer) Compare(body map[string]any, step *model.StepSpec) error {

	return nil
}
