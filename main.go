package main

import (
	"log"
	"os"

	"github.com/getlantern/systray"
)

func main() {
	systray.Run(onReady, onExit)
}

func onReady() {
	systray.SetIcon(getIcon("icon.png"))
	systray.SetTitle("System tray")
	systray.SetTooltip("Exaple for system tray")
	mQuit := systray.AddMenuItem("Quit", "Quit from program")

	go func() {
		<-mQuit.ClickedCh
		systray.Quit()
		log.Println("Quit")
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
