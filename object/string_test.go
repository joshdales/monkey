package object_test

import (
	"monkey/object"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStringHashKey(t *testing.T) {
	hello1 := &object.String{Value: "Hello World"}
	hello2 := &object.String{Value: "Hello World"}
	diff1 := &object.String{Value: "my name is johnny"}
	diff2 := &object.String{Value: "my name is johnny"}

	assert.EqualValuesf(t, hello1.HashKey(), hello2.HashKey(), "strings with the same content have different hash keys")
	assert.EqualValuesf(t, diff1.HashKey(), diff2.HashKey(), "strings with the same content have different hash keys")
	assert.NotEqualValuesf(t, hello1.HashKey(), diff1.HashKey(), "strings with the different content have the same hash keys")
}
