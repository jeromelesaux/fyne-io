package custom_widget

import (
	"image"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type ImageSelectionTable struct {
	widget.Table
	images   []*canvas.Image
	selected []bool
	size     fyne.Size
}

func NewImageSelectionTable(size fyne.Size) *ImageSelectionTable {
	table := &ImageSelectionTable{
		images:   make([]*canvas.Image, 0),
		selected: make([]bool, 0),
		size:     size,
	}

	table.UpdateCell = table.updateCell
	table.Length = table.length
	table.CreateCell = table.createCell
	table.OnSelected = table.onSelect
	table.ExtendBaseWidget(table)
	return table

}

func NewImageSelectionTableWithImages(imgs []image.Image, size fyne.Size) *ImageSelectionTable {
	table := &ImageSelectionTable{
		images:   make([]*canvas.Image, len(imgs)),
		selected: make([]bool, len(imgs)),
		size:     size,
	}

	for i := 0; i < len(imgs); i++ {
		table.selected[i] = true
		table.images[i] = canvas.NewImageFromImage(imgs[i])
	}

	table.UpdateCell = table.updateCell
	table.Length = table.length
	table.CreateCell = table.createCell
	table.OnSelected = table.onSelect
	table.ExtendBaseWidget(table)

	return table
}

func (t *ImageSelectionTable) Images() []image.Image {
	im := make([]image.Image, 0)
	for i := 0; i < len(t.images); i++ {
		if t.selected[i] {
			im = append(im, t.images[i].Image)
		}
	}
	return im
}
func (t *ImageSelectionTable) Substitue(indice int, img *canvas.Image) {
	if indice >= len(t.images) {
		return
	}
	t.images[indice] = img
	t.Refresh()
	canvas.Refresh(t)
}

func (t *ImageSelectionTable) Remove(pos int) {
	t.images = append(t.images[:pos], t.images[pos+1:]...)
	t.Refresh()
	canvas.Refresh(t)
}

func (t *ImageSelectionTable) Reset() {
	t.images = make([]*canvas.Image, 0)
	t.selected = make([]bool, 0)
	canvas.Refresh(t)
}

func (t *ImageSelectionTable) onSelect(id widget.TableCellID) {
	t.selected[id.Col] = !t.selected[id.Col]
	t.Refresh()
	canvas.Refresh(t)
}

func (t *ImageSelectionTable) Append(img *canvas.Image) {
	t.images = append(t.images, img)
	t.selected = append(t.selected, true)
	t.Refresh()
	canvas.Refresh(t)
}

func (t *ImageSelectionTable) createCell() fyne.CanvasObject {
	c := container.NewMax(emptyCell(t.size), widget.NewIcon(nil))
	//	c.Resize(fyne.NewSize(t.size.Width, t.size.Height*2))
	return c
}

func (t *ImageSelectionTable) updateCell(id widget.TableCellID, o fyne.CanvasObject) {
	i := o.(*fyne.Container).Objects[0].(*canvas.Image)
	c := o.(*fyne.Container).Objects[1].(*widget.Icon)

	switch id.Row {
	case 0:
		c.Hide()
		i.Image = t.images[id.Col].Image
		i.SetMinSize(t.size)
		i.Refresh()
		i.Show()
	case 1:
		if t.selected[id.Col] {
			c.SetResource(theme.ConfirmIcon())
		} else {
			c.SetResource(theme.CancelIcon())
		}
		c.Show()
	}

}

func (t *ImageSelectionTable) length() (int, int) {
	return 2, len(t.images)
}
