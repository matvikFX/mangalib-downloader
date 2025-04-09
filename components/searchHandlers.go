package components

import (
	"context"
	"log/slog"
	"strings"

	"mangalib-downloader/api"
	"mangalib-downloader/components/utils"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func (p *SearchModal) setHandlers(ctx context.Context) {
	p.form.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		log := slog.With("SearchModal", "SetInputCapture")

		switch event.Key() {
		case tcell.KeyEscape: // Закрытие страницы поиска
			p.app.Pages.RemovePage(utils.SearchModalID)
		case tcell.KeyEnter:
			searchInput := p.form.GetFormItemByLabel(utils.SearchModalLabel).(*tview.InputField)
			formText := searchInput.GetText()

			if strings.HasPrefix(formText, "http://") || strings.HasPrefix(formText, "https://") {
				id, slug := utils.ParseURL(formText)
				branches, err := api.GetMangaBranches(ctx, id)
				if err != nil {
					log.Error("Error getting manga branches", "Error", err)
					log.Warn("Неправильная ссылка")
					break
				}

				log.Debug("Manga branches",
					"length", len(branches),
					"branches", branches.GetTeams(),
				)
				log.Debug("Manga",
					"id", id,
					"slug", slug,
				)

				p.app.Pages.RemovePage(utils.SearchModalID)
				if len(branches) == 0 {
					p.app.ShowMangaPage(ctx, slug)
				} else {
					p.app.ShowBranchModal(ctx, slug, branches)
				}

				break
			}

			p.app.query = formText
			p.app.page = 1

			searchInput.SetText("")
			p.app.Pages.RemovePage(utils.SearchModalID)
			p.app.ShowListPage()
		}
		return event
	})
}
