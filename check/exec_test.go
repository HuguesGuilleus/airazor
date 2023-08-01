package check

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAssert(t *testing.T) {
	// assert function
	assert.Nil(t, Run(`assert(1, 1)`, nil))
	assert.Nil(t, Run(`assert(1.12345, 1.12345)`, nil))
	assert.Nil(t, Run(`assert("a", "a")`, nil))
	assert.Nil(t, Run(`assert(true, true)`, nil))
	assert.Nil(t, Run(`assert(false, false)`, nil))
	assert.Nil(t, Run(`assert()`, nil))

	// Simple fail
	assert.Equal(t, []string{"expected: 1\n\nreceived: 2\n\ncallstack:\n[native]\n1:20 | f()\n3:2 | <anonymous>()\n"},
		Run(`function f(){assert(1, 2);}
function g() {f();}
f();
	`, nil))

	// assert specific strict function
	assert.Nil(t, Run(`assert.true(true)`, nil))
	assert.Nil(t, Run(`assert.false(false)`, nil))
	assert.Nil(t, Run(`assert.null(null)`, nil))

	assert.Equal(t, "assert true, but found:\n1", Run(`assert.true(1)`, nil)[0])
	assert.Equal(t, "assert false, but found:\n0", Run(`assert.false(0)`, nil)[0])
	assert.Equal(t, "assert null, but found:\n1", Run(`assert.null(1)`, nil)[0])
	assert.Equal(t, "assert null, but found:\nundefined", Run(`assert.null(undefined)`, nil)[0])

	// assert specific flexible function
	assert.Nil(t, Run(`assert.truely(true)`, nil))
	assert.Nil(t, Run(`assert.truely(1)`, nil))
	assert.Nil(t, Run(`assert.truely({})`, nil))
	assert.Nil(t, Run(`assert.falsy(false)`, nil))
	assert.Nil(t, Run(`assert.falsy(0)`, nil))
	assert.Nil(t, Run(`assert.falsy(null)`, nil))

	assert.Equal(t, "assert truely, but found:\n0", Run(`assert.truely(0)`, nil)[0])
	assert.Equal(t, "assert falsy, but found:\n1", Run(`assert.falsy(1)`, nil)[0])

	// check that assert function object if sealed.
	assert.NotEmpty(t, Run("assert.a = 42;", nil))
}
