package constants

func TemplateSource() string {
	return `package main

import (
	"os"

	"github.com/abdfnx/renio/core" 
	"github.com/abdfnx/renio/config"
	"github.com/abdfnx/renio/core/options"
)

func main() {
	snap, _ := Asset("%s")
	toml, _ := Asset("%s")
	config, _ := config.ConfigParse(toml)
	env := options.Environment{
		NoColor: config.Options.NoColor,
		Args:    os.Args[1:],
	}

	opt := options.Options{
		SourceFile: "renio.js",
		Source:     string(snap),
		Perms:      &options.Perms{Fs: true},
		Env:        env,
	}

	core.Run(opt)
}`
}
