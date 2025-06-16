package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"github.com/sangnt1552314/mangadex-tui/internal/ui/interfaces"
	"github.com/sangnt1552314/mangadex-tui/internal/ui/pages"
)

type App struct {
	*tview.Application
	Pages       *tview.Pages
	pageObjects map[string]interfaces.Page
}

var _ interfaces.AppInterface = (*App)(nil)

func NewApp() *App {
	app := &App{
		Application: tview.NewApplication(),
		Pages:       tview.NewPages(),
		pageObjects: make(map[string]interfaces.Page),
	}

	app.setupBindings()
	app.setupPages()

	return app
}

func (a *App) Run() error {
	return a.Application.Run()
}

func (a *App) Stop() {
	a.Application.Stop()
}

func (a *App) setupBindings() {
	a.Application.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyCtrlC:
			a.Stop()
			return nil
		}
		return event
	})
}

func (a *App) setupPages() {
	a.SetRoot(a.Pages, true)

	a.RegisterPage(pages.NewHomePage(a))
	a.RegisterPage(pages.NewAboutPage(a))
	a.RegisterPage(pages.NewDetailPage(a))
	a.RegisterPage(pages.NewSearchPage(a))

	a.SwitchToPage("home")
}

func (a *App) AddPage(name string, page tview.Primitive, visible bool) {
	a.Pages.AddPage(name, page, true, visible)
}

func (a *App) GetPageObject(name string) interfaces.Page {
	return a.pageObjects[name]
}

func (a *App) SwitchToPage(name string) {
	a.Pages.SwitchToPage(name)
}

func (a *App) RegisterPage(page interfaces.Page) {
	page.Init(a)
	a.AddPage(page.Name(), page.View(), page.Name() == "home")
	a.pageObjects[page.Name()] = page
}

func (a *App) EnableMouse(enable bool) {
	a.Application.EnableMouse(enable)
}

func (a *App) SetInputCapture(fn func(event *tcell.EventKey) *tcell.EventKey) {
	a.Application.SetInputCapture(fn)
}

func (a *App) SetRoot(root tview.Primitive, fullscreen bool) {
	a.Application.SetRoot(root, fullscreen)
}

func (a *App) RestorePages() {
	a.Application.SetRoot(a.Pages, true)
}
