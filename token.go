package fscli

import "strings"

const (
	EOF     = "EOF"
	ILLEGAL = "ILLEGAL"

	GET    = "GET"
	QUERY  = "QUERY"
	SELECT = "SELECT"

	WHERE              = "WHERE"
	EQ                 = "=="
	NOT_EQ             = "!="
	IN                 = "IN"
	ARRAY_CONTAINS     = "ARRAY_CONTAINS"
	ARRAY_CONTAINS_ANY = "ARRAY_CONTAINS_ANY"
	ORDER              = "ORDER"
	BY                 = "BY"

	ASC  = "ASC"
	DESC = "DESC"

	IDENT  = "IDENT"
	STRING = "STRING"
	INT    = "INT"
	FLOAT  = "FLOAT"

	AND = "AND"

	LBRACKET = "["
	RBRACKET = "]"
	COMMA    = ","
)

type TokenType = string

type Token struct {
	Type    TokenType
	Literal string
}

var keywards = map[string]TokenType{
	"GET":    GET,
	"QUERY":  QUERY,
	"SELECT": SELECT,
	"WHERE":  WHERE,
	"AND":    AND,
	"ORDER":  ORDER,
	"BY":     BY,
	"ASC":    ASC,
	"DESC":   DESC,
}

var operators = map[string]TokenType{
	"==":                 EQ,
	"!=":                 NOT_EQ,
	"IN":                 IN,
	"ARRAY_CONTAINS":     ARRAY_CONTAINS,
	"ARRAY_CONTAINS_ANY": ARRAY_CONTAINS_ANY,
}

func LookupIdent(ident string) TokenType {
	u := strings.ToUpper(ident)
	if tok, ok := keywards[u]; ok {
		return tok
	}
	if tok, ok := operators[u]; ok {
		return tok
	}
	return IDENT
}
