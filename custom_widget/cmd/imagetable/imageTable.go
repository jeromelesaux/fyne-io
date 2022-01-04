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
	ui "github.com/jeromelesaux/fyne-io/custom_widget"
)

var win *fyne.Window

func main() {
	a := app.New()
	w := a.NewWindow("image table")
	w.Resize(fyne.NewSize(900, 400))
	table := ui.NewImageTable(NewImagesColletion(20, 30), fyne.NewSize(40, 40), 20, 30, nil, indexSelected, nil)
	table.Refresh()
	w.SetContent(
		container.New(
			layout.NewAdaptiveGridLayout(2),
			container.New(
				layout.NewAdaptiveGridLayout(1),
				canvas.NewImageFromImage(OpenFyne()),
			),
			container.NewScroll(
				table,
			),
			container.New(
				layout.NewAdaptiveGridLayout(1),
				widget.NewButton("update", func() {
					img := OpenFyne()
					table.SubstitueImage(0, 0, *canvas.NewImageFromImage(img))
					w.Canvas().Refresh(w.Content())
				}),
			),
		),
	)
	win = &w
	w.ShowAndRun()
}

func indexSelected(row, col int) {
	fmt.Printf("row;%d,col:%d\n", row, col)
}

func NewImagesColletion(nbRows, nbCols int) *[][]canvas.Image {
	images := make([][]canvas.Image, nbRows)
	for i := 0; i < nbRows; i++ {
		images[i] = make([]canvas.Image, nbCols)
	}
	img := canvas.NewImageFromImage(OpenImage())
	for i := 0; i < nbRows; i++ {
		for j := 0; j < nbCols; j++ {
			images[i][j] = *img
		}
	}
	return &images
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
