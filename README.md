# CHIP-8 / SCHIP Emulator

CHIP-8 / SCHIP Emulator writen in GO with simple code to be easy to be studied and understood.

**CHIP-8 Space Invaders Game:**

<img width="430" alt="invaders" src="https://github.com/cassianoperin/CHIP-8_GO/blob/master/images/invaders.png">

**Superchip (SCHIP) Space Fight 2091:**

<img width="430" alt="invaders" src="https://github.com/cassianoperin/CHIP-8_GO/blob/master/images/spacefight2091.png">

## Features
* Pause and resume emulation
* Reset emulation
* Step Forward CPU Cycles for Debug
* Step Back (Rewind) CPU Cycles for Debug
* Online Debug mode

## Compile Instructions

1) MAC
* Install GO:

	 `brew install go`

* Install library requisites:

	`go get github.com/faiface/pixel/pixelgl`

	`go get github.com/faiface/beep`

	`go get github.com/faiface/beep/mp3`

	`go get github.com/faiface/beep/speaker`

* Compile:

	`go build chip8.go`

2) Windows
* Install GO (32 bits):

	https://golang.org/dl/

* Install GCC / mingw-w64

	https://mingw-w64.org/doku.php/download/mingw-builds

* Add GO and Mingw-64 bin folder in PATH variable

* Install library requisites:

	`go get github.com/faiface/pixel/pixelgl`

	`go get github.com/faiface/beep`

	`go get github.com/faiface/beep/mp3`

	`go get github.com/faiface/beep/speaker`

* If you receive some glfw missing file error, download the zip file from https://github.com/go-gl/glfw and extract the archive into the folder $GOPATH\vendor\github.com\go-gl\glfw\v3.2

* Compile:

	`go build chip8.go`


## Usage

1. Run:
	`$chip8 ROM_NAME`

2. Keys
- Original COSMAC Keyboard Layout:

	`1` `2` `3` `C`

	`4` `5` `6` `D`

	`7` `8` `9` `E`

	`A` `0` `B` `F`

- **Keys used in this emulator:**

	`1` `2` `3` `4`

	`Q` `W` `E` `R`

	`A` `S` `D` `F`

	`Z` `X` `C` `V`

	`P`: Pause and Resume emulation

	`[`: Step back (rewind) one CPU cycle **in Pause Mode** (for debug and study purposes)

	`]`: Step forward one CPU cycle in **Pause Mode** (for debug and study purposes)

	`9`: Enable / Disable Debug Mode

	`0`: Reset

	`ESC`: Exit emulator


## Documentation
[Cowgod's Chip-8 Technical Reference](http://devernay.free.fr/hacks/chip8/C8TECH10.HTM#0.0)

[How to write an emulator (CHIP-8 interpreter) — Multigesture.net](http://www.multigesture.net/articles/how-to-write-an-emulator-chip-8-interpreter/)

[Wikipedia - CHIP-8](https://en.wikipedia.org/wiki/CHIP-8)

[HP48 Superchip](https://github.com/Chromatophore/HP48-Superchip)

[SCHIP](http://devernay.free.fr/hacks/chip8/schip.txt)

[trapexit chip-8 documentation](https://github.com/trapexit/chip-8_documentation)

[CHIP‐8-Extensions-Reference](https://github.com/mattmikolay/chip-8/wiki/CHIP%E2%80%908-Extensions-Reference)

[Thomas Daley Wiki (Game Compatibility)](https://github.com/tomdaley92/kiwi-8/issues/9)

[David Winter Documentation](http://vanbeveren.byethost13.com/stuff/CHIP8.pdf?i=2)



## TODO LIST

1. CHIP8 - Key pressing cause slowness
2. CHIP8 - Improve draw method (Rewrite graphics mode to just draw the differences from each frame)
3. CHIP8 - Implement a correct 60 FPS control
4. ALL - Rewind mode make emulation slow due to arrays and graphics processing
5. Implement a Key Timer
6. Migrate from pixel to SDL2
7. Games that uses low res AND schip draw functions should shift 2 bytes (instead od 4) and scroll N/2 lines
8. Sound glitches in Windows
