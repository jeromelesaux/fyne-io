package custom_widget

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type ProgressInfiniteDialog struct {
	bar     *widget.ProgressBarInfinite
	parent  fyne.Window
	win     *widget.PopUp
	message string
}

func NewProgressInfinite(message string, parent fyne.Window) *ProgressInfiniteDialog {
	bar := widget.NewProgressBarInfinite()
	rect := canvas.NewRectangle(color.Transparent)
	rect.SetMinSize(fyne.NewSize(200, 0))

	p := &ProgressInfiniteDialog{bar: bar, parent: parent, message: message}
	p.setButtons(container.NewMax(rect, bar))
	return p
}

func (d *ProgressInfiniteDialog) setButtons(buttons fyne.CanvasObject) {
	content := container.New(layout.NewVBoxLayout(),
		widget.NewLabel(d.message),
		layout.NewSpacer(),
		buttons,
	)
	d.win = widget.NewModalPopUp(content, d.parent.Canvas())
	d.Refresh()
}

// Hide this dialog and stop the infinite progress goroutine
func (d *ProgressInfiniteDialog) Hide() {
	d.bar.Hide()
	d.win.Hide()
}

func (d *ProgressInfiniteDialog) Refresh() {
	d.win.Refresh()
}

func (d *ProgressInfiniteDialog) MinSize() fyne.Size {
	return d.win.MinSize()
}

func (d *ProgressInfiniteDialog) Move(position fyne.Position) {
	d.win.Move(position)
}

func (d *ProgressInfiniteDialog) Position() fyne.Position {
	return d.win.Position()
}

func (d *ProgressInfiniteDialog) Resize(size fyne.Size) {
	d.win.Resize(size)
}

func (d *ProgressInfiniteDialog) Show() {
	d.win.Show()
}

func (d *ProgressInfiniteDialog) Size() fyne.Size {
	return d.win.Size()
}

func (d *ProgressInfiniteDialog) Visible() bool {
	return d.win.Visible()
}
