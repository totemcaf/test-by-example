package jsonx

const (
	Null Type = "null"
)

type nullType struct {
}

var NullX = &nullType{}

func (n *nullType) Eval(_ Context) JsonX {
	return n
}

func (n *nullType) Type() Type {
	return Null
}

func (n *nullType) String() string {
	return "<null>"
}

func (n *nullType) Equals(other JsonX) bool {
	_, ok := other.(*nullType)
	return ok
}

func (n *nullType) Diff(_ Context, actual JsonX) Differences {
	if n.Equals(actual) {
		return nil
	}

	return Differences{{nil, n, n, actual, "different"}}
}
