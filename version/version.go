package version

import (
	_ "embed"
	"fmt"
	"github.com/Slava02/Involvio/config"
)

//go:embed VERSION
var VERSION string

func PrintVersion(cfg *config.Config) {
	println(fmt.Sprintf("%s version %s", cfg.App.Name, VERSION))
}
