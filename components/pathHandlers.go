package components

import (
	"manga-downloader/components/utils"

	"github.com/gdamore/tcell/v2"
)

func (p *PathModal) setHandlers() {
	p.form.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEscape:
			p.app.Pages.RemovePage(utils.PathsModalID)
		}
		return event
	})
}
