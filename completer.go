package fscli

import (
	"fmt"
	"strings"

	"github.com/c-bata/go-prompt"
)

var (
	getSuggestion   = prompt.Suggest{Text: "GET", Description: "GET [docPath]"}
	querySuggestion = prompt.Suggest{Text: "QUERY", Description: "QUERY [collection]"}
)

var rootSuggestions = []prompt.Suggest{
	getSuggestion,
	querySuggestion,
}

var (
	selectSuggestion  = prompt.Suggest{Text: "SELECT", Description: "SELECT [field...]"}
	whereSuggestion   = prompt.Suggest{Text: "WHERE", Description: "WHERE [field] [operator] [value]"}
	orderBySuggestion = prompt.Suggest{Text: "ORDER BY", Description: "ORDER BY [field] [ASC/DESC]"}
	limitSuggestion   = prompt.Suggest{Text: "LIMIT", Description: "LIMIT [count]"}
)

var querySuggestions = []prompt.Suggest{
	selectSuggestion,
	whereSuggestion,
	orderBySuggestion,
	limitSuggestion,
}

var (
	ascSuggestion  = prompt.Suggest{Text: "ASC", Description: "ASC"}
	descSuggestion = prompt.Suggest{Text: "DESC", Description: "DESC"}
)

func newCollectionSuggestion(baseDoc string, name string) prompt.Suggest {
	var t string
	if baseDoc == "" {
		t = name
	} else {
		t = fmt.Sprintf("%s/%s", baseDoc, name)
	}
	return prompt.Suggest{Text: t, Description: name}
}

type Completer struct {
	l               *Lexer
	curToken        Token
	peekToken       Token
	findCollections func(baseDoc string) ([]string, error)
}

func NewCompleter(l *Lexer, findCollections func(baseDoc string) ([]string, error)) *Completer {
	c := &Completer{l: l, findCollections: findCollections}

	// Set both curToken and peekToken
	c.nextToken()
	c.nextToken()

	return c
}

func (c *Completer) Parse() ([]prompt.Suggest, error) {
	if c.curTokenIs(QUERY) {
		return c.parseQueryOperation()
	}
	if c.curTokenIs(GET) {
		return c.parseGetOperation()
	}

	if c.curTokenIs(IDENT) {
		return prompt.FilterHasPrefix(rootSuggestions, c.curToken.Literal, true), nil
	}
	return []prompt.Suggest{}, nil
}

