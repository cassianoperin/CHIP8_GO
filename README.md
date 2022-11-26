# CHIP-8 / CHIP-8 HIRES / SCHIP Emulator

CHIP-8 / CHIP-8 HIRES / SCHIP Emulator writen in GO with simple code to be easy to be studied and understood.

Optional GUI (Graphical user interface) made with fyne.io.


**CHIP-8 Space Invaders Game** | **Superchip (SCHIP) Spacefight 2091!**
:-------------------------:|:-------------------------:
<img width="430" alt="invaders" src="https://github.com/cassianoperin/CHIP8/blob/master/images/invaders.gif">  |  <img width="430" alt="invaders" src="https://github.com/cassianoperin/CHIP8/blob/master/images/spacefight2091.gif">

**Blinky (Color theme)** | **Single Dragon (Color Theme)**
:-------------------------:|:-------------------------:
<img width="430" alt="blinky" src="https://github.com/cassianoperin/CHIP8/blob/master/images/blinky.png"> | <img width="430" alt="singledragon" src="https://github.com/cassianoperin/CHIP8/blob/master/images/singledragon.png">

**CHIP-8 HiRes Astro Dodge (Color theme)** | **On-screen Debug (Color theme)**
:-------------------------:|:-------------------------:
<img width="430" alt="astro" src="https://github.com/cassianoperin/CHIP8/blob/master/images/astro.png"> | <img width="430" alt="debug" src="https://github.com/cassianoperin/CHIP8/blob/master/images/debug.png">

## Features
* ![100%](https://progress-bar.dev/100) Pause and resume emulation
* ![100%](https://progress-bar.dev/100) Reset emulation
* ![100%](https://progress-bar.dev/100) Step Forward CPU Cycles for Debug
* ![100%](https://progress-bar.dev/100) Step Back (Rewind) CPU Cycles for Debug (Need to be activated on command line with -Rewind)
* ![100%](https://progress-bar.dev/100) On-screen Debug mode
* ![100%](https://progress-bar.dev/100) Increase and Decrease CPU Clock Speed
* ![100%](https://progress-bar.dev/100) Color Themes
* ![100%](https://progress-bar.dev/100) Save States
* ![100%](https://progress-bar.dev/100) Fullscreen
* ![100%](https://progress-bar.dev/100) Binary and Hexadecimal rom format
* ![100%](https://progress-bar.dev/100) Emulation Status from all games I have to test


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

## EMULATOR Build Instructions

### MAC
* Install GO:

	 `brew install go`

* Compile:

	`go build -ldflags="-s -w" chip8.go`

### Windows

GO allows to create a Windows executable file using a MacOS:

#### 1. Install mingw-w64 (support the GCC compiler on Windows systems):

`brew install mingw-w64`

#### 2. Compile:

- 32 bits:

`env GOOS="windows" GOARCH="386"   CGO_ENABLED="1" CC="i686-w64-mingw32-gcc"   go build -ldflags="-s -w"`

- 64 bits:

`env GOOS="windows" GOARCH="amd64" CGO_ENABLED="1" CC="x86_64-w64-mingw32-gcc" go build -ldflags="-s -w"`

* If you receive the message when running the executable, you need to ensure that the video drivers supports OpenGL (or the virtual driver in the case of virtualization).

* If you receive this message : "APIUnavailable: WGL: The driver does not appear to support OpenGL", please update your graphics driver os just copy the Mesa3D library from https://fdossena.com/?p=mesa/index.frag  (opengl32.dll) to the executable folder (really slow).

#### 4. Compress binaries

`brew install upx`

`upx <binary_file>`


### Linux

Instructions to build using Ubuntu.

#### 1. Install requisites:

`sudo apt install pkg-config libgl1-mesa-dev licxi-dev libxcursor-dev libxrandr-dev libxinerama-dev libxi-dev mesa-utils build-essential xorg-dev upx`

#### 2. Build:

- 32 bits:

`env GOOS="linux" GOARCH="386"   CGO_ENABLED="1" go build -ldflags="-s -w"`

- 64 bits:

`env GOOS="linux" GOARCH="amd64" CGO_ENABLED="1" go build -ldflags="-s -w"`

#### 3. Compress binaries:

`upx <binary_file>`


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



