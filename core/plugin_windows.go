// +build windows

package core

import (
	"os"

	"github.com/abdfnx/renio/tools"
)

// OpenPlugin Go plugins have not yet been implemented on windows
func OpenPlugin(path string, arg interface{}) interface{} {
	tools.LogError("Go Plugins're not supported for windows.", "See https://github.com/golang/go/issues/19282 `plugin: add Windows support` issue...")
	os.Exit(1)

	return nil
}
