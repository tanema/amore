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
    * ~~WAV~~
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
  - make it work with go-bindata
* ~~Timer~~
* Cmd
  - amore new [name]
    - create 
      - main.go //with imports
      - conf.toml //with all default settings
      - assets/
        - audio/
        - fonts/
        - shaders/
        - images/
  - amore run 
    - bundle assets
    - go run
  - amore build [OS]
    - bundle assets
    - go build every OS
  - amore version
* Wiki
* Mobile Support
* Simplify Cross-Compilation possibly with [shared libraries already linked to c libs](http://blog.ralch.com/tutorial/golang-sharing-libraries/)
* Optimize

CMD notes
=========

* [Way of setting cmd version](http://technosophos.com/2014/06/11/compile-time-string-in-go.html)
* [Asset Bundling](https://github.com/jteeuwen/go-bindata) to make deployment easier

