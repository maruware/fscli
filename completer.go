package fscli

import (
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

type Completer struct {
	l         *Lexer
	curToken  Token
	peekToken Token
}

func NewCompleter(l *Lexer) *Completer {
	c := &Completer{l: l}

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
		return []prompt.Suggest{}, nil
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
			if c.curTokenIs(EQ) || c.curTokenIs(NOT_EQ) {
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

func (c *Completer) nextToken() {
	c.curToken = c.peekToken
	c.peekToken = c.l.NextToken()
}

func (c *Completer) curTokenIs(t TokenType) bool {
	return c.curToken.Type == t
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
