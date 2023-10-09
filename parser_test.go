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
			input: `QUERY user`,
			want:  NewQueryOperation("user", nil, nil),
		},
		{
			desc:  "query with where",
			input: `QUERY user WHERE name == "John Doe"`,
			want: NewQueryOperation("user", nil, []Filter{
				NewStringFilter("name", "==", "John Doe"),
			}),
		},
		{
			desc:  "query with int",
			input: `QUERY user WHERE age == 20`,
			want: NewQueryOperation("user", nil, []Filter{
				NewIntFilter("age", OPERATOR_EQ, 20),
			}),
		},
		{
			desc:  "query with float",
			input: `QUERY user WHERE age == 20.5`,
			want: NewQueryOperation("user", nil, []Filter{
				NewFloatFilter("age", OPERATOR_EQ, 20.5),
			}),
		},
		{
			desc:  "query with IN",
			input: `QUERY user WHERE age IN [20, 21, 22]`,
			want: NewQueryOperation("user", nil, []Filter{
				NewArrayFilter("age", OPERATOR_IN, []any{20, 21, 22}),
			}),
		},
		{
			desc:  "query with IN by mutiple types",
			input: `QUERY user WHERE age IN [20, 21.5, "22"]`,
			want: NewQueryOperation("user", nil, []Filter{
				NewArrayFilter("age", OPERATOR_IN, []any{20, 21.5, "22"}),
			}),
		},
		{
			desc:  "query with array-contains",
			input: `QUERY user WHERE nicknames ARRAY_CONTAINS "Doe"`,
			want: NewQueryOperation("user", nil, []Filter{
				NewStringFilter("nicknames", OPERATOR_ARRAY_CONTAINS, "Doe"),
			}),
		},
		{
			desc:  "query with array-contains-any",
			input: `QUERY user WHERE nicknames ARRAY_CONTAINS_ANY ["Doe", "John"]`,
			want: NewQueryOperation("user", nil, []Filter{
				NewArrayFilter("nicknames", OPERATOR_ARRAY_CONTAINS_ANY, []any{"Doe", "John"}),
			}),
		},
		{
			desc:  "query with trim head slash",
			input: `QUERY /user`,
			want:  NewQueryOperation("user", nil, nil),
		},
		{
			desc:  "query with select",
			input: `QUERY user SELECT name, age`,
			want:  NewQueryOperation("user", []string{"name", "age"}, nil),
		},
		{
			desc:  "query with select and where",
			input: `QUERY user SELECT name, age WHERE age == 20`,
			want:  NewQueryOperation("user", []string{"name", "age"}, []Filter{NewIntFilter("age", OPERATOR_EQ, 20)}),
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
