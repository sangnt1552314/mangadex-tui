package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// SearchPage represents the manga search interface
type SearchPage struct {
	*tview.Flex
	app           *App
	searchInput   *tview.InputField
	resultsList   *tview.List
	searchResults *tview.Flex
}

// NewSearchPage creates a new search page
func NewSearchPage(app *App) *SearchPage {
	search := &SearchPage{
		Flex:          tview.NewFlex().SetDirection(tview.FlexRow),
		app:           app,
		searchInput:   tview.NewInputField(),
		resultsList:   tview.NewList().ShowSecondaryText(true),
		searchResults: tview.NewFlex().SetDirection(tview.FlexRow),
	}

	// Configure search input
	search.searchInput.
		SetLabel("Search manga: ").
		SetFieldWidth(0).
		SetDoneFunc(func(key tcell.Key) {
			if key == tcell.KeyEnter {
				query := search.searchInput.GetText()
				if query != "" {
					search.performSearch(query)
				}
			}
		}).
		SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			switch event.Key() {
			case tcell.KeyDown, tcell.KeyTab:
				// Move to results list when pressing down or tab in search input
				if search.resultsList.GetItemCount() > 0 {
					search.app.SetFocus(search.resultsList)
					return nil
				}
			}
			return event
		})

	// Configure results list
	search.resultsList.SetBorder(true).
		SetTitle("Results").
		SetTitleAlign(tview.AlignLeft).
		SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			switch event.Key() {
			case tcell.KeyUp:
				// If at the top of the list, move back to search input
				if search.resultsList.GetCurrentItem() == 0 {
					search.app.SetFocus(search.searchInput)
					return nil
				}
			case tcell.KeyTab:
				// Move back to search input when pressing tab in results
				search.app.SetFocus(search.searchInput)
				return nil
			}
			return event
		})

	// Add placeholder text for results
	search.resultsList.AddItem("No results yet", "Enter a search term and press Enter", 0, nil)

	// Layout
	search.AddItem(search.searchInput, 3, 0, true).
		AddItem(search.resultsList, 0, 1, false)

	return search
}

// performSearch executes the manga search
func (s *SearchPage) performSearch(query string) {
	// Clear previous results
	s.resultsList.Clear()

	// TODO: Implement actual API call to MangaDex
	// For now, add some dummy results
	s.resultsList.AddItem("Loading...", "Searching for: "+query, 0, nil)
	// Add some dummy results for testing navigation
	s.resultsList.AddItem("Test Manga 1", "Description 1", 0, nil)
	s.resultsList.AddItem("Test Manga 2", "Description 2", 0, nil)
	s.resultsList.AddItem("Test Manga 3", "Description 3", 0, nil)
}

// Focus implements the Primitive interface for SearchPage
func (s *SearchPage) Focus(delegate func(p tview.Primitive)) {
	delegate(s.searchInput)
}

// HasFocus implements the Primitive interface for SearchPage
func (s *SearchPage) HasFocus() bool {
	return s.searchInput.HasFocus() || s.resultsList.HasFocus()
}
