package check

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAssert(t *testing.T) {
	assert.Nil(t, Run(`assert(1, 1)`))
	assert.Nil(t, Run(`assert(1.12345, 1.12345)`))
	assert.Nil(t, Run(`assert("a", "a")`))
	assert.Nil(t, Run(`assert(true, true)`))
	assert.Nil(t, Run(`assert(false, false)`))
	assert.Nil(t, Run(`assert()`))

	// Simple fail
	assert.Equal(t, []string{"expected: 1\n\nreceived: 2\n\ncallstack:\n[native]\n1:20 | f()\n3:2 | <anonymous>()\n"},
		Run(`function f(){assert(1, 2);}
function g() {f();}
f();
	`))
}
