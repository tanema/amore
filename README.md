# amore

An experimental/WIP game framework based on the API and workflow of Love 2D. It
is by no means stable

Objectives
==========
* Enable making games fast 
* Making games portable
 
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
    * WAV
    * MP3
    * Ogg Vorbis
    * MOD
* GFX [ref](https://love2d.org/wiki/love.graphics)
  - primitives
    * ~~polygons~~
    * ~~polyline~~
    * line styles
    * ~~color~~
  - ~~Transforms~~
    * ~~rotate~~
    * ~~scale~~
    * ~~shear~~
    * ~~offset~~
  - ~~Textures~~
  - ~~Images~~
  - ~~Font~~
  - Framebuffer / Canvas
  - Quad
  - Mesh
  - screenshot
  - SpriteBatch
  - Shader 
    * ~~Uniforms~~
    * ~~default shaders~~
    * sendTexture / texture pool
    * temporary attach to send variables to a non attached shader
* ~~Events~~
* Window [ref](https://love2d.org/wiki/love.window)
  - non simple message box
  - request attention
* ~~System~~
* ~~Mouse~~
* ~~Keyboard~~
* ~~Joystick~~
* File [ref](https://love2d.org/wiki/love.filesystem)
  - making assets come from config working dir
* ~~Timer~~
* Wiki
* Remove sdl_image dep
* Mobile Support
* Optimize

Experiments To Check out
========================

* [Asset Bundling](https://github.com/jteeuwen/go-bindata)
* amore as a lib so dependancies are easier



