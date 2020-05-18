package CPU

import (
	"os"
	"fmt"
	"time"
	"encoding/json"
	"io/ioutil"
	"Chip8/Global"
)

type SaveState struct {
	PC				uint16
	Stack			[16]uint16
	SP				uint16
	V				[16]byte
	I				uint16
	Graphics			[128 * 64]byte
	DelayTimer		byte
	SoundTimer		byte
	Cycle			uint16
	Rewind_index		uint16
	SCHIP			bool
	SCHIP_LORES		bool
	SizeX			float64
	SizeY			float64
	CPU_Clock_Speed	time.Duration
	Opcode			uint16
	Memory			[4096]byte
}


func SaveStateRead() {

	saveFile := Global.SavestateFolder + string(os.PathSeparator) + Global.Game_signature

	file, _ := ioutil.ReadFile(saveFile)

	data := SaveState{}

	data = SaveState{
	   PC:			PC,
        Stack:			Stack,
        SP:			SP,
	   V:			V,
	   I:			I,
	   Graphics:		Graphics,
	   DelayTimer:		DelayTimer,
	   SoundTimer:		SoundTimer,
	   Cycle:			Cycle,
	   Rewind_index:	Rewind_index,
	   SCHIP:			SCHIP,
	   SCHIP_LORES:	SCHIP_LORES,
	   SizeX:			Global.SizeX,
	   SizeY:			Global.SizeY,
	   CPU_Clock_Speed:	CPU_Clock_Speed,
   	   Opcode:		Opcode,
	   Memory:		Memory,
    }

	_ = json.Unmarshal([]byte(file), &data)


	Opcode			= data.Opcode
	PC				= data.PC
	Stack			= data.Stack
	SP				= data.SP
	V				= data.V
	I				= data.I
	Graphics			= data.Graphics
	DelayTimer		= data.DelayTimer
	SoundTimer		= data.SoundTimer
	Cycle			= data.Cycle
	Rewind_index		= data.Rewind_index
	SCHIP			= data.SCHIP
	SCHIP_LORES		= data.SCHIP_LORES
	Global.SizeX		= data.SizeX
	Global.SizeY		= data.SizeY
	CPU_Clock_Speed	= data.CPU_Clock_Speed
	Memory			= data.Memory

}


func SaveStateWrite() {

	saveFile := Global.SavestateFolder + string(os.PathSeparator) + Global.Game_signature

	// Create the Savestate Folder
	err := os.Mkdir(Global.SavestateFolder, 0777)
	if err == nil {
		if Debug {
			fmt.Println("Created Savestate Folder")
		}
	} else if os.IsExist(err){
		if Debug {
			fmt.Println("Savestate folder already exist")
		}
	} else {
		fmt.Println("Directory creation failed with error: " + err.Error())
		os.Exit(1)
	}

	data := SaveState{
	   PC:			PC,
        Stack:			Stack,
        SP:			SP,
	   V:			V,
	   I:			I,
	   Graphics:		Graphics,
	   DelayTimer:		DelayTimer,
	   SoundTimer:		SoundTimer,
	   Cycle:			Cycle,
	   Rewind_index:	Rewind_index,
	   SCHIP:			SCHIP,
	   SCHIP_LORES:	SCHIP_LORES,
	   SizeX:			Global.SizeX,
	   SizeY:			Global.SizeY,
	   CPU_Clock_Speed:	CPU_Clock_Speed,
   	   Opcode:		Opcode,
	   Memory:		Memory,
    }

    file, _ := json.MarshalIndent(data, "", " ")

    _ = ioutil.WriteFile(saveFile, file, 0777)
}
