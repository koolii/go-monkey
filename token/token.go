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
	PLUS = "+"

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
)

var keywords = map[string]TokenType{
	"fn":  FUNCTION,
	"let": LET,
}

func LookupIdent(ident string) TokenType {
	// 予約語かどうか判定
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	// それ以外は変数として定義
	return IDENT
}
