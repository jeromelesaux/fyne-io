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
	default10x10Size = fyne.NewSize(10, 10)
)

type Magnify struct {
	Display      float32
	Value        int
	WidthPixels  int
	HeightPixels int
}

var (
	MagnifyX2 = Magnify{
		Display:      7,
		Value:        2,
		WidthPixels:  64,
		HeightPixels: 64,
	}
	MagnifyX4 = Magnify{
		Display:      15,
		Value:        4,
		WidthPixels:  32,
		HeightPixels: 32,
	}
	MagnifyX8 = Magnify{
		Display:      25,
		Value:        8,
		WidthPixels:  16,
		HeightPixels: 16,
	}
)

type Editor struct {
	co *fyne.Container
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
	o    *ClickableImage
	m    *PixelsMap // pixels  map pointer
	sv   func(i image.Image, p color.Palette)
	w    fyne.Window
}

func (e *Editor) onTypedKey(k *fyne.KeyEvent) {
	switch k.Name {
	case "Down":
		e.goDown()
	case "Up":
		e.goUp()
	case "Right":
		e.goRight()
	case "Left":
		e.goLeft()
	case "M":
		switch e.mg {
		case MagnifyX2:
			e.mg = MagnifyX4
		case MagnifyX4:
			e.mg = MagnifyX8
		case MagnifyX8:
			e.mg = MagnifyX2
		}
		e.syncMap()
	case "Space":
		x, y := e.m.SelectedCell()
		e.m.onSelected(widget.TableCellID{Col: x, Row: y})
	}

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
			e.ip[x][y] = c
		}
	}
	e.m.SetColors(e.ip)
}

func (e *Editor) goUp() {
	if e.py > 0 {
		e.py--
	}
	e.syncMap()
}
func (e *Editor) goDown() {
	if e.py < e.oi.Bounds().Max.Y {
		e.py++
	}
	e.syncMap()
}
func (e *Editor) goLeft() {
	if e.px > 0 {
		e.px--
	}
	e.syncMap()
}
func (e *Editor) goRight() {
	if e.px < e.oi.Bounds().Max.X {
		e.px++
	}
	e.syncMap()
}

func (e *Editor) setColor(x, y int, c color.Color) {
	switch e.oi.(type) {
	case *image.NRGBA:
		e.oi.(*image.NRGBA).Set(e.px+x, e.py+y, c)
	case *image.RGBA:
		e.oi.(*image.RGBA).Set(e.px+x, e.py+y, c)
	case *image.Alpha:
		e.oi.(*image.Alpha).Set(e.px+x, e.py+y, c)
	case *image.Alpha16:
		e.oi.(*image.Alpha16).Set(e.px+x, e.py+y, c)
	case *image.RGBA64:
		e.oi.(*image.RGBA64).Set(e.px+x, e.py+y, c)
	case *image.NRGBA64:
		e.oi.(*image.NRGBA64).Set(e.px+x, e.py+y, c)

	default:
		return
	}
}

func (e *Editor) selectColorPalette(id widget.TableCellID) {
	y := id.Col
	x := id.Row
	e.pi = y + (x * 64)
	e.m.SetColor(e.p[y])
	e.setPaletteColor()
}

func (e *Editor) replaceOneColor(cOld, cNew color.Color) {
	for x := 0; x < e.oi.Bounds().Max.X; x++ {
		for y := 0; y < e.oi.Bounds().Max.Y; y++ {
			c := e.oi.At(x, y)
			if colorsAreEqual(c, cOld) {
				switch e.oi.(type) {
				case *image.NRGBA:
					e.oi.(*image.NRGBA).Set(x, y, cNew)
				case *image.RGBA:
					e.oi.(*image.RGBA).Set(x, y, cNew)
				case *image.Alpha:
					e.oi.(*image.Alpha).Set(x, y, cNew)
				case *image.Alpha16:
					e.oi.(*image.Alpha16).Set(x, y, cNew)
				case *image.RGBA64:
					e.oi.(*image.RGBA64).Set(x, y, cNew)
				case *image.NRGBA64:
					e.oi.(*image.NRGBA64).Set(x, y, cNew)
				}
			}
		}
	}
	e.syncMap()
}

