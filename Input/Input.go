package Input

import (
	"fmt"
	"time"
	"Chip8/CPU"
	"Chip8/Global"
	"github.com/faiface/pixel/pixelgl"
)

const (
	increase_rate = 100	// CPU Clock increase rate
	decrease_rate = 100	// CPU Clock decrease rate
	maxCPUClockAllowed = 5000
)

var (
	// Control the Keys Pressed (CHIP8/SCHIP 16 Keys)
	KeyPressedCHIP8 = map[uint16]pixelgl.Button{
		0:	pixelgl.KeyX,
		1:	pixelgl.Key1,
		2:	pixelgl.Key2,
		3:	pixelgl.Key3,
		4:	pixelgl.KeyQ,
		5:	pixelgl.KeyW,
		6:	pixelgl.KeyE,
		7:	pixelgl.KeyA,
		8:	pixelgl.KeyS,
		9:	pixelgl.KeyD,
		10:	pixelgl.KeyZ,
		11:	pixelgl.KeyC,
		12:	pixelgl.Key4,
		13:	pixelgl.KeyR,
		14:	pixelgl.KeyF,
		15:	pixelgl.KeyV,
	}

	// Control the Keys Pressed (Emulator Features, without repetition)
	KeyPressedUtils = map[uint16]pixelgl.Button{
		0:	pixelgl.KeyP,			// Pause
		1:	pixelgl.Key9,			// Debug
		2:	pixelgl.Key0,			// Reset
		3:	pixelgl.Key6,			// Change Color Theme
		4:	pixelgl.KeyK,			// Create Savestate
		5:	pixelgl.KeyL,			// Load Savestate
	}

	// Control the Keys Pressed (Emulator Features, with repetition)
	KeyPressedUtilsRep = map[uint16]pixelgl.Button{
		0:	pixelgl.KeyI,			// CPU Cycle Rewind
		1:	pixelgl.KeyO,			// CPU Cycle Forward
		2:	pixelgl.Key7,			// Decrease CPU Clock
		3:	pixelgl.Key8,			// Increase CPU Clock
	}

)


func Remap_keys() {
	// Platform: SCHIP
	// Game: "Blinky [Hans Christian Egeberg, 1991].ch8"
	// MD5: fb3284205c90d80c3b17aeea2eedf0e4
	if (Global.Game_signature == "121A322E303020432E20") {
		KeyPressedCHIP8[3] = pixelgl.KeyUp
		KeyPressedCHIP8[6] = pixelgl.KeyDown
		KeyPressedCHIP8[7] = pixelgl.KeyLeft
		KeyPressedCHIP8[8] = pixelgl.KeyRight
		Global.WindowTitle = "                                         |     Chip-8     |     Keys:     Left: ←     Right: →     Up: ↑     Down: ↓"
		fmt.Printf("Keys Remaped:\tLeft: ←\t\tRight: →\tUp: ↑\t\tDown: ↓\n\n")
	}

	// Platform: SCHIP
	// Game: "Spacefight 2091 [Carsten Soerensen, 1992].ch8"
	// MD5: f99d0e82a489b8aff1c7203d90f740c3
	if (Global.Game_signature == "12245370616365466967") {
		KeyPressedCHIP8[10] = pixelgl.KeySpace
		KeyPressedCHIP8[3] = pixelgl.KeyLeft
		KeyPressedCHIP8[12] = pixelgl.KeyRight
		Global.WindowTitle = "                                         |     Chip-8     |     Keys:     Left: ←     Right: →     Shoot: Space"
		fmt.Printf("Keys Remaped\tLeft: ←\t\tRight: →\tShoot: Space\n\n")
	}

	// Platform: CHIP-8
	// Game: "Space Invaders [David Winter].ch8"
	// MD5: a67f58742cff77702cc64c64413dc37d
	if (Global.Game_signature == "1225535041434520494E") {
		KeyPressedCHIP8[5] = pixelgl.KeySpace
		KeyPressedCHIP8[4] = pixelgl.KeyLeft
		KeyPressedCHIP8[6] = pixelgl.KeyRight
		Global.WindowTitle = "                                         |     Chip-8     |     Keys:     Left: ←     Right: →     Shoot: Space"
		fmt.Printf("Keys Remaped\tLeft: ←\t\tRight: →\tShoot: Space\n\n")
	}

	// Platform: SCHIP
	// Game: "Ant - In Search of Coke [Erin S. Catto].ch8"
	// MD5: ec7856f9db5917eb6ca14adf1f8d0df2
	if (Global.Game_signature == "12E5B20416E74207631") {
		KeyPressedCHIP8[10] = pixelgl.KeySpace
		KeyPressedCHIP8[3]  = pixelgl.KeyLeft
		KeyPressedCHIP8[12]  = pixelgl.KeyRight
		Global.WindowTitle  = "                                         |     Chip-8     |     Keys:     Left: ←     Right: →     Shoot: Space"
		fmt.Printf("Keys Remaped\tLeft: ←\t\tRight: →\tShoot: Space\n\n")
	}

}



