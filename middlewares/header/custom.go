package header

import (
	"github.com/Felamande/lotdb/settings"
	"github.com/lunny/tango"
)

func CustomHeaders() tango.HandlerFunc {
	return func(ctx *tango.Context) {
		for key, value := range settings.Headers {
			if key == "" || value == "" {
				continue
			}
			ctx.Header().Set(key, value)
		}
		ctx.Next()

	}
}
