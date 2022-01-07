package custom_widget

import (
	"image"
	"image/draw"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type ImageTable struct {
	widget.Table
	images                *[][]canvas.Image
	ImageCallbackFunc     func(*canvas.Image)
	IndexCallbackFunc     func(int, int)
	SetImagesCallbackFunc func(*[][]canvas.Image)
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
	canvasImg := emptyCell(imageSize)
	images := make([][]canvas.Image, 1)
	for i := 0; i < 1; i++ {
		images[i] = make([]canvas.Image, 1)
		for j := 0; j < 1; j++ {
			images[i][j] = *canvasImg
		}
	}
	imageTable.images = &images
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
	images *[][]canvas.Image,
	imageSize fyne.Size,
	nbRows, nbCols int,
	imageSelected func(*canvas.Image),
	indexSelected func(int, int),
	setImages func(*[][]canvas.Image)) *ImageTable {

	if len(*images) != nbRows || len((*images)[0]) != nbCols {
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

func (i *ImageTable) SubstitueImage(row, col int, newImage canvas.Image) {
	if row < 0 || row > len(*i.images) {
		return
	}
	if col < 0 || col > len((*i.images)[0]) {
		return
	}
	(*i.images)[row][col] = newImage
	i.UpdateCell(widget.TableCellID{Row: row, Col: col}, &newImage)
	canvas.Refresh(i)
}

func (i *ImageTable) Update(images *[][]canvas.Image, rowNumber, colNumber int) {
	i.images = images
	i.rowsNumber = rowNumber
	i.colsNumber = colNumber
	for x := 0; x < i.rowsNumber; x++ {
		for y := 0; y < i.colsNumber; y++ {
			i.UpdateCell(widget.TableCellID{Row: x, Col: y}, &(*i.images)[x][y])
		}
	}
	canvas.Refresh(i)
}

func (i *ImageTable) ImagesLength() (row int, col int) {
	return i.rowsNumber, i.colsNumber
}

func (i *ImageTable) Reset() {
	imgs := make([][]canvas.Image, 1)
	imgs[0] = make([]canvas.Image, 1)
	imgs[0][0] = *emptyCell(i.imageSize)
	i.Update(&imgs, 1, 1)
	canvas.Refresh(i)
}

func (i *ImageTable) AppendImage(image canvas.Image, rowNumber int) {
	if rowNumber >= len(*i.images) {
		return
	}
	for x := 0; x < i.rowsNumber; x++ {
		if x != rowNumber {
			(*i.images)[x] = append((*i.images)[x], *emptyCell(i.imageSize))
		}
	}
	(*i.images)[rowNumber] = append((*i.images)[rowNumber], image)
	i.colsNumber++
	i.Refresh()
	//	i.UpdateCell(widget.TableCellID{Row: rowNumber, Col: i.colsNumber}, &(*i.images)[rowNumber][i.colsNumber])
}

func (i *ImageTable) ImageCreate() fyne.CanvasObject {
	img := image.NewRGBA(image.Rectangle{image.Point{0, 0}, image.Point{int(i.Size().Width), int(i.Size().Width)}})
	r := canvas.NewImageFromImage(img)
	r.SetMinSize(i.imageSize)
	return r
}

func (i *ImageTable) ImageUpdate(id widget.TableCellID, o fyne.CanvasObject) {
	c := (*i.images)[id.Row][id.Col]
	o.(*canvas.Image).Image = c.Image
	canvas.Refresh(o)
}

func (i *ImageTable) ImageSelect(id widget.TableCellID) {
	c := (*i.images)[id.Row][id.Col]
	if i.ImageCallbackFunc != nil {
		i.ImageCallbackFunc(&c)
	}
	if i.IndexCallbackFunc != nil {
		i.IndexCallbackFunc(id.Row, id.Col)
	}
	if i.SetImagesCallbackFunc != nil {
		i.SetImagesCallbackFunc(i.images)
	}
}
