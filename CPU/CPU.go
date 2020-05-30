package CPU

import (
	"os"
	"fmt"
	"time"
	"strconv"
	"Chip8/Fontset"
	"Chip8/Global"
)

const (
	KeyArraySize	byte		= 16		// Control the number of Keys mapped in Key Array
	CPU_Clock_increase_rate		= 100	// CPU Clock increase rate
	CPU_Clock_decrease_rate		= 100	// CPU Clock decrease rate
)


var (
	// Components
	Memory			[4096]byte		// Memory
	PC			uint16			// Program Counter
	Opcode			uint16			// CPU Operation Code
	Stack			[16]uint16		// Stack
	SP			uint16			// Stack Pointer
	V			[16]byte		// V Register
	I			uint16			// I Register
	DelayTimer		byte			// Delay Timer
	SoundTimer		byte			// Sound Timer
	Graphics		[128 * 64]byte		// Graphic Array
	// ETI-660 HW
	P			byte			// Pitch (frequency) value of the tone generator (beeper)

	// Key Control
	Key			[KeyArraySize]byte	// Control the Keys Pressed

	// Timers
	FPS			*time.Ticker		// Frames per second
	FPSCounter		*time.Ticker		// Count frames per second
	TimersClock		*time.Ticker		// Delay and Sound Timer
	KeyboardClock		*time.Ticker		// Keyboard Timer to be used with emulator keys
	CPU_Clock		*time.Ticker		// CPU Clock
	CPU_Clock_Speed		time.Duration		// Value defined to CPU Clock
	CPU_Clock_Speed_Max	time.Duration		// Max value of CPU Clock
	SCHIP_TimerClockHack	*time.Ticker		// SCHIP used to decrease DT faster than 60HZ to gain speed
	MessagesClock		*time.Ticker		// Clock used to display messages on screen

	// General Variables and flags
	MemoryCleanSnapshot	[4096]byte		// Some games like Single Dragon changes memory, so to reset its necessary to reload game

	// SCHIP Specific Variables
	SCHIP			bool			// SCHIP MODE ENABLED OR DISABLED
	SCHIP_LORES		bool			// SCHIP in Low Resolution mode (00FE)
	SCHIP_TimerHack		bool			// Enable or disable SCHIP DelayTimer Hack
	RPL			[8]byte			// HP-48 RPL user flags

	// Counters
	Cycle			uint16			// CPU Cycle Counter
	DrawFlagCounter		uint16			// Graphical opcodes Counter per second
	CyclesCounter		uint16			// CPU Cycle Counter per second

	// DEBUG
	Pause			bool			// Pause (Used to Forward and Rewind CPU Cycles)
	Debug			bool			// DEBUG mode
	Debug_L2		bool			// DEBUG Rewind Mode

)


// Initialization
func Initialize() {
	// Components
	Memory			= [4096]byte{}
	if Global.Hybrid_ETI_660_HW {
		PC		= 0x600	// start at 0x600 for ETI-600 HW (Hybrid)
	} else {
		PC		= 0x200	// start at 0x200 (default CHIP-8)
	}
	Opcode			= 0
	Stack			= [16]uint16{}
	SP			= 0
	V			= [16]byte{}
	I			= 0
	DelayTimer		= 0
	SoundTimer		= 0
	Graphics		= [128 * 64]byte{}

	// Timers
	FPS			= time.NewTicker(time.Second / 60)			// FPS Clock
	FPSCounter		= time.NewTicker(time.Second)				// FPS Counter Clock
	CPU_Clock_Speed		= 500							// Initial CPU Clock Speed: CHIP-8=500, SCHIP=2000
	CPU_Clock_Speed_Max	= 10000
	CPU_Clock		= time.NewTicker(time.Second / CPU_Clock_Speed)
	SCHIP_TimerClockHack	= time.NewTicker(time.Second / (CPU_Clock_Speed * 10) )
	KeyboardClock		= time.NewTicker(time.Second / 60)
	TimersClock		= time.NewTicker(time.Second / 60)			// Decrease SoundTimer and DelayTimer
	MessagesClock		= time.NewTicker(time.Second * 5)			// Clock used to display messages on screen

	// General Variables and flags
	Cycle			= 0
	DrawFlagCounter		= 0
	CyclesCounter		= 0
	// Keys
	Key			= [KeyArraySize]byte{}
	// Graphics
	Global.SizeX		= 64
	Global.SizeY		= 32

	// ETI-660 Graphics
	// Update screen size if in ETI-660 HW mode
	if Global.Hybrid_ETI_660_HW {
		Global.SizeX		= 64
		Global.SizeY		= 48
	}

	// SCHIP Specific Variables
	SCHIP			= false
	SCHIP_LORES		= false

	// Load CHIP-8 8x5 fontset (Memory address 0-79)
	for i := 0; i < len(Fontset.Chip8Fontset); i++ {
		Memory[i] = Fontset.Chip8Fontset[i]
	}

	// Load SCHIP 8x10 fontset (Memory address 80-240)
	for i := 0; i < len(Fontset.SCHIPFontset); i++ {
		Memory[i+80] = Fontset.SCHIPFontset[i]
	}
}


