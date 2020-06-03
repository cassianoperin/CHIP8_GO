# CHIP-8 / CHIP-8 HIRES / SCHIP Emulator

CHIP-8 / CHIP-8 HIRES / SCHIP Emulator writen in GO with simple code to be easy to be studied and understood.

Optional GUI (Graphical user interface) made with fyne.io.


**CHIP-8 Space Invaders Game** | **Superchip (SCHIP) Spacefight 2091!**
:-------------------------:|:-------------------------:
<img width="430" alt="invaders" src="https://github.com/cassianoperin/CHIP8/blob/master/images/invaders.gif">  |  <img width="430" alt="invaders" src="https://github.com/cassianoperin/CHIP8/blob/master/images/spacefight2091.gif">


**Blinky (Color theme)** | **Single Dragon (Color Theme)**
:-------------------------:|:-------------------------:
<img width="430" alt="blinky" src="https://github.com/cassianoperin/CHIP8/blob/master/images/blinky.png"> | <img width="430" alt="singledragon" src="https://github.com/cassianoperin/CHIP8/blob/master/images/singledragon.png">

**Optional GUI** | **Game selector screen**
:-------------------------:|:-------------------------:
<img width="430" alt="MainDark" src="https://github.com/cassianoperin/CHIP8/blob/master/GUI/images/git/MainDark.png">  |  <img width="430" alt="GamesDarkSchip" src="https://github.com/cassianoperin/CHIP8/blob/master/GUI/images/git/GamesDarkSchip.png">


## Features
* Pause and resume emulation
* Reset emulation
* Step Forward CPU Cycles for Debug
* Step Back (Rewind) CPU Cycles for Debug (Need to be activated on command line with -Rewind)
* Online Debug mode
* Increase and Decrease CPU Clock Speed
* Color Themes
* Save States
* Fullscreen
* Optional GUI

## EMULATOR Build Instructions

1) MAC
* Install GO:

	 `brew install go`

* Install library requisites:

	`go get github.com/faiface/pixel/pixelgl`

	`go get github.com/faiface/beep`

	`go get github.com/faiface/beep/wav`

	`go get github.com/faiface/beep/speaker`

* Compile:

	`go build -ldflags="-s -w" chip8.go`

2) Windows
* Install GO (32 bits):

	https://golang.org/dl/

* Install GCC / mingw-w64

	https://mingw-w64.org/doku.php/download/mingw-builds

* Add GO and Mingw-64 bin folder in PATH variable

* Install library requisites:

	`go get github.com/faiface/pixel/pixelgl`

	`go get github.com/faiface/beep`

	`go get github.com/faiface/beep/wav`

	`go get github.com/faiface/beep/speaker`

* If you receive some glfw missing file error, download the zip file from https://github.com/go-gl/glfw and extract the archive into the folder $GOPATH\vendor\github.com\go-gl\glfw\v3.2

* Compile:

	`go build -ldflags="-s -w" chip8.go`


## EMULATOR Usage


1. Run:

	`$./chip8 [options] ROM_NAME`


- Options:

	`-Debug`	Enable Debug Mode

	`-DrawFlag`	Enable Draw Graphics on each Drawflag instead @60Hz

	`-Rewind Mode`	Enable Rewind Mode

	`-SchipHack`	Enable SCHIP DelayTimer hack mode to improve speed

	`-help`		Show Command Line Interface options

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

	`I`: Step back (rewind) one CPU cycle **in Pause Mode**

	`O`: Step forward one CPU cycle in **in Pause Mode**

	`K`: Create Save State

	`L`: Load Save State

	`6`: Change Color Themes

	`7`: Decrease CPU Clock Speed

	`8`: Increase CPU Clock Speed

	`9`: Enable / Disable Debug Mode

	`0`: Reset

	`J`: Show / Hide FPS Counter

	`H`: Switch DrawMode (@DrawFlag OR @60Hz)

	`N`: Enter / Exit Fullscreen mode

	`M`: Change resolution

	`ESC`: Exit emulator

## GUI Build Instructions
1) MAC
* Install GO:

	 `brew install go`

