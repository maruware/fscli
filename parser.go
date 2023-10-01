package fscli

import (
	"strconv"
)

func Parse(input string) ([]Token, error) {
	l := NewLexer(input)
	tokens := []Token{}
	for {
		tok := l.NextToken()
		if tok.Type == EOF {
			break
		}
		tokens = append(tokens, tok)
	}
	return tokens, nil
}

func parseInt(word string) (int, error) {
	return strconv.Atoi(word)
}

func parseFloat(word string) (float64, error) {
	return strconv.ParseFloat(word, 64)
}

func parseString(word string) (string, error) {
	return strconv.Unquote(word)
}
