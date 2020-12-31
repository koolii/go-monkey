package lexer

import (
	"fmt"

	"github.com/koolii/go-monkey/token"
)

// Lexer 字句解析器の１単位
// position/readPositionは入力を「覗き見」して、
// 現在の文字に続いて何が来るかを考慮するため
// readPositionは常に入力における「次の」１を指し示す
// positionは現在検査中のバイトchの位置を示す
type Lexer struct {
	input        string
	position     int  // 入力における現在の位置(現在の文字を指し示す)
	readPosition int  // これから読み込む位置(現在の文字の次)
	ch           byte // 現在検査中の文字
}

// New is create Lexer pointer
func New(input string) *Lexer {
	fmt.Println(input)
	l := &Lexer{input: input}
	l.readChar()
	return l
}

// readChar 次の1文字を呼んでinput文字列の現在位置(postiion)を進める
// Unicodeに対応していない
// 対応する場合は l.chをbyteからruneに変更し、次の文字を読む処理を変更する必要がある
// (※ 次の文字が複数のバイトから構成される可能性があるため l.input[l.readPosition]は使えない
func (l *Lexer) readChar() {
	// 終端チェック
	if l.readPosition >= len(l.input) {
		// ASCIIで言うところの "NUL"
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	// positionの更新処理
	l.position = l.readPosition
	l.readPosition++
}

// NextToken is increments Lexer
func (l *Lexer) NextToken() (tok token.Token) {
	// case: 0って数字の時はどうする?
	switch l.ch {
	case '=':
		tok = newToken(token.ASSIGN, l.ch)
	case ';':
		tok = newToken(token.SEMICOLON, l.ch)
	case '(':
		tok = newToken(token.LPAREN, l.ch)
	case ')':
		tok = newToken(token.RPAREN, l.ch)
	case ',':
		tok = newToken(token.COMMA, l.ch)
	case '+':
		tok = newToken(token.PLUS, l.ch)
	case '{':
		tok = newToken(token.LBRACE, l.ch)
	case '}':
		tok = newToken(token.RBRACE, l.ch)
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	}

	l.readChar()
	return tok
}

func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}
