package fscli

import (
	"testing"

	"cloud.google.com/go/firestore"
	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	tests := []struct {
		desc  string
		input string
		want  ParseResult
	}{
		{
			desc:  "simple query",
			input: `QUERY user`,
			want:  &QueryOperation{collection: "user"},
		},
		{
			desc:  "query with where",
			input: `QUERY user WHERE name = "John Doe"`,
			want: &QueryOperation{collection: "user", filters: []Filter{
				NewStringFilter("name", "==", "John Doe"),
			}},
		},
		{
			desc:  "query with int",
			input: `QUERY user WHERE age = 20`,
			want: &QueryOperation{collection: "user", filters: []Filter{
				NewIntFilter("age", OPERATOR_EQ, 20),
			}},
		},
		{
			desc:  "query with float",
			input: `QUERY user WHERE age = 20.5`,
			want: &QueryOperation{collection: "user", filters: []Filter{
				NewFloatFilter("age", OPERATOR_EQ, 20.5),
			}},
		},
		{
			desc:  "query with not equal",
			input: `QUERY user WHERE age != 20`,
			want: &QueryOperation{collection: "user", filters: []Filter{
				NewIntFilter("age", OPERATOR_NOT_EQ, 20),
			}},
		},
		{
			desc:  "query with greater than",
			input: `QUERY user WHERE age > 20`,
			want: &QueryOperation{collection: "user", filters: []Filter{
				NewIntFilter("age", OPERATOR_GT, 20),
			}},
		},
		{
			desc:  "query with greater than or equal",
			input: `QUERY user WHERE age >= 20`,
			want: &QueryOperation{collection: "user", filters: []Filter{
				NewIntFilter("age", OPERATOR_GTE, 20),
			}},
		},
		{

			desc:  "query with less than",
			input: `QUERY user WHERE age < 20`,
			want: &QueryOperation{collection: "user", filters: []Filter{
				NewIntFilter("age", OPERATOR_LT, 20),
			}},
		},
		{
			desc:  "query with less than or equal",
			input: `QUERY user WHERE age <= 20`,
			want: &QueryOperation{collection: "user", filters: []Filter{
				NewIntFilter("age", OPERATOR_LTE, 20),
			}},
		},
		{
			desc:  "query with IN",
			input: `QUERY user WHERE age IN [20, 21, 22]`,
			want: &QueryOperation{collection: "user", filters: []Filter{
				NewArrayFilter("age", OPERATOR_IN, []any{20, 21, 22}),
			}},
		},
		{
			desc:  "query with IN by mutiple types",
			input: `QUERY user WHERE age IN [20, 21.5, "22"]`,
			want: &QueryOperation{collection: "user", filters: []Filter{
				NewArrayFilter("age", OPERATOR_IN, []any{20, 21.5, "22"}),
			}},
		},
		{
			desc:  "query with array-contains",
			input: `QUERY user WHERE nicknames ARRAY_CONTAINS "Doe"`,
			want: &QueryOperation{collection: "user", filters: []Filter{
				NewStringFilter("nicknames", OPERATOR_ARRAY_CONTAINS, "Doe"),
			}},
		},
		{
			desc:  "query with array-contains-any",
			input: `QUERY user WHERE nicknames ARRAY_CONTAINS_ANY ["Doe", "John"]`,
			want: &QueryOperation{collection: "user", filters: []Filter{
				NewArrayFilter("nicknames", OPERATOR_ARRAY_CONTAINS_ANY, []any{"Doe", "John"}),
			}},
		},
		{
			desc:  "query with multiple filters",
			input: `QUERY user WHERE age = 20 AND name = "John Doe"`,
			want: &QueryOperation{collection: "user", filters: []Filter{
				NewIntFilter("age", OPERATOR_EQ, 20),
				NewStringFilter("name", OPERATOR_EQ, "John Doe"),
			}},
		},
		{
			desc:  "query with trim head slash",
			input: `QUERY /user`,
			want:  &QueryOperation{collection: "user"},
		},
		{
			desc:  "query with select",
			input: `QUERY user SELECT name, age`,
			want:  &QueryOperation{collection: "user", selects: []string{"name", "age"}},
		},
		{
			desc:  "query with select and where",
			input: `QUERY user SELECT name, age WHERE age = 20`,
			want:  &QueryOperation{collection: "user", selects: []string{"name", "age"}, filters: []Filter{NewIntFilter("age", OPERATOR_EQ, 20)}},
		},
		{
			desc:  "query with order by",
			input: `QUERY user ORDER BY age ASC`,
			want:  &QueryOperation{collection: "user", orderBys: []OrderBy{{"age", firestore.Asc}}},
		},
		{
			desc:  "query with multiple order by",
			input: `QUERY user ORDER BY age ASC, name DESC`,
			want:  &QueryOperation{collection: "user", orderBys: []OrderBy{{"age", firestore.Asc}, {"name", firestore.Desc}}},
		},
		{
			desc:  "query with order by and select and where",
			input: `QUERY user SELECT name, age WHERE age = 20 ORDER BY age ASC`,
			want: &QueryOperation{collection: "user", selects: []string{"name", "age"}, filters: []Filter{
				NewIntFilter("age", OPERATOR_EQ, 20),
			}, orderBys: []OrderBy{{"age", firestore.Asc}}},
		},
		{
			desc:  "query with limit",
			input: `QUERY user LIMIT 10`,
			want:  &QueryOperation{collection: "user", limit: 10},
		},
		{
			desc:  "query with limit and order by and select and where",
			input: `QUERY user SELECT name, age WHERE age = 20 ORDER BY age LIMIT 10`,
			want: &QueryOperation{collection: "user", selects: []string{"name", "age"}, filters: []Filter{
				NewIntFilter("age", OPERATOR_EQ, 20),
			}, orderBys: []OrderBy{{"age", firestore.Asc}}, limit: 10},
		},
		{
			desc:  "query with limit and order by desc and select and where",
			input: `QUERY user SELECT name, age WHERE age = 20 ORDER BY age DESC LIMIT 10`,
			want: &QueryOperation{collection: "user", selects: []string{"name", "age"}, filters: []Filter{
				NewIntFilter("age", OPERATOR_EQ, 20),
			}, orderBys: []OrderBy{{"age", firestore.Desc}}, limit: 10},
		},
		{
			desc:  "count",
			input: `COUNT user WHERE age = 20`,
			want: &CountOperation{collection: "user", filters: []Filter{
				NewIntFilter("age", OPERATOR_EQ, 20),
			}},
		},
		{
			desc:  "count with multiple filters",
			input: `COUNT user WHERE age = 20 AND name = "John Doe"`,
			want: &CountOperation{collection: "user", filters: []Filter{
				NewIntFilter("age", OPERATOR_EQ, 20),
				NewStringFilter("name", OPERATOR_EQ, "John Doe"),
			}},
		},
		{
			desc:  "list collections",
			input: `\d`,
			want:  &MetacommandListCollections{},
		},
		{
			desc:  "list subcollections",
			input: `\d user/1`,
			want: &MetacommandListCollections{
				baseDoc: "user/1",
			},
		},
		{
			desc:  "pager on",
			input: `\pager on`,
			want: &MetacommandPager{
				on: true,
			},
		},
		{
			desc:  "pager off",
			input: `\pager off`,
			want: &MetacommandPager{
				on: false,
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