func Keyboard() {

	// Handle 16 keys from Chip8 / Schip
	for index, key := range KeyPressedCHIP8 {
		if Global.Win.Pressed(key) {
			CPU.Key[index] = 1
		}else {
			CPU.Key[index] = 0
		}
	}


	// Handle other emulator Keys
	for index, key := range KeyPressedUtils {
		if Global.Win.JustPressed(key) {

			// CPU.Pause Key
			if index == 0 {
				if CPU.Pause {
					CPU.Pause = false
					fmt.Printf("\t\tPAUSE mode Disabled\n")
				} else {
					CPU.Pause = true
					fmt.Printf("\t\tPAUSE mode Enabled\n")
				}
			}


			// Debug
			if index == 1 {
				if CPU.Debug {
					CPU.Debug = false
					fmt.Printf("\t\tDEBUG mode Disabled\n")
				} else {
					CPU.Debug = true
					fmt.Printf("\t\tDEBUG mode Enabled\n")
				}
			}


			// Reset
			if index == 2 {
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

			// Color Theme
			if index == 3 {
				Global.Color_theme += 1

				if Global.Color_theme > 7 {
					Global.Color_theme = 0
				}
			}

			// Create Save State
			if index == 4 {
				CPU.Opcode_savestate		= CPU.Opcode
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
				CPU.SCHIP_savestate		= CPU.SCHIP
				CPU.SCHIP_LORES_savestate	= CPU.SCHIP_LORES
				CPU.SizeX_savestate		= CPU.SizeX
				CPU.SizeY_savestate		= CPU.SizeY
				CPU.CPU_Clock_Speed_savestate = CPU.CPU_Clock_Speed
				CPU.Memory_savestate 		= CPU.Memory
				fmt.Printf("\n\t\tSavestate Created\n")
				// Register that have a savestate
				CPU.Savestate_created		= 1
			}

			// Load Save State
 			if index == 5 {
				if CPU.Savestate_created == 1 {
					CPU.Opcode			= CPU.Opcode_savestate
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
					CPU.SCHIP_LORES		= CPU.SCHIP_LORES_savestate
					CPU.SizeX			= CPU.SizeX_savestate
					CPU.SizeY			= CPU.SizeY_savestate
					CPU.CPU_Clock_Speed	= CPU.CPU_Clock_Speed_savestate
					CPU.Memory 			= CPU.Memory_savestate
					CPU.DrawFlag		= true
					fmt.Printf("\n\t\tSavestate Loaded\n")
				} else {
					fmt.Printf("\n\t\tSavestate not loaded - No Savestate created\n")
				}

			}


		}
	}

	// Handle 16 keys from Chip8 / Schip
	for index, key := range KeyPressedUtilsRep {

		select {
			case <- CPU.KeyboardClock.C:

				if Global.Win.Pressed(key) {

					// Rewind CPU
					if index == 0 {
						if CPU.Pause {
							// Search for track limit history
							// Rewind_buffer size minus [0] used for current value
							// (-2 because I use Rewind_buffer +1 to identify the last vector number)
							if CPU.Rewind_index < CPU.Rewind_buffer -2 {
								// Take care of the first loop
								if (CPU.Cycle == 1) {
									fmt.Printf("\t\tRewind mode - Nothing to rewind (Cycle 0)\n")
									Global.InputDrawFlag = true // Sinalize Graphics to Update the screen
									Global.Win.Update()
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
									fmt.Printf("\t\tRewind mode - Rewind_index:= %d\n\n", CPU.Rewind_index)
								}
							} else {
								fmt.Printf("\t\tRewind mode - END OF TRACK HISTORY!!!\n")
							}
						}
					}

					// Cycle Step Forward Key
					if index == 1 {
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
								fmt.Printf("\t\tForward mode - Rewind_index := %d\n\n", CPU.Rewind_index)
							// Return to real time, forward CPU normally and UPDATE de tracks
							} else {
								CPU.Interpreter()
								fmt.Printf("\t\tForward mode\n\n")
							}
						}
					}


					// Decrease CPU Clock Speed
					if index == 2 {
						tmp	:= CPU.CPU_Clock_Speed
						if (CPU.CPU_Clock_Speed - time.Duration(decrease_rate)) > 0 {
							CPU.CPU_Clock_Speed -= time.Duration(decrease_rate)
							CPU.CPU_Clock = time.NewTicker(time.Second / CPU.CPU_Clock_Speed)
							fmt.Printf("\t\tDecreasing CPU Clock: %d Hz  -->  %d Hz\n", tmp, CPU.CPU_Clock_Speed)
						} else {
							// Reached minimum CPU Clock Speed (1 Hz)
							CPU.CPU_Clock_Speed = 1
							CPU.CPU_Clock = time.NewTicker(time.Second / CPU.CPU_Clock_Speed)
							fmt.Printf("\t\tDecreasing CPU Clock: %d Hz  -->  %d Hz\n", tmp, CPU.CPU_Clock_Speed)
						}
					}

					// Increase CPU Clock Speed
					if index == 3 {
						tmp := CPU.CPU_Clock_Speed
						if (CPU.CPU_Clock_Speed + time.Duration(increase_rate)) <= maxCPUClockAllowed {
							// If Clock Speed = 1, return to multiples of 'increase_rate'
							if CPU.CPU_Clock_Speed == 1 {
								CPU.CPU_Clock_Speed += time.Duration(increase_rate - 1)
								CPU.CPU_Clock.Stop()
								CPU.CPU_Clock = time.NewTicker(time.Second / CPU.CPU_Clock_Speed)
								fmt.Printf("\t\tIncreasing CPU Clock: %d Hz  -->  %d Hz\n", tmp, CPU.CPU_Clock_Speed)
							} else {
								CPU.CPU_Clock_Speed += time.Duration(increase_rate)
								CPU.CPU_Clock.Stop()
								CPU.CPU_Clock = time.NewTicker(time.Second / CPU.CPU_Clock_Speed)
								fmt.Printf("\t\tIncreasing CPU Clock: %d Hz  -->  %d Hz\n", tmp, CPU.CPU_Clock_Speed)
							}
						} else {
							// Reached Maximum CPU Clock Speed (maxCPUClockAllowed Hz)
							CPU.CPU_Clock_Speed = maxCPUClockAllowed
							CPU.CPU_Clock.Stop()
							CPU.CPU_Clock = time.NewTicker(time.Second / CPU.CPU_Clock_Speed)
							fmt.Printf("\t\tIncreasing CPU Clock: Maximum CPU Clock Allowed reached: %d Hz\n", CPU.CPU_Clock_Speed)
						}
					}

				}

			default:
				// No timer to handle
		}

	}
}
