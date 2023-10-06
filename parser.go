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

func (p *Parser) Parse() (Operation, error) {

	if p.curTokenIs(QUERY) {
		return p.parseQueryOperation()
	}
	return nil, fmt.Errorf("invalid")
}

func (p *Parser) parseQueryOperation() (*QueryOperation, error) {
	op := &QueryOperation{opType: QUERY}

	if !p.expectPeek(IDENT) {
		return nil, fmt.Errorf("invalid")
	}

	op.collection = p.curToken.Literal

	if !p.expectPeek(WHERE) {
		return nil, fmt.Errorf("invalid")
	}

	p.nextToken()
	for !p.curTokenIs(EOF) {
		filter, err := p.parseFilter()
		if err != nil {
			p.errors = append(p.errors, err.Error())
			return nil, err
		}
		if filter != nil {
			op.filters = append(op.filters, filter)
		}

		if !p.expectPeek(AND) {
			break
		}
		p.nextToken()
	}

	return op, nil
}

func (p *Parser) parseFilter() (Filter, error) {
	field := p.curToken.Literal

	p.nextToken()

	var operator Operator
	if p.curTokenIs(EQ) {
		operator = OPERATOR_EQ
	} else if p.curTokenIs(NOT_EQ) {
		operator = OPERATOR_NOT_EQ
	} else if p.curTokenIs(IN) {
		operator = OPERATOR_IN
	} else if p.curTokenIs(ARRAY_CONTAINS) {
		operator = OPERATOR_ARRAY_CONTAINS
	} else if p.curTokenIs(ARRAY_CONTAINS_ANY) {
		operator = OPERATOR_ARRAY_CONTAINS_ANY
	} else {
		return nil, fmt.Errorf("invalid filter operator: %s", p.curToken.Literal)
	}

	p.nextToken()

	if p.curTokenIs(LBRACKET) {
		return p.parseArrayFilter(field, operator)
	}

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

func (p *Parser) parseArrayFilter(field string, operator Operator) (Filter, error) {
	var values []any
	p.nextToken()
	for !p.curTokenIs(RBRACKET) {
		if p.curTokenIs(STRING) {
			values = append(values, p.curToken.Literal)
		} else if p.curTokenIs(INT) {
			n, err := strconv.Atoi(p.curToken.Literal)
			if err != nil {
				return nil, fmt.Errorf("invalid int value: %s", p.curToken.Literal)
			}
			values = append(values, n)
		} else if p.curTokenIs(FLOAT) {
			n, err := strconv.ParseFloat(p.curToken.Literal, 64)
			if err != nil {
				return nil, fmt.Errorf("invalid float value: %s", p.curToken.Literal)
			}
			values = append(values, n)
		} else {
			return nil, fmt.Errorf("invalid array filter value: %s", p.curToken.Literal)
		}

		if !p.expectPeek(COMMA) {
			break
		}
		p.nextToken()
	}
	return NewArrayFilter(field, operator, values), nil
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
