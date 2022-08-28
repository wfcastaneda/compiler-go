package lexer

import "gopherlang-interpreter/token"

type Lexer struct {
    input         string
    position      int  // current position in the input (current charecter)
    readPosition  int  // current reading position in the input (after current charecter
    ch            byte // current charection being examined
}

// Get a Lexer for the input string
func New(input string) *Lexer {
    l := &Lexer{input: input}
    l.readChar()
    return l
}

// Lexer helper to get the next character and advance position in the input
func (l *Lexer) readChar() {
    if l.readPosition >= len(l.input) {
        l.ch = 0
    } else {
        l.ch = l.input[l.readPosition]
    }
    l.position = l.readPosition
    l.readPosition += 1
}

// Lexer helper to read through a set of sequential characters and return them, used to identify if a set of
// characters are an identifier or keyword
func (l *Lexer) readIdentifier() string {
    position := l.position
    for isLetter(l.ch) {
        l.readChar()
    }
    return l.input[position:l.position]
}

// Lexer helper to check whether a character is a letter
func isLetter(ch byte) bool {
    return 'a' <= ch && ch <= 'z' ||  'A' <= ch && ch <= 'Z' || ch == '_'
}

// Lexer helper which reads through a sequential set of digits, used to identify a number
func (l *Lexer) readNumber() string {
    position := l.position
    for isDigit(l.ch) {
        l.readChar()
    }
    return l.input[position:l.position]
}

func isDigit(ch byte) bool {
    return '0' <= ch && ch <= '9'
}

// Lexer helper to skip whitespace characters
func (l *Lexer) skipWhitespace() {
    for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
        l.readChar()
    }
}

func (l *Lexer) NextToken() token.Token {
    var tok token.Token

    // Skip whitespace as it only acts as a separator of tokens
    l.skipWhitespace()

    switch l.ch {
    case '=':
        tok = newToken(token.ASSIGN, l.ch)
    case ';':
        tok = newToken(token.SEMICOLON, l.ch)
    case '(':
        tok = newToken(token.LPAREN, l.ch)
    case ')':
        tok = newToken(token.RPAREN, l.ch)
    case ',':
        tok = newToken(token.COMMA, l.ch)
    case '+':
        tok = newToken(token.PLUS, l.ch)
    case '{':
        tok = newToken(token.LBRACE, l.ch)
    case '}':
        tok = newToken(token.RBRACE, l.ch)
    case 0:
        tok.Literal = ""
        tok.Type = token.EOF
    default:
        if isLetter(l.ch) {
            // If the character is a letter, use Lexer.readIdentifier to get the full literal up to the next empty character
            // check if the Literal is a KEYWORD or IDENT type by calling token.LookupIdent
            tok.Literal = l.readIdentifier()
            tok.Type = token.LookupIdent(tok.Literal)
            
            // Early return required because Lexer.readIdentifier advances position using readChar()
            return tok
        } else if isDigit(l.ch) {
            // If the character is a digit, use Lexer.readNumber to get the full number literal up to the next emty character
            tok.Type = token.INT
            tok.Literal = l.readNumber()
            return tok
        } else {
            // Otherwise, the character is ILLEGAL
            tok = newToken(token.ILLEGAL, l.ch)
        }
    }

    l.readChar()
    return tok
}

func newToken(tokenType token.TokenType, ch byte) token.Token {
    return token.Token{Type: tokenType, Literal: string(ch)}
}