func Show() {
	fmt.Printf("Cycle: %d\tOpcode: %04X(%04X)\tPC: %d(0x%X)\tSP: %d\tStack: %d\tV: %d\tI: %d\tDT: %d\tST: %d\tKey: %d\n", Cycle, Opcode, Opcode & 0xF000, PC, PC,  SP, Stack, V, I, DelayTimer, SoundTimer, Key)
}


// SCHIP HI-RES MODE
// If in SCHIP mode will draw 16x16 sprites
func DXY0_SCHIP_HiRes(x, y, n, byte, gpx_position uint16) {

	// Turn n in 16 (pixel size in SCHIP Mode)
	n = 16
	if Debug {
		fmt.Printf("\t\tSCHIP - Opcode Dxy0 HI-RES MODE (%X) DRAW GRAPHICS! - Address I: %d Position V[x(%d)]: %d V[y(%d)]: %d\n" , Opcode, I, x, V[x], y, V[y])
	}

	// Print N Bytes from address I in V[x]V[y] position of the screen
	for byte = 0 ; byte < n ; byte++ {

		var (
			binary	string = ""
			sprite	uint8  = 0
			sprite2	uint8  = 0
		)

		// DOCUMENT SPRITES
		sprite  = Memory[I + (byte * 2)]
		sprite2 = Memory[I + (byte * 2) + 1]

		// Sprite in binary format
		binary = fmt.Sprintf("%.8b%.8b", sprite,sprite2)

		// Always print 8 bits
		for bit := 0; bit < 16 ; bit++ {

			// Convert the binary[bit] variable into an INT using Atoi method
			bit_binary, err := strconv.Atoi(fmt.Sprintf("%c", binary[bit]))
			if err == nil {

			}

			// Set the index to write the 8 bits of each pixel
			gfx_index := uint16(gpx_position) + uint16(bit) + (byte*uint16(Global.SizeX))

			// If tryes to draw bits outside the vector size, ignore
			if ( gfx_index >= uint16(Global.SizeX) * uint16(Global.SizeY) ) {
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

		}
	}

}


// SCHIP LOW-RES MODE
// If NOT in SCHIP mode will draw 16x8 sprites
func DXY0_SCHIP_LoRes(x, y, n, byte, gpx_position uint16) {

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

		// DOCUMENT SPRITES
		sprite = Memory[I + byte]

		// Sprite in binary format
		binary = fmt.Sprintf("%.8b", sprite)

		// Always print 8 bits
		for bit := 0; bit < 8 ; bit++ {

			// Convert the binary[bit] variable into an INT using Atoi method
			bit_binary, err := strconv.Atoi(fmt.Sprintf("%c", binary[bit]))
			if err == nil {

			}

			// Set the index to write the 8 bits of each pixel
			gfx_index := uint16(gpx_position) + uint16(bit) + (byte*uint16(Global.SizeX))

			// If tryes to draw bits outside the vector size, ignore
			if ( gfx_index >= uint16(Global.SizeX) * uint16(Global.SizeY) ) {
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

		}
	}

}


// Dxyn - DRW Vx, Vy, nibble
// Display n-byte sprite starting at memory location I at (Vx, Vy), set VF = collision.
func DXYN_CHIP8(x, y, n, byte, gpx_position uint16) {
	// Draw in Chip-8 Low Resolution mode

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
		sprite = Memory[I + byte]

		// Sprite in binary format
		binary = fmt.Sprintf("%.8b", sprite)

		// Always print 8 bits
		for bit := 0; bit < 8 ; bit++ {

			// Convert the binary[bit] variable into an INT using Atoi method
			bit_binary, err := strconv.Atoi(fmt.Sprintf("%c", binary[bit]))
			if err == nil {

			}

			// Set the index to write the 8 bits of each pixel
			gfx_index := uint16(gpx_position) + uint16(bit) + (byte*uint16(Global.SizeX))

			// If tryes to draw bits outside the vector size, ignore
			if ( gfx_index >= uint16(Global.SizeX) * uint16(Global.SizeY) ) {
				//fmt.Printf("Bigger than 2048 or 8192\n")
				continue
			}

			// If bit=1, test current graphics[index], if is already set, mark v[F]=1 (collision)
			if (bit_binary  == 1){
				// Set colision case graphics[index] is already 1
				if (Graphics[gfx_index] == 1){
					V[0xF] = 1
				}
				// After, XOR the graphics[index] (DRAW)
				Graphics[gfx_index] ^= 1
			}

		}
	}

}


