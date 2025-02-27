#[derive(Debug)]
#[derive(PartialEq)]
pub enum TokenType {
     // Special tokens
    Illegal,
    EOF,
    Newline,
    Indent,
    Dedent,

    // Identifiers and literals
    Ident,
    Int,
    Float,
    String,
    DocString,

    // Operators
    Assign,          // "="
    Plus,            // "+"
    Minus,           // "-"
    Asterisk,        // "*"
    Slash,           // "/"
    Mod,             // "%"
    Exponent,        // "**"
    Increment,       // "+="
    Decrement,       // "-="
    MultAssgn,       // "*="
    DivAssgn,        // "/="
    PlusIncrement,   // "++"
    MinusDecrement,  // "--"
    Eq,              // "=="
    NotEq,           // "!="
    Lt,              // "<"
    Gt,              // ">"
    Le,              // "<="
    Ge,              // ">="
    Bang,            // "!"
    Ampersand,       // "&"
    Hash,            // "#"
    At,              // "@"

    // Delimiters
    Comma,       // ","
    Colon,       // ":"
    Pipe,        // "|"
    Dot,         // "."
    LShift,      // "<<"
    RShift,      // ">>"
    Xor,         // "^"
    Tilde,       // "~"
    LParen,      // "("
    RParen,      // ")"
    LBrace,      // "{"
    RBrace,      // "}"
    LBrack,      // "["
    RBrack,      // "]"
    Underscore,  // "_"

    // Keywords
    Init,
    Self_,
    Spell,
    Spellbook,
    True_,
    False_,
    If,
    Otherwise,
    Else,
    For,
    In,
    While,
    Stop,
    Skip,
    Ignore,
    Return,
    Import,
    Match,
    Case,
    Attempt,
    Resolve,
    Ensnare,
    Raise,
    As,
    Arcane,
    ArcaneSpell,
    Super,
    FString,
    Check,
    None,
    And,
    Or,
    Not, // keyword NOT 
}

pub struct Token {
    pub token_type: TokenType,
    pub literal: String,
}


pub fn lookup_ident(ident: &str) -> TokenType {
    match ident {
        "import"    => TokenType::Import,
        "match"     => TokenType::Match,
        "case"      => TokenType::Case,
        "spell"     => TokenType::Spell,
        "self"      => TokenType::Self_,
        "init"      => TokenType::Init,
        "spellbook" => TokenType::Spellbook,
        "True"      => TokenType::True_,
        "False"     => TokenType::False_,
        "if"        => TokenType::If,
        "otherwise" => TokenType::Otherwise,
        "else"      => TokenType::Else,
        "for"       => TokenType::For,
        "in"        => TokenType::In,
        "while"     => TokenType::While,
        "stop"      => TokenType::Stop,
        "skip"      => TokenType::Skip,
        "ignore"    => TokenType::Ignore,
        "and"       => TokenType::And,
        "or"        => TokenType::Or,
        "not"       => TokenType::Not,
        "return"    => TokenType::Return,
        "attempt"   => TokenType::Attempt,
        "resolve"   => TokenType::Resolve,
        "ensnare"   => TokenType::Ensnare,
        "raise"     => TokenType::Raise,
        "as"        => TokenType::As,
        "arcane"    => TokenType::Arcane,
        "arcanespell" => TokenType::ArcaneSpell,
        "super"     => TokenType::Super,
        "check"     => TokenType::Check,
        "None"      => TokenType::None,
        _           => TokenType::Ident,
    }
}

pub fn lookup_indent(ident: &str) -> TokenType {
    match ident.len() {
        0      => TokenType::Dedent,
        4 | 8  => TokenType::Ident,
        _      => TokenType::Illegal,
    }
}

pub fn new_token(token_type: TokenType, ch: u8) -> Token {
    Token {
        token_type,
        literal: (ch as char).to_string(),
    }
}

