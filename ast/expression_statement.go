package ast

import "monkey/token"

type ExpressionStatement struct {
	// The first token of the expression
	Token      token.Token
	Expression Expression
}

func (es *ExpressionStatement) statementNode() {

}

func (es *ExpressionStatement) TokenLiteral() string {
	return es.Token.Literal
}

func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.TokenLiteral()
	}

	return ""
}
