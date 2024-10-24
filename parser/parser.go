// parser/parser.go
package parser

import (
    "carrionlang/ast"
    "carrionlang/lexer"
    "carrionlang/token"
    "strconv"
)

type (
    prefixParseFn func() ast.Expression
    infixParseFn  func(ast.Expression) ast.Expression
)

const (
    _ int = iota
    LOWEST
    EQUALS      // ==
    LESSGREATER // > or <
    SUM         // +
    PRODUCT     // *
    PREFIX      // -X or !X
    CALL        // myFunction(X)
    MEMBER      // object.property
)

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
    token.DOT:      MEMBER,
}

type Parser struct {
    l      *lexer.Lexer
    errors []string

    curToken  token.Token
    peekToken token.Token

    prefixParseFns map[token.TokenType]prefixParseFn
    infixParseFns  map[token.TokenType]infixParseFn
}

func New(l *lexer.Lexer) *Parser {
    p := &Parser{
        l:      l,
        errors: []string{},
    }

    p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
    p.registerPrefix(token.IDENT, p.parseIdentifier)
    p.registerPrefix(token.INT, p.parseIntegerLiteral)
    p.registerPrefix(token.STRING, p.parseStringLiteral)
    p.registerPrefix(token.MINUS, p.parsePrefixExpression)
    p.registerPrefix(token.BANG, p.parsePrefixExpression)
    p.registerPrefix(token.RPAREN, p.parsePrefixExpression)
    p.registerPrefix(token.NEWLINE, p.parseNewline)
    p.infixParseFns = make(map[token.TokenType]infixParseFn)
    p.registerInfix(token.PLUS, p.parseInfixExpression)
    p.registerInfix(token.MINUS, p.parseInfixExpression)
    p.registerInfix(token.SLASH, p.parseInfixExpression)
    p.registerInfix(token.ASTERISK, p.parseInfixExpression)
    p.registerInfix(token.EQ, p.parseInfixExpression)
    p.registerInfix(token.NOT_EQ, p.parseInfixExpression)
    p.registerInfix(token.LT, p.parseInfixExpression)
    p.registerInfix(token.GT, p.parseInfixExpression)
    p.registerInfix(token.LPAREN, p.parseCallExpression)
    p.registerInfix(token.DOT, p.parseMemberExpression)

    // Read two tokens to initialize curToken and peekToken
    p.nextToken()
    p.nextToken()

    return p
}

func (p *Parser) Errors() []string {
    return p.errors
}

func (p *Parser) parseNewline() ast.Expression {
    p.errors = append(p.errors, "unexpected token 'NEWLINE'")
    return nil
}

func (p *Parser) peekError(t token.TokenType) {
    msg := "expected next token to be " + string(t) + ", got " + string(p.peekToken.Type) + " instead"
    p.errors = append(p.errors, msg)
}

func (p *Parser) noPrefixParseFnError(t token.TokenType) {
    msg := "no prefix parse function for " + string(t) + " found"
    p.errors = append(p.errors, msg)
}

func (p *Parser) parseRPAREN() ast.Expression {
    p.errors = append(p.errors, "unexpected token ')'")
    return nil
}

func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
    p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
    p.infixParseFns[tokenType] = fn
}

func (p *Parser) nextToken() {
    p.curToken = p.peekToken
    p.peekToken = p.l.NextToken()
}

