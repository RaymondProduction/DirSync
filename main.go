package main

import (
	"log"
	"os"

	"github.com/getlantern/systray"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

func main() {
	systray.Run(onReady, onExit)
}

func onReady() {

	gtk.Init(nil)
	win := initGTKWindow()

	systray.SetIcon(getIcon("icon.png"))
	systray.SetTitle("System tray")
	systray.SetTooltip("Exaple for system tray")
	mOpen := systray.AddMenuItem("Open Window", "Open Window")
	mQuit := systray.AddMenuItem("Quit", "Quit the whole app")

	systray.AddSeparator()

	go func() {
		select {
		case <-mOpen.ClickedCh:
			glib.IdleAdd(func() {
				win.ShowAll()
			})
			//win.ShowAll()
		case <-mQuit.ClickedCh:
			glib.IdleAdd(func() {
				gtk.MainQuit()
			})
			systray.Quit()
			log.Println("Quit")
		}
	}()
}

func onExit() {
	// clear if needed
}

func getIcon(filePath string) []byte {
	icon, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Error during downloading icon: %v", err)
	}
	return icon
}

func initGTKWindow() *gtk.Window {

	// Create new window of top level
	win, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err != nil {
		log.Fatal("Failed to create window:", err)
	}
	win.SetTitle("Run GTK")
	win.Connect("destroy", gtk.MainQuit)

	// Create container VBox
	box, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	if err != nil {
		log.Fatal("Failed to create container:", err)
	}
	win.Add(box)

	// Create a label
	label, err := gtk.LabelNew("Hello")
	if err != nil {
		log.Fatal("Failed to create a label:", err)
	}
	box.PackStart(label, true, true, 0)

	// Create a button to open a file selection dialog
	button, err := gtk.ButtonNewWithLabel("Вибрати файл")
	if err != nil {
		log.Fatal("Failed to create a button:", err)
	}
	box.PackStart(button, true, true, 0)

	// Add a button click event handler
	button.Connect("clicked", func() {
		dialog, err := gtk.FileChooserDialogNewWith2Buttons("Select file", win, gtk.FILE_CHOOSER_ACTION_OPEN, "Cancel", gtk.RESPONSE_CANCEL, "Select", gtk.RESPONSE_ACCEPT)
		if err != nil {
			log.Fatal("Failed to create dialog box:", err)
		}
		defer dialog.Destroy()

		filter, err := gtk.FileFilterNew()
		if err != nil {
			log.Fatal("Failed to create file filter:", err)
		}
		filter.AddPattern("*.txt")
		filter.SetName("Text files")
		dialog.AddFilter(filter)

		response := dialog.Run()
		if response == gtk.RESPONSE_ACCEPT {
			filename := dialog.GetFilename()
			log.Println("Selected file:", filename)
		}
	})

	// We set the size of the window
	win.SetDefaultSize(300, 200)

	return win
}
