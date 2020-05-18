package CPU

import (
	"fmt"
	"Chip8/Global"
)

var (
	// Legacy Opcodes and Quirks
	Legacy_Fx55_Fx65		bool	= false		// Enable original Chip-8 Fx55 and Fx65 opcodes (increases I)
	Legacy_8xy6_8xyE		bool	= false		// Enable original Chip-8 8xy6 and 8xyE opcodes
	FX1E_spacefight2091		bool	= false		// FX1E undocumented feature needed by Spacefight 2091!
	DXYN_bowling_wrap		bool	= false		// DXYN sprite wrap in Bowling game
	Resize_Quirk_00FE_00FF		bool	= true		// Resize_Quirk_00FE_00FF - Clears the screen - Must be set to True always
	DXY0_loresWideSpriteQuirks	bool	= false		// DXY0_loresWideSpriteQuirks - Draws a 16x16 sprite even in low-resolution (64x32) mode, row-major
	scrollQuirks_00CN_00FB_00FC	bool	= false		// Shift only 2 lines
)


func Handle_legacy_opcodes() {

	// Quirks needed by specific games

	// Enable Fx55 and Fx65 legacy mode
	// Game "Animal Race [Brian Astle]"
	// MD5: 46497c35ce549cd7617462fe7c9fc284
	if (Global.Game_signature == "6DA6E268E69BA5B5") {
		Legacy_Fx55_Fx65 = true
		fmt.Printf("Legacy mode Fx55/Fx65 enabled.\n")
	}
	// Enable 2nd legacy mode
	//if (Global.Game_signature == "xxxxxxxxxxxxx") {
	//	Legacy_8xy6_8xyE = true
	//	fmt.Printf("Legacy mode 8xy6/8xyE enabled.\n")
	//}

	// Enable undocumented FX1E feature needed by Spacefight 2091!
	// Game "Spacefight 2091 [Carsten Soerensen, 1992].ch8"
	// MD5: f99d0e82a489b8aff1c7203d90f740c3
	if (Global.Game_signature == "12245370616365466967") {
		FX1E_spacefight2091 = true
		fmt.Printf("FX1E undocumented feature enabled.\n")
	}
	// Enable undocumented FX1E feature needed by Spacefight 2091!
	// SCHIP Test Program "sctest_12"
	// MD5: 3ff053faaf994c051ed9b432f412b551
	if (Global.Game_signature == "12122054726F6E697820") {
		FX1E_spacefight2091 = true
		fmt.Printf("FX1E undocumented feature enabled.\n")
	}

	// Enable Pixel Wrap Fix for Bowling game
	// Game: "Bowling [Gooitzen van der Wal]"
	// MD5: b56e0e6e3930011049fcf6cf3384e964
	if (Global.Game_signature == "6314640255E60525B4") {
		DXYN_bowling_wrap = true
		fmt.Printf("DXYN pixel wrap fix enabled.\n")
	}

	// Enable Low Res 16x16 Pixel Draw in Robot.ch8 DEMO
	// SCHIP Demo: "Robot"
	// MD5: e2cd0812b43fb46e4b8abbb3a8d30f4b
	if (Global.Game_signature == "0FEA23A60061062F") {
		DXY0_loresWideSpriteQuirks = true
		fmt.Printf("DXY0 SCHIP Low Res 16x16 Pixel fix enabled.\n")
	}

}