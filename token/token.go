package token


type TokenType string

type Token struct{
  Type TokenType
  Literal string
}

const (
  ILLEGAL TokenType = "ILLEGAL"
  EOF TokenType = "EOF"
  IDENT TokenType = "IDENT"
  INT TokenType = "INT"
  FLOAT TokenType = "FLOAT"
  STRING TokenType = "STRING"

  ASSIGN TokenType = "="
  COLON TokenType = ":"
  COMMA TokenType = ","
  NEWLINE TokenType = "NEWLINE"
  INDENT TokenType = "INDENT"
  DEDENT TokenType = "DEDENT"
  LPAREN TokenType = "("
  RPAREN TokenType = ")"
  LBRAC TokenType = "["
  RBRAC TokenType = "]"
  ARROW TokenType = "->"

  PLUS TokenType = "+"
  MINUS TokenType = "-"
  ASTERISK TokenType = "*"
  SLASH TokenType = "/"
  MODULO TokenType = "%"

  LT TokenType = "<"
  GT TokenType = ">"

  //KeyWord
  SPELLBOOK TokenType = "SPELLBOOK"
  SPELL TokenType = "SPELL"
  BEGIN TokenType = "BEGIN"
  SHARED TokenType = "SHARED"
  IF TokenType = "IF"
  ELIF TokenType = "ELIF"
  ELSE TokenType = "ELSE"
  FOR TokenType = "FOR"
  IN TokenType = "IN"
  RETURN TokenType = "RETURN"
  RANGE TokenType = "RANGE"
  WHILE TokenType = "WHILE"
)
