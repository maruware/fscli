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
