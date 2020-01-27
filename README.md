# GOLANG CHIP-8 Emulator

CHIP-8 Emulator writen in GO with simple code to be easy to be studied and understood.

<img width="430" alt="invaders" src="https://github.com/cassianoperin/CHIP-8_GO/blob/master/images/invaders.png">

## Features
* Pause and resume emulation
* Step Forward CPU Cycle for Debug

## Requirements
* GO
* go get github.com/faiface/pixel/pixelgl
* go get github.com/faiface/beep
* go get github.com/faiface/beep/mp3
* go get github.com/faiface/beep/speaker

## Usage

1. Run:
	`$ go run chip8 ROM_NAME`

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
	
	`[`: Step forward one CPU cycle with paused emulation (for debug and study purposes)

	`9`: Reset

	`ESC`: Exit emulator


## Documentation
[Cowgod's Chip-8 Technical Reference](http://devernay.free.fr/hacks/chip8/C8TECH10.HTM#0.0)

[How to write an emulator (CHIP-8 interpreter) — Multigesture.net](http://www.multigesture.net/articles/how-to-write-an-emulator-chip-8-interpreter/)

[Wikipedia - CHIP-8](https://en.wikipedia.org/wiki/CHIP-8)

## TODO LIST

1. Equalize game speed (some games runs too fast, other slow)
2. Key pressing cause slowness
3. Improve draw method
4. Implement a correct 60 FPS control
5. Rewrite graphics mode to just draw the differences from each frame
6. Test on Windows and Linux
7. Test Clock Program (not working)