func colorsAreEqual(c0, c1 color.Color) bool {
	r0, g0, b0, _ := c0.RGBA()
	r1, g1, b1, _ := c1.RGBA()
	if r0 != r1 {
		return false
	}
	if b0 != b1 {
		return false
	}
	if g0 != g1 {
		return false
	}
	return true
}

func (e *Editor) selectAvailableColor(id widget.TableCellID) {
	y := id.Col
	x := id.Row
	e.pci = y + (x * 64)
	c := e.c[e.pci]
	c0 := e.p[e.pi]
	e.p[e.pi] = c
	e.replaceOneColor(c0, c) // replace all the initial color by the new one
	e.m.SetColor(e.p[e.pi])
	e.setSelectedAvailableColor()

	cell := canvas.NewImageFromImage(fillImageColor(e.p[e.pi], fyne.NewSize(5, 5)))
	cell.SetMinSize(default20x20Size)
	id = widget.TableCellID{
		Row: 0,
		Col: e.pi,
	}
	e.pt.UpdateCell(id, cell)
	e.pt.Refresh()

	e.setPaletteColor()
	e.csii.Refresh()
}

func (e *Editor) posSquareSelect(x, y float32) {
	e.px = int(x) - (e.mg.WidthPixels / 2)
	e.py = int(y) - (e.mg.HeightPixels / 2)
	e.syncMap()
}

func NewEditor(i image.Image, m Magnify, p color.Palette, ca color.Palette, s func(image.Image, color.Palette), w fyne.Window) *Editor {

	if len(p) == 0 {
		p = append(p, color.Black, color.Black, color.Black, color.Black)
	}

	e := &Editor{
		oi:   i,
		mg:   m,
		p:    p,
		c:    ca,
		ip:   make([][]color.Color, m.WidthPixels),
		csi:  canvas.NewImageFromImage(fillImageColor(p[0], default20x20Size)),
		csii: canvas.NewImageFromImage(fillImageColor(ca[0], default20x20Size)),
		sv:   s,
		w:    w,
	}

	e.o = NewClickableImage(e.oi, e.posSquareSelect)
	for i := 0; i < m.WidthPixels; i++ {
		e.ip[i] = make([]color.Color, m.HeightPixels)
	}
	e.w.Canvas().SetOnTypedKey(e.onTypedKey)
	e.m = NewPixelsMap(e.mg, fyne.NewSize(5, 5), e.setColor)

	e.setImagePortion()
	return e
}

func (e *Editor) NewImage(i image.Image) {
	e.oi = i
	e.syncMap()
}

func (e *Editor) NewPalette(p color.Palette) {
	e.p = p
	e.syncMap()
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
		if row == 0 {
			row = 64
		}
		return col, row
	}, func() fyne.CanvasObject {
		o := canvas.NewImageFromImage(fillImageColor(color.Black, fyne.NewSize(5, 5)))
		if len(p) > 64 {
			o.SetMinSize(default10x10Size)
		} else {
			o.SetMinSize(default20x20Size)
		}
		return o
	}, func(id widget.TableCellID, cell fyne.CanvasObject) {
		y := id.Col
		x := id.Row
		cell.(*canvas.Image).Image = fillImageColor(p[y+(x*64)], fyne.NewSize(5, 5))
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

func (e *Editor) squareSelect() {
	i := image.NewNRGBA(e.oi.Bounds())
	draw.Draw(i, i.Bounds(), e.oi, image.Point{0, 0}, draw.Src)
	for x := e.px; x < e.px+e.mg.WidthPixels; x++ {
		i.Set(x, e.py, color.Black)
		i.Set(x, e.py+e.mg.HeightPixels, color.Black)
	}
	for y := e.py; y < e.py+e.mg.HeightPixels; y++ {
		i.Set(e.px, y, color.Black)
		i.Set(e.px+e.mg.WidthPixels, y, color.Black)
	}
	e.o.SetImage(i)
	e.o.Refresh()
}

func (e *Editor) syncMap() {
	// modify the size of ip
	e.ip = e.ip[0:e.mg.WidthPixels]
	for i := 0; i < e.mg.WidthPixels; i++ {
		e.ip[i] = e.ip[i][0:e.mg.HeightPixels]
	}
	e.m.SetMagnify(e.mg)
	e.setImagePortion()
	e.m.px.Refresh()
	e.squareSelect()
}

func (e *Editor) setPaletteTable(t *widget.Table) {
	e.pt = t
}

func (e *Editor) NewEditor() *fyne.Container {

	e.co = container.New(
		layout.NewGridLayoutWithColumns(2),

		container.New(
			layout.NewGridLayoutWithRows(2),
			e.m.NewPixelsMap(),
			e.o,
		),
		container.New(
			layout.NewGridLayoutWithRows(9),

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
					e.syncMap()
				}),
			),
			e.newDirectionsContainer(),
			container.New(
				layout.NewGridLayoutWithColumns(2),
				widget.NewButtonWithIcon("Save", theme.FileImageIcon(), func() {
					if e.sv != nil {
						e.sv(e.oi, e.p)
					}
					e.co.Hide()
				}),
				widget.NewButtonWithIcon("Cancel", theme.CancelIcon(), func() { e.co.Hide() }),
			),
		),
	)
	return e.co
}

