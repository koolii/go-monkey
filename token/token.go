package token

// TokenType is alias of string
type TokenType string

// Token is keyword type of Monkey
type Token struct {
	Type    TokenType
	Literal string
}

const (
	// ILLEGAL トークンや文字が未知であることを表す
	ILLEGAL = "ILLEGAL"
	// EOF どこで読み込みを停止するか構文解析器に伝える
	EOF = "EOF"

	// IDENT literal
	IDENT = "IDENT" // add, foobar, x, y...
	// INT number
	INT = "INT" // 23142432

	// ASSIGN define various
	ASSIGN = "="
	// PLUS add
	PLUS     = "+"
	MINUS    = "-"
	BANG     = "!"
	ASTERISK = "*"
	SLASH    = "/"

	LT = "<"
	GT = ">"

	// COMMA delimiter
	COMMA = ","
	// SEMICOLON end of line
	SEMICOLON = ";"

	LPAREN = "("
	RPAREN = ")"
	LBRACE = "{"
	RBRACE = "}"

	// keyword
	FUNCTION = "FUNCTION"
	LET      = "LET"
	TRUE     = "TRUE"
	FALSE    = "FALSE"
	IF       = "IF"
	ELSE     = "ELSE"
	RETURN   = "RETURN"

	EQ     = "=="
	NOT_EQ = "!="
)

var keywords = map[string]TokenType{
	"fn":     FUNCTION,
	"let":    LET,
	"true":   TRUE,
	"false":  FALSE,
	"if":     IF,
	"else":   ELSE,
	"return": RETURN,
}

func LookupIdent(ident string) TokenType {
	// 予約語かどうか判定
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	// それ以外は変数として定義
	return IDENT
}
