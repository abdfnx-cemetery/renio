package core

import (
	"encoding/json"
	"fmt"

	quickjs "github.com/abdfnx/qjs"
	"github.com/fatih/color"
)

// Build a new console formatter from console_util.go
var f = NewF()

// console.log bindings to quickjs engine
func ConsoleLog(ctx *quickjs.Context, value []quickjs.Value) quickjs.Value {
	data := value[2]
	// dataType is the JavaScript type of the data => `typeof arg`
	dataType := value[1].String()
	var result interface{}

	switch dataType {
    case "string":
      // Prints a string (without color)
      fmt.Println(data.String())

    case "function":
      // Prints String(myFunc)
      fmt.Fprintln(color.Output, color.New(color.FgCyan).SprintFunc()(data.String()))

    case "bigint":
      // Prints bigint corresponding to number
      fmt.Fprintln(color.Output, color.New(color.FgYellow).SprintFunc()(data.BigInt()))

    case "number":
      // Prints a number
      fmt.Fprintln(color.Output, color.New(color.FgYellow).SprintFunc()(data.Int32()))

    default:
      json.Unmarshal([]byte(data.String()), &result)
      prty, _ := f.Marshal(result)
      fmt.Fprintln(color.Output, string(prty))
	}

	return ctx.Null()
}