/*
 PixelsMap widget which displays pixels part of image
*/

type PixelsMap struct {
	mg       Magnify
	sz       fyne.Size
	px       *widget.Table
	sc       color.Color
	mc       [][]color.Color
	setColor func(x, y int, c color.Color)
	x        int
	y        int
}

func (p *PixelsMap) SetColors(cs [][]color.Color) {
	p.mc = cs
}

func (p *PixelsMap) SetColor(c color.Color) {
	p.sc = c
}

func (p *PixelsMap) SetMagnify(m Magnify) {
	p.mg = m
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
	o.SetMinSize(fyne.NewSize(p.mg.Display, p.mg.Display))
	return o
}

func (p *PixelsMap) updateCell(id widget.TableCellID, cell fyne.CanvasObject) {
	x := id.Col
	y := id.Row
	cell.(*canvas.Image).Image = fillImageColor(p.mc[x][y], p.sz)
	cell.Refresh()
}
func (p *PixelsMap) onSelected(id widget.TableCellID) {
	p.x = id.Col
	p.y = id.Row
	p.mc[p.x][p.y] = p.sc
	if p.setColor != nil {
		p.setColor(p.x, p.y, p.sc)
	}
	p.px.Refresh()
}

func (p *PixelsMap) onUnselected(id widget.TableCellID) {

}

func (p *PixelsMap) SelectedCell() (x int, y int) {
	x = p.x
	y = p.y
	return
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

type ClickableImage struct {
	*widget.Icon
	i      *canvas.Image
	tapped func(float32, float32)
}

func NewClickableImage(i image.Image, tapped func(float32, float32)) *ClickableImage {
	c := &ClickableImage{
		Icon:   &widget.Icon{},
		i:      canvas.NewImageFromImage(i),
		tapped: tapped,
	}
	c.ExtendBaseWidget(c)
	return c
}

func (c *ClickableImage) SetImage(i image.Image) {
	c.i.Image = i
	c.i.Refresh()
}

func (c *ClickableImage) Tapped(pe *fyne.PointEvent) {
	if c.tapped != nil {
		x := pe.Position.X * float32(c.i.Image.Bounds().Max.X) / c.i.Size().Width
		y := pe.Position.Y * float32(c.i.Image.Bounds().Max.Y) / c.i.Size().Height
		(c.tapped)(x, y)
	}
}

func (c *ClickableImage) TappedSecondary(_ *fyne.PointEvent) {

}

func (ci *ClickableImage) CreateRenderer() fyne.WidgetRenderer {
	//ci.BaseWidget.ExtendBaseWidget(ci)
	return &clickableImageRenderer{
		image: ci.i,
		objs:  []fyne.CanvasObject{ci.i},
	}
}

func (ci *ClickableImage) Move(position fyne.Position) {
	ci.Icon.Move(position)
}

type clickableImageRenderer struct {
	image *canvas.Image
	objs  []fyne.CanvasObject
}

func (ci *clickableImageRenderer) Destroy() {
	ci.image = nil
}

func (ci *clickableImageRenderer) MinSize() fyne.Size {
	return ci.image.MinSize()
}

func (ci *clickableImageRenderer) Objects() []fyne.CanvasObject {
	return ci.objs
}

func (ci *clickableImageRenderer) Refresh() {
	ci.image.Refresh()
}

func (ci *clickableImageRenderer) Layout(size fyne.Size) {
	ci.image.Resize(size)
}

func (ci *clickableImageRenderer) Resize(size fyne.Size) {
	ci.image.Resize(size)
}
