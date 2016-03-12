package settings

import (
	"os"
	"testing"
)

func TestJs(t *testing.T) {
	os.Chdir(`D:\Dev\gopath\src\github.com\Felamande\lotdb`)
	Init("./settings/settings.js")
	s := Get("template.delime.left").String("fail")
	t.Error("=", s)
}
