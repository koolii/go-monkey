package parser

import (
	"fmt"
	"testing"

	"github.com/koolii/go-monkey/ast"
	"github.com/koolii/go-monkey/lexer"
)

func TestLetStatement(t *testing.T) {
	input := `
let x = 5;
let y = 10;
let foobar = 838383;
`

	l := lexer.New(input)
	p := New(l)

	fmt.Println("parse program")
	program := p.ParseProgram()
	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}
	// ここはどういう動き？
	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements. got=%d\n", len(program.Statements))
	}

	tests := []struct {
		expectedIdentifier string
	}{
		{"x"},
		{"y"},
		{"foobar"},
	}

	for i, tt := range tests {
		stmt := program.Statements[i]
		fmt.Printf("%+v\n", stmt)
		if !testLetStatement(t, stmt, tt.expectedIdentifier) {
			return
		}
	}
}

func testLetStatement(t *testing.T, s ast.Statement, name string) bool {
	if s.TokenLiteral() != "let" {
		t.Errorf("s.TokenLiteral not 'let'. got=%q\n", s.TokenLiteral())
		return false
	}

	// cast
	fmt.Println("cast")
	letStmt, ok := s.(*ast.LetStatement)
	if !ok {
		t.Errorf("s not *ast.LetStatement. got=%T\n", s)
		return false
	}

	// 変数名チェック(LetStatement.Name.Valueは変数名が入る事をチェック)
	// LetStatement.Name.Tokenは問答無用でtoken.IDENT
	if letStmt.Name.Value != name {
		t.Errorf("letStmt.Name.Value not '%s'. got=%s\n", name, letStmt.Name.Value)
		return false
	}

	// これは何をやっている？TokenLiteral()だから入力値になる気がしている。TokenはIDENTになると思う
	fmt.Printf("letStmt.Name.TokenLiteral() is '%s'\n", letStmt.Name.TokenLiteral())
	if letStmt.Name.TokenLiteral() != name {
		t.Errorf("letStmt.Name.TokenLiteral() not '%s'. got=%s\n", name, letStmt.Name.TokenLiteral())
		return false
	}

	return true
}
