package jsonx

import "fmt"

const (
	Extractor = "Extractor"
)

type extractorType struct {
	varName string
}

func (e extractorType) Type() Type {
	return Extractor
}

func (e extractorType) Eval(_ Context) JsonX {
	panic(fmt.Sprintf("evaluators cannot evaluate. Use '${%s}' instead", e.varName))
}

func (e extractorType) Diff(context Context, actual JsonX) Differences {
	context.Set(e.varName, actual)

	return nil
}

func (e extractorType) Equals(_ JsonX) bool {
	panic("implement extractor.Equals")
}

func (e extractorType) String() string {
	panic("implement extractor.String")
}
