package components

import (
	"mangalib-downloader/components/utils"
	"mangalib-downloader/core"

	"github.com/rivo/tview"
)

type SearchModal struct {
	form  *tview.Form
	modal tview.Primitive
}

func ShowSearchModal() {
	searchModal := newSearchModal()
	searchModal.setHandlers()

	core.App.TView.SetFocus(searchModal.form)
	core.App.PageHolder.AddPage(utils.SearchModalID, searchModal.modal, true, true)
}

func newSearchModal() *SearchModal {
	modal := func(p tview.Primitive, width, height int) tview.Primitive {
		return tview.NewGrid().
			SetColumns(0, width, 0).SetRows(0, height, 0).
			AddItem(p, 1, 1, 1, 1, 0, 0, true)
	}

	form := newSearchForm()
	return &SearchModal{
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
