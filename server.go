package main

import (
	"github.com/Felamande/lotdb/modules/utils"
	"github.com/Felamande/lotdb/settings"
	"github.com/lunny/tango"

	//middlewares
	"github.com/Felamande/lotdb/middlewares/time"
	"github.com/tango-contrib/binding"
	"github.com/tango-contrib/events"
	"github.com/tango-contrib/renders"
	//routers
	"github.com/Felamande/lotdb/routers/page"
	"github.com/Felamande/lotdb/routers/query"
)

func init() {
	settings.Init("./settings/settings.toml")
}

func main() {

	t := tango.New()

	t.Use(new(time.TimeHandler))
	t.Use(tango.Static(tango.StaticOptions{
		RootPath: settings.Static.LocalRoot,
	}))
	t.Use(binding.Bind())
	t.Use(tango.ClassicHandlers...)
	t.Use(renders.New(renders.Options{
		Reload:      settings.Template.Reload,
		Directory:   settings.Template.Home,
		Charset:     settings.Template.Charset,
		DelimsLeft:  settings.Template.DelimesLeft,
		DelimsRight: settings.Template.DelimesRight,
		Funcs:       utils.DefaultFuncs(),
	}))
	t.Use(events.Events())

	t.Post(query.Url, new(query.QueryRouter))
	t.Get(page.HomeUrl, new(page.HomeRouter))

	t.Run(settings.Server.Host)
}
