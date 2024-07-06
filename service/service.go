package service

import (
	"log"

	"mangalib-downlaoder/components"
	"mangalib-downlaoder/core"
)

func Start() {
	core.App = core.NewApp()
	core.App.Init()

	components.ShowListPage()
	components.SetHandlers()

	log.Println("Starting app")
	if err := core.App.TView.Run(); err != nil {
		log.Println(err)
	}
}

func Stop() { core.App.Close() }
