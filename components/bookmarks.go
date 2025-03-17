package components

import (
	"github.com/rivo/tview"
)

func ShowBookmarksModal() {
	modal := func(p tview.Primitive, width, height int) tview.Primitive {
		return tview.NewGrid().
			SetColumns(0, width, 0).SetRows(0, height, 0).
			AddItem(p, 1, 1, 1, 1, 0, 0, true)
	}

	_ = modal
	// form := newBranchForm(ctx, bookmarks)
}
