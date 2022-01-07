package main

import (
	"fmt"
	"image"
	"image/draw"
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
var size float32 = 40.

func main() {
	a := app.New()
	w := a.NewWindow("image table")
	w.Resize(fyne.NewSize(900, 400))
	table := ui.NewImageTable(NewImagesColletion(1, 4, OpenImage()), fyne.NewSize(size, size), 1, 4, nil, indexSelected, nil)
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
					table.AppendImage(*canvas.NewImageFromImage(img), 0)
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

func NewNumImagesCollection(nbRows, nbCols int) *[][]canvas.Image {
	images := make([][]canvas.Image, nbRows)
	for i := 0; i < nbRows; i++ {
		images[i] = make([]canvas.Image, nbCols)
	}
	for i := 0; i < nbRows; i++ {
		for j := 0; j < nbCols; j++ {
			img := canvas.NewImageFromImage(ImageLabel(fmt.Sprintf(("%d-%d"), i, j)))
			images[i][j] = *img
		}
	}
	return &images
}

func NewImagesColletion(nbRows, nbCols int, in image.Image) *[][]canvas.Image {
	images := make([][]canvas.Image, nbRows)
	for i := 0; i < nbRows; i++ {
		images[i] = make([]canvas.Image, nbCols)
	}
	img := canvas.NewImageFromImage(in)
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

func ImageLabel(text string) image.Image {
	img := image.NewRGBA(image.Rectangle{image.Point{0, 0}, image.Point{int(size), int(size)}})
	draw.Draw(img, img.Bounds(), image.White, image.ZP, draw.Src)
	/*fontBytes, err := ioutil.ReadFile("UbuntuMono-R.ttf")
	if err != nil {
		fmt.Println(err)
	}
	f, err := freetype.ParseFont(fontBytes)
	if err != nil {
		fmt.Println(err)
	}

	ctx := freetype.NewContext()
	ctx.SetDPI(72)
	ctx.SetDst(img)
	ctx.SetFont(f)
	ctx.SetFontSize(8.)
	ctx.SetSrc(image.Black)

	pt := freetype.Pt(5, 5+int(ctx.PointToFixed(8.0)>>6))
	_, err = ctx.DrawString(text, pt)
	if err != nil {
		fmt.Println(err)
	}*/
	return img
}
