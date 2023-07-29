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
		r.fails = append(r.fails, "run fail: "+err.Error())
	}

	return r.fails
}

type runtime struct {
	fails []string
	vm    *goja.Runtime
}

func newRuntime() *runtime {
	r := &runtime{
		vm: goja.New(),
	}
	global := r.vm.GlobalObject()

	global.DefineDataProperty("assert", r.vm.ToValue(r.assert), goja.FLAG_FALSE, goja.FLAG_FALSE, goja.FLAG_TRUE)

	return r
}

func (r *runtime) assert(a, b goja.Value) {
	if a != b && !a.Equals(b) {
		buff := &bytes.Buffer{}

		fmt.Fprintf(buff, "expected: %s\n\n", a)
		fmt.Fprintf(buff, "received: %s\n\n", b)

		buff.WriteString("callstack:\n")
		for _, frame := range r.vm.CaptureCallStack(-1, nil) {
			if frame.SrcName() != filename {
				buff.WriteString("[native]\n")
				continue
			}
			p := frame.Position()
			fmt.Fprintf(buff, "%d:%d | %s()\n", p.Line, p.Column, frame.FuncName())
		}

		r.fails = append(r.fails, buff.String())
	}
}
