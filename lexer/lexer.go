// lexer/lexer.go
package lexer

import (
    "carrionlang/token"
    "unicode"
)

type Lexer struct {
    input        string
    position     int  // current position in input (points to current char)
    readPosition int  // current reading position (after current char)
    ch           rune // current char under examination
    indentStack  []int
    lineStart    bool
}

func New(input string) *Lexer {
    l := &Lexer{
        input:       input,
        indentStack: []int{0},
        lineStart:   true,
    }
    l.readChar()
    return l
}

func (l *Lexer) readChar() {
    if l.readPosition >= len(l.input) {
        l.ch = 0 // ASCII code for NUL character
    } else {
        l.ch = rune(l.input[l.readPosition])
    }
    l.position = l.readPosition
    l.readPosition++
}

func (l *Lexer) NextToken() token.Token {
    var tok token.Token

    if l.lineStart {
        indentLevel := l.countIndent()
        lastIndent := l.indentStack[len(l.indentStack)-1]

        if indentLevel > lastIndent {
            l.indentStack = append(l.indentStack, indentLevel)
            tok.Type = token.INDENT
            tok.Literal = ""
            l.lineStart = false
            return tok
        }

        for indentLevel < lastIndent {
            l.indentStack = l.indentStack[:len(l.indentStack)-1]
            lastIndent = l.indentStack[len(l.indentStack)-1]
            tok.Type = token.DEDENT
            tok.Literal = ""
            return tok
        }

        l.lineStart = false
    }

    l.skipWhitespace()

    switch l.ch {
    case '=':
        tok = newToken(token.ASSIGN, l.ch)
    case ':':
        tok = newToken(token.COLON, l.ch)
    case ',':
        tok = newToken(token.COMMA, l.ch)
    case '(':
        tok = newToken(token.LPAREN, l.ch)
    case ')':
        tok = newToken(token.RPAREN, l.ch)
    case '+':
        tok = newToken(token.PLUS, l.ch)
    case '-':
        if l.peekChar() == '>' {
            ch := l.ch
            l.readChar()
            tok.Type = token.ARROW
            tok.Literal = string(ch) + string(l.ch)
        } else {
            tok = newToken(token.MINUS, l.ch)
        }
    case '*':
        tok = newToken(token.ASTERISK, l.ch)
    case '/':
        tok = newToken(token.SLASH, l.ch)
    case '<':
        tok = newToken(token.LT, l.ch)
    case '>':
        tok = newToken(token.GT, l.ch)
    case '.':
        tok = newToken(token.DOT, l.ch)
    case '!':
        tok = newToken(token.BANG, l.ch) // Handle BANG token
    case '\n':
        tok.Type = token.NEWLINE
        tok.Literal = ""
        l.lineStart = true
        return tok
    case '"':
        tok.Type = token.STRING
        tok.Literal = l.readString()
        return tok
    case 0:
        tok.Literal = ""
        tok.Type = token.EOF
        return tok
    default:
        if isLetter(l.ch) {
            tok.Literal = l.readIdentifier()
            tok.Type = lookupIdent(tok.Literal)
            return tok
        } else if isDigit(l.ch) {
            tok.Type, tok.Literal = l.readNumber()
            return tok
        } else {
            tok = newToken(token.ILLEGAL, l.ch)
        }
    }
    l.readChar()
    return tok
}

func (l *Lexer) skipWhitespace() {
    for l.ch == ' ' || l.ch == '\t' || l.ch == '\r' {
        l.readChar()
    }
}

func (l *Lexer) countIndent() int {
    count := 0
    for l.ch == ' ' {
        count++
        l.readChar()
    }
    return count
}

func (l *Lexer) readIdentifier() string {
    position := l.position
    for isLetter(l.ch) || isDigit(l.ch) {
        l.readChar()
    }
    return l.input[position:l.position]
}

func (l *Lexer) readNumber() (token.TokenType, string) {
    position := l.position
    tokType := token.INT
    hasDot := false
    for isDigit(l.ch) || l.ch == '.' {
        if l.ch == '.' {
            if hasDot {
                break // Second dot encountered, stop reading number
            }
            hasDot = true
            tokType = token.FLOAT
        }
        l.readChar()
    }
    return tokType, l.input[position:l.position]
}

func (l *Lexer) readString() string {
    l.readChar() // skip opening quote
    position := l.position
    for l.ch != '"' && l.ch != 0 {
        l.readChar()
    }
    str := l.input[position:l.position]
    l.readChar() // skip closing quote
    return str
}

func (l *Lexer) peekChar() rune {
    if l.readPosition >= len(l.input) {
        return 0
    } else {
        return rune(l.input[l.readPosition])
    }
}

func newToken(tokenType token.TokenType, ch rune) token.Token {
    return token.Token{Type: tokenType, Literal: string(ch)}
}

func isLetter(ch rune) bool {
    return unicode.IsLetter(ch) || ch == '_'
}

func isDigit(ch rune) bool {
    return unicode.IsDigit(ch)
}

func lookupIdent(ident string) token.TokenType {
    keywords := map[string]token.TokenType{
        "spellbook": token.SPELLBOOK,
        "spell":     token.SPELL,
        "begin":     token.BEGIN,
        "shared":    token.SHARED,
        "for":       token.FOR,
        "in":        token.IN,
        "return":    token.RETURN,
    }
    if tok, ok := keywords[ident]; ok {
        return tok
    }
    return token.IDENT
}
