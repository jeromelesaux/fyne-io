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

var (
	default20x20Size = fyne.NewSize(20, 20)
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

	csi  *canvas.Image // current selected color in palette
	pi   int           // color position in current palette
	csii *canvas.Image // current selected color available
	pci  int           // color position in current available colors
	pt   *widget.Table
	o    *canvas.Image
	m    *PixelsMap // pixels  map pointer
}

func (e *Editor) setPaletteColor() {
	e.csi.Image = fillImageColor(e.p[e.pi], default20x20Size)
	e.csi.Refresh()
}

func (e *Editor) setSelectedAvailableColor() {
	e.csii.Image = fillImageColor(e.c[e.pci], default20x20Size)
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
	if e.py > 0 {
		e.py--
	}
}
func (e *Editor) goDown() {
	if e.py < e.oi.Bounds().Max.Y {
		e.py++
	}
}
func (e *Editor) goLeft() {
	if e.px > 0 {
		e.px--
	}
}
func (e *Editor) goRight() {
	if e.px < e.oi.Bounds().Max.X {
		e.px++
	}
}

func (e *Editor) setColor(x, y int, c color.Color) {
	e.oi.(*image.NRGBA).Set(e.px+x, e.py+y, c)
}

func (e *Editor) selectColorPalette(id widget.TableCellID) {
	y := id.Col
	x := id.Row
	e.pi = y + (x * 64)
	e.m.SetColor(e.c[y])
	e.setPaletteColor()
}

func (e *Editor) selectAvailableColor(id widget.TableCellID) {
	y := id.Col
	x := id.Row
	e.pci = y + (x * 64)
	c := e.c[y]
	e.p[e.pi] = c
	e.setSelectedAvailableColor()

	cell := canvas.NewImageFromImage(fillImageColor(e.p[e.pi], fyne.NewSize(5, 5)))
	cell.SetMinSize(default20x20Size)
	e.pt.UpdateCell(id, cell)
	e.pt.Refresh()

	e.setPaletteColor()
	e.csii.Refresh()
}

func NewEditor(i image.Image, m Magnify, p color.Palette, ca color.Palette) *Editor {

	e := &Editor{
		oi:   i,
		mg:   m,
		p:    p,
		c:    ca,
		ip:   make([][]color.Color, m.HeightPixels),
		o:    canvas.NewImageFromImage(i),
		csi:  canvas.NewImageFromImage(fillImageColor(p[0], default20x20Size)),
		csii: canvas.NewImageFromImage(fillImageColor(ca[0], default20x20Size)),
	}
	for i := 0; i < m.HeightPixels; i++ {
		e.ip[i] = make([]color.Color, m.WidthPixels)
	}
	e.m = NewPixelsMap(e.mg, fyne.NewSize(5, 5), e.setColor)
	e.setImagePortion()
	return e
}

func (e *Editor) newDirectionsContainer() *fyne.Container {
	return container.New(
		layout.NewAdaptiveGridLayout(3),
		widget.NewButtonWithIcon("LEFT", theme.NavigateBackIcon(), e.goLeft),
		container.New(
			layout.NewAdaptiveGridLayout(1),
			widget.NewButtonWithIcon("UP", theme.MoveUpIcon(), e.goUp),
			widget.NewButtonWithIcon("DOWN", theme.MoveDownIcon(), e.goDown),
		),
		widget.NewButtonWithIcon("RIGHT", theme.NavigateNextIcon(), e.goRight),
	)
}

func (e *Editor) newPaletteContainer(p color.Palette, setTable func(t *widget.Table), sel func(id widget.TableCellID)) *fyne.Container {
	t := widget.NewTable(func() (int, int) {
		col := len(p) / 64
		if col == 0 {
			col = 1
		}
		row := len(p) % 64

		return col, row
	}, func() fyne.CanvasObject {
		o := canvas.NewImageFromImage(fillImageColor(color.Black, fyne.NewSize(5, 5)))
		o.SetMinSize(default20x20Size)
		return o
	}, func(id widget.TableCellID, cell fyne.CanvasObject) {
		y := id.Col

		cell.(*canvas.Image).Image = fillImageColor(p[y], fyne.NewSize(5, 5))
		cell.Refresh()
	})
	t.OnSelected = sel
	if setTable != nil {
		setTable(t)
	}

	return container.New(
		layout.NewGridLayout(1),
		t,
	)
}

func (e *Editor) applyMagnify() {}

func (e *Editor) setPaletteTable(t *widget.Table) {
	e.pt = t
}

func (e *Editor) NewEditor() *fyne.Container {

	return container.New(
		layout.NewGridLayoutWithColumns(2),

		container.New(
			layout.NewGridLayoutWithRows(2),
			e.m.NewPixelsMap(),
			e.o,
		),
		container.New(
			layout.NewGridLayoutWithRows(8),

			widget.NewLabel("Your palette :"),
			e.newPaletteContainer(e.p, e.setPaletteTable, e.selectColorPalette),
			container.New(
				layout.NewAdaptiveGridLayout(1),
				widget.NewLabel("Selected color from your palette :"),
				e.csi,
			),

			widget.NewLabel("Color available :"),
			e.newPaletteContainer(e.c, nil, e.selectAvailableColor),
			container.New(
				layout.NewAdaptiveGridLayout(1),
				widget.NewLabel("Selected color from available colors :"),
				e.csii,
			),
			container.New(
				layout.NewGridLayoutWithColumns(2),
				widget.NewLabel("Magnify :"),
				widget.NewSelect([]string{"x2", "x4", "x8"}, func(s string) {
					switch s {
					case "x2":
						e.mg = MagnifyX2
					case "x4":
						e.mg = MagnifyX4
					case "x8":
						e.mg = MagnifyX8
					default:
						return
					}
					e.applyMagnify()
				}),
			),
			e.newDirectionsContainer(),
		),
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
	o.SetMinSize(fyne.NewSize(7, 7))
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
