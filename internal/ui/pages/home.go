package pages

import (
	"fmt"
	"log"
	"strconv"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"github.com/sangnt1552314/mangadex-tui/internal/api"
	"github.com/sangnt1552314/mangadex-tui/internal/models"
	"github.com/sangnt1552314/mangadex-tui/internal/services"
	"github.com/sangnt1552314/mangadex-tui/internal/ui/interfaces"
)

type HomePage struct {
	app            interfaces.AppInterface
	rootView       *tview.Flex
	popularContent *tview.Flex
}

func NewHomePage(app interfaces.AppInterface) *HomePage {

	return &HomePage{
		app:            app,
		rootView:       tview.NewFlex(),
		popularContent: tview.NewFlex(),
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
	menuFlex.AddItem(aboutButton, 9, 1, false)
	menuFlex.AddItem(exitButton, 9, 1, false)

	return menuFlex
}

func (p *HomePage) setupMainContent() tview.Primitive {
	mainContent := tview.NewFlex().SetDirection(tview.FlexRow)
	mainContent.SetBorder(false).SetTitleAlign(tview.AlignLeft)

	// Search Component
	// searchBox := p.setInputSearchComponent()

	// Popular Flex
	popularFlex := p.setupPoplarFlex(tview.NewFlex().SetDirection(tview.FlexRow))

	// Feature & Latest Manga Flex
	featureLatestFlex := tview.NewFlex().SetDirection(tview.FlexColumn)

	featureFlex := tview.NewFlex().SetDirection(tview.FlexRow)
	featureFlex.SetBorder(true).SetTitle("Featured").SetTitleAlign(tview.AlignLeft)

	latestFlex := tview.NewFlex().SetDirection(tview.FlexRow)
	latestFlex.SetBorder(true).SetTitle("Latest").SetTitleAlign(tview.AlignLeft)

	featureLatestFlex.AddItem(featureFlex, 0, 1, false)
	featureLatestFlex.AddItem(latestFlex, 0, 1, false)

	// Manga List Tables
	featureParams := models.MangaQueryParams{
		Limit: 9,
		Order: map[string]string{
			models.OrderByFollowCount: "desc",
		},
	}
	featureMangaList := tview.NewTable()
	p.setTableHeaderManga(featureMangaList)
	featureFlex.AddItem(featureMangaList, 0, 1, false)
	p.setMangaListData(featureMangaList, featureParams)

	latestParams := models.MangaQueryParams{
		Limit: 9,
		Order: map[string]string{
			models.OrderByCreatedAt: "desc",
		},
	}
	latestMangaList := tview.NewTable()
	p.setTableHeaderManga(latestMangaList)
	p.setMangaListData(latestMangaList, latestParams)
	latestFlex.AddItem(latestMangaList, 0, 1, false)

	//Setup Components
	// mainContent.AddItem(searchBox, 0, 1, false)
	mainContent.AddItem(popularFlex, 0, 6, false)
	mainContent.AddItem(featureLatestFlex, 0, 4, false)

	return mainContent
}
func (p *HomePage) setupPoplarFlex(popularFlex *tview.Flex) tview.Primitive {
	limit := 5
	popularParams := models.MangaQueryParams{
		Limit: limit,
		Order: map[string]string{
			models.OrderByRating: "desc",
		},
	}

	currentIndex := 0
	popularManga, err := api.GetManga(popularParams)

	if err != nil {
		log.Println("Error fetching popular manga:", err)
		return nil
	}

	popularFlex.SetBorder(true).SetTitle("Popular").SetTitleAlign(tview.AlignLeft)

	// Create a content area for popular manga
	p.popularContent.SetDirection(tview.FlexColumn)
	p.buildPopularContent(p.popularContent, popularManga[currentIndex])

	// Create a navigation flex for popular manga
	popularNavigationFlex := tview.NewFlex().SetDirection(tview.FlexColumn)
	popularNavigationFlex.SetBorder(false)
	leftButton := tview.NewButton("◀ Previous")
	viewButton := tview.NewButton("View Detail")
	rightButton := tview.NewButton("Next ▶")
	leftButton.SetStyle(tcell.StyleDefault.Foreground(tcell.ColorYellow).Background(tcell.ColorBlack))
	viewButton.SetStyle(tcell.StyleDefault.Foreground(tcell.ColorGreen).Background(tcell.ColorBlack))
	rightButton.SetStyle(tcell.StyleDefault.Foreground(tcell.ColorYellow).Background(tcell.ColorBlack))
	leftButton.SetSelectedFunc(func() {
		if currentIndex > 0 {
			currentIndex--
			p.buildPopularContent(p.popularContent, popularManga[currentIndex])
		}
	})
	rightButton.SetSelectedFunc(func() {
		if currentIndex < len(popularManga)-1 {
			currentIndex++
			p.buildPopularContent(p.popularContent, popularManga[currentIndex])
		}
	})
	viewButton.SetSelectedFunc(func() {
		if currentIndex < len(popularManga) {
			selectedManga := popularManga[currentIndex]
			p.showMangaDetailModal(&selectedManga)
		}
	})
	popularNavigationFlex.AddItem(leftButton, 0, 1, false)
	popularNavigationFlex.AddItem(viewButton, 0, 1, false)
	popularNavigationFlex.AddItem(rightButton, 0, 1, false)

	// Add popular manga to the content area
	popularFlex.AddItem(p.popularContent, 0, 9, false)
	popularFlex.AddItem(popularNavigationFlex, 0, 1, false)

	return popularFlex
}

func (p *HomePage) buildPopularContent(popularContent *tview.Flex, manga models.Manga) {
	popularContent.Clear()

	// Create a box for cover art placeholder
	imageFlex := tview.NewImage()
	// imageFlex.SetSize(30, 30)

	// Get and set the image
	if img := services.GetMangaImage(manga.ID, 256); img != nil {
		imageFlex.SetImage(img)
	}

	infoFlex := tview.NewFlex().SetDirection(tview.FlexRow)
	infoFlex.SetBorder(true).SetTitle("Information").SetTitleAlign(tview.AlignLeft)
	title := tview.NewTextView().
		SetText(fmt.Sprintf("Title: %s", manga.Attributes.Title["en"])).
		SetTextColor(tcell.ColorOrange).
		SetTextAlign(tview.AlignLeft).
		SetDynamicColors(true)
	status := tview.NewTextView().
		SetText(fmt.Sprintf("Status: %s", services.FormatTextStatus(manga.Attributes.Status))).
		SetTextColor(services.GetColorStatus(manga.Attributes.Status)).
		SetTextAlign(tview.AlignLeft).
		SetDynamicColors(true)
	year := tview.NewTextView().
		SetText(fmt.Sprintf("Year: %s", services.FormatTextYear(manga.Attributes.Year))).
		SetTextColor(tcell.ColorWhite).
		SetTextAlign(tview.AlignLeft).
		SetDynamicColors(true)
	description := tview.NewTextView().
		SetText(fmt.Sprintf("Description: %s", manga.Attributes.Description["en"])).
		SetTextColor(tcell.ColorWhite).
		SetTextAlign(tview.AlignLeft).
		SetDynamicColors(true)
	tags := manga.Attributes.Tags
	tagsView := tview.NewTextView().
		SetText(fmt.Sprintf("Tags: %s", services.FormatTags(tags))).
		SetTextAlign(tview.AlignLeft).
		SetDynamicColors(true)
	infoFlex.AddItem(title, 0, 1, false)
	infoFlex.AddItem(status, 0, 1, false)
	infoFlex.AddItem(year, 0, 1, false)
	infoFlex.AddItem(tagsView, 0, 2, false)
	infoFlex.AddItem(description, 0, 4, false)

	popularContent.AddItem(imageFlex, 0, 3, false)
	popularContent.AddItem(infoFlex, 0, 7, false)
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

func (p *HomePage) setMangaListData(mangaList *tview.Table, params models.MangaQueryParams) {
	mangas, err := api.GetManga(params)

	if err != nil {
		log.Println("Error fetching manga data:", err)
		return
	}

	for i, manga := range mangas {
		mangaCopy := manga
		titleCell := tview.NewTableCell(manga.Attributes.Title["en"]).SetReference(&mangaCopy).SetMaxWidth(30)
		mangaList.SetCell(i+1, 0, titleCell)
		mangaList.SetCell(i+1, 1, p.formatTableStatus(manga.Attributes.Status))
		mangaList.SetCell(i+1, 2, tview.NewTableCell(strconv.Itoa(manga.Attributes.Year)))
	}

	mangaList.SetSelectedFunc(func(row, column int) {
		if row == 0 {
			return // Skip header row
		}
		// Get the cell at the selected row
		cell := mangaList.GetCell(row, 0)
		if cell == nil {
			log.Printf("Error: Invalid cell at row %d", row)
			return
		}

		// Get manga reference
		ref := cell.GetReference()
		if ref == nil {
			log.Printf("Error: No manga reference at row %d", row)
			return
		}

		selectedManga, ok := ref.(*models.Manga)
		if !ok || selectedManga == nil {
			log.Printf("Error: Invalid manga reference at row %d", row)
			return
		}

		p.showMangaDetailModal(selectedManga)
		// p.buildPopularContent(p.popularContent, *selectedManga)
	})

	mangaList.SetSelectable(true, false)
}

func (p *HomePage) showMangaDetailModal(manga *models.Manga) {
	// Create content area
	content := fmt.Sprintf(`Title: [orange]%s[-]
		Status: [%s]%s[-]
		Year: %s
		Description: 
		%s
		Tags: %s`,
		manga.Attributes.Title["en"],
		services.GetColorStatus(manga.Attributes.Status).String(),
		services.FormatTextStatus(manga.Attributes.Status),
		services.FormatTextYear(manga.Attributes.Year),
		services.ShortenDescription(manga.Attributes.Description["en"], 300),
		services.FormatTags(manga.Attributes.Tags))

	// Create and configure modal
	modal := tview.NewModal().
		SetText(content).
		SetBackgroundColor(tcell.ColorBlack).
		AddButtons([]string{"Close", "View Detail"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "Close" {
				p.app.SetRoot(p.rootView, true)
			} else if buttonLabel == "View Detail" {
				// TODO: Implement chapter view
				p.app.SetRoot(p.rootView, true)
			}
		})

	// Show the modal page
	p.app.SetRoot(modal, false)
}

func (p *HomePage) formatTableStatus(status string) *tview.TableCell {
	switch status {
	case "ongoing":
		return tview.NewTableCell("Ongoing").SetTextColor(tcell.ColorGreen)
	case "completed":
		return tview.NewTableCell("Completed").SetTextColor(tcell.ColorOrange)
	case "hiatus":
		return tview.NewTableCell("Hiatus").SetTextColor(tcell.ColorYellow)
	case "cancelled":
		return tview.NewTableCell("Cancelled").SetTextColor(tcell.ColorRed)
	default:
		return tview.NewTableCell(status).SetTextColor(tcell.ColorWhite)
	}
}
