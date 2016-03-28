// +build js

package ui

import (
	"github.com/gopherjs/gopherjs/js"
)

var CurrentWindow *Window

func InitWindowAndContext(config WindowConfig) (*Window, Context, error) {
	return &Window{}, Context(js.MakeWrapper(0)), nil
}

func (window *Window) SetTitle(title string)                                         {}
func (window *Window) Minimize()                                                     {}
func (window *Window) Maximize()                                                     {}
func (window *Window) WarpMouseInWindow(x, y int)                                    {}
func (window *Window) SetGrab(grabbed bool)                                          {}
func (window *Window) SetMinimumSize(w, h int)                                       {}
func (window *Window) SetPosition(x, y int)                                          {}
func (window *Window) GetPosition() (x, y int)                                       { return 0, 0 }
func (window *Window) SwapBuffers()                                                  {}
func (window *Window) Raise()                                                        {}
func (window *Window) Destroy()                                                      {}
func (window *Window) GetDrawableSize() (int, int)                                   { return 0, 0 }
func (window *Window) SetIcon(path string) error                                     { return nil }
func (window *Window) IsMouseGrabbed() bool                                          { return false }
func (window *Window) IsVisible() bool                                               { return false }
func (window *Window) Alert(title, message string) error                             { return nil }
func (window *Window) Confirm(title, message string) bool                            { return false }
func (window *Window) ShowMessageBox(title, message string, buttons []string) string { return "" }
func (window *Window) HasFocus() bool                                                { return false }
func (window *Window) HasMouseFocus() bool                                           { return false }
func (window *Window) RequestAttention(continuous bool)                              {}
