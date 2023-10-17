package fscli

import "strings"

const (
	EOF     = "EOF"
	ILLEGAL = "ILLEGAL"

	GET    = "GET"
	QUERY  = "QUERY"
	COUNT  = "COUNT"
	SELECT = "SELECT"

	WHERE              = "WHERE"
	EQ                 = "="
	NOT_EQ             = "!="
	GT                 = ">"
	GTE                = ">="
	LT                 = "<"
	LTE                = "<="
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

	LIMIT = "LIMIT"

	LBRACKET = "["
	RBRACKET = "]"
	COMMA    = ","

	LIST_COLLECTIONS = "LIST_COLLECTIONS"
	PAGER            = "PAGER"
)

type TokenType = string

type Token struct {
	Type    TokenType
	Literal string
}

var keywards = map[string]TokenType{
	"GET":    GET,
	"QUERY":  QUERY,
	"COUNT":  COUNT,
	"SELECT": SELECT,
	"WHERE":  WHERE,
	"AND":    AND,
	"ORDER":  ORDER,
	"BY":     BY,
	"ASC":    ASC,
	"DESC":   DESC,
	"LIMIT":  LIMIT,
}

var operators = map[string]TokenType{
	"==":                 EQ,
	"!=":                 NOT_EQ,
	">":                  GT,
	">=":                 GTE,
	"<":                  LT,
	"<=":                 LTE,
	"IN":                 IN,
	"ARRAY_CONTAINS":     ARRAY_CONTAINS,
	"ARRAY_CONTAINS_ANY": ARRAY_CONTAINS_ANY,
}

var metacommands = map[string]TokenType{
	`\d`:     LIST_COLLECTIONS,
	`\pager`: PAGER,
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

func LookupMetacommand(s string) TokenType {
	if tok, ok := metacommands[s]; ok {
		return tok
	}
	return ILLEGAL
}
