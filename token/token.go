// token/token.go
package token

type TokenType string

type Token struct {
    Type    TokenType
    Literal string
}

const (
    ILLEGAL TokenType = "ILLEGAL"
    EOF     TokenType = "EOF"

    IDENT   TokenType = "IDENT"
    INT     TokenType = "INT"
    FLOAT   TokenType = "FLOAT"
    STRING  TokenType = "STRING"

    ASSIGN  TokenType = "="
    COLON   TokenType = ":"
    COMMA   TokenType = ","
    INDENT  TokenType = "INDENT"
    DEDENT  TokenType = "DEDENT"
    LPAREN  TokenType = "("
    RPAREN  TokenType = ")"
    ARROW   TokenType = "->"

    PLUS     TokenType = "+"
    MINUS    TokenType = "-"
    ASTERISK TokenType = "*"
    SLASH    TokenType = "/"

    LT TokenType = "<"
    GT TokenType = ">"
    EQ TokenType = "=="
    NOT_EQ TokenType = "!="

    DOT  TokenType = "."
    BANG TokenType = "!"
    NEWLINE TokenType = "NEWLINE"
    // Keywords
    SPELLBOOK TokenType = "SPELLBOOK"
    SPELL     TokenType = "SPELL"
    BEGIN     TokenType = "BEGIN"
    SHARED    TokenType = "SHARED"
    FOR       TokenType = "FOR"
    IN        TokenType = "IN"
    RETURN    TokenType = "RETURN"
)
