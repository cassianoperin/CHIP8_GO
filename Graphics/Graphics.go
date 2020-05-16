package Graphics

import (
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font/basicfont"
	"github.com/faiface/pixel/imdraw"

	"Chip8/CPU"
	"Chip8/Sound"
	"Chip8/Input"
	"Chip8/Global"
)



const (
	screenWidth	= float64(1024)
	screenHeight	= float64(768)
)

var (
	// FPS
	drawCounter	= 0			// imd.Draw per second counter
	updateCounter	= 0			// Win.Updates per second counter
	textFPS		*text.Text	// On screen FPS counter
	textFPSstr	string		// String with the FPS counter

	// Video modes and Fullscreen
	atlas = text.NewAtlas(basicfont.Face7x13, text.ASCII)
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
		Title:  Global.WindowTitle,
		Bounds: pixel.R(0, 0, screenWidth, screenHeight),
		VSync:  false,
		Resizable: false,
		Undecorated: false,
		NoIconify: false,
		AlwaysOnTop: true,
	}
	var err error
	Global.Win, err = pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	// Retrieve all monitors.
	monitors := pixelgl.Monitors()

	// Map the video modes available
	for i := 0; i < len(monitors); i++ {
		// Retrieve all video modes for a specific monitor.
		modes := monitors[i].VideoModes()
		for j := 0; j < len(modes); j++ {
			Global.Settings = append(Global.Settings, Global.Setting{
				Monitor: monitors[i],
				Mode:    &modes[j],
			})
		}

	}

	Global.ActiveSetting = &Global.Settings[0]

	//Initialize FPS Text
	textFPS = text.New(pixel.V(10, 740), atlas)
}


func drawGraphics(graphics [128 * 64]byte) {

	// Background color
	Global.Win.Clear(colornames.Black)
	imd := imdraw.New(nil)
	imd.Color = pixel.RGB(1, 1, 1)
	textFPS.Color = colornames.Red

	//Select Color Schema
	if Global.Color_theme != 0 {

		switch color_theme := Global.Color_theme ; {

		case color_theme == 1:
			Global.Win.Clear(colornames.White)

			imd.Color = colornames.Black

		case color_theme == 2:
			imd.Color = colornames.Lightgreen

		case color_theme == 3:
			Global.Win.Clear(colornames.Dimgray)
			imd.Color = colornames.Lightgreen

		case color_theme == 4:
			imd.Color = colornames.Steelblue

		case color_theme == 5:
			Global.Win.Clear(colornames.Darkgray)
			imd.Color = colornames.Steelblue

		case color_theme == 6:
			imd.Color = colornames.Indianred
			textFPS.Color = colornames.Steelblue

		case color_theme == 7:
			Global.Win.Clear(colornames.Darkgray)
			imd.Color = colornames.Indianred
			textFPS.Color = colornames.Steelblue
		}

	}

	screenWidth	:= Global.Win.Bounds().W()
	width		:= screenWidth/CPU.SizeX
	height		:= ((screenHeight)/CPU.SizeY)

	// If in SCHIP mode, read the entire vector. If in Chip8 mode, read from 0 to 2047 only
	for gfxindex := 0 ; gfxindex < int(CPU.SizeX) * int(CPU.SizeY) ; gfxindex++ {
		if (CPU.Graphics[gfxindex] == 1 ) {

			// Column
			x := gfxindex % int(CPU.SizeX)
			// Line
			y := gfxindex / int(CPU.SizeX)
			// Needs to be inverted to IMD Draw function before
			y = (int(CPU.SizeY) - 1) - y

			//draw_rectangle(10, 10, 50, 50, red)
			imd.Push(pixel.V ( width * float64(x)         , height * float64(y)          ) )
			imd.Push(pixel.V ( width * float64(x) + width , height * float64(y) + height ) )
			imd.Rectangle(0)
		}

	}

	// Draw Graphics to the screen
	imd.Draw(Global.Win)
	drawCounter ++	// Increment the draws per second counter

	// Draw text to the screen
	if Global.ShowFPS {
		textFPS.Clear()
		fmt.Fprintf(textFPS, textFPSstr)
		textFPS.Draw(Global.Win, pixel.IM.Scaled(textFPS.Orig, 2))
	}

}






func Run() {

	// Get game signature
	CPU.Get_game_signature()

	// Set up render system
	renderGraphics()

	// Print initial resolution
	// if debug {
		fmt.Printf("Resolution mode[%d]: %dx%d @ %dHz\n",Global.ResolutionCounter ,Global.ActiveSetting.Mode.Width, Global.ActiveSetting.Mode.Height, Global.ActiveSetting.Mode.RefreshRate)
	// }

	// Print Message if using SCHIP Hack
	if CPU.SCHIP_TimerHack {
		fmt.Printf("SCHIP DelayTimer Clock Hack ENABLED\n")
	}

	// Create a clean memory needed by some games on reset
	CPU.MemoryCleanSnapshot = CPU.Memory

	// Identify special games that needs legacy opcodes
	CPU.Handle_legacy_opcodes()

	// Remap keys to a better experience
	Input.Remap_keys()



	// Main Infinite Loop
	for !Global.Win.Closed() {

		// Esc to quit program
		if Global.Win.Pressed(pixelgl.KeyEscape) {
			break
		}

		// Handle Keys pressed
		Input.Keyboard()

		// Handle Input flags
		if Global.InputDrawFlag {
			drawGraphics(CPU.Graphics)
		}

		// ---------- Every Cycle Control the clocks!!! ---------- //

		// CPU Clock
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
				// if CPU.DrawFlag {
				// 	drawGraphics(CPU.Graphics)
				// }

				// Draw Graphics on Console
				//drawGraphicsConsole()

			// Independent of CPU CLOCK, Sound and Delay Timers runs at 60Hz
			case <-CPU.TimersClock.C:
				// When ticker run (60 times in a second, check de DelayTimer)
				// SCHIP Uses a hack to decrease DT faster to gain speed
				if !CPU.SCHIP_TimerHack {
					if CPU.DelayTimer > 0 {
						CPU.DelayTimer--
					}
				}

				// When ticker run (60 times in a second, check de SoundTimer)
				if CPU.SoundTimer > 0 {
					if CPU.SoundTimer == 1 {
						go Sound.PlaySound(Sound.Beep_buffer)
					}
					CPU.SoundTimer--
				}

			//SCHIP Speed hack, decrease DT faster
			case <-CPU.SCHIP_TimerClockHack.C:
				if CPU.SCHIP_TimerHack {
					// Decrease faster than usual 60Hz
					if CPU.DelayTimer > 0 {
						CPU.DelayTimer--
					}
				}


			// 60 FPS Control - Update the screen
			case <-CPU.FPS.C:
				// Instead of draw screen every time drawflag is set, draw at 60Hz
				drawGraphics(CPU.Graphics)
				// Update the screen after draw
				Global.Win.Update()
				updateCounter++	// Increment the updates per second counter


			// Once per second count the number of draws and Win Updates
			case <-CPU.FPSCounter.C:

				// Update the values to print on screen
				textFPSstr = fmt.Sprintf("FPS: %d\tDraws: %d\tCPU Speed: %d Hz", updateCounter, drawCounter, CPU.CPU_Clock_Speed)
				drawCounter = 0
				updateCounter = 0


			default:
				// No timer to handle
		}


		// Update Input Events
		Global.Win.UpdateInput()

	}

}
