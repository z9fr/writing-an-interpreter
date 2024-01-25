package lexter

import (
	"monkey-lang.z9fr.xyz/internal/token"
)

/*
lexer

input is the source code and outputs the tokens. we dont need to save code or anything
we only need a itter and methord called NextToken to get the next token

we do need to initiate lexer and call our itter
*/

type Lexer struct {
	input        string
	position     int  // current position in input (points to current char)
	readPosition int  // current reading position in input (after current char)
	ch           byte // current char under examination
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

// readChar is used to get the next character and move in to next char of the input
func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition += 1
}

func (l *Lexer) NextToken() token.Token {
	var t token.Token

	l.skipWhitespace()

	switch l.ch {
	// we are going to extend `=` and `!` to support operations like
	// `==`, `!=`
	case '=':
		if l.PeekChar() == '=' {
			// we basically assign the current char in to variable
			// and we read the next possition to build the token type
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			t = token.Token{Type: token.EQ, Literal: literal}
		} else {
			t = NewToken(token.ASSIGN, l.ch)
		}
	case '!':
		if l.PeekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			t = token.Token{Type: token.NOT_EQ, Literal: literal}
		} else {
			t = NewToken(token.BANG, l.ch)
		}
	case ';':
		t = NewToken(token.SEMICOLON, l.ch)
	case '(':
		t = NewToken(token.LPAREN, l.ch)
	case ')':
		t = NewToken(token.RPAREN, l.ch)
	case ',':
		t = NewToken(token.COMMA, l.ch)
	case '+':
		t = NewToken(token.PLUS, l.ch)
	case '-':
		t = NewToken(token.MINUS, l.ch)
	case '*':
		t = NewToken(token.ASTERISK, l.ch)
	case '/':
		t = NewToken(token.SLASH, l.ch)
	case '<':
		t = NewToken(token.LT, l.ch)
	case '>':
		t = NewToken(token.GT, l.ch)
	case '{':
		t = NewToken(token.LBRACE, l.ch)
	case '}':
		t = NewToken(token.RBRACE, l.ch)
	case 0:
		t.Type = token.EOF
		t.Literal = ""
	default:
		if isLetter(l.ch) {
			t.Literal = l.readIdentifier()
			t.Type = token.LookupIdent(t.Literal)
			return t
		} else if isDigit(l.ch) {
			t.Literal = l.readNumber()
			t.Type = token.INT
			return t
		} else {
			t = NewToken(token.ILLEGAL, l.ch)
		}
	}

	l.readChar()
	return t
}

func NewToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}

func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

// the whitespace charactors in between values are there for example
// let value = 5;
// so we do need to skip these charts
// in this lang whitespaces are just seperator
func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func (l *Lexer) readNumber() string {
	// TODO: need to get float numbers etc
	position := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

// peek char is same as readChar but it doesnt increment the current possition
func (l *Lexer) PeekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	} else {
		return l.input[l.readPosition]
	}
}
