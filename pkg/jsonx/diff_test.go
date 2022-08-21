package jsonx

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDiff(t *testing.T) {
	context := NewContext()
	context.Set("name", "Chrisjen Avasarala")
	context.Set("age", 42)
	context.Set("is_cool", true)

	tests := []struct {
		name     string
		expected JsonX
		actual   JsonX
		want     Differences
	}{
		{"Same null", p(nil), p(nil), nil},
		{"String same as string", p("hello"), p("hello"), nil},
		{"int same as int", p(42), p(42), nil},

		{"expression same as string", p("$name"), p("Chrisjen Avasarala"), nil},
		{"complex expression same as string", p("Ms $name is Secretary-General"), p("Ms Chrisjen Avasarala is Secretary-General"), nil},
		{"expression same as int", p("$age"), p(42), nil},

		{"Null different than string", p(nil), p("hello"), Differences{{nil, p(nil), p(nil), p("hello"), "different"}}},
		{"String different than null", p("hello"), p(nil), Differences{{nil, p("hello"), p("hello"), p(nil), "different"}}},
		{"String different than string", p("hello"), p("chao"), Differences{{nil, p("hello"), p("hello"), p("chao"), "different"}}},
		{"Int different than int", p(42), p(24), Differences{{nil, p(42), p(42), p(24), "different"}}},
		{"String different than int", p("42"), p(24), Differences{{nil, p("42"), p("42"), p(24), "different"}}},

		{
			"expression different as string",
			p("$name"),
			p("Penelope Cruz"),
			Differences{{nil, p("$name"), p("Chrisjen Avasarala"), p("Penelope Cruz"), "different"}},
		},
		{
			"complex expression different as string",
			p("Ms $name is Secretary-General"),
			p("Ms Penelope Cruz is Secretary-General"),
			Differences{{nil, p("Ms $name is Secretary-General"), p("Ms Chrisjen Avasarala is Secretary-General"), p("Ms Penelope Cruz is Secretary-General"), "different"}},
		},
		{
			"expression different as int",
			p("$age"),
			p(76),
			Differences{{nil, p("$age"), p(42), p(76), "different"}},
		},
		{
			"array different as string",
			p([]interface{}{"hello", "world"}),
			p("find"),
			Differences{{nil, p([]interface{}{"hello", "world"}), p([]interface{}{"hello", "world"}), p("find"), "expected array"}},
		},
		{
			"array different length",
			p([]interface{}{"hello", "world"}),
			p([]interface{}{"find"}),
			Differences{{nil, p([]interface{}{"hello", "world"}), p([]interface{}{"hello", "world"}), p([]interface{}{"find"}), "different array lengths"}},
		},
		{
			"array one different value",
			p([]interface{}{"hello", "world"}),
			p([]interface{}{"hello", "moon"}),
			Differences{{[]string{"1"}, p("world"), p("world"), p("moon"), "different"}},
		},
		{
			"array different values",
			p([]interface{}{"hello", "world"}),
			p([]interface{}{"by", "moon"}),
			Differences{
				{[]string{"0"}, p("hello"), p("hello"), p("by"), "different"},
				{[]string{"1"}, p("world"), p("world"), p("moon"), "different"},
			},
		},

		{
			name:     "map different as string",
			expected: p(map[string]interface{}{"hello": "world"}),
			actual:   p("find"),
			want:     Differences{{nil, p(map[string]interface{}{"hello": "world"}), p(map[string]interface{}{"hello": "world"}), p("find"), "expected map"}},
		},
		{
			"same empty structs",
			p(struct{}{}),
			p(struct{}{}),
			nil,
		},
		{
			"same structs with one value",
			p(struct{ name string }{name: "Chrisjen"}),
			p(struct{ name string }{name: "Chrisjen"}),
			nil,
		},
		{
			"same structs with expression",
			p(map[string]any{"name": "$name", "age": "$age"}),
			p(struct {
				name string
				age  int
			}{name: "Chrisjen Avasarala", age: 42}),
			nil,
		},
		{
			"structs with one different value",
			p(struct{ name string }{name: "Chrisjen"}),
			p(struct{ name string }{name: "James"}),
			Differences{{[]string{"name"}, p("Chrisjen"), p("Chrisjen"), p("James"), "different"}},
		},
		{
			"structs with two different value",
			p(struct {
				name, lastName string
				age            int
			}{name: "Chrisjen", lastName: "Avasarala", age: 42}),
			p(struct {
				name, lastName string
				age            int
			}{name: "James", lastName: "Holden", age: 42}),
			Differences{
				{[]string{"name"}, p("Chrisjen"), p("Chrisjen"), p("James"), "different"},
				{[]string{"lastName"}, p("Avasarala"), p("Avasarala"), p("Holden"), "different"},
			},
		},
		{
			"struct: expected value missing on actual",
			p(struct{ name, lastName string }{name: "Chrisjen", lastName: "Avasarala"}),
			p(struct{ name string }{name: "Chrisjen"}),
			Differences{
				{[]string{"lastName"}, p("Avasarala"), p("Avasarala"), p(nil), "missing value"},
			},
		},
		{
			"struct: actual has additional value not in expected",
			p(struct{ name string }{name: "Chrisjen"}),
			p(struct{ name, position string }{name: "Chrisjen", position: "Secretaries-General"}),
			Differences{
				{[]string{"position"}, p(nil), p(nil), p("Secretaries-General"), "extra value"},
			},
		},
		// struct: with expression
		{
			"multi levels structs with different values",
			p(struct{ address struct{ city, state string } }{address: struct{ city, state string }{city: "New York", state: "NY"}}),
			p(struct{ address struct{ city, state string } }{address: struct{ city, state string }{city: "Morón", state: "NY"}}),
			Differences{
				{[]string{"city", "address"}, p("New York"), p("New York"), p("Morón"), "different"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			differences := tt.expected.Diff(context, tt.actual)
			assert.Equalf(t, tt.want, differences, "Diff(context, %v, %v)", tt.expected, tt.actual)
		})
	}
}

var aParser = NewParser()

func p(v interface{}) JsonX {
	return aParser.Parse(v)
}
