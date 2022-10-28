package parser

import (
	"fmt"
	"strconv"
	"gopherlang-interpreter/ast"
	"gopherlang-interpreter/lexer"
	"gopherlang-interpreter/token"
)

// iota const gives options incrementing values (precedence)
const (
	_ int = iota
	LOWEST
	EQUALS		// ==
	LESSGREATER // > or <
	SUM 		// +
	PRODUCT		// *
	PREFIX		// -X or !X
	CALL		// myFunction(X)
)

type Parser struct {
	l *lexer.Lexer

	curToken  token.Token
	peekToken token.Token

	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
	errors []string
}

type (
	prefixParseFn func() ast.Expression
	infixParseFn func(ast.Expression) ast.Expression
)

// Init new Parser
func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}

	// Read two tokens so that both curToken and peekToken are set
	p.nextToken()
	p.nextToken()

	// Map TokenType's to prefix parse fn's
	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	p.registerPrefix(token.IDENT, p.parseIdentifier)
	p.registerPrefix(token.INT, p.parseIntegerLiteral)

	return p
}

// Get parser errors
func (p *Parser) Errors() []string {
	return p.errors
}

// Advance both curToken and peekToken from the Parser's lexer
func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) ParseProgram() *ast.Program {
	// Construct root node
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	// Parse all statements until reaching EOF token. (Recursive Descent)
	for !p.curTokenIs(token.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}

	return program
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}

// Parse LET statement
func (p *Parser) parseLetStatement() *ast.LetStatement {
	// Construct Statement Node
	stmt := &ast.LetStatement{Token: p.curToken}

	// First expect IDENT token for let statement
	if !p.expectPeek(token.IDENT) {
		return nil
	}

	// Construct Identifier Node
	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	// Next expect ASSIGN token for let statement
	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	// TODO:
	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

// Parse RETURN statement
func (p* Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.curToken}

	p.nextToken()

	//TODO:
	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

// Parse EXPRESSION STATEMENT
func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.curToken}
	
	stmt.Expression = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.curToken.Type]

	if prefix == nil {
		return nil
	}

	leftExp := prefix()
	
	return leftExp
}

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

// Helper to check parser.peekToken is of TokenType t
func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}

// Helper to check parser.peekToken is of TokenType t
func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead",
		t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

// Check that next token (parser.peekToken) is type t and advance using parser.nextToken()
func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		p.peekError(t)
		return false
	}
}

// Helpers for registering tokens to prefix/infix parse fn's
func (p* Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}
func (p* Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}
