package evaluators

import (
	"reflect"

	"bitbucket.org/altscore/test-by-example.git/internal/contexts"
)

type Difference struct {
	Actual   interface{}
	Expected interface{}
}
type Comparator interface {
	// Compare does a structural comparison the expected against the actual
	// In Placeholder are replaced by context values
	// Out placeholders updates the context with their values
	Compare(expected any, actual any) []Difference
}

type comparator struct {
	context     contexts.RunningContext
	differences []Difference
}

func NewComparatorFor(context contexts.RunningContext) Comparator {
	return &comparator{
		context,
		[]Difference{},
	}
}

func (c *comparator) Compare(expected any, actual any) []Difference {
	// Wrap the expected and actual values in reflect.Value
	originalExpected := reflect.ValueOf(expected)
	originalActual := reflect.ValueOf(actual)

	c.compareValue(originalExpected, originalActual)

	// Remove the reflection wrapper
	return c.differences
}

func (c *comparator) compareValue(expected reflect.Value, actual reflect.Value) {
	switch expected.Kind() {
	// The first cases handle nested structures and Evaluate them recursively

	case reflect.Ptr:
		c.evaluatePtr(expected, actual)
	case reflect.Interface:
		c.evaluateInterface(expected, actual)
	case reflect.Struct:
		c.evaluateStruct(expected, actual)
	case reflect.Slice:
		c.evaluateSlice(expected, actual)
	case reflect.Map:
		if c.isBas64Map(expected, actual) {
			c.evaluateBase64Map(expected, actual)
		} else {
			c.evaluateMap(expected, actual)
		}

	// Otherwise, we cannot traverse anywhere so this finishes the recursion

	// If it is a string Evaluate it (yay finally we're doing what we came for)
	case reflect.String:
		c.evaluateStr(expected, actual)

	// And everything else will simply be taken from the expected	default:
	default:
		c.evaluateOthers(expected, actual)
	}
}

func (c *comparator) evaluatePtr(expected reflect.Value, actual reflect.Value) {

}

func (c *comparator) evaluateInterface(expected reflect.Value, actual reflect.Value) {

}

func (c *comparator) evaluateStruct(expected reflect.Value, actual reflect.Value) {

}

func (c *comparator) evaluateSlice(expected reflect.Value, actual reflect.Value) {

}

func (c *comparator) isBas64Map(expected reflect.Value, actual reflect.Value) bool {
	return false
}

func (c *comparator) evaluateBase64Map(expected reflect.Value, actual reflect.Value) {

}

func (c *comparator) evaluateMap(expected reflect.Value, actual reflect.Value) {

}

func (c *comparator) evaluateStr(expected reflect.Value, actual reflect.Value) {

}

func (c *comparator) evaluateOthers(expected reflect.Value, actual reflect.Value) {
	if !expected.IsValid() && !actual.IsValid() {
		return
	} else if !expected.IsValid() {
		c.addDifference(
			false,
			nil,
			actual.Interface(),
		)
	} else if !actual.IsValid() {
		c.addDifference(
			false,
			expected.Interface(),
			nil,
		)
	} else if expected.Type() == actual.Type() {
		c.addDifference(
			expected.Interface() == actual.Interface(),
			expected.Interface(),
			actual.Interface(),
		)
	}
}

func (c *comparator) addDifference(areEqual bool, expected any, actual any) {
	if !areEqual {
		c.differences = append(c.differences,
			Difference{
				Actual:   actual,
				Expected: expected,
			},
		)
	}

}
