package object_test

import (
	"monkey/object"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIntegerHashKey(t *testing.T) {
	hello1 := &object.Integer{Value: 1}
	hello2 := &object.Integer{Value: 1}
	diff1 := &object.Integer{Value: 2}
	diff2 := &object.Integer{Value: 2}

	assert.EqualValuesf(t, hello1.HashKey(), hello2.HashKey(), "integers with the same content have different hash keys")
	assert.EqualValuesf(t, diff1.HashKey(), diff2.HashKey(), "integers with the same content have different hash keys")
	assert.NotEqualValuesf(t, hello1.HashKey(), diff1.HashKey(), "integers with the different content have the same hash keys")
}
