package main

import "mangalib-downloader/service"

func main() {
	service.Start()
	defer service.Stop()
}
