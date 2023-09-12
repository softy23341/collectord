// Licensed under the LGPLv3, see LICENCE file for details.

package erry

import (
	"reflect"
	"runtime"
	"strings"
)

// goPath is the deduced path based on the location of this file as compiled.
var goPath string

type emptyStub struct{}

func init() {
	if _, file, _, ok := runtime.Caller(0); ok {
		pkgPath := reflect.TypeOf(emptyStub{}).PkgPath()
		if idx := strings.LastIndex(file, pkgPath); idx > 0 {
			goPath = file[:idx]
		}
	}
}

func trimGoPath(filename string) string {
	return strings.TrimPrefix(filename, goPath)
}
