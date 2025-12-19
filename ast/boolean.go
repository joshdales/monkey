package ast

import "monkey/token"

type Boolean struct {
	Token token.Token // The `token.TRUE` or `token.FALSE` value
	Value bool
}

func (b *Boolean) expressionNode() {

}

func (b *Boolean) TokenLiteral() string {
	return b.Token.Literal
}

func (b *Boolean) String() string {
	return b.Token.Literal
}
