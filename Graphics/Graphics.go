package Graphics

import (
	"time"
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
	"Chip8/CPU"
)

var (
	win		*pixelgl.Window
)

const (
	sizeX		= 64
	sizeY		= 32
	screenWidth	= float64(1024)
	screenHeight	= float64(768)
	keyboard_tmout	= 10	// Milliseconds
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


func renderGraphics() {
	cfg := pixelgl.WindowConfig{
		Title:  "Chip8",
		Bounds: pixel.R(0, 0, screenWidth, screenHeight),
		VSync:  true,
	}
	var err error
	win, err = pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}
}


func drawGraphics(graphics [64 * 32]byte) {
	// Background color
	win.Clear(colornames.Black)
	imd := imdraw.New(nil)
	imd.Color = pixel.RGB(1, 1, 1)
	screenWidth := win.Bounds().W()
	width := screenWidth/sizeX
	height := screenHeight/sizeY

	for gfxindex := 0 ; gfxindex < len(CPU.Graphics) ; gfxindex++ {
		if (CPU.Graphics[gfxindex] ==1 ){
			if (gfxindex < 64) {
				x := gfxindex
				y := 31
				imd.Push(pixel.V ( width*float64(x), height*float64(y) ) )
				imd.Push(pixel.V ( width*float64(x)+width, height*float64(y)+height ) )
				imd.Rectangle(0)
			} else {
				y := 31
				x := 0
				nro := gfxindex
				for nro >= 64 {
					nro -= 64
					y = y - 1
					x = nro
				}
				imd.Push(pixel.V ( width*float64(x), height*float64(y) ) )
				imd.Push(pixel.V ( width*float64(x)+width, height*float64(y)+height ) )
				imd.Rectangle(0)
			}
		}
	}

	imd.Draw(win)
	win.Update()
}


func Keyboard() {
	for index, key := range CPU.KeyPressed {
		if win.Pressed(key) {
			CPU.Key[index] = 1

			// CPU.Pause Key
			if index == 16 {
				if CPU.Pause {
					CPU.Pause = false
					fmt.Printf("\t\tPAUSE mode Disabled\n")
					time.Sleep(100 * keyboard_tmout * time.Millisecond)
				} else {
					CPU.Pause = true
					fmt.Printf("\t\tPAUSE mode Enabled\n")
					time.Sleep(100 * keyboard_tmout * time.Millisecond)
				}
			}

			// Rewind CPU
			if index == 17 {
				if CPU.Pause {
					// Search for track limit history
					// Rewind_buffer size minus [0] used for current value
					// (-2 because I use Rewind_buffer +1 to identify the last vector number)
					if CPU.Rewind_index < CPU.Rewind_buffer -2 {
						// Take care of the first loop
						if (CPU.Cycle == 1) {
							fmt.Printf("\t\tRewind mode - Nothing to rewind (Cycle 0)\n")
							drawGraphics(CPU.Graphics)
							time.Sleep(100 * keyboard_tmout * time.Millisecond)
						} else {
							// Update values, reading the track records
							CPU.PC		= CPU.PC_track[CPU.Rewind_index +1]
							CPU.Stack	= CPU.Stack_track[CPU.Rewind_index +1]
							CPU.SP		= CPU.SP_track[CPU.Rewind_index +1]
							CPU.V		= CPU.V_track[CPU.Rewind_index +1]
							CPU.I		= CPU.I_track[CPU.Rewind_index +1]
							CPU.Graphics	= CPU.GFX_track[CPU.Rewind_index +1]
							CPU.DrawFlag	= CPU.DF_track[CPU.Rewind_index +1]
							CPU.DelayTimer	= CPU.DT_track[CPU.Rewind_index +1]
							CPU.SoundTimer	= CPU.ST_track[CPU.Rewind_index +1]
							// For now, reset the keys every rewind cycle
							CPU.Key		= [32]byte{}
							CPU.Cycle	= CPU.Cycle - 2
							CPU.Rewind_index= CPU.Rewind_index +1
							// Call a CPU Cycle
							CPU.Interpreter()
							time.Sleep(keyboard_tmout * time.Millisecond)
							fmt.Printf("\t\tRewind mode - Rewind_index:= %d\n\n", CPU.Rewind_index)
						}
					} else {
						fmt.Printf("\t\tRewind mode - END OF TRACK HISTORY!!!\n")
						time.Sleep(100 * keyboard_tmout * time.Millisecond)
					}
				}
			}

			// Cycle Step Forward Key
			if index == 18 {
				if CPU.Pause {
					// If inside the rewind loop, search for cycles inside it
					// DO NOT update the track records in this stage
					if CPU.Rewind_index > 0 {
						CPU.PC		= CPU.PC_track[CPU.Rewind_index -1]
						CPU.Stack	= CPU.Stack_track[CPU.Rewind_index -1]
						CPU.SP		= CPU.SP_track[CPU.Rewind_index -1]
						CPU.V		= CPU.V_track[CPU.Rewind_index -1]
						CPU.I		= CPU.I_track[CPU.Rewind_index -1]
						CPU.Graphics	= CPU.GFX_track[CPU.Rewind_index -1]
						CPU.DrawFlag	= CPU.DF_track[CPU.Rewind_index -1]
						CPU.DelayTimer	= CPU.DT_track[CPU.Rewind_index -1]
						CPU.SoundTimer	= CPU.ST_track[CPU.Rewind_index -1]
						CPU.Key		= [32]byte{}
						CPU.Rewind_index	-= 1
						CPU.Interpreter()
						time.Sleep(keyboard_tmout * time.Millisecond)
						fmt.Printf("\t\tForward mode - Rewind_index := %d\n\n", CPU.Rewind_index)
					// Return to real time, forward CPU normally and UPDATE de tracks
					} else {
						CPU.Interpreter()
						time.Sleep(keyboard_tmout * time.Millisecond)
						fmt.Printf("\t\tForward mode\n\n")
					}
				}
			}

			// Reset
			if index == 19 {
				CPU.PC			= 0x200
				CPU.Stack		= [16]uint16{}
				CPU.SP			= 0
				CPU.V			= [16]byte{}
				CPU.I			= 0
				CPU.Graphics		= [64 * 32]byte{}
				CPU.DrawFlag		= false
				CPU.DelayTimer		= 0
				CPU.SoundTimer		= 0
				CPU.Key			= [32]byte{}
				CPU.Cycle		= 0
				CPU.Rewind_index	= 0
				// If paused, remove the pause to continue CPU Loop
				if CPU.Pause {
					CPU.Pause = false
				}
			}

		}else {
			CPU.Key[index] = 0
		}
	}

}


func Run() {

	// Set up render system
	renderGraphics()

	// Main Infinite Loop
	for !win.Closed() {

		// Esc to quit program
		if win.Pressed(pixelgl.KeyEscape) {
			break
		}

		// Handle Keys pressed
		Keyboard()

		//// Calls CPU Interpreter ////
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

		// If necessary, DRAW
		if CPU.DrawFlag {
			drawGraphics(CPU.Graphics)
		}

		// Draw Graphics on Console
		//drawGraphicsConsole()

		// Update Input Events
		win.UpdateInput()
	}

}
