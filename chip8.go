package main

import (
	"io"
	"os"
	"fmt"
	"log"
	"flag"
	"runtime"
	"strconv"
	"crypto/md5"
	"Chip8/CPU"
	"Chip8/Global"
	"Chip8/Graphics"
	"Chip8/Sound"
	"github.com/faiface/pixel/pixelgl"
)

var (
	hexFlag	bool
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

	// Don't run with files bigger than 4KB (binary) or 10KB (Hexadecimal)
	if hexFlag {
		if romsize >= 10240 {
			fmt.Printf("Hexadecimal bigger than 10KB, invalid ROM.\n")
			os.Exit(0)
		}
	} else {
		if romsize >= 4096 {
			fmt.Printf("Binary file bigger than 4KB, invalid ROM.\n")
			os.Exit(0)
		}
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
	// fmt.Printf("%d\n", data)
	// fmt.Printf("%X\n", data)
	// fmt.Printf("%s\n", data)


	// If rom format is HEXADECIMAL (.hex)
	if hexFlag {
		var (
			rom_raw	[]byte
			rom		[]byte
		)

		// Filter only hexadecimal characters
		for i := 0; i < len(data); i++ {

			// Put into rom_raw slice if is a number [0-9]
			for j := 0; j <= 9; j++ {
				// Compare the value in data[i] with [0-9]
				if string(data[i]) == strconv.Itoa(j) {
					tmp := fmt.Sprintf("0x%s",string(data[i]) )
					d, _ := strconv.ParseInt(tmp, 0, 10)
					rom_raw = append(rom_raw, byte(d))

				}
			}

			// Put into rom_raw slice if a letter [A-F]
			for j := 'A'; j <= 'F'; j++ {
				// Compare the value in data[i] with [A-F]
				if string(data[i]) == string(j) {
					tmp := fmt.Sprintf("0x%s",string(data[i]) )
					d, _ := strconv.ParseInt(tmp, 0, 10)
					rom_raw = append(rom_raw, byte(d))
				}
			}

			// Put into rom_raw slice if a letter [a-f]
			for j := 'a'; j <= 'f'; j++ {
				// Compare the value in data[i] with [a-f]
				if string(data[i]) == string(j) {
					tmp := fmt.Sprintf("0x%s",string(data[i]) )
					d, _ := strconv.ParseInt(tmp, 0, 10)
					rom_raw = append(rom_raw, byte(d))
				}
			}
		}

		// Agroup each 2 bytes into one
		for i := 0; i < len(rom_raw); i+=2 {
				tmp := fmt.Sprintf("0x%02X", uint8(rom_raw[i])<<4 | uint8(rom_raw[i+1]) )
				d, _ := strconv.ParseInt(tmp, 0, 10)
				rom = append(rom, byte(d))
		}

		// Put ROM into the memory (starting at 0x200)
		for i := 0; i < len(rom); i++ {
			if Global.Hybrid_ETI_660_HW {
				CPU.Memory[i+1536] = rom[i]	// start at 0x600
			} else {
				CPU.Memory[i+512] = rom[i]	// start at 0x200
			}
		}

		// fmt.Printf("\nROM (only Hex characters):\n%d\n", rom_raw)
		// fmt.Printf("\nROM:\n%02X\n", rom)

	// If rom format is BINARY (.ch8)
	} else {
		// Load ROM from 0x200 address in memory, or 0x600 for hybrid hardware ETI-600
		for i := 0; i < len(data); i++ {
			if Global.Hybrid_ETI_660_HW {
				CPU.Memory[i+1536] = data[i]	// start at 0x600
			} else {
				CPU.Memory[i+512] = data[i]	// start at 0x200
			}
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
	cliHybridETI660	:= flag.Bool("ETI660", false, "Enable ETI-660 mode for hybrid games made for this hardware")
	cliPause			:= flag.Bool("Pause", false, "Start emulation Paused")
	cliHex			:= flag.Bool("Hex", false, "Open roms in Hexadecimal format")

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
		fmt.Printf("Usage: %s [options] ROM_FILE\n  -Debug\n    	Enable Debug Mode\n  -DrawFlag\n    	Enable Draw Graphics on each Drawflag instead @60Hz\n  -ETI660\n    	Enable ETI-660 mode for hybrid games made for this hardware\n  -Hex\n    	Open roms in Hexadecimal format\n  -Pause\n    	Start emulation Paused\n  -Rewind Mode\n    	Enable Rewind Mode\n  -SchipHack\n    	Enable SCHIP DelayTimer hack mode to improve speed\n  -help\n    	Show this menu\n\n", os.Args[0])
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
		Global.SizeYused = 0.7 //Reserve debug screen area
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

	if *cliHex {
		hexFlag = true
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
