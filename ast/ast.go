package ast

import "github.com/koolii/go-monkey/token"

// 結局の所ASTを作成し、木構造となる

// let <identifier> = <expression>
// <identifier>/<expression>は可変で、<expression>にはいろいろな式が適用される
// e.g. 10, 関数リテラル(fn())

// すべてのノードはNodeインターフェイスを実装する
// TokenLiteral()はテスト・デバッグ用で利用する
type Node interface {
	TokenLiteral() string
}

// statementNode()はダミーでExpressionインターフェイス
// との間違いをコンパイラが指摘しやすいようにする
type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

// Program すべてのASTのルートノードとなる
// すべての有効なMonkeyプログラムはStatementsにキャッシュされる
type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		// ここはなぜ0インデックスになっている？
		return p.Statements[0].TokenLiteral()
	}
	return ""
}

// LetStatement let x = 10: などの文のNode
type LetStatement struct {
	Token token.Token // token.Let
	Name  *Identifier
	Value Expression // 値を生成する式を保持するため 値リテラル以外にも add(1, 10) * 100等がある
}

func (ls *LetStatement) statementNode()       {}
func (ls *LetStatement) TokenLiteral() string { return ls.Token.Literal }

// Identifier 変数名の確保のためのNode(識別子)
// token.IDENTのTypeはIDENT、Literalは実際の変数名文字列
type Identifier struct {
	Token token.Token // token.IDENT
	Value string
}

// IdentifierはExpressionを実装しているが、これはNodeの種類を少なくするために
// 簡単にするため
// そもそもIdentifierはNodeの種類を少なくするため、変数束縛の名前を表現のために作っている
func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }

type ReturnStatement struct {
	Token       token.Token // token.RETURN
	ReturnValue Expression
}

func (rs *ReturnStatement) statementNode()       {}
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal } // return
