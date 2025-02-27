mod token;
use token;

#[test]
fn test_lookup_ident_keyword() {
    assert_eq!(token::lookup_ident("spell"), token::TokenType::Spell);
}

#[test]
fn test_lookup_spellbook_terms(){
    assert_eq!(token::lookup_ident("init"), token::TokenType::Init);
}
