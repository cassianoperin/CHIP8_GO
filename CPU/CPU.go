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

// Components
var (
	Memory		[4096]byte // Memory
	PC			uint16     // Program Counter
	Opcode		uint16     // CPU Operation Code
	Stack		[16]uint16 // Stack
	SP			uint16     // Stack Pointer
	V			[16]byte
	I			uint16
	DelayTimer	byte
	SoundTimer	byte
	TimerClock	*time.Ticker
	Graphics		[64 * 32]byte

	// True if the screen must be drawn
	DrawFlag		bool
	Key			[16]byte
	Cycle		uint16

	// Control the Keys Pressed
	KeyPressed = map[uint16]pixelgl.Button{
		0:	pixelgl.Key0,
		1:	pixelgl.Key1,
		2:	pixelgl.Key4,
		3:	pixelgl.Key3,
		4:	pixelgl.Key2,
		5:	pixelgl.Key5,
		6:	pixelgl.Key6,
		7:	pixelgl.Key7,
		8:	pixelgl.Key8,
		9:	pixelgl.Key9,
		10:	pixelgl.KeyLeft,
		11:	pixelgl.KeyRight,
		12:	pixelgl.KeyUp,
		13:	pixelgl.KeyDown,
		14:	pixelgl.KeyQ,
		15:	pixelgl.KeyW,
	}

	// Beep sound file
	sound_file string

)


func Initialize() {
	// Initialization
	Memory		= [4096]byte{}
	PC			= 0x200
	//CPU.PC		=	0x2FE
	Opcode		= 0
	Stack		= [16]uint16{}
	SP			= 0
	V			= [16]byte{}
	I 			= 0
	Graphics		= [64 * 32]byte{}
	DrawFlag		= false
	DelayTimer	= 0
	SoundTimer	= 0
	// Create a ticker at 60Hz
	TimerClock	= time.NewTicker(time.Second / 60)
	// Load fontset
	for i := 0; i < len(Fontset.Chip8Fontset); i++ {
		Memory[i] = Fontset.Chip8Fontset[i]
	}
	Key			= [16]byte{}
	Cycle		= 0

}

func Debug() {

	fmt.Printf("Cycle: %d\tOpcode: %04X(%04X)\tPC: %d(0x%X)\tSP: %d\tStack: %d\tV: %d\tI: %d\tDT: %d\tST: %d\tKey: %d\n", Cycle, Opcode, Opcode & 0xF000, PC, PC,  SP, Stack, V, I, DelayTimer, SoundTimer, Key)
}