* Install library requisites (need fyne beta/develop 1.3):

	`go get gopkg.in/ini.v1`

   	`go get golang.org/x/image/colornames`

	cd $GOPATH/src

	mkdir fyne.io && cd fyne.io

	git clone https://github.com/fyne-io/fyne.git && cd fyne && git checkout develop

* Compile:

	`go build -ldflags="-s -w" Chip8GUI.go`

2) Windows (**NOT TESTED YET!!**)

* Install GO:

	https://golang.org/dl/

* Install GCC / mingw-w64

	https://mingw-w64.org/doku.php/download/mingw-builds

* Add GO and Mingw-64 bin folder in PATH variable

* Install library requisites (need fyne beta/develop 1.3):

	 `go get gopkg.in/ini.v1`

 	 `go get golang.org/x/image/colornames`

        cd $GOPATH/src

        mkdir fyne.io && cd fyne.io

        git clone https://github.com/fyne-io/fyne.git && cd fyne && git checkout develop

* Compile:

	`go build -ldflags="-s -w" Chip8GUI.go`


## GUI Usage

1. Open the program, configure emulator ROM paths

2. Restart the program to update the Games Tab


## Documentation

### Chip-8
[Cowgod's Chip-8 Technical Reference](http://devernay.free.fr/hacks/chip8/C8TECH10.HTM#0.0)

[How to write an emulator (CHIP-8 interpreter) — Multigesture.net](http://www.multigesture.net/articles/how-to-write-an-emulator-chip-8-interpreter/)

[Wikipedia - CHIP-8](https://en.wikipedia.org/wiki/CHIP-8)

[trapexit chip-8 documentation](https://github.com/trapexit/chip-8_documentation)

[Thomas Daley Wiki (Game Compatibility)](https://github.com/tomdaley92/kiwi-8/issues/9)

[David Winter Documentation](http://vanbeveren.byethost13.com/stuff/CHIP8.pdf?i=2)

[Chip-8 Database](https://chip-8.github.io/database/)

[Archive chip8.com](https://web.archive.org/web/20160401091945/http://www.chip8.com/)

[Columbia University - Columbia University](http://www.cs.columbia.edu/~sedwards/classes/2016/4840-spring/designs/Chip8.pdf)

[Tom Swan Documentation](https://github.com/TomSwan/pips-for-vips)


### Superchip (SCHIP)
[HP48 Superchip](https://github.com/Chromatophore/HP48-Superchip)

[SCHIP](http://devernay.free.fr/hacks/chip8/schip.txt)

[Chip-8 Reference Manual](http://chip8.sourceforge.net/chip8-1.1.pdf)


### Chip-8 Extensions
[CHIP‐8 Extensions Reference](https://github.com/mattmikolay/chip-8/wiki/CHIP%E2%80%908-Extensions-Reference)

[CHIP‐8 Extensions (Github)](https://chip-8.github.io/extensions/)

[Chip8 Hybrids](https://www.ngemu.com/threads/chip8-thread.114578/page-22/)

[EMMA02 Opcodes Documentation](https://www.emma02.hobby-site.com/pseudo_chip8.html)

[MegaChip](https://github.com/gcsmith/gchip/blob/master/docs/megachip10.txt)

### GUI

[Fyne.io](https://fyne.io/)

[Go-ini](https://github.com/go-ini/ini)

## Emulator TODO LIST:

1. Turn CPU into a function / core (keep just the instructions, put debug messages and other stuff out of core?)
2. When Fullscreen on and off, need to change resolution to stabilize
3. Center Window
4. Improve debug mode
5. Migrate Draw functions to use bitwise operations
6. Implement automatic detection draw mode (if Draw > 60 use @60hz, otherwise just draw on drawflag)
7. Finish interpreter improvements in DXYN and DXY0
8. Improve Fx0A (Keyboard)
9. Vertical Brix reset without ball
10. Emulate CHIP-8X
11. Emulate MEGA-CHIP

## GUI TODO LIST:
1. Games Tab: Show image of the games
2. Crete an .APP to run directly from Finder on Mac
3. Reload game list when update paths
4. Make it work on Windows
5. Add tab for CHIP8 HiRes Games
6. Add configuration for CHIP8 HIRES ROMs Path
