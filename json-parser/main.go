package main

import "fmt"

type TokenType int

const (
	TOKEN_ILLEGAL TokenType = iota
	TOKEN_EOF
	TOKEN_LBRACE
	TOKEN_RBRACE
	TOKEN_STRING
	TOKEN_COLON
	TOKEN_COMMA
	TOKEN_NULL
	TOKEN_BOOLEAN
	TOKEN_NUMBER
	TOKEN_LBRACKET
	TOKEN_RBRACKET
)

const (
	NULL_LITERAL  = "null"
	TRUE_LITERAL  = "true"
	FALSE_LITERAL = "false"
)

type Token struct {
	Type    TokenType
	Literal string
}

// String representation for debugging
func (t Token) String() string {
	names := map[TokenType]string{
		TOKEN_ILLEGAL:  "ILLEGAL",
		TOKEN_EOF:      "EOF",
		TOKEN_LBRACE:   "LBRACE",
		TOKEN_RBRACE:   "RBRACE",
		TOKEN_STRING:   "STRING",
		TOKEN_COLON:    "COLON",
		TOKEN_COMMA:    "COMMA",
		TOKEN_NULL:     "NULL",
		TOKEN_BOOLEAN:  "BOOLEAN",
		TOKEN_NUMBER:   "NUMBER",
		TOKEN_LBRACKET: "LBRACKET",
		TOKEN_RBRACKET: "RBRACKET",
	}
	return fmt.Sprintf("Token{Type: %s, Literal: '%s'}", names[t.Type], t.Literal)
}

type Lexer struct {
	input        string
	position     int
	readPosition int
	ch           byte
}

