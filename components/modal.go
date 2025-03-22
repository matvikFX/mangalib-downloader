package components

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func (a *TViewApp) ShowModal(id, text string) {
	modal := tview.NewModal()
	modal.SetText(text).
		SetBackgroundColor(tcell.ColorBlack).
		AddButtons([]string{"Ok"}).
		SetFocus(0).
		SetDoneFunc(func(_ int, _ string) {
			a.Pages.RemovePage(id)
		})

	a.App.SetFocus(modal)
	a.Pages.AddPage(id, modal, true, true)
}
