package orderer

import (
	"fmt"
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

const (
	colorsByColumn = 8
)

var (
	default20x20Size = fyne.NewSize(20, 20)
	default30x30Size = fyne.NewSize(30, 30)
	default10x10Size = fyne.NewSize(10, 10)
)

type Orderer struct {
	np     color.Palette // new color palette, final result
	cp     color.Palette // current color palette to order
	c      color.Color   //  selected color
	wnp    *widget.Table
	wcp    *widget.Table
	cc     *canvas.Image
	l      *widget.Label
	setPal func(color.Palette) // set the new palette callback
}

func NewOrderer(p color.Palette, setPalette func(color.Palette)) *Orderer {
	o := &Orderer{
		cp:     p,
		setPal: setPalette,
		l:      widget.NewLabel("Select an upper color "),
		cc:     canvas.NewImageFromImage(fillImageColor(color.Black, default10x10Size)),
	}
	// fill the new palette with black color
	// same size as the input palette
	for i := 0; i < len(p); i++ {
		o.np = append(o.np, color.Black)
	}

	o.wcp = widget.NewTable(func() (int, int) {
		col := len(o.cp)
		if len(o.cp) > colorsByColumn {
			col = colorsByColumn
		}

		row := len(o.cp) / colorsByColumn
		if len(o.cp)%colorsByColumn != 0 {
			row++
		}
		if row == 0 {
			row = 1
		}
		return row, col
	}, func() fyne.CanvasObject {
		im := canvas.NewImageFromImage(fillImageColor(color.Black, fyne.NewSize(20, 20)))
		im.SetMinSize(default30x30Size)

		return im
	}, func(id widget.TableCellID, cell fyne.CanvasObject) {
		y := id.Col
		x := id.Row
		if y+(x*colorsByColumn) >= len(o.cp) {
			cell.(*canvas.Image).Image = fillImageColor(color.Black, fyne.NewSize(20, 20))
		} else {
			cell.(*canvas.Image).Image = fillImageColor((o.cp)[y+(x*colorsByColumn)], fyne.NewSize(20, 20))
		}
		cell.Refresh()
	})
	o.wcp.OnSelected = o.selectedColor

	o.wnp = widget.NewTable(func() (int, int) {
		col := len(o.np)
		if len(o.np) > colorsByColumn {
			col = colorsByColumn
		}

		row := len(o.np) / colorsByColumn
		if len(o.np)%colorsByColumn != 0 {
			row++
		}
		if row == 0 {
			row = 1
		}
		return row, col
	}, func() fyne.CanvasObject {
		im := canvas.NewImageFromImage(fillImageColor(color.Black, fyne.NewSize(20, 20)))
		im.SetMinSize(default30x30Size)

		return im
	}, func(id widget.TableCellID, cell fyne.CanvasObject) {
		y := id.Col
		x := id.Row
		if y+(x*colorsByColumn) >= len(o.cp) {
			cell.(*canvas.Image).Image = fillImageColor(color.Black, fyne.NewSize(20, 20))
		} else {
			cell.(*canvas.Image).Image = fillImageColor((o.np)[y+(x*colorsByColumn)], fyne.NewSize(20, 20))
		}
		cell.Refresh()
	})
	o.wnp.OnSelected = o.positionSelectedColor

	return o
}

func (o *Orderer) NewOrderer() *fyne.Container {
	return container.New(
		layout.NewVBoxLayout(),
		container.New(
			layout.NewGridLayoutWithRows(4),
			o.wcp,
			container.New(
				layout.NewGridLayoutWithRows(2),
				widget.NewButtonWithIcon("Clone", theme.ContentCopyIcon(), func() {
					copy(o.np, o.cp)
					o.wnp.Refresh()
				}),
				container.New(
					layout.NewGridLayoutWithColumns(2),
					o.l,
					o.cc,
				)),
			o.wnp,

			widget.NewButton("Apply", func() {
				if o.setPal != nil {
					o.setPal(o.np)
				}
			}),
		))
}

func (o *Orderer) positionSelectedColor(id widget.TableCellID) {
	if id.Col < 0 || id.Row < 0 {
		return
	}
	y := id.Col
	x := id.Row
	if len(o.np) <= y+(x*colorsByColumn) {
		return
	}
	o.np[y+(x*colorsByColumn)] = o.c
	o.wnp.Refresh()
}

func (o *Orderer) selectedColor(id widget.TableCellID) {
	if id.Col < 0 || id.Row < 0 {
		return
	}
	y := id.Col
	x := id.Row
	if len(o.cp) <= y+(x*colorsByColumn) {
		return
	}
	o.c = o.cp[y+(x*colorsByColumn)]
	o.l.SetText(fmt.Sprintf("You select the color [%d]", y+(x*colorsByColumn)+1))
	o.cc.Image = fillImageColor(o.c, default10x10Size)
	o.l.Refresh()
	o.cc.Refresh()
}

func fillImageColor(c color.Color, s fyne.Size) image.Image {
	im := image.NewRGBA(image.Rect(
		0, 0,
		int(s.Height), int(s.Width),
	))
	draw.Draw(im, im.Bounds(), &image.Uniform{c}, image.Pt(0, 0), draw.Src)
	return im
}
