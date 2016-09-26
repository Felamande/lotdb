package settings

import (
	"fmt"
	"io/ioutil"
	"sync"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/Maxgis/tree"
	"github.com/kardianos/osext"
)

type staticCfg struct {
	RemoteRoot  string `toml:"remote"`
	VirtualRoot string `toml:"vstatic"`
	LocalRoot   string `toml:"lstatic"`
	CompressDef string `toml:"compress"`
}

type serverCfg struct {
	Port string `toml:"port"`
	Host string `toml:"host"`
}

type templateCfg struct {
	Home         string `toml:"home"`
	DelimesLeft  string `toml:"ldelime"`
	DelimesRight string `toml:"rdelime"`
	Charset      string `toml:"charset"`
	Reload       bool   `toml:"reload"`
}
type defaultVar struct {
	AppName string `toml:"appname"`
}

type adminCfg struct {
	Passwd string `toml:"passwd"`
}

type logCfg struct {
	Path   string `toml:"path"`
	Format string `toml:"format"`
}
type dbCfg struct {
	Type string
	Uri  string
}
type timeCfg struct {
	ZoneString string         `toml:"zone"`
	Location   *time.Location `toml:"-"`
}

type tlsCfg struct {
	Enable bool   `toml:"enable"`
	Cert   string `toml:"cert"`
	Key    string `toml:"key"`
}

type debugCfg struct {
	Port   int  `toml:"port"`
	Enable bool `toml:"enable"`
}

type setting struct {
	Static      staticCfg         `toml:"static"`
	Server      serverCfg         `toml:"server"`
	DB          dbCfg             `toml:"database"`
	Template    templateCfg       `toml:"template"`
	DefaultVars defaultVar        `toml:"defaultvars"`
	Admin       adminCfg          `toml:"admin"`
	Log         logCfg            `toml:"log"`
	Time        timeCfg           `toml:"time"`
	TLS         tlsCfg            `toml:"tls"`
	Debug       debugCfg          `toml:"debug"`
	Headers     map[string]string `toml:"headers"`
}

var (
	Folder string
	// IsInit = false

	//GlobalSettings
	Static      staticCfg
	Server      serverCfg
	Template    templateCfg
	DefaultVars defaultVar
	Admin       adminCfg
	Log         logCfg
	DB          dbCfg
	TLS         tlsCfg
	Time        timeCfg
	Debug       debugCfg
	Headers     map[string]string
)

func init() {
	var err error
	Folder, err = osext.ExecutableFolder()
	if err != nil {
		panic(err)
	}

}

var once sync.Once

func Init(cfgFile string) {
	once.Do(func() {
		settingStruct := setting{}
		b, err := ioutil.ReadFile(cfgFile)
		if err != nil {
			panic(err)
		}
		if err := toml.Unmarshal(b, &settingStruct); err != nil {
			panic(err)
		}
		settingStruct.Time.Location, err = time.LoadLocation(settingStruct.Time.ZoneString)
		if err != nil {
			fmt.Println(err)
			settingStruct.Time.Location = time.UTC
		}

		Static = settingStruct.Static
		Server = settingStruct.Server

		Template = settingStruct.Template
		DefaultVars = settingStruct.DefaultVars
		Admin = settingStruct.Admin
		Log = settingStruct.Log
		DB = settingStruct.DB
		Time = settingStruct.Time
		Headers = settingStruct.Headers
		Debug = settingStruct.Debug
		TLS = settingStruct.TLS
		tree.Print(settingStruct)
	})

}
