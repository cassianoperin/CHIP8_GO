package CPU

import (
	"fmt"
	"math/rand"
	"os"
	"time"
	"strconv"
	"github.com/faiface/pixel/pixelgl"
	"Chip8/Fontset"
	"Chip8/Sound"

)

const (
	// Rewind Buffer Size
	// FIX IT, INCREASING IT TO 15000, emulation became slow due to vector and matrix processing
	// Exponentially increased by size of graphics array!
	// NEED TO FIND A WAY TO PROCESSES IT QUICKLY
	Rewind_buffer	uint16 = 2
	// Control the number of Keys mapped in Key Array
	KeyArraySize	byte	= 23
)

// Components
var (
	Memory		[4096]byte // Memory
	PC		uint16     // Program Counter
	Opcode		uint16     // CPU Operation Code
	Stack		[16]uint16 // Stack
	SP		uint16     // Stack Pointer
	V		[16]byte
	I		uint16
	DelayTimer	byte
	SoundTimer	byte
	TimerClock	*time.Ticker
	CPU_Clock	*time.Ticker
	CPU_Clock_Speed	time.Duration
	//Graphics	[64 * 32]byte
	Graphics	[128 * 64]byte

	// True if the screen must be drawn
	DrawFlag	bool
	Key		[KeyArraySize]byte
	Cycle		uint16

	// Control the Keys Pressed
	KeyPressed = map[uint16]pixelgl.Button{
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
		16:	pixelgl.KeyP,			// Pause
		17:	pixelgl.KeyLeftBracket,		// CPU Cycle Rewind
		18:	pixelgl.KeyRightBracket,	// CPU Cycle Forward
		19:	pixelgl.Key9,			// Debug
		20:	pixelgl.Key0,			// Reset
		21:	pixelgl.Key7,			// Decrease CPU Clock
		22:	pixelgl.Key8,			// Increase CPU Clock
	}

	// Pause (Used to Forward and Rewind CPU Cycles)
	Pause		bool =  false

	// DEBUG modes
	Debug		bool = false
	// Debug Rewind Mode
	Debug_v2	bool = false
	// Debug Draw Graphics function
	Debug_v3	bool = false
	// Time tracks
	Debug_time	bool = true

	// Rewind Variables
	Rewind_index	uint16 = 0
	PC_track	[Rewind_buffer]uint16
	SP_track	[Rewind_buffer]uint16
	I_track		[Rewind_buffer]uint16
	DT_track	[Rewind_buffer]byte
	ST_track	[Rewind_buffer]byte
	DF_track	[Rewind_buffer]bool
	V_track		[Rewind_buffer][16]byte
	Stack_track	[Rewind_buffer][16]uint16
	GFX_track	[Rewind_buffer][128 * 64]byte
	// Beep sound file
	sound_file string
	// GRAPHICS
	SizeX		float64 = 64
	SizeY		float64 = 32

	// SCHIP
	SCHIP = false
	RPL	[8]byte // HP-48 RPL user flags

	// LEGACY OPCODES / Quirks
	// Game Signature (identify games that needs legacy opcodes)
	Game_signature	string = ""
	// Enable original Chip-8 Fx55 and Fx65 (increases I)
	Legacy_Fx55_Fx65	bool	= false
	Legacy_8xy6_8xyE	bool	= false
	// FX1E undocumented feature needed by Spacefight 2091!
	FX1E_spacefight2091	bool	= false
	// DXYN sprite wrap in Bowling game
	DXYN_bowling_wrap	bool	= false
	// Resize_Quirk_00FE_00FF - Clears the screen - Must be set to True always
	Resize_Quirk_00FE_00FF	bool	= true

)


func Initialize() {
	// Initialization
	Memory		= [4096]byte{}
	PC		= 0x200
	Opcode		= 0
	Stack		= [16]uint16{}
	SP		= 0
	V		= [16]byte{}
	I		= 0
	Graphics	= [128 * 64]byte{}
	//Graphics	= [64 * 32]byte{}
	DrawFlag	= false
	DelayTimer	= 0
	SoundTimer	= 0

	// Create a ticker at 60Hz
	TimerClock	= time.NewTicker(time.Second / 60)

	// CPU Clock Speed
	CPU_Clock_Speed	= 800
	CPU_Clock	= time.NewTicker(time.Second / CPU_Clock_Speed)


	// Load CHIP-8 8x5 fontset
	// Memory address 0-79
	for i := 0; i < len(Fontset.Chip8Fontset); i++ {
		Memory[i] = Fontset.Chip8Fontset[i]
	}
	// Load SCHIP 8x10 fontset
	// Memory address 80-240
	for i := 0; i < len(Fontset.SCHIPFontset); i++ {
		Memory[i+80] = Fontset.SCHIPFontset[i]
	}

	Key		= [KeyArraySize]byte{}
	Cycle		= 0

}

func Show() {
	fmt.Printf("Cycle: %d\tOpcode: %04X(%04X)\tPC: %d(0x%X)\tSP: %d\tStack: %d\tV: %d\tI: %d\tDT: %d\tST: %d\tKey: %d\n", Cycle, Opcode, Opcode & 0xF000, PC, PC,  SP, Stack, V, I, DelayTimer, SoundTimer, Key)
}


