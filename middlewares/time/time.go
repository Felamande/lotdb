package time

import (
	"regexp"
	"time"

	"github.com/lunny/tango"
)

var uaRegexp = regexp.MustCompile(`(Android|Windows|iPhone|iPad|MicroMessenger|MSIE|Chrome|QQ)`)

type TimeHandler struct {
}

func (h *TimeHandler) Handle(ctx *tango.Context) {
	t1 := time.Now()
	// ctx.Logger.Info(ctx.Header())
	ctx.Next()
	ctx.Logger.Infof("Completed %v %v %v in %v for %v %v",
		ctx.Req().Method,
		ctx.Req().URL.Path,
		ctx.Status(),
		time.Since(t1),
		ctx.Req().RemoteAddr,
		uaRegexp.FindAllString(ctx.Req().Header.Get("User-Agent"), -1),
	)
}
