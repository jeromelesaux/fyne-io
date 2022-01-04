package main

import (
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	ui "github.com/jeromelesaux/fyne-io/custom_widget"
)

func main() {
	a := app.New()
	w := a.NewWindow("Colored Button")
	w.Resize(fyne.NewSize(50, 50))
	w.SetContent(
		container.NewVBox(
			ui.NewButtonColored(color.Black, tappedBlack),
			ui.NewButtonColored(color.White, tappedWhite),
		),
	)
	w.ShowAndRun()
}

func tappedBlack() {
	fmt.Println("black")
}

func tappedWhite() {
	fmt.Println("white")
}