// CPU Interpreter
func Interpreter() {

	// Reset Flag every cycle
	Global.DrawFlag = false

	// Read the Opcode from PC and PC+1 bytes
	Opcode = uint16(Memory[PC])<<8 | uint16(Memory[PC+1])

	// Print Cycle and Debug Information
	if Debug {
		Show()
	}

	// Enable tracking to Rewind function
	if Rewind_mode {
		rewind()
	}

	// Start timer to measure procedures inside Interpreter
	start := time.Now()

	// Map Opcode Family
	switch Opcode & 0xF000 {

		// ---------------------------- CHIP-8 0xxx instruction set ---------------------------- //
		case 0x0000: //0NNN

			x := Opcode & 0x000F

			switch Opcode & 0x00F0 {

				case 0x00E0:

					// 00E0 (CHIP-8)
					if x == 0x0000 {
						opc_chip8_00E0()
						break
					}

					// 00EE (CHIP-8)
					if x == 0x000E {
						opc_chip8_00EE()
						break
					}

				// 02D8 (CHIP-8 NON DOCUMENTED)
				case 0x00D0:
					opc_chip8_ND_02D8()
					break

					// 00FF (SCHIP)
					// 00FF - In ETI-660, 00FF is a NO OP (do nothing)
				case 0x00F0:
					if x == 0x000F {
						// ETI-660 Do Nothing
						if Global.Hybrid_ETI_660_HW {
							PC += 2
							if Debug {
								fmt.Printf("\t\tHybrid ETI-660 - Opcode 00FF executed: No Operation (do nothing)\tPC+=2\n")
							}
							break

						// Enable SCHIP Mode
						} else {
							opc_schip_00FF()
							break
						}

					// 00FE (SCHIP)
					} else if x == 0x000E {
						opc_schip_00FE()

					// 00FD (SCHIP)
					} else if x == 0x000D {
						opc_schip_00FD()

					// 00FC (SCHIP)
					// 00FC (ETI-660) - Turn display off
					} else if x == 0x000C {
						// ETI-660 Opcode
						if Global.Hybrid_ETI_660_HW {
							PC+=2
							if Debug {
								fmt.Printf("\t\tHybrid ETI-660 - Opcode 00FC executed: Turn display off (Do nothing)\t PC+=2\n")
								fmt.Printf("\n\nPROPOSITAL EXIT TO MAP 00FC USAGE!!!!\n\n")
								os.Exit(2)
							}

						// SCHIP Opcode
						} else {
							opc_schip_00FC()
						}

					// 00FB (SCHIP)
					} else if x == 0x000B {
						opc_schip_00FB()


					// ETI-660 0x00F8
					// Turn display on
					} else if x == 0x0008 {
						PC += 2
						if Debug {
							fmt.Printf("\t\tHybrid ETI-660 - Opcode 00F8 executed: Turn display on (Do nothing)\t PC+=2\n")
							fmt.Printf("\n\nPROPOSITAL EXIT TO MAP 00F8 USAGE!!!!\n\n")
							os.Exit(2)
						}


					// Two-page display for CHIP-8X (Extension of CHIP-8x) 0x00F0
					// 00F0: Return from subroutine (replaces 00EE)
					// Also used in some Hybrid ETI-660 programs like "Music Maker"
					} else if x == 0x0000 {
						PC = Stack[SP] + 2
						SP --
						if Debug {
							fmt.Printf("\t\tCHIP-8X Two-page display (Extension) - Opcode 00F0 executed. - Return from subroutine (replaces 00EE in CHIP-8x).\n")
						}


					} else {
						fmt.Printf("\t\tOpcode %04X NOT IMPLEMENTED.\n", Opcode)
						os.Exit(2)
					}


				// 00CN (SCHIP)
				case 0x00C0:
					opc_schip_00CN(x)
					break

				// CHIP-8 HIRES - 0230
				// Clear screen used by Hi Resolution Chip8
				case 0x0030:

					// Clear display
					Graphics = [128 * 64]byte{}

					// Set CHIP-8 HIRES Resolution
					Global.SizeX = 64
					Global.SizeY = 64

					PC += 2
					if Debug {
						fmt.Println("\t\tHIRES - Opcode 0230 executed. - Clear the display\n")
					}
					break

				// ETI-660 - 0x0000
				// Return to monitor (exit interpreter)
				case 0x0000:
					if Debug {
						fmt.Printf("\t\tHybrid ETI-660 - Opcode 0000 executed: Return to monitor (exit interpreter)\n")
					}
					break

				default:
					if Debug {
						fmt.Printf("\t\tOpcode 0x%04X NOT IMPLEMENTED!!!!\n", Opcode)
					}
					os.Exit(0)
			}


		// ---------------------------- CHIP-8 1xxx instruction set ---------------------------- //
		// 1nnn (CHIP-8)
		case 0x1000:

			// HI-RES CHIP-8 EMULATION
			// If PC=0x200 AND Opcode=0x1260, update Opcode to 0x12C0 (Jump to address 0x2c0)
			// Need to add Opcode 0x0230 to handle the clearscreen event for 64x64 hires
			if PC == 0x200 && Opcode == 0x1260 {
				// Execute the operation
				PC = 0x2C0

				// After show the execution time
				if Debug {
					fmt.Printf("\t\tHIRES - Opcode 1260 WITH PC=0x200. Init 64x64 Chip8 hires mode. Opcode=0x12C0, jump to address 0x2c0 -> (PC=0x2c0)\n")
				}
				break

			// Or start the regular code from 1nnn
			} else {
				opc_chip8_1NNN()
				break
			}


		// ---------------------------- CHIP-8 2xxx instruction set ---------------------------- //
		// 2nnn (CHIP-8)
		case 0x2000:
			opc_chip8_2NNN()
			break

		// ---------------------------- CHIP-8 3xxx instruction set ---------------------------- //
		// 3xnn (CHIP-8)
		case 0x3000:
			opc_chip8_3XNN()
			break

		// ---------------------------- CHIP-8 4xxx instruction set ---------------------------- //
		// 4xnn (CHIP-8)
		case 0x4000:
			opc_chip8_4XNN()
			break

		// ---------------------------- CHIP-8 5xxx instruction set ---------------------------- //
		// 5xy0 (CHIP-8)
		case 0x5000:
			opc_chip8_5XY0()
			break

		// ---------------------------- CHIP-8 6xxx instruction set ---------------------------- //
		// 6xnn (CHIP-8)
		case 0x6000:
			opc_chip8_6XNN()
			break

		// ---------------------------- CHIP-8 7xxx instruction set ---------------------------- //
		// 7xnn (CHIP-8)
		case 0x7000:
			opc_chip8_7XNN()
			break

		// ---------------------------- CHIP-8 8xxx instruction set ---------------------------- //
		// 0x8000 instruction set
		case 0x8000:
			x := (Opcode & 0x0F00) >> 8
			y := (Opcode & 0x00F0) >> 4

			switch Opcode & 0x000F {

			// 8xy0 (CHIP-8)
			case 0x0000:
				opc_chip8_8XY0(x, y)
				break

			// 8xy1 (CHIP-8)
			case 0x0001:
				opc_chip8_8XY1(x, y)
				break

			// 8xy2 (CHIP-8)
			case 0x0002:
				opc_chip8_8XY2(x, y)
				break

			// 8xy3 (CHIP-8)
			case 0x0003:
				opc_chip8_8XY3(x, y)
				break

			// 8xy4 (CHIP-8)
			case 0x0004:
				opc_chip8_8XY4(x, y)
				break

			// 8xy5 (CHIP-8)
			case 0x0005:
				opc_chip8_8XY5(x, y)
				break

			// 8xy6 (CHIP-8)
			case 0x0006:
				opc_chip8_8XY6(x, y)
				break


			// 8xy7 (CHIP-8)
			case 0x0007:
				opc_chip8_8XY7(x, y)
				break

			// 8xyE (CHIP-8)
			case 0x000E:
				opc_chip8_8XYE(x, y)
				break

			default:
				if Debug {
					fmt.Printf("\t\tOpcode 0x8000 NOT IMPLEMENTED!!!!\n")
				}
				os.Exit(0)
			}

		// ---------------------------- CHIP-8 9xxx instruction set ---------------------------- //
		// 9xy0 (CHIP-8)
		case 0x9000:
			opc_chip8_9XY0()
			break

		// ---------------------------- CHIP-8 Axxx instruction set ---------------------------- //
		// Annn (CHIP-8)
		case 0xA000:
			opc_chip8_ANNN()
			break

		// ---------------------------- CHIP-8 Bxxx instruction set ---------------------------- //
		// Bnnn (CHIP-8)
		case 0xB000:
			opc_chip8_BNNN()
			break

		// ---------------------------- CHIP-8 Cxxx instruction set ---------------------------- //
		// Cxnn (CHIP-8)
		case 0xC000:
			opc_chip8_CXNN()
			break

		// ---------------------------- CHIP-8 Dxxx instruction set ---------------------------- //
		case 0xD000:

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
				if V[x] + uint8(n) > (uint8(Global.SizeX) +1)  {
					PC += 2
					break
				}
			}

			// Translate the x and Y to the Graphics Vector
			gpx_position = (uint16(V[x]) + (uint16(Global.SizeX) * uint16(V[y])))

			// SCHIP Dxy0
			// When in high res mode show a 16x16 sprite at (VX, VY)
			// If N=0, Draw in SCHIP High Resolution mode
			if n == 0 {
				// SCHIP HI-RES MODE
				// If in SCHIP mode will draw 16x16 sprites
				if SCHIP {
					DXY0_SCHIP_HiRes(x, y, n, byte, gpx_position)
				// SCHIP LOW-RES MODE
				// If NOT in SCHIP mode will draw 16x8 sprites
				} else {
					// Quirk to SCHIP Robot DEM)
					// Even in SCHIP Mode this game needs to draw 16x16 Pixels
					if DXY0_loresWideSpriteQuirks {
						DXY0_SCHIP_HiRes(x, y, n, byte, gpx_position)
					} else {
						DXY0_SCHIP_LoRes(x, y, n, byte, gpx_position)
					}
				}
			// Dxyn - DRW Vx, Vy, nibble
			// Display n-byte sprite starting at memory location I at (Vx, Vy), set VF = collision.
			} else {
				DXYN_CHIP8(x, y, n, byte, gpx_position)
			}

			PC += 2
			Global.DrawFlag = true
			DrawFlagCounter ++


		// ---------------------------- CHIP-8 Exxx instruction set ---------------------------- //
		// 0xE000 instruction set
		case 0xE000:

			x := (Opcode & 0x0F00) >> 8
			switch Opcode & 0x00FF {

			// Ex9E (CHIP-8)
			case 0x009E:
				opc_chip8_EX9E(x)
				break

			// ExA1 (CHIP-8)
			case 0x00A1:
				opc_chip8_EXA1(x)
				break
			default:
				fmt.Printf("Opcode Family E000 - Not mapped opcote: E000\n")
				os.Exit(0)
			}


		// ---------------------------- CHIP-8 Fxxx instruction set ---------------------------- //
		case 0xF000:

			x := (Opcode & 0x0F00) >> 8

			switch Opcode & 0x00FF {

			// Fx07 (CHIP-8)
			case 0x0007:
				opc_chip8_FX07(x)
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

			// Fx15 (CHIP-8)
			case 0x0015:
				opc_chip8_FX15(x)
				break

			// Fx18 (CHIP-8)
			case 0x0018:
				opc_chip8_FX18(x)
				break

			// Fx1E (CHIP-8)
			case 0x001E:
				opc_chip8_FX1E(x)
				break

			// Fx29 (CHIP-8)
			case 0x0029:
				opc_chip8_FX29(x)
				break

			// Fx30 (SCHIP)
			case 0x0030:
				opc_schip_FX30(x)
				break

			// Fx33 (CHIP-8)
			case 0x0033:
				opc_chip8_FX33(x)
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
			case 0x0065:
				opc_chip8_FX65(x)
				break

			// FX75 (SCHIP)
			case 0x0075:
				opc_schip_FX75(x)
				break

			// FX85 (SCHIP)
			case 0x0085:
				opc_schip_FX85(x)
				break

			// CHIP-8 ETI-660 Hybrid - Fx00
			// Set the pitch (frequency) of the tone generator (beeper) to Vx
			case 0x0000:
				P = V[x]	// NOT USED YET!!! Need to implement sound library to handle it
				PC +=2

				if Debug {
					fmt.Printf("\t\tHybrid ETI-660 - Opcode Fx00 executed: Set the pitch (frequency) of the tone generator (beeper) to value of V[%d]\t\tP=%d\n", x, V[x])
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


	Cycle ++	// Increment overall cycle Counter
	CyclesCounter ++	// Increment cycle counter to measure with FPS

	// Debug time execution - Opcode Handling
	if Debug {
		elapsed := time.Since(start)
		fmt.Printf("\t\tTime track - Opcode took: %s\n\n", elapsed)
	}

}
