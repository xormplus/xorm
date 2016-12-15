package tablib

// Sheet represents a sheet in a Databook, holding a title (if any) and a dataset.
type Sheet struct {
	title   string
	dataset *Dataset
}

// Title return the title of the sheet.
func (s Sheet) Title() string {
	return s.title
}

// Dataset returns the dataset of the sheet.
func (s Sheet) Dataset() *Dataset {
	return s.dataset
}

// Databook represents a Databook which is an array of sheets.
type Databook struct {
	sheets map[string]Sheet
}

// NewDatabook constructs a new Databook.
func NewDatabook() *Databook {
	return &Databook{make(map[string]Sheet)}
}

// Sheets returns the sheets in the Databook.
func (d *Databook) Sheets() map[string]Sheet {
	return d.sheets
}

// Sheet returns the sheet with a specific title.
func (d *Databook) Sheet(title string) Sheet {
	return d.sheets[title]
}

// AddSheet adds a sheet to the Databook.
func (d *Databook) AddSheet(title string, dataset *Dataset) {
	d.sheets[title] = Sheet{title, dataset}
}

// Size returns the number of sheets in the Databook.
func (d *Databook) Size() int {
	return len(d.sheets)
}

// Wipe removes all Dataset objects from the Databook.
func (d *Databook) Wipe() {
	for k := range d.sheets {
		delete(d.sheets, k)
	}
}
