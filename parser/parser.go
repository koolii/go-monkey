package parser

import (
	"fmt"
	"strconv"

	"github.com/koolii/go-monkey/ast"
	"github.com/koolii/go-monkey/lexer"
	"github.com/koolii/go-monkey/token"
)

// 優先順位の順番の管理もしている
const (
	_ int = iota
	LOWEST
	EQUALS      // ==
	LESSGREATER // > or <
	SUM         // +
	PRODUCT     // *
	PREFIX      // -X or !X
	CALL        // myFunc(x)
)

type Parser struct {
	l      *lexer.Lexer
	errors []string

	// Lexerで言うところの position/readPositionのような動き
	// Lexerは次に読み込む無加工の1文字だったが、今回は文字ではなくtokenになる
	// curTokenだけで判断が出来ない時にpeekTokenを利用する
	curToken  token.Token
	peekToken token.Token

	// curToken.Typeに関連付けられた構文解析関数がマップにあるかどうかがすぐにチェックできる
	// 規約
	// - 構文解析関数に関連付けられたトークンが curToken にセットされている状態で動作を開始する
	// - この関数の処理対象である式の一番最後のトークンがcurTokenにセットされた状態になるまで進んで終了する
	// - トークンを進めすぎては行けない
	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l, errors: []string{}}
	fmt.Printf("Parser: %+v", p)

	// curToken/peekTokenを読み込む
	p.nextToken()
	p.nextToken()

	// mapの初期化
	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	p.infixParseFns = make(map[token.TokenType]infixParseFn)

	// 構文解析関数を登録
	// expressionをParseする際にここに登録してある関数を実行し、Expressionを取得
	p.registerPrefix(token.IDENT, p.parseIdentifier)
	p.registerPrefix(token.INT, p.parseIntegerLiteral)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)
	p.registerPrefix(token.BANG, p.parsePrefixExpression)

	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.SLASH, p.parseInfixExpression)
	p.registerInfix(token.ASTERISK, p.parseInfixExpression)
	p.registerInfix(token.EQ, p.parseInfixExpression)
	p.registerInfix(token.NOT_EQ, p.parseInfixExpression)
	p.registerInfix(token.LT, p.parseInfixExpression)
	p.registerInfix(token.GT, p.parseInfixExpression)

	return p
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead", t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
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
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
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

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.curToken}

	// returnトークンの次のexpressionのセクションまで移動させる
	p.nextToken()

	// 一旦セミコロンまでスキップ
	for !p.curTokenIs(token.SEMICOLON) {
		fmt.Printf("token is not semicolon: %+v\n", p)
		p.nextToken()
	}

	return stmt
}

func (p *Parser) expectPeek(tokenType token.TokenType) bool {
	if p.peekTokenIs(tokenType) {
		// トークンを一つ進める
		// peekTokenの型をチェックし、その型が正しい場合に限りnextToken()で次のトークンを読み出す
		p.nextToken()
		return true
	}
	p.peekError(tokenType)
	return false
}

func (p *Parser) curTokenIs(tokenType token.TokenType) bool {
	return p.curToken.Type == tokenType
}
func (p *Parser) peekTokenIs(tokenType token.TokenType) bool {
	return p.peekToken.Type == tokenType
}

// Expression section

// Prattでは、トークンタイプに最大2つの構文解析関数を割り当てられる
// 前置・中置
type (
	prefixParseFn func() ast.Expression
	// 構文解析中の中置演算子の「左側」
	infixParseFn func(ast.Expression) ast.Expression
)

// Parse内のマップにエントリを追加するヘルパーメソッド
func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}
func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}
func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.curToken}
	stmt.Expression = p.parseExpression(LOWEST)
	// もしもセミコロンがなかったとしても問題はない
	// 構文解析器にエラーを追加しない
	// なぜなら式分のセミコロンを省略できるようにするため
	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	return stmt
}

// p.curToken.Typeの前置に関連付けられた構文解析関数があるかを確認している
// TODO ここの中身分かってない
func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.curToken.Type]

	if prefix == nil {
		// 該当する prefxの演算子が存在しなかった場合(プログラムが解釈出来ないと判断される)
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}
	leftExp := prefix()
	fmt.Printf("parseExpression(leftExp): %+v\n", leftExp)

	// 2.6.9 中置演算子対応
	// TODO マジでなんで動くの？
	// 他に算術演算子(中置演算子)があるかどうか検索(precedenceで優先度付も実施)
	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}

		p.nextToken()

		leftExp = infix(leftExp)
	}

	return leftExp
}

func (p *Parser) noPrefixParseFnError(t token.TokenType) {
	msg := fmt.Sprintf("no prefix parse function for %s found", t)
	p.errors = append(p.errors, msg)
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

// 2.6.7
func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{Token: p.curToken}

	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}

	lit.Value = value

	return lit
}

// 2.6.8
func (p *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{Token: p.curToken, Operator: p.curToken.Literal}

	// getting Right expression
	p.nextToken()

	expression.Right = p.parseExpression(PREFIX)

	return expression
}

// 2.6.9

// 優先順位のマッピング
var precedences = map[token.TokenType]int{
	token.EQ:       EQUALS,
	token.NOT_EQ:   EQUALS,
	token.LT:       LESSGREATER,
	token.GT:       LESSGREATER,
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.SLASH:    PRODUCT,
	token.ASTERISK: PRODUCT,
}

// 次のトークン
func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}
	return LOWEST
}

// 現在のトークン
func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.curToken.Type]; ok {
		return p
	}
	return LOWEST
}

// parsePrefixExpression()との違いは、引数引数としてleftを取ること
// InfixExpressionを生成するのに必要
func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{Token: p.curToken, Operator: p.curToken.Literal, Left: left}

	precedence := p.curPrecedence()
	p.nextToken()
	expression.Right = p.parseExpression(precedence)

	return expression
}
