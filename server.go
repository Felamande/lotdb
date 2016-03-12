package main

import (
	"flag"
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

var conf = flag.String("-conf", "./settings/settings.js", "conf file, javascript.")

func init() {
	flag.Parse()
	settings.Init(*conf)
}

func main() {
	l := log.New(os.Stdout, "[tango]", log.Llevel|log.LstdFlags)
	l.SetLocation(settings.Location())
	t := tango.NewWithLog(l)

	t.Use(new(timemw.TimeHandler))
	t.Use(tango.Static(tango.StaticOptions{
		RootPath: settings.Get("static.localroot").String("./public"),
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
		Reload:      settings.Get("template.reload").Bool(false),
		Directory:   settings.Get("template.home").String("./templates"),
		Charset:     settings.Get("template.charset").String("UTF-8"),
		DelimsLeft:  settings.Get("template.delimes.Left").String("{%"),
		DelimsRight: settings.Get("template.delimes.right").String("%}"),
		Funcs:       utils.DefaultFuncs(),
	}))
	t.Use(events.Events())
	t.Use(header.CustomHeaders())

	t.Post(query.Url, new(query.QueryRouter))
	t.Get(page.HomeUrl, new(page.HomeRouter))
	host := settings.Get("server.port").String(":9000")
	go settings.MustWatch(*conf)
	if settings.Get("tls.use").Bool(false) {
		t.RunTLS(settings.Get("tls.cert").String(""), settings.Get("tls.key").String(""), host)
	} else {
		t.Run(host)
	}

}
