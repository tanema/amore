package main

import (
	"fmt"

	// These are lua wrapped code that will be made accessible to lua
	_ "github.com/tanema/amore/gfx/wrap"
	_ "github.com/tanema/amore/input"

	"github.com/tanema/amore/runtime"
)

func main() {
	if err := runtime.Run("main.lua"); err != nil {
		fmt.Println(err)
	}
}
