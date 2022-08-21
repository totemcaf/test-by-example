package jsonx

const (
	VarExpansion Type = "varExpansion"
)

type varExpansionType struct {
	varName string
}

func (n *varExpansionType) Equals(actual JsonX) bool {
	panic("implement varExpansion.Equals")
}

func (n *varExpansionType) Diff(context Context, actual JsonX) Differences {
	expected := n.Eval(context)

	if expected.Equals(actual) {
		return nil
	}

	return Differences{{nil, n, expected, actual, "different"}}
}

func (n *varExpansionType) Eval(context Context) JsonX {
	return NewParser().Parse(context.Get(n.varName))
}

func (n *varExpansionType) Type() Type {
	return VarExpansion
}

func (n *varExpansionType) String() string {
	return "${" + n.varName + "}"
}
