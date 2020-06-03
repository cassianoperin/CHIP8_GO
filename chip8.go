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
		os.Exit(0)
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

	// Load ROM from 0x200 address in memory, or 0x600 for hybrid hardware ETI-600
	for i := 0; i < len(data); i++ {
		if Global.Hybrid_ETI_660_HW {
			CPU.Memory[i+1536] = data[i]	// start at 0x600
		} else {
			CPU.Memory[i+512] = data[i]	// start at 0x200
		}
	}

}

func checkArgs() {

	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s [options] ROM_FILE\n\n%s -help for a list of available options\n\n", os.Args[0], os.Args[0] )
		os.Exit(0)
	}

	cliHelp			:= flag.Bool("help", false, "Show this menu")
	cliSchipHack		:= flag.Bool("SchipHack", false, "Enable SCHIP DelayTimer hack mode to improve speed")
	cliDrawFlag		:= flag.Bool("DrawFlag", false, "Enable Draw Graphics on each Drawflag instead @60Hz")
	cliDebug			:= flag.Bool("Debug", false, "Enable Debug Mode")
	cliRewind			:= flag.Bool("Rewind", false, "Enable Rewind Mode")
	cliHybridETI660	:= flag.Bool("Hybrid-ETI-660", false, "Enable ETI-660 mode for hybrid games made for this hardware")
	cliPause			:= flag.Bool("Pause", false, "Start emulation Paused")

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
		fmt.Printf("Usage: %s [options] ROM_FILE\n  -Debug\n    	Enable Debug Mode\n  -DrawFlag\n    	Enable Draw Graphics on each Drawflag instead @60Hz\n  -Hybrid-ETI-660\n    	Enable ETI-660 mode for hybrid games made for this hardware\n  -Pause\n    	Start emulation Paused\n  -Rewind Mode\n    	Enable Rewind Mode\n  -SchipHack\n    	Enable SCHIP DelayTimer hack mode to improve speed\n  -help\n    	Show this menu\n\n", os.Args[0])
		os.Exit(0)
	}

	if *cliRewind {
		CPU.Rewind_mode = true
	}

	if *cliSchipHack {
		CPU.SCHIP_TimerHack = true
	}

	if *cliDebug {
		CPU.Debug = true
	}

	if *cliDrawFlag {
		// Enable Draw at DrawFlag instead of @60Hz
		Global.OriginalDrawMode = true
	}

	if *cliHybridETI660 {
		// Enable ETI-660 Hardware mode (hybrid)
		// Store rom at 0x600 instead of default 0x200
		// The ETI 660 had 64 x 48 OR 64 x 64 with a modification
		Global.Hybrid_ETI_660_HW = true
	}

	if *cliPause {
		CPU.Pause = true
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
