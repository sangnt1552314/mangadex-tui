package pages

import (
	"fmt"
	"log"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"github.com/sangnt1552314/mangadex-tui/internal/api"
	"github.com/sangnt1552314/mangadex-tui/internal/models"
	"github.com/sangnt1552314/mangadex-tui/internal/services"
	"github.com/sangnt1552314/mangadex-tui/internal/ui/interfaces"
)

type DetailPage struct {
	app      interfaces.AppInterface
	rootView *tview.Flex
	manga    *models.Manga
	limit    int
	offset   int
	total    int
}

func NewDetailPage(app interfaces.AppInterface) *DetailPage {
	return &DetailPage{
		app:      app,
		rootView: tview.NewFlex(),
		manga:    nil,
		limit:    100,
		offset:   0,
		total:    0,
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
		case tcell.KeyCtrlH:
			app.SwitchToPage("home")
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

	p.rootView.AddItem(mainContent, 0, 1, false)
	p.rootView.AddItem(menu, 3, 0, false)
}

func (p *DetailPage) setupMenu() tview.Primitive {
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
	menuFlex.AddItem(searchButton, 9, 1, false)
	menuFlex.AddItem(aboutButton, 9, 1, false)
	menuFlex.AddItem(exitButton, 9, 1, false)

	return menuFlex
}

func (p *DetailPage) setupMainContent() tview.Primitive {
	if p.manga == nil {
		noMangaText := tview.NewTextView().SetText("No manga selected")
		return noMangaText
	}

	mainContent := tview.NewFlex().SetDirection(tview.FlexColumn)
	mainContent.SetBorder(false)

	imageFlex := tview.NewImage()
	if img := services.GetMangaImageByFilename(p.manga.ID, services.GetCoverFileName(*p.manga), 512); img != nil {
		imageFlex.SetImage(img)
		// imageFlex.SetAlign(tview.AlignCenter, tview.AlignCenter)
	}

	mangaDataFlex := tview.NewFlex().SetDirection(tview.FlexRow)
	mangaDataFlex.SetBorder(false)

	topMangaDataFlex := tview.NewFlex()
	p.setupTopMangaDataFlex(topMangaDataFlex)

	bottomMangaDataFlex := tview.NewFlex().SetDirection(tview.FlexColumn)
	bottomMangaDataFlex.SetBorder(false)

	categoryDataFlex := tview.NewFlex()
	p.setupCategoryDataFlex(categoryDataFlex)

	chapterDataFlex := tview.NewFlex()
	p.setupChapterDataFlex(chapterDataFlex)

	bottomMangaDataFlex.AddItem(categoryDataFlex, 0, 3, false)
	bottomMangaDataFlex.AddItem(chapterDataFlex, 0, 7, false)

	mangaDataFlex.AddItem(topMangaDataFlex, 0, 4, false)
	mangaDataFlex.AddItem(bottomMangaDataFlex, 0, 6, false)

	mainContent.AddItem(imageFlex, 0, 3, false)
	mainContent.AddItem(mangaDataFlex, 0, 7, false)

	return mainContent
}

func (p *DetailPage) setupTopMangaDataFlex(flex *tview.Flex) {
	flex.SetDirection(tview.FlexColumn)
	flex.SetBorder(false)

	leftFlex := tview.NewFlex().SetDirection(tview.FlexRow)
	leftFlex.SetBorder(true).SetTitle("Manga Info").SetTitleAlign(tview.AlignLeft)
	rightFlex := tview.NewFlex().SetDirection(tview.FlexRow)
	rightFlex.SetBorder(true).SetTitle("Description").SetTitleAlign(tview.AlignLeft)

	// Main Title and Alternative Title
	mainTitle := tview.NewTextView().
		SetText(p.manga.Attributes.Title["en"]).
		SetTextColor(tcell.ColorOrange)

	var altTitle string
	if len(p.manga.Attributes.AltTitles) > 0 {
		for _, title := range p.manga.Attributes.AltTitles {
			if engTitle, ok := title["en"]; ok {
				altTitle = engTitle
				break
			}
		}

		if altTitle == "" && len(p.manga.Attributes.AltTitles) > 0 {
			for _, val := range p.manga.Attributes.AltTitles[0] {
				altTitle = val
				break
			}
		}
	}

	altTitleView := tview.NewTextView().
		SetText(altTitle).
		SetTextColor(tcell.ColorLightGrey)

	// Year and Status
	yearText := tview.NewTextView().
		SetText(fmt.Sprintf("Year: %s", services.FormatTextYear(p.manga.Attributes.Year)))

	status := p.manga.Attributes.Status
	statusColor := services.GetColorStatus(status)
	statusText := tview.NewTextView().
		SetText(fmt.Sprintf("Status: [%s]%s[-]",
			statusColor.TrueColor(),
			services.FormatTextStatus(status))).
		SetDynamicColors(true)

	// Author
	authorName := services.GetAuthorName(*p.manga)
	authorText := tview.NewTextView().
		SetText(fmt.Sprintf("Author: %s", authorName)).
		SetTextColor(tcell.ColorLightCyan)

	// Artist
	artistName := services.GetArtistName(*p.manga)
	artistText := tview.NewTextView().
		SetText(fmt.Sprintf("Artist: %s", artistName)).
		SetTextColor(tcell.ColorLightCyan)

	// Description
	var description string
	if desc, ok := p.manga.Attributes.Description["en"]; ok {
		description = desc
	} else if len(p.manga.Attributes.Description) > 0 {
		// Get first description in any language
		for _, desc := range p.manga.Attributes.Description {
			description = desc
			break
		}
	}

	descText := tview.NewTextView().
		SetText(description).
		SetWrap(true).
		SetDynamicColors(true)

	leftFlex.AddItem(mainTitle, 0, 1, false)
	leftFlex.AddItem(altTitleView, 0, 1, false)
	leftFlex.AddItem(yearText, 0, 1, false)
	leftFlex.AddItem(statusText, 0, 1, false)
	leftFlex.AddItem(authorText, 0, 1, false)
	leftFlex.AddItem(artistText, 0, 1, false)

	rightFlex.AddItem(descText, 0, 1, false)
	flex.AddItem(leftFlex, 0, 1, false)
	flex.AddItem(rightFlex, 0, 2, false)
}

func (p *DetailPage) setupCategoryDataFlex(flex *tview.Flex) {
	flex.SetDirection(tview.FlexRow)
	flex.SetBorder(true).SetTitle("Categories").SetTitleAlign(tview.AlignLeft)

	// Tags
	tagsText := tview.NewTextView().
		SetText(fmt.Sprintf("Tags: %s", services.FormatTags(p.manga.Attributes.Tags))).
		SetDynamicColors(true)

	flex.AddItem(tagsText, 0, 1, false)
}

func (p *DetailPage) setupChapterDataFlex(flex *tview.Flex) {
	flex.SetDirection(tview.FlexRow)
	flex.SetBorder(true).SetTitle("Chapters").SetTitleAlign(tview.AlignLeft)

	simpleChapterResp, err := api.GetChapterListResponse(models.ChapterQueryParams{
		MangaId:            p.manga.ID,
		Limit:              1,
		Offset:             0,
		TranslatedLanguage: []string{"en"},
		Order: map[string]string{
			"volume":  "asc",
			"chapter": "asc",
		},
	})

	if err != nil {
		log.Println("Error fetching manga data:", err)
		return
	}

	p.total = simpleChapterResp.Total

	params := models.ChapterQueryParams{
		MangaId:            p.manga.ID,
		Limit:              p.limit,
		Offset:             p.offset,
		TranslatedLanguage: []string{"en"},
		Order: map[string]string{
			"volume":  "asc",
			"chapter": "asc",
		},
	}

	chapterList := tview.NewTable().SetFixed(1, 0)
	p.setChapterListData(chapterList, params)

	navigationFlex := tview.NewFlex().SetDirection(tview.FlexColumn)
	navigationFlex.SetBorder(false)
	leftButton := tview.NewButton("◀ Previous")
	rightButton := tview.NewButton("Next ▶")
	leftButton.SetStyle(tcell.StyleDefault.Foreground(tcell.ColorYellow).Background(tcell.ColorBlack))
	rightButton.SetStyle(tcell.StyleDefault.Foreground(tcell.ColorYellow).Background(tcell.ColorBlack))
	leftButton.SetSelectedFunc(func() {
		if p.offset > 0 {
			p.offset -= p.limit
			params.Offset = p.offset
			p.setChapterListData(chapterList, params)
		}
		if p.offset < 0 {
			p.offset = 0
			params.Offset = p.offset
			p.setChapterListData(chapterList, params)
		}
		if p.offset >= p.total {
			p.offset = p.total - p.limit
			params.Offset = p.offset
			p.setChapterListData(chapterList, params)
		}
	})
	rightButton.SetSelectedFunc(func() {
		if p.offset+p.limit < p.total {
			p.offset += p.limit
			params.Offset = p.offset
			p.setChapterListData(chapterList, params)
		}
		if p.offset >= p.total {
			p.offset = p.total - p.limit
			params.Offset = p.offset
			p.setChapterListData(chapterList, params)
		}
		if p.offset < 0 {
			p.offset = 0
			params.Offset = p.offset
			p.setChapterListData(chapterList, params)
		}
	})
	navigationFlex.AddItem(leftButton, 0, 1, false)
	navigationFlex.AddItem(rightButton, 0, 1, false)

	flex.AddItem(chapterList, 0, 1, false)
	flex.AddItem(navigationFlex, 1, 0, false)

}

func (p *DetailPage) setChapterListData(list *tview.Table, params models.ChapterQueryParams) {
	list.Clear()

	p.setChapterListHeader(list)

	chapters, err := api.GetChapters(params)

	if err != nil {
		log.Println("Error fetching manga data:", err)
		return
	}

	for i, chapter := range chapters {
		chapterCopy := chapter
		titleCell := tview.NewTableCell(chapter.Attributes.Title).SetReference(&chapterCopy).SetMaxWidth(30)
		chapterCell := tview.NewTableCell(fmt.Sprintf("Chapter %s", chapter.Attributes.Chapter)).SetReference(&chapterCopy).SetMaxWidth(20)
		list.SetCell(i+1, 0, chapterCell)
		list.SetCell(i+1, 1, titleCell)
	}

	list.SetSelectable(true, false)

	list.SetSelectedFunc(func(row, column int) {
		if row == 0 {
			return // Skip header row
		}

		cell := list.GetCell(row, 0)
		if cell == nil {
			log.Printf("Error: Invalid cell at row %d", row)
			return
		}

		ref := cell.GetReference()
		if ref == nil {
			log.Printf("Error: No chapter reference at row %d", row)
			return
		}

		selectedChapter := ref.(*models.Chapter)
		if selectedChapter != nil {
			readerPage := p.app.GetPageObject("reader").(*ReaderPage)
			readerPage.SetData(p.manga, selectedChapter)
			p.app.RestorePages()
			p.app.SwitchToPage("reader")
		}
	})
}

func (p *DetailPage) setChapterListHeader(list *tview.Table) {
	list.SetCell(0, 0, tview.NewTableCell("Chapter").
		SetSelectable(false).
		SetTextColor(tcell.ColorOrange))

	list.SetCell(0, 1, tview.NewTableCell("Title").
		SetSelectable(false).
		SetTextColor(tcell.ColorYellow))
}
