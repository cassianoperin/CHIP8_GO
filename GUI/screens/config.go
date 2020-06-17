package screens

import (
	"os"
	"fmt"
	"fyne.io/fyne"
	"image/color"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"

	"gopkg.in/ini.v1"
)

const (
	ConfigFile	= "GUI/Chip8GUI.ini"
)

var (
	Debug bool = false	// Debug flag

	// Images
	images_Pong		string = "GUI/images/pong.png"
	images_Chip8		string = "GUI/images/chip8.png"
	images_Blinky		string = "GUI/images/blinky.jpg"

	// Button Run object that needs to be reused in other functions
	buttonRun			*widget.Button

	// Extension of the file opened in "FileOpen Dialog"
	extension			string
	filename			string
	filename_short		string

	// Status bar
	status			string = "No loaded rom"
	statusLabel		*widget.Label

	files_chip8		[]os.FileInfo
	files_schip		[]os.FileInfo
	EmulatorFolder		= ""
	EMULATOR_PATH		= ""
	CHIP8_PATH		string = ""
	CHIP8HIRES_PATH	string = ""
	SCHIP_PATH		string = ""
)

// func LoadIni(ConfigFile string) {
func LoadIni(ConfigFile string) (*ini.File, error) {
	// ---- Read Configs from .ini file ---- //
	cfg, err := ini.Load(ConfigFile)

	return cfg, err
}

func makeTextGrid(text string) *widget.TextGrid {
	grid := widget.NewTextGridFromString(text)
	grid.SetStyleRange(0, 0, 0, 20,
		&widget.CustomTextGridStyle{FGColor: color.NRGBA{R: 0, G: 153, B: 255, A: 255}})

	grid.ShowLineNumbers = false
	grid.ShowWhitespace = false

	return grid
}


func ConfigScreen(a fyne.App, win fyne.Window) fyne.CanvasObject {

	labelNewLine := widget.NewLabel("\n")
	labelNewLine.Alignment = fyne.TextAlignCenter
	// labelNewLine.Alignment = fyne.AlignTrailing

	labelEmulator := makeTextGrid("EMULATOR PATH:")
	// labelEmulator := widget.NewLabel("Emulator Path:")
	// labelEmulator.Alignment = fyne.TextAlignLeading
	inputEmulator := widget.NewEntry()
	inputEmulator.SetPlaceHolder("Ex.: /emulators/chip8, C:" + "\\" + "emulators" + "\\" + "chip8.exe")
	inputEmulator.SetText(EMULATOR_PATH)

	labelChip8 := makeTextGrid("CHIP-8 ROM PATH:")
	// labelChip8 := widget.NewLabel("Chip8 Rom Path:")
	// labelChip8.Alignment = fyne.TextAlignLeading
	inputChip8 := widget.NewEntry()
	inputChip8.SetPlaceHolder("Ex.: /roms/chip8, C:" + "\\" + "roms" + "\\" + "chip8")
	inputChip8.SetText(CHIP8_PATH)

	labelChip8HiRes := makeTextGrid("CHIP-8 HiRes ROM PATH:")
	// labelChip8HiRes := widget.NewLabel("Chip8 HiRes Rom Path:")
	// labelChip8HiRes.Alignment = fyne.TextAlignLeading
	inputChip8HiRes := widget.NewEntry()
	inputChip8HiRes.SetPlaceHolder("Ex.: /roms/chip8HiRes, C:" + "\\" + "roms" + "\\" + "chip8HiRes")
	inputChip8HiRes.SetText(CHIP8HIRES_PATH)

	labelSchip := makeTextGrid("SCHIP ROM PATH:")
	// labelSchip := widget.NewLabel("Schip Rom Path:")
	// labelSchip.Alignment = fyne.TextAlignLeading
	inputSchip := widget.NewEntry()
	inputSchip.SetPlaceHolder("Ex.: /roms/schip, C:" + "\\" + "roms" + "\\" + "schip")
	inputSchip.SetText(SCHIP_PATH)

	buttonSave := widget.NewButtonWithIcon("Save", theme.DocumentSaveIcon(), func() {

		// // ---- Read Configs from .ini file ---- //
		cfg, err := LoadIni(ConfigFile)
		if err != nil {
			fmt.Printf("CONFIG TAB:\tFail to %v file. Exiting.\n", err)
			os.Exit(1)
		}

		// Demandas space before Comments
		cfg, err = ini.LoadSources(ini.LoadOptions{
			SpaceBeforeInlineComment: true,
		}, ConfigFile)

		// Update Emulator Settings on .ini
		// tmp := fmt.Sprintf("\"%s\"", inputEmulator.Text)
		cfg.Section("emulator").Key("chip8emu").SetValue(inputEmulator.Text)
		EMULATOR_PATH = inputEmulator.Text

		// tmp = fmt.Sprintf("\"%s\"", inputChip8.Text)
		cfg.Section("paths").Key("chip8_roms").SetValue(inputChip8.Text)
		CHIP8_PATH = inputChip8.Text

		// tmp = fmt.Sprintf("\"%s\"", inputChip8.Text)
		cfg.Section("paths").Key("chip8hires_roms").SetValue(inputChip8HiRes.Text)
		CHIP8HIRES_PATH = inputChip8HiRes.Text

		// tmp = fmt.Sprintf("\"%s\"", inputSchip.Text)
		cfg.Section("paths").Key("schip_roms").SetValue(inputSchip.Text)
		SCHIP_PATH = inputSchip.Text

		cfg.SaveTo(ConfigFile)

		dialog.ShowInformation("Information", "Configuration Updated!\n", win)

	})



	return widget.NewVBox(

		widget.NewGroup("   THEME   ",
			fyne.NewContainerWithLayout(layout.NewGridLayout(2),
				widget.NewButton("Dark", func() {
					a.Settings().SetTheme(theme.DarkTheme())
				}),
				widget.NewButton("Light", func() {
					a.Settings().SetTheme(theme.LightTheme())
				}),
			),
		),

		labelNewLine,

		widget.NewGroup("   EMULATOR   ",
			labelEmulator,
			inputEmulator,
		),

		labelNewLine,

		widget.NewGroup("   ROM LOCATION   ",
			labelChip8,
			inputChip8,
			labelChip8HiRes,
			inputChip8HiRes,
			labelSchip,
			inputSchip,
		),

		layout.NewSpacer(),

		buttonSave,

		statusLabel,
	)

}
