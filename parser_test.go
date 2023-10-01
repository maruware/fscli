package fscli

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
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
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			got, err := Parse(tt.input)
			if err != nil {
				t.Fatalf("Parse(%q) returned error: %v", tt.input, err)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
