package main

import (
	"image"
	"image/color"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"github.com/jeromelesaux/fyne-io/widget/editor"
)

func main() {
	a := app.New()
	w := a.NewWindow("image table")
	w.Resize(fyne.NewSize(900, 400))
	f, err := os.Open("./image.png")
	if err != nil {
		panic(err)
	}
	i, _, err := image.Decode(f)
	if err != nil {
		panic(err)
	}
	e := editor.NewEditor(i, editor.MagnifyX2, CpcOldPalette[0:8], CpcOldPalette)
	w.SetContent(e.NewEditor())
	w.ShowAndRun()
}

var (
	White         = color.RGBA{A: 0xFF, R: 111, G: 125, B: 107}
	SeaGreen      = color.RGBA{A: 0xFF, R: 0, G: 243, B: 107}
	PastelYellow  = color.RGBA{A: 0xFF, R: 244, G: 244, B: 109}
	Blue          = color.RGBA{A: 0xFF, R: 0, G: 2, B: 107}
	Purple        = color.RGBA{A: 0xFF, R: 241, G: 2, B: 104}
	Cyan          = color.RGBA{A: 0xFF, R: 0, G: 120, B: 104}
	Pink          = color.RGBA{A: 0xFF, R: 243, G: 125, B: 107}
	BrightYellow  = color.RGBA{A: 0xFF, R: 244, G: 244, B: 13}
	BrightWhite   = color.RGBA{A: 0xFF, R: 255, G: 244, B: 250}
	BrightRed     = color.RGBA{A: 0xFF, R: 244, G: 5, B: 6}
	BrightMagenta = color.RGBA{A: 0xFF, R: 243, G: 2, B: 245}
	Orange        = color.RGBA{A: 0xFF, R: 243, G: 125, B: 13}
	PastelMagenta = color.RGBA{A: 0xFF, R: 251, G: 125, B: 250}
	BrightGreen   = color.RGBA{A: 0xFF, R: 2, G: 241, B: 1}
	BrightCyan    = color.RGBA{A: 0xFF, R: 15, G: 243, B: 242}
	Black         = color.RGBA{A: 0xFF, R: 0, G: 2, B: 1}
	BrightBlue    = color.RGBA{A: 0xFF, R: 12, G: 2, B: 245}
	Green         = color.RGBA{A: 0xFF, R: 2, G: 120, B: 1}
	SkyBlue       = color.RGBA{A: 0xFF, R: 12, G: 123, B: 245}
	Magenta       = color.RGBA{A: 0xFF, R: 106, G: 2, B: 104}
	PastelGreen   = color.RGBA{A: 0xFF, R: 113, G: 243, B: 107}
	Lime          = color.RGBA{A: 0xFF, R: 113, G: 243, B: 4}
	PastelCyan    = color.RGBA{A: 0xFF, R: 113, G: 243, B: 245}
	Red           = color.RGBA{A: 0xFF, R: 108, G: 2, B: 1}
	Mauve         = color.RGBA{A: 0xFF, R: 108, G: 2, B: 242}
	Yellow        = color.RGBA{A: 0xFF, R: 111, G: 123, B: 1}
	PastelBlue    = color.RGBA{A: 0xFF, R: 111, G: 125, B: 247}

	CpcOldPalette = color.Palette{White,
		SeaGreen,
		PastelYellow,
		Blue,
		Purple,
		Cyan,
		Pink,
		BrightYellow,
		BrightWhite,
		BrightRed,
		BrightMagenta,
		Orange,
		PastelMagenta,
		BrightGreen,
		BrightCyan,
		Black,
		BrightBlue,
		Green,
		SkyBlue,
		Magenta,
		PastelGreen,
		Lime,
		PastelCyan,
		Red,
		Mauve,
		Yellow,
		PastelBlue,
	}
)
