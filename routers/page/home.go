package page

import (
	"github.com/Felamande/lotdb/routers/base"
	"github.com/tango-contrib/renders"
)

const HomeUrl = "/"

type HomeRouter struct {
	base.BaseTplRouter
}

func (r *HomeRouter) Get() {
	if r.Data == nil {
		r.Data = make(renders.T)
	}
	r.Data["title"] = "query lottery db"
	r.Tpl = "home.html"

	r.Render(r.Tpl, r.Data)
}
