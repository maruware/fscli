package fscli

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	tests := []struct {
		desc  string
		input string
		want  Operation
	}{
		{
			desc:  "simple query",
			input: `QUERY user WHERE name == "John Doe"`,
			want: &QueryOperation{
				opType:     QUERY,
				collection: "user",
				filters: []Filter{
					NewStringFilter("name", "==", "John Doe"),
				},
			},
		},
		{
			desc:  "query with int",
			input: `QUERY user WHERE age == 20`,
			want: &QueryOperation{
				opType:     QUERY,
				collection: "user",
				filters: []Filter{
					NewIntFilter("age", "==", 20),
				},
			},
		},
		{
			desc:  "query with float",
			input: `QUERY user WHERE age == 20.5`,
			want: &QueryOperation{
				opType:     QUERY,
				collection: "user",
				filters: []Filter{
					NewFloatFilter("age", "==", 20.5),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			l := NewLexer(tt.input)
			p := NewParser(l)
			got := p.Parse()

			assert.Equal(t, tt.want, got)
		})
	}
}
