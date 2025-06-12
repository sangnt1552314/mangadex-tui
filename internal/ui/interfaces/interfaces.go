package interfaces

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// AppInterface defines what pages need from the app
type AppInterface interface {
	SwitchToPage(name string)
	AddPage(name string, page tview.Primitive, visible bool)
	Stop()
	EnableMouse(enable bool)
	SetInputCapture(fn func(event *tcell.EventKey) *tcell.EventKey)
	SetRoot(root tview.Primitive, fullscreen bool)
}

// Page defines what the app needs from pages
type Page interface {
	Name() string
	View() tview.Primitive
	Init(AppInterface)
}
