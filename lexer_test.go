package fscli

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLexer(t *testing.T) {
	tests := []struct {
		desc  string
		input string
		want  []Token
	}{
		{
			desc:  "simple query",
			input: `QUERY user WHERE name = "John Doe"`,
			want: []Token{
				{Type: QUERY, Literal: "QUERY"},
				{Type: IDENT, Literal: "user"},
				{Type: WHERE, Literal: "WHERE"},
				{Type: IDENT, Literal: "name"},
				{Type: EQ, Literal: "="},
				{Type: STRING, Literal: "John Doe"},
			},
		},
		{
			desc:  "query with single quote",
			input: `QUERY user WHERE name = 'John Doe'`,
			want: []Token{
				{Type: QUERY, Literal: "QUERY"},
				{Type: IDENT, Literal: "user"},
				{Type: WHERE, Literal: "WHERE"},
				{Type: IDENT, Literal: "name"},
				{Type: EQ, Literal: "="},
				{Type: STRING, Literal: "John Doe"},
			},
		},
		{
			desc:  "query with int",
			input: `QUERY user WHERE age = 20`,
			want: []Token{
				{Type: QUERY, Literal: "QUERY"},
				{Type: IDENT, Literal: "user"},
				{Type: WHERE, Literal: "WHERE"},
				{Type: IDENT, Literal: "age"},
				{Type: EQ, Literal: "="},
				{Type: INT, Literal: "20"},
			},
		},
		{
			desc:  "query with float",
			input: `QUERY user WHERE age = 20.5`,
			want: []Token{
				{Type: QUERY, Literal: "QUERY"},
				{Type: IDENT, Literal: "user"},
				{Type: WHERE, Literal: "WHERE"},
				{Type: IDENT, Literal: "age"},
				{Type: EQ, Literal: "="},
				{Type: FLOAT, Literal: "20.5"},
			},
		},
		{
			desc:  "query subcollection",
			input: `QUERY users/abc/posts WHERE title = "Hello World"`,
			want: []Token{
				{Type: QUERY, Literal: "QUERY"},
				{Type: IDENT, Literal: "users/abc/posts"},
				{Type: WHERE, Literal: "WHERE"},
				{Type: IDENT, Literal: "title"},
				{Type: EQ, Literal: "="},
				{Type: STRING, Literal: "Hello World"},
			},
		},
		{
			desc:  "query with not equal",
			input: `QUERY user WHERE name != "John Doe"`,
			want: []Token{
				{Type: QUERY, Literal: "QUERY"},
				{Type: IDENT, Literal: "user"},
				{Type: WHERE, Literal: "WHERE"},
				{Type: IDENT, Literal: "name"},
				{Type: NOT_EQ, Literal: "!="},
				{Type: STRING, Literal: "John Doe"},
			},
		},
		{
			desc:  "query with lower case",
			input: `query user where name = 'John Doe'`,
			want: []Token{
				{Type: QUERY, Literal: "query"},
				{Type: IDENT, Literal: "user"},
				{Type: WHERE, Literal: "where"},
				{Type: IDENT, Literal: "name"},
				{Type: EQ, Literal: "="},
				{Type: STRING, Literal: "John Doe"},
			},
		},
		{
			desc:  "query with lower and upper case",
			input: `QuEry user WhERE name = 'John Doe'`,
			want: []Token{
				{Type: QUERY, Literal: "QuEry"},
				{Type: IDENT, Literal: "user"},
				{Type: WHERE, Literal: "WhERE"},
				{Type: IDENT, Literal: "name"},
				{Type: EQ, Literal: "="},
				{Type: STRING, Literal: "John Doe"},
			},
		},
		{
			desc:  "query with and",
			input: `QUERY user WHERE name = "John Doe" AND age = 20`,
			want: []Token{
				{Type: QUERY, Literal: "QUERY"},
				{Type: IDENT, Literal: "user"},
				{Type: WHERE, Literal: "WHERE"},
				{Type: IDENT, Literal: "name"},
				{Type: EQ, Literal: "="},
				{Type: STRING, Literal: "John Doe"},
				{Type: AND, Literal: "AND"},
				{Type: IDENT, Literal: "age"},
				{Type: EQ, Literal: "="},
				{Type: INT, Literal: "20"},
			},
		},
		{
			desc:  "query with in",
			input: `QUERY user WHERE name IN ["John Doe", "Jane Doe"]`,
			want: []Token{
				{Type: QUERY, Literal: "QUERY"},
				{Type: IDENT, Literal: "user"},
				{Type: WHERE, Literal: "WHERE"},
				{Type: IDENT, Literal: "name"},
				{Type: IN, Literal: "IN"},
				{Type: LBRACKET, Literal: "["},
				{Type: STRING, Literal: "John Doe"},
				{Type: COMMA, Literal: ","},
				{Type: STRING, Literal: "Jane Doe"},
				{Type: RBRACKET, Literal: "]"},
			},
		},
		{
			desc:  "query with GT",
			input: `QUERY user WHERE age > 20`,
			want: []Token{
				{Type: QUERY, Literal: "QUERY"},
				{Type: IDENT, Literal: "user"},
				{Type: WHERE, Literal: "WHERE"},
				{Type: IDENT, Literal: "age"},
				{Type: GT, Literal: ">"},
				{Type: INT, Literal: "20"},
			},
		},
		{
			desc:  "query with GTE",
			input: `QUERY user WHERE age >= 20`,
			want: []Token{
				{Type: QUERY, Literal: "QUERY"},
				{Type: IDENT, Literal: "user"},
				{Type: WHERE, Literal: "WHERE"},
				{Type: IDENT, Literal: "age"},
				{Type: GTE, Literal: ">="},
				{Type: INT, Literal: "20"},
			},
		},
		{
			desc:  "query with LT",
			input: `QUERY user WHERE age < 20`,
			want: []Token{
				{Type: QUERY, Literal: "QUERY"},
				{Type: IDENT, Literal: "user"},
				{Type: WHERE, Literal: "WHERE"},
				{Type: IDENT, Literal: "age"},
				{Type: LT, Literal: "<"},
				{Type: INT, Literal: "20"},
			},
		},
		{
			desc:  "query with LTE",
			input: `QUERY user WHERE age <= 20`,
			want: []Token{
				{Type: QUERY, Literal: "QUERY"},
				{Type: IDENT, Literal: "user"},
				{Type: WHERE, Literal: "WHERE"},
				{Type: IDENT, Literal: "age"},
				{Type: LTE, Literal: "<="},
				{Type: INT, Literal: "20"},
			},
		},
		{
			desc:  "query with array-contains",
			input: `QUERY user WHERE nicknames ARRAY_CONTAINS "Doe"`,
			want: []Token{
				{Type: QUERY, Literal: "QUERY"},
				{Type: IDENT, Literal: "user"},
				{Type: WHERE, Literal: "WHERE"},
				{Type: IDENT, Literal: "nicknames"},
				{Type: ARRAY_CONTAINS, Literal: "ARRAY_CONTAINS"},
				{Type: STRING, Literal: "Doe"},
			},
		},
		{
			desc:  "query with array-contains-any",
			input: `QUERY user WHERE nicknames ARRAY_CONTAINS_ANY ["Doe", "John"]`,
			want: []Token{
				{Type: QUERY, Literal: "QUERY"},
				{Type: IDENT, Literal: "user"},
				{Type: WHERE, Literal: "WHERE"},
				{Type: IDENT, Literal: "nicknames"},
				{Type: ARRAY_CONTAINS_ANY, Literal: "ARRAY_CONTAINS_ANY"},
				{Type: LBRACKET, Literal: "["},
				{Type: STRING, Literal: "Doe"},
				{Type: COMMA, Literal: ","},
				{Type: STRING, Literal: "John"},
				{Type: RBRACKET, Literal: "]"},
			},
		},
		{
			desc:  "query with timestamp",
			input: `QUERY user WHERE created_at > TIMESTAMP("2006-01-02T15:04:05+09:00")`,
			want: []Token{
				{Type: QUERY, Literal: "QUERY"},
				{Type: IDENT, Literal: "user"},
				{Type: WHERE, Literal: "WHERE"},
				{Type: IDENT, Literal: "created_at"},
				{Type: GT, Literal: ">"},
				{Type: IDENT, Literal: "TIMESTAMP"},
				{Type: LPAREN, Literal: "("},
				{Type: STRING, Literal: "2006-01-02T15:04:05+09:00"},
				{Type: RPAREN, Literal: ")"},
			},
		},
		{
			desc:  "get",
			input: `GET users/abc`,
			want: []Token{
				{Type: GET, Literal: "GET"},
				{Type: IDENT, Literal: "users/abc"},
			},
		},
		{
			desc:  "get with letter and digit",
			input: `GET users/abc123`,
			want: []Token{
				{Type: GET, Literal: "GET"},
				{Type: IDENT, Literal: "users/abc123"},
			},
		},
		{
			desc:  "query with emoji",
			input: `QUERY users WHERE name = "👍"`,
			want: []Token{
				{Type: QUERY, Literal: "QUERY"},
				{Type: IDENT, Literal: "users"},
				{Type: WHERE, Literal: "WHERE"},
				{Type: IDENT, Literal: "name"},
				{Type: EQ, Literal: "="},
				{Type: STRING, Literal: "👍"},
			},
		},
		{
			desc:  "query with select",
			input: `QUERY users SELECT name, age`,
			want: []Token{
				{Type: QUERY, Literal: "QUERY"},
				{Type: IDENT, Literal: "users"},
				{Type: SELECT, Literal: "SELECT"},
				{Type: IDENT, Literal: "name"},
				{Type: COMMA, Literal: ","},
				{Type: IDENT, Literal: "age"},
			},
		},
		{
			desc:  "query with uuid path",
			input: `QUERY users/b92ac3fc-8045-41d9-b1a6-9d0a033757c6/posts WHERE title = "Hello World"`,
			want: []Token{
				{Type: QUERY, Literal: "QUERY"},
				{Type: IDENT, Literal: "users/b92ac3fc-8045-41d9-b1a6-9d0a033757c6/posts"},
				{Type: WHERE, Literal: "WHERE"},
				{Type: IDENT, Literal: "title"},
				{Type: EQ, Literal: "="},
				{Type: STRING, Literal: "Hello World"},
			},
		},
		{
			desc:  "query with order by",
			input: `QUERY users ORDER BY name ASC`,
			want: []Token{
				{Type: QUERY, Literal: "QUERY"},
				{Type: IDENT, Literal: "users"},
				{Type: ORDER, Literal: "ORDER"},
				{Type: BY, Literal: "BY"},
				{Type: IDENT, Literal: "name"},
				{Type: ASC, Literal: "ASC"},
			},
		},
		{
			desc:  "query with order by desc",
			input: `QUERY users ORDER BY name DESC`,
			want: []Token{
				{Type: QUERY, Literal: "QUERY"},
				{Type: IDENT, Literal: "users"},
				{Type: ORDER, Literal: "ORDER"},
				{Type: BY, Literal: "BY"},
				{Type: IDENT, Literal: "name"},
				{Type: DESC, Literal: "DESC"},
			},
		},
		{
			desc:  "query with multiple order by",
			input: `QUERY users ORDER BY name ASC, age DESC`,
			want: []Token{
				{Type: QUERY, Literal: "QUERY"},
				{Type: IDENT, Literal: "users"},
				{Type: ORDER, Literal: "ORDER"},
				{Type: BY, Literal: "BY"},
				{Type: IDENT, Literal: "name"},
				{Type: ASC, Literal: "ASC"},
				{Type: COMMA, Literal: ","},
				{Type: IDENT, Literal: "age"},
				{Type: DESC, Literal: "DESC"},
			},
		},
		{
			desc:  "query with limit",
			input: `QUERY users LIMIT 10`,
			want: []Token{
				{Type: QUERY, Literal: "QUERY"},
				{Type: IDENT, Literal: "users"},
				{Type: LIMIT, Literal: "LIMIT"},
				{Type: INT, Literal: "10"},
			},
		},
		{
			desc:  "query with limit and order by and select and where",
			input: `QUERY users SELECT name, age WHERE name = "John Doe" ORDER BY name ASC, age DESC LIMIT 10`,
			want: []Token{
				{Type: QUERY, Literal: "QUERY"},
				{Type: IDENT, Literal: "users"},
				{Type: SELECT, Literal: "SELECT"},
				{Type: IDENT, Literal: "name"},
				{Type: COMMA, Literal: ","},
				{Type: IDENT, Literal: "age"},
				{Type: WHERE, Literal: "WHERE"},
				{Type: IDENT, Literal: "name"},
				{Type: EQ, Literal: "="},
				{Type: STRING, Literal: "John Doe"},
				{Type: ORDER, Literal: "ORDER"},
				{Type: BY, Literal: "BY"},
				{Type: IDENT, Literal: "name"},
				{Type: ASC, Literal: "ASC"},
				{Type: COMMA, Literal: ","},
				{Type: IDENT, Literal: "age"},
				{Type: DESC, Literal: "DESC"},
				{Type: LIMIT, Literal: "LIMIT"},
				{Type: INT, Literal: "10"},
			},
		},
		{
			desc:  "count",
			input: `COUNT users WHERE name = "John Doe"`,
			want: []Token{
				{Type: COUNT, Literal: "COUNT"},
				{Type: IDENT, Literal: "users"},
				{Type: WHERE, Literal: "WHERE"},
				{Type: IDENT, Literal: "name"},
				{Type: EQ, Literal: "="},
				{Type: STRING, Literal: "John Doe"},
			},
		},
		{
			desc:  "list collections",
			input: `\d`,
			want: []Token{
				{Type: LIST_COLLECTIONS, Literal: `\d`},
			},
		},
		{
			desc:  "pager on",
			input: `\pager on`,
			want: []Token{
				{Type: PAGER, Literal: `\pager`},
				{Type: IDENT, Literal: "on"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			l := NewLexer(tt.input)
			tokens := []Token{}
			for {
				tok := l.NextToken()
				if tok.Type == EOF {
					break
				}
				tokens = append(tokens, tok)
			}

			assert.Equal(t, tt.want, tokens)
		})
	}
}
