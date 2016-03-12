package header

import (
	"github.com/Felamande/lotdb/settings"
	"github.com/lunny/tango"
	"github.com/robertkrimen/otto"
)

func CustomHeaders() tango.HandlerFunc {
	return func(ctx *tango.Context) {
		settings.Get("headers").ForEach(func(key string, value otto.Value) {
			if !value.IsString() {
				return
			}
			ctx.Header().Set(key, value.String())
		})

		ctx.Next()

	}
}
