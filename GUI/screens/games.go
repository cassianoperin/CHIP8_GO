package screens

import (
	"os"
	"fmt"
	"io/ioutil"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
	"fyne.io/fyne/theme"

	"path/filepath"
)

// GamesScreen loads a tab panel for containers and layouts
func GamesScreen() fyne.CanvasObject {

	// -------------------- List CHIP-8 ROMS -------------------- //
	// Read the directory (path "CHIP8_PATH")

	files_chip8, err := ioutil.ReadDir(CHIP8_PATH)
	if err != nil {
		fmt.Println("GAMES TAB:\t Error reading CHIP8 Directory: ", err)
	}

	// -------------------- List CHIP-8 HiRes ROMS -------------------- //
	// Read the directory (path "CHIP8_PATH")

	files_chip8hires, err := ioutil.ReadDir(CHIP8HIRES_PATH)
	if err != nil {
		fmt.Println("GAMES TAB:\t Error reading CHIP-8 HiRes Directory: ", err)
	}

	// -------------------- List SCHIP ROMS -------------------- //
	// Read the directory (path "SCHIP_PATH")
	files_schip, err := ioutil.ReadDir(SCHIP_PATH)
	if err != nil {
		fmt.Println("GAMES TAB:\t Error reading SCHIP Directory: ", err)
	}

	// -------------------- List CHIP-8 ROMS -------------------- //
	// Image
	chip8Background := canvas.NewImageFromFile(images_Pong)
	chip8Background.SetMinSize(fyne.NewSize(300, 282))

	// For each file create a button
	gameButtonList := makeButtonList(files_chip8, CHIP8_PATH)

	// Create a container with the buttons
	gameContainer := widget.NewVScrollContainer(widget.NewVBox(gameButtonList...))


	// -------------------- Tab CHIP-8 Content -------------------- //

	tabChip8 := widget.NewVBox(

		// Button RUN
		fyne.NewContainerWithLayout(layout.NewGridLayout(6), buttonRun, ),

		// Game Container
		widget.NewVBox(
			layout.NewSpacer(),

			fyne.NewContainerWithLayout(
					layout.NewAdaptiveGridLayout(1),
					fyne.NewContainerWithLayout(layout.NewBorderLayout(nil, nil, nil, nil),  gameContainer),
					chip8Background ),

			layout.NewSpacer() ),
	)




	// -------------------- List CHIP-8 HiRes ROMS -------------------- //
	// Image
	chip8hiresBackground := canvas.NewImageFromFile(images_Pong)
	chip8hiresBackground.SetMinSize(fyne.NewSize(300, 282))

	// For each file create a button
	gameButtonList = makeButtonList(files_chip8hires, CHIP8HIRES_PATH)

	// Create a container with the buttons
	gameContainer = widget.NewVScrollContainer(widget.NewVBox(gameButtonList...))


	// -------------------- Tab CHIP-8 HiRes Content -------------------- //

	tabChip8HiRes := widget.NewVBox(

		// Button RUN
		fyne.NewContainerWithLayout(layout.NewGridLayout(6), buttonRun, ),

		// Game Container
		widget.NewVBox(
			layout.NewSpacer(),

			fyne.NewContainerWithLayout(
					layout.NewAdaptiveGridLayout(1),
					fyne.NewContainerWithLayout(layout.NewBorderLayout(nil, nil, nil, nil),  gameContainer),
					chip8Background ),

			layout.NewSpacer() ),
	)



	// -------------------- List SCHIP ROMS -------------------- //

	// Image
	schipBackground := canvas.NewImageFromFile(images_Blinky)
	schipBackground.SetMinSize(fyne.NewSize(300, 282))

	// For each file create a button
	gameButtonList = makeButtonList(files_schip, SCHIP_PATH)

	// Create a container with the buttons
	gameContainer = widget.NewVScrollContainer(widget.NewVBox(gameButtonList...))


	// -------------------- Tab SCHIP Content -------------------- //

	tabSchip := widget.NewVBox(

		// Button RUN
		fyne.NewContainerWithLayout(layout.NewGridLayout(6), buttonRun, ),

		// Game Container
		widget.NewVBox(
			layout.NewSpacer(),

			fyne.NewContainerWithLayout(
					layout.NewAdaptiveGridLayout(1),
					fyne.NewContainerWithLayout(layout.NewBorderLayout(nil, nil, nil, nil), gameContainer),
					schipBackground ),

			layout.NewSpacer() ),
	)




	// Tabs container (CHIP-8/CHIP-8HiRes/SCHIP)
	tabs := widget.NewTabContainer(
		widget.NewTabItemWithIcon("Chip8", theme.HomeIcon(), tabChip8 ),
		widget.NewTabItemWithIcon("Chip8 HiRes", theme.ContentAddIcon(), tabChip8HiRes ),
		widget.NewTabItemWithIcon("SCHIP", theme.ViewRestoreIcon(), tabSchip) )
	tabs.SetTabLocation(widget.TabLocationTop)



	return widget.NewVBox(

		// Render the tabs (with their contents)
		tabs,

		// Render the Status Bar at the end of screen
		fyne.NewContainerWithLayout(layout.NewGridLayout(1), statusLabel,),
	)

}



func makeButtonList(files []os.FileInfo, path string) []fyne.CanvasObject {

	var items []fyne.CanvasObject

	for _, file := range files {
		if file.Mode().IsRegular() {
			if filepath.Ext(file.Name()) == ".ch8" || filepath.Ext(file.Name()) == ".CH8" {
			  if Debug {
				  fmt.Println("Tab GAME\tGame: ", file.Name())
			  }

			  index := file.Name() // capture

			  items = append(items, widget.NewButton(fmt.Sprintf("%s",file.Name()), func() {
				if Debug {
					fmt.Println("Tab GAME\tTapped: ",index )
				}
				// If opened a valid CH8 file, enable the run button
				buttonRun.Enable()

				filename = fmt.Sprintf("%s%s", path, index)

				// Print the File Name
				if Debug {
					fmt.Println("Tab GAME\tFile name: ", filename)
				}

				// Update the status bar
				status = fmt.Sprintf("Loaded rom: %s", index)

				statusLabel.SetText(status)
				}))
			}
		}
	}

	return items
}
