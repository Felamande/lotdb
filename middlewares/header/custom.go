package header

import (
	"github.com/Felamande/lotdb/settings"
	"github.com/lunny/tango"
)

func CustomHeaders() tango.HandlerFunc {
	return func(ctx *tango.Context) {
		settings.Get("headers").ForEach(func(key string, value interface{}) {

			ctx.Header().Set(key, value.(string))
		})

		ctx.Next()

	}
}
