package tools

import (
	"errors"
	"fmt"
	"os"

	quickjs "github.com/abdfnx/qjs"
)

func Check(err error) {
	if err != nil {
		var evalErr *quickjs.Error
		if errors.As(err, &evalErr) {
			fmt.Println(evalErr.Cause)
			fmt.Println(evalErr.Stack)
		}

		Panic(err)
	}
}

// Panic pretty print the error and exit with status code 1
func Panic(err error) {
	LogError("Error", fmt.Sprintf("%v", err))
	os.Exit(1)
}
