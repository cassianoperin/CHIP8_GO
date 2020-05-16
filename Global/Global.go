// Global Variables
package Global

import (
	"github.com/faiface/pixel/pixelgl"
)

var (
	// Global
	Game_signature				string	// Game Signature (identify games that needs legacy opcodes)

	// Graphics
	Win		*pixelgl.Window
	WindowTitle	string = "Chip-8"
	Color_theme	= 0

	//Input
	InputDrawFlag	bool		// Force draw, necessary in some emulator rewind and forward status

)