// CPU Interpreter
func Interpreter() {

	// Start timer to measure procedures inside Interpreter
	start := time.Now()

	// Reset Flag every cycle
	DrawFlag = false

	// Read the Opcode from PC and PC+1 bytes
	Opcode = uint16(Memory[PC])<<8 | uint16(Memory[PC+1])

	// Print Cycle and Debug Information
	if Debug {
		Show()
	}

	// REWIND MODE - ARRAYS
	// Just update when not inside a reward loop
	// Otherwise navigate inside the track arrays
	if Rewind_index == 0 {
		// Every new value is recorded on the first array value or first line of the Matrix for vectors
		// Each new cycle, Shift Right the values (for arrays) and Shift lines to the end for Matrix
		for i := Rewind_buffer - 1 ; i > 0 ; i-- {
			// PC
			PC_track[i]=PC_track[i-1]
			// SP
			SP_track[i]=SP_track[i-1]
			// I
			I_track[i]=I_track[i-1]
			// DelayTimer
			DT_track[i]=DT_track[i-1]
			// SoundTimer
			ST_track[i]=ST_track[i-1]
			// DrawFlag
			DF_track[i]=DF_track[i-1]
			// V Matrix
			V_track[i]=V_track[i-1]
			// Stack Matrix
			Stack_track[i]=Stack_track[i-1]
			// Graphics Matrix
			GFX_track[i]=GFX_track[i-1]

		}

		// After store the current value in the first Array position or Matrix line
		// PC - ADD new value
		PC_track[0] = PC
		// SP - ADD new value
		SP_track[0] = SP
		// I - ADD new value
		I_track[0] = I
		// DelayTimer - ADD new value
		DT_track[0] = DelayTimer
		// SoundTimer - ADD new value
		ST_track[0] = SoundTimer
		// DrawFlag - ADD new value
		DF_track[0] = DrawFlag
		// V Matrix - ADD new vector
		V_track[0] = V
		// Stack Matrix - ADD new vector
		Stack_track[0] = Stack
		// Graphics Matrix - ADD new vector
		GFX_track[0] = Graphics

		if Debug_v2 {
			fmt.Printf("\tPC_track: %d\n", PC_track)
			fmt.Printf("\tSP_Track: %d\n", SP_track)
			fmt.Printf("\tI_Track: %d\n", I_track)
			fmt.Printf("\tDT_Track: %d\n", DT_track)
			fmt.Printf("\tST_Track: %d\n", ST_track)
			fmt.Printf("\tDF_Track: %t\n", DF_track)
			fmt.Printf("\tV_Track: %d\n", V_track)
			fmt.Printf("\tStack_Track: %d\n", Stack_track)
			//fmt.Printf("\tGFX_Track: %d\n", GFX_track)
			fmt.Printf("\n")
		}

	}

	// Debug time execution - Rewind Mode
	if Debug {
		elapsed := time.Since(start)
		fmt.Printf("\t\tTime track - Rewind Mode took: %s\n", elapsed)
	}

	// Debug time execution - Opcode Handling
	if Debug {
		start = time.Now()
	}


	// Map Opcode Family
	switch Opcode & 0xF000 {

		// ############################ 0x0000 instruction set ############################
		case 0x0000: //0NNN

			x := Opcode & 0x000F
			//fmt.Printf("\t\t%X", x)
			//os.Exit(2)


			switch Opcode & 0x00F0 {


			case 0x00E0:
				// 00E0 - CLS
				// Clear the display.
				if x == 0x0000 {
					// Clear display
					//Graphics = [64 * 32]byte{}
					Graphics = [128 * 64]byte{}
					PC += 2
					if Debug {
						fmt.Println("\t\tOpcode 00E0 executed. - Clear the display\n")
					}
					break
				}

				// 00EE - RET
				// Return from a subroutine
				// The interpreter sets the program counter to the address at the top of the stack, then subtracts 1 from the stack pointer.
				// MUST MOVE TO NEXT ADDRESS AFTER THIS (PC+=2)
				if x == 0x000E {
					//fmt.Println("   RCA 1802 Opcode 0x00EE - Return")
					PC = Stack[SP] + 2
					SP --
					if Debug {
						fmt.Printf("\t\tOpcode 00EE executed. - Return from a subroutine (PC=%d)\n", PC)
					}
					break
				}

			// 02D8
			// NON DOCUMENTED OPCODED, USED BY DEMO CLOCK Program
			// LDA 02, I // Load from memory at address I into V[00] to V[02]
		case 0x00D0:
				x := (Opcode & 0x0F00) >> 8

				if x != 2 {
					//Map if this opcode can receive a different value here
					os.Exit(2)
				}

				V[0] = byte(I)
				V[1] = byte(I) + 1
				V[2] = byte(I) + 2

				PC += 2
				if Debug {
					fmt.Printf("\t\tOpcode 02DB executed (NON DOCUMENTED). - Load from memory at address I(%d) into V[0]= %d, V[1]= %d and V[2]= %d.\n", I, I , I+1, I+2)
				}
				break

				// SCHIP - 00FF
				// Enable High-Res Mode (128 x 64 resolution)
			case 0x00F0:
				if x == 0x000F {
					SCHIP = true

					SizeX = 128
					SizeY = 64

					if Resize_Quirk_00FE_00FF {
						// Clear the screen when changing graphic mode
						Graphics	= [128 * 64]byte{}
					}

					PC += 2
					if Debug {
						fmt.Printf("\t\tSCHIP - Opcode 00FF executed. - Enable high res (128 x 64) mode.\n")
					}
					break
				// SCHIP - 00FE
				// Enable Low-Res Mode (64 x 32 resolution)
				} else if x == 0x000E {
					SCHIP = false

					SizeX = 64
					SizeY = 32

					if Resize_Quirk_00FE_00FF {
						// Clear the screen when changing graphic mode
						Graphics	= [128 * 64]byte{}
					}

					PC += 2
					if Debug {
						fmt.Printf("\t\tSCHIP - Opcode 00FE executed. - Enable low res (64 x 32) mode.\n")
					}

				// SCHIP - 00FD
				// Exit Emulator
				} else if x == 0x000D {
					fmt.Printf("SCHIP - Opcode 00FD executed. - Exit emulator.\n")
					os.Exit(0)

				// SCHIP - 00FC
				// Scroll display 4 pixels left
				} else if x == 0x000C {
					//SCHIP = true
					shift := 4
					rowsize := int(SizeX)
					//fmt.Printf(Graphics[])

					gfx_len := 0
					if SCHIP {
						gfx_len = (128 * 64)
					} else {
						gfx_len = (64 * 32)
					}


					// Run all the array
					for i := 0 ; i < gfx_len ; i++ {

						// Shift values until the last shift(4) bytes for each line
						if i < rowsize - shift{
							Graphics[i] = Graphics[i+shift]
						}

						if i == rowsize -1 {
							//Change the last 4 bytes of each line to zero
							for i := rowsize - shift ; i < rowsize ; i++ {
								Graphics[i] = 0
							}
							// Update index to next line
							rowsize += int(SizeX)
						}
					}

					DrawFlag	= true
					PC += 2
					if Debug {
						fmt.Printf("\t\tSCHIP - Opcode 00FC executed. - Scroll display 4 pixels left.\n")
					}

					// SCHIP - 00FB
					// Scroll display 4 pixels right
					} else if x == 0x000B {

						shift := 4	// Number of bytes to be shifted
						rowsize := int(SizeX)
						index := 0
						gfx_len := 0

						// Calculate the values because I'm using the same array
						if SCHIP {
							gfx_len = 128 * 64
							index = (128 * 64) - rowsize
						} else {
							gfx_len = 64 * 32
							index = (64 * 32) - rowsize
						}

						// Run all the array
						for i := gfx_len -1  ; i >= 0  ; i-- {

							// Shift values until the last shift bytes for each line
							if i >=  index + shift {
								Graphics[i] = Graphics[i - shift]
							}

							// If find the index, change the last shift(4) bytes to zero and update the index
							// To process the next line
							if i == index {
								//Change the first 4 bytes of each line to zero
								for j := index + shift - 1; j >= index  ; j-- {
									Graphics[j] = 0
								}
								// Update index to next line
								index -= int(SizeX)
							}
						}

						DrawFlag	= true
						PC += 2
						if Debug {
							fmt.Printf("\t\tSCHIP - Opcode 00FB executed. - Scroll display 4 pixels right.\n")
						}

				} else {
					fmt.Printf("\t\tOpcode 00F%X NOT IMPLEMENTED.\n", x)
					os.Exit(2)

				}

				// SCHIP - 00CN
				// Scroll display N lines down
			case 0x00C0:
					SCHIP = true

					shift := int(x) * 128

					// Shift Right N lines on Graphics Array
					for i:=len(Graphics) -1 ; i >= shift ; i-- {
						Graphics[i] = Graphics[i - shift]
					}

					// Clean the shifted display bytes
					for i:=0 ; i < shift ; i++ {
						Graphics[i] = 0
					}

					DrawFlag	= true
					PC += 2
					if Debug {
						fmt.Printf("\t\tSCHIP - Opcode 00CN executed. - Scroll display %d lines down.\n", int(x))
					}

					break


			default:
				if Debug {
					fmt.Printf("\t\tOpcode 0x0000 NOT IMPLEMENTED!!!!\n", Opcode)
				}
				os.Exit(2)
			}


		// ############################ 0x1000 instruction set ############################
		// 1nnn - JP addr
		// Jump to location nnn.
		// The interpreter sets the program counter to nnn.
		case 0x1000:
			PC = Opcode & 0x0FFF
			if Debug {
				fmt.Printf("\t\tOpcode 1nnn executed: Jump to location 0x%d\n", Opcode & 0x0FFF)
			}
			break


		// ############################ 0x2000 instruction set ############################
		// 2nnn - CALL addr
		// Call subroutine at nnn.
		// The interpreter increments the stack pointer, then puts the current PC on the top of the stack. The PC is then set to nnn.
		case 0x2000:
			SP++
			Stack[SP] = PC
			PC = Opcode & 0x0FFF
			if Debug {
				fmt.Printf("\t\tOpcode 2nnn executed: Call Subroutine at 0x%d\n", PC)
			}
			break

		// ############################ 0x3000 instruction set ############################
		// 3xkk - SE Vx, byte
		// Skip next instruction if Vx = kk.
		// The interpreter compares register Vx to kk, and if they are equal, increments the program counter by 2.
		case 0x3000:
			x := (Opcode & 0x0F00) >> 8
			kk := byte(Opcode & 0x00FF)
			if V[x] == kk {
				PC += 4
				if Debug {
					fmt.Printf("\t\tOpcode 3xk executed: V[x(%d)]:(%d) = kk(%d), skip one instruction.\n", x, V[x], kk)
				}
			} else {
				PC += 2
				if Debug {
					fmt.Printf("\t\tOpcode 3xk executed: V[x(%d)]:(%d) != kk(%d), do NOT skip one instruction.\n", x, V[x], kk)
				}
			}
			break


		// ############################ 0x4000 instruction set ############################
		// 4xkk - SNE Vx, byte
		// Skip next instruction if Vx != kk.
		// The interpreter compares register Vx to kk, and if they are not equal, increments the program counter by 2.
		case 0x4000:
			x := (Opcode & 0x0F00) >> 8
			kk := byte(Opcode & 0x00FF)
			if V[x] != kk {
				PC += 4
				if Debug {
					fmt.Printf("\t\tOpcode 4xkk executed: V[x(%d)]:%d != kk(%d), skip one instruction\n", x, V[x], kk)
				}
			} else {
				if Debug {
					fmt.Printf("\t\tOpcode 4xkk executed: V[x(%d)]:%d = kk(%d), NOT skip one instruction\n", x, V[x], kk)
				}
				PC += 2
			}
			break


		// ############################ 0x5000 instruction set ############################
		// 5xy0 - SE Vx, Vy
		// Skip next instruction if Vx = Vy.
		// The interpreter compares register Vx to register Vy, and if they are equal, increments the program counter by 2.
		case 0x5000:
			x := (Opcode & 0x0F00) >> 8
			y := (Opcode & 0x00F0) >> 4

			if (V[x] == V[y]){
				PC += 4
				if Debug {
					fmt.Printf("\t\tOpcode 5xy0 executed: V[x(%d)]:%d EQUAL V[y(%d)]:%d, SKIP one instruction\n", x, V[x], y, V[y])
				}
			} else {
				PC += 2
				if Debug {
					fmt.Printf("\t\tOpcode 5xy0 executed: V[x(%d)]:%d NOT EQUAL V[y(%d)]:%d, DO NOT SKIP one instruction\n", x, V[x], y, V[y])
				}
			}
			break


		// ############################ 0x6000 instruction set ############################
		// 6xkk - LD Vx, byte
		// Set Vx = kk.
		// The interpreter puts the value kk into register Vx.
		case 0x6000:
			x := (Opcode & 0x0F00) >> 8
			kk := byte(Opcode)

			V[x] = kk
			PC += 2
			if Debug {
				fmt.Printf("\t\tOpcode 6xkk executed: Set V[x(%d)] = %d\n", x, kk)
			}
			break


		// ############################ 0x7000 instruction set ############################
		// 7xkk - ADD Vx, byte
		// Set Vx = Vx + kk.
		// Adds the value kk to the value of register Vx, then stores the result in Vx.
		case 0x7000:
			x := (Opcode & 0x0F00) >> 8
			kk := byte(Opcode)

			V[x] += kk

			PC += 2
			if Debug {
				fmt.Printf("\t\tOpcode 7xkk executed: Add the value kk(%d) to V[x(%d)]\n", kk, x)
			}
			break


		//############################ 0x8000 instruction set ############################
		// 0x8000 instruction set
		case 0x8000:
			x := (Opcode & 0x0F00) >> 8
			y := (Opcode & 0x00F0) >> 4
			switch Opcode & 0x000F {

			// 8xy0 - LD Vx, Vy
			// Set Vx = Vy.
			// Stores the value of register Vy in register Vx.
			case 0x0000:
				V[x] = V[y]
				PC += 2
				if Debug {
					fmt.Printf("\t\tOpcode 8xy0 executed: Set V[x(%d)] = V[y(%d)]:%d\n", x, y, V[y])
				}
				break

			// Set Vx = Vx OR Vy.
			// Performs a bitwise OR on the values of Vx and Vy, then stores the result in Vx. A bitwise OR compares the corrseponding bits from two values,
			// and if either bit is 1, then the same bit in the result is also 1. Otherwise, it is 0.
			case 0x0001:
				V[x] |= V[y]
				PC += 2
				if Debug {
					fmt.Printf("\t\tOpcode 8xy1 executed: Set V[x(%d)]:%d OR V[y(%d)]:%d\n", x, V[x], y, V[y])
				}
				break

			// 8xy2 - AND Vx, Vy
			// Set Vx = Vx AND Vy.
			// Performs a bitwise AND on the values of Vx and Vy, then stores the result in Vx. A bitwise AND compares the corrseponding bits from two values, and if both bits are 1, then the same bit in the result is also 1. Otherwise, it is 0.
			case 0x0002:
				V[x] &= V[y]
				PC += 2
				if Debug {
					fmt.Printf("\t\tOpcode 8xy2 executed: Set V[x(%d)] = V[x(%d)] AND V[y(%d)]\n", x, x, y)
				}
				break

			// 8xy3 - XOR Vx, Vy
			// Set Vx = Vx XOR Vy.
			// Performs a bitwise exclusive OR on the values of Vx and Vy, then stores the result in Vx. An exclusive OR compares the corrseponding bits from two values,
			// and if the bits are not both the same, then the corresponding bit in the result is set to 1. Otherwise, it is 0.
			case 0x0003:
				if Debug {
					fmt.Printf("\t\tOpcode 8xy3 executed:  V[x(%d)]:%d XOR V[y(%d)]:%d \n", x, V[x], y, V[y])
				}
				V[x] ^= V[y]
				PC += 2
				break

			// 8xy4 - ADD Vx, Vy
			// Set Vx = Vx + Vy, set VF = carry.
			// The values of Vx and Vy are added together. If the result is greater than 8 bits (i.e., > 255,) VF is set to 1, otherwise 0.
			// Only the lowest 8 bits of the result are kept, and stored in Vx.
			case 0x0004:
				if ( V[x] + V[y] < V[x]) {
					V[0xF] = 1
				} else {
					V[0xF] = 0
				}
				if Debug {
					fmt.Printf("\t\tOpcode 8xy4 executed: Set V[x(%d)] = V[x(%d)] + V[y(%d)]\n", x, x, y)
				}
				// Old implementation, sum values, READ THE DOCS IN CASE OF PROBLEMS
				V[x] += V[y]

				PC += 2
				break


			// 8xy5 - SUB Vx, Vy
			// Set Vx = Vx - Vy, set VF = NOT borrow.
			// If Vx > Vy, then VF is set to 1, otherwise 0. Then Vy is subtracted from Vx, and the results stored in Vx.
			case 0x0005:
				if V[x] >= V[y] {
					V[0xF] = 1
				} else {
					V[0xF] = 0
				}

				V[x] -= V[y]
				PC += 2
				if Debug {
					fmt.Printf("\t\tOpcode 8xy5 executed: Set V[x(%d)] = V[x(%d)]:%d - V[y(%d)]:%d\n", x, x, V[x], y, V[y])
				}
				break

			// 8xy6 - SHR Vx {, Vy}
			// Set Vx = Vx SHR 1.
			// If the least-significant bit of Vx is 1, then VF is set to 1, otherwise 0. Then Vx is divided by 2 (SHR).
			case 0x0006:
				V[0xF] = V[x] & 0x01

				if Legacy_8xy6_8xyE {
					V[x] = V[y] >> 1
				} else {
					V[x] = V[x] >> 1
				}

				PC += 2
				if Debug {
					fmt.Printf("\t\tOpcode 8xy6 executed: Set V[x(%d)]:%d SHIFT RIGHT 1\n", x, V[x])
				}
				// Original Chip8 INCREMENT I in this instruction ###
				break


			// 8xy7 - SUBN Vx, Vy
			// Set Vx = Vy - Vx, set VF = NOT borrow.
			// If Vy > Vx, then VF is set to 1, otherwise 0. Then Vx is subtracted from Vy, and the results stored in Vx.
			case 0x0007:
				if V[x] > V[y] {
					V[0xF] = 0
				} else {
					V[0xF] = 1
				}
				if Debug {
					fmt.Printf("\t\tOpcode 8xy7 executed: Set V[x(%d)]:%d = V[y(%d)]:%d - V[x(%d)]:%d\t\t = %d \n", x, V[x], y, V[y], x, V[x], V[y] - V[x])
				}
				V[x] = V[y] - V[x]

				PC += 2
				break

			// 8xyE - SHL Vx {, Vy}
			// Set Vx = Vx SHL 1.
			// If the most-significant bit of Vx is 1, then VF is set to 1, otherwise to 0. Then Vx is multiplied by 2.
			case 0x000E:
				V[0xF] = V[x] >> 7 // Set V[F] to the Most Important Bit

				if Legacy_8xy6_8xyE {
					V[x] = V[y] << 1
				} else {
					V[x] = V[x] << 1
				}

				PC += 2
				if Debug {
					fmt.Printf("\t\tOpcode 8xyE executed: Set V[x(%d)]:%d SHIFT LEFT 1\n", x, V[x])
				}
				break

			default:
				if Debug {
					fmt.Printf("\t\tOpcode 0x8000 NOT IMPLEMENTED!!!!\n")
				}
				os.Exit(2)
			}


		// ############################ 0x9000 instruction set ############################
		// 9xy0 - SNE Vx, Vy
		// Skip next instruction if Vx != Vy.
		// The values of Vx and Vy are compared, and if they are not equal, the program counter is increased by 2.
		case 0x9000:
			x := (Opcode & 0x0F00) >> 8
			y := (Opcode & 0x00F0) >> 4

			if ( V[x] != V[y] ) {
				PC += 4
				if Debug {
					fmt.Printf("\t\tOpcode 9xy0 executed: V[x(%d)]:%d != V[y(%d)]:%d, SKIP one instruction\n", x, V[x], y, V[y])
				}
			} else {
				PC += 2
				if Debug {
					fmt.Printf("\t\tOpcode 9xy0 executed: V[x(%d)]:%d = V[y(%d)]:%d, DO NOT SKIP one instruction\n", x, V[x], y, V[y])
				}
			}
			break

		// ############################ 0xA000 instruction set ############################
		// Annn - LD I, addr
		// Set I = nnn.
		// The value of register I is set to nnn.
		case 0xA000:
			//fmt.Println("   Opcode Family: 0xA000 - Sets I to the address NNN")
			I = Opcode & 0x0FFF
			PC += 2
			if Debug {
				fmt.Printf("\t\tOpcode Annn executed: Set I = %d\n", I)
			}
			break


		// ############################ 0xB000 instruction set ############################
		// Bnnn - JP V0, addr
		// Jump to location nnn + V0.
		// The program counter is set to nnn plus the value of V0.
		case 0xB000:

			nnn := Opcode & 0x0FFF
			PC = nnn + uint16(V[0])
			if Debug {
				print ("\t\tOpcode Bnnn executed: Jump to location nnn(%d) + V[0(%d)]\n", nnn, V[0])
			}
			break


		// ############################ 0xC000 instruction set ############################
		// Cxkk - RND Vx, byte
		// Set Vx = random byte AND kk.
		// The interpreter generates a random number from 0 to 255, which is then ANDed with the value kk. The results are stored in Vx. See instruction 8xy2 for more information on AND.
		case 0xC000: // CNNN
			x := uint16(Opcode&0x0F00) >> 8
			kk := Opcode & 0x00FF
			V[x] = byte(rand.Float32()*255) & byte(kk)
			PC += 2
			if Debug {
				fmt.Printf("\t\tOpcode Cxkk executed: V[x(%d)] = %d (random byte AND kk(%d)) = %d\n", x, V[x], kk, V[x])
			}
			break


		// ############################ 0xD000 instruction set ############################
		case 0xD000: // DXYN

			var (
				x		uint16 = (Opcode & 0x0F00) >> 8
				y		uint16 = (Opcode & 0x00F0) >> 4
				n		uint16 = (Opcode & 0x000F)
				byte		uint16 = 0
				gpx_position	uint16 = 0
			)

			// Clean the colision flag
			V[0xF] = 0

			// Fix for Bowling game where the pins wrap the screen
			if DXYN_bowling_wrap {
				if V[x] + uint8(n) > (uint8(SizeX) +1)  {
					PC += 2
					break
				}
			}

			// Translate the x and Y to the Graphics Vector
			gpx_position = (uint16(V[x]) + (uint16(SizeX) * uint16(V[y])))
			// DEBUG
			//fmt.Printf ("\tGraphic vector position: %d\tValue: %d\n", gpx_position, Graphics[x + (128 * y)] )

			// SCHIP Dxy0
			// When in high res mode show a 16x16 sprite at (VX, VY)
			// If N=0, Draw in SCHIP High Resolution mode
			if n == 0 {
				// SCHIP HI-RES MODE
				// If in SCHIP mode will draw 16x16 sprites
				if SCHIP {
					// Turn n in 16 (pixel size in SCHIP Mode)
					n = 16
					if Debug {
						fmt.Printf("\t\tSCHIP - Opcode Dxy0 HI-RES MODE (%X) DRAW GRAPHICS! - Address I: %d Position V[x(%d)]: %d V[y(%d)]: %d\n" , Opcode, I, x, V[x], y, V[y])
					}

					// Print N Bytes from address I in V[x]V[y] position of the screen
					for byte = 0 ; byte < n ; byte++ {

						var (
							binary string = ""
							sprite uint8 = 0
							sprite2 uint8 = 0
						)

						// Set the sprite
						//fmt.Printf("Memory I + Byte(%d): %d (%d)\n", byte, Memory[I+ byte], Memory[I + byte])

						// DOCUMENT SPRITES
						sprite = Memory[I + (byte * 2)]
						sprite2 = Memory[I + (byte * 2) + 1]


						// Sprite in binary format
						//binary = fmt.Sprintf("%.8b", sprite)
						binary = fmt.Sprintf("%.8b%.8b", sprite,sprite2)
						// fmt.Printf("BINARY = %b\n",sprite)
						// fmt.Printf("TESTE: %.8b%.8b", sprite,sprite2)
						// fmt.Printf("\n")

						// Always print 8 bits
						for bit := 0; bit < 16 ; bit++ {

							// Convert the binary[bit] variable into an INT using Atoi method
							bit_binary, err := strconv.Atoi(fmt.Sprintf("%c", binary[bit]))
							if err == nil {
								// fmt.Println(bit_binary)
							}


							// Set the index to write the 8 bits of each pixel
							gfx_index := uint16(gpx_position) + uint16(bit) + (byte*uint16(SizeX))

							// If tryes to draw bits outside the vector size, ignore
							if ( gfx_index >= uint16(SizeX) * uint16(SizeY) ) {
								//fmt.Printf("Bigger than 2048 or 8192\n")
								continue
							}


							// If bit=1, test current graphics[index], if is already set, mark v[F]=1 (colision)
							if (bit_binary  == 1){
								// Set colision case graphics[index] is already 1
								if (Graphics[gfx_index] == 1){
									V[0xF] = 1
								}
								// After, XOR the graphics[index] (DRAW)
								Graphics[gfx_index] ^= 1
							}

						// DEBUG 2
						//fmt.Printf ("\n\tByte: %d,\tSprite: %d\tBinary: %s\tbit: %d\tIndex: %d\tbinary[bit]: %c\tGraphics[index]: %d",byte, sprite, binary, bit, index, binary[bit], Graphics[index])
						}
					}
				// SCHIP LOW-RES MODE
				// If NOT in SCHIP mode will draw 16x8 sprites
				} else {
					n = 16
					if Debug {
						fmt.Printf("\t\tSCHIP - Opcode Dxy0 LOW-RES MODE (%X DRAW GRAPHICS! - Address I: %d Position V[x(%d)]: %d V[y(%d)]: %d\n" , Opcode, I, x, V[x], y, V[y])
					}

					// Print N Bytes from address I in V[x]V[y] position of the screen
					for byte = 0 ; byte < n ; byte++ {

						var (
							binary string = ""
							sprite uint8 = 0
							//sprite2 uint8 = 0
						)

						// Set the sprite
						//fmt.Printf("Memory I + Byte(%d): %d (%d)\n", byte, Memory[I+ byte], Memory[I + byte])

						// DOCUMENT SPRITES
						sprite = Memory[I + byte]
						//sprite2 = Memory[I + (byte * 2) + 1]


						// Sprite in binary format
						binary = fmt.Sprintf("%.8b", sprite)
						//binary = fmt.Sprintf("%.8b%.8b", sprite,sprite2)
						// fmt.Printf("BINARY = %b\n",sprite)
						// fmt.Printf("TESTE: %.8b%.8b", sprite,sprite2)
						// fmt.Printf("\n")

						// Always print 8 bits
						for bit := 0; bit < 8 ; bit++ {

							// Convert the binary[bit] variable into an INT using Atoi method
							bit_binary, err := strconv.Atoi(fmt.Sprintf("%c", binary[bit]))
							if err == nil {
								// fmt.Println(bit_binary)
							}


							// Set the index to write the 8 bits of each pixel
							gfx_index := uint16(gpx_position) + uint16(bit) + (byte*uint16(SizeX))

							// If tryes to draw bits outside the vector size, ignore
							if ( gfx_index >= uint16(SizeX) * uint16(SizeY) ) {
								//fmt.Printf("Bigger than 2048 or 8192\n")
								continue
							}


							// If bit=1, test current graphics[index], if is already set, mark v[F]=1 (colision)
							if (bit_binary  == 1){
								// Set colision case graphics[index] is already 1
								if (Graphics[gfx_index] == 1){
									V[0xF] = 1
								}
								// After, XOR the graphics[index] (DRAW)
								Graphics[gfx_index] ^= 1
							}

						// DEBUG 2
						//fmt.Printf ("\n\tByte: %d,\tSprite: %d\tBinary: %s\tbit: %d\tIndex: %d\tbinary[bit]: %c\tGraphics[index]: %d",byte, sprite, binary, bit, index, binary[bit], Graphics[index])
						}
					}

				}

			// Dxyn - DRW Vx, Vy, nibble
			// Display n-byte sprite starting at memory location I at (Vx, Vy), set VF = collision.
			} else {
			// Else, Draw in Chip-8 Low Resolution mode
				if Debug {
					fmt.Printf("\t\tOpcode Dxyn(%X) DRAW GRAPHICS! - Address I: %d Position V[x]: %d V[y]: %d N: %d bytes\n" , Opcode, I, V[x], V[y], n)

				}


				// Print N Bytes from address I in V[x]V[y] position of the screen
				for byte = 0 ; byte < n ; byte++ {

					var (
						binary string = ""
						sprite uint8 = 0
					)

					// Set the sprite
					//fmt.Printf("Memory I + Byte(%d): %d (%d)\n", byte, Memory[I+ byte], Memory[I + byte])
					sprite = Memory[I + byte]

					// Sprite in binary format
					binary = fmt.Sprintf("%.8b", sprite)

					// Always print 8 bits
					for bit := 0; bit < 8 ; bit++ {

						// Convert the binary[bit] variable into an INT using Atoi method
						bit_binary, err := strconv.Atoi(fmt.Sprintf("%c", binary[bit]))
						if err == nil {
							// fmt.Println(bit_binary)
						}

						// Set the index to write the 8 bits of each pixel
						gfx_index := uint16(gpx_position) + uint16(bit) + (byte*uint16(SizeX))

						// If tryes to draw bits outside the vector size, ignore
						if ( gfx_index >= uint16(SizeX) * uint16(SizeY) ) {
							//fmt.Printf("Bigger than 2048 or 8192\n")
							continue
						}

						// If bit=1, test current graphics[index], if is already set, mark v[F]=1 (colision)
						if (bit_binary  == 1){
							// Set colision case graphics[index] is already 1
							if (Graphics[gfx_index] == 1){
								V[0xF] = 1
							}
							// After, XOR the graphics[index] (DRAW)
							Graphics[gfx_index] ^= 1
						}

					// DEBUG 2
					//fmt.Printf ("\n\tByte: %d,\tSprite: %d\tBinary: %s\tbit: %d\tIndex: %d\tbinary[bit]: %c\tGraphics[index]: %d",byte, sprite, binary, bit, index, binary[bit], Graphics[index])
					}
				}

			}

			PC += 2
			DrawFlag = true


		// ############################ 0xE000 instruction set ############################
		// 0xE000 instruction set
		case 0xE000:

			x := (Opcode & 0x0F00) >> 8
			switch Opcode & 0x00FF {


			// Ex9E - SKP Vx
			// Skip next instruction if key with the value of Vx is pressed.
			// Checks the keyboard, and if the key corresponding to the value of Vx is currently in the down position, PC is increased by 2.
			case 0x009E:
				if Key[V[x]] == 1 {
					PC += 4
					if Debug {
						fmt.Printf("\t\tOpcode Ex9E executed: Key[%d] pressed, skip one instruction\n",V[x])
					}
				} else {
					PC += 2
					if Debug {
						fmt.Printf("\t\tOpcode Ex9E executed: Key[%d] NOT pressed, continue\n",V[x])
					}
				}
				break

			// ExA1 - SKNP Vx
			// Skip next instruction if key with the value of Vx is not pressed.
			// Checks the keyboard, and if the key corresponding to the value of Vx is currently in the up position, PC is increased by 2.
			case 0x00A1:
				if Key[V[x]] == 0 {
					PC += 4
					if Debug {
						fmt.Printf("\t\tOpcode ExA1 executed: Key[%d] NOT pressed, skip one instruction\n",V[x])
					}
				} else {
					Key[V[x]] = 0
					PC += 2
					if Debug {
						fmt.Printf("\t\tOpcode ExA1 executed: Key[%d] pressed, continue\n",V[x])
					}
				}
				break
			default:
				fmt.Printf("Opcode Family E000 - Not mapped opcote: E000\n")
				os.Exit(2)
			}


		// ############################ 0xF000 instruction set ############################
		case 0xF000:

			x := (Opcode & 0x0F00) >> 8

			switch Opcode & 0x00FF {

			// Fx07 - LD Vx, DT
			// Set Vx = delay timer value.
			// The value of DT is placed into Vx.
			case 0x0007:
				V[x] = DelayTimer
				PC += 2
				if Debug {
					fmt.Printf("\t\tOpcode Fx07 executed: Set V[x(%d)] with value of DelayTimer(%d)\n", x, DelayTimer)
				}
				break

			// Fx0A - LD Vx, K
			// Wait for a key press, store the value of the key in Vx.
			// All execution stops until a key is pressed, then the value of that key is stored in Vx.
			case 0x000A:
				pressed := 0
				for i := 0 ; i < len(Key) ; i++ {
					if (Key[i] == 1){
						V[x] = byte(i)
						pressed = 1
						PC +=2
						if Debug {
							fmt.Printf("\t\tOpcode Fx0A executed: Wait for a key (Key[%d]) press -  (PRESSED)\n", i)
						}
						// Stop after find the first key pressed
						break

					}
				}
				if pressed == 0 {
					if Debug {
						fmt.Printf("\t\tOpcode Fx0A executed: Wait for a key press - (NOT PRESSED)\n")
					}
				}
				break


			// Fx15 - LD DT, Vx
			// Set delay timer = Vx.
			// DT is set equal to the value of Vx.
			case 0x0015:
				DelayTimer = V[x]
				PC += 2
				if Debug {
					fmt.Printf("\t\tOpcode Fx15 executed: Set delay timer = V[x(%d):%d]\n", x, V[x])
				}
				break

			// Fx18 - LD ST, Vx
			// Set sound timer = Vx.
			// ST is set equal to the value of Vx.
			case 0x0018:
				SoundTimer = V[x]
				PC += 2
				if Debug {
					fmt.Printf("\t\tOpcode Fx18 executed: Set sound timer = V[x(%d)]:%d\n",x, V[x])
				}
				break

			// Fx1E - ADD I, Vx
			// Set I = I + Vx.
			// The values of I and Vx are added, and the results are stored in I.
			// ***
			// Check FX1E (I = I + VX) buffer overflow. If buffer overflow, register
			// VF must be set to 1, otherwise 0. As a result, register VF not set to 1.
			// This undocumented feature of the Chip-8 and used by Spacefight 2091!
			case 0x001E:
				if Debug {
					fmt.Printf("\t\tOpcode Fx1E executed: Add the value of V[x(%d)]:%d to I(%d)\n",x, V[x], I)
				}

				// *** Implement the undocumented feature used by Spacefight 2091
				if FX1E_spacefight2091 {
					if ( I + uint16(V[x]) > 0xFFF ) { //4095 - Buffer overflow
						V[0xF] = 1
						I = ( I + uint16(V[x]) ) - 4095
						fmt.Printf("\n\t\tPAUSE mode ENABLED\n\t\tProposital Pause to map when FX1E fix is used in Spacefight 2091!\n")
						fmt.Printf("\n\t\tPress \"P\" to continue.\n")
						Pause = true	// Put here to try to identify usage in the game
					} else {
						V[0xF] = 0
						I += uint16(V[x])
					}
				// Normal opcode pattern
				} else {
					I += uint16(V[x])
				}

				PC += 2
				break

			// Fx29 - LD F, Vx
			// Set I = location of sprite for digit Vx.
			// The value of I is set to the location for the hexadecimal sprite corresponding to the value of Vx.
			case 0x0029:
				// Load CHIP-8 font. Start from Memory[0]
				I = uint16(V[x]) * 5
				PC += 2
				if Debug {
					fmt.Printf("\t\tOpcode Fx29 executed: Set I(%X) = location of sprite for digit V[x(%d)]:%d (*5)\n", I, x, V[x])
				}
				break

			// SCHIP Fx30 - LD F, Vx
			// Set I = location of sprite for digit Vx.
			// The value of I is set to the location for the hexadecimal sprite corresponding to the value of Vx.
			case 0x0030:
				// Load SCHIP font. Start from Memory[80]
				I = 80 + uint16(V[x]) * 10
				PC += 2
				if Debug {
					fmt.Printf("\t\tSCHIP Opcode Fx30 executed: Set I(%X) = location of sprite for digit V[x(%d)]:%d (*10)\n", I, x, V[x])
				}
				break
			// Fx33 - LD B, Vx
			// BCD - Binary Code hexadecimal
			// Store BCD representation of Vx in memory locations I, I+1, and I+2.
			// set_BCD(Vx);
			// Ex. V[x] = ff (maximum value) = 255
			// memory[i+0] = 2
			// memory[i+1] = 5
			// memory[i+2] = 5
			// % = modulus operator:
			// 3 % 1 would equal zero (since 3 divides evenly by 1)
			// 3 % 2 would equal 1 (since dividing 3 by 2 results in a remainder of 1).
			case 0x0033:
				Memory[I]   = V[x] / 100
				Memory[I+1] = (V[x] / 10) % 10
				Memory[I+2] = (V[x] % 100) % 10
				PC += 2
				if Debug {
					fmt.Printf("\t\tOpcode Fx33 executed: Store BCD representation of V[x(%d)]:%d in memory locations I(%X):%d, I+1(%X):%d, and I+2(%X):%d\n", x, V[x], I, Memory[I], I+1, Memory[I+1], I+2, Memory[I+2])
				}
				break

			// Fx55 - LD [I], Vx
			// Store registers V0 through Vx in memory starting at location I.
			// The interpreter copies the values of registers V0 through Vx into memory, starting at the address in I.
			//
			// Stores V0 to VX (including VX) in memory starting at address I. The offset from I is increased by 1 for each value written, but I itself is left unmodified.[d]
			// In the original CHIP-8 implementation, and also in CHIP-48, I is left incremented after this instruction had been executed. In SCHIP, I is left unmodified.
			case 0x0055:
				for i := uint16(0); i <= x; i++ {
					Memory[I+i] = V[i]
				}
				PC += 2

				// If needed, run the original Chip-8 opcode (not used in recent games)
				if Legacy_Fx55_Fx65 {
					I = I + x + 1
				}

				if Debug {
					fmt.Printf("\t\tOpcode Fx55 executed: Registers V[0] through V[x(%d)] in memory starting at location I(%d)\n",x, I)
				}
				break

			// Fx65 - LD Vx, [I]
			// Read registers V0 through Vx from memory starting at location I.
			// The interpreter reads values from memory starting at location I into registers V0 through Vx.
			//// I is set to I + X + 1 after operation²
			//// ² Erik Bryntse’s S-CHIP documentation incorrectly implies this instruction does not modify
			//// the I register. Certain S-CHIP-compatible emulators may implement this instruction in this manner.
			//// MAYBE NEED TO IMPLEMENT NO S-CHIP8 ***
			case 0x0065:

				for i := uint16(0); i <= x; i++ {
					V[i] = Memory[I+i]
					//fmt.Printf("   New value of V[%d] = %d\n", i, V[i])
				}

				// Fix I implementation TEST SUGGESTED
				//fmt.Println(I)
				//fmt.Println(x)
				// I = I + x + 1
				//fmt.Println(I)
				//fmt.Println(x)
				// Increment Program Counter
				PC += 2

				// If needed, run the original Chip-8 opcode (not used in recent games)
				if Legacy_Fx55_Fx65 {
					I = I + x + 1
				}

				if Debug {
					fmt.Printf("\t\tOpcode Fx65 executed: Read registers V[0] through V[x(%d)] from memory starting at location I(%X)\n",x, I)
				}
				break

			// SCHIP FX75
			// Store V0 through VX to HP-48 RPL user flags (X <= 7).
			case 0x0075:

				// Temporary, to check
				if x >= 8 {
					fmt.Printf("FX75 X VALUE CONTROL!!!")
					os.Exit(2)
				}

				for i := 0; i <= int(x); i++ {
					RPL[i] = V[i]
				}

				PC += 2
				if Debug {
					fmt.Printf("\t\tSCHIP - Opcode Fx75 executed: Read RPL user flags from 0 to %d and store in V[0] through V[x(%d)]\n",x,x)
				}

				break

			// SCHIP FX85
			// Read V0 through VX to HP-48 RPL user flags (X <= 7).
			case 0x0085:

				// Temporary, to check
				if x >= 8 {
					fmt.Printf("FX85 X VALUE CONTROL!!!")
					os.Exit(2)
				}

				for i := 0; i <= int(x); i++ {
					V[i] = RPL[i]
				}

				PC += 2
				if Debug {
					fmt.Printf("\t\tSCHIP - Opcode Fx85 executed: Read registers V[0] through V[x(%d)] and store in RPL user flags\n",x)
				}

				break
			default:
				fmt.Printf("\t\tOpcode Family F000 - Not mapped opcode: 0x%X\n", Opcode)
				os.Exit(2)

			}
			break



		default:
			fmt.Printf("\t\tOPCODE FAMILY %X NOT IMPLEMENTED!\n", Opcode & 0xF000)
			os.Exit(3)
	}


	// Independent of CPU CLOCK, Sound and Delay Timers runs at 60Hz
	select {
		case <-TimerClock.C:
			// When ticker run (60 times in a second, check de DelayTimer)
			if DelayTimer > 0 {
				DelayTimer--
			}

			// When ticker run (60 times in a second, check de SoundTimer)
			if SoundTimer > 0 {
				if SoundTimer == 1 {
					go Sound.PlaySound(Sound.Beep_buffer)
				}
				SoundTimer--
			}

		default:
			// No timer to handle
	}

	Cycle ++

	// Debug time execution - Opcode Handling
	if Debug {
		elapsed := time.Since(start)
		fmt.Printf("\t\tTime track - Opcode took: %s\n\n", elapsed)
	}

}
