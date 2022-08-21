package jsonx

import "fmt"

const (
	Boolean Type = "bool"
)

type boolType struct {
	value bool
}

func (n *boolType) Equals(_ JsonX) bool {
	panic("implement boolean.Equal")
}

func (n *boolType) Diff(_ Context, expected JsonX) Differences {
	if expectedBool, ok := expected.(*boolType); !ok {
		return []*Difference{{
			Path:        []string{},
			ExpectedRaw: expected,
			Expected:    expected,
			Actual:      n,
			Message:     "expected a boolean",
		}}
	} else if n.value != expectedBool.value {
		return []*Difference{{
			Path:        []string{},
			ExpectedRaw: expected,
			Expected:    expected,
			Actual:      n,
			Message:     "different value",
		}}
	}
	return nil
}

func (n *boolType) Eval(_ Context) JsonX {
	return n
}

func (n *boolType) Type() Type {
	return Boolean
}

func (n *boolType) String() string {
	return fmt.Sprintf("%v", n.value)
}
