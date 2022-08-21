package jsonx

import "encoding/json"

const (
	String Type = "string"
)

type stringType struct {
	value string
}

func (n *stringType) MarshalJSON() ([]byte, error) {
	return json.Marshal(n.value)
}

func (n *stringType) Equals(other JsonX) bool {
	otherStr, ok := other.(*stringType)
	return ok && otherStr.value == n.value
}

func (n *stringType) Diff(_ Context, actual JsonX) Differences {
	if n.Equals(actual) {
		return nil
	}

	return Differences{{nil, n, n, actual, "different"}}
}

func (n *stringType) EvalStr(_ Context) string {
	return n.value
}

func (n *stringType) Eval(_ Context) JsonX {
	return n
}

func (n *stringType) Type() Type {
	return String
}

func (n *stringType) String() string {
	return n.value
}
