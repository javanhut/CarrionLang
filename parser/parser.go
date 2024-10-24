// parser/parser.go
package parser

import (
    "carrionlang/ast"
    "carrionlang/lexer"
    "carrionlang/token"
    "strconv"
)

type Parser struct {
    l      *lexer.Lexer
    errors []string

    curToken  token.Token
    peekToken token.Token
}

func New(l *lexer.Lexer) *Parser {
    p := &Parser{
        l:      l,
        errors: []string{},
    }
    // Read two tokens to initialize curToken and peekToken
    p.nextToken()
    p.nextToken()
    return p
}

func (p *Parser) Errors() []string {
    return p.errors
}

func (p *Parser) nextToken() {
    p.curToken = p.peekToken
    p.peekToken = p.l.NextToken()
}

func (p *Parser) ParseProgram() *ast.Program {
    program := &ast.Program{}
    program.Statements = []ast.Statement{}

    for p.curToken.Type != token.EOF {
        stmt := p.parseStatement()
        if stmt != nil {
            program.Statements = append(program.Statements, stmt)
        }
        p.nextToken()
    }
    return program
}

func (p *Parser) parseStatement() ast.Statement {
    switch p.curToken.Type {
    case token.IDENT:
        return p.parseVariableDeclaration()
    case token.SPELLBOOK:
        return p.parseSpellbookDeclaration()
    case token.NEWLINE:
        return nil
    default:
        return p.parseExpressionStatement()
    }
}

func (p *Parser) parseVariableDeclaration() *ast.VariableDeclaration {
    stmt := &ast.VariableDeclaration{Token: p.curToken}
    stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

    p.nextToken()

    if p.curToken.Type == token.COLON {
        // Type hint is provided
        p.nextToken()
        stmt.TypeHint = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
        p.nextToken()
    }

    if p.curToken.Type != token.ASSIGN {
        p.errors = append(p.errors, "expected '=' after variable declaration")
        return nil
    }

    p.nextToken()

    stmt.Value = p.parseExpression(LOWEST)

    return stmt
}

func (p *Parser) parseSpellbookDeclaration() *ast.SpellbookDeclaration {
    stmt := &ast.SpellbookDeclaration{Token: p.curToken}

    p.nextToken()

    if p.curToken.Type != token.IDENT {
        p.errors = append(p.errors, "expected identifier after 'spellbook'")
        return nil
    }

    stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

    p.nextToken()

    if p.curToken.Type != token.COLON {
        p.errors = append(p.errors, "expected ':' after spellbook name")
        return nil
    }

    p.nextToken()

    if p.curToken.Type != token.NEWLINE {
        p.errors = append(p.errors, "expected newline after ':'")
        return nil
    }

    p.nextToken()

    if p.curToken.Type != token.INDENT {
        p.errors = append(p.errors, "expected indentation")
        return nil
    }

    p.nextToken()

    stmt.Body = []ast.Statement{}
    for p.curToken.Type != token.DEDENT && p.curToken.Type != token.EOF {
        bodyStmt := p.parseStatement()
        if bodyStmt != nil {
            stmt.Body = append(stmt.Body, bodyStmt)
        }
        p.nextToken()
    }

    return stmt
}

// Add parsing functions for spells, expressions, etc.

// Precedence levels
const (
    _ int = iota
    LOWEST
    EQUALS      // ==
    LESSGREATER // > or <
    SUM         // +
    PRODUCT     // *
    PREFIX      // -X or !X
    CALL        // myFunction(X)
)

// Add more parsing functions as needed...

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
    stmt := &ast.ExpressionStatement{Token: p.curToken}
    stmt.Expression = p.parseExpression(LOWEST)
    return stmt
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
    var leftExp ast.Expression

    switch p.curToken.Type {
    case token.IDENT:
        leftExp = p.parseIdentifier()
    case token.INT:
        leftExp = p.parseIntegerLiteral()
    case token.STRING:
        leftExp = p.parseStringLiteral()
    default:
        return nil
    }

    for !p.peekTokenIs(token.NEWLINE) && precedence < p.peekPrecedence() {
        p.nextToken()
        leftExp = p.parseInfixExpression(leftExp)
    }

    return leftExp
}

func (p *Parser) parseIdentifier() ast.Expression {
    return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
    lit := &ast.IntegerLiteral{Token: p.curToken}

    value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
    if err != nil {
        msg := "could not parse %q as integer"
        p.errors = append(p.errors, msg)
        return nil
    }
    lit.Value = value
    return lit
}

func (p *Parser) parseStringLiteral() ast.Expression {
    return &ast.StringLiteral{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
    expr := &ast.InfixExpression{
        Token:    p.curToken,
        Operator: p.curToken.Literal,
        Left:     left,
    }
    precedence := p.curPrecedence()
    p.nextToken()
    expr.Right = p.parseExpression(precedence)
    return expr
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
    return p.peekToken.Type == t
}

func (p *Parser) curTokenIs(t token.TokenType) bool {
    return p.curToken.Type == t
}

func (p *Parser) expectPeek(t token.TokenType) bool {
    if p.peekTokenIs(t) {
        p.nextToken()
        return true
    } else {
        p.peekError(t)
        return false
    }
}

func (p *Parser) peekError(t token.TokenType) {
    msg := "expected next token to be %s, got %s instead"
    p.errors = append(p.errors, msg)
}

var precedences = map[token.TokenType]int{
    token.EQ:       EQUALS,
    token.NOT_EQ:   EQUALS,
    token.LT:       LESSGREATER,
    token.GT:       LESSGREATER,
    token.PLUS:     SUM,
    token.MINUS:    SUM,
    token.SLASH:    PRODUCT,
    token.ASTERISK: PRODUCT,
    token.LPAREN:   CALL,
}

func (p *Parser) peekPrecedence() int {
    if p, ok := precedences[p.peekToken.Type]; ok {
        return p
    }
    return LOWEST
}

func (p *Parser) curPrecedence() int {
    if p, ok := precedences[p.curToken.Type]; ok {
        return p
    }
    return LOWEST
}

