package Input

import (
	"CHIP8_GO/CPU"
	"CHIP8_GO/Global"
	"fmt"
	"os"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

var (
	forward_count    int
	Keyboard_timeout bool

	// Control the Keys Pressed (CHIP8/SCHIP 16 Keys)
	KeyPressedCHIP8 = map[uint16]pixelgl.Button{
		0:  pixelgl.KeyX,
		1:  pixelgl.Key1,
		2:  pixelgl.Key2,
		3:  pixelgl.Key3,
		4:  pixelgl.KeyQ,
		5:  pixelgl.KeyW,
		6:  pixelgl.KeyE,
		7:  pixelgl.KeyA,
		8:  pixelgl.KeyS,
		9:  pixelgl.KeyD,
		10: pixelgl.KeyZ,
		11: pixelgl.KeyC,
		12: pixelgl.Key4,
		13: pixelgl.KeyR,
		14: pixelgl.KeyF,
		15: pixelgl.KeyV,
	}

	// Control the Keys Pressed (Emulator Features, without repetition)
	KeyPressedUtils = map[uint16]pixelgl.Button{
		0: pixelgl.KeyP, // Pause
		1: pixelgl.Key9, // Debug
		2: pixelgl.Key0, // Reset
		3: pixelgl.Key6, // Change Color Theme
		4: pixelgl.KeyK, // Create Savestate
		5: pixelgl.KeyL, // Load Savestate
		6: pixelgl.KeyM, // Change video resolution
		7: pixelgl.KeyN, // Fullscreen
		8: pixelgl.KeyJ, // Show / Hide FPS
		9: pixelgl.KeyH, // Draw Mode (Drawflag or at 60Hz)
	}

	// Control the Keys Pressed (Emulator Features, with repetition)
	KeyPressedUtilsRep = map[uint16]pixelgl.Button{
		0: pixelgl.KeyI, // CPU Cycle Rewind
		1: pixelgl.KeyO, // CPU Cycle Forward
		2: pixelgl.Key7, // Decrease CPU Clock
		3: pixelgl.Key8, // Increase CPU Clock
	}
)

func Keyboard() {

	// Handle 16 keys from Chip8 / Schip
	for index, key := range KeyPressedCHIP8 {
		if Global.Win.Pressed(key) {
			// Used to slow down key press to some programs
			if !Keyboard_timeout {

				CPU.Key[index] = 1

				// Just enable if this quirk is enabled
				if CPU.Keyboard_slow_press {
					Keyboard_timeout = true
				}

			}
		} else {
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
					// Show messages
					if CPU.Debug {
						fmt.Printf("\t\tPAUSE mode Disabled\n")
					}
					Global.TextMessageStr = "PAUSE mode Disabled"
					Global.ShowMessage = true

				} else {
					CPU.Pause = true
					// Show messages
					if CPU.Debug {
						fmt.Printf("\t\tPAUSE mode Enabled\n")
					}
					Global.TextMessageStr = "PAUSE mode Enabled"
					Global.ShowMessage = true
					forward_count = 0
				}
			}

			// Debug
			if index == 1 {
				if CPU.Debug {
					CPU.Debug = false
					Global.SizeYused = 1.0
					// Show messages
					if CPU.Debug {
						fmt.Printf("\t\tDEBUG mode Disabled\n")
					}
					Global.TextMessageStr = "DEBUG mode Disabled"
					Global.ShowMessage = true
				} else {
					CPU.Debug = true
					Global.SizeYused = 0.7
					// Show messages
					if CPU.Debug {
						fmt.Printf("\t\tDEBUG mode Enabled\n")
					}
					Global.TextMessageStr = "DEBUG mode Enabled"
					Global.ShowMessage = true
				}
			}

			// Reset
			if index == 2 {

				CPU.PC = 0x200 // start at 0x200 (default CHIP8)
				CPU.Stack = [16]uint16{}
				CPU.SP = 0
				CPU.V = [16]byte{}
				CPU.I = 0
				CPU.Graphics = [128 * 64]byte{}
				Global.DrawFlag = false
				CPU.DelayTimer = 0
				CPU.SoundTimer = 0
				CPU.Key = [CPU.KeyArraySize]byte{}
				CPU.Cycle = 0
				CPU.Rewind_index = 0
				// If paused, remove the pause to continue CPU Loop
				if CPU.Pause {
					CPU.Pause = false
				}
				CPU.SCHIP = false
				Global.SizeX = 64
				Global.SizeY = 32
				CPU.CPU_Clock_Speed = 500
				CPU.Memory = CPU.MemoryCleanSnapshot

				// If in ETI-660 HW Mode, update some values
				if Global.Hybrid_ETI_660_HW {
					CPU.PC = 0x600               // start at 0x600 for ETI-600 HW (Hybrid)
					if CPU.ETI660_64x32_screen { // test if 64x32 quirk is necessary
						Global.SizeX = 64
						Global.SizeY = 32
					} else {
						Global.SizeX = 64 // Default ETI-660 screen size
						Global.SizeY = 48
					}
				}

				// Show messages
				if CPU.Debug {
					fmt.Printf("\t\tReset\n")
				}
				Global.TextMessageStr = "Reset"
				Global.ShowMessage = true
			}

			// Color Theme
			if index == 3 {
				Global.Color_theme += 1

				if Global.Color_theme > 7 {
					Global.Color_theme = 0
				}
				// Show messages
				if CPU.Debug {
					fmt.Printf("\t\tColor theme: %d\n", Global.Color_theme)
				}
				Global.TextMessageStr = fmt.Sprintf("Color theme: %d", Global.Color_theme)
				Global.ShowMessage = true
			}

			// Create Save State
			if index == 4 {
				// Create SaveState
				CPU.SaveStateWrite()

				// Show messages
				if CPU.Debug {
					fmt.Printf("\n\t\tSavestate Created\n")
				}
				Global.TextMessageStr = "Savestate Created"
				Global.ShowMessage = true

			}

			// Load Save State
			if index == 5 {

				saveFile := Global.SavestateFolder + string(os.PathSeparator) + Global.Game_signature

				if _, err := os.Stat(saveFile); err == nil {
					// File Exist

					// Read SaveState
					CPU.SaveStateRead()

					// Refresh the screen
					Global.DrawFlag = true
					// Show messages
					if CPU.Debug {
						fmt.Printf("\n\t\tSavestate Loaded\n")
					}
					Global.TextMessageStr = "Savestate Loaded"
					Global.ShowMessage = true

				} else if os.IsNotExist(err) {
					// File does NOT exist

					// Show messages
					if CPU.Debug {
						fmt.Printf("\n\t\tSavestate not found\n")
					}
					Global.TextMessageStr = "Savestate not found"
					Global.ShowMessage = true

				} else {
					fmt.Println("ERROR checking Savestate file!")
					os.Exit(0)
				}

			}

			// Change video resolution
			if index == 6 {

				// If the mode is smaller than the number of resolutions available increment (-4 to avoid the biggest ones)
				if Global.ResolutionCounter < len(Global.Settings)-4 {
					Global.ResolutionCounter++
				} else {
					Global.ResolutionCounter = 0 // reset Global.ResolutionCounter
				}

				Global.ActiveSetting = &Global.Settings[Global.ResolutionCounter]

				if Global.IsFullScreen {
					Global.Win.SetMonitor(Global.ActiveSetting.Monitor)
				} else {
					Global.Win.SetMonitor(nil)
				}
				Global.Win.SetBounds(pixel.R(0, 0, float64(Global.ActiveSetting.Mode.Width), float64(Global.ActiveSetting.Mode.Height)))

				// Show messages
				if CPU.Debug {
					fmt.Printf("\t\tResolution mode[%d]: %dx%d @ %dHz\n", Global.ResolutionCounter, Global.ActiveSetting.Mode.Width, Global.ActiveSetting.Mode.Height, Global.ActiveSetting.Mode.RefreshRate)
				}
				Global.TextMessageStr = fmt.Sprintf("Resolution mode[%d]: %dx%d @ %dHz", Global.ResolutionCounter, Global.ActiveSetting.Mode.Width, Global.ActiveSetting.Mode.Height, Global.ActiveSetting.Mode.RefreshRate)
				Global.ShowMessage = true

				Global.CenterWindow()

			}

			// Fullscreen
			if index == 7 {
				if Global.IsFullScreen {
					// Switch to windowed and backup the correct monitor.
					Global.Win.SetMonitor(nil)
					Global.IsFullScreen = false

					Global.CenterWindow()

					// Show messages
					if CPU.Debug {
						fmt.Printf("\n\t\tFullscreen Disabled\n")
					}
					Global.TextMessageStr = "Fullscreen Disabled"
					Global.ShowMessage = true

				} else {
					// Switch to fullscreen.
					Global.Win.SetMonitor(Global.ActiveSetting.Monitor)
					Global.IsFullScreen = true

					// Show messages
					if CPU.Debug {
						fmt.Printf("\n\t\tFullscreen Enabled\n")
					}
					Global.TextMessageStr = "Fullscreen Enabled"
					Global.ShowMessage = true
				}
				Global.Win.SetBounds(pixel.R(0, 0, float64(Global.ActiveSetting.Mode.Width), float64(Global.ActiveSetting.Mode.Height)))

			}

			// FPS
			if index == 8 {
				Global.ShowFPS = !Global.ShowFPS
			}

			// Draw Mode (DrawFlag OR at 60Hz)
			if index == 9 {
				Global.OriginalDrawMode = !Global.OriginalDrawMode

				// Show messages
				if Global.OriginalDrawMode {
					if CPU.Debug {
						fmt.Printf("\n\t\tDrawMode: DrawFlag\n")
					}
					Global.TextMessageStr = "DrawMode: DrawFlag"
					Global.ShowMessage = true
				} else {
					if CPU.Debug {
						fmt.Printf("\n\t\tDrawMode: @60Hz\n")
					}
					Global.TextMessageStr = "@60Hz"
					Global.ShowMessage = true
				}
			}

		}
	}

	// Handle other emulator keys (with key repetition)
	for index, key := range KeyPressedUtilsRep {

		select {
		case <-CPU.KeyboardClock.C:

			if Global.Win.Pressed(key) {

				// Rewind CPU
				if index == 0 {
					if CPU.Rewind_mode {
						if CPU.Pause {
							// Clear forward_count
							forward_count = 0
							// Search for track limit history
							// Rewind_buffer size minus [0] used for current value
							// (-2 because I use Rewind_buffer +1 to identify the last vector number)
							if CPU.Rewind_index < CPU.Rewind_buffer-2 {
								// Take care of the first loop
								if CPU.Cycle == 1 {
									Global.InputDrawFlag = true // Sinalize Graphics to Update the screen
									Global.Win.Update()
									// Show messages
									if CPU.Debug {
										fmt.Printf("\t\tRewind mode - Nothing to rewind (Cycle 0)\n")
									}
									Global.TextMessageStr = "Rewind mode - Nothing to rewind (Cycle 0)"
									Global.ShowMessage = true
								} else {
									// Update values, reading the track records
									CPU.PC = CPU.PC_track[CPU.Rewind_index+1]
									CPU.Stack = CPU.Stack_track[CPU.Rewind_index+1]
									CPU.SP = CPU.SP_track[CPU.Rewind_index+1]
									CPU.V = CPU.V_track[CPU.Rewind_index+1]
									CPU.I = CPU.I_track[CPU.Rewind_index+1]
									CPU.Graphics = CPU.GFX_track[CPU.Rewind_index+1]
									Global.DrawFlag = CPU.DF_track[CPU.Rewind_index+1]
									CPU.DelayTimer = CPU.DT_track[CPU.Rewind_index+1]
									CPU.SoundTimer = CPU.ST_track[CPU.Rewind_index+1]
									CPU.Key = [CPU.KeyArraySize]byte{}
									CPU.Cycle = CPU.Cycle - 2
									CPU.Rewind_index = CPU.Rewind_index + 1
									// Call a CPU Cycle
									CPU.Interpreter()
									// Show messages
									if CPU.Debug {
										fmt.Printf("\t\tRewind mode - Rewind_index:= %d\n", CPU.Rewind_index)
									}
									Global.TextMessageStr = fmt.Sprintf("Rewind mode - Rewind_index:= %d", CPU.Rewind_index)
									Global.ShowMessage = true
								}
							} else {
								// Show messages
								if CPU.Debug {
									fmt.Printf("\t\tRewind mode - END OF TRACK HISTORY!!!\n")
								}
								Global.TextMessageStr = fmt.Sprintf("Rewind mode - END OF TRACK HISTORY!!!")
								Global.ShowMessage = true
							}
						}
					} else {
						// Show messages
						if CPU.Debug {
							fmt.Printf("\t\tRewind mode DISABLED!!!\n")
						}
						Global.TextMessageStr = fmt.Sprintf("Rewind mode DISABLED!!!")
						Global.ShowMessage = true
					}
				}

				// Cycle Step Forward Key
				if index == 1 {
					if CPU.Pause {
						// If inside the rewind loop, search for cycles inside it
						// DO NOT update the track records in this stage
						if CPU.Rewind_index > 0 {
							CPU.PC = CPU.PC_track[CPU.Rewind_index-1]
							CPU.Stack = CPU.Stack_track[CPU.Rewind_index-1]
							CPU.SP = CPU.SP_track[CPU.Rewind_index-1]
							CPU.V = CPU.V_track[CPU.Rewind_index-1]
							CPU.I = CPU.I_track[CPU.Rewind_index-1]
							CPU.Graphics = CPU.GFX_track[CPU.Rewind_index-1]
							Global.DrawFlag = CPU.DF_track[CPU.Rewind_index-1]
							CPU.DelayTimer = CPU.DT_track[CPU.Rewind_index-1]
							CPU.SoundTimer = CPU.ST_track[CPU.Rewind_index-1]
							CPU.Key = [CPU.KeyArraySize]byte{}
							CPU.Rewind_index -= 1
							CPU.Interpreter()
							// Show messages
							if CPU.Debug {
								fmt.Printf("\t\tForward mode - Rewind_index := %d\n", CPU.Rewind_index)
							}
							Global.TextMessageStr = fmt.Sprintf("Forward mode - Rewind_index := %d", CPU.Rewind_index)
							Global.ShowMessage = true
							// Return to real time, forward CPU normally and UPDATE de tracks
						} else {
							CPU.Interpreter()
							// Show messages
							if CPU.Debug {
								fmt.Printf("\t\tForward mode - Forward %d cycles\n", forward_count)
							}
							forward_count++
							Global.TextMessageStr = fmt.Sprintf("Forward mode - Forward %d cycles", forward_count)
							Global.ShowMessage = true
						}
					}
				}

				// Decrease CPU Clock Speed
				if index == 2 {
					tmp := CPU.CPU_Clock_Speed
					if (CPU.CPU_Clock_Speed - time.Duration(CPU.CPU_Clock_decrease_rate)) > 0 {
						CPU.CPU_Clock_Speed -= time.Duration(CPU.CPU_Clock_decrease_rate)
						CPU.CPU_Clock = time.NewTicker(time.Second / CPU.CPU_Clock_Speed)
						// Show messages
						if CPU.Debug {
							fmt.Printf("\t\tDecreasing CPU Clock: %d Hz  -->  %d Hz\n", tmp, CPU.CPU_Clock_Speed)
						}
						Global.TextMessageStr = fmt.Sprintf("Decreasing CPU Clock: %d Hz  -->  %d Hz", tmp, CPU.CPU_Clock_Speed)
						Global.ShowMessage = true
					} else {
						// Reached minimum CPU Clock Speed (1 Hz)
						CPU.CPU_Clock_Speed = 1
						CPU.CPU_Clock = time.NewTicker(time.Second / CPU.CPU_Clock_Speed)
						// Show messages
						if CPU.Debug {
							fmt.Printf("\t\tDecreasing CPU Clock: %d Hz  -->  %d Hz\n", tmp, CPU.CPU_Clock_Speed)
						}
						Global.TextMessageStr = fmt.Sprintf("Decreasing CPU Clock: %d Hz  -->  %d Hz", tmp, CPU.CPU_Clock_Speed)
						Global.ShowMessage = true
					}
				}

				// Increase CPU Clock Speed
				if index == 3 {
					tmp := CPU.CPU_Clock_Speed
					if (CPU.CPU_Clock_Speed + time.Duration(CPU.CPU_Clock_increase_rate)) <= CPU.CPU_Clock_Speed_Max {
						// If Clock Speed = 1, return to multiples of 'CPU.CPU_Clock_increase_rate'
						if CPU.CPU_Clock_Speed == 1 {
							CPU.CPU_Clock_Speed += time.Duration(CPU.CPU_Clock_increase_rate - 1)
							CPU.CPU_Clock.Stop()
							CPU.CPU_Clock = time.NewTicker(time.Second / CPU.CPU_Clock_Speed)
							// Show messages
							if CPU.Debug {
								fmt.Printf("\t\tIncreasing CPU Clock: %d Hz  -->  %d Hz\n", tmp, CPU.CPU_Clock_Speed)
							}
							Global.TextMessageStr = fmt.Sprintf("Increasing CPU Clock: %d Hz  -->  %d Hz", tmp, CPU.CPU_Clock_Speed)
							Global.ShowMessage = true
						} else {
							CPU.CPU_Clock_Speed += time.Duration(CPU.CPU_Clock_increase_rate)
							CPU.CPU_Clock.Stop()
							CPU.CPU_Clock = time.NewTicker(time.Second / CPU.CPU_Clock_Speed)
							// Show messages
							if CPU.Debug {
								fmt.Printf("\t\tIncreasing CPU Clock: %d Hz  -->  %d Hz\n", tmp, CPU.CPU_Clock_Speed)
							}
							Global.TextMessageStr = fmt.Sprintf("Increasing CPU Clock: %d Hz  -->  %d Hz", tmp, CPU.CPU_Clock_Speed)
							Global.ShowMessage = true
						}
					} else {
						// Reached Maximum CPU Clock Speed (maxCPUClockAllowed Hz)
						CPU.CPU_Clock_Speed = CPU.CPU_Clock_Speed_Max
						CPU.CPU_Clock.Stop()
						CPU.CPU_Clock = time.NewTicker(time.Second / CPU.CPU_Clock_Speed)
						// Show messages
						if CPU.Debug {
							fmt.Printf("\t\tIncreasing CPU Clock: Maximum CPU Clock Allowed reached: %d Hz\n", CPU.CPU_Clock_Speed)
						}
						Global.TextMessageStr = fmt.Sprintf("Increasing CPU Clock: Maximum CPU Clock Allowed reached: %d Hz", CPU.CPU_Clock_Speed)
						Global.ShowMessage = true
					}
				}

			}

		default:
			// No timer to handle
		}

	}
}
