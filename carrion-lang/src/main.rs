mod token;

use token::{lookup_ident, lookup_indent, TokenType};


fn main() {
    let ident = "spell";
    let token_type = lookup_ident(ident);
    println!("Token for '{}': {:?}", ident, token_type);

    let indent = "    ";
    let indent_token = lookup_indent(indent);
    dbg!(&indent_token);
    println!("Token for {:?}", indent_token);
    assert_eq!(lookup_ident("spell"), TokenType::Spell);
}
