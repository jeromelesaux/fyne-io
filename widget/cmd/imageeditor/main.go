package main

import (
	"image"
	"image/color"
	"image/png"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/jeromelesaux/fyne-io/widget/editor"
)

func main() {
	a := app.New()
	w := a.NewWindow("image table")
	w.Resize(fyne.NewSize(900, 400))
	f, err := os.Open("./red.png")
	if err != nil {
		panic(err)
	}
	i, _, err := image.Decode(f)
	if err != nil {
		panic(err)
	}

	e := editor.NewEditor(i, editor.MagnifyX2, CpcOldPalette[0:8], NewCpcPlusPalette(), save, w)
	w.SetContent(
		container.New(
			layout.NewHBoxLayout(),
			widget.NewButton("New", func() {
				f, err := os.Open("image.png")
				if err != nil {
					panic(err)
				}
				i, _, err := image.Decode(f)
				if err != nil {
					panic(err)
				}
				e.NewImageAndPalette(i, CpcOldPalette[2:6])
				e.NewAvailablePalette(CpcOldPalette)
			}),
			e.NewEmbededEditor("Export"),
		),
	)

	w.ShowAndRun()
}

func save(i image.Image, p color.Palette) {
	f, err := os.Create("new.png")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	if err := png.Encode(f, i); err != nil {
		panic(err)
	}
}

func NewCpcPlusPalette() color.Palette {
	plusPalette := color.Palette{}
	var r, g, b uint8
	for g = 0; g < 0x10; g++ {
		for r = 0; r < 0x10; r++ {
			for b = 0; b < 0x10; b++ {
				plusPalette = append(plusPalette, color.RGBA{R: r * 16, B: b * 16, G: g * 16, A: 0xFF})
			}
		}
	}
	return plusPalette
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
