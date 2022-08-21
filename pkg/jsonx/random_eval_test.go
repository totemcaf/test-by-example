package jsonx

import (
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"
)

func Test_random_eval(t *testing.T) {

	context := &SimpleContext{
		vars: map[string]any{
			"someVar": "hello",
			"someInt": 987,
			"result":  "I don't know",
		},
	}
	tests := []struct {
		name       string
		jsonObject interface{}
		want       JsonX
	}{
		{"Random string", "${:random.string}", p("hrUKPt")},
		{"Random name", "${:random.name}", p("Jeromy Schmeler")},
		{"Random email", "${:random.email}", p("jeromyschmeler@ziemann.biz")},
		{"Random phone", "${:random.phone}", p("5780357683")},
		{"Random address", "${:random.address}", p("803 Harbor burgh, Chula Vista, New Hampshire 97582, Chula Vista, New Hampshire, 97582, Mali")},
		{"Random company name", "${:random.companyName}", p("GovTribe")},
		{"Random regex", `${:random.regex:/[a-zA-Z]{3\}/}`, p("h{3}")},
	}
	defer gofakeit.Seed(0)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewParser()
			parsed := p.Parse(tt.jsonObject)

			gofakeit.Seed(42) // Ensure repeatable values
			evaluated := parsed.Eval(context)

			fmt.Printf("%v\n", evaluated)

			assert.Equalf(t, tt.want, evaluated, "Eval(%v)", tt.jsonObject)
		})
	}
}

func Test_random_eval_sets_context(t *testing.T) {
	// GIVEN empty context
	context := NewContext()

	// WHEN evaluating random.string
	parsed := p("${result:random.string}")

	gofakeit.Seed(42) // Ensure repeatable values
	_ = parsed.Eval(context)

	// THEN context is set
	assert.Equal(t, p("hrUKPt"), context.Get("result"))
}
