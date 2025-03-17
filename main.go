package main

import "manga-downloader/components"

func main() {
	app := components.NewTViewApp()

	app.Start()
	defer app.Stop()
}
