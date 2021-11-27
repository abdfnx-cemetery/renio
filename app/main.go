package main

import (
	"fmt"
	"runtime"
	"errors"
	"os"
	
	"github.com/abdfnx/renio/context"
	"github.com/abdfnx/renio/contract"
	"github.com/abdfnx/renio/cmd/renio"
	"github.com/abdfnx/renio/cmd/factory"
	"github.com/abdfnx/renio/core"
	"github.com/abdfnx/renio/tools"
	"github.com/abdfnx/renio/shared"

	surveyCore "github.com/AlecAivazis/survey/v2/core"
	"github.com/AlecAivazis/survey/v2/terminal"

	"github.com/mgutz/ansi"
	"github.com/spf13/cobra"
)

type exitCode int

const (
	exitOK     exitCode = 0
	exitError  exitCode = 1
	exitCancel exitCode = 2
)

func main() {
	code := mainRun()
	os.Exit(int(code))
}

func mainRun() exitCode {
	runtime.LockOSThread()

	cmdFactory := factory.New()
	hasDebug := os.Getenv("DEBUG") != ""
	stderr := cmdFactory.IOStreams.ErrOut

	if !cmdFactory.IOStreams.ColorEnabled() {
		surveyCore.DisableColor = true
	} else {
		surveyCore.TemplateFuncsWithColor["color"] = func(style string) string {
			switch style {
			case "white":
				if cmdFactory.IOStreams.ColorSupport256() {
					return fmt.Sprintf("\x1b[%d;5;%dm", 38, 242)
				}

				return ansi.ColorCode("default")

			default:
				return ansi.ColorCode(style)
			}
		}
	}

	if len(os.Args) > 1 && os.Args[1] != "" {
		cobra.MousetrapHelpText = ""
	}

	RootCmd := renio.Execute(renio.Renio{
		Run:    core.Run,
		Bundle: contract.BundleModule,
		Dev:    context.RunDev,
	}, cmdFactory)

	if cmd, err := RootCmd.ExecuteC(); err != nil {
		if err == tools.SilentError {
			return exitError
		} else if tools.IsUserCancellation(err) {
			if errors.Is(err, terminal.InterruptErr) {
				fmt.Fprint(stderr, "\n")
			}

			return exitCancel
		}

		shared.PrintError(stderr, err, cmd, hasDebug)

		return exitError
	}

	if renio.HasFailed() {
		return exitError
	}

	return exitOK
}
