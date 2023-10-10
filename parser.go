package fscli

import (
	"fmt"
	"strconv"
	"strings"

	"cloud.google.com/go/firestore"
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
	if p.curTokenIs(GET) {
		return p.parseGetOperation()
	}
	return nil, fmt.Errorf("invalid operation: %s", p.curToken.Literal)
}

func (p *Parser) parseQueryOperation() (*QueryOperation, error) {
	op := &QueryOperation{}

	if !p.expectPeek(IDENT) {
		return nil, fmt.Errorf("invalid: expected collection but got %s", p.curToken.Literal)
	}

	op.collection = p.trimHeadSlash(p.curToken.Literal)

	p.nextToken()

	if p.curTokenIs(EOF) {
		return op, nil
	}

	if p.curTokenIs(SELECT) {
		p.nextToken()
		selects, err := p.parseSelects()
		if err != nil {
			return nil, err
		}
		op.selects = selects

		if p.curTokenIs(EOF) {
			return op, nil
		}

		p.nextToken()
	}

	if p.curTokenIs(WHERE) {
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
		}

		if p.curTokenIs(EOF) {
			return op, nil
		}
		p.nextToken()
	}

	if p.curTokenIs(ORDER) {
		p.nextToken()
		if !p.curTokenIs(BY) {
			return nil, fmt.Errorf("invalid: expected by but got %s", p.curToken.Type)
		}
		orderBys, err := p.parseOrderBy()
		if err != nil {
			return nil, err
		}
		op.orderBys = orderBys

		if p.curTokenIs(EOF) {
			return op, nil
		}
		p.nextToken()
	}

	if p.curTokenIs(LIMIT) {
		p.nextToken()
		if !p.curTokenIs(INT) {
			return nil, fmt.Errorf("invalid: expected int but got %s", p.curToken.Type)
		}
		limit, err := strconv.Atoi(p.curToken.Literal)
		if err != nil {
			return nil, fmt.Errorf("invalid: expected int but got %s", p.curToken.Literal)
		}
		op.limit = limit
	}

	return op, nil
}

func (p *Parser) parseGetOperation() (*GetOperation, error) {
	op := &GetOperation{}

	if !p.expectPeek(IDENT) {
		return nil, fmt.Errorf("invalid: expected path but got %s", p.curToken.Literal)
	}

	path := p.curToken.Literal
	lastSlash := strings.LastIndex(path, "/")
	if lastSlash == -1 {
		return nil, fmt.Errorf("invalid")
	}
	op.collection = path[:lastSlash]
	op.docId = path[lastSlash+1:]
	return op, nil
}

func (p *Parser) parseSelects() ([]string, error) {
	var selects []string
	for {
		if !p.curTokenIs(IDENT) {
			return nil, fmt.Errorf("invalid: expected field but got %s", p.curToken.Literal)
		}
		selects = append(selects, p.curToken.Literal)

		if !p.expectPeek(COMMA) {
			break
		}
		p.nextToken()
	}

	return selects, nil
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

func (p *Parser) parseOrderBy() ([]OrderBy, error) {
	orderBys := []OrderBy{}
	for !p.curTokenIs(EOF) {
		p.nextToken()
		if !p.curTokenIs(IDENT) {
			return nil, fmt.Errorf("invalid: expected field but got %s", p.curToken.Literal)
		}
		field := p.curToken.Literal
		p.nextToken()

		var fsDir firestore.Direction
		if p.curTokenIs(EOF) {
			fsDir = firestore.Asc
		} else if !p.curTokenIs(ASC) && !p.curTokenIs(DESC) {
			return nil, fmt.Errorf("invalid: expected direction but got %s", p.curToken.Literal)
		}
		if p.curTokenIs(ASC) {
			fsDir = firestore.Asc
		} else if p.curTokenIs(DESC) {
			fsDir = firestore.Desc
		}

		orderBys = append(orderBys, OrderBy{field, fsDir})

		if !p.expectPeek(COMMA) {
			return orderBys, nil
		}
	}
	return orderBys, nil
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

func (p *Parser) trimHeadSlash(s string) string {
	if strings.HasPrefix(s, "/") {
		return s[1:]
	}
	return s
}
