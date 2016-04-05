# Amore

An WIP game framework based on the API and workflow of Love 2D with usage of sdl2 and opengl
 
This project uses [goxjs/gl](https://github.com/goxjs/gl) so that it supports:

- **OS X**, **Linux** and **Windows** via OpenGL 2.1 backend,
- **iOS** and **Android** via OpenGL ES 2.0 backend,

Objectives
==========
* Enable making games easy, fast and fun
* Making games portable
* single executable deployment strategy.
 
Installation
============

```bash
go get -u github.com/tanema/amore/...
go install github.com/tanema/amore/cmd
```

Command
=======

* `amore new` will generate initial files for a game in the current folder
* `amore bundle` will generate a file called `asset_bundle.go` with all the assets and config in ziped byte array to be included in the binary
 
Requirements
============
Amore requires [SDL2](http://libsdl.org/download-2.0.php) to operate on PC. You can install it by doing the following.

__Ubuntu 14.04 and above__, type: `apt-get install libsdl2-dev`  

__Fedora 20 and above__, type: `yum install SDL2-devel`

__Arch Linux__, type: `pacman -S sdl2`

__Mac OS X__, via [Homebrew](http://brew.sh): `brew install sdl2`

__Windows__, via [Msys2](https://msys2.github.io): `pacman -S mingw-w64-x86_64-gcc mingw-w64-x86_64-SDL2`

_Note: Ubuntu 14.04 currently has broken header file in the SDL2 package that disables people from compiling against it. It will be needed to either patch the header file or install SDL2 from source._

Example
=======

```golang
package main

import (
	"github.com/tanema/amore"
	"github.com/tanema/amore/gfx"
)

func main() {
	amore.Start(update, draw)
}

func update(deltaTime float32) {
}

func draw() {
	gfx.Print("Hello World",50, 50)
}
```

TODO
=====
* ~~Audio (wav, ogg, flac)~~
* GFX [ref](https://love2d.org/wiki/love.graphics)
  - ~~primitives (polygons, lines, color, stencil, scissor)~~
  - ~~Transforms (rotate, scale, shear, offset)~~
  - ~~Textures~~
  - ~~Font (image, ttf)~~
  - ~~Canvas (Just RGBA8/Normal)~~
  - ~~Quad~~
  - ~~Shader~~
  - ~~Images (gif, jpeg, png)~~
  - ~~Particle System~~
  - ~~Mesh~~
  - ~~SpriteBatch~~
  - ~~Text~~
  - Video
* ~~Physics~~ (just use [github.com/neguse/go-box2d-lite](https://github.com/neguse/go-box2d-lite))
* ~~Events~~
* ~~Window~~ [ref](https://love2d.org/wiki/love.window)
* ~~System~~
* ~~Mouse~~
* ~~Keyboard~~
* ~~Joystick~~
* ~~File~~ [ref](https://love2d.org/wiki/love.filesystem)
* ~~Timer~~
* ~~Asset Bundling~~
* Full Platform Support Testing
* Wiki

