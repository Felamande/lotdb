package jsvm

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"

	"github.com/Felamande/otto"
)

// type Object interface {
// 	Extend(funcName string, Func func(call otto.FunctionCall) otto.Value) error
// 	Obj() *otto.Object
// 	Type() ObjType
// }

type module struct {
	obj      *otto.Object
	methods  map[string]interface{}
	required bool
	register string
}

type builtin struct {
	obj *otto.Object
}

var modules map[string]*module
var builtins map[string]*builtin

var vm *otto.Otto

func init() {
	modules = make(map[string]*module)
	builtins = make(map[string]*builtin)
	vm = otto.New()
	vm.Set("require", require)

}

func Vm() *otto.Otto {
	return vm
}

func Builtin(obj string) *builtin {
	if b, exist := builtins[obj]; exist {
		return b
	}
	if _, exist := modules[obj]; exist {
		return nil
	}
	v, err := vm.Get(obj)
	if err != nil {
		return nil
	}
	return &builtin{v.Object()}

}

func (b *builtin) Extend(name string, fn interface{}) error {
	return b.obj.Set(name, fn)
}

func require(call otto.FunctionCall) otto.Value {
	mod := call.Argument(0).String()
	p, exist := modules[mod]
	if !exist {
		return otto.UndefinedValue()
	}

	if p.required {
		return p.obj.Value()
	}
	for name, method := range p.methods {
		p.obj.Set(name, method)
	}
	p.required = true

	return p.obj.Value()

}

func Module(name string) *module {
	_, file, _, _ := runtime.Caller(1)
	reg := filepath.Dir(file)
	if m, exist := modules[name]; exist {
		if m.register != reg {
			return nil
		}
		return m

	}
	o, _ := vm.Object(`({})`)

	vm.Set(name, o)

	m := &module{o, make(map[string]interface{}), false, reg}
	modules[name] = m
	return m
}

func (m *module) Extend(obj string, Func func(call otto.FunctionCall) otto.Value) {
	m.methods[obj] = Func
}
func (m *module) Obj() *otto.Object {
	return m.obj
}

func Run(src interface{}) error {
	n := src.(string)
	if fi, _ := os.Stat(n); fi == nil {
		_, err := vm.Run(n)
		return err
	}
	b, err := ioutil.ReadFile(n)
	if err != nil {
		return err
	}
	_, err = vm.Run(string(b))
	return err
}
