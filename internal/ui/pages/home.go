package pages

import (
	"log"
	"strconv"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"github.com/sangnt1552314/mangadex-tui/internal/api"
	"github.com/sangnt1552314/mangadex-tui/internal/models"
	"github.com/sangnt1552314/mangadex-tui/internal/ui/interfaces"
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
	p.rootView.AddItem(mainContent, 0, 9, false)
	p.rootView.AddItem(menu, 0, 1, false)
}

func (p *HomePage) setupMenu() tview.Primitive {
	// Replace List with Flex set to horizontal direction
	menuFlex := tview.NewFlex().SetDirection(tview.FlexColumn)
	menuFlex.SetBackgroundColor(tcell.ColorBlack).SetBorder(true).SetTitle("Menu").SetTitleAlign(tview.AlignLeft)

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
	mainContent := tview.NewFlex().SetDirection(tview.FlexRow)
	mainContent.SetBorder(false).SetTitleAlign(tview.AlignLeft)

	// Search Component
	searchBox := p.setInputSearchComponent()

	// Popular Flex
	popularFlex := tview.NewFlex().SetDirection(tview.FlexRow)
	popularFlex.SetBorder(true).SetTitle("Popular").SetTitleAlign(tview.AlignLeft)

	// Feature & Latest Manga Flex
	featureLatestFlex := tview.NewFlex().SetDirection(tview.FlexColumn)

	featureFlex := tview.NewFlex().SetDirection(tview.FlexRow)
	featureFlex.SetBorder(true).SetTitle("Featured").SetTitleAlign(tview.AlignLeft)

	latestFlex := tview.NewFlex().SetDirection(tview.FlexRow)
	latestFlex.SetBorder(true).SetTitle("Latest").SetTitleAlign(tview.AlignLeft)

	featureLatestFlex.AddItem(featureFlex, 0, 1, false)
	featureLatestFlex.AddItem(latestFlex, 0, 1, false)

	// Manga List Table
	featureMangaList := tview.NewTable()
	p.setTableHeaderManga(featureMangaList)
	p.getMangas(featureMangaList, 10, "feature")
	featureFlex.AddItem(featureMangaList, 0, 1, false)

	latestMangaList := tview.NewTable()
	p.setTableHeaderManga(latestMangaList)
	p.getMangas(latestMangaList, 10, "latest")
	latestFlex.AddItem(latestMangaList, 0, 1, false)

	mainContent.AddItem(searchBox, 0, 1, false)
	mainContent.AddItem(popularFlex, 0, 4, false)
	mainContent.AddItem(featureLatestFlex, 0, 4, false)

	return mainContent
}

func (p *HomePage) setInputSearchComponent() tview.Primitive {
	search := tview.NewInputField()
	search.SetTitle("Search").SetTitleAlign(tview.AlignLeft)
	search.SetBorder(true)
	search.SetFieldBackgroundColor(tcell.ColorNone).SetFieldTextColor(tcell.ColorWhite)

	return search
}

func (p *HomePage) setTableHeaderManga(mangaList *tview.Table) {
	mangaList.SetCell(0, 0, tview.NewTableCell("Title").
		SetSelectable(false).
		SetTextColor(tcell.ColorOrange))

	mangaList.SetCell(0, 1, tview.NewTableCell("Status").
		SetSelectable(false).
		SetTextColor(tcell.ColorYellow))

	mangaList.SetCell(0, 2, tview.NewTableCell("Year").
		SetSelectable(false).
		SetTextColor(tcell.ColorDarkOrange))

	mangaList.SetFixed(1, 0)
}

func (p *HomePage) getMangas(mangaList *tview.Table, limit int, mangaType string) {
	mangas, err := api.GetManga(limit, mangaType)

	if err != nil {
		log.Println("Error fetching manga data:", err)
		return
	}

	for i, manga := range mangas {
		titleCell := tview.NewTableCell(manga.Title).SetReference(&manga).SetMaxWidth(30)
		mangaList.SetCell(i+1, 0, titleCell)
		mangaList.SetCell(i+1, 1, p.formatStatus(manga.Status))
		mangaList.SetCell(i+1, 2, tview.NewTableCell(strconv.Itoa(manga.Year)))
	}

	mangaList.SetSelectedFunc(func(row, _ int) {
		if row == 0 {
			return // Skip header row
		}
		selectedManga := mangaList.GetCell(row, 0).GetReference().(*models.Manga)
		log.Printf("Selected Manga: %s (ID: %s)", selectedManga.Title, selectedManga.ID)
	})

	mangaList.SetSelectable(true, false)
}

func (p *HomePage) formatStatus(status string) *tview.TableCell {
	switch status {
	case "ongoing":
		return tview.NewTableCell("Ongoing").SetTextColor(tcell.ColorGreen)
	case "completed":
		return tview.NewTableCell("Completed").SetTextColor(tcell.ColorBlue)
	case "hiatus":
		return tview.NewTableCell("Hiatus").SetTextColor(tcell.ColorYellow)
	case "cancelled":
		return tview.NewTableCell("Cancelled").SetTextColor(tcell.ColorRed)
	default:
		return tview.NewTableCell(status).SetTextColor(tcell.ColorWhite)
	}
}
