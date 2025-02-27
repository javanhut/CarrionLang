mod token;

use token::{lookup_ident, lookup_indent, TokenType};


fn main() {
    let ident = "spell";
    let token_type = lookup_ident(ident);
    println!("Token for '{}': {:?}", ident, token_type);
}
