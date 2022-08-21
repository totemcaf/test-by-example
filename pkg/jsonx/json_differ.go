package jsonx

import "fmt"

type Differ struct {
	context     Context
	parser      Parser
	differences Differences
}

func NewDiffer(context Context) *Differ {
	return &Differ{context: context, parser: NewParser()}
}

func (d *Differ) Compare(expected, actual any) error {
	actualJsonX := d.parser.Parse(actual)
	expectedJsonX := d.parser.Parse(expected)

	// Compare first to process all extractors so varExpansions will have the corresponding values
	// TODO improve this to do everything in one pass (visit map entries in original order, allow extractors in expressions, etc)
	_ = expectedJsonX.Diff(d.context, actualJsonX)

	d.differences = expectedJsonX.Diff(d.context, actualJsonX)

	if len(d.differences) == 0 {
		return nil
	}
	return fmt.Errorf("%s", d.String())
}

func (d *Differ) Differences() Differences {
	return d.differences
}

func (d *Differ) String() string {
	return d.differences.String()
}
