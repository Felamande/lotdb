package settings

import (
	"io/ioutil"
	"sync"

	"github.com/BurntSushi/toml"
	"github.com/kardianos/osext"
)

type staticCfg struct {
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

type setting struct {
	Static      staticCfg   `toml:"static"`
	Server      serverCfg   `toml:"server"`
	DB          dbCfg       `toml:"database"`
	Template    templateCfg `toml:"template"`
	DefaultVars defaultVar  `toml:"defaultvars"`
	Admin       adminCfg    `toml:"admin"`
	Log         logCfg      `toml:"log"`
}

var (
	Folder        string
	settingStruct = new(setting)
	IsInit        = false

	//GlobalSettings
	Static      staticCfg
	Server      serverCfg
	Template    templateCfg
	DefaultVars defaultVar
	Admin       adminCfg
	Log         logCfg
	DB          dbCfg
)

var lock = new(sync.Mutex)
var InitOnce = new(sync.Once)

func init() {
	var err error
	Folder, err = osext.ExecutableFolder()
	if err != nil {
		panic(err)
	}

}

func Init(cfgFile string) {
	InitOnce.Do(func() {
		b, err := ioutil.ReadFile(cfgFile)
		if err != nil {
			panic(err)
		}
		toml.Unmarshal(b, settingStruct)

		Static = settingStruct.Static
		Server = settingStruct.Server

		Template = settingStruct.Template
		DefaultVars = settingStruct.DefaultVars
		Admin = settingStruct.Admin
		Log = settingStruct.Log
		DB = settingStruct.DB

	})

	IsInit = true
}
