package toolate

import (
	"github.com/Felamande/lotdb/models"
	"github.com/Felamande/lotdb/routers/base"
	"github.com/tango-contrib/binding"
)

const Url = "/toolate"

type TooLateRouter struct {
	base.BaseJSONRouter
	binding.Binder

	Form models.TooLateForm
}

func (r *TooLateRouter) Post() {
	if err := r.Bind(&r.Form); err != nil {
		r.Logger.Error(err)
		return
	}

	if r.Form.Clicked {
		r.Logger.Info("clicked true.")
	}
}