func (p *Parser) expectPeek(t token.TokenType) bool {
    if p.peekToken.Type == t {
        p.nextToken()
        return true
    } else {
        p.peekError(t)
        return false
    }
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
        if p.peekToken.Type == token.ASSIGN || p.peekToken.Type == token.COLON {
            return p.parseVariableDeclaration()
        }
        return p.parseExpressionStatement()
    case token.SPELLBOOK:
        return p.parseSpellbookDeclaration()
    case token.SPELL:
        return p.parseSpellDeclaration()
    case token.RETURN:
        return p.parseReturnStatement()
    case token.NEWLINE:
        return nil // Skip NEWLINE tokens
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
        if p.curToken.Type != token.IDENT {
            p.errors = append(p.errors, "expected type identifier after ':'")
            return nil
        }
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
        p.errors = append(p.errors, "expected indentation after spellbook declaration")
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

func (p *Parser) parseSpellDeclaration() *ast.SpellDeclaration {
    sd := &ast.SpellDeclaration{Token: p.curToken}

    p.nextToken() // Move to spell name

    if p.curToken.Type != token.IDENT {
        p.errors = append(p.errors, "expected spell name after 'spell'")
        return nil
    }

    sd.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

    p.nextToken()

    if p.curToken.Type != token.LPAREN {
        p.errors = append(p.errors, "expected '(' after spell name")
        return nil
    }

    sd.Parameters = p.parseFunctionParameters()

    if p.peekToken.Type == token.ARROW {
        p.nextToken() // Move to '->'
        p.nextToken() // Move to return type
        if p.curToken.Type != token.IDENT {
            p.errors = append(p.errors, "expected return type identifier after '->'")
            return nil
        }
        sd.ReturnType = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
    }

    if !p.expectPeek(token.COLON) {
        return nil
    }

    if !p.expectPeek(token.NEWLINE) {
        return nil
    }

    if !p.expectPeek(token.INDENT) {
        return nil
    }

    p.nextToken()

    sd.Body = p.parseBlockStatement()

    return sd
}

func (p *Parser) parseFunctionParameters() []*ast.Identifier {
    identifiers := []*ast.Identifier{}

    if p.peekToken.Type == token.RPAREN {
        p.nextToken()
        return identifiers
    }

    p.nextToken() // Move to first parameter

    ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
    identifiers = append(identifiers, ident)

    for p.peekToken.Type == token.COMMA {
        p.nextToken() // Skip ','
        p.nextToken() // Move to next parameter
        ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
        identifiers = append(identifiers, ident)
    }

    if !p.expectPeek(token.RPAREN) {
        return nil
    }

    return identifiers
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
    stmt := &ast.ReturnStatement{Token: p.curToken}

    p.nextToken()

    stmt.ReturnValue = p.parseExpression(LOWEST)

    return stmt
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
    block := &ast.BlockStatement{Token: p.curToken}
    block.Statements = []ast.Statement{}

    for p.curToken.Type != token.DEDENT && p.curToken.Type != token.EOF {
        stmt := p.parseStatement()
        if stmt != nil {
            block.Statements = append(block.Statements, stmt)
        }
        p.nextToken()
    }

    return block
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
    stmt := &ast.ExpressionStatement{Token: p.curToken}
    stmt.Expression = p.parseExpression(LOWEST)
    return stmt
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
    prefix := p.prefixParseFns[p.curToken.Type]
    if prefix == nil {
        p.noPrefixParseFnError(p.curToken.Type)
        return nil
    }
    leftExp := prefix()

    for p.peekToken.Type != token.NEWLINE && precedence < p.peekPrecedence() {
        infix := p.infixParseFns[p.peekToken.Type]
        if infix == nil {
            return leftExp
        }

        p.nextToken()

        leftExp = infix(leftExp)
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
        msg := "could not parse " + p.curToken.Literal + " as integer"
        p.errors = append(p.errors, msg)
        return nil
    }
    lit.Value = value
    return lit
}

func (p *Parser) parseStringLiteral() ast.Expression {
    return &ast.StringLiteral{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parsePrefixExpression() ast.Expression {
    expression := &ast.PrefixExpression{
        Token:    p.curToken,
        Operator: p.curToken.Literal,
    }

    p.nextToken()

    expression.Right = p.parseExpression(PREFIX)

    return expression
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
    expression := &ast.InfixExpression{
        Token:    p.curToken,
        Operator: p.curToken.Literal,
        Left:     left,
    }

    precedence := p.curPrecedence()
    p.nextToken()
    expression.Right = p.parseExpression(precedence)

    return expression
}

func (p *Parser) parseCallExpression(function ast.Expression) ast.Expression {
    exp := &ast.CallExpression{Token: p.curToken, Function: function}
    exp.Arguments = p.parseCallArguments()
    return exp
}

func (p *Parser) parseCallArguments() []ast.Expression {
    args := []ast.Expression{}

    if p.peekToken.Type == token.RPAREN {
        p.nextToken()
        return args
    }

    p.nextToken()
    args = append(args, p.parseExpression(LOWEST))

    for p.peekToken.Type == token.COMMA {
        p.nextToken()
        p.nextToken()
        args = append(args, p.parseExpression(LOWEST))
    }

    if !p.expectPeek(token.RPAREN) {
        return nil
    }

    return args
}

func (p *Parser) parseMemberExpression(object ast.Expression) ast.Expression {
    me := &ast.MemberExpression{
        Token:    p.curToken,
        Object:   object,
        Property: &ast.Identifier{Token: p.peekToken, Value: p.peekToken.Literal},
    }

    p.nextToken() // Consume the DOT token
    p.nextToken() // Move to property identifier

    return me
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

