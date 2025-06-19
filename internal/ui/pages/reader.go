package pages

import (
	"image"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/sangnt1552314/mangadex-tui/internal/models"
	"github.com/sangnt1552314/mangadex-tui/internal/services"
	"github.com/sangnt1552314/mangadex-tui/internal/ui/interfaces"
)

type ReaderPage struct {
	app      interfaces.AppInterface
	rootView *tview.Flex
	manga    *models.Manga
	chapter  *models.Chapter
	images   []image.Image
}

func NewReaderPage(app interfaces.AppInterface) *ReaderPage {
	return &ReaderPage{
		app:      app,
		rootView: tview.NewFlex(),
		manga:    nil,
		chapter:  nil,
		images:   nil,
	}
}

func (p *ReaderPage) Name() string {
	return "reader"
}

func (p *ReaderPage) View() tview.Primitive {
	return p.rootView
}

func (p *ReaderPage) SetData(manga *models.Manga, chapter *models.Chapter) {
	p.manga = manga
	p.chapter = chapter
	p.updateUI()
}

func (p *ReaderPage) Init(app interfaces.AppInterface) {
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

func (p *ReaderPage) updateUI() {
	if p.manga == nil || p.chapter == nil {
		return
	}

	p.rootView.Clear()

	// Layout
	p.rootView.SetDirection(tview.FlexRow).
		SetBorder(false)

	// Layout - Main Content
	mainContent := p.setupMainContent()

	// Layout - Menu
	menu := p.setupMenu()

	// Add components to the root view
	p.rootView.AddItem(mainContent, 0, 1, true)
	p.rootView.AddItem(menu, 3, 0, false)

}

func (p *ReaderPage) setupMenu() tview.Primitive {
	menuFlex := tview.NewFlex().SetDirection(tview.FlexColumn)
	menuFlex.SetBackgroundColor(tcell.ColorBlack).SetBorder(true).SetTitle("Options").SetTitleAlign(tview.AlignLeft)

	homeButton := tview.NewButton("⌂ Home")
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

	menuFlex.AddItem(homeButton, 9, 1, false)
	menuFlex.AddItem(searchButton, 9, 1, false)
	menuFlex.AddItem(aboutButton, 9, 1, false)
	menuFlex.AddItem(exitButton, 9, 1, false)

	return menuFlex
}

func (p *ReaderPage) setupMainContent() tview.Primitive {
	currentImageIndex := 0

	images, err := services.GetImagesByChapterId(p.chapter.ID)
	p.images = images

	if err != nil {
		mainContent := tview.NewTextView().
			SetText("Error loading chapter images: " + err.Error()).
			SetTextColor(tcell.ColorRed).
			SetBackgroundColor(tcell.ColorBlack).
			SetBorder(true)
		return mainContent
	}

	mainContent := tview.NewFlex().SetDirection(tview.FlexRow)
	mainContent.SetBorder(true)

	imageFlex := tview.NewImage()
	p.showImage(currentImageIndex, imageFlex)

	navigationFlex := tview.NewFlex().SetDirection(tview.FlexColumn)
	navigationFlex.SetBorder(false)
	leftButton := tview.NewButton("◀ Previous")
	rightButton := tview.NewButton("Next ▶")
	leftButton.SetStyle(tcell.StyleDefault.Foreground(tcell.ColorYellow).Background(tcell.ColorBlack))
	rightButton.SetStyle(tcell.StyleDefault.Foreground(tcell.ColorYellow).Background(tcell.ColorBlack))
	leftButton.SetSelectedFunc(func() {
		currentImageIndex--
		if currentImageIndex < 0 {
			currentImageIndex = 0
		}
		p.showImage(currentImageIndex, imageFlex)
	})
	rightButton.SetSelectedFunc(func() {
		currentImageIndex++
		if currentImageIndex >= len(images) {
			currentImageIndex = len(images) - 1
		}
		p.showImage(currentImageIndex, imageFlex)
	})

	navigationFlex.AddItem(leftButton, 0, 1, false)
	navigationFlex.AddItem(rightButton, 0, 1, false)

	mainContent.AddItem(imageFlex, 0, 1, false)
	mainContent.AddItem(navigationFlex, 1, 0, false)

	return mainContent
}

func (p *ReaderPage) showImage(i int, imageFlex *tview.Image) {
	image := p.images[i]
	if image == nil {
		imageFlex.SetImage(nil)
		return
	}

	imageFlex.SetImage(image)
	imageFlex.SetBackgroundColor(tcell.ColorBlack)
}
