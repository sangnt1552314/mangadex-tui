package components

import (
	"image"
	_ "image/jpeg"
	_ "image/png"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// ImageView is a component that displays images in the terminal
type ImageView struct {
	*tview.Box
	image image.Image
}

// NewImageView creates and returns a new image view
func NewImageView() *ImageView {
	return &ImageView{
		Box: tview.NewBox(),
	}
}

// SetImage sets the image to be displayed
func (i *ImageView) SetImage(img image.Image) *ImageView {
	i.image = img
	return i
}

// Draw draws this primitive onto the screen
func (i *ImageView) Draw(screen tcell.Screen) {
	if i.image == nil {
		return
	}

	i.Box.Draw(screen)

	x, y, width, height := i.GetInnerRect()
	imgWidth := i.image.Bounds().Dx()
	imgHeight := i.image.Bounds().Dy()

	// Calculate scaling factors
	scaleX := float64(imgWidth) / float64(width)
	scaleY := float64(imgHeight) / float64(height*2) // Each terminal cell is roughly twice as tall as wide

	for cy := 0; cy < height; cy++ {
		for cx := 0; cx < width; cx++ {
			// Sample the image (using the top and bottom half of where this cell represents)
			topY := int(float64(cy*2) * scaleY)
			bottomY := int(float64(cy*2+1) * scaleY)
			imageX := int(float64(cx) * scaleX)

			if topY >= imgHeight || bottomY >= imgHeight || imageX >= imgWidth {
				continue
			}

			topColor := i.image.At(imageX, topY)
			bottomColor := i.image.At(imageX, bottomY)

			// Convert to tcell colors
			r1, g1, b1, _ := topColor.RGBA()
			r2, g2, b2, _ := bottomColor.RGBA()

			topTcellColor := tcell.NewRGBColor(int32(r1>>8), int32(g1>>8), int32(b1>>8))
			bottomTcellColor := tcell.NewRGBColor(int32(r2>>8), int32(g2>>8), int32(b2>>8))

			// Use Unicode block characters for improved resolution
			// '▀' (U+2580) is a half block where the top half is the foreground color
			// and the bottom half is the background color
			screen.SetContent(x+cx, y+cy, '▀', nil, tcell.StyleDefault.
				Background(bottomTcellColor).
				Foreground(topTcellColor))
		}
	}
}
