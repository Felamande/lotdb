package debug

import (
	"net/http"
	_ "net/http/pprof"
	"strconv"
)

func On(port int) {
	http.ListenAndServe(":"+strconv.Itoa(port), nil)
}
