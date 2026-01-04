package evaluator_test

import (
	"monkey/object"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestQuote(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`quote(5)`, `5`},
		{`quote(5 + 8)`, `(5 + 8)`},
		{`quote(foobar)`, `foobar`},
		{`quote(foobar + barfoo)`, `(foobar + barfoo)`},
	}

	for _, tt := range tests {
		evaluated := testEval(t, tt.input)
		quote, ok := evaluated.(*object.Quote)
		require.Truef(t, ok, "expected *object.Quote, got %T (%+v)", evaluated, evaluated)
		require.NotNil(t, quote.Node)
		assert.Equal(t, tt.expected, quote.Node.String())
	}
}
