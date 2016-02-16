# amore

An experimental/WIP game framework based on the API and workflow of Love 2D with
usage of sdl2 and opengl 2.1. It is by no means stable.

Objectives
==========
* Enable making games easy, fast and fun with access still available to underlying mechanics.
* Making games portable, cross platform and with a single executable deployment strategy.
 
Requirements
============
* [OpenGL 2.1+ / OpenGL ES 2+](https://www.opengl.org/wiki/Getting_Started)
* [SDL2](http://libsdl.org/download-2.0.php)
* [SDL2_image](http://www.libsdl.org/projects/SDL_image/)
* libopenal-dev

TODO
=====
* Audio [ref](https://love2d.org/wiki/love.audio)
  - formats:
    * ~~WAV~~
    * MP3
    * Ogg Vorbis
    * MOD
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
  - Text
  - Video
* ~~Physics~~ (just use [github.com/neguse/go-box2d-lite](https://github.com/neguse/go-box2d-lite))
* ~~Events~~
* ~~Window~~ [ref](https://love2d.org/wiki/love.window) (Need support from SDL library)
* ~~System~~
* ~~Mouse~~
* ~~Keyboard~~
* ~~Joystick~~
* ~~File~~ [ref](https://love2d.org/wiki/love.filesystem)
* ~~Timer~~
* Wiki
* Full Platform Support
* Optimize

Notes and ideas
====

* Simplify Cross-Compilation possibly with [shared libraries already linked to c libs](http://blog.ralch.com/tutorial/golang-sharing-libraries/)
  - it would be nice to be able to provide amore as a library for each platform and not have the user install libs

Asset Bundling design:
 
[https://github.com/rakyll/statik](https://github.com/rakyll/statik) bundles everything into 
zip data and registers it with the fs when it is included. I could use this strategy
to do the same only not dependant on http.File but rather os.File. The when no assets
are bundled I can just use fs asset folder.

```golang
import (
  "github.com/tanema/amore/file"
  //dont include it and assets come from file system from the assets directory
  _ "{{project_path}}/static_assets" //include this and assets will come from bundle
)
func main() {
	tree, _ = gfx.NewImage("images/palm_tree.png") //no change in usage
}
```

