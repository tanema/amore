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
  - ~~primitives~~
    * ~~polygons~~
    * ~~polyline~~
    * ~~line styles~~
    * ~~color~~
  - ~~Transforms~~
    * ~~rotate~~
    * ~~scale~~
    * ~~shear~~
    * ~~offset~~
  - ~~Textures~~
  - ~~Images~~
  - ~~Font~~
  - ~~Canvas~~
  - ~~Quad~~
  - SpriteBatch
  - Mesh
  - Stencil
  - Shader 
    * ~~Uniforms~~
    * ~~default shaders~~
    * sendTexture / texture pool
    * temporary attach to send variables to a non attached shader
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
* [Asset Bundling](https://github.com/jteeuwen/go-bindata) to make deployment easier
  - it would be nice to embed the assets and config into the binary so it is just single file deploy
  - I want to make it work with the file package though so that there is a single entry point to all assets
  - hard to do since the assets would be bundled into the final game and not amore espectially if amore is precompiled

