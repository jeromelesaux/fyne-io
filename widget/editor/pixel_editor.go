package editor

import (
	"image"
	"image/color"
	"image/draw"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	w "github.com/jeromelesaux/fyne-io/widget"
)

type Magnify struct {
	Value        int
	WidthPixels  int
	HeightPixels int
}

var (
	MagnifyX2 = Magnify{
		Value:        2,
		WidthPixels:  64,
		HeightPixels: 64,
	}
	MagnifyX4 = Magnify{
		Value:        4,
		WidthPixels:  32,
		HeightPixels: 32,
	}
	MagnifyX8 = Magnify{
		Value:        8,
		WidthPixels:  16,
		HeightPixels: 16,
	}
)

type Editor struct {
	mg Magnify       // magnify used
	p  color.Palette // current palette used
	c  color.Palette // colors available
	oi image.Image
	ip [][]color.Color // portion extracted from original image to display
	px int             // position in X from the original image
	py int             // position in Y from the original image

	// Editor widgets
	up    w.ButtonColored // go upper in the original image
	down  w.ButtonColored // go down in the original image
	right w.ButtonColored // go right
	left  w.ButtonColored // go left
	m     *PixelsMap      // pixels  map pointer
}

func NewEditor() *fyne.Container {

	return container.New(
		layout.NewAdaptiveGridLayout(0),
	)
}

type PixelsMap struct {
	mg       Magnify
	sz       fyne.Size
	px       widget.Table
	sc       color.Color
	mc       [][]color.Color
	setColor func(x, y int, c color.Color)
}

func (p *PixelsMap) SetColors(cs [][]color.Color) {
	p.mc = cs
}

func (p *PixelsMap) SetColor(c color.Color) {
	p.sc = c
}

func (p *PixelsMap) length() (int, int) {
	return p.mg.HeightPixels, p.mg.WidthPixels
}

func fillImageColor(c color.Color, s fyne.Size) image.Image {
	im := image.NewRGBA(image.Rect(
		0, 0,
		int(s.Height), int(s.Width),
	))
	draw.Draw(im, im.Bounds(), &image.Uniform{c}, image.Pt(0, 0), draw.Src)
	return im
}

func (p *PixelsMap) createCell() fyne.CanvasObject {
	return canvas.NewImageFromImage(fillImageColor(color.Black, p.sz))
}

func (p *PixelsMap) updateCell(id widget.TableCellID, cell fyne.CanvasObject) {
	x := id.Col
	y := id.Row
	cell.(*canvas.Image).Image = fillImageColor(p.mc[x][y], p.sz)
	p.px.Refresh()
}
func (p *PixelsMap) onSelected(id widget.TableCellID) {
	x := id.Col
	y := id.Row
	p.mc[x][y] = p.sc
	if p.setColor != nil {
		p.setColor(x, y, p.sc)
	}
	p.px.Refresh()
}

func (p *PixelsMap) onUnselected(id widget.TableCellID) {

}

func NewPixelsMap(m Magnify, s func(x, y int, c color.Color)) *PixelsMap {
	p := &PixelsMap{mg: m, setColor: s}
	p.px = *widget.NewTable(
		p.length,
		p.createCell,
		p.updateCell,
	)
	p.px.OnSelected = p.onSelected
	p.px.OnUnselected = p.onUnselected

	p.mc = make([][]color.Color, p.mg.HeightPixels)
	for i := 0; i < p.mg.HeightPixels; i++ {
		p.mc[i] = make([]color.Color, p.mg.WidthPixels)
	}

	return p
}
