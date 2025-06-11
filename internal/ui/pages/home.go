package pages

import (
	"log"
	"strconv"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"github.com/sangnt1552314/mangadex-tui/internal/api"
	"github.com/sangnt1552314/mangadex-tui/internal/models"
	"github.com/sangnt1552314/mangadex-tui/internal/ui/interfaces"
	// "github.com/sangnt1552314/mangadex-tui/internal/models"
)

type HomePage struct {
	app      interfaces.AppInterface
	rootView *tview.Flex
}

func NewHomePage(app interfaces.AppInterface) *HomePage {

	return &HomePage{
		app:      app,
		rootView: tview.NewFlex(),
	}
}

func (p *HomePage) Name() string {
	return "home"
}

func (p *HomePage) View() tview.Primitive {
	return p.rootView
}

func (p *HomePage) Init(app interfaces.AppInterface) {
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
	p.rootView.AddItem(mainContent, 0, 8, false)
	p.rootView.AddItem(menu, 0, 1, false)
}

func (p *HomePage) setupMenu() tview.Primitive {
	// Replace List with Flex set to horizontal direction
	menuFlex := tview.NewFlex().SetDirection(tview.FlexColumn)
	menuFlex.SetBackgroundColor(tcell.ColorBlack).SetBorder(true)

	exitButton := tview.NewButton("⏻ Exit")
	exitButton.SetStyle(tcell.StyleDefault.Foreground(tcell.ColorYellow).Background(tcell.ColorBlack))
	exitButton.SetSelectedFunc(func() {
		p.app.Stop()
	})

	aboutButton := tview.NewButton("ℹ About")
	aboutButton.SetStyle(tcell.StyleDefault.Foreground(tcell.ColorYellow).Background(tcell.ColorBlack))
	aboutButton.SetSelectedFunc(func() {
		p.app.SwitchToPage("about")
	})

	// Add buttons to the flex container with equal proportion
	menuFlex.AddItem(aboutButton, 9, 1, false)
	menuFlex.AddItem(exitButton, 9, 1, false)

	return menuFlex
}

func (p *HomePage) setupMainContent() tview.Primitive {
	search := p.setInputSearchComponent()

	current_poplar := tview.NewFlex().SetDirection(tview.FlexRow)
	current_poplar.SetBorder(true).SetTitle("Popular").SetTitleAlign(tview.AlignLeft)
	popular_mangalist := tview.NewTable()
	p.setTableHeaderManga(popular_mangalist)

	mainContent := tview.NewFlex().SetDirection(tview.FlexRow)
	mainContent.SetBorder(false).SetTitleAlign(tview.AlignLeft)

	mainContent.AddItem(search, 0, 1, false)
	mainContent.AddItem(current_poplar, 0, 5, false)

	current_poplar.AddItem(popular_mangalist, 0, 1, false)

	// Initialize manga data
	p.initMangaData(popular_mangalist, 100)

	return mainContent
}

func (p *HomePage) setInputSearchComponent() tview.Primitive {
	search := tview.NewInputField()
	search.SetTitle("Search").SetTitleAlign(tview.AlignLeft)
	search.SetBorder(true)
	search.SetFieldBackgroundColor(tcell.ColorNone).SetFieldTextColor(tcell.ColorWhite)

	return search
}

func (p *HomePage) setTableHeaderManga(manga_list *tview.Table) {
	manga_list.SetCell(0, 0, tview.NewTableCell("Title").
		SetMaxWidth(18).SetSelectable(false).
		SetTextColor(tcell.ColorYellow).
		SetAttributes(tcell.AttrBold))

	manga_list.SetCell(0, 1, tview.NewTableCell("Status").
		SetMaxWidth(18).SetSelectable(false).
		SetTextColor(tcell.ColorYellow).
		SetAttributes(tcell.AttrBold))

	manga_list.SetCell(0, 2, tview.NewTableCell("Year").
		SetMaxWidth(18).SetSelectable(false).
		SetTextColor(tcell.ColorYellow).
		SetAttributes(tcell.AttrBold))
	manga_list.SetFixed(1, 0)
}

func (p *HomePage) initMangaData(manga_list *tview.Table, limit int) {
	mangas, err := api.GetLatestManga(limit)

	if err != nil {
		log.Println("Error fetching manga data:", err)
		return
	}
	for i, manga := range mangas {
		titleCell := tview.NewTableCell(manga.Title).SetReference(&manga)
		manga_list.SetCell(i+1, 0, titleCell)
		manga_list.SetCell(i+1, 1, tview.NewTableCell(manga.Status))
		manga_list.SetCell(i+1, 2, tview.NewTableCell(strconv.Itoa(manga.Year)))
		manga_list.SetSelectedFunc(func(row, _ int) {
			if row == 0 {
				return // Skip header row
			}
			selectedManga := manga_list.GetCell(row, 0).GetReference().(*models.Manga)
			log.Printf("Selected Manga: %s (ID: %s)", selectedManga.Title, selectedManga.ID)
		})
	}

	manga_list.SetSelectable(true, false)
}
