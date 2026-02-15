package table

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
)

type TableCell struct {
	widget.BaseWidget
	bg    *canvas.Rectangle
	shape *canvas.Circle
	label *widget.Label
}

func NewTableCell() *TableCell {
	bg := canvas.NewRectangle(color.Transparent)
	circle := canvas.NewCircle(color.NRGBA{R: 0, G: 150, B: 255, A: 255})
	label := widget.NewLabel("..")

	c := &TableCell{
		bg:    bg,
		shape: circle,
		label: label,
	}
	c.ExtendBaseWidget(c)
	return c
}

func (tc *TableCell) CreateRenderer() fyne.WidgetRenderer {

	objects := []fyne.CanvasObject{tc.bg, tc.shape, tc.label}
	return &tableCellRenderer{
		cell:    tc,
		objects: objects,
	}
}

func (tc *TableCell) MinSize() fyne.Size {
	return fyne.NewSize(30, 30)
}

type tableCellRenderer struct {
	cell    *TableCell
	objects []fyne.CanvasObject
}

func (tcr *tableCellRenderer) Layout(size fyne.Size) {
	tcr.cell.bg.Resize(size)

	// 40% of smaller dimension
	diameter := fyne.Min(size.Width, size.Height) * 0.4

	tcr.cell.shape.Resize(fyne.NewSize(diameter, diameter))
	tcr.cell.shape.Move(fyne.NewPos(
		//	(size.Width-diameter)/2,	// centered
		diameter,
		(size.Height-diameter)/2,
	))
}

func (tcr *tableCellRenderer) MinSize() fyne.Size {
	return tcr.cell.MinSize()
}

func (tcr *tableCellRenderer) Refresh() {
	canvas.Refresh(tcr.cell)
}

func (tcr *tableCellRenderer) Destroy() {}

func (tcr *tableCellRenderer) Objects() []fyne.CanvasObject {
	return tcr.objects
}
