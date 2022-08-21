package jsonx

import "fmt"

const Concatenation Type = "concatenation"

type concatenationType struct {
	values []JsonX
}

func (n *concatenationType) String() string {
	return fmt.Sprintf("%v", n.values)
}

func (n *concatenationType) Equals(_ JsonX) bool {
	panic("implement concatenation.Equals")
}

func (n *concatenationType) Diff(context Context, actual JsonX) Differences {
	expected := n.Eval(context)

	if expected.Equals(actual) {
		return nil
	}

	return Differences{{nil, n, expected, actual, "different"}}

}

func (n *concatenationType) Eval(context Context) JsonX {
	var result string
	for _, value := range n.values {
		eval := value.Eval(context)

		result += fmt.Sprintf("%v", eval)
	}
	return &stringType{result} // NewParser().Parse(result)
}

func (n *concatenationType) Type() Type {
	return Concatenation
}
