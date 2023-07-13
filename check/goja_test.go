package check

import (
	"encoding/hex"
	"testing"

	"github.com/dop251/goja"
	"github.com/stretchr/testify/assert"
)

func TestGojaFunctionAdd(t *testing.T) {
	vm := goja.New()
	vm.GlobalObject().Set("add", func(a, b int) int { return a + b })
	v, err := vm.RunString(`add(1,2)`)
	assert.NoError(t, err)
	assert.Equal(t, 3, int(v.ToInteger()))
}

func TestGojaFunctionBytes(t *testing.T) {
	vm := goja.New()
	vm.GlobalObject().Set("hex", func(src []byte) string { return hex.EncodeToString(src) })
	v, err := vm.RunString(`hex(new Uint8Array([1,2,3]).buffer)`)
	assert.NoError(t, err)
	assert.Equal(t, "010203", v.String())
}

func TestGojaBytes(t *testing.T) {
	vm := goja.New()
	vm.GlobalObject().Set("b", func() []byte { return []byte("abc") })
	v, err := vm.RunString(`Object.prototype.toString.call(b())`)
	assert.NoError(t, err)
	assert.Equal(t, "[object Array]", v.String())
}

func TestGojaGetObjectField(t *testing.T) {
	vm := goja.New()
	_, err := vm.RunString(`a = {v:36}; b = {v:36};`)
	assert.NoError(t, err)

	a := vm.GlobalObject().Get("a").ToObject(vm).Get("v")
	b := vm.Get("b").ToObject(vm).Get("v")
	assert.True(t, a.Equals(b))
}

func TestGojaThrow(t *testing.T) {
	v, err := goja.New().RunString("throw 'fuck!';")
	assert.Nil(t, v)
	assert.NotNil(t, err)

	assert.Equal(t, "fuck!", err.(*goja.Exception).Value().Export())
}

func TestGojaPassValue(t *testing.T) {
	vm := goja.New()
	vm.GlobalObject().Set("i", 5)
	vm.GlobalObject().Set("response", map[string]any{"code": 200, "body": "Hello World!"})

	v, err := vm.RunString("response.code + i")
	assert.NoError(t, err)
	assert.Equal(t, 205, int(v.ToInteger()))
}

func TestGojaAdd4(t *testing.T) {
	vm := goja.New()
	v, err := vm.RunString("2 + 2")
	assert.NoError(t, err)
	assert.Equal(t, 4, int(v.ToInteger()))
}

func TestGojaEmptyScript(t *testing.T) {
	v, err := goja.New().RunString("")
	assert.True(t, v.StrictEquals(goja.Undefined()))
	assert.NoError(t, err)
}
