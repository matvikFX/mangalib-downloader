package components

import (
	"github.com/gdamore/tcell/v2"
)

func (t *TViewApp) setHandlers() {
	t.App.EnableMouse(true)
	t.App.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Rune() {
		case 'P': // Настройки
			t.ShowPathModal()
		case 'H': // Помощь
			t.ShowHelpPage()
		case 'B': // Закладки
			// t.ShowBookmarksModal()
		}

		switch event.Key() {
		case tcell.KeyCtrlS: // Поиск по названию
			t.ShowSearchModal()
		case tcell.KeyCtrlC: // Завершение работы
			t.Stop()
		}

		return event
	})
}
