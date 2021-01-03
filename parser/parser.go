package parser

import (
	"fmt"

	"github.com/koolii/go-monkey/ast"
	"github.com/koolii/go-monkey/lexer"
	"github.com/koolii/go-monkey/token"
)

type Parser struct {
	l *lexer.Lexer

	// Lexerで言うところの position/readPositionのような動き
	// Lexerは次に読み込む無加工の1文字だったが、今回は文字ではなくtokenになる
	// curTokenだけで判断が出来ない時にpeekTokenを利用する
	curToken  token.Token
	peekToken token.Token
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l}
	fmt.Printf("Parser: %+v", p)

	// curToken/peekTokenを読み込む
	p.nextToken()
	p.nextToken()

	return p
}

// 次のtokenに移動する
func (p *Parser) nextToken() {
	// ?構造体を生成したタイミングで peekToken等も初期化される？
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) ParseProgram() *ast.Program {
	// ルートノードを作成
	program := &ast.Program{}
	// 空のスライスで初期化
	program.Statements = []ast.Statement{}

	for p.curToken.Type != token.EOF {
		// 一文をパース
		stmt := p.parseStatement()
		if stmt != nil {
			// 追加
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}
	return program
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.LET:
		// ここのReturnTypeが ast.Statementになっているが、
		// これを *ast.Statementにするとエラーとなる
		// よく分かっていないが、 Statement < LetStatementの構成だが、だが、ポインタを利用すると継承？がうまく出来ない？
		return p.parseLetStatement()
	default:
		return nil
	}
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: p.curToken}

	if !p.expectPeek(token.IDENT) {
		return nil
	}

	// この時点でtokenが一つ進んでいる
	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	// 一旦セミコロンまでスキップ
	for !p.curTokenIs(token.SEMICOLON) {
		fmt.Printf("token is not semicolon: %+v\n", p)
		p.nextToken()
	}

	return stmt
}

func (p *Parser) expectPeek(tokenType token.TokenType) bool {
	if p.peekTokenIs(tokenType) {
		// if p.peekToken.Type == tokenType {
		// トークンを一つ進める
		p.nextToken()
		return true
	}
	return false
}

func (p *Parser) curTokenIs(tokenType token.TokenType) bool {
	return p.curToken.Type == tokenType
}
func (p *Parser) peekTokenIs(tokenType token.TokenType) bool {
	return p.peekToken.Type == tokenType
}
