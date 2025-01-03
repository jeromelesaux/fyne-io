package custom_widget

import (
	"image"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type CheckedImage struct {
	c        *widget.Check
	Selected bool
	i        *canvas.Image
	s        fyne.Size
}

func (c *CheckedImage) tapped(b bool) {
	c.Selected = b
}

func NewCheckedImageWithImage(i *canvas.Image, s fyne.Size) *fyne.Container {
	i.SetMinSize(s)
	c := &CheckedImage{
		Selected: true,
		i:        i,
		s:        s,
	}

	cb := widget.NewCheck("", c.tapped)
	cb.SetChecked(true)
	c.c = cb
	ct := container.New(
		layout.NewVBoxLayout(),
		c.i,
		c.c,
	)
	ct.Resize(s)
	return ct
}

func NewCheckedImage(s fyne.Size) *fyne.Container {
	i := emptyCell(s)
	i.SetMinSize(s)
	c := &CheckedImage{
		Selected: true,
		i:        i,
		s:        s,
	}

	cb := widget.NewCheck("", c.tapped)
	cb.SetChecked(true)
	c.c = cb
	ct := container.New(
		layout.NewVBoxLayout(),
		c.i,
		c.c,
	)
	ct.Resize(s)
	return ct
}

type ImageSelectionTable struct {
	*fyne.Container
	size fyne.Size
}

func NewImageSelectionTable(size fyne.Size) *ImageSelectionTable {
	t := &ImageSelectionTable{
		size:      size,
		Container: container.NewHBox(),
	}
	return t

}

func NewImageSelectionTableWithImages(imgs []image.Image, size fyne.Size) *ImageSelectionTable {
	t := &ImageSelectionTable{
		Container: container.NewHBox(),
		size:      size,
	}

	for i := 0; i < len(imgs); i++ {
		ci := canvas.NewImageFromImage(imgs[i])
		t.Add(NewCheckedImageWithImage(ci, t.size))
	}

	return t
}

func (t *ImageSelectionTable) Images() []image.Image {
	im := make([]image.Image, 0)
	for _, v := range t.Container.Objects {
		c := v.(*fyne.Container).Objects[1].(*widget.Check)
		i := v.(*fyne.Container).Objects[0].(*canvas.Image)
		if c.Checked {
			im = append(im, i.Image)
		}
	}

	return im
}

func (t *ImageSelectionTable) Reset() {
	t.Container.Objects = t.Container.Objects[:0]
	t.Refresh()
	canvas.Refresh(t)
}

func (t *ImageSelectionTable) Append(img *canvas.Image) {
	i := NewCheckedImageWithImage(img, t.size)

	t.Add(i)
	t.Refresh()
	canvas.Refresh(t)
}
