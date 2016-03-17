package os

import (
	"bytes"
	"errors"
	"io"
	"os"
	// "errors"

	"io/ioutil"
	"os/exec"
	"strings"

	"github.com/Felamande/lotdb/settings/jsvm"
	"github.com/Felamande/otto"
)

func init() {
	p := jsvm.Module("os")
	p.Extend("readFile", readFile)
	p.Extend("readFileAsync", readFileAsync)
	p.Extend("system", system)
	p.Extend("getwd", getwd)
	p.Extend("output", sysOutput)
	p.Extend("writeFile", writeFile)
}

func readFile(call otto.FunctionCall) otto.Value {

	file, err := call.Argument(0).ToString()
	if err != nil {
		e, _ := otto.ToValue(readError{err.Error()})
		return e
	}
	b, err := ioutil.ReadFile(file)
	if err != nil {
		e, _ := otto.ToValue(readError{err.Error()})
		return e
	}
	v, _ := otto.ToValue(string(b))
	return v

}

func readFileAsync(call otto.FunctionCall) otto.Value {
	fileArg := call.Argument(0)
	contentCb := call.Argument(1)
	errCb := call.Argument(2)

	file := fileArg.String()
	if !fileArg.IsString() {
		return callback(errCb, "invalid file name "+file)
	}

	f, err := os.Open(file)
	if err != nil {
		return callback(errCb, err.Error())
	}

	go func() {
		defer f.Close()
		buf := bytes.NewBuffer([]byte(""))
		io.Copy(buf, f)
		callback(contentCb, buf.String())
	}()
	return otto.UndefinedValue()

}

func getwd(call otto.FunctionCall) otto.Value {
	dir, err := os.Getwd()
	if err != nil {
		return otto.UndefinedValue()
	}
	v, _ := otto.ToValue(dir)

	return v
}

func system(call otto.FunctionCall) otto.Value {
	var cmdList []string

	arg0 := call.Argument(0)
	iarg0, err := arg0.Export()
	if err != nil {
		return otto.UndefinedValue()
	}

	switch arg0t := iarg0.(type) {
	case []string:
		cmdList = arg0t
	case string:
		cmdList = strings.Split(arg0t, " ")
	default:
		return otto.UndefinedValue()
	}

	cmdLen := len(cmdList)
	if cmdLen == 0 {
		panic("no cmd")
	}
	c := exec.Command(cmdList[0], cmdList[1:]...)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	err = c.Run()
	if err != nil {
		v, _ := otto.ToValue(err)
		return v
	}
	return otto.TrueValue()
}

func sysOutput(call otto.FunctionCall) otto.Value {
	var cmdList []string

	arg0 := call.Argument(0)
	iarg0, err := arg0.Export()
	if err != nil {
		v, _ := otto.ToValue("")
		return v
	}

	switch arg0t := iarg0.(type) {
	case []string:
		cmdList = arg0t
	case string:
		cmdList = strings.Split(arg0t, " ")
	default:
		v, _ := otto.ToValue("")
		return v
	}

	cmdLen := len(cmdList)
	if cmdLen == 0 {
		panic("no cmd")
	}
	b := new(bytes.Buffer)
	c := exec.Command(cmdList[0], cmdList[1:]...)
	c.Stdout = b
	c.Stderr = b
	err = c.Run()
	if err != nil {
		v, _ := otto.ToValue("")
		return v
	}
	v, _ := otto.ToValue(b.String())
	return v
}

//writeFile function writeFile(file,content,flag,errCb,formatter)
//flag a:append, c:create, t:truncate, l:newline,s:sync
//formatter function, format content to custom string.
//errCb error callback
func writeFile(call otto.FunctionCall) otto.Value {
	fileArg := call.Argument(0)
	content := call.Argument(1)
	flagArg := call.Argument(2)
	errCb := call.Argument(3) //error callback
	formatterCb := call.Argument(4)
	// callbackArg := call.Argument(2)
	if !fileArg.IsString() {
		return callback(errCb, "invalid file name "+fileArg.String())
	}
	if content.IsUndefined() {
		return callback(errCb, "content is not provided.")
	}

	var flagStr string
	if flagArg.IsString() {
		flagStr = flagArg.String()
	}
	var flag = os.O_WRONLY
	var newline string
	for _, m := range flagStr {
		switch m {
		case 'a':
			flag |= os.O_APPEND
		case 'c':
			flag |= os.O_CREATE
		case 't':
			flag |= os.O_TRUNC
		case 's':
			flag |= os.O_SYNC
		case 'l':
			newline = "\n"
		}
	}
	f, err := os.OpenFile(fileArg.String(), flag, 0777)
	if err != nil {
		return callback(errCb, err.Error())
	}
	defer f.Close()
	formatted, err := cbGetValue(formatterCb, content)
	if err != nil {
		return callback(errCb, err.Error())
	}
	bytes := []byte(formatted + newline)

	_, err = f.Write(bytes)
	if err != nil {
		return callback(errCb, err.Error())
	}
	return otto.UndefinedValue()
}

func stringValue(s string) otto.Value {
	v, _ := otto.ToValue(s)
	return v
}

func errorValue(err error) otto.Value {
	value, _ := otto.ToValue(err)
	return value
}

func callback(cb otto.Value, arg interface{}) otto.Value {
	if cb.Class() != "Function" {
		return otto.UndefinedValue()
	}
	cb.Call(cb, arg)
	return otto.UndefinedValue()

}
func cbGetValue(cb otto.Value, arg otto.Value) (string, error) {
	if cb.IsUndefined() {
		return arg.String(), nil
	}
	if !cb.IsFunction() {
		return "", errors.New("invalid formatter")
	}
	v, err := cb.Call(cb, arg)
	if err != nil {
		return "", err
	}
	return v.String(), nil
}