// CPU Interpreter
func Interpreter() {

	// Reset Flag every cycle
	DrawFlag = false

	// Read the Opcode from PC and PC+1 bytes
   	Opcode = uint16(Memory[PC])<<8 | uint16(Memory[PC+1])

	// Print Cycle and Debug Information
	Debug()

	// Map Opcode Family
	switch Opcode & 0xF000 {

		// ############################ 0x0000 instruction set ############################
	   	case 0x0000: //0NNN

	   		switch Opcode & 0x00FF {

			// 00E0 - CLS
			// Clear the display.
	   		case 0x00E0:
	   			// Clear display
	   			Graphics = [64 * 32]byte{}
	   			PC += 2
				fmt.Println("\t\tOpcode 00E0 executed. - Clear the display\n\n")
				break

			// 00EE - RET
			// Return from a subroutine
			// The interpreter sets the program counter to the address at the top of the stack, then subtracts 1 from the stack pointer.
			// MUST MOVE TO NEXT ADDRESS AFTER THIS (PC+=2)
	   		case 0x00EE:
	   			//fmt.Println("   RCA 1802 Opcode 0x00EE - Return")
	   			PC = Stack[SP] + 2
	   			SP --
				fmt.Printf("\t\tOpcode 00EE executed. - Return from a subroutine (PC=%d)\n\n", PC)
				break

	   		default:
	   			fmt.Printf("\t\tOpcode 0x0000 NOT IMPLEMENTED!!!!\n\n", Opcode)
	   			os.Exit(2)
	   		}


		// ############################ 0x1000 instruction set ############################
		// 1nnn - JP addr
		// Jump to location nnn.
		// The interpreter sets the program counter to nnn.
		case 0x1000:
			PC = Opcode & 0x0FFF
			fmt.Printf("\t\tOpcode 1nnn executed: Jump to location 0x%d\n\n", Opcode & 0x0FFF)
			break


		// ############################ 0x2000 instruction set ############################
		// 2nnn - CALL addr
		// Call subroutine at nnn.
		// The interpreter increments the stack pointer, then puts the current PC on the top of the stack. The PC is then set to nnn.
		case 0x2000:
			SP++
			Stack[SP] = PC
			PC = Opcode & 0x0FFF
			fmt.Printf("\t\tOpcode 2nnn executed: Call Subroutine at 0x%d\n\n", PC)
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
				fmt.Printf("\t\tOpcode 3xk executed: V[x(%d)]:(%d) = kk(%d), skip one instruction.\n\n", x, V[x], kk)
			} else {
				PC += 2
				fmt.Printf("\t\tOpcode 3xk executed: V[x(%d)]:(%d) != kk(%d), do NOT skip one instruction.\n\n", x, V[x], kk)
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
				fmt.Printf("\t\tOpcode 4xkk executed: V[x(%d)]:%d != kk(%d), skip one instruction\n\n", x, V[x], kk)
			} else {
				fmt.Printf("\t\tOpcode 4xkk executed: V[x(%d)]:%d = kk(%d), NOT skip one instruction\n\n", x, V[x], kk)
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
				fmt.Printf("\t\tOpcode 5xy0 executed: V[x(%d)]:%d EQUAL V[y(%d)]:%d, SKIP one instruction\n\n", x, V[x], y, V[y])
			} else {
				PC += 2
				fmt.Printf("\t\tOpcode 5xy0 executed: V[x(%d)]:%d NOT EQUAL V[y(%d)]:%d, DO NOT SKIP one instruction\n\n", x, V[x], y, V[y])
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
			fmt.Printf("\t\tOpcode 6xkk executed: Set V[x(%d)] = %d\n\n", x, kk)
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
			fmt.Printf("\t\tOpcode 7xkk executed: Add the value kk(%d) to V[x(%d)]\n\n", kk, x)
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
				fmt.Printf("\t\tOpcode 8xy0 executed: Set V[x(%d)] = V[y(%d)]:%d\n\n", x, y, V[y])
				break

			// Set Vx = Vx OR Vy.
			// Performs a bitwise OR on the values of Vx and Vy, then stores the result in Vx. A bitwise OR compares the corrseponding bits from two values,
			// and if either bit is 1, then the same bit in the result is also 1. Otherwise, it is 0.
			case 0x0001:
				V[x] |= V[y]
				PC += 2
				fmt.Printf("\t\tOpcode 8xy1 executed: Set V[x(%d)]:%d OR V[y(%d)]:%d\n\n", x, V[x], y, V[y])
				break

			// 8xy2 - AND Vx, Vy
			// Set Vx = Vx AND Vy.
			// Performs a bitwise AND on the values of Vx and Vy, then stores the result in Vx. A bitwise AND compares the corrseponding bits from two values, and if both bits are 1, then the same bit in the result is also 1. Otherwise, it is 0.
			case 0x0002:
				V[x] &= V[y]
				PC += 2
				fmt.Printf("\t\tOpcode 8xy2 executed: Set V[x(%d)] = V[x(%d)] AND V[y(%d)]\n\n", x, x, y)
				break

			// 8xy3 - XOR Vx, Vy
			// Set Vx = Vx XOR Vy.
			// Performs a bitwise exclusive OR on the values of Vx and Vy, then stores the result in Vx. An exclusive OR compares the corrseponding bits from two values,
			// and if the bits are not both the same, then the corresponding bit in the result is set to 1. Otherwise, it is 0.
			case 0x0003:
				fmt.Printf("\t\tOpcode 8xy3 executed:  V[x(%d)]:%d XOR V[y(%d)]:%d \n\n", x, V[x], y, V[y])
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
				fmt.Printf("\t\tOpcode 8xy4 executed: Set V[x(%d)] = V[x(%d)] + V[y(%d)]\n\n", x, x, y)

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
				fmt.Printf("\t\tOpcode 8xy5 executed: Set V[x(%d)] = V[x(%d)]:%d - V[y(%d)]:%d\n\n", x, x, V[x], y, V[y])
				break

			// 8xy6 - SHR Vx {, Vy}
			// Set Vx = Vx SHR 1.
			// If the least-significant bit of Vx is 1, then VF is set to 1, otherwise 0. Then Vx is divided by 2 (SHR).
			case 0x0006:
				V[0xF] = V[x] & 0x01
				V[x] = V[x] >> 1
				PC += 2
				fmt.Printf("\t\tOpcode 8xy6 executed: Set V[x(%d)]:%d SHIFT RIGHT 1\n\n", x, V[x])
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

				fmt.Printf("\t\tOpcode 8xy7 executed: Set V[x(%d)]:%d = V[y(%d)]:%d - V[x(%d)]:%d\t\t = %d \n\n", x, V[x], y, V[y], x, V[x], V[y] - V[x])
				V[x] = V[y] - V[x]

				PC += 2
				break

			// 8xyE - SHL Vx {, Vy}
			// Set Vx = Vx SHL 1.
			// If the most-significant bit of Vx is 1, then VF is set to 1, otherwise to 0. Then Vx is multiplied by 2.
			case 0x000E:
				V[0xF] = V[x] >> 7 // Set V[F] to the Most Important Bit
				V[x] = V[x] << 1
				PC += 2
				fmt.Printf("\t\tOpcode 8xyE executed: Set V[x(%d)]:%d SHIFT LEFT 1\n\n", x, V[x])
				break

			default:
				fmt.Printf("\t\tOpcode 0x8000 NOT IMPLEMENTED!!!!\n\n")
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
				fmt.Printf("\t\tOpcode 9xy0 executed: V[x(%d)]:%d != V[y(%d)]:%d, SKIP one instruction\n\n", x, V[x], y, V[y])
			} else {
				PC += 2
				fmt.Printf("\t\tOpcode 9xy0 executed: V[x(%d)]:%d = V[y(%d)]:%d, DO NOT SKIP one instruction\n\n", x, V[x], y, V[y])
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
			fmt.Printf("\t\tOpcode Annn executed: Set I = %d\n\n", I)
			break


		// ############################ 0xB000 instruction set ############################
		// Bnnn - JP V0, addr
		// Jump to location nnn + V0.
		// The program counter is set to nnn plus the value of V0.
		case 0xB000:
			nnn := Opcode & 0x0FFF
			PC = nnn + uint16(V[0])
			print ("\t\tOpcode Bnnn executed: Jump to location nnn(%d) + V[0(%d)]\n\n", nnn, V[0])
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
			fmt.Printf("\t\tOpcode Cxkk executed: V[x(%d)] = %d (random byte AND kk(%d)) = \n\n", x, V[x], kk)
			break


		// ############################ 0xD000 instruction set ############################
		// Dxyn - DRW Vx, Vy, nibble
		// Display n-byte sprite starting at memory location I at (Vx, Vy), set VF = collision.
		case 0xD000: // DXYN

			var (
				x			uint16 = (Opcode & 0x0F00) >> 8
				y			uint16 = (Opcode & 0x00F0) >> 4
				n			uint16 = (Opcode & 0x000F)
				byte			uint16 = 0
				gpx_position	uint16 = 0
			)

			// Clean the colision flag
			V[0xF] = 0

			fmt.Printf("\t\tOpcode Dxyn(%X DRAW GRAPHICS! - Address I: %d Position V[x]: %d V[y]: %d N: %d bytes\n\n" , Opcode, I, V[x], V[y], n)

			// Check if y is out of range
			if (V[y] > 31) {
				V[y] = V[y] % 2
				fmt.Printf("\t\tV[y] > 31, modulus applied")
			}

			// Check if x is out of range
			if (V[x] > 63) {
				V[x] = V[x] % 64
				fmt.Printf("\t\tV[x] > 63, modulus applied")
			}

			// Translate the x and Y to the Graphics Vector
			gpx_position = (uint16(V[x]) + (64 * uint16(V[y])))

			// DEBUG
			//fmt.Printf ("\tGraphic vector position: %d\tValue: %d\n", gpx_position, Graphics[x + (64 * y)] )

			// Print N Bytes from address I in V[x]V[y] position of the screen
			for byte = 0 ; byte < n ; byte++ {

				var (
					binary string = ""
					sprite uint8 = 0
				)

				// Set the sprite
				sprite = uint8(Memory[I + byte])

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
					index := uint16(gpx_position) + uint16(bit) + (byte*64)

					// If tryes to draw bits outside the vector size, ignore
					if ( index > 2047) {
					   continue
					}

					// If bit=1, test current graphics[index], if is already set, mark v[F]=1 (colision)
					if (bit_binary  == 1){
						// Set colision case graphics[index] is already 1
						if (Graphics[index] == 1){
							V[0xF] = 1
						}
						// After, XOR the graphics[index] (DRAW)
						Graphics[index] ^= 1
					}

				// DEBUG 2
				//fmt.Printf ("\n\tByte: %d,\tSprite: %d\tBinary: %s\tbit: %d\tIndex: %d\tbinary[bit]: %c\tGraphics[index]: %d",byte, sprite, binary, bit, index, binary[bit], Graphics[index])
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
					fmt.Printf("\t\tOpcode Ex9E executed: Key[%d] pressed, skip one instruction\n\n",V[x])
				} else {
					PC += 2
					fmt.Printf("\t\tOpcode Ex9E executed: Key[%d] NOT pressed, continue\n\n",V[x])
				}
			 	break

			// ExA1 - SKNP Vx
			// Skip next instruction if key with the value of Vx is not pressed.
			// Checks the keyboard, and if the key corresponding to the value of Vx is currently in the up position, PC is increased by 2.
			case 0x00A1:
				if Key[V[x]] == 0 {
					PC += 4
					fmt.Printf("\t\tOpcode ExA1 executed: Key[%d] NOT pressed, skip one instruction\n\n",V[x])
				} else {
					Key[V[x]] = 0
					PC += 2
					fmt.Printf("\t\tOpcode ExA1 executed: Key[%d] pressed, continue\n\n",V[x])
				}
				break
			default:
				fmt.Printf("Opcode Family E000 - Not mapped opcote: E000\n")
				os.Exit(3)
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
				fmt.Printf("\t\tOpcode Fx07 executed: Set V[x(%d)] with value of DelayTimer(%d)\n\n", x, DelayTimer)
				break

			// Fx0A - LD Vx, K
			// Wait for a key press, store the value of the key in Vx.
		 	// All execution stops until a key is pressed, then the value of that key is stored in Vx.
			case 0x000A:
				for i := 0 ; i < len(Key) ; i++ {
					if (Key[i] == 1){
						V[x] = byte(i)
						PC +=2
						fmt.Printf("\tOpcode Fx0A executed: Wait for a key (Key[%d]) press -  (PRESSED)\n\n", i)
						// Stop after find the first key pressed
						break

					} else {
						fmt.Printf("\tOpcode Fx0A executed: Wait for a key (Key[%d]) press - (NOT PRESSED)\n\n", i)
					}
				}
				break


			// Fx15 - LD DT, Vx
			// Set delay timer = Vx.
			// DT is set equal to the value of Vx.
			case 0x0015:
				DelayTimer = V[x]
				PC += 2
				fmt.Printf("\t\tOpcode Fx15 executed: Set delay timer = V[x(%d):%d]\n\n", x, V[x])
				break

			// Fx18 - LD ST, Vx
			// Set sound timer = Vx.
			// ST is set equal to the value of Vx.
			case 0x0018:
				SoundTimer = V[x]
				PC += 2
				fmt.Printf("\t\tOpcode Fx18 executed: Set sound timer = V[x(%d)]:%d\n\n",x, V[x])
				break

			// Fx1E - ADD I, Vx
			// Set I = I + Vx.
			// The values of I and Vx are added, and the results are stored in I.
			case 0x001E:
				fmt.Printf("\t\tOpcode Fx1E executed: Add the value of V[x(%d)]:%d to I(%d)\n\n",x, V[x], I)
				I += uint16(V[x])
				PC += 2
				break

			// Fx29 - LD F, Vx
			// Set I = location of sprite for digit Vx.
			// The value of I is set to the location for the hexadecimal sprite corresponding to the value of Vx.
			case 0x0029:
				I = uint16(V[x]) * 5
				PC += 2
				fmt.Printf("\t\tOpcode Fx29 executed: Set I(%X) = location of sprite for digit V[x(%d)]:%d (*5)\n\n", I, x, V[x])
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
				fmt.Printf("\t\tOpcode Fx33 executed: Store BCD representation of V[x(%d)]:%d in memory locations I(%X):%d, I+1(%X):%d, and I+2(%X):%d\n\n", x, V[x], I, Memory[I], I+1, Memory[I+1], I+2, Memory[I+2])
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
				fmt.Printf("\t\tOpcode Fx55 executed: Registers V[0] through V[x(%d)] in memory starting at location I(%d)\n\n",x, I)
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
				fmt.Printf("\t\tOpcode Fx65 executed: Read registers V[0] through V[x(%d)] from memory starting at location I(%X)\n\n",x, I)
				break

			default:
				fmt.Printf("\t\tOpcode Family F000 - Not mapped opcode: 0x%X\n\n", Opcode)
				os.Exit(2)

			}
			break



		default:
			fmt.Printf("\t\tOPCODE FAMILY %X NOT IMPLEMENTED!\n\n", Opcode & 0xF000)
			os.Exit(3)
	}


	// Every Cycle Control the clock!!!
	select {
	case <-TimerClock.C:
	// When ticker run (60 times in a second, check de DelayTimer)
		if DelayTimer > 0 {
			//fmt.Printf("##############################DelayTimer= %d", DelayTimer)
			DelayTimer--
		}

		if SoundTimer > 0 {
			if SoundTimer == 1 {
				Sound.PlaySound(Sound.Beep_buffer)
		   }
		   SoundTimer--
		}

		default:
			// No timer to handle
		}

	Cycle ++

}
