package components

import (
	"manga-downloader/components/utils"

	"github.com/gdamore/tcell/v2"
)

func (p *HelpPage) setHandlers() {
	p.Grid.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEscape:
			p.app.Pages.RemovePage(utils.HelpPageID)
		}
		return event
	})
}
