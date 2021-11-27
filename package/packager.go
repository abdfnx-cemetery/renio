package _package

import (
	"os"
	"path/filepath"

	"github.com/abdfnx/renio/config"
	"github.com/abdfnx/renio/tools"

	"github.com/go-bindata/go-bindata"
)

// PkgSource pack bundled js source into an executable
func PkgSource(source string) {
	c := bindata.NewConfig()

	input := parseInput(source)

	if config.ConfigExists() {
		config := parseInput(config.DefaultConfigPath)
		c.Input = []bindata.InputConfig{input, config}
	} else {
		c.Input = []bindata.InputConfig{input}
	}

	c.Output = "target/renio-package/asset.go"

	err := bindata.Translate(c)
	tools.Check(err)

	entry := GeneratePkgSource(source)
	f, _ := os.Create("target/renio-package/main.go")

	defer f.Close()
	_, err = f.WriteString(entry)

	tools.Check(err)

	ExecBuild("target/renio-package")
}

func parseInput(path string) bindata.InputConfig {
	return bindata.InputConfig{
		Path:      filepath.Clean(path),
		Recursive: false,
	}
}
