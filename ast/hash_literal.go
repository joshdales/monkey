package ast

import (
	"bytes"
	"fmt"
	"monkey/token"
	"strings"
)

type HashLiteral struct {
	Token token.Token // the `token.LBRACE` token
	Pairs map[Expression]Expression
}

func (hl *HashLiteral) expressionNode() {}

func (hl *HashLiteral) TokenLiteral() string {
	return hl.Token.Literal
}

func (hl *HashLiteral) String() string {
	var out bytes.Buffer

	pairs := make([]string, len(hl.Pairs))
	for key, value := range hl.Pairs {
		pairs = append(pairs, fmt.Sprintf("%s: %s", key.String(), value.String()))
	}

	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")

	return out.String()
}
