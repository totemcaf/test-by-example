package jsonx

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_extractors_set_value_when_compute_diff(t *testing.T) {
	// GIVEN empty context
	context := NewContext()

	// WHEN compute diff
	expected := p("$(extractedName)")
	actual := p("Robin Hood")

	_ = expected.Diff(context, actual)

	// THEN extractor set value
	assert.Equal(t, p("Robin Hood"), context.Get("extractedName"))
}

func Test_extractors_does_not_generate_diff(t *testing.T) {
	// GIVEN empty context, actual value, and expected result
	context := NewContext()

	expected := p("$(extractedName)")
	actual := p("Robin Hood")

	// WHEN compute diff
	differences := expected.Diff(context, actual)

	// THEN extractor set value
	assert.Len(t, differences, 0)
}
