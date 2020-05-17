package Input

import (
	"fmt"
	"Chip8/Global"
	"github.com/faiface/pixel/pixelgl"
)


func Remap_keys() {
	// Platform: SCHIP
	// Game: "Blinky [Hans Christian Egeberg, 1991].ch8"
	// MD5: fb3284205c90d80c3b17aeea2eedf0e4
	if (Global.Game_signature == "121A322E303020432E20") {
		KeyPressedCHIP8[3] = pixelgl.KeyUp
		KeyPressedCHIP8[6] = pixelgl.KeyDown
		KeyPressedCHIP8[7] = pixelgl.KeyLeft
		KeyPressedCHIP8[8] = pixelgl.KeyRight
		// Show messages
		fmt.Printf("Keys Remaped:\tLeft: ←\t\tRight: →\tUp: ↑\t\tDown: ↓\n\n")
		Global.TextMessageStr = "Keys Remaped\tLeft: LEFT KEY\t\tRight: RIGHT KEY\t\tUp: UP KEY\t\tDown: DOWN KEY"
		Global.ShowMessage = true
	}

	// Platform: SCHIP
	// Game: "Spacefight 2091 [Carsten Soerensen, 1992].ch8"
	// MD5: f99d0e82a489b8aff1c7203d90f740c3
	if (Global.Game_signature == "12245370616365466967") {
		KeyPressedCHIP8[10] = pixelgl.KeySpace
		KeyPressedCHIP8[3] = pixelgl.KeyLeft
		KeyPressedCHIP8[12] = pixelgl.KeyRight
		// Show messages
		fmt.Printf("Keys Remaped\tLeft: ←\t\tRight: →\tShoot: Space\n\n")
		Global.TextMessageStr = "Keys Remaped\tLeft: LEFT KEY\t\tRight: RIGHT KEY\t\tShoot: SPACE"
		Global.ShowMessage = true
	}

	// Platform: CHIP-8
	// Game: "Space Invaders [David Winter].ch8"
	// MD5: a67f58742cff77702cc64c64413dc37d
	if (Global.Game_signature == "1225535041434520494E") {
		KeyPressedCHIP8[5] = pixelgl.KeySpace
		KeyPressedCHIP8[4] = pixelgl.KeyLeft
		KeyPressedCHIP8[6] = pixelgl.KeyRight
		// Show messages
		fmt.Printf("Keys Remaped\tLeft: ←\t\tRight: →\tShoot: Space\n\n")
		Global.TextMessageStr = "Keys Remaped\tLeft: LEFT KEY\t\tRight: RIGHT KEY\t\tShoot: SPACE"
		Global.ShowMessage = true
	}

	// Platform: SCHIP
	// Game: "Ant - In Search of Coke [Erin S. Catto].ch8"
	// MD5: ec7856f9db5917eb6ca14adf1f8d0df2
	if (Global.Game_signature == "12E5B20416E74207631") {
		KeyPressedCHIP8[10] = pixelgl.KeySpace
		KeyPressedCHIP8[3]  = pixelgl.KeyLeft
		KeyPressedCHIP8[12]  = pixelgl.KeyRight
		// Show messages
		fmt.Printf("Keys Remaped\tLeft: ←\t\tRight: →\tShoot: Space\n\n")
		Global.TextMessageStr = "Keys Remaped\tLeft: LEFT KEY\t\tRight: RIGHT KEY\t\tShoot: SPACE"
		Global.ShowMessage = true
	}

	// Platform: SCHIP
	// Game: "Car [Klaus von Sengbusch, 1994].ch8"
	// MD5: c497bb692ea4b32a4a7b11b1373ef92f
	if (Global.Game_signature == "121E312E3030204B2E76") {
		KeyPressedCHIP8[7] = pixelgl.KeyLeft
		KeyPressedCHIP8[8]  = pixelgl.KeyRight
		// Show messages
		fmt.Printf("Keys Remaped\tLeft: ←\t\tRight: →\tShoot: Space\n\n")
		Global.TextMessageStr = "Keys Remaped\tLeft: LEFT KEY\t\tRight: RIGHT KEY"
		Global.ShowMessage = true
	}

}
