package CPU

import (
	"time"
)

var (
	// Savestates
	Savestate_created			int = 0
	// CPU Components
	PC_savestate				uint16
	Stack_savestate				[16]uint16
	SP_savestate				uint16
	V_savestate				[16]byte
	I_savestate				uint16
	Graphics_savestate			[128 * 64]byte
	DelayTimer_savestate			byte
	SoundTimer_savestate			byte
	Cycle_savestate				uint16
	Rewind_index_savestate			uint16
	SCHIP_savestate				bool
	SCHIP_LORES_savestate			bool
	SizeX_savestate				float64
	SizeY_savestate				float64
	CPU_Clock_Speed_savestate		time.Duration
	Opcode_savestate			uint16
	Memory_savestate			[4096]byte
)
