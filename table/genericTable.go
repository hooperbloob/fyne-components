package table

import (
	"image/color"
	"sort"
	"strings"
	"tableApp1/meta"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

type Column[T any] struct {
	width         int
	field         *meta.FieldDescriptor[T]
	alignment     fyne.TextAlign
	colorSelector func(T) color.Color
}

func (col *Column[T]) IsIcon() bool {
	return col.colorSelector != nil
}

func (col *Column[T]) ColorFor(item T) color.Color {
	return col.colorSelector(item)
}

func (col *Column[T]) StringValueFor(item T) string {
	return col.field.StringValueFor(item)
}

func NewColumn[T any](width int, field *meta.FieldDescriptor[T], valueAlignment fyne.TextAlign, colorSelector func(T) color.Color) Column[T] {

	return Column[T]{
		width:         width,
		field:         field,
		alignment:     valueAlignment,
		colorSelector: colorSelector,
	}
}

// ========================================================================================================================================
type GenericTable[T any] struct {
	widget.BaseWidget
	data         []T
	columns      []Column[T]
	table        *widget.Table
	selectedRows map[int]bool
	newItemFunc  func() T
	sortCol      int
	sortAsc      bool
}

func (gTable *GenericTable[T]) SetColumnWidths() {

	for index, col := range gTable.columns {
		gTable.table.SetColumnWidth(index, float32(col.width))
	}
}

var selectionColor = color.NRGBA{R: 100, G: 150, B: 255, A: 80}

func NewGenericTable[T any](columns []Column[T], newItemFunc func() T) *GenericTable[T] {
	gt := &GenericTable[T]{
		columns:      columns,
		newItemFunc:  newItemFunc,
		selectedRows: make(map[int]bool),
	}

	gt.table = widget.NewTable(
		func() (int, int) {
			return len(gt.data), len(gt.columns)
		},

		func() fyne.CanvasObject {
			return NewTableCell()
		},
		func(id widget.TableCellID, obj fyne.CanvasObject) {
			cell := obj.(*TableCell)
			label := cell.label

			column := gt.columns[id.Col]
			item := gt.data[id.Row]

			label.SetText(column.StringValueFor(item))
			label.Alignment = column.alignment

			if column.IsIcon() {
				clr := column.ColorFor(item)
				cell.shape.FillColor = clr
				cell.shape.Show()
				cell.label.Hide()
			} else {
				cell.shape.Hide()
				cell.label.Show()
			}
			cell.shape.Refresh()
		},
	)

	gt.setupHeaders()
	gt.setupHandlers()

	gt.ExtendBaseWidget(gt)
	return gt
}

func (gt *GenericTable[T]) setupHeaders() {

	gt.table.ShowHeaderColumn = false
	gt.table.ShowHeaderRow = true
	gt.table.CreateHeader = func() fyne.CanvasObject {
		return widget.NewLabel("Header")
	}
	gt.table.UpdateHeader = func(id widget.TableCellID, cell fyne.CanvasObject) {
		label := cell.(*widget.Label)
		if id.Col >= 0 && id.Col < len(gt.columns) {
			label.SetText(gt.columns[id.Col].field.Label)
			label.TextStyle = fyne.TextStyle{Bold: true}
		}
	}
}

func (gt *GenericTable[T]) setupHandlers() {

	gt.table.OnSelected = func(id widget.TableCellID) {
		gt.selectedRows[id.Row] = true
	}

	gt.table.OnUnselected = func(id widget.TableCellID) {
		delete(gt.selectedRows, id.Row)
	}
}

func (gt *GenericTable[T]) SetData(data []T) {
	gt.data = data
	gt.selectedRows = make(map[int]bool)
	gt.table.Refresh()
}

func (gt *GenericTable[T]) GetData() []T {
	return gt.data
}

func (gt *GenericTable[T]) AddItem(item T) {
	gt.data = append(gt.data, item)
	gt.table.Refresh()
}

// Replaces the item at the index with the new one
func (gt *GenericTable[T]) ItemEdited(idx int, item T) {
	gt.data[idx] = item
	gt.table.Refresh()
}

func (gt *GenericTable[T]) SelectedItems() map[int]T {

	if len(gt.selectedRows) == 0 {
		return map[int]T{}
	}

	// Build a map of indices to delete
	selected := make(map[int]T)
	for key := range gt.selectedRows {
		selected[key] = gt.data[key]
	}
	return selected
}

// DeleteSelected removes all selected items from the table
func (gt *GenericTable[T]) DeleteSelected() int {
	if len(gt.selectedRows) == 0 {
		return 0
	}

	// Build a map of indices to delete
	toDelete := make(map[int]bool)
	for key := range gt.selectedRows {
		toDelete[key] = true
	}

	// Create new slice without deleted items
	newData := make([]T, 0, len(gt.data)-len(toDelete))
	for i, item := range gt.data {
		if !toDelete[i] {
			newData = append(newData, item)
		}
	}

	deleted := len(gt.data) - len(newData)
	gt.data = newData
	gt.selectedRows = make(map[int]bool)
	gt.table.Refresh()

	return deleted
}

func (gt *GenericTable[T]) GetSelectedCount() int {
	return len(gt.selectedRows)
}

func (gt *GenericTable[T]) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(gt.table)
}

func (gt *GenericTable[T]) asLineOn(sb *strings.Builder, item T, separator string) {

	sb.WriteString(gt.columns[0].field.Accessor(item))

	for i := 1; i < len(gt.columns); i++ {
		sb.WriteString(separator)
		sb.WriteString(gt.columns[i].field.Accessor(item))
	}
}

func (gt *GenericTable[T]) SelectionAsString(columnSeparator string, lineSeparator string) string {

	var sb strings.Builder

	for _, item := range gt.SelectedItems() {
		gt.asLineOn(&sb, item, columnSeparator)
		sb.WriteString(lineSeparator)
	}
	return sb.String()
}

func (gt *GenericTable[T]) SelectAll() {

	for r := 0; r < len(gt.data); r++ {
		for c := 0; c < len(gt.columns); c++ {
			gt.table.Select(widget.TableCellID{Row: r, Col: c})
		}
	}
	gt.table.Refresh()
}

func (gt *GenericTable[T]) SortBy(col int) {

	if col < 0 || col >= len(gt.columns) {
		return
	}

	column := gt.columns[col]
	if column.field.LessThan == nil {
		return
	}

	if gt.sortCol == col {
		gt.sortAsc = !gt.sortAsc
	} else {
		gt.sortCol = col
		gt.sortAsc = true
	}

	sort.Slice(gt.data, func(i, j int) bool {
		if gt.sortAsc {
			return column.field.LessThan(gt.data[i], gt.data[j])
		}
		return column.field.LessThan(gt.data[j], gt.data[i])
	})

	gt.table.Refresh()
}
