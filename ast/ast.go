package ast

import (
	"bytes"

	"github.com/koolii/go-monkey/token"
)

// 結局の所ASTを作成し、木構造となる

// let <identifier> = <expression>
// <identifier>/<expression>は可変で、<expression>にはいろいろな式が適用される
// e.g. 10, 関数リテラル(fn())

// すべてのノードはNodeインターフェイスを実装する
// TokenLiteral()はテスト・デバッグ用で利用する
type Node interface {
	TokenLiteral() string
	String() string
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

// バッファを作成、それぞれのXXXStatementのString()メソッドの戻り値をバッファに書き込む
func (p *Program) String() string {
	var out bytes.Buffer
	for _, s := range p.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}

// LetStatement let x = 10: などの文のNode
type LetStatement struct {
	Token token.Token // token.Let
	Name  *Identifier
	Value Expression // 値を生成する式を保持するため 値リテラル以外にも add(1, 10) * 100等がある
}

func (ls *LetStatement) statementNode()       {}
func (ls *LetStatement) TokenLiteral() string { return ls.Token.Literal }
func (ls *LetStatement) String() string {
	var out bytes.Buffer

	out.WriteString(ls.TokenLiteral() + " ")
	out.WriteString(ls.Name.String())
	out.WriteString(" = ")

	if ls.Value != nil {
		out.WriteString(ls.Value.String())
	}

	out.WriteString(";")

	return out.String()
}

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
func (i *Identifier) String() string       { return i.Value }

type ReturnStatement struct {
	Token       token.Token // token.RETURN
	ReturnValue Expression
}

func (rs *ReturnStatement) statementNode()       {}
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal } // return
func (rs *ReturnStatement) String() string {
	var out bytes.Buffer

	out.WriteString(rs.TokenLiteral() + " ")

	if rs.ReturnValue != nil {
		out.WriteString(rs.ReturnValue.String())
	}

	out.WriteString(";")

	return out.String()
}

// ExpressionStatement 式文でProgram.Statementsに追加できる
type ExpressionStatement struct {
	Token      token.Token // 式の最初のトークン
	Expression Expression  // 式を保存
}

func (es *ExpressionStatement) statementNode()       {}
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal }
func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}

// IntegerLiteral is for token.INT
type IntegerLiteral struct {
	Token token.Token // token.INT
	Value int64
}

func (il *IntegerLiteral) expressionNode()      {}
func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Literal }
func (il *IntegerLiteral) String() string       { return il.Token.Literal }
