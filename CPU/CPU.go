package CPU

import (
	"os"
	"fmt"
	"time"
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
	OpcMessage		string			// Debug Message returned after Opcodes

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

			switch Opcode & 0x0F00 {

				case 0x0000: //00NN

					switch Opcode & 0x00FF { //00NN

						// 0000 (ETI-660)
						case 0x0000:
							opc_chip8_ETI660_0000()
							break

						// 00E0 (CHIP-8)
						case 0x00E0:
							opc_chip8_00E0()
							break

						// 00EE (CHIP-8)
						case 0x00EE:
							opc_chip8_00EE()
							break

						// 00F0 (CHIP-8x HiRes)
						case 0x00F0:
							opc_chip8HiRes_00F0()
							break

						// 00F8 (ETI-660)
						case 0x00F8:
							opc_chip8_ETI660_00F8()
							break

						// 00FB (SCHIP)
						case 0x00FB:
							opc_schip_00FB()
							break

						// 00FC (SCHIP)
						// 00FC (ETI-660) - Turn display off
						case 0x00FC:
							// ETI-660 Opcode
							if Global.Hybrid_ETI_660_HW {
								opc_chip8_ETI660_00FC()
								break

							// SCHIP Opcode
							} else {
								opc_schip_00FC()
								break
							}

						// 00FD (SCHIP)
						case 0x00FD:
							opc_schip_00FD()
							break

						// 00FE (SCHIP)
						case 0x00FE:
							opc_schip_00FE()
							break

						// 00FF (SCHIP)
						// 00FF - In ETI-660, 00FF is a NO OP (do nothing)
						case 0x00FF:
							// 00FF - ETI-660
							if Global.Hybrid_ETI_660_HW {
								opc_chip8_ETI660_00FF()
								break

							// 00FF - SCHIP
							} else {
								opc_schip_00FF()
								break
							}
					}

					switch Opcode & 0x00F0 { //00N0
						// 00CN (SCHIP)
						case 0x00C0:
							n := Opcode & 0x000F
							opc_schip_00CN(n)
							break
					}

				case 0x0200: //02NN

					switch Opcode & 0x0FFF {

						// 0230 (CHIP-8 HIRES)
						case 0x0230:
							opc_chip8HiRes_0230()
							break

						// 02D8 (CHIP-8 NON DOCUMENTED)
						case 0x02D8:
							opc_chip8_ND_02D8()
							break

						default:
							fmt.Printf("\t\tOpcode 0x%04X NOT IMPLEMENTED!!!!\n", Opcode)
							os.Exit(0)
					}

				default:
					fmt.Printf("\t\tOpcode 0x%04X NOT IMPLEMENTED!!!!\n", Opcode)
					os.Exit(0)
			}

		// ---------------------------- CHIP-8 1xxx instruction set ---------------------------- //
		// 1nnn (CHIP-8)
		case 0x1000:

			// 1nnn (CHIP-8 HIRES)
			// If PC=0x200 AND Opcode=0x1260, update Opcode to 0x12C0 (Jump to address 0x2c0)
			if PC == 0x200 && Opcode == 0x1260 {
				opc_chip8HiRes_1NNN()
				break

			// 1nnn (CHIP-8)
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
					fmt.Printf("\t\tOpcode 0x%04X NOT IMPLEMENTED!!!!\n", Opcode)
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

			switch Opcode & 0x000F {

				// DXY0 (SCHIP)
				case 0x0000:

					if !SCHIP {
						// Quirk to SCHIP Robot DEMO)
						// Even in SCHIP Mode this game needs to draw 16x16 Pixels
						if DXY0_loresWideSpriteQuirks {
							SCHIP_LORES = false
						} else {
							// If NOT in SCHIP mode will draw 16x8 sprites
							SCHIP_LORES = true
						}
					}
					// If in SCHIP mode will draw 16x16 sprites
					opc_schip_DXY0(Opcode)

				// DXYN (CHIP-8, Draw n-byte sprites)
				default:
					opc_chip8_DXYN(Opcode)
					break
			}

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
					fmt.Printf("\t\tOpcode 0x%04X NOT IMPLEMENTED!!!!\n", Opcode)
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

				// Fx0A (CHIP-8)
				case 0x000A:
					opc_chip8_FX0A(x)
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

				// Fx55 (CHIP-8)
				case 0x0055:
					opc_chip8_FX55(x)
					break

				// Fx65 (CHIP-8)
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

				// Fx00 (ETI-660)
				case 0x0000:
					opc_chip8_ETI660_FX00(x)
					break

				default:
					fmt.Printf("\t\tOpcode 0x%04X NOT IMPLEMENTED!!!!\n", Opcode)
					os.Exit(0)
				}
				break


		// End of main switch (0xN000)
		default:
			fmt.Printf("\t\tOpcode Family %X not implemented!\n", Opcode & 0xF000)
			os.Exit(0)
	}


	Cycle ++	// Increment overall cycle Counter
	CyclesCounter ++	// Increment cycle counter to measure with FPS

	// Debug time execution - Opcode Handling
	if Debug {
		elapsed := time.Since(start)
		fmt.Printf("\t\tTime track - Opcode took: %s\n\n", elapsed)
	}

}
