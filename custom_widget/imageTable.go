package custom_widget

import (
	"image"
	"image/draw"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type ImagesTableCollection struct {
	images          []*canvas.Image
	ImagesPerRow    int
	ImagesPerColumn int
	Left            int
}

func NewImagesTableCollection(row, col int, size fyne.Size) *ImagesTableCollection {
	itc := &ImagesTableCollection{
		ImagesPerRow:    row,
		ImagesPerColumn: col,
		images:          make([]*canvas.Image, row*col)}
	for i := 0; i < itc.ImagesPerColumn*itc.ImagesPerRow; i++ {
		itc.images[i] = emptyCell(size)
	}
	return itc
}

func (i *ImagesTableCollection) At(row, col int) *canvas.Image {
	if row < i.ImagesPerRow && col < i.ImagesPerColumn {
		return i.images[row+(row*col)]
	}
	return &canvas.Image{}
}

func (i *ImagesTableCollection) Set(row, col int, img *canvas.Image) {
	if row < i.ImagesPerRow && col < i.ImagesPerColumn {
		i.images[row+(row*col)] = img
	}
}

type ImageTable struct {
	widget.Table
	images                *ImagesTableCollection
	ImageCallbackFunc     func(*canvas.Image)
	IndexCallbackFunc     func(int, int)
	SetImagesCallbackFunc func(*ImagesTableCollection)
	imageSize             fyne.Size
	rowsNumber            int
	colsNumber            int
}

func emptyCell(imageSize fyne.Size) *canvas.Image {
	img := image.NewRGBA(image.Rectangle{image.Point{0, 0}, image.Point{int(imageSize.Width), int(imageSize.Height)}})
	bg := theme.BackgroundColor()
	draw.Draw(img, img.Bounds(), &image.Uniform{bg}, image.Point{0, 0}, draw.Src)
	canvasImg := canvas.NewImageFromImage(img)
	return canvasImg
}

func NewEmptyImageTable(imageSize fyne.Size) *ImageTable {
	imageTable := &ImageTable{}
	imageTable.images = NewImagesTableCollection(1, 1, imageSize)
	imageTable.ImageCallbackFunc = nil
	imageTable.IndexCallbackFunc = nil
	imageTable.SetImagesCallbackFunc = nil
	imageTable.CreateCell = imageTable.ImageCreate
	imageTable.Length = imageTable.ImagesLength
	imageTable.UpdateCell = imageTable.ImageUpdate
	imageTable.OnSelected = imageTable.ImageSelect
	imageTable.imageSize = imageSize
	imageTable.rowsNumber = 1
	imageTable.colsNumber = 1
	imageTable.ExtendBaseWidget(imageTable)
	return imageTable
}

func NewImageTable(
	images *ImagesTableCollection,
	imageSize fyne.Size,
	nbRows, nbCols int,
	imageSelected func(*canvas.Image),
	indexSelected func(int, int),
	setImages func(*ImagesTableCollection)) *ImageTable {

	if images.ImagesPerRow != nbRows || images.ImagesPerColumn != nbCols {
		panic("images matrix must corresponds to number of rows and columns")
	}

	imageTable := &ImageTable{}
	imageTable.images = images
	imageTable.ImageCallbackFunc = imageSelected
	imageTable.IndexCallbackFunc = indexSelected
	imageTable.SetImagesCallbackFunc = setImages
	imageTable.CreateCell = imageTable.ImageCreate
	imageTable.Length = imageTable.ImagesLength
	imageTable.UpdateCell = imageTable.ImageUpdate
	imageTable.OnSelected = imageTable.ImageSelect
	imageTable.imageSize = imageSize
	imageTable.rowsNumber = nbRows
	imageTable.colsNumber = nbCols
	imageTable.ExtendBaseWidget(imageTable)

	return imageTable
}

func (i *ImageTable) Images() [][]image.Image {
	imgs := make([][]image.Image, i.rowsNumber)
	for x := 0; x < i.rowsNumber; x++ {
		imgs[x] = make([]image.Image, i.colsNumber)
		for y := 0; y < i.colsNumber; y++ {
			imgs[x][y] = (*i.images).At(x, y).Image
		}
	}
	return imgs
}

func (i *ImageTable) SubstitueImage(row, col int, newImage *canvas.Image) {
	if row < 0 || row > i.images.ImagesPerRow {
		return
	}
	if col < 0 || col > i.images.ImagesPerColumn {
		return
	}
	i.images.Set(row, col, newImage)
	i.UpdateCell(widget.TableCellID{Row: row, Col: col}, newImage)
	i.Refresh()
	canvas.Refresh(i)
}

func (i *ImageTable) Update(images *ImagesTableCollection, rowNumber, colNumber int) {
	i.images = images
	i.rowsNumber = rowNumber
	i.colsNumber = colNumber
	for x := 0; x < i.rowsNumber; x++ {
		for y := 0; y < i.colsNumber; y++ {
			i.UpdateCell(widget.TableCellID{Row: x, Col: y}, i.images.At(x, y))
		}
	}
	i.Refresh()
	canvas.Refresh(i)
}

func (i *ImageTable) ImagesLength() (row int, col int) {
	return i.rowsNumber, i.colsNumber
}

func (i *ImageTable) Reset() {
	imgs := NewImagesTableCollection(1, 1, i.imageSize)
	i.Update(imgs, 1, 1)
	i.Refresh()
	canvas.Refresh(i)
}

func (i *ImageTable) ImageCreate() fyne.CanvasObject {
	img := image.NewRGBA(image.Rectangle{image.Point{0, 0}, image.Point{int(i.Size().Width), int(i.Size().Width)}})
	r := canvas.NewImageFromImage(img)
	r.SetMinSize(i.imageSize)
	return r
}

func (i *ImageTable) ImageUpdate(id widget.TableCellID, o fyne.CanvasObject) {
	c := i.images.At(id.Row, id.Col)
	o.(*canvas.Image).Image = c.Image
	canvas.Refresh(o)
}

func (i *ImageTable) ImageSelect(id widget.TableCellID) {
	c := i.images.At(id.Row, id.Col)
	if i.ImageCallbackFunc != nil {
		i.ImageCallbackFunc(c)
	}
	if i.IndexCallbackFunc != nil {
		i.IndexCallbackFunc(id.Row, id.Col)
	}
	if i.SetImagesCallbackFunc != nil {
		i.SetImagesCallbackFunc(i.images)
	}
}
