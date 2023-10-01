package fscli

import (
	"fmt"
	"strconv"
)

type Parser struct {
	l         *Lexer
	curToken  Token
	peekToken Token
	errors    []string
}

func NewParser(l *Lexer) *Parser {
	p := &Parser{l: l, errors: []string{}}

	// Set both curToken and peekToken
	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) Parse() Operation {

	if p.curTokenIs(QUERY) {
		return p.parseQueryOperation()
	}
	return nil
}

func (p *Parser) parseQueryOperation() *QueryOperation {
	op := &QueryOperation{opType: QUERY}

	if !p.expectPeek(IDENT) {
		return nil
	}

	op.collection = p.curToken.Literal

	if !p.expectPeek(WHERE) {
		return nil
	}

	p.nextToken()
	for !p.curTokenIs(EOF) {
		filter, err := p.parseFilter()
		if err != nil {
			p.errors = append(p.errors, err.Error())
			return nil
		}
		if filter != nil {
			op.filters = append(op.filters, filter)
		}
		p.nextToken()
	}

	return op
}

func (p *Parser) parseFilter() (Filter, error) {
	field := p.curToken.Literal

	p.nextToken()
	operator := p.curToken.Literal

	p.nextToken()

	if p.curTokenIs(INT) {
		n, err := strconv.Atoi(p.curToken.Literal)
		if err != nil {
			return nil, fmt.Errorf("invalid int value: %s", p.curToken.Literal)
		}
		return NewIntFilter(field, operator, n), nil
	}
	if p.curTokenIs(FLOAT) {
		n, err := strconv.ParseFloat(p.curToken.Literal, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid float value: %s", p.curToken.Literal)
		}
		return NewFloatFilter(field, operator, n), nil
	}
	if p.curTokenIs(STRING) {
		return NewStringFilter(field, operator, p.curToken.Literal), nil
	}
	return nil, fmt.Errorf("invalid filter value: %s", p.curToken.Literal)
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) curTokenIs(t TokenType) bool {
	return p.curToken.Type == t
}
func (p *Parser) peekTokenIs(t TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) expectPeek(t TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		p.peekError(t)
		return false
	}
}

func (p *Parser) peekError(t TokenType) {
	msg := fmt.Sprintf("expected next token tobe %s, got %s instead", t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}
