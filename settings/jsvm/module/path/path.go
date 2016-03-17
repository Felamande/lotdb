package path

import (

	// "errors"

	"path/filepath"

	"github.com/Felamande/lotdb/settings/jsvm"
	"github.com/Felamande/otto"
)

func init() {
	p := jsvm.Module("path")
	p.Extend("join", join)
}

func join(call otto.FunctionCall) otto.Value {
	var pathList []string

	for _, arg := range call.ArgumentList {
		v, _ := arg.Export()
		switch val := v.(type) {
		case []string:
			pathList = append(pathList, val...)
		case string:
			pathList = append(pathList, val)
		}
	}
	path := filepath.Join(pathList...)
	v, _ := otto.ToValue(path)
	return v
}
