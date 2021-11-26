package context

import (
	"fmt"

	"github.com/abdfnx/qjs"
)

// ReportDiagnostics report TypeScript diagnostics
func ReportDiagnostics(diagnostics quickjs.Value) {
	diag := diagnostics.String()
	if diag != "" {
		fmt.Println(diagnostics.String())
	}
}
