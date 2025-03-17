package components

import (
	"context"
	"time"

	"manga-downloader/components/utils"
	"manga-downloader/models"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Pages interface {
	SetFocus(p tview.Primitive)
	AddAndSwitchToPage(id string, item tview.Primitive, resize bool)
	ChangePath(path string)
	ChangePage(d_page int)
}

var (
	timer *time.Timer

	selectedManga = &models.MangaInfo{}
)

func (t *TViewApp) SetHandlers() {
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
			t.App.Stop()
		}

		return event
	})
}

func (p *ListPage) setHandlers(ctx context.Context, cancel context.CancelFunc) {
	p.table.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		reload := func() {
			cancel()
			go p.setListTable()
		}

		if event.Rune() == ' ' {
			if timer != nil {
				timer.Stop()
			}
			p.app.ShowBranchModal(ctx)
		}

		switch event.Key() {
		case tcell.KeyEscape: // Обнулить поисковую строку и вернуться на первую страницу
			if p.app.Downloader.Query != "" || p.app.Downloader.Page != 1 {
				p.app.Downloader.Query = ""
				p.app.Downloader.Page = 1
				reload()
			}
		case tcell.KeyCtrlF: // Предыдущая страница
			p.app.Downloader.Page++
			reload()
		case tcell.KeyCtrlB: // Следующая страница
			if p.app.Downloader.Page == 1 {
				p.app.ShowModal(utils.NoMangaID, "Ниже первой страницы опуститься нельзя")
				break
			}
			p.app.Downloader.Page--
			reload()
		}
		return event
	})

	p.table.SetSelectedFunc(func(_, _ int) {
		if timer != nil {
			timer.Stop()
		}
		p.app.ShowBranchModal(ctx)
	})

	p.table.SetSelectionChangedFunc(func(row, column int) {
		manga := p.getMangaFromCell(row)
		selectedManga = &models.MangaInfo{
			Manga: *manga,
		}

		p.textView.SetTitle("Загрузка информации о манге...")
		p.textView.SetText("")

		if timer != nil {
			timer.Stop()
		}
		timer = time.AfterFunc(800*time.Millisecond, func() {
			info, err := p.app.Downloader.GetInfo(ctx, manga.Slug)
			if err != nil {
				p.app.Logger.Write(err.Error())
				return
			}
			branches, err := p.app.Downloader.GetMangaBranches(ctx, selectedManga.ID)
			if err != nil {
				p.app.Logger.Write(err.Error())
				return
			}
			info.Branches = branches
			selectedManga = info

			infoText := utils.InfoText(info, nil)
			p.app.App.QueueUpdateDraw(func() {
				p.textView.SetTitle("Информация о манге")
				p.textView.SetText(infoText)
			})
		})
	})
}

func (p *MangaPage) setHandlers(ctx context.Context, cancel context.CancelFunc) {
	select_change_row_color := func(row int) {
		// Я не знаю почему только так работает выделение нескольких столбцов
		// Пока оставлю так, если найду способ лучше, поменяю
		cols := []int{0, 1, 1, 2, 2}
		for _, col := range cols {
			cell := p.table.GetCell(row, col)
			if p.selected[row] {
				cell.SetBackgroundColor(tcell.ColorBlack)
				delete(p.selected, row)
			} else {
				cell.SetBackgroundColor(tcell.ColorRed)
				p.selected[row] = true
			}
			p.table.SetCell(row, col, cell)
		}
	}

	p.grid.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == ' ' {
			row, _ := p.table.GetSelection()
			select_change_row_color(row)
		}

		switch event.Key() {
		case tcell.KeyEscape: // Выход со страницы манги
			p.app.Downloader.Branch = 0
			timer.Reset(1 * time.Second)
			p.app.Pages.RemovePage(utils.MangaPageID)
			cancel()
		case tcell.KeyCtrlD: // Скачивание выделенных
			if len(p.selected) != 0 {
				go p.downloadSelected(ctx)
			}
		case tcell.KeyCtrlA: // Скачивание всех глав
			go func() {
				p.app.Downloader.DownloadManga(ctx, selectedManga)

				<-p.app.Downloader.Downloaded
				p.app.ShowModal(utils.DownloadSuccessID,
					"Манга '"+selectedManga.RusName+"' успешно скачана")

				go p.setChapters(ctx)
			}()
		case tcell.KeyCtrlT: // Выбор ветки перевода
			if len(selectedManga.Branches) > 0 {
				p.app.ShowBranchModal(ctx)
			} else {
				p.app.ShowModal(utils.NoBranchesID,
					"У данной менги нет другой ветки перевода")
			}
		}
		return event
	})

	// Выбор глав для скачивания
	p.table.SetSelectedFunc(func(row, _ int) {
		select_change_row_color(row)
	})
}

func (p *SearchModal) setHandlers() {
	p.form.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEscape: // Закрытие страницы поиска
			p.app.Pages.RemovePage(utils.SearchModalID)
		case tcell.KeyEnter:
			searchInput := p.form.GetFormItemByLabel(utils.SearchModalLabel).(*tview.InputField)
			formText := searchInput.GetText()
			p.app.Downloader.Query = formText
			p.app.Downloader.Page = 1

			searchInput.SetText("")
			p.app.Pages.RemovePage(utils.SearchModalID)
			p.app.ShowListPage()
		}
		return event
	})
}

func (p *PathModal) setHandlers() {
	p.form.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEscape:
			p.app.Pages.RemovePage(utils.PathsModalID)
		}
		return event
	})
}

func (p *HelpPage) setHandlers() {
	p.Grid.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEscape:
			p.app.Pages.RemovePage(utils.HelpPageID)
		}
		return event
	})
}
