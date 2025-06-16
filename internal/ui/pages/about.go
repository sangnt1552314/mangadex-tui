package pages

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/sangnt1552314/mangadex-tui/internal/ui/interfaces"
)

type AboutPage struct {
	app      interfaces.AppInterface
	rootView *tview.Flex
}

func NewAboutPage(app interfaces.AppInterface) *AboutPage {
	return &AboutPage{
		app:      app,
		rootView: tview.NewFlex(),
	}
}

func (p *AboutPage) Name() string {
	return "about"
}

func (p *AboutPage) View() tview.Primitive {
	return p.rootView
}

func (p *AboutPage) Init(app interfaces.AppInterface) {
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

	// Layout - Content
	content := p.setupContent()
	p.rootView.AddItem(content, 0, 9, true)

	// Layout - Menu
	menu := p.setupMenu()
	p.rootView.AddItem(menu, 0, 1, false)
}

func (p *AboutPage) setupMenu() tview.Primitive {
	menuFlex := tview.NewFlex().SetDirection(tview.FlexColumn)
	menuFlex.SetBackgroundColor(tcell.ColorBlack).SetBorder(true).SetTitle("Options").SetTitleAlign(tview.AlignLeft)

	homeButton := tview.NewButton("⌂ Home")
	homeButton.SetStyle(tcell.StyleDefault.Foreground(tcell.ColorDodgerBlue).Background(tcell.ColorBlack))
	homeButton.SetSelectedFunc(func() {
		p.app.SwitchToPage("home")
	})

	searchButton := tview.NewButton("🔍 Search")
	searchButton.SetStyle(tcell.StyleDefault.Foreground(tcell.ColorPurple).Background(tcell.ColorBlack))
	searchButton.SetSelectedFunc(func() {
		p.app.SwitchToPage("search")
	})

	exitButton := tview.NewButton("⏻ Exit")
	exitButton.SetStyle(tcell.StyleDefault.Foreground(tcell.ColorRed).Background(tcell.ColorBlack))
	exitButton.SetSelectedFunc(func() {
		p.app.Stop()
	})

	// Add buttons to the flex container with equal proportion
	menuFlex.AddItem(homeButton, 9, 1, false)
	menuFlex.AddItem(searchButton, 9, 1, false)
	menuFlex.AddItem(exitButton, 9, 1, false)

	return menuFlex
}

func (p *AboutPage) setupContent() tview.Primitive {
	mainContent := tview.NewFlex().SetDirection(tview.FlexRow)
	mainContent.SetBackgroundColor(tcell.ColorBlack).
		SetBorder(true).
		SetTitle("About MangaDex TUI").
		SetTitleColor(tcell.ColorGreen)
	mainContent.SetBorderPadding(1, 1, 2, 2)
	mainContent.AddItem(
		tview.NewTextView().
			SetText("MangaDex TUI is a terminal-based client for MangaDex.\n\n"+
				"Developed by Sang Nguyen.\n\n"+
				"Visit the project on GitHub: https://github.com/sangnt1552314/mangadex-tui\n\n"+
				"Support MangaDex at https://mangadex.org/").SetTextColor(tcell.ColorGreen).
			SetDynamicColors(true),
		0, 1, false,
	)
	mainContent.SetBorderPadding(1, 1, 2, 2)
	return mainContent
}
