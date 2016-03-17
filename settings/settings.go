package settings

import (
	"strings"
	"time"

	"github.com/Felamande/lotdb/settings/jsvm"
	"github.com/Felamande/otto"

	//module
	_ "github.com/Felamande/lotdb/settings/jsvm/module/gofunc"
	_ "github.com/Felamande/lotdb/settings/jsvm/module/math"
	_ "github.com/Felamande/lotdb/settings/jsvm/module/os"
	_ "github.com/Felamande/lotdb/settings/jsvm/module/path"
	_ "github.com/Felamande/lotdb/settings/jsvm/module/str"
	_ "github.com/robertkrimen/otto/underscore"

	"gopkg.in/fsnotify.v1"
)

var lastReloadTime = time.Now()

type result struct {
	value otto.Value
	err   error
}

func Init(file string) {
	jsvm.Run(file)
}

func get(val string) *result {
	value, err := jsvm.Vm().Get(val)
	return &result{value, err}
}

func (r *result) get(val string) *result {
	if !r.value.IsObject() {
		return r
	}

	value, err := r.value.Object().Get(val)

	return &result{value, err}
}

func Get(path string) *result {

	vals := strings.Split(path, ".")
	r := get(vals[0])

	for _, val := range vals[1:] {
		r = r.get(val)
	}
	return r
}

func (r *result) String(Default string) string {
	if r.err != nil || r.value.IsUndefined() || !r.value.IsString() {
		return Default
	}

	return r.value.String()
}

func (r *result) Bool(Default bool) bool {
	if r.err != nil || r.value.IsUndefined() || !r.value.IsBoolean() {
		return Default
	}

	b, err := r.value.ToBoolean()

	if err != nil {
		return Default
	}
	return b
}

func (r *result) Int(Default int) int {
	if r.err != nil || r.value.IsUndefined() || !r.value.IsNumber() {
		return Default
	}

	i, err := r.value.ToInteger()
	if err != nil {
		return Default
	}

	return int(i)
}

func (r *result) Float(Default float64) float64 {
	if r.err != nil || r.value.IsUndefined() || !r.value.IsNumber() {
		return Default
	}

	f, err := r.value.ToFloat()
	if err != nil {
		return Default
	}

	return f
}

func (r *result) ForEach(foreach func(k string, v interface{}, err error)) {
	if r.err != nil || !r.value.IsObject() {
		return
	}

	o := r.value.Object()
	o.ForEach(func(key string) {
		value, err := o.Get(key)
		if err != nil {
			foreach(key, value, err)
			return
		}
		val, err := value.Export()
		foreach(key, val, err)

	})
}

func Location() *time.Location {
	loc := Get("time.zone").String("Hongkong")
	loca, err := time.LoadLocation(loc)
	if err != nil {
		loca = time.UTC
	}
	return loca
}

func MustWatch(file string) {
	if err := watch(file); err != nil {
		panic(err)
	}
}

func Watch(file string) error {
	return watch(file)
}

func watch(file string) error {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	err = w.Add(file)
	if err != nil {
		return err
	}
	for {
		select {
		case e := <-w.Events:
			if e.Op == fsnotify.Remove || e.Op == fsnotify.Rename || e.Op == fsnotify.Chmod {
				continue
			}
			if time.Now().Sub(lastReloadTime) < 3 {
				continue
			}
			reload(file)
			lastReloadTime = time.Now()
		}
	}
}

func reload(file string) error {
	err := jsvm.Run(file)
	if err != nil {
		return err
	}
	return nil
}
