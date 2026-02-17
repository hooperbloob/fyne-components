package table

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type ItemAction[T any] struct {
	Label   string
	Icon    fyne.Resource
	Action  func(T)
	Enabler func([]T) bool
}

// TableContainer wraps the GenericTable with controls
type TableContainer[T any] struct {
	widget.BaseWidget
	table          *GenericTable[T]
	addButton      *widget.Button
	editButton     *widget.Button
	deleteButton   *widget.Button
	deleteAction   ItemAction[T]
	customActions  []ItemAction[T]
	customControls []*widget.Button
	container      *fyne.Container
	window         fyne.Window
	editItemFunc   func(T, bool, int, func(T)) // Function to show add/edit dialog
}

// NewTableContainer creates a container with table and controls
func NewTableContainer[T any](
	table *GenericTable[T],
	window fyne.Window,
	editItemFunc func(T, bool, int, func(T)),
	actions []ItemAction[T],
) *TableContainer[T] {
	tc := &TableContainer[T]{
		table:        table,
		window:       window,
		editItemFunc: editItemFunc,
	}

	tc.addButton = widget.NewButtonWithIcon("", theme.ContentAddIcon(), tc.handleAdd)
	tc.editButton = widget.NewButtonWithIcon("", theme.SettingsIcon(), tc.handleEdit)
	tc.deleteButton = widget.NewButtonWithIcon("", theme.DeleteIcon(), tc.handleDelete)
	//tc.deleteButton.Importance = widget.DangerImportance

	tc.customActions = actions
	tc.customControls = make([]*widget.Button, len(actions))

	for idx, act := range actions {
		tc.customControls[idx] = tc.buttonFor(act, idx) // TODO allow for other control types
		tc.customActions[idx] = act
	}

	tc.editButton.Disable()
	tc.deleteButton.Disable()

	// Update delete button state when selection changes
	table.table.OnSelected = func(id widget.TableCellID) {
		table.newSelectedRow(id.Row)
		tc.updateEditButtons()
	}

	table.table.OnUnselected = func(id widget.TableCellID) {
		table.selectedRows.Remove(id.Row)
		tc.updateEditButtons()
	}

	controls := tc.createControls()
	controlContainer := container.NewVBox(controls...)

	tc.container = container.NewBorder(
		nil,
		nil,
		nil,
		controlContainer,
		table,
	)

	tc.ExtendBaseWidget(tc)
	return tc
}

func (tc *TableContainer[T]) buttonFor(act ItemAction[T], idx int) *widget.Button {

	var button *widget.Button

	if act.Icon == nil {
		button = widget.NewButton(act.Label, func() { tc.handleCustom(idx) })
	} else {
		button = widget.NewButtonWithIcon("", act.Icon, func() { tc.handleCustom(idx) })
	}

	button.Disable()
	return button
}

func (tc *TableContainer[T]) AdjustColumns() {
	tc.table.SetColumnWidths()
}

func (tc *TableContainer[T]) updateEditButtons() {

	selections := tc.table.SelectedItems()
	if len(selections) > 0 {
		tc.editButton.Enable()
		tc.deleteButton.Enable()
	} else {
		tc.deleteButton.Disable()
		tc.editButton.Disable()
	}
	tc.enableCustom(valuesOf(selections))
}

func (tc *TableContainer[T]) createControls() []fyne.CanvasObject {

	cbCount := len(tc.customControls)
	controls := make([]fyne.CanvasObject, cbCount+5)
	controls[0] = tc.addButton
	controls[1] = tc.editButton
	controls[2] = widget.NewSeparator()
	for i, cb := range tc.customControls {
		controls[i+3] = cb
	}
	controls[cbCount+3] = layout.NewSpacer()
	controls[cbCount+4] = tc.deleteButton
	return controls
}

func valuesOf[T any](mapp map[int]T) []T {

	values := make([]T, len(mapp))
	i := 0
	for _, val := range mapp {
		values[i] = val
		i++
	}
	return values
}

func (tc *TableContainer[T]) enableCustom(values []T) {

	for idx, control := range tc.customControls {
		if len(values) > 0 {
			enabler := tc.customActions[idx].Enabler
			if enabler != nil {
				enabled := enabler(values)
				if enabled {
					control.Enable()
					continue
				}
			}
		} else {
			control.Disable()
		}
	}
}

// handleAdd shows dialog to add new item
func (tc *TableContainer[T]) handleAdd() {
	newItem := tc.table.newItemFunc()
	tc.editItemFunc(newItem, true, -1, func(edited T) {
		tc.table.AddItem(edited)
	})
}

func firstValueOf[T any](items map[int]T) (int, T) {
	for i, item := range items {
		return i, item
	}
	var zero T // won't get here, allows it to compile
	return -1, zero
}

func (tc *TableContainer[T]) handleEdit() {
	selectedItems := tc.table.SelectedItems()
	idx, item := firstValueOf(selectedItems)

	tc.editItemFunc(item, false, idx, func(edited T) {
		tc.table.ItemEdited(idx, edited)
	})
}

// handleDelete deletes selected items with confirmation
func (tc *TableContainer[T]) handleDelete() {
	count := tc.table.GetSelectedCount()
	if count == 0 {
		return
	}

	message := fmt.Sprintf("Delete %d selected item(s)?", count)
	dialog.ShowConfirm("Confirm Delete", message, func(confirmed bool) {
		if confirmed {
			tc.table.DeleteSelected()
			tc.updateEditButtons()
		}
	}, tc.window)
}

// handleAdd shows dialog to add new item
func (tc *TableContainer[T]) handleCustom(index int) {

	action := tc.customActions[index]
	selectedItems := tc.table.SelectedItems()

	action.Action(selectedItems[0])
}

// HandleKeyboard processes keyboard shortcuts
func (tc *TableContainer[T]) HandleKeyboard(event *fyne.KeyEvent) {
	// Simple Delete key for all platforms
	if event.Name == fyne.KeyDelete || event.Name == fyne.KeyBackspace {
		if tc.table.GetSelectedCount() > 0 {
			tc.handleDelete()
		}
	}
}

// TypedShortcut handles keyboard shortcuts with modifiers
func (tc *TableContainer[T]) TypedShortcut(shortcut fyne.Shortcut) {
	if _, ok := shortcut.(*desktop.CustomShortcut); ok {
		// Handle custom shortcuts if needed
		return
	}

	// Check for Ctrl+N or Cmd+N
	if _, ok := shortcut.(*fyne.ShortcutCopy); !ok {
		// Try to match our add shortcut
		if key, ok := shortcut.(fyne.KeyboardShortcut); ok {
			if key.Key() == fyne.KeyN {
				tc.handleAdd()
			}
		}
	}
}

func (tc *TableContainer[T]) CopySelectionToClipboard(app fyne.App) {

	content := tc.table.SelectionAsString("\t", "\n")
	app.Clipboard().SetContent(content)
}

func (tc *TableContainer[T]) SelectAll() {
	tc.table.SelectAll()
}

// CreateRenderer implements fyne.Widget
func (tc *TableContainer[T]) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(tc.container)
}

func (tc *TableContainer[T]) SetReadOnlyEnabler(enabler func([]T) bool) {

}

func (tc *TableContainer[T]) NewItemValidator(validator func([]T) error) {

}
