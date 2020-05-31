// ---------------------------- 05 CHIP-8 for ETI-660 opcodes ---------------------------- //
// 0000: Return to monitor (exit interpreter)
// 00F8: Turn display on
// 00FC: Turn display off
// 00FF: Do nothing
// FX00: Set pitch of sound signal

package CPU

import (
	"os"
	"fmt"
)

// ----------------------- CHIP-8 for ETI-660 0xxx instruction set ----------------------- //

// 0000
// Return to monitor (exit interpreter)
func opc_chip8_ETI660_0000() {
	if Debug {
		fmt.Printf("\t\tHybrid ETI-660 - Opcode 0000 executed: Return to monitor (exit interpreter)\n")
	}
}

// 00F8
// Turn display on
func opc_chip8_ETI660_00F8() {
	PC += 2
	if Debug {
		fmt.Printf("\t\tHybrid ETI-660 - Opcode 00F8 executed: Turn display on (Do nothing)\t PC+=2\n")
		fmt.Printf("\n\nPROPOSITAL EXIT TO MAP 00F8 USAGE!!!!\n\n")
		os.Exit(2)
	}
}

// 00FC
// Turn display off
func opc_chip8_ETI660_00FC() {
	PC+=2
	if Debug {
		fmt.Printf("\t\tHybrid ETI-660 - Opcode 00FC executed: Turn display off (Do nothing)\t PC+=2\n")
		fmt.Printf("\n\nPROPOSITAL EXIT TO MAP 00FC USAGE!!!!\n\n")
		os.Exit(2)
	}
}

// 00FF
// NO OP (do nothing)
func opc_chip8_ETI660_00FF() {
	PC += 2
	if Debug {
		fmt.Printf("\t\tHybrid ETI-660 - Opcode 00FF executed: No Operation (do nothing)\tPC+=2\n")
	}
}

// ----------------------- CHIP-8 for ETI-660 Fxxx instruction set ----------------------- //

// Fx00
// Set the pitch (frequency) of the tone generator (beeper) to Vx
func opc_chip8_ETI660_FX00(x uint16) {
	P = V[x]	// NOT USED YET!!! Need to implement sound library to handle it
	PC +=2

	if Debug {
		fmt.Printf("\t\tHybrid ETI-660 - Opcode Fx00 executed: Set the pitch (frequency) of the tone generator (beeper) to value of V[%d]\t\tP=%d\n", x, V[x])
	}
}
