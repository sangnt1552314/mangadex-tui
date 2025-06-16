package pages

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/sangnt1552314/mangadex-tui/internal/ui/interfaces"
)

type SearchPage struct {
	app      interfaces.AppInterface
	rootView *tview.Flex
}

func NewSearchPage(app interfaces.AppInterface) *SearchPage {
	return &SearchPage{
		app:      app,
		rootView: tview.NewFlex(),
	}
}

func (p *SearchPage) Name() string {
	return "search"
}

func (p *SearchPage) View() tview.Primitive {
	return p.rootView
}

func (p *SearchPage) Init(app interfaces.AppInterface) {
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

	// Layout - Main Content
	mainContent := p.setupMainContent()

	// Layout - Menu
	menu := p.setupMenu()

	// Add components to the root view
	p.rootView.AddItem(mainContent, 0, 9, true)
	p.rootView.AddItem(menu, 0, 1, false)
}

func (p *SearchPage) setupMenu() tview.Primitive {
	menuFlex := tview.NewFlex().SetDirection(tview.FlexColumn)
	menuFlex.SetBackgroundColor(tcell.ColorBlack).SetBorder(true).SetTitle("Options").SetTitleAlign(tview.AlignLeft)

	homeButton := tview.NewButton("⌂ Home")
	homeButton.SetStyle(tcell.StyleDefault.Foreground(tcell.ColorDodgerBlue).Background(tcell.ColorBlack))
	homeButton.SetSelectedFunc(func() {
		p.app.SwitchToPage("home")
	})

	aboutButton := tview.NewButton("ℹ About")
	aboutButton.SetStyle(tcell.StyleDefault.Foreground(tcell.ColorGreen).Background(tcell.ColorBlack))
	aboutButton.SetSelectedFunc(func() {
		p.app.SwitchToPage("about")
	})

	exitButton := tview.NewButton("⏻ Exit")
	exitButton.SetStyle(tcell.StyleDefault.Foreground(tcell.ColorRed).Background(tcell.ColorBlack))
	exitButton.SetSelectedFunc(func() {
		p.app.Stop()
	})

	// Add buttons to the flex container with equal proportion
	menuFlex.AddItem(homeButton, 9, 1, false)
	menuFlex.AddItem(aboutButton, 9, 1, false)
	menuFlex.AddItem(exitButton, 9, 1, false)

	return menuFlex
}

func (p *SearchPage) setupMainContent() tview.Primitive {
	mainContent := tview.NewFlex().SetDirection(tview.FlexRow)
	mainContent.SetBorder(false).SetTitleAlign(tview.AlignLeft)

	// Search Component
	searchBox := p.setInputSearchComponent()

	mainContent.AddItem(searchBox, 0, 1, false)

	return mainContent
}

func (p *SearchPage) setInputSearchComponent() tview.Primitive {
	search := tview.NewInputField()
	search.SetTitle("Search").SetTitleAlign(tview.AlignLeft)
	search.SetBorder(true)
	search.SetFieldBackgroundColor(tcell.ColorNone).SetFieldTextColor(tcell.ColorWhite)

	return search
}
