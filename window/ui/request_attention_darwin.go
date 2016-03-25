// +build darwin, !js

package ui

/*
#cgo CFLAGS: -x objective-c -fobjc-arc
#cgo LDFLAGS: -framework Foundation
#cgo LDFLAGS: -framework Cocoa
#include <stdbool.h>

void requestAttention(bool continuous);
*/
import "C"

func requestAttention(continuous bool) {
	C.requestAttention(C._Bool(continuous))
}
