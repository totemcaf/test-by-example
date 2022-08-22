package evaluators

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/totemcaf/test-by-example.git/internal/collections/lists"
	"github.com/totemcaf/test-by-example.git/internal/contexts"
)

// Evaluator evaluates an object against the provided context
// This implementation is based on code from: https://gist.github.com/hvoecking/10772475
type Evaluator interface {
	Evaluate(obj interface{}) interface{}
	EvaluateStr(expression string) string
}

type evaluator struct {
	contexts.RunningContext
}

// For is Deprecated, use JsonX instead
func For(context contexts.RunningContext) Evaluator {
	return &evaluator{
		context,
	}
}

// EvaluateStr takes the string expression and replace the placeholders with context values
// Placeholders:
//      ${varName}  replace the Placeholder with the value of the var name
func (c evaluator) EvaluateStr(expression string) string {
	placeholders := c.breakExpression(expression)

	return strings.Join(lists.Map(placeholders, func(ph Placeholder) string { return ph.Value(c) }), "")
}

// breakExpression breaks the given string in the placeholders. It returns a list of strings, uneven
// values are string constants, even values are placeholders
func (c evaluator) breakExpression(expression string) []Placeholder {
	parts := strings.Split(expression, "${")

	results := make([]Placeholder, 0)

	results = append(results, literalPlaceholder{parts[0]})
	for _, part := range parts[1:] {
		subParts := strings.SplitN(part, "}", 2)

		results = append(results, variablePlaceholder{subParts[0]})

		if len(subParts) > 0 {
			results = append(results, literalPlaceholder{subParts[1]})
		}
	}

	return results
}

// Evaluate given a value, replaces all placeholders with the variable values
func (c evaluator) Evaluate(obj interface{}) interface{} {
	// Wrap the original in a reflect.Value
	original := reflect.ValueOf(obj)

	theCopy := c.evaluateValue(original)

	// Remove the reflection wrapper
	return theCopy.Interface()
}

func (c evaluator) evaluateValue(original reflect.Value) reflect.Value {
	switch original.Kind() {
	// The first cases handle nested structures and Evaluate them recursively

	case reflect.Ptr:
		return c.evaluatePtr(original)
	case reflect.Interface:
		return c.evaluateInterface(original)
	case reflect.Struct:
		return c.evaluateStruct(original)
	case reflect.Slice:
		return c.evaluateSlice(original)
	case reflect.Map:
		if c.isBas64Map(original) {
			return c.evaluateBase64Map(original)
		}
		return c.evaluateMap(original)

	// Otherwise, we cannot traverse anywhere so this finishes the recursion

	// If it is a string Evaluate it (yay finally we're doing what we came for)
	case reflect.String:
		return c.evaluateStr(original)

	// And everything else will simply be taken from the original	default:
	default:
		theCopy := reflect.New(original.Type()).Elem()
		theCopy.Set(original)
		return theCopy
	}
}

// If it is a pointer we need to unwrap and call once again
func (c evaluator) evaluatePtr(original reflect.Value) reflect.Value {
	theCopy := reflect.New(original.Type()).Elem()

	// To get the actual value of the original we have to call Elem()
	// At the same time this unwraps the pointer, so we don't end up in
	// an infinite recursion
	originalValue := original.Elem()
	// Check if the pointer is nil
	if !originalValue.IsValid() {
		return theCopy
	}

	// Allocate a new object and set the pointer to it
	theCopy.Set(reflect.New(originalValue.Type()))

	// Unwrap the newly created pointer
	child := c.evaluateValue(originalValue)

	theCopy.Elem().Set(child)

	return theCopy
}

// If it is an interface (which is very similar to a pointer), do basically the
// same as for the pointer. Though a pointer is not the same as an interface so
// note that we have to call Elem() after creating a new object because otherwise
// we would end up with an actual pointer
func (c evaluator) evaluateInterface(original reflect.Value) reflect.Value {
	theCopy := reflect.New(original.Type()).Elem()

	// Get rid of the wrapping interface
	originalValue := original.Elem()
	// Create a new object. Now new gives us a pointer, but we want the value it
	// points to, so we have to call Elem() to unwrap it
	childVale := c.evaluateValue(originalValue)

	theCopy.Set(childVale)

	return theCopy
}

func (c evaluator) evaluateStruct(original reflect.Value) reflect.Value {
	theCopy := reflect.New(original.Type()).Elem()

	for i := 0; i < original.NumField(); i += 1 {
		copyValue := c.evaluateValue(original.Field(i))
		theCopy.Field(i).Set(copyValue)

		fmt.Printf("--> %v\n", theCopy.Interface())
	}

	return theCopy
}

// If it is a slice we create a new slice and Evaluate each element
func (c evaluator) evaluateSlice(original reflect.Value) reflect.Value {
	theCopy := reflect.New(original.Type()).Elem()

	theCopy.Set(reflect.MakeSlice(original.Type(), original.Len(), original.Cap()))
	for i := 0; i < original.Len(); i += 1 {
		copyValue := c.evaluateValue(original.Index(i))
		theCopy.Index(i).Set(copyValue)
	}

	return theCopy
}

// If it is a map we create a new map and Evaluate each value
func (c evaluator) evaluateMap(original reflect.Value) reflect.Value {
	theType := original.Type()

	if theType.Key().Kind() != reflect.String {
		// Convert maps map[any]string to map[string]string
		// Ensure the map is keyed by strings
		theType = reflect.MapOf(reflect.TypeOf(""), theType.Elem())
	}
	theCopy := reflect.New(theType).Elem()

	// mapType := original.Type()
	// Ensure the map is keyed by strings
	mapType := theCopy.Type() //   reflect.MapOf(reflect.TypeOf(""), originalMapType.Elem())
	theCopy.Set(reflect.MakeMap(mapType))

	for _, key := range original.MapKeys() {
		originalValue := original.MapIndex(key)
		// New gives us a pointer, but again we want the value
		copyValue := c.evaluateValue(originalValue)
		keyStr := fmt.Sprintf("%v", key.String())
		keyValueStr := reflect.ValueOf(keyStr)
		theCopy.SetMapIndex(keyValueStr, copyValue)
	}

	return theCopy
}

func (c evaluator) evaluateStr(original reflect.Value) reflect.Value {
	translatedString := c.EvaluateStr(original.Interface().(string))

	return reflect.ValueOf(translatedString)
}

func (c evaluator) isBas64Map(original reflect.Value) bool {
	return original.Len() == 1 && original.MapKeys()[0].Interface().(string) == "$base64"
}

func (c evaluator) evaluateBase64Map(_ reflect.Value) reflect.Value {
	return reflect.ValueOf("this is the base 64 representation")
}

/*
   case "$base64":
       theCopy.SetMapIndex(keyValueStr, reflect.ValueOf("a base 64"))

*/
