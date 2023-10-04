package fscli

import "strings"

const (
	EOF     = "EOF"
	ILLEGAL = "ILLEGAL"

	QUERY  = "QUERY"
	WHERE  = "WHERE"
	EQ     = "=="
	NOT_EQ = "!="

	IDENT  = "IDENT"
	STRING = "STRING"
	INT    = "INT"
	FLOAT  = "FLOAT"
)

type TokenType = string

type Token struct {
	Type    TokenType
	Literal string
}

var keywards = map[string]TokenType{
	"QUERY": QUERY,
	"WHERE": WHERE,
}

var operators = map[string]TokenType{
	"==": EQ,
	"!=": NOT_EQ,
}

func LookupIdent(ident string) TokenType {
	u := strings.ToUpper(ident)
	if tok, ok := keywards[u]; ok {
		return tok
	}
	if tok, ok := operators[ident]; ok {
		return tok
	}
	return IDENT
}
