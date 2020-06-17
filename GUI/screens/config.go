package screens

import (
	"os"
	"fmt"
	"fyne.io/fyne"
	"image/color"
	"fyne.io/fyne/canvas"
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

func ConfigScreen(a fyne.App, win fyne.Window) fyne.CanvasObject {

	// --------------------- LABELS --------------------- //

	labelNewLine := widget.NewLabel("\n")
	labelNewLine.Alignment = fyne.TextAlignCenter
	// labelNewLine.Alignment = fyne.AlignTrailing

	labelTheme := widget.NewLabel("\nTHEME")
	labelTheme.Alignment = fyne.TextAlignCenter

	labelOptions := widget.NewLabel("OPTIONS")
	labelOptions.Alignment = fyne.TextAlignCenter

	labelEmulator := widget.NewLabel("EMULATOR")
	labelEmulator.Alignment = fyne.TextAlignCenter

	labelROM := widget.NewLabel("ROM LOCATION")
	labelROM.Alignment = fyne.TextAlignCenter

	// --------------------- OPTIONS --------------------- //

	checkDebug := widget.NewCheck("Debug", func(on bool) { fmt.Println("Debug", on) })
	checkDebug.Disable()

	checkPause := widget.NewCheck("Start Paused", func(on bool) { fmt.Println("Start Paused", on) })
	checkPause.Disable()

	checkRewind:= widget.NewCheck("Rewind Mode", func(on bool) { fmt.Println("Rewind mode", on) })
	checkRewind.Disable()

	checkSchipHack := widget.NewCheck("SCHIP Hack", func(on bool) { fmt.Println("SCHIP Hack", on) })
	checkSchipHack.Disable()


	// --------------------- FORMS --------------------- //

	inputEmulator := widget.NewEntry()
	inputEmulator.SetPlaceHolder("Ex.: /emulators/chip8, C:" + "\\" + "emulators" + "\\" + "chip8.exe")
	inputEmulator.SetText(EMULATOR_PATH)

	inputChip8 := widget.NewEntry()
	inputChip8.SetPlaceHolder("Ex.: /roms/chip8, C:" + "\\" + "roms" + "\\" + "chip8")
	inputChip8.SetText(CHIP8_PATH)

	inputChip8HiRes := widget.NewEntry()
	inputChip8HiRes.SetPlaceHolder("Ex.: /roms/chip8HiRes, C:" + "\\" + "roms" + "\\" + "chip8HiRes")
	inputChip8HiRes.SetText(CHIP8HIRES_PATH)

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

		labelTheme,

		fyne.NewContainerWithLayout(layout.NewGridLayout(2),
			widget.NewButton("Dark", func() {
				a.Settings().SetTheme(theme.DarkTheme())
			}),
			widget.NewButton("Light", func() {
				a.Settings().SetTheme(theme.LightTheme())
			}),
		),


		labelNewLine,

		labelEmulator,

		fyne.NewContainerWithLayout(layout.NewFormLayout(),
			canvas.NewText("EMULATOR PATH:", color.NRGBA{0, 153, 255, 0xff}),
			inputEmulator,
		),

		labelNewLine,

		labelROM,

		fyne.NewContainerWithLayout(layout.NewFormLayout(),
			canvas.NewText("CHIP-8 ROM PATH:", color.NRGBA{0, 153, 255, 0xff}),
			inputChip8,
			canvas.NewText("CHIP-8 HiRes ROM PATH:", color.NRGBA{0, 153, 255, 0xff}),
			inputChip8HiRes,
			canvas.NewText("SCHIP ROM PATH:", color.NRGBA{0, 153, 255, 0xff}),
			inputSchip,
		),

		layout.NewSpacer(),

		labelOptions,

		fyne.NewContainerWithLayout(layout.NewGridLayout(2),
			checkDebug,
			checkPause,
			checkRewind,
			checkSchipHack,
		),

		layout.NewSpacer(),

		buttonSave,

		statusLabel,
	)

}
