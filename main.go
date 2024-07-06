package main

import (
	"mangalib-downlaoder/service"
)

func main() {
	service.Start()
	defer service.Stop()
}
