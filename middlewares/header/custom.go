package header

import (
	"github.com/Felamande/lotdb/settings"
	"github.com/lunny/tango"
)

func CustomHeaders() tango.HandlerFunc {
	return func(ctx *tango.Context) {
		ctx.Next()
		settings.Get("headers").ForEach(func(key string, value interface{}, err error) {
			if err != nil {
				ctx.Logger.Warn("set header:", err)
				return
			}
			if v, ok := value.(string); ok {
				ctx.Header().Set(key, v)
			}

		})

	}
}
