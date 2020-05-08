package Input

import (
	"fmt"
	"os"
	"time"
	"Chip8/CPU"
	"github.com/hajimehoshi/ebiten"
)

var (
	// Initial Color Theme
	Color_theme = 2

	// Flag used to set fullscreen
	pressedKeys = map[ebiten.Key]bool{}	//Map used to make manually IsKeyJustPressed
)

// ------------------------ Remap Keys ------------------------ //

// Platform: SCHIP
// Game: "Blinky [Hans Christian Egeberg, 1991].ch8"
// MD5: fb3284205c90d80c3b17aeea2eedf0e4
func remap_blinky() {
	// UP
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		CPU.Key[3] = 1
	} else {
		CPU.Key[3] = 0
	}
	// Down
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		CPU.Key[6] = 1
	} else {
		CPU.Key[6] = 0
	}
	// Left
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		CPU.Key[7] = 1
	} else {
		CPU.Key[7] = 0
	}
	// Right
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		CPU.Key[8] = 1
	} else {
		CPU.Key[8] = 0
	}
}

// Platform: SCHIP
// Game: "Spacefight 2091 [Carsten Soerensen, 1992].ch8"
// MD5: f99d0e82a489b8aff1c7203d90f740c3
func remap_spacefight() {
	// Left
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		CPU.Key[3] = 1
	} else {
		CPU.Key[3] = 0
	}
	// Right
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		CPU.Key[12] = 1
	} else {
		CPU.Key[12] = 0
	}
	// Shoot
	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		CPU.Key[10] = 1
	} else {
		CPU.Key[10] = 0
	}
}

// Platform: CHIP-8
// Game: "Space Invaders [David Winter].ch8"
// MD5: a67f58742cff77702cc64c64413dc37d
func remap_invaders() {
	// Left
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		CPU.Key[4] = 1
	} else {
		CPU.Key[4] = 0
	}
	// Right
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		CPU.Key[6] = 1
	} else {
		CPU.Key[6] = 0
	}
	// Shoot
	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		CPU.Key[5] = 1
	} else {
		CPU.Key[5] = 0
	}
}


// ------------------------ Main Input Functions ------------------------ //

// Control the 16 buttons of Chip8 / SCHIP
func Keyboard_chip8() {

	// ESC - Exit Emulator
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		os.Exit(0)
	}

	// ----------------- Remaps ----------------- //

	// Platform: SCHIP
	// Game: "Blinky [Hans Christian Egeberg, 1991].ch8"
	// MD5: fb3284205c90d80c3b17aeea2eedf0e4
	if (CPU.Game_signature == "121A322E303020432E20") {
		remap_blinky()

	} else if (CPU.Game_signature == "12245370616365466967") {
		remap_spacefight()

	} else if (CPU.Game_signature == "1225535041434520494E") {
		remap_invaders()

	// -------------- Default Keys -------------- //
	} else {
		// Key 0
		if ebiten.IsKeyPressed(ebiten.KeyX) {
			CPU.Key[0] = 1
		} else {
			CPU.Key[0] = 0
		}
		// Key 1
		if ebiten.IsKeyPressed(ebiten.Key1) {
			CPU.Key[1] = 1
		} else {
			CPU.Key[1] = 0
		}
		// Key 2
		if ebiten.IsKeyPressed(ebiten.Key2) {
			CPU.Key[2] = 1
		} else {
			CPU.Key[2] = 0
		}
		// Key 3
		if ebiten.IsKeyPressed(ebiten.Key3) {
			CPU.Key[3] = 1
		} else {
			CPU.Key[3] = 0
		}
		// Key 4
		if ebiten.IsKeyPressed(ebiten.KeyQ) {
			CPU.Key[4] = 1
		} else {
			CPU.Key[4] = 0
		}
		// Key 5
		if ebiten.IsKeyPressed(ebiten.KeyW) {
			CPU.Key[5] = 1
		} else {
			CPU.Key[5] = 0
		}
		// Key 6
		if ebiten.IsKeyPressed(ebiten.KeyE) {
			CPU.Key[6] = 1
		} else {
			CPU.Key[6] = 0
		}
		// Key 7
		if ebiten.IsKeyPressed(ebiten.KeyA) {
			CPU.Key[7] = 1
		} else {
			CPU.Key[7] = 0
		}
		// Key 8
		if ebiten.IsKeyPressed(ebiten.KeyS) {
			CPU.Key[8] = 1
		} else {
			CPU.Key[8] = 0
		}
		// Key 9
		if ebiten.IsKeyPressed(ebiten.KeyD) {
			CPU.Key[9] = 1
		} else {
			CPU.Key[9] = 0
		}
		// Key 10
		if ebiten.IsKeyPressed(ebiten.KeyZ) {
			CPU.Key[10] = 1
		} else {
			CPU.Key[10] = 0
		}
		// Key 11
		if ebiten.IsKeyPressed(ebiten.KeyC) {
			CPU.Key[11] = 1
		} else {
			CPU.Key[11] = 0
		}
		// Key 12
		if ebiten.IsKeyPressed(ebiten.Key4) {
			CPU.Key[12] = 1
		} else {
			CPU.Key[12] = 0
		}
		// Key 13
		if ebiten.IsKeyPressed(ebiten.KeyR) {
			CPU.Key[13] = 1
		} else {
			CPU.Key[13] = 0
		}
		// Key 14
		if ebiten.IsKeyPressed(ebiten.KeyF) {
			CPU.Key[14] = 1
		} else {
			CPU.Key[14] = 0
		}
		// Key 15
		if ebiten.IsKeyPressed(ebiten.KeyV) {
			CPU.Key[15] = 1
		} else {
			CPU.Key[15] = 0
		}
	}

}


