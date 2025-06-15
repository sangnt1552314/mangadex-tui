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
	p.updateUI()
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

	p.updateUI()
}

func (p *DetailPage) updateUI() {
	p.rootView.Clear()

	p.rootView.SetDirection(tview.FlexRow).
		SetBorder(false)

	menu := p.setupMenu()

	mainContent := p.setupMainContent()

	p.rootView.AddItem(mainContent, 0, 9, false)
	p.rootView.AddItem(menu, 0, 1, false)
}

func (p *DetailPage) setupMenu() tview.Primitive {
	menuFlex := tview.NewFlex().SetDirection(tview.FlexColumn)
	menuFlex.SetBackgroundColor(tcell.ColorBlack).SetBorder(true).SetTitle("Options").SetTitleAlign(tview.AlignLeft)

	homeButton := tview.NewButton("⌂ Home")
	homeButton.SetStyle(tcell.StyleDefault.Foreground(tcell.ColorBlue).Background(tcell.ColorBlack))
	homeButton.SetSelectedFunc(func() {
		p.app.SwitchToPage("home")
	})

	exitButton := tview.NewButton("⏻ Exit")
	exitButton.SetStyle(tcell.StyleDefault.Foreground(tcell.ColorRed).Background(tcell.ColorBlack))
	exitButton.SetSelectedFunc(func() {
		p.app.Stop()
	})

	aboutButton := tview.NewButton("ℹ About")
	aboutButton.SetStyle(tcell.StyleDefault.Foreground(tcell.ColorGreen).Background(tcell.ColorBlack))
	aboutButton.SetSelectedFunc(func() {
		p.app.SwitchToPage("about")
	})

	// Add buttons to the flex container with equal proportion
	menuFlex.AddItem(homeButton, 9, 1, false)
	menuFlex.AddItem(aboutButton, 9, 1, false)
	menuFlex.AddItem(exitButton, 9, 1, false)

	return menuFlex
}

func (p *DetailPage) setupMainContent() tview.Primitive {
	if p.manga == nil {
		noMangaText := tview.NewTextView().SetText("No manga selected")
		return noMangaText
	}

	mainContent := tview.NewFlex().SetDirection(tview.FlexRow)
	mainContent.SetBorder(true).SetTitle(p.manga.Attributes.Title["en"]).SetTitleAlign(tview.AlignLeft)

	return mainContent
}
