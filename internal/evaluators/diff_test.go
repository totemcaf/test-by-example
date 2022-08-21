package evaluators

import (
	"testing"

	"bitbucket.org/altscore/test-by-example.git/internal/contexts"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func Test_comparator_Compare_no_differences(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	type fields struct {
		definitions []string
	}
	type args struct {
		expected any
		actual   any
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "nil with nil",
			fields: fields{
				[]string{},
			},
			args: args{
				expected: nil,
				actual:   nil,
			},
		},
		{
			name: "int with int",
			fields: fields{
				[]string{},
			},
			args: args{
				expected: 42,
				actual:   42,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			context := contexts.NewRunningContext(logger.Sugar())

			c := NewComparatorFor(context)
			assert.Len(t, c.Compare(tt.args.expected, tt.args.actual), 0, "Compare(%v, %v)", tt.args.expected, tt.args.actual)
		})
	}
}
func Test_comparator_Compare_with_differences(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	type fields struct {
		definitions []string
	}
	type args struct {
		expected any
		actual   any
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []Difference
	}{
		{
			name: "nil with int",
			fields: fields{
				[]string{},
			},
			args: args{
				expected: nil,
				actual:   123,
			},
			want: []Difference{
				{
					Expected: nil,
					Actual:   123,
				},
			},
		},
		{
			name: "int with nil",
			fields: fields{
				[]string{},
			},
			args: args{
				expected: nil,
				actual:   123,
			},
			want: []Difference{
				{
					Expected: nil,
					Actual:   123,
				},
			},
		},
		{
			name: "int with int",
			fields: fields{
				[]string{},
			},
			args: args{
				expected: 42,
				actual:   123,
			},
			want: []Difference{
				{
					Expected: 42,
					Actual:   123,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			context := contexts.NewRunningContext(logger.Sugar())

			c := NewComparatorFor(context)
			assert.Equalf(t, tt.want, c.Compare(tt.args.expected, tt.args.actual), "Compare(%v, %v)", tt.args.expected, tt.args.actual)
		})
	}
}
