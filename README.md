# Amore

[![GoDoc](https://godoc.org/github.com/tanema/amore?status.svg)](http://godoc.org/github.com/tanema/amore)

A game library for Go, inspired by Love 2D. Currently in Beta. This is a hobby library and support should not be expected.

## Objectives

* Enable making games easy, fast and fun
* Making games portable
* Single executable deployment strategy.

## Aimed Platform Support:

- **OS X**
- **Linux**
- **Windows**
- **iOS**
- **Android**

## Installation

It is recommended that you use [glide](https://github.com/Masterminds/glide) for working
with this project so that it is certain which dependancies you are vendoring.

You can get this project with glide by running

```bash
glide get github.com/tanema/amore
```

Install the amore package by running the go get command

```bash
go get -u github.com/tanema/amore/...
```

### Requirements

Amore requires [SDL2](http://libsdl.org/download-2.0.php) to operate on PC. You can install it by doing the following.

__Ubuntu 14.04 and above__, type: `apt-get install libsdl2-dev`

__Fedora 20 and above__, type: `yum install SDL2-devel`

__Arch Linux__, type: `pacman -S sdl2`

__openSUSE__, type: `zypper in libSDL2-devel`

__Mac OS X__, via [Homebrew](http://brew.sh): `brew install sdl2`

__Windows__, via [Msys2](https://msys2.github.io): `pacman -S mingw-w64-x86_64-gcc mingw-w64-x86_64-SDL2`

## Amore Command

Install the command line helper with the go install tool

```bash
go install github.com/tanema/amore/cmd
```

* `amore new` will generate initial files for a game in the current folder
* `amore bundle` will generate a file called `asset_bundle.go` with all the assets and config in ziped byte array to be included in the binary


## Example

See more examples at [github.com/tanema/amore-examples](https://github.com/tanema/amore-examples)

```golang
// Basic hello world program
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

## TODO

- [ ] font rendering, especially spacing is bad
- [ ] audio seeking broken
- [ ] use `runtime.SetFinalizer` to release assets so that the develop does not have to manage volatiles.
