package pages

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/sangnt1552314/mangadex-tui/internal/models"
	"github.com/sangnt1552314/mangadex-tui/internal/ui/interfaces"
)

type ReaderPage struct {
	app      interfaces.AppInterface
	rootView *tview.Flex
	manga    *models.Manga
	chapter  *models.Chapter
}

func NewReaderPage(app interfaces.AppInterface) *ReaderPage {
	return &ReaderPage{
		app:      app,
		rootView: tview.NewFlex(),
		manga:    nil,
		chapter:  nil,
	}
}

func (p *ReaderPage) Name() string {
	return "reader"
}

func (p *ReaderPage) View() tview.Primitive {
	return p.rootView
}

func (p *ReaderPage) SetManga(manga *models.Manga) {
	p.manga = manga
}

func (p *ReaderPage) SetChapter(chapter *models.Chapter) {
	p.chapter = chapter
}

func (p *ReaderPage) SetData(manga *models.Manga, chapter *models.Chapter) {
	p.manga = manga
	p.chapter = chapter
	p.updateUI()
}

func (p *ReaderPage) Init(app interfaces.AppInterface) {
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

	p.updateUI()
}

func (p *ReaderPage) updateUI() {
	if p.manga == nil || p.chapter == nil {
		return
	}

	// Layout
	p.rootView.SetDirection(tview.FlexRow).
		SetBorder(false)

	// Layout - Main Content
	mainContent := p.setupMainContent()

	// Layout - Menu
	menu := p.setupMenu()

	// Add components to the root view
	p.rootView.AddItem(mainContent, 0, 1, true)
	p.rootView.AddItem(menu, 3, 0, false)

}

func (p *ReaderPage) setupMenu() tview.Primitive {
	menuFlex := tview.NewFlex().SetDirection(tview.FlexColumn)
	menuFlex.SetBackgroundColor(tcell.ColorBlack).SetBorder(true).SetTitle("Options").SetTitleAlign(tview.AlignLeft)

	homeButton := tview.NewButton("⌂ Home")
	homeButton.SetSelectedFunc(func() {
		p.app.SwitchToPage("home")
	})

	menuFlex.AddItem(homeButton, 0, 1, false)

	return menuFlex
}

func (p *ReaderPage) setupMainContent() tview.Primitive {
	mainContent := tview.NewTextView().
		SetText("Reader Page Content").
		SetTextAlign(tview.AlignCenter).
		SetDynamicColors(true).
		SetBorder(true).
		SetTitle("Reader")

	mainContent.SetBackgroundColor(tcell.ColorBlack)
	mainContent.SetBorderColor(tcell.ColorWhite)

	return mainContent
}
