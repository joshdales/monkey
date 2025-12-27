package ast

import (
	"bytes"
	"monkey/token"
	"strings"
)

type ArrayLiteral struct {
	Token    token.Token // the `token.LBRACKET` token
	Elements []Expression
}

func (al *ArrayLiteral) expressionNode() {}

func (al *ArrayLiteral) TokenLiteral() string {
	return al.Token.Literal
}

func (al *ArrayLiteral) String() string {
	var out bytes.Buffer

	elements := make([]string, 0, len(al.Elements))
	for _, el := range al.Elements {
		elements = append(elements, el.String())
	}

	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")

	return out.String()
}
