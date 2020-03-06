package main

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"Chip8/CPU"
	"Chip8/Graphics"
	"Chip8/Sound"
	"github.com/faiface/pixel/pixelgl"
)


// Function used by readROM to avoid 'bytesread' return
func ReadContent(file *os.File, bytes_number int) []byte {

	bytes := make([]byte, bytes_number)

	_, err := file.Read(bytes)
	if err != nil {
		log.Fatal(err)
	}

	return bytes
}


// Read ROM and write it to the RAM
func readROM(filename string) {

	var (
		fileInfo os.FileInfo
		err      error
	)

	// Get ROM info
	fileInfo, err = os.Stat(filename)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Loading ROM:", filename)
	romsize := fileInfo.Size()
	fmt.Printf("Size in bytes: %d\n", romsize)

	// Open ROM file, insert all bytes into memory
	file, err := os.Open(filename)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}

	// Call ReadContent passing the total size of bytes
	data := ReadContent(file, int(romsize))
	// Print raw data
	//fmt.Printf("%d\n", data)
	//fmt.Printf("%X\n", data)

	// Load ROM from 0x200 address in memory
	for i := 0; i < len(data); i++ {
		CPU.Memory[i+512] = data[i]
	}

	// //Print Memory
	// for i := 0; i < len(CPU.Memory); i++ {
	// 	fmt.Printf("%X ", CPU.Memory[i])
	// }
	// os.Exit(2)
}

func checkArgs() {

	if len(os.Args) != 2 {
		fmt.Printf("Usage: %s ROM_FILE\n\n", os.Args[0] )
		os.Exit(2)
	}

}

func MaxParallelism() int {
	maxProcs := runtime.GOMAXPROCS(0)
	numCPU := runtime.NumCPU()
	if maxProcs < numCPU {
		return maxProcs
	}
	return numCPU
}


func get_game_signature() {

	// Used to identify games that needs legacy opcodes
	// Read the first 10 game instructions in memory
	signature_size := 10
	for i:=0 ; i < signature_size ; i++ {
		CPU.Game_signature += fmt.Sprintf("%X", CPU.Memory[int(CPU.PC)+i])
	}
	fmt.Printf("Game signature: %s\n", CPU.Game_signature)
}

func handle_legacy_opcodes() {

	// Quirks needed by specific games

	// Enable Fx55 and Fx65 legacy mode
	// Game "Animal Race [Brian Astle]"
	// MD5: 46497c35ce549cd7617462fe7c9fc284
	if (CPU.Game_signature == "6DA6E268E69BA5B5") {
		CPU.Legacy_Fx55_Fx65 = true
		fmt.Printf("Legacy mode Fx55/Fx65 enabled.\n")
	}
	// Enable 2nd legacy mode
	//if (CPU.Game_signature == "xxxxxxxxxxxxx") {
	//	CPU.Legacy_8xy6_8xyE = true
	//	fmt.Printf("Legacy mode 8xy6/8xyE enabled.\n")
	//}

	// Enable undocumented FX1E feature needed by Spacefight 2091!
	// Game "Spacefight 2091 [Carsten Soerensen, 1992].ch8"
	// MD5: f99d0e82a489b8aff1c7203d90f740c3
	if (CPU.Game_signature == "12245370616365466967") {
		CPU.FX1E_spacefight2091 = true
		fmt.Printf("FX1E undocumented feature enabled.\n")
	}
	// Enable undocumented FX1E feature needed by Spacefight 2091!
	// SCHIP Test Program "sctest_12"
	// MD5: 3ff053faaf994c051ed9b432f412b551
	if (CPU.Game_signature == "12122054726F6E697820") {
		CPU.FX1E_spacefight2091 = true
		fmt.Printf("FX1E undocumented feature enabled.\n")
	}

	// Enable Pixel Wrap Fix for Bowling game
	// Game: "Bowling [Gooitzen van der Wal]"
	// MD5: b56e0e6e3930011049fcf6cf3384e964
	if (CPU.Game_signature == "6314640255E60525B4") {
		CPU.DXYN_bowling_wrap = true
		fmt.Printf("DXYN pixel wrap fix enabled.\n")
	}

	// Enable Low Res 16x16 Pixel Draw in Robot.ch8 DEMO
	// SCHIP Demo: "Robot"
	// MD5: e2cd0812b43fb46e4b8abbb3a8d30f4b
	if (CPU.Game_signature == "0FEA23A60061062F") {
		CPU.DXY0_loresWideSpriteQuirks = true
		fmt.Printf("DXY0 SCHIP Low Res 16x16 Pixel fix enabled.\n")
	}

}


func testFile(filename string) {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		fmt.Printf("File '%s' not found.\n\n", os.Args[1])
		os.Exit(2)
	}
}



// Main function
func main() {

	// Validate the Arguments
	checkArgs()

	// Check if file exist
	testFile(os.Args[1])

	// Check the number of CPUS to create threads
	// fmt.Println("MaxParallelism: ", MaxParallelism())
	runtime.GOMAXPROCS( MaxParallelism() )


	// Set initial variables values
	CPU.Initialize()
	Sound.Initialize("Sound/beep.mp3")
	readROM(os.Args[1])

	// Identify special games that needs legacy opcodes
	get_game_signature()
	handle_legacy_opcodes()

	// Start Window System and draw Graphics
	pixelgl.Run(Graphics.Run)

}
