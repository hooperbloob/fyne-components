package domains

import (
	"fmt"
	"fyne-components/meta"
	"fyne-components/table"
	"image/color"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type Person struct {
	Name  string
	Email string
	Age   int
}

var people = []Person{
	{Name: "Alice Smith", Email: "alice@peanuts.com", Age: 30},
	{Name: "Bob Johnson", Email: "", Age: 25},
	{Name: "Carol Williams", Email: "carol@doughnuts.com", Age: 35},
	{Name: "June Smith", Email: "jsmith@peanuts.com", Age: 12},
	{Name: "Rob Johnson", Email: "", Age: 85},
	{Name: "Mitch Sommerset", Email: "mitch@cranky.com", Age: 39},
}

var statusField = meta.NewFieldDescriptor("?", func(p Person) string {
	if p.Email == "" {
		return "n"
	} else {
		return "y "
	}
}, nil, nil)
var personNameField = meta.NewFieldDescriptor("Name", func(p Person) string { return p.Name }, nil, nil)
var personEmailField = meta.NewFieldDescriptor("EMail", func(p Person) string { return p.Email }, nil, nil)
var personAgeField = meta.NewFieldDescriptor("Age", func(p Person) string { return fmt.Sprintf("%d", p.Age) }, nil, func(a, b Person) bool { return a.Age < b.Age })

var colorSetter = func(person Person) color.Color {

	if person.Email == "" {
		return color.NRGBA{R: 240, G: 80, B: 0, A: 255}
	} else {
		return color.Transparent
	}
}

var personColumns = []table.Column[Person]{
	table.NewColumn(40, statusField, fyne.TextAlignTrailing, colorSetter),
	table.NewColumn(40, personAgeField, fyne.TextAlignLeading, nil),
	table.NewColumn(120, personNameField, fyne.TextAlignLeading, nil),
	table.NewColumn(190, personEmailField, fyne.TextAlignLeading, nil),
}

var ageValidator = func(s string) error {
	if s == "" {
		return fmt.Errorf("Age is required")
	}
	age, err := strconv.Atoi(s)
	if err != nil {
		return fmt.Errorf("Age must be a number")
	}
	if age < 0 || age > 150 {
		return fmt.Errorf("Age must be between 0 and 150")
	}
	return nil
}

func SetupPeopleTable(window fyne.Window) *table.TableContainer[Person] {

	newPersonFunc := func() Person { // Function to create a new empty Person
		return Person{}
	}

	gTable := table.NewGenericTable(personColumns, newPersonFunc)

	gTable.SetData(people)

	editPersonFunc := func(person Person, isAdd bool, idx int, callback func(Person)) {
		nameEntry := widget.NewEntry()
		nameEntry.SetText(person.Name)

		emailEntry := widget.NewEntry()
		emailEntry.SetText(person.Email)

		ageEntry := widget.NewEntry()
		ageEntry.SetText(fmt.Sprintf("%d", person.Age))
		ageEntry.Validator = ageValidator

		formItems := []*widget.FormItem{
			{Text: "Name", Widget: nameEntry},
			{Text: "Email", Widget: emailEntry},
			{Text: "Age", Widget: ageEntry},
		}

		var title = "Edit Person"
		if isAdd {
			title = "Add Person"
		}

		dialog.ShowForm(title, "Save", "Cancel", formItems, func(confirmed bool) {
			if confirmed {
				var age int
				fmt.Sscanf(ageEntry.Text, "%d", &age)

				edited := Person{
					Name:  nameEntry.Text,
					Email: emailEntry.Text,
					Age:   age,
				}
				callback(edited)
			}
		}, window)
	}
	customFunctions := []table.ItemAction[Person]{
		{
			Label:   "E",
			Icon:    theme.MailSendIcon(),
			Action:  func(p Person) { println("Email clicked for: " + p.Name) },
			Enabler: func(p []Person) bool { return len(p[0].Email) > 0 }, // TODO loop through & check all
		},
	}

	return table.NewTableContainer(gTable, window, editPersonFunc, customFunctions) // Create the container with controls
}
