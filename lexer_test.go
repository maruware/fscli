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
			input: `QUERY user WHERE name == "John Doe"`,
			want: []Token{
				{Type: QUERY, Literal: "QUERY"},
				{Type: IDENT, Literal: "user"},
				{Type: WHERE, Literal: "WHERE"},
				{Type: IDENT, Literal: "name"},
				{Type: EQ, Literal: "=="},
				{Type: STRING, Literal: "John Doe"},
			},
		},
		{
			desc:  "query with single quote",
			input: `QUERY user WHERE name == 'John Doe'`,
			want: []Token{
				{Type: QUERY, Literal: "QUERY"},
				{Type: IDENT, Literal: "user"},
				{Type: WHERE, Literal: "WHERE"},
				{Type: IDENT, Literal: "name"},
				{Type: EQ, Literal: "=="},
				{Type: STRING, Literal: "John Doe"},
			},
		},
		{
			desc:  "query with int",
			input: `QUERY user WHERE age == 20`,
			want: []Token{
				{Type: QUERY, Literal: "QUERY"},
				{Type: IDENT, Literal: "user"},
				{Type: WHERE, Literal: "WHERE"},
				{Type: IDENT, Literal: "age"},
				{Type: EQ, Literal: "=="},
				{Type: INT, Literal: "20"},
			},
		},
		{
			desc:  "query with float",
			input: `QUERY user WHERE age == 20.5`,
			want: []Token{
				{Type: QUERY, Literal: "QUERY"},
				{Type: IDENT, Literal: "user"},
				{Type: WHERE, Literal: "WHERE"},
				{Type: IDENT, Literal: "age"},
				{Type: EQ, Literal: "=="},
				{Type: FLOAT, Literal: "20.5"},
			},
		},
		{
			desc:  "query subcollection",
			input: `QUERY users/abc/posts WHERE title == "Hello World"`,
			want: []Token{
				{Type: QUERY, Literal: "QUERY"},
				{Type: IDENT, Literal: "users/abc/posts"},
				{Type: WHERE, Literal: "WHERE"},
				{Type: IDENT, Literal: "title"},
				{Type: EQ, Literal: "=="},
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
			input: `query user where name == 'John Doe'`,
			want: []Token{
				{Type: QUERY, Literal: "query"},
				{Type: IDENT, Literal: "user"},
				{Type: WHERE, Literal: "where"},
				{Type: IDENT, Literal: "name"},
				{Type: EQ, Literal: "=="},
				{Type: STRING, Literal: "John Doe"},
			},
		},
		{
			desc:  "query with lower and upper case",
			input: `QuEry user WhERE name == 'John Doe'`,
			want: []Token{
				{Type: QUERY, Literal: "QuEry"},
				{Type: IDENT, Literal: "user"},
				{Type: WHERE, Literal: "WhERE"},
				{Type: IDENT, Literal: "name"},
				{Type: EQ, Literal: "=="},
				{Type: STRING, Literal: "John Doe"},
			},
		},
		{
			desc:  "query with and",
			input: `QUERY user WHERE name == "John Doe" AND age == 20`,
			want: []Token{
				{Type: QUERY, Literal: "QUERY"},
				{Type: IDENT, Literal: "user"},
				{Type: WHERE, Literal: "WHERE"},
				{Type: IDENT, Literal: "name"},
				{Type: EQ, Literal: "=="},
				{Type: STRING, Literal: "John Doe"},
				{Type: AND, Literal: "AND"},
				{Type: IDENT, Literal: "age"},
				{Type: EQ, Literal: "=="},
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
