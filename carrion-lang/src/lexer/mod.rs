use crate::token;

pub struct Lexer {
    input: String,
    position: usize,
    read_postion: usize,
    ch: u8,
}

impl Lexer {
    pub fn new(input: String) -> Self {
       let mut l = Self {
            input,
            position: 0,
            read_postion: 0,
            ch: b'\0'
        }
       l.read_char();
       l
    }


    pub fn read_char(&mut self){
        if self.read_postion >= self.input.len() {
            self.ch = 0;
        } else {
            self.ch = self.input.as_bytes()[self.read_postion];
        }
        self.position = self.read_postion;
        self.read_postion += 1;
    }

    pub fn next_token(&mut self) -> token::Token {
        let tok = match self.ch {
            b'=' => token::new_token(token::TokenType::Assign, self.ch as char),
            b'+' => token::new_token(token::TokenType::Plus, self.ch as char),
            0    => token::Token {
                token_type: token::TokenType::EOF,
                literal: String::new(),
            },
            _   => token::new_token(token::TokenType::Illegal, self.ch),
        };
        self.read_char();
        tok
    }
}

