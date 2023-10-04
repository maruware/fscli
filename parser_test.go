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
					NewIntFilter("age", OP_EQ, 20),
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
					NewFloatFilter("age", OP_EQ, 20.5),
				},
			},
		},
		{
			desc:  "query with IN",
			input: `QUERY user WHERE age IN [20, 21, 22]`,
			want: &QueryOperation{
				opType:     QUERY,
				collection: "user",
				filters: []Filter{
					NewArrayFilter("age", OP_IN, []any{20, 21, 22}),
				},
			},
		},
		{
			desc:  "query with IN by mutiple types",
			input: `QUERY user WHERE age IN [20, 21.5, "22"]`,
			want: &QueryOperation{
				opType:     QUERY,
				collection: "user",
				filters: []Filter{
					NewArrayFilter("age", OP_IN, []any{20, 21.5, "22"}),
				},
			},
		},
		{
			desc:  "query with array-contains",
			input: `QUERY user WHERE nicknames ARRAY_CONTAINS "Doe"`,
			want: &QueryOperation{
				opType:     QUERY,
				collection: "user",
				filters: []Filter{
					NewStringFilter("nicknames", OP_ARRAY_CONTAINS, "Doe"),
				},
			},
		},
		{
			desc:  "query with array-contains-any",
			input: `QUERY user WHERE nicknames ARRAY_CONTAINS_ANY ["Doe", "John"]`,
			want: &QueryOperation{
				opType:     QUERY,
				collection: "user",
				filters: []Filter{
					NewArrayFilter("nicknames", OP_ARRAY_CONTAINS_ANY, []any{"Doe", "John"}),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			l := NewLexer(tt.input)
			p := NewParser(l)
			got, err := p.Parse()
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, tt.want, got)
		})
	}
}
