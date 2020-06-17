package screens

import (
	"fmt"
	"strings"
	"errors"
	"os/exec"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/storage"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
)


func fileOpened(f fyne.URIReadCloser) {
	if f == nil {
		if Debug {
			fmt.Println("HOME TAB:\tFileOpen Dialog Cancelled")
		}
		return
	}

	// Extract the extension of the file
	extension := f.URI().Extension()

	// Print the Extension
	if Debug {
		fmt.Println("HOME TAB:\tExtension: ", extension)
	}

	// Save the full path and file name opened
	//filename = f.URI()[7:]
	filename = f.URI().String()[7:]

	// If opened a valid CH8 file, enable the run button
	// Save the file name opened
	lastBin := strings.LastIndex( filename, "/" )
	filename_short = fmt.Sprintf(filename[lastBin+1:])

	// Enable Button Run
	buttonRun.Enable()

	// Print the File Name
	if Debug {
		fmt.Println("HOME TAB:\tFile name: ",filename)
	}

	// Update the status bar
	status = fmt.Sprintf("Loaded rom: %s", filename_short)
	statusLabel.SetText(status)

	err := f.Close()
	if err != nil {
		fyne.LogError("\nHOME TAB:\tFailed to close stream", err)
	}
}



// DialogScreen loads a panel that lists the dialog windows that can be tested.
func HomeScreen(win fyne.Window) fyne.CanvasObject {

	// Logo object
	logo := canvas.NewImageFromFile(images_Chip8)
	logo.SetMinSize(fyne.NewSize(800, 600))

	// Status bar label object
	statusLabel = widget.NewLabel(status)

	// Run button object
	buttonRun = widget.NewButtonWithIcon("Run", theme.MediaPlayIcon(), func() {
		if Debug {
			fmt.Println("HOME TAB:\tButton Run Tapped! Filename: ", filename)
		}
		buttonRun.Disable()
		// Update the status bar
		status = fmt.Sprintf("Started game: %s", filename_short)
		statusLabel.SetText(status)

		// Open the file and send the stdout and stderr to a Dialog
		cmd := exec.Command(EMULATOR_PATH, filename)
		if Debug {
			fmt.Println("Filename: ", filename)
		}
		cmd.Dir = EmulatorFolder

		stdoutStderr, err := cmd.CombinedOutput()
		if err != nil {
			errorMessage:= fmt.Sprintf("Error running app!\n\n%s",err)
			execErr := errors.New(errorMessage)
			dialog.ShowError(execErr, win)


		if stdoutStderr != nil {
			errorMessage:= fmt.Sprintf("Error from Emulator!\n\n%s",stdoutStderr)
			execErr := errors.New(errorMessage)
			dialog.ShowError(execErr, win)
		}

		// Update the status bar
		statusLabel.SetText("")

		}
		if Debug {
			fmt.Println("StdoutStderr: ", stdoutStderr)
		}


	})
	// Start Disabled
	buttonRun.Disable()

	// Button Load Rom object
	buttonLoadRom := widget.NewButton("Load ROM", func() {
				fd := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
					if err == nil && reader == nil {
						return
					}
					if err != nil {
						dialog.ShowError(err, win)
						return
					}

					fileOpened(reader)

				}, win)
				fd.SetFilter(storage.NewExtensionFileFilter([]string{".ch8", ".CH8", ".cH8", ".Ch8"}))
				fd.Show()

			})



	return	widget.NewVBox(

			fyne.NewContainerWithLayout(layout.NewGridLayout(6),

				// Button LOAD ROM
				buttonLoadRom,

				// Button RUN
				buttonRun,

			), widget.NewHBox(layout.NewSpacer(), logo, layout.NewSpacer()),

			fyne.NewContainerWithLayout(layout.NewGridLayout(1),

				// Status bar
				statusLabel,

			),
		)

}
