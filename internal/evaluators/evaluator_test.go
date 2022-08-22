package evaluators

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/totemcaf/test-by-example.git/internal/contexts"
	"go.uber.org/zap"
)

func TestEvaluate_empty(t *testing.T) {
	context := contexts.NewRunningContext(makeLogger())
	evaluator := For(context)

	result := evaluator.Evaluate("")

	assert.Equal(t, "", result)
}

func makeLogger() *zap.SugaredLogger {
	return zap.NewNop().Sugar()
}

func TestEvaluate_literal_string(t *testing.T) {
	context := contexts.NewRunningContext(makeLogger())
	evaluator := For(context)

	result := evaluator.Evaluate("a value")

	assert.Equal(t, "a value", result)
}

func TestEvaluate_pattern_no_value(t *testing.T) {
	context := contexts.NewRunningContext(makeLogger())
	evaluator := For(context)

	result := evaluator.Evaluate("${aVar}")

	assert.Equal(t, "<nil>", result)
}

func TestEvaluate_pattern_with_value(t *testing.T) {
	context := contexts.NewRunningContext(makeLogger())
	evaluator := For(context)
	context.Set("aVar", "sample-value")

	result := evaluator.Evaluate("${aVar}")

	assert.Equal(t, "sample-value", result)
}

func TestEvaluate_pattern_with_value_and_literals(t *testing.T) {
	context := contexts.NewRunningContext(makeLogger())
	evaluator := For(context)
	context.Set("aVar", "sample-value")

	result := evaluator.Evaluate("prefix ${aVar} middle ${aVar} suffix")

	assert.Equal(t, "prefix sample-value middle sample-value suffix", result)
}

func TestEvaluate_pattern_with_numeric_value(t *testing.T) {
	context := contexts.NewRunningContext(makeLogger())
	evaluator := For(context)
	context.Set("aVar", 1234)

	result := evaluator.Evaluate("${aVar}")

	assert.Equal(t, "1234", result)
}

func TestEvaluate_pattern_array_with_values(t *testing.T) {
	context := contexts.NewRunningContext(makeLogger())
	evaluator := For(context)
	context.Set("aVar", "sample-value")
	context.Set("anotherVar", "a-different-string")

	result := evaluator.Evaluate([]string{
		"${aVar}",
		"${anotherVar}",
		"--${aVar}++",
		"some//${anotherVar}/run",
	})

	expected := []string{"sample-value", "a-different-string", "--sample-value++", "some//a-different-string/run"}

	assert.Equal(t, expected, result)
}

func TestEvaluate_pattern_map_with_string_values(t *testing.T) {
	context := contexts.NewRunningContext(makeLogger())
	evaluator := For(context)
	context.Set("aVar", "sample-value")
	context.Set("anotherVar", "a-different-string")

	result := evaluator.Evaluate(map[string]string{
		"aa": "${aVar}",
		"bb": "${anotherVar}",
		"cc": "--${aVar}++",
		"dd": "some//${anotherVar}/run",
	})

	expected := map[string]string{
		"aa": "sample-value",
		"bb": "a-different-string",
		"cc": "--sample-value++",
		"dd": "some//a-different-string/run",
	}

	assert.Equal(t, expected, result)
}

type sampleStruct struct {
	FirstValue       string
	SecondsValue     string
	AnInt            int
	ABoolean         bool
	APointerToString *string
	ANil             *bool
}

func TestEvaluate_pattern_struct_with_string_values(t *testing.T) {
	context := contexts.NewRunningContext(makeLogger())
	evaluator := For(context)
	context.Set("aVar", "sample-value")
	context.Set("anotherVar", "a-different-string")

	complexPattern := "some//${anotherVar}/${aVar}/run"
	input := &sampleStruct{
		FirstValue:       "${aVar}",
		SecondsValue:     "${anotherVar}",
		AnInt:            42,
		ABoolean:         true,
		APointerToString: &complexPattern,
		ANil:             nil,
	}

	result := evaluator.Evaluate(input)

	aComplexResult := "some//a-different-string/sample-value/run"

	var expected = &sampleStruct{
		FirstValue:       "sample-value",
		SecondsValue:     "a-different-string",
		AnInt:            42,
		ABoolean:         true,
		APointerToString: &aComplexResult,
		ANil:             nil,
	}

	assert.Equal(t, expected, result)
	assert.EqualValues(t, expected, result)
}

func TestEvaluate_field_encoded_base64(t *testing.T) {
	context := contexts.NewRunningContext(makeLogger())
	evaluator := For(context)

	data := map[string]any{
		"$base64": map[string]any{
			"field1": "a value",
			"field2": 123,
		},
	}

	result := evaluator.Evaluate(data)

	assert.Equal(t, "this is the base 64 representation", result)
}

// as int: ${aVar:int}
