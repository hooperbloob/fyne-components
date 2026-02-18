package domains

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/hooperbloob/fyne-components/meta"
	"github.com/hooperbloob/fyne-components/table"

	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

type File struct {
	Name   string
	Time   time.Time
	Folder string
	Size   int64
}

var fileStatusField = meta.NewFieldDescriptor("?", func(f File) string { return f.Name }, nil, nil)
var fileNameField = meta.NewFieldDescriptor("Name", func(f File) string { return f.Name }, nil, nil)
var fileTimeField = meta.NewFieldDescriptor("Time", func(f File) string { return f.Time.Format("2006 01 02150405") }, nil, nil)
var fileSizeField = meta.NewFieldDescriptor("Size", func(f File) string { return fmt.Sprintf("%d", f.Size) }, nil, nil)

var fileColumns = []table.Column[File]{
	table.NewColumn(130, fileSizeField, fyne.TextAlignTrailing, nil),
	table.NewColumn(130, fileTimeField, fyne.TextAlignLeading, nil),
	table.NewColumn(300, fileNameField, fyne.TextAlignLeading, nil),
}

func expandPath(path string) (string, error) {

	if len(path) > 0 && path[0] == '~' {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		return filepath.Join(home, path[1:]), nil
	}
	return path, nil
}

func FilesFrom(folder string) []File {

	path, err := expandPath(folder)
	if err != nil {
		log.Fatalf("?")
	}

	dir, err := os.Open(path)
	if err != nil {
		log.Fatalf("Failed to open directory: %v", err)
	}
	defer dir.Close()

	// Read the directory contents
	localFiles, err := dir.Readdir(-1)

	var files []File
	for _, entry := range localFiles {
		files = append(files, File{
			Name: entry.Name(),
			Size: entry.Size(),
			Time: entry.ModTime(),
		})
	}

	return files
}

func SetupFileTable(window fyne.Window, folder string) *table.TableContainer[File] {

	newFileFunc := func() File { // Function to create a new empty Person
		return File{}
	}

	gTable := table.NewGenericTable(fileColumns, newFileFunc)

	gTable.SetData(FilesFrom(folder))

	editFileFunc := func(file *File, isAdd bool, idx int, callback func(File)) {
		nameEntry := widget.NewEntry()
		nameEntry.SetText(file.Name)

		folderEntry := widget.NewEntry()
		folderEntry.SetText(file.Folder)

		formItems := []*widget.FormItem{
			{Text: "Name", Widget: nameEntry},
			{Text: "Folder", Widget: folderEntry},
		}

		var title = "Edit File"
		if isAdd {
			title = "Add File"
		}

		dialog.ShowForm(title, "Save", "Cancel", formItems, func(confirmed bool) {
			if confirmed {

				edited := File{
					Name:   nameEntry.Text,
					Folder: folderEntry.Text,
				}
				callback(edited)
			}
		}, window)
	}
	// customFunctions := []table.ItemAction[Person]{
	// 	{
	// 		Label:   "E",
	// 		Icon:    theme.MailSendIcon(),
	// 		Action:  func(p []*Person) { println("Email clicked for: " + p[0].Name) },
	// 		Enabler: func(p []*Person) bool { return len(p[0].Email) > 0 }, // TODO loop through & check all
	// 	},
	// }

	return table.NewTableContainer(gTable, window, editFileFunc, nil) // Create the container with controls
}
