package object_test

import (
	"monkey/object"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBooleanHashKey(t *testing.T) {
	true1 := &object.Boolean{Value: true}
	true2 := &object.Boolean{Value: true}
	false1 := &object.Boolean{Value: false}
	false2 := &object.Boolean{Value: false}

	assert.EqualValuesf(t, true1.HashKey(), true2.HashKey(), "booleans with the same content have different hash keys")
	assert.EqualValuesf(t, false1.HashKey(), false2.HashKey(), "booleans with the same content have different hash keys")
	assert.NotEqualValuesf(t, true1.HashKey(), false1.HashKey(), "booleans with the different content have the same hash keys")
}
