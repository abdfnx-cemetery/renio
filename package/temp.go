package package

import (
	"fmt"

	"github.com/abdfnx/renio/config"
	"github.com/abdfnx/renio/constants"
)

func GeneratePkgSource(path string) string {
	return fmt.Sprintf(constants.TemplateSource(), path, config.DefaultConfigPath)
}
