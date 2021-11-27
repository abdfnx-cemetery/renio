package core

import (
	"io"

	"github.com/abdfnx/renio/core/options"
	"github.com/abdfnx/renio/tools"

	"github.com/abdfnx/qjs"
)

// PrepareRuntimeContext prepare the runtime and context with Renio's internal ops
// injects `__send` and `__recv` global dispatch functions into runtime
func PrepareRuntimeContext(cxt *quickjs.Context, jsruntime quickjs.Runtime, args []string, flags *options.Perms, mode string) {
	// Assign perms
	renio := &options.Renio{Perms: flags}

	globals := cxt.Globals()
	// Attach send & recv global ops
	globals.SetFunction("__send", RenioSendNameSpace(renio))
	globals.SetFunction("__recv", RenioRecvNameSpace(renio))

	// Prepare runtime context with namespace and client op code
	// The snapshot is generated at bootstrap process
	snap, _ := Asset("target/renio.js")
	k, err := cxt.Eval(string(snap))
	tools.Check(err)
	defer k.Free()

	ns := globals.Get("Renio")

	defer ns.Free()
	// Assign `Renio.args` with the os args
	_Args := cxt.Array()

	for i, arg := range args {
		_Arg := cxt.String(arg)
		_Args.SetByUint32(uint32(i), _Arg)
	}

	ns.Set("args", _Args)
	// Assing `Renio.mode` with current environment mode
	_Mode := cxt.String(mode)
	ns.Set("mode", _Mode)

	// Runtime check to execute async jobs
	for {
		_, err = jsruntime.ExecutePendingJob()
		if err == io.EOF {
			err = nil
			break
		}

		tools.Check(err)
	}
}

// Run create and dispatch a QuickJS runtime binded with Renio's OPs configurable using options
func Run(opt options.Options) {
	// Create a new quickJS runtime
	jsruntime := quickjs.NewRuntime()
	defer jsruntime.Free()

	// Create a new runtime context
	cxt := jsruntime.NewContext()
	defer cxt.Free()

	// mode is not configurable directly and is to be determined based on RunTests
	// defaults to `run`
	mode := "run"

	if opt.Env.RunTests {
		mode = "test"
	}

	// Prepare runtime and context with Renio namespace
	PrepareRuntimeContext(cxt, jsruntime, opt.Env.Args, opt.Perms, mode)

	// Evaluate the source
	result, err := cxt.EvalFile(opt.Source, opt.SourceFile)
	tools.Check(err)
	defer result.Free()

	// Check for exceptions
	if result.IsException() {
		err = cxt.Exception()
		tools.Check(err)
	}

	for {
		_, err = jsruntime.ExecutePendingJob()
		if err == io.EOF {
			err = nil
			break
		}

		tools.Check(err)
	}
}
