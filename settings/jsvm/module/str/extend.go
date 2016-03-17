package str

import (
	"strings"

	"github.com/Felamande/lotdb/settings/jsvm"
	"github.com/Felamande/otto"
)

func init() {
	if b := jsvm.Builtin("String"); b != nil {
		b.Extend("split", split)
	}

}

func split(call otto.FunctionCall) otto.Value {
	if !call.This.IsString() {

		return call.This
	}
	str := call.This.String()
	sep := call.Argument(0).String()
	splits := strings.Split(str, sep)
	v, err := otto.ToValue(splits)
	if err != nil {

	}
	return v
}
