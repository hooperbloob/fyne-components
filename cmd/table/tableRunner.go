package main

import (
	"tableApp1/domains"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/driver/desktop"
)

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Generic Table Demo")
	myWindow.Resize(fyne.NewSize(500, 200))

	//tableContainer := domains.SetupPeopleTable(myWindow)
	tableContainer := domains.SetupFileTable(myWindow, "~/Downloads")

	myWindow.SetContent(tableContainer)

	tableContainer.AdjustColumns()

	myWindow.Canvas().SetOnTypedKey(func(event *fyne.KeyEvent) { // Set up keyboard handling for Delete/Backspace
		tableContainer.HandleKeyboard(event)
	})

	CtrlC := &desktop.CustomShortcut{fyne.KeyC, fyne.KeyModifierControl}
	myWindow.Canvas().AddShortcut(CtrlC, func(_ fyne.Shortcut) {
		tableContainer.CopySelectionToClipboard(myApp)
	})

	CtrlA := &desktop.CustomShortcut{fyne.KeyA, fyne.KeyModifierControl}
	myWindow.Canvas().AddShortcut(CtrlA, func(_ fyne.Shortcut) {
		tableContainer.SelectAll()
	})

	myWindow.ShowAndRun()
}
