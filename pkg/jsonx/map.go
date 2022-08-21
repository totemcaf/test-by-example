package jsonx

import (
	"encoding/json"
	"fmt"
)

const (
	Map Type = "map"
)

type mapType struct {
	values map[string]JsonX
}

func (n *mapType) MarshalJSON() ([]byte, error) {
	return json.Marshal(n.values)
}

func (n *mapType) String() string {
	return fmt.Sprintf("Map(%s)", n.values)
}

func (n *mapType) Equals(_ JsonX) bool {
	panic("implement map.Equals")
}

func (n *mapType) Diff(context Context, actual JsonX) Differences {
	otherMap, ok := actual.(*mapType)
	if !ok {
		return Differences{{nil, n, n, actual, "expected map"}}
	}

	var differences []*Difference

	for key, value := range n.values {
		otherValue, ok := otherMap.values[key]
		if !ok {
			differences = append(differences, &Difference{[]string{key}, value, value, NullX, "missing value"})
			continue
		} else {
			differences = append(differences, value.Diff(context, otherValue).addPath(key)...)
		}
	}

	for key, value := range otherMap.values {
		if _, ok := n.values[key]; !ok {
			differences = append(differences, &Difference{[]string{key}, NullX, NullX, value, "extra value"})
		}
	}

	return differences
}

func (n *mapType) Eval(context Context) JsonX {
	values := make(map[string]JsonX)

	for key, value := range n.values {
		values[key] = value.Eval(context)
	}

	return &mapType{values: values}
}

func (n *mapType) Type() Type {
	return Map
}
