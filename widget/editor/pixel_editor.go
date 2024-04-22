package editor

import (
	"image"
	"image/color"
	"image/draw"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
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
	mg Magnify         // magnify used
	p  color.Palette   // current palette used
	c  color.Palette   // colors available
	oi image.Image     // original image
	ip [][]color.Color // portion extracted from original image to display
	px int             // position in X from the original image
	py int             // position in Y from the original image

	// Editor widgets
	up    *widget.Button // go upper in the original image
	down  *widget.Button // go down in the original image
	right *widget.Button // go right
	left  *widget.Button // go left
	m     *PixelsMap     // pixels  map pointer
}

func (e *Editor) setImagePortion() {
	for y := 0; y < e.mg.HeightPixels; y++ {
		for x := 0; x < e.mg.WidthPixels; x++ {
			c := e.oi.At(e.px+x, e.py+y)
			e.ip[y][x] = c
		}
	}
	e.m.SetColors(e.ip)
}

func (e *Editor) goUp() {

}
func (e *Editor) goDown() {

}
func (e *Editor) goLeft() {

}
func (e *Editor) goRight() {

}

func (e *Editor) setColor(x, y int, c color.Color) {

}

func NewEditor(i image.Image, m Magnify, p color.Palette, ca color.Palette) *Editor {

	e := &Editor{
		oi: i,
		mg: m,
		p:  p,
		c:  ca,
		ip: make([][]color.Color, m.HeightPixels),
	}
	for i := 0; i < m.HeightPixels; i++ {
		e.ip[i] = make([]color.Color, m.WidthPixels)
	}
	e.m = NewPixelsMap(e.mg, fyne.NewSize(5, 5), e.setColor)
	e.setImagePortion()
	return e
}

func (e *Editor) newDirectionsContainer() *fyne.Container {
	e.up = widget.NewButtonWithIcon("UP", theme.MoveUpIcon(), e.goUp)
	e.down = widget.NewButtonWithIcon("DOWN", theme.MoveDownIcon(), e.goDown)
	e.left = widget.NewButtonWithIcon("LEFT", theme.NavigateBackIcon(), e.goLeft)
	e.right = widget.NewButtonWithIcon("LEFT", theme.NavigateNextIcon(), e.goRight)
	return container.New(
		layout.NewAdaptiveGridLayout(3),
		e.left,
		container.New(
			layout.NewAdaptiveGridLayout(1),
			e.up,
			e.down,
		),
		e.right,
	)
}

func (e *Editor) NewEditor() *fyne.Container {

	return container.New(
		layout.NewGridLayoutWithColumns(2),

		e.m.NewPixelsMap(),
		e.newDirectionsContainer(),
	)
}

type PixelsMap struct {
	mg       Magnify
	sz       fyne.Size
	px       *widget.Table
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
	o := canvas.NewImageFromImage(fillImageColor(color.Black, p.sz))
	o.SetMinSize(fyne.NewSize(8, 8))
	return o
}

func (p *PixelsMap) updateCell(id widget.TableCellID, cell fyne.CanvasObject) {
	x := id.Col
	y := id.Row
	cell.(*canvas.Image).Image = fillImageColor(p.mc[x][y], p.sz)
	cell.Refresh()
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

func (pm *PixelsMap) NewPixelsMap() *fyne.Container {
	return container.New(
		layout.NewGridLayout(1),
		pm.px,
	)
}

func NewPixelsMap(m Magnify, sz fyne.Size, s func(x, y int, c color.Color)) *PixelsMap {
	p := &PixelsMap{mg: m, setColor: s, sz: sz, sc: color.Black}
	p.px = widget.NewTable(
		p.length,
		p.createCell,
		p.updateCell,
	)
	p.px.OnSelected = p.onSelected
	p.px.OnUnselected = p.onUnselected

	p.mc = make([][]color.Color, p.mg.HeightPixels)
	for i := 0; i < p.mg.HeightPixels; i++ {
		p.mc[i] = make([]color.Color, p.mg.WidthPixels)
		for j := 0; j < p.mg.WidthPixels; j++ {
			p.mc[i][j] = color.Black
		}
	}

	return p
}
