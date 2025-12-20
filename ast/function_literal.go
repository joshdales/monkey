package ast

import (
	"bytes"
	"monkey/token"
)

type FunctionLiteral struct {
	Token      token.Token // The `token.FUNCTION` token
	Parameters []*Identifier
	Body       *BlockStatement
}

func (fl *FunctionLiteral) expressionNode() {

}

func (fl *FunctionLiteral) TokenLiteral() string {
	return fl.Token.Literal
}

func (fl *FunctionLiteral) String() string {
	var out bytes.Buffer

	return out.String()
}
