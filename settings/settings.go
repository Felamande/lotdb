package settings

import (
	"io/ioutil"
	"strings"
	"sync"
	"time"

	"github.com/robertkrimen/otto"
	_ "github.com/robertkrimen/otto/underscore" //register underscore
	"gopkg.in/fsnotify.v1"
)

var vm *otto.Otto
var confFile string

var lock = new(sync.Mutex)

type result struct {
	value otto.Value
	err   error
}

func Init(file string) {

	// InitOnce.Do(func() {

	vm = otto.New()

	b, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}

	_, err = vm.Run(string(b))
	if err != nil {
		panic(err)
	}
	confFile = file
	// })

}

func get(val string) *result {
	value, err := vm.Get(val)
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
	lock.Lock()
	defer lock.Unlock()

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

func (r *result) ForEach(foreach func(key string, val otto.Value)) {
	if r.err != nil || !r.value.IsObject() {
		return
	}

	o := r.value.Object()
	keys := o.Keys()
	for _, key := range keys {
		value, err := o.Get(key)
		if err != nil {
			continue
		}

		foreach(key, value)

	}
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
			reload(file)

		}
	}
}

func reload(file string) error {
	lock.Lock()
	defer lock.Unlock()
	b, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	_, err = vm.Run(string(b))
	if err != nil {
		return err
	}
	return nil
}
