package fscli

type Lexer struct {
	input        []rune
	position     int
	readPosition int
	ch           rune
}

func NewLexer(input string) *Lexer {
	l := &Lexer{input: []rune(input)}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition += 1
}

func (l *Lexer) peekChar() rune {
	if l.readPosition >= len(l.input) {
		return 0
	} else {
		return l.input[l.readPosition]
	}
}

func (l *Lexer) NextToken() Token {
	var tok Token

	l.skipWhitespace()

	switch l.ch {
	case '=':
		tok = newToken(EQ, l.ch)
	case '!':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = Token{Type: NOT_EQ, Literal: literal}
		} else {
			// illegal
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = Token{Type: ILLEGAL, Literal: literal}
		}
	case '"':
		tok.Type = STRING
		tok.Literal = l.readString('"')
	case '\'':
		tok.Type = STRING
		tok.Literal = l.readString('\'')
	case '[':
		tok = newToken(LBRACKET, l.ch)
	case ']':
		tok = newToken(RBRACKET, l.ch)
	case ',':
		tok = newToken(COMMA, l.ch)
	case 0:
		tok.Literal = ""
		tok.Type = EOF
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = LookupIdent(tok.Literal)
			return tok
		} else if isDigit(l.ch) {
			tok.Type, tok.Literal = l.readNumber()
			return tok
		} else {
			tok = newToken(ILLEGAL, l.ch)
		}
	}

	l.readChar()
	return tok
}

func newToken(tokenType TokenType, ch rune) Token {
	return Token{Type: tokenType, Literal: string(ch)}
}

func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) || isDigit(l.ch) {
		l.readChar()
	}
	return string(l.input[position:l.position])
}

func (l *Lexer) readString(quote rune) string {
	l.readChar()
	position := l.position
	for l.ch != quote && l.ch != 0 {
		l.readChar()
	}

	r := string(l.input[position:l.position])
	return r
}

func isLetter(ch rune) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_' || ch == '/' || ch == '-' || ch == '.'
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

func (l *Lexer) readNumber() (TokenType, string) {
	position := l.position
	for isDigit(l.ch) {
		l.readChar()
	}

	if l.ch != '.' {
		return INT, string(l.input[position:l.position])
	}

	// float
	if l.ch == '.' {
		l.readChar()
	}
	for isDigit(l.ch) {
		l.readChar()
	}

	return FLOAT, string(l.input[position:l.position])
}

func isDigit(ch rune) bool {
	return '0' <= ch && ch <= '9'
}
