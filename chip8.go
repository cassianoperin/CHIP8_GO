package main

import (
	"fmt"
	"log"
	"os"
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
	}

}


func get_game_signature() {

	// Used to identify games that needs legacy opcodes
	// Read the first 10 game instructions in memory
	signature_size := 10
	for i:=0 ; i < signature_size ; i++ {
		CPU.Game_signature += fmt.Sprintf("%X", CPU.Memory[int(CPU.PC)+i])
	}
	fmt.Printf("Game signature: %s\n\n", CPU.Game_signature)
}

func handle_legacy_opcodes() {
	// Game "Animal Race [Brian Astle].ch8"
	// MD5: 46497c35ce549cd7617462fe7c9fc284
	if (CPU.Game_signature == "6DA6E268E69BA5B5") {
		CPU.Legacy_Fx55_Fx65 = true
		fmt.Printf("Legacy mode Fx55/Fx65 enabled.\n")
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

	checkArgs()
	testFile(os.Args[1])

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
