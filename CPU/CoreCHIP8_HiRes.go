// ---------------------------- 02 Two-page display for CHIP-8 opcodes ---------------------------- //
// 0230: Clear screen
// 1nnn: If PC=0x200 AND Opcode=0x1260, update Opcode to 0x12C0 (Jump to address 0x2c0)
// --------------------------- 01 Two-page display for CHIP-8X opcodes ---------------------------- //
// 00F0: // Return from subroutine (replaces 00EE)

package CPU

import (
	"fmt"
	"Chip8/Global"
)

// ---------------------------- CHIP-8 HiRes 0xxx instruction set ---------------------------- //

// CHIP-8 HIRES - 0230
// Clear screen used by Hi Resolution Chip8
func opc_chip8HiRes_0230() {
	Graphics = [128 * 64]byte{}

	// Set CHIP-8 HIRES Resolution
	Global.SizeX = 64
	Global.SizeY = 64

	PC += 2

	if Debug {
		OpcMessage = fmt.Sprintf("HIRES 0230: Clear the display")
		fmt.Printf("\t\t%s\n" , OpcMessage)
	}
}

// ---------------------------- CHIP-8 HiRes 1xxx instruction set ---------------------------- //

// HI-RES CHIP-8 EMULATION
// If PC=0x200 AND Opcode=0x1260, update Opcode to 0x12C0 (Jump to address 0x2c0)
func opc_chip8HiRes_1NNN() {
	PC = 0x2C0

	// After show the execution time
	if Debug {
		OpcMessage = fmt.Sprintf("HIRES 1260 WITH PC=0x200: Init 64x64 mode. PC=0x2C0 (jump to address 0x2c0)")
		fmt.Printf("\t\t%s\n" , OpcMessage)
	}
}

// ---------------------------- CHIP-8X HiRes 0xxx instruction set ---------------------------- //
// 00F0 Two-page display for CHIP-8X (Extension of CHIP-8x)
// Return from subroutine (replaces 00EE)
// Also used in some Hybrid ETI-660 programs like "Music Maker"
func opc_chip8HiRes_00F0() {
	PC = Stack[SP] + 2
	SP --
	if Debug {
		OpcMessage = fmt.Sprintf("CHIP-8X Two-page display 00F0: Return from subroutine (replaces 00EE in CHIP-8x)")
		fmt.Printf("\t\t%s\n" , OpcMessage)
	}
}
