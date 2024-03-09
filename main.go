package main

import (
	"fmt"
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
	systray.SetTooltip("Exaple for system tray")
	mOpen := systray.AddMenuItem("Open Window", "Open Window")
	mQuit := systray.AddMenuItem("Quit", "Quit the whole app")

	systray.AddSeparator()

	go func() {
		for {
			select {
			case <-mOpen.ClickedCh:
				glib.IdleAdd(func() {
					win.ShowAll()
				})
			case <-mQuit.ClickedCh:
				glib.IdleAdd(func() {
					win.Destroy()
					gtk.MainQuit()
				})
				systray.Quit()
				log.Println("Quit")
				return // exit the goroutine after the program is finished
			}
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

	// Create builder
	b, err := gtk.BuilderNew()
	if err != nil {
		log.Fatal("Error bulder:", err)
	}

	// Lload the window from the Glade file into the builder
	err = b.AddFromFile("main.glade")
	if err != nil {
		log.Fatal("Error when loading glade file:", err)
	}

	// We get the object of the main window by ID
	obj, err := b.GetObject("setting-window")
	if err != nil {
		log.Fatal("Error:", err)
	}

	win := obj.(*gtk.Window)

	// We get the object of the main window by ID
	objOpenFolder, err := b.GetObject("open_folder")
	if err != nil {
		log.Fatal("Error:", err)
	}

	button := objOpenFolder.(*gtk.Button)

	objPath, err := b.GetObject("path")
	if err != nil {
		log.Fatal("Error:", err)
	}

	entry := objPath.(*gtk.Entry)

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
		filter.AddPattern("*.*")
		filter.SetName("Text files")
		dialog.AddFilter(filter)

		response := dialog.Run()
		if response == gtk.RESPONSE_ACCEPT {
			filename := dialog.GetFilename()
			log.Println("Selected file:", filename)
			entry.SetText(filename)
		}
	})

	win.Connect("destroy", func() {
		fmt.Println("Destroy")
	})

	win.Connect("delete-event", func() bool {
		win.Hide()  // Hide the window.
		return true // Returning true prevents further propagation of the signal and stops the window from closing.
	})

	return win
}
