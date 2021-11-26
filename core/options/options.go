package options

import (
	"github.com/abdfnx/qjs"
)

type Recv func(id quickjs.Value, val quickjs.Value) quickjs.Value

type Renio struct {
  // Renio -> represents general data for the runtime
	// Permissions
	Perms *Perms
	// Async recv function
	Recv Recv
}

type Environment struct {
  // Environment -> configure the runtime environment
	// Enable or disable color logging
	NoColor bool
	// Command-line args to pass into Renio.args
	Args []string
	// Whether to run tests associated with `Renio.tests()`
	RunTests bool
}

type Perms struct {
  // Perms -> permissions available for Renio
	// File system access
	Fs bool
	// Net access
	Net bool
	// Env access
	Env bool
}

type Options struct {
  // Options -> options for dispatching a new Renio + QuickJS runtime
	// File name of the source (used for debuging purposes)
	SourceFile string
	// Source code
	Source string
	// Permission
	Perms *Perms
	// Configure Environment
	Env Environment
}
