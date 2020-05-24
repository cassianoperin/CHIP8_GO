package CPU

import (
	"fmt"
	"time"
	"Chip8/Global"
)

const (
	Rewind_buffer	uint16 = 100	 // Rewind Buffer Size
)

var (
	// Rewind Variables
	rewind_mode		bool	= false		// Enable and Disable Rewind Mode to increase emulation speed
	Rewind_index		uint16	= 0		// Rewind Index
	// CPU Components
	PC_track		= new([Rewind_buffer]uint16)
	SP_track		= new([Rewind_buffer]uint16)
	I_track			= new([Rewind_buffer]uint16)
	DT_track		= new([Rewind_buffer]byte)
	ST_track		= new([Rewind_buffer]byte)
	DF_track		= new([Rewind_buffer]bool)
	V_track			= new([Rewind_buffer][16]byte)
	Stack_track		= new([Rewind_buffer][16]uint16)
	GFX_track		= new([Rewind_buffer][128 * 64]byte)
)


func rewind() {

	// Start timer to measure procedures inside Interpreter
	start := time.Now()

	// REWIND MODE - SLICES
	// Just update when not inside a Rewind loop
	if Rewind_index == 0 {
		// PC
		copy(PC_track[1:], PC_track[0:])
		PC_track[0]	= PC
		// SP
		copy(SP_track[1:], SP_track[0:])
		SP_track[0]	= SP
		// I
		copy(I_track[1:], I_track[0:])
		I_track[0]	= I
		// DelayTimer
		copy(DT_track[1:], DT_track[0:])
		DT_track[0]	= DelayTimer
		// SoundTimer
		copy(ST_track[1:], ST_track[0:])
		ST_track[0]	= SoundTimer
		// DrawFlag
		copy(DF_track[1:], DF_track[0:])
		DF_track[0]	= Global.DrawFlag
		// V
		copy(V_track[1:], V_track[0:])
		V_track[0]	= V
		// Stack
		copy(Stack_track[1:], Stack_track[0:])
		Stack_track[0]	= Stack
		// GFX_track
		copy(GFX_track[1:], GFX_track[0:])
		GFX_track[0]	= Graphics
	}


	if Debug_L2 {
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


	// Debug time execution - Rewind Mode
	if Debug {
		elapsed := time.Since(start)
		fmt.Printf("\t\tTime track - Rewind Mode took: %s\n", elapsed)
	}

}
