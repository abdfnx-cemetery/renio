package tools

import (
	"runtime"

	"github.com/abdfnx/shell"
)

func CheckDotRenioDir() {
	if runtime.GOOS == "windows" {
		shell.PWSLCmd(`
			if (!(Test-Path -path $HOME/.renio)) {
				mkdir $HOME/.renio
			}
		`)
	} else {
		shell.ShellCmd(`
			if [ ! -d $HOME/.renio ]; then
				mkdir $HOME/.renio
			fi
		`)
	}
}
