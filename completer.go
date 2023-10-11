package fscli

import (
	"github.com/c-bata/go-prompt"
)

var rootSuggestions = []prompt.Suggest{
	{Text: "GET", Description: "GET [docPath]"},
	{Text: "QUERY", Description: "QUERY [collection]"},
}

var querySuggestions = []prompt.Suggest{
	{Text: "SELECT", Description: "SELECT [field...]"},
	{Text: "WHERE", Description: "WHERE [field] [operator] [value]"},
	{Text: "ORDER BY", Description: "ORDER BY [field] [ASC/DESC]"},
	{Text: "LIMIT", Description: "LIMIT [count]"},
}

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
	return rootSuggestions, nil
}

func (c *Completer) parseQueryOperation() ([]prompt.Suggest, error) {
	if !c.expectPeek(IDENT) {
		return []prompt.Suggest{}, nil
	}

	// collection := c.curToken.Literal

	c.nextToken()

	if c.curTokenIs(EOF) {
		return querySuggestions, nil
	}

	if c.curTokenIs(SELECT) {
		c.nextToken()

		if c.curTokenIs(EOF) {
			return []prompt.Suggest{}, nil
		}

		for !c.curTokenIs(EOF) && !c.curTokenIs(WHERE) {
			c.nextToken()
		}
	}

	if c.curTokenIs(EOF) {
		// without select
		return querySuggestions[1:], nil
	}

	if c.curTokenIs(WHERE) {
		c.nextToken()

		if c.curTokenIs(EOF) {
			return []prompt.Suggest{}, nil
		}

		for !c.curTokenIs(EOF) && !c.curTokenIs(ORDER) {
			c.nextToken()
		}
	}

	if c.curTokenIs(EOF) {
		// without where
		return querySuggestions[2:], nil
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

		for !c.curTokenIs(EOF) && !c.curTokenIs(LIMIT) {
			c.nextToken()
		}
	}

	if c.curTokenIs(EOF) {
		// without order by
		return querySuggestions[3:], nil
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
