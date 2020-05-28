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
	ETI660_64x32_screen			bool = false		// Enable screen adjustment to 64x32 instead of default 64x48 ETI-660 HW
)


func Handle_legacy_opcodes() {

	// Quirks needed by specific games

	// Enable Fx55 and Fx65 legacy mode
	// Game "Animal Race [Brian Astle]"
	if (Global.Game_signature == "46497c35ce549cd7617462fe7c9fc284") {
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
	if (Global.Game_signature == "f99d0e82a489b8aff1c7203d90f740c3") {
		FX1E_spacefight2091 = true
		fmt.Printf("FX1E undocumented feature enabled.\n")
	}
	// Enable undocumented FX1E feature needed by Spacefight 2091!
	// SCHIP Test Program "sctest_12 (SC Test.ch8)"
	if (Global.Game_signature == "3ff053faaf994c051ed9b432f412b551") {
		FX1E_spacefight2091 = true
		fmt.Printf("FX1E undocumented feature enabled.\n")
	}

	// Enable Pixel Wrap Fix for Bowling game
	// Game: "Bowling [Gooitzen van der Wal]"
	if (Global.Game_signature == "b56e0e6e3930011049fcf6cf3384e964") {
		DXYN_bowling_wrap = true
		fmt.Printf("DXYN pixel wrap fix enabled.\n")
	}

	// Enable Low Res 16x16 Pixel Draw in Robot.ch8 DEMO
	// SCHIP Demo: "Robot"
	if (Global.Game_signature == "e2cd0812b43fb46e4b8abbb3a8d30f4b") {
		DXY0_loresWideSpriteQuirks = true
		fmt.Printf("DXY0 SCHIP Low Res 16x16 Pixel fix enabled.\n")
	}

	// This game uses 64x32 screen size
	// CHIP-8 ETI-660 Hybrid: "Pong"
	if (Global.Game_signature == "8e835a8da5e7d713819d7d70279cf998") {
		ETI660_64x32_screen = true
		fmt.Printf("ETI-660 Quirk 64 x 32 resolution Enabled.\n")
		Global.SizeX = 64
		Global.SizeY = 32
	}


}
