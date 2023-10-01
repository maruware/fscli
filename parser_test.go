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
					&StringFilter{
						field:    "name",
						operator: "==",
						value:    "John Doe",
					},
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
