package Input

import (
	"fmt"
	"Chip8/Global"
	"github.com/faiface/pixel/pixelgl"
)


func Remap_keys() {

	// -------------- Just shoot games -------------- //

	// Platform: CHIP-8
	// Game: "Airplane.ch8"
	// Game: "Landing.ch8"
	// Game: "Missile [David Winter].ch8"
	if (Global.Game_signature == "02c2993360d186d3ae265ee7388481de" || Global.Game_signature == "9072c134d284f809fcf15f07a44aac01" ||
		Global.Game_signature == "ccf8de6e6fe799b56e27b6761251b107") {
		KeyPressedCHIP8[8] = pixelgl.KeySpace
		// Show messages
		fmt.Printf("Keys Remaped:\tShoot: SPACE\n\n")
		Global.TextMessageStr = "Keys Remaped\tShoot: SPACE"
	}

	// Platform: CHIP-8
	// Game: "Blitz [David Winter].ch8"
	// Game: "Submarine [Carmelo Cortez, 1978].ch8"
	if (Global.Game_signature == "8180b836eeb629ba93583519a5fb7b38" || Global.Game_signature == "f2ed4039d738eb118edd9b9de52960e6") {
		KeyPressedCHIP8[5] = pixelgl.KeySpace
		// Show messages
		fmt.Printf("Keys Remaped:\tShoot: SPACE\n\n")
		Global.TextMessageStr = "Keys Remaped\tShoot: SPACE"
	}

	// Platform: CHIP-8
	// Game: "Rocket Launcher.ch8"
	// Game: "Rocket [Joseph Weisbecker, 1978].ch8"
	if (Global.Game_signature == "e551a30b32b45451b0b67931744244e1" || Global.Game_signature == "bb64919f8eaef896e41ae94561b05fc8") {
		KeyPressedCHIP8[15] = pixelgl.KeySpace
		// Show messages
		fmt.Printf("Keys Remaped:\tShoot: SPACE\n\n")
		Global.TextMessageStr = "Keys Remaped\tShoot: SPACE"
	}

	// Platform: CHIP-8
	// Game: "Slide [Joyce Weisbecker].ch8"
	if (Global.Game_signature == "c6fa9e7a3f6dba491d1dfc1fe7b5df4e") {
		KeyPressedCHIP8[0] = pixelgl.KeySpace
		// Show messages
		fmt.Printf("Keys Remaped:\tShoot: SPACE\n\n")
		Global.TextMessageStr = "Keys Remaped\tShoot: SPACE"
	}

	// ------------ Left and right games ------------ //

	// Platform: CHIP-8
	// Game: "Breakout (Brix hack) [David Winter, 1997].ch8"
	// Game: "Breakout [Carmelo Cortez, 1979].ch8"
	// Game: "Brick (Brix hack, 1990).ch8"
	// Game: "Brix [Andreas Gustafsson, 1990].ch8"
	// Game: "Wipe Off [Joseph Weisbecker].ch8"
	if (Global.Game_signature == "3b26819c641e08cce4559aa1c68b09b1" || Global.Game_signature == "1d2a47947d6d46b1ceb41cf38b8cfc7e" ||
		Global.Game_signature == "008335a1292130403ddd5f222fa56944" || Global.Game_signature == "d677c1b9de941484d718799aebafebf3" ||
		Global.Game_signature == "41d64f82dc3c457e4f8543e081ae8e85" ) {
		KeyPressedCHIP8[4] = pixelgl.KeyLeft
		KeyPressedCHIP8[6] = pixelgl.KeyRight
		// Show messages
		fmt.Printf("Keys Remaped:\tLeft: ←\t\tRight: →\n\n")
		Global.TextMessageStr = "Keys Remaped\tLeft: LEFT KEY\t\tRight: RIGHT KEY"
	}

	// Platform: SCHIP
	// Game: "Car [Klaus von Sengbusch, 1994].ch8"
	if (Global.Game_signature == "c497bb692ea4b32a4a7b11b1373ef92f") {
		KeyPressedCHIP8[7] = pixelgl.KeyLeft
		KeyPressedCHIP8[8]  = pixelgl.KeyRight
		// Show messages
		fmt.Printf("Keys Remaped\tLeft: ←\t\tRight: →\tShoot: Space\n\n")
		Global.TextMessageStr = "Keys Remaped\tLeft: LEFT KEY\t\tRight: RIGHT KEY"
	}

	// -------- Left, Right and Action games ------- //

	// Platform: CHIP-8
	// Game: "Connect 4 [David Winter].ch8"
	// Game: "Space Invaders [David Winter].ch8"
	// Game: "Space Invaders [David Winter] (alt).ch8"
	// Game: "Rocket Launch [Jonas Lindstedt].ch8"
	if (Global.Game_signature == "29c5bf2d8f754dfe923923934513fb2d" || Global.Game_signature == "a67f58742cff77702cc64c64413dc37d" ||
		Global.Game_signature == "4fe20b951dbc801d7f682b88e672626c" || Global.Game_signature == "550dfbf87cf1dc62104e22def4378b3b" ) {
		KeyPressedCHIP8[4] = pixelgl.KeyLeft
		KeyPressedCHIP8[6] = pixelgl.KeyRight
		if Global.Game_signature == "550dfbf87cf1dc62104e22def4378b3b" {
			KeyPressedCHIP8[11] = pixelgl.KeySpace
		} else {
			KeyPressedCHIP8[5] = pixelgl.KeySpace
		}
		// Show messages
		fmt.Printf("Keys Remaped:\tLeft: ←\t\tRight: →\tStart/Action/Shoot: SPACE\n\n")
		Global.TextMessageStr = "Keys Remaped\tLeft: LEFT KEY\t\tRight: RIGHT KEY\t\tStart/Action/Shoot: SPACE"
	}

	// Platform: SCHIP
	// Game: "Spacefight 2091 [Carsten Soerensen, 1992].ch8"
	// Game: "Ant - In Search of Coke [Erin S. Catto].ch8"
	if (Global.Game_signature == "f99d0e82a489b8aff1c7203d90f740c3" || Global.Game_signature == "ec7856f9db5917eb6ca14adf1f8d0df2") {
		KeyPressedCHIP8[10] = pixelgl.KeySpace
		KeyPressedCHIP8[3]  = pixelgl.KeyLeft
		KeyPressedCHIP8[12] = pixelgl.KeyRight
		// Show messages
		fmt.Printf("Keys Remaped\tLeft: ←\t\tRight: →\tStart/Action/Shoot: Space\n\n")
		Global.TextMessageStr = "Keys Remaped\tLeft: LEFT KEY\t\tRight: RIGHT KEY\t\tStart/Action/Shoot: SPACE"
	}

	// Platform: CHIP-8
	// Game: "Space Intercept [Joseph Weisbecker, 1978].ch8"
	// Game: "UFO\ \[Lutz\ V,\ 1992\].ch8"
	if (Global.Game_signature == "2d7ab275b39ca46d9c7228b9cee46b63" || Global.Game_signature == "7e35f93c5a788e7e0027c78e8b76c8fb") {
		KeyPressedCHIP8[4] = pixelgl.KeyLeft
		KeyPressedCHIP8[5] = pixelgl.KeyUp
		KeyPressedCHIP8[6] = pixelgl.KeyRight
		// Show messages
		fmt.Printf("Keys Remaped\tLeft: ←\t\tRight: →\tUp: ↑\n\n")
		Global.TextMessageStr = "Keys Remaped\tLeft: LEFT KEY\t\tRight: RIGHT KEY\t\tUp: UP KEY"
	}


	// ---------- Up, Down and Action games -------- //

	// Platform: CHIP-8
	// Game: "Space Flight.ch8"
	// Game: "Vertical Brix [Paul Robson, 1996].ch8"
	if (Global.Game_signature == "dba8c50789808184c96f0173930bc81e" || Global.Game_signature == "9ec0fe6b275220f2c8821889a5a7fcab") {
		KeyPressedCHIP8[1] = pixelgl.KeyUp
		KeyPressedCHIP8[4] = pixelgl.KeyDown
		if Global.Game_signature == "9ec0fe6b275220f2c8821889a5a7fcab" {
			KeyPressedCHIP8[7] = pixelgl.KeySpace
		} else {
			KeyPressedCHIP8[15] = pixelgl.KeySpace
		}
		// Show messages
		fmt.Printf("Keys Remaped:\tUp: ↑\t\tDown: \tStart/Action/Shoot: SPACE\n\n")
		Global.TextMessageStr = "Keys Remaped\tUp: UP KEY\t\tDown: DOWN KEY\t\tStart/Action/Shoot: SPACE"
	}

	// -------------- Just arrow games -------------- //

	// Platform: CHIP-8
	// Game: "Blinky [Hans Christian Egeberg, 1991].ch8"
	// Platform: SCHIP
	// Game: "Blinky [Hans Christian Egeberg, 1991].ch8"
	if (Global.Game_signature == "e1c84e1156174070661c1f6ca0481ba5" || Global.Game_signature == "fb3284205c90d80c3b17aeea2eedf0e4") {
		KeyPressedCHIP8[3] = pixelgl.KeyUp
		KeyPressedCHIP8[6] = pixelgl.KeyDown
		KeyPressedCHIP8[7] = pixelgl.KeyLeft
		KeyPressedCHIP8[8] = pixelgl.KeyRight
		// Show messages
		fmt.Printf("Keys Remaped:\tLeft: ←\t\tRight: →\tUp: ↑\t\tDown: ↓\n\n")
		Global.TextMessageStr = "Keys Remaped\tLeft: LEFT KEY\t\tRight: RIGHT KEY\t\tUp: UP KEY\t\tDown: DOWN KEY"
	}

	// Platform: CHIP-8
	// Game: "Cave.ch8"
	// Game: "Shooting Stars [Philip Baltzer, 1978].ch8"
	// Game: "Worm V4 [RB-Revival Studios, 2007].ch8"
	// Game: "X-Mirror.ch8"
	if (Global.Game_signature == "bea7fdce1ef1733f9298dbbe2257cb9c" || Global.Game_signature == "8a202caf3b4f0fe3194276b8f8e508b7" ||
		Global.Game_signature == "4c0c381ac4942462b41876bb75e8a20a" || Global.Game_signature == "e868f362a4a91cf331753c55545dc271" ) {
		KeyPressedCHIP8[2] = pixelgl.KeyUp
		KeyPressedCHIP8[8] = pixelgl.KeyDown
		KeyPressedCHIP8[4] = pixelgl.KeyLeft
		KeyPressedCHIP8[6] = pixelgl.KeyRight
		// Show messages
		fmt.Printf("Keys Remaped:\tLeft: ←\t\tRight: →\tUp: ↑\t\tDown: ↓\n\n")
		Global.TextMessageStr = "Keys Remaped\tLeft: LEFT KEY\t\tRight: RIGHT KEY\t\tUp: UP KEY\t\tDown: DOWN KEY"
	}

	// ---------- Arrows and Action games ---------- //

	// Platform: CHIP-8
	// Game: "Astro Dodge [Revival Studios, 2008].ch8"
	// Game: "Kaleidoscope [Joseph Weisbecker, 1978].ch8"
	// Game: "Hidden [David Winter, 1996].ch8"
	// Game: "Most Dangerous Game [Peter Maruhnic].ch8"
	if (Global.Game_signature == "a7b171e6f738913f89153262b01581ba" || Global.Game_signature == "f330d48b32a2fd77cf41939f1d40ac06" ||
		Global.Game_signature == "d3f623110c962a28b86fc63e64bf33f0" || Global.Game_signature == "d8b3ccd5151d4b08edc0d5d87bb70603" ||
		Global.Game_signature == "dadaaf440809732d51485a2bc9410781" ) {
		KeyPressedCHIP8[2] = pixelgl.KeyUp
		KeyPressedCHIP8[8] = pixelgl.KeyDown
		KeyPressedCHIP8[4] = pixelgl.KeyLeft
		KeyPressedCHIP8[6] = pixelgl.KeyRight
		if Global.Game_signature == "d8b3ccd5151d4b08edc0d5d87bb70603" {	//Most Dangerous Game [Peter Maruhnic].ch8
			KeyPressedCHIP8[0] = pixelgl.KeySpace
		} else if Global.Game_signature == "dadaaf440809732d51485a2bc9410781" {	//Tapeworm [JDR, 1999].ch8
			KeyPressedCHIP8[15] = pixelgl.KeySpace
		} else {
			KeyPressedCHIP8[5] = pixelgl.KeySpace
		}
		// Show messages
		fmt.Printf("Keys Remaped:\tLeft: ←\t\tRight: →\tUp: ↑\t\tDown: ↓\t\tStart/Action/Shoot: SPACE\n\n")
		Global.TextMessageStr = "Keys Remaped\tLeft: LEFT KEY\t\tRight: RIGHT KEY\t\tUp: UP KEY\t\tDown: DOWN KEY\t\tStart/Action/Shoot: SPACE"
	}

	// Platform: CHIP-8
	// Game: "Tank.ch8"
	if (Global.Game_signature == "0ac0fd0d309c21a6cad14bb28217f5e4") {
		KeyPressedCHIP8[8] = pixelgl.KeyUp
		KeyPressedCHIP8[2] = pixelgl.KeyDown
		KeyPressedCHIP8[4] = pixelgl.KeyLeft
		KeyPressedCHIP8[6] = pixelgl.KeyRight
		KeyPressedCHIP8[5] = pixelgl.KeySpace
		// Show messages
		fmt.Printf("Keys Remaped:\tLeft: ←\t\tRight: →\tUp: ↑\t\tDown: ↓\t\tStart/Action/Shoot: SPACE\n\n")
		Global.TextMessageStr = "Keys Remaped\tLeft: LEFT KEY\t\tRight: RIGHT KEY\t\tUp: UP KEY\t\tDown: DOWN KEY\t\tStart/Action/Shoot: SPACE"
	}

	// Platform: CHIP-8
	// Game: "Tetris [Fran Dachille, 1991].ch8"
	if (Global.Game_signature == "aef4fc8c2a5e8431f5e0736ab281f2ee") {
		KeyPressedCHIP8[7] = pixelgl.KeyDown
		KeyPressedCHIP8[5] = pixelgl.KeyLeft
		KeyPressedCHIP8[6] = pixelgl.KeyRight
		KeyPressedCHIP8[4] = pixelgl.KeySpace
		// Show messages
		fmt.Printf("Keys Remaped:\tLeft: ←\t\tRight: →\t\tDown: ↓\t\tStart/Action/Shoot: SPACE\n\n")
		Global.TextMessageStr = "Keys Remaped\tLeft: LEFT KEY\t\tRight: RIGHT KEY\t\tDown: DOWN KEY\t\tStart/Action/Shoot: SPACE"
	}


	// Show Message on screen
	Global.ShowMessage = true
}
