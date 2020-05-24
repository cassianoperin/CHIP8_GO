package main

import (
	"io"
	"os"
	"fmt"
	"log"
	"flag"
	"runtime"
	"crypto/md5"
	"Chip8/CPU"
	"Chip8/Global"
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

	// Don't run with files bigger than 4KB
	if romsize >= 4096 {
		fmt.Printf("File bigger than 4KB, invalid ROM.\n")
		os.Exit(1)
	}

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

}

func checkArgs() {

	if len(os.Args) < 2 {
		fmt.Printf("%d\n\n", len(os.Args) )

		fmt.Printf("Usage: %s [options] ROM_FILE\n%s -help for a list of available options\n\n", os.Args[0], os.Args[0] )
		os.Exit(0)
	}

	cliHelp		:= flag.Bool("help", false, "Show this menu")
	cliSchipHack	:= flag.Bool("SchipHack", false, "Enable SCHIP timer hack mode to improve speed")
	cliDrawFlag	:= flag.Bool("DrawFlag", false, "Enable Draw Graphics on each Drawflag instead @60Hz")

	// wordPtr := flag.String("word", "foo", "a string")
	// numbPtr := flag.Int("numb", 42, "an int")
	// var svar string
	// flag.StringVar(&svar, "ROM_FILE", "bar", "ROM_FILE")
	// fmt.Println("word:", *wordPtr)
	// fmt.Println("numb:", *numbPtr)
	// fmt.Println("svar:", svar)
	// fmt.Println("tail:", flag.Arg(0))
	flag.Parse()

	if *cliHelp {
		fmt.Printf("Usage: %s [options] ROM_FILE\n  -DrawFlag\n    	Enable Draw Graphics on each Drawflag instead @60Hz\n  -SchipHack\n    	Enable SCHIP timer hack mode to improve speed\n  -help\n    	Show this menu\n\n", os.Args[0])
		os.Exit(0)
	}

	if *cliSchipHack {
		CPU.SCHIP_TimerHack = true
	}

	if *cliDrawFlag {
		// Enable Draw at DrawFlag instead of @60Hz
		Global.OriginalDrawMode = true
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


func testFile(filename string) {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		fmt.Printf("File '%s' not found.\n\n", filename)
		os.Exit(0)
	}
}

func get_game_signature(filename string) {

	// Generate Game Signature (MD5)
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		log.Fatal(err)
	}
	Global.Game_signature = fmt.Sprintf("%x", h.Sum(nil))
	fmt.Printf("Game signature: %s\n", Global.Game_signature)
}


// Main function
func main() {

	// Validate the Arguments
	checkArgs()

	// Check if file exist
	testFile(flag.Arg(0))

	// Check the number of CPUS to create threads
	// fmt.Println("MaxParallelism: ", MaxParallelism())
	runtime.GOMAXPROCS( MaxParallelism() )

	// Set initial variables values
	CPU.Initialize()

	// Initialize sound buffer the sound function
	Sound.AudioDaemonStart("Sound/beep.wav")

	// Read ROM into Memory
	readROM(flag.Arg(0))

	// Get game signature
	get_game_signature(flag.Arg(0))

	// Start Window System and draw Graphics
	pixelgl.Run(Graphics.Run)

}
