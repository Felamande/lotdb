package gofunc

import (
	"runtime"
	"strings"

	"github.com/Felamande/lotdb/settings/jsvm"
	"github.com/Felamande/otto"
)

func init() {
	p := jsvm.Module("go")
	p.Extend("version", version)
	p.Extend("ostype", ostype)
}

func ostype(otto.FunctionCall) otto.Value {
	os := runtime.GOOS
	v, _ := otto.ToValue(os)
	return v

}

func version(otto.FunctionCall) otto.Value {
	s := runtime.Version()
	ver := strings.Split(s, "go")[1]
	v, _ := otto.ToValue(ver)
	return v
}
