// ast/ast.go
package ast

import (
    "carrionlang/token"
    "strings"
)

type Node interface {
    TokenLiteral() string
    String() string
}

type Statement interface {
    Node
    statementNode()
}

type Expression interface {
    Node
    expressionNode()
}

type Program struct {
    Statements []Statement
}

func (p *Program) TokenLiteral() string {
    if len(p.Statements) > 0 {
        return p.Statements[0].TokenLiteral()
    } else {
        return ""
    }
}

func (p *Program) String() string {
    var out strings.Builder
    for _, s := range p.Statements {
        out.WriteString(s.String())
    }
    return out.String()
}

// Statements

type ExpressionStatement struct {
    Token      token.Token // First token of the expression
    Expression Expression
}

func (es *ExpressionStatement) statementNode()       {}
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal }
func (es *ExpressionStatement) String() string       { return es.Expression.String() }

// Variable Declaration

type VariableDeclaration struct {
    Token    token.Token // The IDENT token
    Name     *Identifier
    TypeHint *Identifier // May be nil
    Value    Expression
}

func (vd *VariableDeclaration) statementNode()       {}
func (vd *VariableDeclaration) TokenLiteral() string { return vd.Token.Literal }
func (vd *VariableDeclaration) String() string {
    var out strings.Builder
    out.WriteString(vd.Name.String())
    if vd.TypeHint != nil {
        out.WriteString(":")
        out.WriteString(vd.TypeHint.String())
    }
    out.WriteString(" = ")
    out.WriteString(vd.Value.String())
    return out.String()
}

// Spellbook Declaration

type SpellbookDeclaration struct {
    Token token.Token // The 'spellbook' token
    Name  *Identifier
    Body  []Statement
}

func (sd *SpellbookDeclaration) statementNode()       {}
func (sd *SpellbookDeclaration) TokenLiteral() string { return sd.Token.Literal }
func (sd *SpellbookDeclaration) String() string {
    var out strings.Builder
    out.WriteString("spellbook ")
    out.WriteString(sd.Name.String())
    out.WriteString(":\n")
    for _, stmt := range sd.Body {
        out.WriteString("    ")
        out.WriteString(stmt.String())
        out.WriteString("\n")
    }
    return out.String()
}

// Spell Declaration

type SpellDeclaration struct {
    Token      token.Token // The 'spell' token
    Name       *Identifier
    Parameters []*Identifier
    Body       *BlockStatement
    ReturnType *Identifier // May be nil
}

func (sd *SpellDeclaration) statementNode()       {}
func (sd *SpellDeclaration) TokenLiteral() string { return sd.Token.Literal }
func (sd *SpellDeclaration) String() string {
    var out strings.Builder
    out.WriteString("spell ")
    out.WriteString(sd.Name.String())
    out.WriteString("(")
    params := []string{}
    for _, p := range sd.Parameters {
        params = append(params, p.String())
    }
    out.WriteString(strings.Join(params, ", "))
    out.WriteString(")")
    if sd.ReturnType != nil {
        out.WriteString(" -> ")
        out.WriteString(sd.ReturnType.String())
    }
    out.WriteString(":\n")
    out.WriteString(sd.Body.String())
    return out.String()
}

// Return Statement

type ReturnStatement struct {
    Token       token.Token // The 'return' token
    ReturnValue Expression
}

func (rs *ReturnStatement) statementNode()       {}
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal }
func (rs *ReturnStatement) String() string {
    var out strings.Builder
    out.WriteString("return ")
    if rs.ReturnValue != nil {
        out.WriteString(rs.ReturnValue.String())
    }
    return out.String()
}

// Block Statement

type BlockStatement struct {
    Token      token.Token // The INDENT token or starting token of the block
    Statements []Statement
}

func (bs *BlockStatement) statementNode()       {}
func (bs *BlockStatement) TokenLiteral() string { return bs.Token.Literal }
func (bs *BlockStatement) String() string {
    var out strings.Builder
    for _, s := range bs.Statements {
        out.WriteString(s.String())
        out.WriteString("\n")
    }
    return out.String()
}

// Expressions

type Identifier struct {
    Token token.Token // The IDENT token
    Value string
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
func (i *Identifier) String() string       { return i.Value }

// Integer Literal

type IntegerLiteral struct {
    Token token.Token
    Value int64
}

func (il *IntegerLiteral) expressionNode()      {}
func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Literal }
func (il *IntegerLiteral) String() string       { return il.Token.Literal }

// String Literal

type StringLiteral struct {
    Token token.Token
    Value string
}

func (sl *StringLiteral) expressionNode()      {}
func (sl *StringLiteral) TokenLiteral() string { return sl.Token.Literal }
func (sl *StringLiteral) String() string       { return "\"" + sl.Value + "\"" }

// Prefix Expression

type PrefixExpression struct {
    Token    token.Token // The prefix token, e.g., '!'
    Operator string
    Right    Expression
}

func (pe *PrefixExpression) expressionNode()      {}
func (pe *PrefixExpression) TokenLiteral() string { return pe.Token.Literal }
func (pe *PrefixExpression) String() string {
    var out strings.Builder
    out.WriteString("(")
    out.WriteString(pe.Operator)
    out.WriteString(pe.Right.String())
    out.WriteString(")")
    return out.String()
}

// Infix Expression

type InfixExpression struct {
    Token    token.Token // The operator token, e.g., '+'
    Left     Expression
    Operator string
    Right    Expression
}

func (ie *InfixExpression) expressionNode()      {}
func (ie *InfixExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *InfixExpression) String() string {
    var out strings.Builder
    out.WriteString("(")
    out.WriteString(ie.Left.String())
    out.WriteString(" " + ie.Operator + " ")
    out.WriteString(ie.Right.String())
    out.WriteString(")")
    return out.String()
}

// Call Expression

type CallExpression struct {
    Token     token.Token // The '(' token
    Function  Expression  // Identifier or FunctionLiteral
    Arguments []Expression
}

func (ce *CallExpression) expressionNode()      {}
func (ce *CallExpression) TokenLiteral() string { return ce.Token.Literal }
func (ce *CallExpression) String() string {
    var out strings.Builder
    args := []string{}
    for _, a := range ce.Arguments {
        args = append(args, a.String())
    }
    out.WriteString(ce.Function.String())
    out.WriteString("(")
    out.WriteString(strings.Join(args, ", "))
    out.WriteString(")")
    return out.String()
}

// Member Expression

type MemberExpression struct {
    Token     token.Token // The '.' token
    Object    Expression
    Property  *Identifier
}

func (me *MemberExpression) expressionNode()      {}
func (me *MemberExpression) TokenLiteral() string { return me.Token.Literal }
func (me *MemberExpression) String() string {
    var out strings.Builder
    out.WriteString(me.Object.String())
    out.WriteString(".")
    out.WriteString(me.Property.String())
    return out.String()
}

