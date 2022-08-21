package jsonx

import (
	"encoding/json"
	"strconv"
)

const (
	Array Type = "array"
)

type arrayType struct {
	values []JsonX
}

func (n *arrayType) MarshalJSON() ([]byte, error) {
	return json.Marshal(n.values)
}

func (n *arrayType) String() string {
	panic("implement array.String")
}

func (n *arrayType) Equals(_ JsonX) bool {
	panic("implement array.Equals")
}

func (n *arrayType) Diff(context Context, actual JsonX) Differences {
	diffs := make(Differences, 0)

	actualArray, ok := actual.(*arrayType)

	if !ok {
		return Differences{{nil, n, n, actual, "expected array"}}
	}

	if len(n.values) != len(actualArray.values) {
		return Differences{{nil, n, n, actual, "different array lengths"}}
	}

	actualValues := actualArray.values

	for i, value := range n.values {
		diff := value.Diff(context, actualValues[i])
		if len(diff) > 0 {
			diffs = append(diffs, diff.addPath(strconv.Itoa(i))...)
		}
	}

	return diffs
}

func (n *arrayType) Eval(context Context) JsonX {
	values := make([]JsonX, len(n.values))
	for i, value := range n.values {
		values[i] = value.Eval(context)
	}
	return &arrayType{values: values}
}

func (n *arrayType) Type() Type {
	return Array
}
