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
	l := &Lexer{input: input}
	l.readChar()
	return l
}

func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool {
	// 浮動点小数型は完全にサポートしていなかった
	return '0' <= ch && ch <= '9'
}

// readChar 次の1文字を呼んでinput文字列の現在位置(postiion)を進める
// Unicodeに対応していない
// 対応する場合は l.chをbyteからruneに変更し、次の文字を読む処理を変更する必要がある
// (※ 次の文字が複数のバイトから構成される可能性があるため l.input[l.readPosition]は使えない
// ミュータブルにLexer内を移動させる
func (l *Lexer) readChar() {
	// 終端チェック
	if l.readPosition >= len(l.input) {
		// ASCIIで言うところの "NUL"
		l.ch = 0
		fmt.Printf("readChar(): EOF\n")
	} else {
		l.ch = l.input[l.readPosition]
		fmt.Printf("readChar(): %c\n", l.ch)
	}
	// positionの更新処理
	l.position = l.readPosition
	l.readPosition++
}

// NextToken is increments Lexer
func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhitespace()

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
	default:
		fmt.Println("This is default case")
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			fmt.Printf("This is identifier: %s\n", tok.Literal)
			// readIdentifier()でreadChar()を実行しているため、余分にreadChar()を実行させない
			return tok
		} else if isDigit(l.ch) {
			tok.Type = token.INT
			tok.Literal = l.readNumber()
			return tok
		}
		tok = newToken(token.ILLEGAL, l.ch)
	}

	l.readChar()
	return tok
}

func (l *Lexer) skipWhitespace() {
	// 空白・改行がなくなるまで続ける
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

func (l *Lexer) readIdentifier() string {
	// 読み込む初期位置を取得
	position := l.position
	// 途切れるところまで読み込む(終端位置を取得)
	for isLetter(l.ch) {
		l.readChar()
	}
	// 初期位置-終端位置までの文字列を取得
	return l.input[position:l.position]
}

func (l *Lexer) readNumber() string {
	// 読み込む初期位置を取得
	position := l.position
	// 途切れるところまで読み込む(終端位置を取得)
	for isDigit(l.ch) {
		l.readChar()
	}
	// 初期位置-終端位置までの文字列を取得
	return l.input[position:l.position]
}