// Control the emulator buttons, with a specific timer to be less sensitive
func Keyboard_emulator() {

	// CPU.Pause Key
	// IsKeyJustPressed make the entire emulator slow, so made it manually
	if ebiten.IsKeyPressed(ebiten.KeyP) {
		if !pressedKeys[ebiten.KeyP] {
			// Key has been pressed
			if CPU.Pause {
				CPU.Pause = false
				fmt.Printf("\t\tPAUSE mode Disabled\n")
			} else {
				CPU.Pause = true
				fmt.Printf("\t\tPAUSE mode Enabled\n")
			}
		}
 		pressedKeys[ebiten.KeyP] = true
	} else {
		pressedKeys[ebiten.KeyP] = false
	}


	// // Rewind CPU
	if ebiten.IsKeyPressed(ebiten.KeyI) {
		select {
			case <- CPU.KeyboardClock.C:

				if CPU.Pause {
					// Search for track limit history
					// Rewind_buffer size minus [0] used for current value
					// (-2 because I use Rewind_buffer +1 to identify the last vector number)
					if CPU.Rewind_index < CPU.Rewind_buffer -2 {
						// Take care of the first loop
						if (CPU.Cycle == 1) {
							fmt.Printf("\t\tRewind mode - Nothing to rewind (Cycle 0)\n")
						} else {
							// Update values, reading the track records
							CPU.PC			= CPU.PC_track[CPU.Rewind_index +1]
							CPU.Stack			= CPU.Stack_track[CPU.Rewind_index +1]
							CPU.SP			= CPU.SP_track[CPU.Rewind_index +1]
							CPU.V			= CPU.V_track[CPU.Rewind_index +1]
							CPU.I			= CPU.I_track[CPU.Rewind_index +1]
							CPU.Graphics		= CPU.GFX_track[CPU.Rewind_index +1]
							CPU.DrawFlag		= CPU.DF_track[CPU.Rewind_index +1]
							CPU.DelayTimer		= CPU.DT_track[CPU.Rewind_index +1]
							CPU.SoundTimer		= CPU.ST_track[CPU.Rewind_index +1]
							CPU.Key			= [16]byte{}
							CPU.Cycle			= CPU.Cycle - 2
							CPU.Rewind_index	= CPU.Rewind_index +1
							// Call a CPU Cycle
							CPU.Interpreter()
							fmt.Printf("\t\tRewind mode - Rewind_index:= %d\n\n", CPU.Rewind_index)
						}
					} else {
						fmt.Printf("\t\tRewind mode - END OF TRACK HISTORY!!!\n")
					}
				}

			default:
			// No timer to handle
		}

	}


	// Cycle Step Forward Key
	if ebiten.IsKeyPressed(ebiten.KeyO) {
		select {
			case <- CPU.KeyboardClock.C:

				if CPU.Pause {
					// If inside the rewind loop, search for cycles inside it
					// DO NOT update the track records in this stage
					if CPU.Rewind_index > 0 {
						CPU.PC			= CPU.PC_track[CPU.Rewind_index -1]
						CPU.Stack			= CPU.Stack_track[CPU.Rewind_index -1]
						CPU.SP			= CPU.SP_track[CPU.Rewind_index -1]
						CPU.V			= CPU.V_track[CPU.Rewind_index -1]
						CPU.I			= CPU.I_track[CPU.Rewind_index -1]
						CPU.Graphics		= CPU.GFX_track[CPU.Rewind_index -1]
						CPU.DrawFlag		= CPU.DF_track[CPU.Rewind_index -1]
						CPU.DelayTimer		= CPU.DT_track[CPU.Rewind_index -1]
						CPU.SoundTimer		= CPU.ST_track[CPU.Rewind_index -1]
						CPU.Key			= [16]byte{}
						CPU.Rewind_index	-= 1
						// Call a CPU Cycle
						CPU.Interpreter()
						fmt.Printf("\t\tForward mode - Rewind_index := %d\n\n", CPU.Rewind_index)
					// Return to real time, forward CPU normally and UPDATE de tracks
					} else {
						// Call a CPU Cycle
						CPU.Interpreter()
						fmt.Printf("\t\tForward mode\n\n")
					}
				}

			default:
				// No timer to handle
		}

	}


	// Debug
	// IsKeyJustPressed make the entire emulator slow, so made it manually
	if ebiten.IsKeyPressed(ebiten.Key9) {
		if !pressedKeys[ebiten.Key9] {
			// Key has been pressed
			if CPU.Debug {
				CPU.Debug = false
				fmt.Printf("\t\tDEBUG mode Disabled\n")
			} else {
				CPU.Debug = true
				fmt.Printf("\t\tDEBUG mode Enabled\n")
			}
		}
 		pressedKeys[ebiten.Key9] = true
	} else {
		pressedKeys[ebiten.Key9] = false
	}


	// Reset
	// IsKeyJustPressed make the entire emulator slow, so made it manually
	if ebiten.IsKeyPressed(ebiten.Key0) {
		if !pressedKeys[ebiten.Key0] {
			// Key has been pressed
			CPU.PC			= 0x200
			CPU.Stack			= [16]uint16{}
			CPU.SP			= 0
			CPU.V			= [16]byte{}
			CPU.I			= 0
			CPU.Graphics		= [128 * 64]byte{}
			CPU.DrawFlag		= false
			CPU.DelayTimer		= 0
			CPU.SoundTimer		= 0
			CPU.Key			= [16]byte{}
			CPU.Cycle			= 0
			CPU.Rewind_index	= 0
			// If paused, remove the pause to continue CPU Loop
			if CPU.Pause {
				CPU.Pause = false
			}
			CPU.SCHIP 		= false
			CPU.SizeX			= 64
			CPU.SizeY			= 32
			CPU.CPU_Clock_Speed = 500
			CPU.Memory 		= CPU.MemoryCleanSnapshot
		}
 		pressedKeys[ebiten.Key0] = true
	} else {
		pressedKeys[ebiten.Key0] = false
	}


	// Create Save State
	if ebiten.IsKeyPressed(ebiten.KeyK) {
		if !pressedKeys[ebiten.KeyK] {
			// Key has been pressed
			CPU.Opcode_savestate		= CPU.Opcode
			CPU.PC_savestate			= CPU.PC
			CPU.Stack_savestate			= CPU.Stack
			CPU.SP_savestate			= CPU.SP
			CPU.V_savestate			= CPU.V
			CPU.I_savestate			= CPU.I
			CPU.Graphics_savestate		= CPU.Graphics
			CPU.DelayTimer_savestate		= CPU.DelayTimer
			CPU.SoundTimer_savestate		= CPU.SoundTimer
			CPU.Cycle_savestate			= CPU.Cycle
			CPU.Rewind_index_savestate	= CPU.Rewind_index
			CPU.SCHIP_savestate			= CPU.SCHIP
			CPU.SCHIP_LORES_savestate	= CPU.SCHIP_LORES
			CPU.SizeX_savestate			= CPU.SizeX
			CPU.SizeY_savestate			= CPU.SizeY
			CPU.CPU_Clock_Speed_savestate = CPU.CPU_Clock_Speed
			CPU.Memory_savestate 		= CPU.Memory
			fmt.Printf("\n\t\tSavestate Created\n")
			// Register that have a savestate
			CPU.Savestate_created		= 1
		}
 		pressedKeys[ebiten.KeyK] = true
	} else {
		pressedKeys[ebiten.KeyK] = false
	}


	// Load Save State
	// IsKeyJustPressed make the entire emulator slow, so made it manually
	if ebiten.IsKeyPressed(ebiten.KeyL) {
		if !pressedKeys[ebiten.KeyL] {
			// Key has been pressed
			if CPU.Savestate_created == 1 {
				CPU.Opcode		= CPU.Opcode_savestate
				CPU.PC			= CPU.PC_savestate
				CPU.Stack			= CPU.Stack_savestate
				CPU.SP			= CPU.SP_savestate
				CPU.V			= CPU.V_savestate
				CPU.I			= CPU.I_savestate
				CPU.Graphics		= CPU.Graphics_savestate
				CPU.DelayTimer		= CPU.DelayTimer_savestate
				CPU.SoundTimer		= CPU.SoundTimer_savestate
				CPU.Cycle			= CPU.Cycle_savestate
				CPU.Rewind_index	= CPU.Rewind_index_savestate
				CPU.SCHIP			= CPU.SCHIP_savestate
				CPU.SCHIP_LORES	= CPU.SCHIP_LORES_savestate
				CPU.SizeX			= CPU.SizeX_savestate
				CPU.SizeY			= CPU.SizeY_savestate
				CPU.CPU_Clock_Speed	= CPU.CPU_Clock_Speed_savestate
				CPU.Memory 		= CPU.Memory_savestate
				CPU.DrawFlag		= true
				fmt.Printf("\n\t\tSavestate Loaded\n")
			} else {
				fmt.Printf("\n\t\tSavestate not loaded - No Savestate created\n")
			}
		}
 		pressedKeys[ebiten.KeyL] = true
	} else {
		pressedKeys[ebiten.KeyL] = false
	}


	// Decrease CPU Clock Speed
	if ebiten.IsKeyPressed(ebiten.Key7) {
		select {
			case <- CPU.KeyboardClock.C:

				decrease_rate := 100
				fmt.Printf("\n\t\tCurrent CPU Clock: %d Hz\n", CPU.CPU_Clock_Speed)
				if (CPU.CPU_Clock_Speed - time.Duration(decrease_rate)) > 0 {
					CPU.CPU_Clock_Speed -= time.Duration(decrease_rate)
					CPU.CPU_Clock = time.NewTicker(time.Second / CPU.CPU_Clock_Speed)
					fmt.Printf("\t\tNew CPU Clock: %d Hz\n\n", CPU.CPU_Clock_Speed)
				} else {
					// Reached minimum CPU Clock Speed (1 Hz)
					CPU.CPU_Clock_Speed = 1
					CPU.CPU_Clock = time.NewTicker(time.Second / CPU.CPU_Clock_Speed)
					fmt.Printf("\t\tNew CPU Clock: %d Hz\n\n", CPU.CPU_Clock_Speed)
				}

			default:
				// No timer to handle
		}
	}

	// Increase CPU Clock Speed
	if ebiten.IsKeyPressed(ebiten.Key8) {
		select {
			case <- CPU.KeyboardClock.C:

				increase_rate := 100
				fmt.Printf("\n\t\tCurrent CPU Clock: %d Hz\n", CPU.CPU_Clock_Speed)
				if (CPU.CPU_Clock_Speed + time.Duration(increase_rate)) <= time.Duration(CPU.TPS) {
					// If Clock Speed = 1, return to multiples of 'increase_rate'
					if CPU.CPU_Clock_Speed == 1 {
						CPU.CPU_Clock_Speed += time.Duration(increase_rate - 1)
						CPU.CPU_Clock.Stop()
						CPU.CPU_Clock = time.NewTicker(time.Second / CPU.CPU_Clock_Speed)
						fmt.Printf("\t\tNew CPU Clock: %d Hz\n\n", CPU.CPU_Clock_Speed)
					} else {
						CPU.CPU_Clock_Speed += time.Duration(increase_rate)
						CPU.CPU_Clock.Stop()
						CPU.CPU_Clock = time.NewTicker(time.Second / CPU.CPU_Clock_Speed)
						fmt.Printf("\t\tNew CPU Clock: %d Hz\n\n", CPU.CPU_Clock_Speed)
					}
				} else {
					// Reached Maximum TPS
					CPU.CPU_Clock_Speed = time.Duration(CPU.TPS)
					CPU.CPU_Clock.Stop()
					CPU.CPU_Clock = time.NewTicker(time.Second / CPU.CPU_Clock_Speed)
					fmt.Printf("\t\tNew CPU Clock: %d Hz\n\n", CPU.CPU_Clock_Speed)
				}

			default:
				// No timer to handle
		}
	}


	// Color Theme
	// IsKeyJustPressed make the entire emulator slow, so made it manually
	if ebiten.IsKeyPressed(ebiten.Key6) {
		if !pressedKeys[ebiten.Key6] {
			// Key has been pressed
			Color_theme += 1

			if Color_theme > 7 {
				Color_theme = 0
			}
		}
 		pressedKeys[ebiten.Key6] = true
	} else {
		pressedKeys[ebiten.Key6] = false
	}


	// Fullscreen
	// IsKeyJustPressed make the entire emulator slow, so made it manually
	if ebiten.IsKeyPressed(ebiten.KeyN) {
		if !pressedKeys[ebiten.KeyN] {
			// Key has been pressed
			ebiten.SetFullscreen(!ebiten.IsFullscreen())
		}
		pressedKeys[ebiten.KeyN] = true
	} else {
		pressedKeys[ebiten.KeyN] = false
	}

	// Cursor Visibility
	// IsKeyJustPressed make the entire emulator slow, so made it manually
	if ebiten.IsKeyPressed(ebiten.KeyM) {
		if !pressedKeys[ebiten.KeyM] {
			// Key has been pressed
			ebiten.SetCursorVisible(!ebiten.IsCursorVisible())
		}
 		pressedKeys[ebiten.KeyM] = true
	} else {
		pressedKeys[ebiten.KeyM] = false
	}

	// Windows Decorated
	// IsKeyJustPressed make the entire emulator slow, so made it manually
	if ebiten.IsKeyPressed(ebiten.KeyH) {
		if !pressedKeys[ebiten.KeyH] {
			// Key has been pressed
			ebiten.SetWindowDecorated(!ebiten.IsWindowDecorated())
		}
 		pressedKeys[ebiten.KeyH] = true
	} else {
		pressedKeys[ebiten.KeyH] = false
	}

	// Runnable On Unfocused
	// IsKeyJustPressed make the entire emulator slow, so made it manually
	if ebiten.IsKeyPressed(ebiten.KeyJ) {
		if !pressedKeys[ebiten.KeyJ] {
			// Key has been pressed
			ebiten.SetRunnableOnUnfocused(!ebiten.IsRunnableOnUnfocused())
		}
 		pressedKeys[ebiten.KeyJ] = true
	} else {
		pressedKeys[ebiten.KeyJ] = false
	}

	// Show TPS and FPS on screen
	// IsKeyJustPressed make the entire emulator slow, so made it manually
	if ebiten.IsKeyPressed(ebiten.KeyU) {
		if !pressedKeys[ebiten.KeyU] {
			// Key has been pressed
			CPU.ShowTPS = !CPU.ShowTPS
		}
 		pressedKeys[ebiten.KeyU] = true
	} else {
		pressedKeys[ebiten.KeyU] = false
	}

}
