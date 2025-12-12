package lexer_test

import (
	"monkey/lexer"
	"monkey/token"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNextToken(t *testing.T) {
	input := "=+(){},;"
	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.ASSIGN, "="},
		{token.PLUS, "+"},
		{token.LPAREN, "("},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.RBRACE, "}"},
		{token.COMMA, ","},
		{token.SEMICOLON, ";"},
		{token.EOF, ""},
	}

	lex := lexer.New(input)
	for _, tt := range tests {
		tok := lex.NextToken()
		require.EqualValues(t, tt.expectedType, tok.Type)
		require.EqualValues(t, tt.expectedLiteral, tok.Literal)
	}
}
