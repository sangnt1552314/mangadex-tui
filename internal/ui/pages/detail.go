package pages

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/sangnt1552314/mangadex-tui/internal/models"
	"github.com/sangnt1552314/mangadex-tui/internal/ui/interfaces"
)

type DetailPage struct {
	app      interfaces.AppInterface
	rootView *tview.Flex
	manga    *models.Manga
}

type DataReceiver interface {
	SetData(data interface{})
}

func NewDetailPage(app interfaces.AppInterface) *DetailPage {
	return &DetailPage{
		app:      app,
		rootView: tview.NewFlex(),
	}
}

func (p *DetailPage) Name() string {
	return "detail"
}

func (p *DetailPage) View() tview.Primitive {
	return p.rootView
}

func (p *DetailPage) SetManga(manga *models.Manga) {
	p.manga = manga
	// p.updateUI()
}

func (p *DetailPage) Init(app interfaces.AppInterface) {
	p.app = app

	// Functionalities
	app.EnableMouse(true)
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyCtrlC:
			app.Stop()
			return nil
		}
		return event
	})

	// Layout
	p.rootView.SetDirection(tview.FlexRow).
		SetBorder(false)
}
