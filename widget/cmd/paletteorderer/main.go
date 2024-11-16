package main

import (
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"github.com/jeromelesaux/fyne-io/widget/palette/orderer"
)

func main() {
	a := app.New()
	w := a.NewWindow("image table")
	w.Resize(fyne.NewSize(200, 400))
	p := NewCpcPlusPalette()[100:116]
	o := orderer.NewOrderer(p, setPalette)
	w.SetContent(o.NewOrderer())
	w.ShowAndRun()
}

func setPalette(p color.Palette) {
	fmt.Println("get new palette")
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
