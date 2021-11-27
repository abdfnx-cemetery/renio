package tools

import "runtime"

// replace a string if it is windows
func isWindows(toChange *string, replacement string) {
	if runtime.GOOS == "windows" {
		*toChange = replacement
	}
}
