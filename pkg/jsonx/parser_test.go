package jsonx

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_parser_Parse(t *testing.T) {
	tests := []struct {
		name       string
		jsonObject interface{}
		want       JsonX
		ofType     Type
	}{
		{"JsonX", p("hello"), p("hello"), String},

		{"Null", nil, &nullType{}, Null},
		{"Boolean false", false, &boolType{false}, Boolean},
		{"Boolean true", true, &boolType{true}, Boolean},
		{"String", "hello", &stringType{"hello"}, String},
		{"Int", 123, &intType{123}, Int},
		{"Array", []string{"hello"}, &arrayType{[]JsonX{&stringType{"hello"}}}, Array},
		{"Map", map[string]string{"value1": "hello"}, &mapType{map[string]JsonX{"value1": &stringType{"hello"}}}, Map},
		{
			"complex Map",
			map[string]string{"value1": "hello", "value2": "chao"},
			&mapType{map[string]JsonX{"value1": &stringType{"hello"}, "value2": &stringType{"chao"}}},
			Map,
		},
		{
			"interface Map",
			map[interface{}]string{"value1": "hello", "value2": "chao"},
			&mapType{map[string]JsonX{"value1": &stringType{"hello"}, "value2": &stringType{"chao"}}},
			Map,
		},

		{"Struct", sampleStruct{
			String: "hi",
			Int:    42,
		}, &mapType{map[string]JsonX{"String": &stringType{"hi"}, "Int": &intType{42}}}, Map},
		{"escaped $", "$$", &stringType{"$"}, String},
		{"doubled escaped $ ", "$$$$", &concatenationType{[]JsonX{
			&stringType{"$"},
			&stringType{"$"},
		}}, Concatenation},
		{"escaped $ and expansion", "$$$someVar", &concatenationType{[]JsonX{
			&stringType{"$"},
			&varExpansionType{"someVar"},
		}}, Concatenation},
		{"escaped $ and expansion 2", "$$${someVar}", &concatenationType{[]JsonX{
			&stringType{"$"},
			&varExpansionType{"someVar"},
		}}, Concatenation},
		{"Simple var expansion", "$someVar", &varExpansionType{varName: "someVar"}, VarExpansion},
		{"Alternate simple var expansion", "${someVar}", &varExpansionType{varName: "someVar"}, VarExpansion},
		{"Text and var expansions", "The values is $someVar",
			&concatenationType{
				values: []JsonX{
					&stringType{"The values is "},
					&varExpansionType{varName: "someVar"},
				},
			},
			Concatenation,
		},
		{"Text and several var expansions",
			"The values is $someVar with an $intValue",
			&concatenationType{
				values: []JsonX{
					&stringType{"The values is "},
					&varExpansionType{varName: "someVar"},
					&stringType{" with an "},
					&varExpansionType{varName: "intValue"},
				},
			},
			Concatenation,
		},
		{
			"Interface",
			struct{ x interface{} }{x: asInterface("hello")},
			&mapType{map[string]JsonX{"x": &stringType{"hello"}}},
			Map,
		},
		{
			"Pointer to",
			struct{ x *string }{x: asPointer("hello")},
			&mapType{map[string]JsonX{"x": &stringType{"hello"}}},
			Map,
		},

		// Generators
		{"Random name", "${:random.name}", &randomValueType{_type: RandomName}, RandomValue},
		{"Generator with config", "${:random.name:10}", &randomValueType{_type: RandomName, config: "10"}, RandomValue},
		{"Generator with regex", "${:random.name:/a+b*/}", &randomValueType{_type: RandomName, config: "a+b*"}, RandomValue},
		{"Generator with escaped regex", `${:random.name:/a+\}*/}`, &randomValueType{_type: RandomName, config: `a+\}*`}, RandomValue},
		{"Generator with escaped regex 2", `${:random.name:/a+[-\}]*/}`, &randomValueType{_type: RandomName, config: `a+[-\}]*`}, RandomValue},

		// Extractors
		{"Extractor", "$(varToSet)", &extractorType{varName: "varToSet"}, Extractor},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewParser()
			parsed := p.Parse(tt.jsonObject)
			assert.Equalf(t, tt.ofType, parsed.Type(), "Parse(%v) type", tt.jsonObject)
			assert.Equalf(t, tt.want, parsed, "Parse(%v)", tt.jsonObject)
		})
	}
}

func asPointer[T any](t T) *T {
	return &t
}

func asInterface(a any) interface{} {
	return a
}

type sampleStruct struct {
	String string
	Int    int
}
