package components

import (
	"manga-downloader/components/utils"

	"github.com/rivo/tview"
)

type SearchModal struct {
	app *TViewApp

	form  *tview.Form
	modal tview.Primitive
}

func (a *TViewApp) ShowSearchModal() {
	searchModal := newSearchModal(a)
	searchModal.setHandlers()

	a.App.SetFocus(searchModal.form)
	a.Pages.AddPage(utils.SearchModalID, searchModal.modal, true, true)
}

func newSearchModal(app *TViewApp) *SearchModal {
	modal := func(p tview.Primitive, width, height int) tview.Primitive {
		return tview.NewGrid().
			SetColumns(0, width, 0).SetRows(0, height, 0).
			AddItem(p, 1, 1, 1, 1, 0, 0, true)
	}

	form := newSearchForm()
	return &SearchModal{
		app: app,

		form:  form,
		modal: modal(form, 100, 5),
	}
}

func newSearchForm() *tview.Form {
	form := tview.NewForm()
	form.AddInputField(utils.SearchModalLabel, "", 87, nil, nil).
		SetTitle("Поиск по названию").SetBorder(true)

	return form
}
