package jsonx

import (
	"encoding/json"
	"fmt"
)

const (
	Int Type = "int"
)

type intType struct {
	value int64
}

func (n *intType) MarshalJSON() ([]byte, error) {
	return json.Marshal(n.value)
}

func (n *intType) Equals(value JsonX) bool {
	other, ok := value.(*intType)
	return ok && other.value == n.value
}

func (n *intType) Diff(_ Context, actual JsonX) Differences {
	if n.Equals(actual) {
		return nil
	}

	return Differences{{nil, n, n, actual, "different"}}
}

func (n *intType) Eval(_ Context) JsonX {
	return n
}

func (n *intType) Type() Type {
	return Int
}

func (n *intType) String() string {
	return fmt.Sprintf("%d", n.value)
}
