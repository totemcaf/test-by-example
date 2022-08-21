package jsonx

type Type string

type JsonX interface {
	Type() Type
	// Eval evaluates this JsonX node and returns a new "literal" JsonX without expressions.
	Eval(context Context) JsonX
	// Diff Compares this JsonX node with another JsonX node and returns a list of differences.
	// The differences are returned in the same order as they appear in the JsonX tree.
	// First value is the expected one, second is the actual value. The expected value can contain
	// placeholders for context variables, but also value extractors.
	// If it is a placeholder, the context value will be used to compare. If it is an extractor,
	// the value will be set in the context and the comparison will success. (Note: When adding format
	// matchers the extractor or matcher will check the expected.)
	Diff(context Context, actual JsonX) Differences

	// Equals compares this JsonX with another JsonX and returns true if they are same type and same value.
	// No evaluations are performed.
	Equals(actual JsonX) bool
	String() string
}

type Parser interface {
	Parse(jsonObject interface{}) JsonX
}
