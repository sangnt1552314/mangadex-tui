package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// App represents the main application UI
type App struct {
	*tview.Application
	pages    *tview.Pages
	menu     *tview.List
	content  *tview.TextView
	mainFlex *tview.Flex

	// Pages
	searchPage   *SearchPage
	settingsPage *SettingsPage
}

// NewApp creates a new instance of the application UI
func NewApp() *App {
	app := &App{
		Application: tview.NewApplication(),
		pages:       tview.NewPages(),
		menu:        tview.NewList().ShowSecondaryText(true),
		content:     tview.NewTextView().SetDynamicColors(true).SetWrap(true),
	}

	// Create pages
	app.searchPage = NewSearchPage(app)
	app.settingsPage = NewSettingsPage(app)

	// Initialize pages
	app.setupPages()

	// Set root and enable mouse
	app.SetRoot(app.pages, true)

	return app
}

// Run starts the application
func (a *App) Run() error {
	return a.Application.Run()
}

// switchToPage switches the content area to the specified primitive
func (a *App) switchToPage(p tview.Primitive, title string) {
	a.content.Clear()
	a.content.SetTitle(title)
	a.content.SetBorder(true)
	a.mainFlex.RemoveItem(a.mainFlex.GetItem(1))
	a.mainFlex.AddItem(p, 0, 3, true)
	a.SetFocus(p)
}

// setupPages initializes all the application pages
func (a *App) setupPages() {
	// Setup menu
	a.menu.SetBorder(true).
		SetTitle("Menu").
		SetTitleAlign(tview.AlignLeft)

	a.menu.AddItem("Search Manga", "Search for manga by title or other fields", 'S', func() {
		a.switchToPage(a.searchPage, "Search")
	}).
		AddItem("My Library", "View your saved manga", 'L', nil).
		AddItem("Settings", "Configure application settings", 'C', func() {
			a.switchToPage(a.settingsPage, "Settings")
		}).
		AddItem("Quit", "Press to exit", 'Q', func() {
			a.Stop()
		})

	// Setup content
	a.content.SetBorder(true).
		SetTitle("Welcome").
		SetTitleAlign(tview.AlignLeft)

	a.content.SetText("Welcome to MangaDex TUI!\n\n" +
		"Use the menu on the left to navigate:\n\n" +
		"• Search for manga using various filters\n" +
		"• Save manga to your library\n" +
		"• Configure application settings\n\n" +
		"Press Ctrl+C to quit at any time")

	a.content.SetTextAlign(tview.AlignLeft)

	// Create main layout
	a.mainFlex = tview.NewFlex().
		AddItem(a.menu, 0, 1, true).    // Menu takes up 1/4 of the space
		AddItem(a.content, 0, 3, false) // Content takes up 3/4 of the space

	// Add the flex layout to pages
	a.pages.AddPage("main", a.mainFlex, true, true)

	// Add global keyboard shortcuts
	a.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyCtrlC:
			a.Stop()
			return nil
		case tcell.KeyTab:
			// Toggle focus between menu and content
			if a.menu.HasFocus() {
				currentContent := a.mainFlex.GetItem(1)
				a.SetFocus(currentContent)
			} else {
				a.SetFocus(a.menu)
			}
			return nil
		case tcell.KeyEsc:
			// Return to welcome page
			currentContent := a.mainFlex.GetItem(1)
			if currentContent != a.content {
				a.mainFlex.RemoveItem(currentContent)
				a.mainFlex.AddItem(a.content, 0, 3, false)
				a.SetFocus(a.menu)
			}
			return nil
		}
		return event
	})
}
