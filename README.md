# amore

An WIP game framework based on the API and workflow of Love 2D with usage of sdl2 and opengl
 
This project uses [goxjs/gl](https://github.com/goxjs/gl) so that it supports:

- **OS X**, **Linux** and **Windows** via OpenGL 2.1 backend,

- **iOS** and **Android** via OpenGL ES 2.0 backend,

and pending a UI wrapper around SDL there will be future support for:

- **Modern Browsers** (desktop and mobile) via WebGL 1.0 backend.

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
* [SDL2](http://libsdl.org/download-2.0.php)
* [SDL2_image](http://www.libsdl.org/projects/SDL_image/)

Below is some commands that can be used to install the required packages in
some Linux distributions. Some older versions of the distributions such as

On __Ubuntu 14.04 and above__, type:  
`apt-get install libsdl2{-image}-dev`  
_Note: Ubuntu 14.04 currently has broken header file in the SDL2 package that disables people from compiling against it. It will be needed to either patch the header file or install SDL2 from source._

On __Fedora 20 and above__, type:  
`yum install SDL2{,_image}-devel`

On __Arch Linux__, type:  
`pacman -S sdl2{,_image}`

On __Mac OS X__, install SDL2 via [Homebrew](http://brew.sh) like so:
`brew install sdl2{,_image}`

On __Windows__, install SDL2 via [Msys2](https://msys2.github.io) like so:
`pacman -S mingw-w64-x86_64-gcc mingw-w64-x86_64-SDL2{,_image}`

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
* Audio
  - formats:
    * ~~WAV~~
    * MP3
    * ~~Ogg Vorbis~~
    * ~~FLAC~~
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
* Wiki
* Full Platform Support (web support)
* Optimize

