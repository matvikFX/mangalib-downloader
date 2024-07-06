package components

import (
	"fmt"

	"mangalib-downlaoder/components/utils"
	"mangalib-downlaoder/core"

	"github.com/rivo/tview"
)

type HelpPage struct {
	Grid *tview.Grid
}

func ShowHelpPage() {
	helpPage := newHelpPage()

	core.App.TView.SetFocus(helpPage.Grid)
	core.App.PageHolder.AddAndSwitchToPage(utils.HelpPageID, helpPage.Grid, true)
}

func newHelpPage() *HelpPage {
	textFormat := fmt.Sprintf("%%-%ds:%%%ds\n", 25, 25)

	text := "Сочетания клавиш\n" +
		"-------------------------\n\n" +
		"Универсальные\n" +
		fmt.Sprintf(textFormat, "Ctrl + O", "Помощь") +
		fmt.Sprintf(textFormat, "Ctrl + S", "Поиск") +
		"\nСтраницы Манги\n" +
		fmt.Sprintf(textFormat, "Enter", "Выбрать главу") +
		fmt.Sprintf(textFormat, "Ctrl + D", "Скачать выбранные главы") +
		fmt.Sprintf(textFormat, "Ctrl + A", "Скачать все главы") +
		"\nТаблица Манги\n" +
		fmt.Sprintf(textFormat, "Escape", "Обнуление поиска") +
		fmt.Sprintf(textFormat, "Ctrl + F/B", "След/Пред страница")

	help := tview.NewTextView()
	help.SetText(text).
		SetTextAlign(tview.AlignCenter).
		SetBorder(true)

	grid := tview.NewGrid()
	grid.AddItem(help, 0, 0, 6, 6, 0, 0, true)

	helpPage := &HelpPage{
		Grid: grid,
	}
	helpPage.setHandlers()

	return helpPage
}
