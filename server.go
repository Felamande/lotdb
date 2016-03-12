package main

import (
	"os"

	"github.com/Felamande/lotdb/settings"
	"github.com/lunny/tango"

	//modules
	"github.com/Felamande/lotdb/modules/log"
	"github.com/Felamande/lotdb/modules/utils"

	//middlewares
	"github.com/Felamande/lotdb/middlewares/header"
	timemw "github.com/Felamande/lotdb/middlewares/time"
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
	l := log.New(os.Stdout, "[tango]", log.Llevel|log.LstdFlags)
	l.SetLocation(settings.Time.Location)
	t := tango.NewWithLog(l)

	t.Use(new(timemw.TimeHandler))
	t.Use(tango.Static(tango.StaticOptions{
		RootPath: settings.Static.LocalRoot,
	}))
	t.Use(binding.Bind())
	t.Use(
		tango.Recovery(false),
		tango.Compresses([]string{}),
		tango.Static(tango.StaticOptions{Prefix: "public"}),
		tango.Return(),
		tango.Param(),
		tango.Contexts(),
	)
	t.Use(renders.New(renders.Options{
		Reload:      settings.Template.Reload,
		Directory:   settings.Template.Home,
		Charset:     settings.Template.Charset,
		DelimsLeft:  settings.Template.DelimesLeft,
		DelimsRight: settings.Template.DelimesRight,
		Funcs:       utils.DefaultFuncs(),
	}))
	t.Use(events.Events())
	t.Use(header.CustomHeaders())

	t.Post(query.Url, new(query.QueryRouter))
	t.Get(page.HomeUrl, new(page.HomeRouter))
	if settings.Tls.Use {
		t.RunTLS(settings.Tls.Cert, settings.Tls.Key, settings.Server.Host)
	} else {
		t.Run(settings.Server.Host)
	}

}
