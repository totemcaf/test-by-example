package jsonx

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_eval(t *testing.T) {
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
		{"Null", nil, &nullType{}},
		{"String", "hello", &stringType{"hello"}},
		{"Int", 123, &intType{123}},
		{"Simple var expansion", "$someVar", &stringType{"hello"}},
		{"Simple var expansion", "$someInt", &intType{987}},
		{"Text and var expansions", "The values is $someVar with an $someInt", &stringType{"The values is hello with an 987"}},
		{"Another expansions", "The sum of $someInt plus '$$${someInt}' is ${result}", &stringType{"The sum of 987 plus '$987' is I don't know"}},
		{"Array without expansions", []string{"hello"}, &arrayType{[]JsonX{&stringType{"hello"}}}},
		{"Array with expansions", []string{"$someVar"}, &arrayType{[]JsonX{&stringType{"hello"}}}},
		{"Map without expansions", map[string]string{"value1": "my world"}, &mapType{map[string]JsonX{"value1": &stringType{"my world"}}}},
		{"Map with expansions", map[string]string{"value2": "$someInt"}, &mapType{map[string]JsonX{"value2": &intType{987}}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewParser()
			parsed := p.Parse(tt.jsonObject)

			evaluated := parsed.Eval(context)

			assert.Equalf(t, tt.want, evaluated, "Eval(%v)", tt.jsonObject)
		})
	}
}

var complexJson = `
{
  "body": {
    "id": "${flowId}",
    "partnerId": "627e50c8112ee12b37cccede",
    "reference": {
      "id": "ALK-00028-B",
      "data": {
        "customerType": "VIP"
      }
    },
    "isApproved": false,
    "flow": {
      "startDate": "${startTime}",
      "state": "in-progress",
      "currentStep": "",
      "lastStep": ""
    },
    "state": "pending",
    "data": {},
    "sourceConnections": []
  }
} 
`

func Test_Eval_a_complex_object(t *testing.T) {
	// GIVEN a json object with a complex structure
	var jsonObject interface{}
	_ = json.Unmarshal([]byte(complexJson), &jsonObject)

	jsonX := NewParser().Parse(jsonObject)

	textContext := &SimpleContext{
		vars: map[string]any{
			"flowId":    "ALK-00028-Actual",
			"startTime": "2018-01-01T00:00:00Z",
		},
	}

	// WHEN it is evaluated
	result := jsonX.Eval(textContext)

	// THEN it should return a JsonX object with the same structure
	assert.Equal(t, &mapType{map[string]JsonX{
		"body": &mapType{map[string]JsonX{
			"id":        &stringType{"ALK-00028-Actual"},
			"partnerId": &stringType{"627e50c8112ee12b37cccede"},
			"reference": &mapType{map[string]JsonX{
				"id": &stringType{"ALK-00028-B"},
				"data": &mapType{map[string]JsonX{
					"customerType": &stringType{"VIP"},
				}},
			}},
			"flow": &mapType{map[string]JsonX{
				"startDate":   &stringType{"2018-01-01T00:00:00Z"},
				"state":       &stringType{"in-progress"},
				"currentStep": &stringType{""},
				"lastStep":    &stringType{""},
			}},
			"isApproved":        &boolType{false},
			"state":             &stringType{"pending"},
			"data":              &mapType{map[string]JsonX{}},
			"sourceConnections": &arrayType{[]JsonX{}},
		}},
	}}, result)
}

func Test_extractor_fails_to_eval(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			assert.Equal(t, "evaluators cannot evaluate. Use '${someVar}' instead", r)
		}
	}()

	p("$(someVar)").Eval(NewContext())

	t.Fail()
}
