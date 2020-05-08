package Graphics

import (
	"fmt"
	"image/color"
	"golang.org/x/image/colornames"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"Chip8/CPU"
	"Chip8/Input"
)

var (
	emptyImage, _	= ebiten.NewImage(16, 16, ebiten.FilterDefault)			// Ebiten Image Declaration
)

const (
	ScreenWidth	= 1024
	ScreenHeight	=  768
)


// Print Graphics on Console
func drawGraphicsConsole() {
	newline := 64
	for index := 0; index < 64*32; index++ {
		switch index {
		case newline:
		  fmt.Printf("\n")
			newline += 64
	  }
    if CPU.Graphics[index] == 0 {
			fmt.Printf(" ")
		} else {
			fmt.Printf("#")
		}
	}
	fmt.Printf("\n")
}


// Ebiten Initialization
func init() {
	emptyImage.Fill(color.White)

	// By default disable run on unfocused
	ebiten.SetRunnableOnUnfocused(false)

	// By default disable mouse cursor into the game window
	ebiten.SetCursorVisible(false)

	// Disable VSync
	ebiten.SetVsyncEnabled(true)
}


// Ebiten Rectangle Draw Function
func rect(x, y, w, h float32, clr color.RGBA) ([]ebiten.Vertex, []uint16) {
	r := float32(clr.R) / 0xff
	g := float32(clr.G) / 0xff
	b := float32(clr.B) / 0xff
	a := float32(clr.A) / 0xff
	x0 := x
	y0 := y
	x1 := x + w
	y1 := y + h

	return []ebiten.Vertex{
		{
			DstX:   x0,
			DstY:   y0,
			SrcX:   1,
			SrcY:   1,
			ColorR: r,
			ColorG: g,
			ColorB: b,
			ColorA: a,
		},
		{
			DstX:   x1,
			DstY:   y0,
			SrcX:   1,
			SrcY:   1,
			ColorR: r,
			ColorG: g,
			ColorB: b,
			ColorA: a,
		},
		{
			DstX:   x0,
			DstY:   y1,
			SrcX:   1,
			SrcY:   1,
			ColorR: r,
			ColorG: g,
			ColorB: b,
			ColorA: a,
		},
		{
			DstX:   x1,
			DstY:   y1,
			SrcX:   1,
			SrcY:   1,
			ColorR: r,
			ColorG: g,
			ColorB: b,
			ColorA: a,
		},
	}, []uint16{0, 1, 2, 1, 2, 3}
}


// Ebiten Update Screen Function (MAIN LOOP)
func Update(screen *ebiten.Image) error {

	// ----------------------- Handle Background and Pixel colors ----------------------- //
	// Default pixel color
	var (
		pixelColor = colornames.White
		width	uint16
		height	uint16
	)


	// Select Color Schema
	if Input.Color_theme != 0 {

		switch Theme := Input.Color_theme ; {

			case Theme == 1:
				// Background
				screen.Fill(colornames.White)
				// Pixel Color
				pixelColor = colornames.Black

			case Theme == 2:
				// Background
				screen.Fill(colornames.Black)
				// Pixel Color
				pixelColor = colornames.Lightgreen

			case Theme == 3:
				// Background
				screen.Fill(colornames.Dimgray)
				// Pixel Color
				pixelColor = colornames.Lightgreen

			case Theme == 4:
				// Background
				screen.Fill(colornames.Black)
				// Pixel Color
				pixelColor = colornames.Steelblue

			case Theme == 5:
				// Background
				screen.Fill(colornames.Darkgray)
				// Pixel Color
				pixelColor = colornames.Steelblue

			case Theme == 6:
				// Background
				screen.Fill(colornames.Black)
				// Pixel Color
				pixelColor = colornames.Indianred

			case Theme == 7:
				// Background
				screen.Fill(colornames.Darkgray)
				// Pixel Color
				pixelColor = colornames.Indianred

		}
	}


	// ---------------------------------- Handle Input ---------------------------------- //

	// Handle Chip8 / SCHIP Keys
	Input.Keyboard_chip8()

	// Handle emulator Keys (specific timer)
	Input.Keyboard_emulator()

	// --------------------------------------- CPU -------------------------------------- //

	// Every Cycle Control the clock, limited by TPS configured
	select {
		case <- CPU.CPU_Clock.C:

			// Ignore if in Pause mode
			if !CPU.Pause {
				// If in Rewind Mode, every new cycle forward decrease the Rewind Index
				if CPU.Rewind_index > 0 {
					CPU.Interpreter()
					CPU.Rewind_index -= 1
					fmt.Printf("\t\tForward mode - Rewind_index := %d\n", CPU.Rewind_index)
				} else {
					// Continue run normally
					CPU.Interpreter()
				}
			}

		default:
			// No timer to handle
	}

	// --------------------------------------- DRAW ------------------------------------- //

	// DRAW When Necessary // NOT USED ANYMORE ONCE EBITEN UPDATE at 60Hz Graphics to Screen
	// if CPU.DrawFlag {
	// 	drawGraphics(CPU.Graphics)
	// }

	// Need to stay HERE
	if ebiten.IsDrawingSkipped() {
		return nil
	}

	width	= ScreenWidth  / CPU.SizeX
	height	= ScreenHeight / CPU.SizeY

	// If in SCHIP mode, read the entire vector. If in Chip8 mode, read from 0 to 2047 only
	for gfxindex := 0 ; gfxindex < int(CPU.SizeX) * int(CPU.SizeY) ; gfxindex++ {
		if (CPU.Graphics[gfxindex] == 1 ) {

			// Column
			x := gfxindex % int(CPU.SizeX)
			// Line
			y := gfxindex / int(CPU.SizeX)

			// X initial position, Y initial position, X lenght, Y height, (color.RGBA{0x00, 0x80, 0x00, 0x80})
			v, i := rect(float32(x) * float32(width), float32(y) * float32(height), float32(width), float32(height), pixelColor)
			screen.DrawTriangles(v, i, emptyImage, nil)
		}
	}

	if CPU.ShowTPS {
		ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f   FPS: %0.2f", ebiten.CurrentTPS(), ebiten.CurrentFPS() ))
	}

	return nil
}
