// Global Variables used to to avoid circular dependencies
package Global

import (
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
)

// Fullscreen / Video Modes
type Setting struct {
	Mode	*pixelgl.VideoMode
	Monitor	*pixelgl.Monitor
}

var (
	// Global
	Game_signature		string	= ""	// Game Signature (identify games that needs legacy opcodes)

	// Graphics
	Win			*pixelgl.Window
	WindowTitle		string = "Chip-8"
	Color_theme		= 0
	// Fullscreen / Video Modes
	Texts			[]*text.Text
	StaticText		*text.Text
	Settings		[]Setting
	ActiveSetting		*Setting
	IsFullScreen		= false		// Fullscrenn flag
	ResolutionCounter	int = 0		// Index of the available video resolution supported
	// FPS
	ShowFPS			bool		// Show or hide FPS counter flag
	// On screen messages
	ShowMessage		bool
	TextMessageStr		string

	//Input
	InputDrawFlag		bool		// Force draw, necessary in some emulator rewind and forward status

)
