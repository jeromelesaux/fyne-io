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
	e := editor.NewEditor(i, editor.MagnifyX2, color.Palette{}, color.Palette{})
	w.SetContent(e.NewEditor())
	w.ShowAndRun()
}
