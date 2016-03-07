package models

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/Felamande/lotdb/settings"
	_ "github.com/mattn/go-sqlite3" //init sqlite3
)

func TestForm(t *testing.T) {
	os.Chdir("D:/Dev/gopath/src/github.com/Felamande/lotdb")
	settings.Init("./settings/settings.toml")
	form := QueryForm{
		Sum: 123,
		Filters: []Filter{
			{
				Type:  "include",
				Value: 3,
			},
			{
				Type:  "exclude",
				Value: 5,
			},
		},
	}
	b, _ := json.Marshal(&form)
	fmt.Println(string(b))
}
