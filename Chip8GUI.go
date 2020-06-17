package main

import (
	"os"
	"fmt"
	"strings"
	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"

	"gopkg.in/ini.v1"

	"Chip8/GUI/screens"
)


const (
	preferenceCurrentTab	= "currentTab"
)


func main() {

	// // ---- Read Configs from .ini file ---- //
	cfg, err := screens.LoadIni(screens.ConfigFile)
	if err != nil {
		fmt.Printf("MAIN TAB:\tFail to %v file. Exiting.\n", err)
		os.Exit(1)
	}

	// Require space before Comments
	cfg, err = ini.LoadSources(ini.LoadOptions{
		SpaceBeforeInlineComment: true,
	}, screens.ConfigFile)

	// Get Emulator Path from .ini
	screens.EMULATOR_PATH = cfg.Section("emulator").Key("chip8emu").String()
	// Get the Emulator Folder
	index := strings.LastIndex( screens.EMULATOR_PATH, "/" )
	screens.EmulatorFolder = fmt.Sprintf(screens.EMULATOR_PATH[:index+1])

	// Get Chip8 Rom Path from .ini
	screens.CHIP8_PATH = cfg.Section("paths").Key("chip8_roms").String()

	// Get Chip8 Rom Path from .ini
	screens.CHIP8HIRES_PATH = cfg.Section("paths").Key("chip8hires_roms").String()

	// Get SCHIP Rom Path from .ini
	screens.SCHIP_PATH = cfg.Section("paths").Key("schip_roms").String()


	// ---- Initialize the application ---- //
	a := app.NewWithID("Chip8GUI")
	a.SetIcon(theme.FyneLogo())

	// w := a.NewWindow("Chip8 Emulator")
	w := fyne.CurrentApp().NewWindow("Chip8 Emulator")

	// Resize
	// w.Resize(fyne.NewSize(800, 600))

	// Fixed window size
	w.SetFixedSize(true)

	// Center Window
	w.CenterOnScreen()

	// ---- Main Menu ---- //
	mainMenu := fyne.NewMainMenu(
		// a quit item will be appended to our first menu
		// fyne.NewMenu("File", newItem, fyne.NewMenuItemSeparator(), settingsItem),
		// fyne.NewMenu("Edit", cutItem, copyItem, pasteItem, fyne.NewMenuItemSeparator(), findItem),
		// helpMenu,
	)
	w.SetMainMenu(mainMenu)
	w.SetMaster()


	// ---- Tabs ---- //
	tabs := widget.NewTabContainer(
		widget.NewTabItemWithIcon("Main", theme.HomeIcon(), screens.HomeScreen(w)),
		widget.NewTabItemWithIcon("Games", theme.ViewRestoreIcon(), screens.GamesScreen()),
		widget.NewTabItemWithIcon("Configuration", theme.SettingsIcon(), screens.ConfigScreen(a, w)) )
	tabs.SetTabLocation(widget.TabLocationTop)
	tabs.OnChanged = func(tab *widget.TabItem) {
		if screens.Debug {
			fmt.Println("Tab selected: ", tab.Text)
		}
	}
	// tabs.SelectTabIndex(a.Preferences().Int(preferenceCurrentTab))
	tabs.SelectTabIndex(0)

	w.SetContent(tabs)

	w.ShowAndRun()
	a.Preferences().SetInt(preferenceCurrentTab, tabs.CurrentTabIndex())
}
