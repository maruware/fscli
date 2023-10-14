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
			want:  []prompt.Suggest{},
		},
		{
			desc:  "middle of query",
			input: `QUE`,
			want:  []prompt.Suggest{querySuggestion},
		},
		{
			desc:  "middle of get",
			input: `GE`,
			want:  []prompt.Suggest{getSuggestion},
		},
		{
			desc:  "query",
			input: `QUERY`,
			want:  []prompt.Suggest{},
		},
		{
			desc:  "get",
			input: `GET`,
			want:  []prompt.Suggest{},
		},
		{
			desc:  "middle of query with collection",
			input: `QUERY us`,
			want:  []prompt.Suggest{newCollectionSuggestion("", "user")},
		},
		{
			desc:  "middle of query with sub collection",
			input: `QUERY user/1/p`,
			want:  []prompt.Suggest{newCollectionSuggestion("user/1", "posts")},
		},
		{
			desc:  "middle of query with select",
			input: `QUERY user S`,
			want:  []prompt.Suggest{selectSuggestion},
		},
		{
			desc:  "query with select and no field",
			input: `QUERY user SELECT `,
			want:  []prompt.Suggest{},
		},
		{
			desc:  "query with select and field",
			input: `QUERY user SELECT name`,
			want:  []prompt.Suggest{},
		},
		{
			desc:  "middle of where",
			input: `QUERY user SELECT name W`,
			want:  []prompt.Suggest{whereSuggestion},
		},
		{
			desc:  "query with select and field and where",
			input: `QUERY user SELECT name WHERE`,
			want:  []prompt.Suggest{},
		},
		{
			desc:  "query with select and field and where and field",
			input: `QUERY user SELECT name WHERE name `,
			// TODO: should return operators
			want: []prompt.Suggest{},
		},
		{
			desc:  "query with select and field and where and field and operator",
			input: `QUERY user SELECT name WHERE name ==`,
			want:  []prompt.Suggest{},
		},
		{
			desc:  "middle of order by",
			input: `QUERY user ORD`,
			want:  []prompt.Suggest{orderBySuggestion},
		},
		{
			desc:  "middle of order by after where",
			input: `QUERY user SELECT name WHERE name = "Doe" ORD`,
			want:  []prompt.Suggest{orderBySuggestion},
		},
		{
			desc:  "query with order by",
			input: `QUERY user ORDER BY`,
			want:  []prompt.Suggest{},
		},
		{
			desc:  "query with order by and field",
			input: `QUERY user ORDER BY name `,
			want:  []prompt.Suggest{},
		},
		{
			desc:  "middle of asc",
			input: `QUERY user ORDER BY name A`,
			want:  []prompt.Suggest{ascSuggestion},
		},
		{
			desc:  "middle of desc",
			input: `QUERY user ORDER BY name D`,
			want:  []prompt.Suggest{descSuggestion},
		},
		{
			desc:  "middle of limit",
			input: `QUERY user ORDER BY name LI`,
			want:  []prompt.Suggest{limitSuggestion},
		},
		{
			desc:  "middle of limit after ASC",
			input: `QUERY user ORDER BY name ASC LI`,
			want:  []prompt.Suggest{limitSuggestion},
		},
		{
			desc:  "middle of limit after multiple order by",
			input: `QUERY user ORDER BY name ASC, age DESC LI`,
			want:  []prompt.Suggest{limitSuggestion},
		},
	}

	findCollections := func(baseDoc string) ([]string, error) {
		if baseDoc == "" {
			return []string{"user", "group"}, nil
		}
		if baseDoc == "user/1" {
			return []string{"posts"}, nil
		}
		return []string{}, nil
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			l := NewLexer(tt.input)
			c := NewCompleter(l, findCollections)
			got, err := c.Parse()
			if err != nil {
				t.Fatal(err)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