func NewLexer(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition++
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

func (l *Lexer) readString() string {
	position := l.position + 1
	for {
		l.readChar()
		if l.ch == '"' || l.ch == 0 {
			break
		}
	}
	return l.input[position:l.position]
}

func (l *Lexer) readNull() string {
	position := l.position
	for l.position < len(l.input) && (l.ch >= 'a' && l.ch <= 'z') {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) readNumber() string {
	position := l.position
	for l.position < len(l.input) && (l.ch >= '0' && l.ch <= '9') {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) NextToken() Token {
	var tok Token

	l.skipWhitespace()

	switch l.ch {
	case '{':
		tok = Token{Type: TOKEN_LBRACE, Literal: string(l.ch)}
	case '}':
		tok = Token{Type: TOKEN_RBRACE, Literal: string(l.ch)}
	case ':':
		tok = Token{Type: TOKEN_COLON, Literal: string(l.ch)}
	case ',':
		tok = Token{Type: TOKEN_COMMA, Literal: string(l.ch)}
	case '"':
		tok.Literal = l.readString()
		tok.Type = TOKEN_STRING
	case 'n':
		literal := l.readNull()
		if literal == NULL_LITERAL {
			tok = Token{Type: TOKEN_NULL, Literal: literal}
		} else {
			tok = Token{Type: TOKEN_ILLEGAL, Literal: literal}
		}
		return tok
	case 't', 'f':
		literal := l.readNull()
		if literal == TRUE_LITERAL || literal == FALSE_LITERAL {
			tok = Token{Type: TOKEN_BOOLEAN, Literal: literal}
		} else {
			tok = Token{Type: TOKEN_ILLEGAL, Literal: literal}
		}
		return tok
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		literal := l.readNumber()
		tok = Token{Type: TOKEN_NUMBER, Literal: literal}
		return tok
	case '[':
		tok = Token{Type: TOKEN_LBRACKET, Literal: string(l.ch)}
	case ']':
		tok = Token{Type: TOKEN_RBRACKET, Literal: string(l.ch)}
	case 0:
		tok = Token{Type: TOKEN_EOF, Literal: ""}
	default:
		tok = Token{Type: TOKEN_ILLEGAL, Literal: string(l.ch)}
	}

	l.readChar()
	return tok
}

type Parser struct {
	lexer        *Lexer
	currentToken Token
	peekToken    Token
}

func NewParser(lexer *Lexer) *Parser {
	p := &Parser{lexer: lexer}
	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parser) nextToken() {
	p.currentToken = p.peekToken
	p.peekToken = p.lexer.NextToken()
}

func (p *Parser) expectPeek(t TokenType) error {
	if p.peekToken.Type == t {
		p.nextToken()
		return nil
	}
	return fmt.Errorf("expected next token to be %s, got %s instead", p.tokenTypeToString(t), p.tokenTypeToString(p.peekToken.Type))
}

// tokenTypeToString converts TokenType to string for error messages
func (p *Parser) tokenTypeToString(t TokenType) string {
	names := map[TokenType]string{
		TOKEN_LBRACE:   "'{'",
		TOKEN_RBRACE:   "'}'",
		TOKEN_LBRACKET: "'['",
		TOKEN_RBRACKET: "']'",
		TOKEN_STRING:   "STRING",
		TOKEN_COLON:    "':'",
		TOKEN_COMMA:    "','",
		TOKEN_EOF:      "EOF",
	}
	if name, exists := names[t]; exists {
		return name
	}
	return "UNKNOWN"
}

// currentTokenIs checks if the current token is of the given type
func (p *Parser) currentTokenIs(t TokenType) bool {
	return p.currentToken.Type == t
}

// Parse starts the parsing process and returns the result
func (p *Parser) Parse() (map[string]any, error) {
	if !p.currentTokenIs(TOKEN_LBRACE) {
		return nil, fmt.Errorf("expected JSON object to start with '{', got %s", p.currentToken.Literal)
	}
	return p.parseObject()
}

// parseObject parses a JSON object: { ... }
func (p *Parser) parseObject() (map[string]any, error) {
	obj := make(map[string]any)

	// Expect opening brace {
	if p.currentToken.Type != TOKEN_LBRACE {
		return nil, fmt.Errorf("expected '{', got %s", p.currentToken.Literal)
	}

	// Check if object is empty: { }
	if p.peekToken.Type == TOKEN_RBRACE {
		p.nextToken()   // consume the closing brace
		return obj, nil // Empty object
	}

	// Move to the first key
	p.nextToken()

	for {
		if !p.currentTokenIs(TOKEN_STRING) {
			return nil, fmt.Errorf("expected string key, got %s", p.currentToken.Literal)
		}

		key := p.currentToken.Literal

		if err := p.expectPeek(TOKEN_COLON); err != nil {
			return nil, err
		}

		p.nextToken()
		value, err := p.parseValue()
		if err != nil {
			return nil, err
		}

		obj[key] = value

		if err := p.expectPeek(TOKEN_RBRACE); err == nil {
			break // End of object
		}

		if err := p.expectPeek(TOKEN_COMMA); err != nil {
			return nil, fmt.Errorf("expected ',' or '}', got %s", p.peekToken.Literal)
		}

		p.nextToken()
	}

	return obj, nil
}

func (p *Parser) parseArray() ([]any, error) {
	arr := make([]any, 0)

	if p.currentToken.Type != TOKEN_LBRACKET {
		return nil, fmt.Errorf("expected '[', got %s", p.currentToken.Literal)
	}

	if p.peekToken.Type == TOKEN_RBRACKET {
		p.nextToken()
		return arr, nil
	}

	p.nextToken()

	for {
		value, err := p.parseValue()
		if err != nil {
			return nil, err
		}

		arr = append(arr, value)

		if err := p.expectPeek(TOKEN_RBRACKET); err == nil {
			break
		}

		if err := p.expectPeek(TOKEN_COMMA); err != nil {
			return nil, fmt.Errorf("expected ',' or ']', got %s", p.peekToken.Literal)
		}

		p.nextToken()
	}
	return arr, nil
}

func (p *Parser) parseValue() (any, error) {
	switch p.currentToken.Type {
	case TOKEN_STRING:
		return p.currentToken.Literal, nil
	case TOKEN_NULL:
		return nil, nil
	case TOKEN_BOOLEAN:
		switch p.currentToken.Literal {
		case TRUE_LITERAL:
			return true, nil
		case FALSE_LITERAL:
			return false, nil
		}
		return nil, fmt.Errorf("invalid boolean value: %s", p.currentToken.Literal)
	case TOKEN_NUMBER:
		return p.currentToken.Literal, nil
	case TOKEN_LBRACKET:
		return p.parseArray()
	case TOKEN_LBRACE:
		return p.parseObject()
	default:
		return nil, fmt.Errorf("unexpected value type: %s", p.currentToken.Literal)
	}
}

func main() {
	// 	input := `{
	//   "key1": true,
	//   "key2": false,
	//   "key3": null,
	//   "key4": "value",
	//   "key5": 101
	// }`

	input := `{"key": [1, 2, 3]}`

	// lexer := NewLexer(input)
	// fmt.Println("Tokens from lexer:")
	// for {
	// 	tok := lexer.NextToken()
	// 	fmt.Printf("%+v\n", tok)
	// 	if tok.Type == TOKEN_EOF {
	// 		break
	// 	}
	// }

	lexer := NewLexer(input)
	parser := NewParser(lexer)
	result, err := parser.Parse()
	if err != nil {
		fmt.Println("Parser error:", err)
		return
	}

	fmt.Printf("\nParsing result: %v\n", result)
	fmt.Printf("Result type: %T\n", result)
	fmt.Println("Success! Empty object parsed correctly.")
}
