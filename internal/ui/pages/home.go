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
	app         interfaces.AppInterface
	rootView    *tview.Flex
	mainContent *tview.Flex
}

func NewHomePage(app interfaces.AppInterface) *HomePage {
	mainContent := tview.NewFlex().SetDirection(tview.FlexRow)
	mainContent.SetBorder(false).SetTitleAlign(tview.AlignLeft)

	return &HomePage{
		app:         app,
		rootView:    tview.NewFlex(),
		mainContent: mainContent,
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
		case 'h':
			app.SwitchToPage("home")
			return nil
		}
		return event
	})

	// Layout
	p.rootView.SetDirection(tview.FlexColumn).
		SetBorder(false)

	menu := tview.NewFlex().SetDirection(tview.FlexRow)
	menu_list := tview.NewList()
	// donate := tview.NewTextView().SetText("Support MangaDex at https://mangadex.org/").SetTextAlign(tview.AlignCenter)

	p.rootView.AddItem(menu, 0, 1, false)
	p.rootView.AddItem(p.mainContent, 0, 8, false)

	// Layout - Menu
	menu.SetBorder(true).
		SetTitle("Menu").
		SetTitleAlign(tview.AlignLeft).
		SetTitleColor(tview.Styles.PrimaryTextColor)
	menu_list.AddItem("⌂ Home", "", 'h', func() {
		app.SwitchToPage("home")
	})
	menu_list.AddItem("⏻ Exit ", "", 'q', func() {
		app.Stop()
	})
	menu_list.AddItem("ℹ About", "", 'a', func() {
		app.SwitchToPage("about")
	})
	menu.AddItem(menu_list, 0, 8, false)

	// Layout - Main Content
	p.setupMainContent()
}

func (p *HomePage) setupMainContent() {
	search := p.setInputSearchComponent()

	current_poplar := tview.NewFlex().SetDirection(tview.FlexRow)
	current_poplar.SetBorder(true).SetTitle("Popular").SetTitleAlign(tview.AlignLeft)
	popular_mangalist := tview.NewTable()
	p.setTableHeaderManga(popular_mangalist)

	p.mainContent.AddItem(search, 0, 1, false)
	p.mainContent.AddItem(current_poplar, 0, 5, false)

	current_poplar.AddItem(popular_mangalist, 0, 1, false)

	// Initialize manga data
	p.initMangaData(popular_mangalist, 100)
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
