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
	WindowTitle	string = "Chip-8"
	color_theme	= 0
)

const (
	screenWidth	= float64(1024)
	screenHeight	= float64(768)
	keyboard_tmout	= 30	// Milliseconds
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
		Title:  WindowTitle,
		Bounds: pixel.R(0, 0, screenWidth, screenHeight),
		VSync:  false,
	}
	var err error
	win, err = pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}
}


func drawGraphics(graphics [128 * 64]byte) {

	// Background color
	win.Clear(colornames.Black)
	imd := imdraw.New(nil)
	imd.Color = pixel.RGB(1, 1, 1)

	//Select Color Schema
	if color_theme != 0 {

		switch color_theme := color_theme ; {

		case color_theme == 1:
			win.Clear(colornames.White)
			imd.Color = colornames.Black

		case color_theme == 2:
			imd.Color = colornames.Lightgreen

		case color_theme == 3:
			win.Clear(colornames.Dimgray)
			imd.Color = colornames.Lightgreen

		case color_theme == 4:
			imd.Color = colornames.Steelblue

		case color_theme == 5:
			win.Clear(colornames.Darkgray)
			imd.Color = colornames.Steelblue

		case color_theme == 6:
			imd.Color = colornames.Indianred

		case color_theme == 7:
			win.Clear(colornames.Darkgray)
			imd.Color = colornames.Indianred
		}

	}




	screenWidth	:= win.Bounds().W()
	width		:= screenWidth/CPU.SizeX
	height		:= screenHeight/CPU.SizeY

	// If in SCHIP mode, read the entire vector. If in Chip8 mode, read from 0 to 2047 only
	for gfxindex := 0 ; gfxindex < int(CPU.SizeX) * int(CPU.SizeY) ; gfxindex++ {
		if (CPU.Graphics[gfxindex] == 1 ) {

			// Column
			x := gfxindex % int(CPU.SizeX)
			// Line
			y := gfxindex / int(CPU.SizeX)
			// Needs to be inverted to IMD Draw functoion before
			y = (int(CPU.SizeY) - 1) - y

			if CPU.Debug_v3 {
				fmt.Printf("\n\t Graphics.drawGraphics Debug: Column(X): %d, Line(Y): %d", x, y)
			}

			//draw_rectangle(10, 10, 50, 50, red)
			imd.Push(pixel.V ( width * float64(x)         , height * float64(y)          ) )
			imd.Push(pixel.V ( width * float64(x) + width , height * float64(y) + height ) )
			imd.Rectangle(0)
		}

	}

	if CPU.Debug_v3 {
		fmt.Printf("\n")
	}

	imd.Draw(win)

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
							CPU.Key		= [CPU.KeyArraySize]byte{}
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
						CPU.Key		= [CPU.KeyArraySize]byte{}
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

			// Debug
			if index == 19 {
				if CPU.Debug {
					CPU.Debug = false
					fmt.Printf("\t\tDEBUG mode Disabled\n")
					time.Sleep(50 * keyboard_tmout * time.Millisecond)
				} else {
					CPU.Debug = true
					fmt.Printf("\t\tDEBUG mode Enabled\n")
					time.Sleep(50 * keyboard_tmout * time.Millisecond)
				}
			}


			// Reset
			if index == 20 {
				CPU.PC			= 0x200
				CPU.Stack		= [16]uint16{}
				CPU.SP			= 0
				CPU.V			= [16]byte{}
				CPU.I			= 0
				CPU.Graphics		= [128 * 64]byte{}
				CPU.DrawFlag		= false
				CPU.DelayTimer		= 0
				CPU.SoundTimer		= 0
				CPU.Key			= [CPU.KeyArraySize]byte{}
				CPU.Cycle		= 0
				CPU.Rewind_index	= 0
				// If paused, remove the pause to continue CPU Loop
				if CPU.Pause {
					CPU.Pause = false
				}
				CPU.SCHIP = false
				CPU.SizeX	= 64
				CPU.SizeY	= 32
				CPU.CPU_Clock_Speed = 500
				CPU.Memory = CPU.MemoryCleanSnapshot
			}


			// Create Savestate
			if index == 24 {
				CPU.PC_savestate			= CPU.PC
				CPU.Stack_savestate		= CPU.Stack
				CPU.SP_savestate			= CPU.SP
				CPU.V_savestate			= CPU.V
				CPU.I_savestate			= CPU.I
				CPU.Graphics_savestate		= CPU.Graphics
				CPU.DelayTimer_savestate	= CPU.DelayTimer
				CPU.SoundTimer_savestate	= CPU.SoundTimer
				CPU.Cycle_savestate		= CPU.Cycle
				CPU.Rewind_index_savestate	= CPU.Rewind_index
				CPU.SCHIP_savestate		= CPU.SCHIP_savestate
				CPU.SizeX_savestate		= CPU.SizeX
				CPU.SizeY_savestate		= CPU.SizeY
				CPU.CPU_Clock_Speed_savestate = CPU.CPU_Clock_Speed
				CPU.Memory_savestate 		= CPU.Memory
				fmt.Printf("\n\t\tSavestate Created")
				time.Sleep(10 * keyboard_tmout * time.Millisecond)
			}

			// Create Load
			if index == 25 {
				CPU.PC			= CPU.PC_savestate
				CPU.Stack			= CPU.Stack_savestate
				CPU.SP			= CPU.SP_savestate
				CPU.V				= CPU.V_savestate
				CPU.I				= CPU.I_savestate
				CPU.Graphics		= CPU.Graphics_savestate
				CPU.DelayTimer		= CPU.DelayTimer_savestate
				CPU.SoundTimer		= CPU.SoundTimer_savestate
				CPU.Cycle			= CPU.Cycle_savestate
				CPU.Rewind_index		= CPU.Rewind_index_savestate
				CPU.SCHIP			= CPU.SCHIP_savestate
				CPU.SizeX			= CPU.SizeX_savestate
				CPU.SizeY			= CPU.SizeY_savestate
				CPU.CPU_Clock_Speed	= CPU.CPU_Clock_Speed_savestate
				CPU.Memory 			= CPU.Memory_savestate
				CPU.DrawFlag		= true
				time.Sleep(10 * keyboard_tmout * time.Millisecond)
				fmt.Printf("\n\t\tSavestate Loaded")
			}


			// Decrease CPU Clock Speed
			if index == 21 {
				decrease_rate := 50
				fmt.Printf("\n\t\tCurrent CPU Clock: %d Hz\n", CPU.CPU_Clock_Speed)
				if (CPU.CPU_Clock_Speed - time.Duration(decrease_rate)) > 0 {
					CPU.CPU_Clock_Speed -= time.Duration(decrease_rate)
					CPU.CPU_Clock = time.NewTicker(time.Second / CPU.CPU_Clock_Speed)
					fmt.Printf("\t\tNew CPU Clock: %d Hz\n\n", CPU.CPU_Clock_Speed)
					time.Sleep(5 * keyboard_tmout * time.Millisecond)
				} else {
					// Reached minimum CPU Clock Speed (1 Hz)
					CPU.CPU_Clock_Speed = 1
					CPU.CPU_Clock = time.NewTicker(time.Second / CPU.CPU_Clock_Speed)
					fmt.Printf("\t\tNew CPU Clock: %d Hz\n\n", CPU.CPU_Clock_Speed)
					time.Sleep(5 * keyboard_tmout * time.Millisecond)
				}
			}

			// Increase CPU Clock Speed
			if index == 22 {
				increase_rate := 50
				fmt.Printf("\n\t\tCurrent CPU Clock: %d Hz\n", CPU.CPU_Clock_Speed)
				if (CPU.CPU_Clock_Speed + time.Duration(increase_rate)) <= 3000 {
					// If Clock Speed = 1, return to multiples of 'increase_rate'
					if CPU.CPU_Clock_Speed == 1 {
						CPU.CPU_Clock_Speed += time.Duration(increase_rate - 1)
						CPU.CPU_Clock.Stop()
						CPU.CPU_Clock = time.NewTicker(time.Second / CPU.CPU_Clock_Speed)
						fmt.Printf("\t\tNew CPU Clock: %d Hz\n\n", CPU.CPU_Clock_Speed)
						time.Sleep(5 * keyboard_tmout * time.Millisecond)
					} else {
						CPU.CPU_Clock_Speed += time.Duration(increase_rate)
						CPU.CPU_Clock.Stop()
						CPU.CPU_Clock = time.NewTicker(time.Second / CPU.CPU_Clock_Speed)
						fmt.Printf("\t\tNew CPU Clock: %d Hz\n\n", CPU.CPU_Clock_Speed)
						time.Sleep(5 * keyboard_tmout * time.Millisecond)
					}
				} else {
					// Reached Maximum CPU Clock Speed (3000 Hz)
					CPU.CPU_Clock_Speed = 3000
					CPU.CPU_Clock.Stop()
					CPU.CPU_Clock = time.NewTicker(time.Second / CPU.CPU_Clock_Speed)
					fmt.Printf("\t\tNew CPU Clock: %d Hz\n\n", CPU.CPU_Clock_Speed)
					time.Sleep(5 * keyboard_tmout * time.Millisecond)
				}
			}

			// Color Theme
			if index == 23 {
				color_theme += 1

				if color_theme > 7 {
					color_theme = 0
				}
				time.Sleep(5 * keyboard_tmout * time.Millisecond)
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

		// Every Cycle Control the clock!!!
		select {
			case <- CPU.CPU_Clock.C:

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



			default:
				// No timer to handle
		}

		//SCHIP Speed hack, decrease DT faster
		if CPU.SCHIP {
			select {
			case <-CPU.SCHIP_TimerClockHack.C:
					// Decrease faster than usual 60Hz
					if CPU.DelayTimer > 0 {
						CPU.DelayTimer--
					}


				default:
					// No timer to handle
			}
		}

		// Update Input Events
		win.UpdateInput()

		// 60 FPS Control - Update the screen
		select {
		case <-CPU.FPS .C:

			win.Update()

			default:
				// No timer to handle
		}

	}

}
