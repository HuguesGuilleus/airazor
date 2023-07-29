package check

import (
	"bytes"
	"fmt"

	"github.com/dop251/goja"
)

const filename = "check"

func Run(src string) []string {
	prog, err := goja.Compile(filename, src, true)
	if err != nil {
		return []string{"compile fail: " + err.Error()}
	}

	r := newRuntime()
	_, err = r.vm.RunProgram(prog)
	if err != nil {
		r.print = append(r.print, "run fail: "+err.Error())
	}

	return r.print
}

type runtime struct {
	print []string
	vm    *goja.Runtime
}

func newRuntime() *runtime {
	r := &runtime{
		vm: goja.New(),
	}

	r.define(nil, "log", r.log)

	r.define(nil, "assert", r.assert)
	assert := r.vm.Get("assert").ToObject(r.vm)
	r.define(assert, "true", r.assertTrue)
	r.define(assert, "false", r.assertFalse)
	r.define(assert, "null", r.assertNull)
	r.define(assert, "truely", r.assertTruely)
	r.define(assert, "falsy", r.assertFalsy)

	r.vm.RunString("Object.seal(assert);")

	return r
}

func (r *runtime) define(obj *goja.Object, name string, i any) {
	if obj == nil {
		obj = r.vm.GlobalObject()
	}
	err := obj.DefineDataProperty(name, r.vm.ToValue(i), goja.FLAG_FALSE, goja.FLAG_FALSE, goja.FLAG_TRUE)
	if err != nil {
		panic(err)
	}
}

// Log the value in the r.print strings.
func (r *runtime) log(values ...goja.Value) {
	buff := bytes.Buffer{}
	for _, v := range values {
		printValue(&buff, v)
		buff.WriteByte(' ')
	}
	r.print = append(r.print, buff.String())
}

func (r *runtime) assert(a, b goja.Value) {
	if a != b && !a.Equals(b) {
		buff := &bytes.Buffer{}

		buff.WriteString("expected: ")
		printValue(buff, a)
		buff.WriteString("\n\n")

		buff.WriteString("received: ")
		printValue(buff, b)
		buff.WriteString("\n\n")

		buff.WriteString("callstack:\n")
		for _, frame := range r.vm.CaptureCallStack(-1, nil) {
			if frame.SrcName() != filename {
				buff.WriteString("[native]\n")
				continue
			}
			p := frame.Position()
			fmt.Fprintf(buff, "%d:%d | %s()\n", p.Line, p.Column, frame.FuncName())
		}

		r.print = append(r.print, buff.String())
	}
}

func (r *runtime) assertTrue(v goja.Value) {
	r.assertSpecific("true", r.vm.ToValue(true).StrictEquals(v), v)
}
func (r *runtime) assertFalse(v goja.Value) {
	r.assertSpecific("false", r.vm.ToValue(false).StrictEquals(v), v)
}
func (r *runtime) assertNull(v goja.Value) {
	r.assertSpecific("null", goja.Null().StrictEquals(v), v)
}

func (r *runtime) assertTruely(v goja.Value) {
	constructor, _ := goja.AssertConstructor(r.vm.Get("Boolean"))
	b, _ := constructor(nil, v)
	r.assertSpecific("truely", r.vm.ToValue(true).Equals(b), v)
}
func (r *runtime) assertFalsy(v goja.Value) {
	constructor, _ := goja.AssertConstructor(r.vm.Get("Boolean"))
	b, _ := constructor(nil, v)
	r.assertSpecific("falsy", r.vm.ToValue(false).Equals(b), v)
}

func (r *runtime) assertSpecific(name string, equal bool, v goja.Value) {
	if !equal {
		buff := bytes.Buffer{}
		buff.WriteString("assert ")
		buff.WriteString(name)
		buff.WriteString(", but found:\n")
		printValue(&buff, v)

		r.print = append(r.print, buff.String())
	}
}

func printValue(buff *bytes.Buffer, value goja.Value) {
	buff.WriteString(value.String())
}
