package main

import (
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"

	custom "github.com/jeromelesaux/fyne-io/widget"
)

func main() {
	a := app.New()
	w := a.NewWindow("infinite progressbar")
	w.Resize(fyne.NewSize(900, 400))
	w.SetContent(
		container.New(
			layout.NewAdaptiveGridLayout(1),
			widget.NewButton("fire me", func() {
				p := custom.NewProgressInfinite("computing please wait...", w)
				p.Show()
				go func() {
					t := time.NewTicker(4 * time.Second)
					<-t.C
					p.Hide()
				}()
			}),
		),
	)
	w.ShowAndRun()
}