func (c *Completer) parseQueryOperation() ([]prompt.Suggest, error) {
	if !c.expectPeek(IDENT) {
		return []prompt.Suggest{}, nil
	}

	// collection := c.curToken.Literal
	if c.peekTokenIs(EOF) {
		collection := normalizeFirestorePath(c.curToken.Literal)
		parts := strings.Split(collection, "/")
		if len(parts)%2 == 0 {
			return []prompt.Suggest{}, nil
		}

		var baseDoc string
		if len(parts) == 1 {
			baseDoc = ""
		} else {
			baseDoc = strings.Join(parts[:len(parts)-1], "/")
		}
		collections, err := c.findCollections(baseDoc)
		if err != nil {
			return []prompt.Suggest{}, nil
		}
		suggestions := make([]prompt.Suggest, 0, len(collections))
		for _, col := range collections {
			suggestions = append(suggestions, newCollectionSuggestion(baseDoc, col))
		}

		return prompt.FilterHasPrefix(suggestions, c.curToken.Literal, false), nil
	}

	c.nextToken()

	if c.curTokenIs(EOF) {
		return []prompt.Suggest{}, nil
	}

	// select / where / order by / limit
	if c.curTokenIs(IDENT) {
		return prompt.FilterHasPrefix(querySuggestions, c.curToken.Literal, true), nil
	}

	if c.curTokenIs(SELECT) {
		c.nextToken()

		if c.curTokenIs(EOF) {
			return []prompt.Suggest{}, nil
		}

		for !c.curTokenIs(EOF) {
			if c.curTokenIs(COMMA) {
				c.nextToken()
			} else {
				c.nextToken()
				break
			}
		}
	}

	if c.curTokenIs(EOF) {
		return []prompt.Suggest{}, nil
	}
	// where / order by / limit
	if c.curTokenIs(IDENT) {
		return prompt.FilterHasPrefix(querySuggestions[1:], c.curToken.Literal, true), nil
	}

	if c.curTokenIs(WHERE) {
		c.nextToken()

		if c.curTokenIs(EOF) {
			return []prompt.Suggest{}, nil
		}

		c.nextToken()
		for !c.curTokenIs(EOF) {
			if c.curTokenIsOperator() {
				c.nextToken()
				// value

				if c.curTokenIs(EOF) {
					return []prompt.Suggest{}, nil
				}
				c.nextToken()
			} else if c.curTokenIs(AND) {
				c.nextToken()
				if c.curTokenIs(EOF) {
					return []prompt.Suggest{}, nil
				}
				c.nextToken()
			} else {
				break
			}
		}
	}

	if c.curTokenIs(EOF) {
		return []prompt.Suggest{}, nil
	}
	// order by / limit
	if c.curTokenIs(IDENT) {
		return prompt.FilterHasPrefix(querySuggestions[2:], c.curToken.Literal, true), nil
	}

	if c.curTokenIs(ORDER) {
		c.nextToken()

		if c.curTokenIs(EOF) {
			return []prompt.Suggest{}, nil
		}
		if c.curTokenIs(BY) {
			c.nextToken()
		}

		if c.curTokenIs(EOF) {
			return []prompt.Suggest{}, nil
		}

		c.nextToken()

		for !c.curTokenIs(EOF) {
			if c.curTokenIs(ASC) || c.curTokenIs(DESC) {
				c.nextToken()
			}
			if c.curTokenIs(COMMA) {
				c.nextToken()
				if c.curTokenIs(EOF) {
					return []prompt.Suggest{}, nil
				}
				c.nextToken()
			} else {
				break
			}
		}
	}

	if c.curTokenIs(EOF) {
		return []prompt.Suggest{}, nil
	}
	// asc/ desc / limit
	if c.curTokenIs(IDENT) {
		return prompt.FilterHasPrefix([]prompt.Suggest{ascSuggestion, descSuggestion, limitSuggestion}, c.curToken.Literal, true), nil
	}

	return []prompt.Suggest{}, nil
}

func (c *Completer) parseGetOperation() ([]prompt.Suggest, error) {
	if !c.expectPeek(IDENT) {
		return []prompt.Suggest{}, nil
	}
	if c.peekTokenIs(EOF) {
		docPath := normalizeFirestorePath(c.curToken.Literal)
		parts := strings.Split(docPath, "/")
		if len(parts)%2 == 0 {
			return []prompt.Suggest{}, nil
		}
		var baseDoc string
		if len(parts) == 1 {
			baseDoc = ""
		} else {
			baseDoc = strings.Join(parts[:len(parts)-1], "/")
		}

		collections, err := c.findCollections(baseDoc)
		if err != nil {
			return []prompt.Suggest{}, nil
		}
		suggestions := make([]prompt.Suggest, 0, len(collections))
		for _, col := range collections {
			suggestions = append(suggestions, newCollectionSuggestion(baseDoc, col))
		}

		return prompt.FilterHasPrefix(suggestions, c.curToken.Literal, false), nil
	}

	return []prompt.Suggest{}, nil
}

func (c *Completer) nextToken() {
	c.curToken = c.peekToken
	c.peekToken = c.l.NextToken()
}

func (c *Completer) curTokenIs(t TokenType) bool {
	return c.curToken.Type == t
}

func (c *Completer) curTokenIsOperator() bool {
	return c.curTokenIs(EQ) ||
		c.curTokenIs(NOT_EQ) ||
		c.curTokenIs(GT) ||
		c.curTokenIs(GTE) ||
		c.curTokenIs(LT) ||
		c.curTokenIs(LTE) ||
		c.curTokenIs(IN) ||
		c.curTokenIs(ARRAY_CONTAINS) ||
		c.curTokenIs(ARRAY_CONTAINS_ANY)
}

func (c *Completer) peekTokenIs(t TokenType) bool {
	return c.peekToken.Type == t
}

func (c *Completer) expectPeek(t TokenType) bool {
	if c.peekTokenIs(t) {
		c.nextToken()
		return true
	} else {
		return false
	}
}
