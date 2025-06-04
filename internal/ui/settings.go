package ui

import (
	"github.com/rivo/tview"
)

// SettingsPage represents the settings interface
type SettingsPage struct {
	*tview.Flex
	app         *App
	settingList *tview.List
}

// NewSettingsPage creates a new settings page
func NewSettingsPage(app *App) *SettingsPage {
	settings := &SettingsPage{
		Flex:        tview.NewFlex().SetDirection(tview.FlexRow),
		app:         app,
		settingList: tview.NewList().ShowSecondaryText(true),
	}

	// Configure settings list
	settings.settingList.SetBorder(true).
		SetTitle("Settings").
		SetTitleAlign(tview.AlignLeft)

	// Add dummy settings
	settings.settingList.
		AddItem("Theme", "Change application color theme (Current: Default)", 'T', nil).
		AddItem("Download Path", "Set manga download location (Current: ~/Downloads)", 'D', nil).
		AddItem("Cache Size", "Set maximum cache size (Current: 1GB)", 'C', nil).
		AddItem("Language", "Set preferred manga language (Current: English)", 'L', nil).
		AddItem("Reading Mode", "Set reading direction (Current: Left to Right)", 'R', nil)

	// Layout
	settings.AddItem(settings.settingList, 0, 1, true)

	return settings
}

// Focus implements the Primitive interface for SettingsPage
func (s *SettingsPage) Focus(delegate func(p tview.Primitive)) {
	delegate(s.settingList)
}

// HasFocus implements the Primitive interface for SettingsPage
func (s *SettingsPage) HasFocus() bool {
	return s.settingList.HasFocus()
}
