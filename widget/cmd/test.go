package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func main() {
	a := app.New()
	w := a.NewWindow("Devices")
	w.Resize(fyne.NewSize(320, 240))

	t := widget.NewTable(func() (int, int) {
		return 100, 4
	}, func() fyne.CanvasObject {
		return container.NewMax(widget.NewLabel("template11"), widget.NewIcon(nil))
	}, func(id widget.TableCellID, o fyne.CanvasObject) {
		l := o.(*fyne.Container).Objects[0].(*widget.Label)
		i := o.(*fyne.Container).Objects[1].(*widget.Icon)
		l.Show()
		i.Hide()

		switch id.Col {
		case 2:
			l.Hide()
			i.Show()

			i.SetResource(getIcon(id.Row))
		case 0:
			l.SetText("hostname")
		case 1:
			l.SetText("127.0.0.1")
		case 3:
			l.SetText("notes...")
		}
	})
	t.SetColumnWidth(2, 24)
	t.SetColumnWidth(3, 156)
	c := container.NewGridWithRows(
		2,
		t,
		widget.NewButtonWithIcon("test", theme.ConfirmIcon(), func() {}),
	)
	w.SetContent(c)

	w.ShowAndRun()
}

// nolint: ireturn
func getIcon(i int) fyne.Resource {
	switch i % 3 {
	case 1:
		return theme.HomeIcon()
	case 2:
		return theme.MailSendIcon()
	default:
		return theme.MediaVideoIcon()
	}
}
