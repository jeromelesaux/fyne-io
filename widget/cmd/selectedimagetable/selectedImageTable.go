package main

import (
	"fmt"
	"image"
	"image/png"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	ui "github.com/jeromelesaux/fyne-io/widget"
)

var size float32 = 70.

func main() {
	a := app.New()
	w := a.NewWindow("image table")
	w.Resize(fyne.NewSize(400, 200))
	imgs := []image.Image{OpenImage(), OpenImage(), OpenImage()}
	t := ui.NewImageSelectionTableWithImages(imgs, fyne.NewSize(size, size))
	w.SetContent(
		container.New(
			layout.NewAdaptiveGridLayout(2),
			container.New(layout.NewGridLayoutWithRows(3),
				container.New(
					layout.NewAdaptiveGridLayout(1),
					widget.NewButton("append", func() {
						t.Append(canvas.NewImageFromImage(OpenFyne()))
					}),
				),
				container.New(
					layout.NewAdaptiveGridLayout(1),
					widget.NewButton("count", func() {
						fmt.Println("nb selected images", len(t.Images()))
					}),
				),
			),

			container.NewScroll(
				t,
			),
		),
	)

	w.ShowAndRun()
}

func OpenFyne() image.Image {
	f, err := os.Open("fyne.png")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	img, err := png.Decode(f)
	if err != nil {
		panic(err)
	}
	return img
}

func OpenImage() image.Image {
	f, err := os.Open("image.png")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	img, err := png.Decode(f)
	if err != nil {
		panic(err)
	}
	return img
}
