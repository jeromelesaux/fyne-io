package main

import (
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	ui "github.com/jeromelesaux/fyne-io/widget"
	"github.com/jeromelesaux/martine/constants"
)

var (
	colorIndex int
	win        *fyne.Window
)

func main() {
	a := app.New()
	w := a.NewWindow("Colored Button")
	w.Resize(fyne.NewSize(900, 400))

	w.SetContent(
		container.NewVBox(
			widget.NewButtonWithIcon("Swap", theme.ColorChromaticIcon(), func() {
				swapColor(constants.CpcOldPalette, *win)
			}),
		),
	)
	win = &w
	w.ShowAndRun()
}

func clr(c color.Color) {
	fmt.Println(c)
}

func index(i int) {
	fmt.Println(i)
	colorIndex = i
}

func swapColor(p color.Palette, w fyne.Window) {
	pt := ui.NewPaletteTable(p, clr, index, setPalette)
	var cont *fyne.Container
	var colorToChange color.Color

	cont = container.NewVBox(
		pt,
		widget.NewButton("Color", func() {
			picker := dialog.NewColorPicker("Pick a Color", "What is your favorite color?", func(c color.Color) {
				colorToChange = c
			}, w)
			picker.Advanced = true
			picker.Show()
		}),
		widget.NewButton("swap", func() {
			p[colorIndex] = colorToChange
			npt := ui.NewPaletteTable(p, clr, index, setPalette)
			pt = npt
			cont.Refresh()
		}))
	cont.Resize(fyne.NewSize(200, 200))
	d := dialog.NewCustom("Swap color", "Ok", cont, w)
	d.Resize(w.Canvas().Size())
	d.Show()
	fmt.Println(pt.Palette)
}

func setPalette(p color.Palette) {
	fmt.Println(p)
}
