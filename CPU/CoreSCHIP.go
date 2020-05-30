// ---------------------------- 06 Super-CHIP 1.0 opcodes ---------------------------- //
// 00FD: Exit interpreter
// 00FE: Disable high-resolution mode
// 00FF: Enable high-resolution mode
// DXY0: Draw 16 x 16 sprite (only if high-resolution mode is enabled)
// FX75: Store V0..VX in RPL user flags (X <= 7)
// FX85: Read V0..VX from RPL user flags (X <= 7)
// ---------------------------- 04 Super-CHIP 1.1 opcodes ---------------------------- //
// 00CN: Scroll display N pixels down; in low resolution mode, N/2 pixels
// 00FB: Scroll right by 4 pixels; in low resolution mode, 2 pixels
// 00FC: Scroll left by 4 pixels; in low resolution mode, 2 pixels
// FX30: Point I to 10-byte font sprite for digit VX (only digits 0-9)

package CPU

import (
	"os"
	"fmt"
	"time"
	"Chip8/Global"
)

// ---------------------------- SCHIP 0xxx instruction set ---------------------------- //

// SCHIP - 00CN
// Scroll display N lines down
func opc_schip_00CN(x uint16) {
	SCHIP = true

	shift := int(x) * 128

	// If in SCHIP Low Res mode, scroll N/2 lines only
	if scrollQuirks_00CN_00FB_00FC {
		shift = (int(x) * 128 ) / 2
	}

	// Shift Right N lines on Graphics Array
	for i:=len(Graphics) -1 ; i >= shift ; i-- {
		Graphics[i] = Graphics[i - shift]
	}

	// Clean the shifted display bytes
	for i:=0 ; i < shift ; i++ {
		Graphics[i] = 0
	}

	Global.DrawFlag = true
	DrawFlagCounter ++

	PC += 2
	if Debug {
		fmt.Printf("\t\tSCHIP - Opcode 00CN executed. - Scroll display %d lines down.\n", int(x))
	}
}

// SCHIP - 00FB
// Scroll display 4 pixels right
func opc_schip_00FB() {
	shift := 4	// Number of bytes to be shifted

	// If in SCHIP Low Res mode, shift 2 pixels only
	if scrollQuirks_00CN_00FB_00FC {
		shift = 2
	}

	rowsize := int(Global.SizeX)
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
			index -= int(Global.SizeX)
		}
	}

	Global.DrawFlag	= true
	DrawFlagCounter ++

	PC += 2
	if Debug {
		fmt.Printf("\t\tSCHIP - Opcode 00FB executed. - Scroll display 4 pixels right.\n")
	}
}

// SCHIP - 00FC
// Scroll display 4 pixels left
func opc_schip_00FC() {
	shift := 4

	// If in SCHIP Low Res mode, shift 2 pixels only
	if scrollQuirks_00CN_00FB_00FC {
		shift = 2
	}

	rowsize := int(Global.SizeX)

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
			rowsize += int(Global.SizeX)
		}
	}

	Global.DrawFlag	= true
	DrawFlagCounter ++

	PC += 2
	if Debug {
		fmt.Printf("\t\tSCHIP - Opcode 00FC executed. - Scroll display 4 pixels left.\n")
	}
}


// SCHIP - 00FD
// Exit Emulator
func opc_schip_00FD() {
	fmt.Printf("SCHIP - Opcode 00FD executed. - Exit emulator.\n")
	os.Exit(0)
}

// SCHIP - 00FE
// Enable Low-Res Mode (64 x 32 resolution)
func opc_schip_00FE() {
	// Disable SCHIP Mode
	SCHIP = false
	SCHIP_LORES = true
	scrollQuirks_00CN_00FB_00FC = true

	// Set the clock to CHIP-8 Speed
	CPU_Clock_Speed = 500
	CPU_Clock.Stop()
	CPU_Clock = time.NewTicker(time.Second / CPU_Clock_Speed)

	// Set CHIP-8 Resolution
	Global.SizeX = 64
	Global.SizeY = 32

	if Resize_Quirk_00FE_00FF {
		// Clear the screen when changing graphic mode
		Graphics	= [128 * 64]byte{}
	}

	PC += 2
	if Debug {
		fmt.Printf("\t\tSCHIP - Opcode 00FE executed. - Enable low res (64 x 32) mode.\n")
	}
}

// SCHIP - 00FF
// Enable High-Res Mode (128 x 64 resolution)
// Enable SCHIP Mode
func opc_schip_00FF() {
	SCHIP = true
	SCHIP_LORES = false
	scrollQuirks_00CN_00FB_00FC = false

	// Set the clock to SCHIP
	CPU_Clock_Speed = 1500
	CPU_Clock.Stop()
	CPU_Clock = time.NewTicker(time.Second / CPU_Clock_Speed)

	// Set SCHIP Resolution
	Global.SizeX = 128
	Global.SizeY = 64

	if Resize_Quirk_00FE_00FF {
		// Clear the screen when changing graphic mode
		Graphics	= [128 * 64]byte{}
	}

	PC += 2
	if Debug {
		fmt.Printf("\t\tSCHIP - Opcode 00FF executed. - Enable high res (128 x 64) mode.\n")
	}
}

// ---------------------------- SCHIP Dxxx instruction set ---------------------------- //
func opc_schip_DXY0() {
	// Migrate later
}

// ---------------------------- SCHIP Fxxx instruction set ---------------------------- //

// SCHIP Fx30 - LD F, Vx
// Set I = location of sprite for digit Vx.
// The value of I is set to the location for the hexadecimal sprite corresponding to the value of Vx.
func opc_schip_FX30(x uint16) {
	// Load SCHIP font. Start from Memory[80]
	I = 80 + uint16(V[x]) * 10
	PC += 2
	if Debug {
		fmt.Printf("\t\tSCHIP Opcode Fx30 executed: Set I(%X) = location of sprite for digit V[x(%d)]:%d (*10)\n", I, x, V[x])
	}
}

// SCHIP FX75
// Store V0 through VX to HP-48 RPL user flags (X <= 7).
func opc_schip_FX75(x uint16) {
	for i := 0; i <= int(x); i++ {
		RPL[i] = V[i]
	}

	PC += 2
	if Debug {
		fmt.Printf("\t\tSCHIP - Opcode Fx75 executed: Read RPL user flags from 0 to %d and store in V[0] through V[x(%d)]\n",x,x)
	}
}

// SCHIP FX85
// Read V0 through VX to HP-48 RPL user flags (X <= 7).
func opc_schip_FX85(x uint16) {
	for i := 0; i <= int(x); i++ {
		V[i] = RPL[i]
	}

	PC += 2
	if Debug {
		fmt.Printf("\t\tSCHIP - Opcode Fx85 executed: Read registers V[0] through V[x(%d)] and store in RPL user flags\n",x)
	}
}
