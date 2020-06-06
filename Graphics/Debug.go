package Graphics

import (
	"fmt"
	"Chip8/CPU"
	"Chip8/Global"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
	"github.com/faiface/pixel/imdraw"
)

// Print Graphics on Console
func drawDebugScreen(imd *imdraw.IMDraw) {

	basePosition := screenHeight * (1 - Global.SizeYused)	// Value reserved for debug on screen

	// -------------------------- Draw Rectangles -------------------------- //
	// Background
	currentColor := imd.Color	// Keep the current collor
	imd.Color = colornames.Dimgray
	imd.Push(pixel.V ( 0 , 0  ) )
	imd.Push(pixel.V ( screenWidth , basePosition ) )
	imd.Rectangle(0)
	// Return to the theme color
	imd.Color = currentColor
	// Up bar
	imd.Push(pixel.V ( 0 , basePosition  ) )
	imd.Push(pixel.V ( screenWidth , basePosition -2 ) )
	imd.Rectangle(0)
	// Down bar
	imd.Push(pixel.V ( 0 , 0  ) )
	imd.Push(pixel.V ( screenWidth , 2 ) )
	imd.Rectangle(0)
	// Left bar
	imd.Push(pixel.V ( 0 , 0  ) )
	imd.Push(pixel.V ( 2 , basePosition ) )
	imd.Rectangle(0)
	// Right bar
	imd.Push(pixel.V ( screenWidth , 0  ) )
	imd.Push(pixel.V ( screenWidth -2 , basePosition ) )
	imd.Rectangle(0)
}

func drawDebugInfo() {

	basePosition := screenHeight * (1 - Global.SizeYused)	// Value reserved for debug on screen
	fontSize := basePosition / 150

	// -------------------------- Draw Text -------------------------- //

	// Debug
	cpuMessage = text.New(pixel.V(20, basePosition - (basePosition * 0.2) ), atlas)	// X, Y
	cpuMessage.Clear()
	cpuMessage.Color = colornames.Black
	fmt.Fprintf(cpuMessage, "DEBUG")
	cpuMessage.Draw(Global.Win, pixel.IM.Scaled(cpuMessage.Orig, fontSize * 1.3))

	cpuMessage = text.New(pixel.V(20, basePosition - (basePosition * 0.35) ), atlas)

	cpuMessage.Clear()
	cpuMessage.LineHeight = atlas.LineHeight() * 1.5

	// Cycle
	cpuMessage.Color = colornames.Black
	fmt.Fprintf(cpuMessage, "Cycle:        ")
	cpuMessage.Color = colornames.White
	text := ""
	if CPU.Cycle == 0 {
		text = fmt.Sprintf("%d  ",CPU.Cycle)
	} else {
		text = fmt.Sprintf("%d  ",CPU.Cycle - 1)
	}
	cpuMessage.Dot.X -= cpuMessage.BoundsOf(text).W()
	fmt.Fprintf(cpuMessage, text)
	// Opcode
	cpuMessage.Color = colornames.Black
	fmt.Fprintf(cpuMessage, "Opcode:")
	cpuMessage.Color = colornames.White
	text = fmt.Sprintf(" %04X  ",CPU.Opcode)
	fmt.Fprintf(cpuMessage, text)
	// PC
	cpuMessage.Color = colornames.Black
	fmt.Fprintf(cpuMessage, "PC:               ")
	cpuMessage.Color = colornames.White
	text = fmt.Sprintf("%d(0x%04X)  ",CPU.PC, CPU.PC)
	cpuMessage.Dot.X -= cpuMessage.BoundsOf(text).W()
	fmt.Fprintf(cpuMessage, text)
	// I
	cpuMessage.Color = colornames.Black
	fmt.Fprintf(cpuMessage, "I:       ")
	cpuMessage.Color = colornames.White
	text = fmt.Sprintf("%d  ",CPU.I)
	cpuMessage.Dot.X -= cpuMessage.BoundsOf(text).W()
	fmt.Fprintf(cpuMessage, text)
	// DT
	cpuMessage.Color = colornames.Black
	fmt.Fprintf(cpuMessage, "DT:      ")
	cpuMessage.Color = colornames.White
	text = fmt.Sprintf("%d  ",CPU.DelayTimer)
	cpuMessage.Dot.X -= cpuMessage.BoundsOf(text).W()
	fmt.Fprintf(cpuMessage, text)
	// ST
	cpuMessage.Color = colornames.Black
	fmt.Fprintf(cpuMessage, "ST:      ")
	cpuMessage.Color = colornames.White
	text = fmt.Sprintf("%d  ",CPU.SoundTimer)
	cpuMessage.Dot.X -= cpuMessage.BoundsOf(text).W()
	fmt.Fprintf(cpuMessage, text)
	// SP
	cpuMessage.Color = colornames.Black
	fmt.Fprintf(cpuMessage, "SP:     ")
	cpuMessage.Color = colornames.White
	text = fmt.Sprintf("%d",CPU.SP)
	cpuMessage.Dot.X -= cpuMessage.BoundsOf(text).W()
	fmt.Fprintf(cpuMessage, text)
	// Stack
	cpuMessage.Color = colornames.Black
	fmt.Fprintf(cpuMessage, "\nStack: ")
	cpuMessage.Color = colornames.White
	fmt.Fprintf(cpuMessage, "[   ")
	for i:=0 ; i  <len(CPU.Stack) ; i++ {
		text = fmt.Sprintf("%d",CPU.Stack[i])
		cpuMessage.Dot.X -= cpuMessage.BoundsOf(text).W()
		fmt.Fprintf(cpuMessage, text)
		if i < 15 {
			fmt.Fprintf(cpuMessage, "     ")
		}
	}
	fmt.Fprintf(cpuMessage, "]")
	// V
	cpuMessage.Color = colornames.Black
	fmt.Fprintf(cpuMessage, "\nV:     ")
	cpuMessage.Color = colornames.White
	fmt.Fprintf(cpuMessage, "[   ")
	for i:=0 ; i  <len(CPU.V) ; i++ {
		text = fmt.Sprintf("%d",CPU.V[i])
		cpuMessage.Dot.X -= cpuMessage.BoundsOf(text).W()
		fmt.Fprintf(cpuMessage, text)
		if i < 15 {
			fmt.Fprintf(cpuMessage, "     ")
		}
	}
	fmt.Fprintf(cpuMessage, "]")
	// Keys
	cpuMessage.Color = colornames.Black
	fmt.Fprintf(cpuMessage, "\nKeys: ")
	cpuMessage.Color = colornames.White
	text = fmt.Sprintf(" %d  ",CPU.Key)
	fmt.Fprintf(cpuMessage, text)
	//Opcode Message
	cpuMessage.Color = colornames.Black
	fmt.Fprintf(cpuMessage,"\nMsg:   ")
	cpuMessage.Color = colornames.White
	text = fmt.Sprintf("%s ",CPU.OpcMessage)
	fmt.Fprintf(cpuMessage, text)

	// Draw Text
	cpuMessage.Draw(Global.Win, pixel.IM.Scaled(cpuMessage.Orig, fontSize))

}
