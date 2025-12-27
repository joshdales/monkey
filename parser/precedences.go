package parser

import "monkey/token"

const (
	_ int = iota
	LOWEST
	EQUALS      // `==`
	LESSGREATER // `>` or `<`
	SUM         // `+`
	PRODUCT     // `*`
	PREFIX      // `-X` or `!X`
	CALL        // `myfunction(X)`
	INDEX       // `array[index]`
)

var precedences = map[token.TokenType]int{
	token.EQ:       EQUALS,
	token.NOT_EQ:   EQUALS,
	token.LT:       LESSGREATER,
	token.GT:       LESSGREATER,
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.SLASH:    PRODUCT,
	token.ASTERISK: PRODUCT,
	token.LPAREN:   CALL,
	token.LBRACKET: INDEX,
}

func (p *Parser) curPrecedence() int {
	if pre, ok := precedences[p.curToken.Type]; ok {
		return pre
	}

	return LOWEST
}

func (p *Parser) peekPrecedence() int {
	if pre, ok := precedences[p.peekToken.Type]; ok {
		return pre
	}

	return LOWEST
}
