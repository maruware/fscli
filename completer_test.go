package fscli

import (
	"testing"

	"github.com/c-bata/go-prompt"
	"github.com/stretchr/testify/assert"
)

func TestCompleter(t *testing.T) {
	tests := []struct {
		desc  string
		input string
		want  []prompt.Suggest
	}{
		{
			desc:  "root",
			input: ``,
			want:  rootSuggestions,
		},
		{
			desc:  "query",
			input: `QUERY`,
			want:  []prompt.Suggest{},
		},
		{
			desc:  "query with collection",
			input: `QUERY user`,
			want:  querySuggestions,
		},
		{
			desc:  "query with select and no field",
			input: `QUERY user SELECT`,
			want:  []prompt.Suggest{},
		},
		{
			desc:  "query with select and field",
			input: `QUERY user SELECT name`,
			want:  querySuggestions[1:],
		},
		{
			desc:  "query with select and field and where",
			input: `QUERY user SELECT name WHERE`,
			want:  []prompt.Suggest{},
		},
		{
			desc:  "query with select and field and where and field",
			input: `QUERY user SELECT name WHERE name`,
			// TODO: should return operators
			want: querySuggestions[2:],
		},
		{
			desc:  "query with select and field and where and field and operator",
			input: `QUERY user SELECT name WHERE name ==`,
			// TODO: should return empty
			want: querySuggestions[2:],
		},
		{
			desc:  "query with order by",
			input: `QUERY user ORDER BY`,
			want:  []prompt.Suggest{},
		},
		{
			desc:  "query with order by and field",
			input: `QUERY user ORDER BY name`,
			want:  querySuggestions[3:],
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			l := NewLexer(tt.input)
			c := NewCompleter(l)
			got, err := c.Parse()
			if err != nil {
				t.Fatal(err)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
