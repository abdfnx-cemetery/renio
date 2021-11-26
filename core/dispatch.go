package core

import (
	"os"

	"github.com/abdfnx/qjs"
	"github.com/abdfnx/renio/core/ops"
	"github.com/abdfnx/renio/core/options"
	"github.com/abdfnx/renio/tools"

	"github.com/spf13/afero"
)

/*
  RenioSendNameSpace Native function corresponding to the JavaScript global `__send`
  It is binded with `__send` and accepts arguments including op ID
*/
func RenioSendNameSpace(renio *options.Renio) func(ctx *quickjs.Context, this quickjs.Value, args []quickjs.Value) quickjs.Value {
	// Create a new file system driver
	var fs = ops.FsDriver{
		// NOTE: afero can also be used to create in-memory file system
		Fs:    afero.NewOsFs(),
		Perms: renio.Perms,
	}

	// The returned function handles the op and execute corresponding native code
	return func(ctx *quickjs.Context, this quickjs.Value, args []quickjs.Value) quickjs.Value {
		switch args[0].Int32() {
		case FSRead:
			FileSystemChecker(renio.Perms)
			file := args[1]
			val := fs.ReadFile(ctx, file)
			return val

		case FSExists:
			FileSystemChecker(renio.Perms)
			file := args[1]
			val := fs.Exists(ctx, file)
			return val

		case FSWrite:
			FileSystemChecker(renio.Perms)
			file := args[1]
			contents := args[2]
			val := fs.WriteFile(ctx, file, contents)
			return val

		case FSCwd:
			FileSystemChecker(renio.Perms)
			val := fs.Cwd(ctx)
			return val

		case FSStat:
			FileSystemChecker(renio.Perms)
			file := args[1]
			val := fs.Stat(ctx, file)
			return val

		case FSRemove:
			FileSystemChecker(renio.Perms)
			file := args[1]
			val := fs.Remove(ctx, file)
			return val

		case Log:
			return ConsoleLog(ctx, args)

		case Fetch:
			NetChecker(renio.Perms)
			one := args[1]
			url := args[2]
			body := ops.Fetch(ctx, url)
			obj := ctx.Object()
			defer obj.Free()
			obj.Set("ok", body)
			renio.Recv(one, obj)
			return ctx.Null()

		case Serve:
			id := args[1]
			url := args[2]
			cb := func(res quickjs.Value) string {
				obj := ctx.Object()
				defer obj.Free()
				obj.Set("ok", res)
				rtrn := renio.Recv(id, res)
				return rtrn.String()
			}

			ops.Serve(ctx, cb, id, url)
			return ctx.Null()

		case FSMkdir:
			FileSystemChecker(renio.Perms)
			file := args[1]
			val := fs.Mkdir(ctx, file)
			return val

		case Env:
			EnvChecker(renio.Perms)
			val := ops.Env(ctx, args)
			return val

		case FSWalk:
			FileSystemChecker(renio.Perms)
			file := args[1]
			val := fs.Walk(ctx, file)
			return val

		default:
			return ctx.Null()
		}
	}
}

// FileSystemChecker utility to check whether file system access is avaliable or not
func FileSystemChecker(perms *options.Perms) {
	if !perms.Fs {
		tools.LogError("Perms Error: ", "Filesystem access is blocked.")
		os.Exit(1)
	}
}

// NetChecker utility to check whether net access is avaliable or not
func NetChecker(perms *options.Perms) {
	if !perms.Net {
		tools.LogError("Perms Error: ", "Net is blocked.")
		os.Exit(1)
	}
}

func EnvChecker(perms *options.Perms) {
	if !perms.Env {
		tools.LogError("Perms Error: ", "Environment Variables is blocked.")
		os.Exit(1)
	}
}

// RenioRecvNameSpace Native function corresponding to the JavaScript global `__recv`
// It is binded with `__recv` and accepts arguments including recv ID of the async function
func RenioRecvNameSpace(renio *options.Renio) func(ctx *quickjs.Context, this quickjs.Value, args []quickjs.Value) quickjs.Value {
	// the returned function handles the __recv behaviour
	// It is capable of calling the callback for a particular async op after it has finished
	return func(ctx *quickjs.Context, this quickjs.Value, args []quickjs.Value) quickjs.Value {
		fn := args[0]

		if renio.Recv != nil {
			ctx.ThrowError(fmt.Errorf("recv cannot be called more than once"))
			return ctx.Null()
		}

		renio.Recv = func(id quickjs.Value, val quickjs.Value) quickjs.Value {
			result := fn.Call(id, val)
			// defer result.Free()
			return result
		}

		return ctx.Null()
	}
}
