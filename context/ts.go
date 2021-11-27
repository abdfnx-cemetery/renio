package context

import (
	"fmt"
	"time"
	"runtime"

	"github.com/abdfnx/renio/core"
	"github.com/abdfnx/renio/core/options"
	"github.com/abdfnx/renio/tools"

	"github.com/abdfnx/qjs"
	"github.com/abdfnx/shell"

	"github.com/briandowns/spinner"
)

func Compile(source string, sourceFile string, fn func(val quickjs.Value), flags *options.Perms, args []string) {
	data, err := core.Asset("typescript/typescript.js")
	if err != nil {
		panic("Asset was not found.")

		// clone abdfnx/renio_typescript
		cmd := "git clone https://github.com/abdfnx/renio_typescript typescrip"
		s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
		s.Suffix = " 📦 Cloning..."
		s.Start()

		shell.Run(cmd)

		s.Stop()
	}

	runtime.LockOSThread()
	jsruntime := quickjs.NewRuntime()
	defer jsruntime.Free()

	context := jsruntime.NewContext()
	defer context.Free()

	core.PrepareRuntimeContext(context, jsruntime, args, flags, "dev")

	globals := context.Globals()
	report := func(ctx *quickjs.Context, this quickjs.Value, args []quickjs.Value) quickjs.Value {
		fn(args[0])
		return ctx.Null()
	}

	d := func(ctx *quickjs.Context, this quickjs.Value, args []quickjs.Value) quickjs.Value {
		asset, er := core.Asset(args[0].String())
		if er != nil {
			panic("Asset was not found.")
		}

		return ctx.String(string(asset))
	}

	globals.Set("Report", context.Function(report))
	globals.Set("Asset", context.Function(d))
	bundle := string(data) + jsCheck(source, sourceFile)
	result, err := context.Eval(bundle)
	tools.Check(err)

	defer result.Free()
}

func jsCheck(source, sourceFile string) string {
	return fmt.Sprintf("typeCheck(`%s`, `%s`);", sourceFile, source)
}
