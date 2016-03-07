package query

import (
	"github.com/Felamande/lotdb/models"
	"github.com/Felamande/lotdb/routers/base"
	"github.com/tango-contrib/binding"
)

const Url = "/query"

type QueryRouter struct {
	base.BaseJSONRouter
	binding.Binder

	Form   models.QueryForm
	Errors binding.Errors
}

func (r *QueryRouter) Before() {
	r.BaseJSONRouter.Before()
	r.Form = models.QueryForm{}
	if e := r.Bind(&r.Form); len(e) != 0 {
		r.Errors = e
	}
	if Formatter, ok := interface{}(r.Form).(models.Formatter); ok {
		Formatter.Format()
	}
}

func (r *QueryRouter) Post() interface{} {
	if _, exist := r.JSON["err"]; exist {
		return r.JSON
	}
	r.Logger.Info(r.Form, r.Header().Get("Content-Type"))

	if r.Errors.Len() != 0 {
		r.JSON["err"] = r.Errors[0].Error()
		r.JSON["sucess"] = false
		r.Logger.Error(r.Errors.ErrorMap())
		return r.JSON
	}

	re, err := models.GetQueryResults(r.Form)
	if err != nil {
		if dberr, ok := err.(models.DatabaseError); ok {
			r.JSON["err"] = "database error."
			r.JSON["sucess"] = false
			go r.Logger.Error(dberr)
			return r.JSON
		}
		r.JSON["err"] = err.Error()
		r.JSON["sucess"] = false
		return r.JSON
	}

	r.JSON["re"] = re
	r.JSON["success"] = true

	return r.JSON
}
