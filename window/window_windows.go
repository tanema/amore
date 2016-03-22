// +build windows, !js

package window

/*
#cgo CFLAGS: -std=c99
#cgo CFLAGS: -DPSAPI_VERSION=1
#cgo LDFLAGS: -lpsapi
#cgo windows LDFLAGS: -lSDL2

#include <stdbool.h>
#include <windows.h>

#if defined(__WIN32)
	#include <SDL2/SDL.h>
	#include <stdlib.h>
#else
	#include <SDL.h>
#endif

void requestAttention(bool continuous) {
	SDL_SysWMinfo wminfo = {};
	SDL_VERSION(&wminfo.version);
	if (SDL_GetWindowWMInfo(window, &wminfo)) {
		FLASHWINFO flashinfo = {};
		flashinfo.cbSize = sizeof(FLASHWINFO);
		flashinfo.hwnd = wminfo.info.win.window;
		flashinfo.uCount = 1;
		flashinfo.dwFlags = FLASHW_ALL;
		if (continuous) {
			flashinfo.uCount = 0;
			flashinfo.dwFlags |= FLASHW_TIMERNOFG;
		}
		FlashWindowEx(&flashinfo);
	}
}
*/
import "C"

func requestAttention(continuous bool) {
	C.requestAttention(C._Bool(continuous))
}
