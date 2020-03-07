# CHIP-8 / SCHIP Emulator

CHIP-8 / SCHIP Emulator writen in GO with simple code to be easy to be studied and understood.

**CHIP-8 Space Invaders Game:**

<img width="430" alt="invaders" src="https://github.com/cassianoperin/CHIP-8_GO/blob/master/images/invaders.gif">

**Superchip (SCHIP) Spacefight 2091!:**

<img width="430" alt="invaders" src="https://github.com/cassianoperin/CHIP-8_GO/blob/master/images/spacefight2091.gif">

## Features
* Pause and resume emulation
* Reset emulation
* Step Forward CPU Cycles for Debug
* Step Back (Rewind) CPU Cycles for Debug
* Online Debug mode
* Increase and Decrease CPU Clock Speed

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
- Original COSMAC Keyboard Layout (CHIP-8):

	`1` `2` `3` `C`

	`4` `5` `6` `D`

	`7` `8` `9` `E`

	`A` `0` `B` `F`

- Original HP48SX Keyboard Layout (SuperChip):

	`7` `8` `9` `/`

	`4` `5` `6` `*`

	`1` `2` `3` `-`

	`0` `.` `_` `+`

- **Keys used in this emulator:**

	`1` `2` `3` `4`

	`Q` `W` `E` `R`

	`A` `S` `D` `F`

	`Z` `X` `C` `V`

	`P`: Pause and Resume emulation

	`[`: Step back (rewind) one CPU cycle **in Pause Mode** (for debug and study purposes)

	`]`: Step forward one CPU cycle in **Pause Mode** (for debug and study purposes)

	`7`: Decrease CPU Clock Speed

	`8`: Increase CPU Clock Speed

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

[Chip-8 Database](https://chip-8.github.io/database/)

[Chip-8 Reference Manual](http://chip8.sourceforge.net/chip8-1.1.pdf)

[Archive chip8.com page](https://web.archive.org/web/20160401091945/http://www.chip8.com/)



## TODO LIST

1. Migrate from pixel to SDL2
2. Sound glitches in Windows
3. Implement Save States
4. Add color schemas
5. Create menus
