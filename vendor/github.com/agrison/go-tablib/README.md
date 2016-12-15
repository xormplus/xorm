# go-tablib: format-agnostic tabular dataset library

[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat-square)][license]
[![Go Documentation](http://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)][godocs]
[![Go Report Card](https://goreportcard.com/badge/github.com/agrison/go-tablib)][goreportcard]
[![Build Status](https://travis-ci.org/agrison/go-tablib.svg?branch=master)](https://travis-ci.org/agrison/go-tablib)

[license]: https://github.com/agrison/go-tablib/blob/master/LICENSE
[godocs]: https://godoc.org/github.com/agrison/go-tablib
[goreportcard]: https://goreportcard.com/report/github.com/agrison/go-tablib

Go-Tablib is a format-agnostic tabular dataset library, written in Go.
This is a port of the famous Python's [tablib](https://github.com/kennethreitz/tablib) by Kenneth Reitz with some new features.

Export formats supported:

* JSON (Sets + Books)
* YAML (Sets + Books)
* XLSX (Sets + Books)
* XML (Sets + Books)
* TSV (Sets)
* CSV (Sets)
* ASCII + Markdown (Sets)
* MySQL (Sets)
* Postgres (Sets)

Loading formats supported:

* JSON (Sets + Books)
* YAML (Sets + Books)
* XML (Sets)
* CSV (Sets)
* TSV (Sets)


## Overview

### tablib.Dataset
A Dataset is a table of tabular data. It must have a header row. Datasets can be exported to JSON, YAML, CSV, TSV, and XML. They can be filtered, sorted and validated against constraint on columns.

### tablib.Databook
A Databook is a set of Datasets. The most common form of a Databook is an Excel file with multiple spreadsheets. Databooks can be exported to JSON, YAML and XML.

### tablib.Exportable
An exportable is a struct that holds a buffer representing the Databook or Dataset after it has been formated to any of the supported export formats.
At this point the Datbook or Dataset cannot be modified anymore, but it can be returned as a `string`, a `[]byte` or written to a `io.Writer` or a file.

## Usage

Creates a dataset and populate it:

```go
ds := NewDataset([]string{"firstName", "lastName"})
```

Add new rows:
```go
ds.Append([]interface{}{"John", "Adams"})
ds.AppendValues("George", "Washington")
```

Add new columns:
```go
ds.AppendColumn("age", []interface{}{90, 67})
ds.AppendColumnValues("sex", "male", "male")
```

Add a dynamic column, by passing a function which has access to the current row, and must
return a value:
```go
func lastNameLen(row []interface{}) interface{} {
	return len(row[1].(string))
}
ds.AppendDynamicColumn("lastName length", lastNameLen)
ds.CSV()
// >>
// firstName, lastName, age, sex, lastName length
// John, Adams, 90, male, 5
// George, Washington, 67, male, 10
```

Delete rows:
```go
ds.DeleteRow(1) // starts at 0
```

Delete columns:
```go
ds.DeleteColumn("sex")
```

Get a row or multiple rows:
```go
row, _ := ds.Row(0)
fmt.Println(row["firstName"]) // George

rows, _ := ds.Rows(0, 1)
fmt.Println(rows[0]["firstName"]) // George
fmt.Println(rows[1]["firstName"]) // Thomas
```

Slice a Dataset:
```go
newDs, _ := ds.Slice(1, 5) // returns a fresh Dataset with rows [1..5[
```


## Filtering

You can add **tags** to rows by using a specific `Dataset` method. This allows you to filter your `Dataset` later. This can be useful to separate rows of data based on arbitrary criteria (e.g. origin) that you don’t want to include in your `Dataset`.
```go
ds := NewDataset([]string{"Maker", "Model"})
ds.AppendTagged([]interface{}{"Porsche", "911"}, "fast", "luxury")
ds.AppendTagged([]interface{}{"Skoda", "Octavia"}, "family")
ds.AppendTagged([]interface{}{"Ferrari", "458"}, "fast", "luxury")
ds.AppendValues("Citroen", "Picasso")
ds.AppendValues("Bentley", "Continental")
ds.Tag(4, "luxury") // Bentley
ds.AppendValuesTagged("Aston Martin", "DB9", /* these are tags */ "fast", "luxury")
```

Filtering the `Dataset` is possible by calling `Filter(column)`:
```go
luxuryCars, err := ds.Filter("luxury").CSV()
fmt.Println(luxuryCars)
// >>>
// Maker,Model
// Porsche,911
// Ferrari,458
// Bentley,Continental
// Aston Martin,DB9
```

```go
fastCars, err := ds.Filter("fast").CSV()
fmt.Println(fastCars)
// >>>
// Maker,Model
// Porsche,911
// Ferrari,458
// Aston Martin,DB9
```

Tags at a specific row can be retrieved by calling `Dataset.Tags(index int)`

## Sorting

Datasets can be sorted by a specific column.
```go
ds := NewDataset([]string{"Maker", "Model", "Year"})
ds.AppendValues("Porsche", "991", 2012)
ds.AppendValues("Skoda", "Octavia", 2011)
ds.AppendValues("Ferrari", "458", 2009)
ds.AppendValues("Citroen", "Picasso II", 2013)
ds.AppendValues("Bentley", "Continental GT", 2003)

sorted, err := ds.Sort("Year").CSV()
fmt.Println(sorted)
// >>
// Maker, Model, Year
// Bentley, Continental GT, 2003
// Ferrari, 458, 2009
// Skoda, Octavia, 2011
// Porsche, 991, 2012
// Citroen, Picasso II, 2013
```

## Constraining

Datasets can have columns constrained by functions and further checked if valid.
```go
ds := NewDataset([]string{"Maker", "Model", "Year"})
ds.AppendValues("Porsche", "991", 2012)
ds.AppendValues("Skoda", "Octavia", 2011)
ds.AppendValues("Ferrari", "458", 2009)
ds.AppendValues("Citroen", "Picasso II", 2013)
ds.AppendValues("Bentley", "Continental GT", 2003)

ds.ConstrainColumn("Year", func(val interface{}) bool { return val.(int) > 2008 })
ds.ValidFailFast() // false
if !ds.Valid() { // validate the whole dataset, errors are retrieved in Dataset.ValidationErrors
	ds.ValidationErrors[0] // Row: 4, Column: 2
}
```

A Dataset with constrained columns can be filtered to keep only the rows satisfying the constraints.
```go
valid := ds.ValidSubset().Tabular("simple") // Cars after 2008
fmt.Println(valid)
```

Will output:
```
------------  ---------------  ---------
      Maker            Model       Year
------------  ---------------  ---------
    Porsche              991       2012

      Skoda          Octavia       2011

    Ferrari              458       2009

    Citroen       Picasso II       2013
------------  ---------------  ---------
```

```go
invalid := ds.InvalidSubset().Tabular("simple") // Cars before 2008
fmt.Println(invalid)
```

Will output:
```
------------  -------------------  ---------
      Maker                Model       Year
------------  -------------------  ---------
    Bentley       Continental GT       2003
------------  -------------------  ---------
```

## Loading

### JSON
```go
ds, _ := LoadJSON([]byte(`[
  {"age":90,"firstName":"John","lastName":"Adams"},
  {"age":67,"firstName":"George","lastName":"Washington"},
  {"age":83,"firstName":"Henry","lastName":"Ford"}
]`))
```

### YAML
```go
ds, _ := LoadYAML([]byte(`- age: 90
  firstName: John
  lastName: Adams
- age: 67
  firstName: George
  lastName: Washington
- age: 83
  firstName: Henry
  lastName: Ford`))
```

## Exports

### Exportable

Any of the following export format returns an `*Exportable` which means you can use:
- `Bytes()` to get the content as a byte array
- `String()` to get the content as a string
- `WriteTo(io.Writer)` to write the content to an `io.Writer`
- `WriteFile(filename string, perm os.FileMode)` to write to a file

It avoids unnecessary conversion between `string` and `[]byte` to output/write/whatever.
Thanks to [@figlief](https://github.com/figlief) for the proposition. 

### JSON
```go
json, _ := ds.JSON()
fmt.Println(json)
```

Will output:
```json
[{"age":90,"firstName":"John","lastName":"Adams"},{"age":67,"firstName":"George","lastName":"Washington"},{"age":83,"firstName":"Henry","lastName":"Ford"}]
```

### XML
```go
xml, _ := ds.XML()
fmt.Println(xml)
```

Will ouput:
```xml
<dataset>
 <row>
   <age>90</age>
   <firstName>John</firstName>
   <lastName>Adams</lastName>
 </row>  <row>
   <age>67</age>
   <firstName>George</firstName>
   <lastName>Washington</lastName>
 </row>  <row>
   <age>83</age>
   <firstName>Henry</firstName>
   <lastName>Ford</lastName>
 </row>
</dataset>
```

### CSV
```go
csv, _ := ds.CSV()
fmt.Println(csv)
```

Will ouput:
```csv
firstName,lastName,age
John,Adams,90
George,Washington,67
Henry,Ford,83
```

### TSV
```go
tsv, _ := ds.TSV()
fmt.Println(tsv)
```

Will ouput:
```tsv
firstName lastName  age
John  Adams  90
George  Washington  67
Henry Ford 83
```

### YAML
```go
yaml, _ := ds.YAML()
fmt.Println(yaml)
```

Will ouput:
```yaml
- age: 90
  firstName: John
  lastName: Adams
- age: 67
  firstName: George
  lastName: Washington
- age: 83
  firstName: Henry
  lastName: Ford
```

### HTML
```go
html := ds.HTML()
fmt.Println(html)
```

Will output:
```html
<table class="table table-striped">
	<thead>
		<tr>
			<th>firstName</th>
			<th>lastName</th>
			<th>age</th>
		</tr>
	</thead>
	<tbody>
		<tr>
			<td>George</td>
			<td>Washington</td>
			<td>90</td>
		</tr>
		<tr>
			<td>Henry</td>
			<td>Ford</td>
			<td>67</td>
		</tr>
		<tr>
			<td>Foo</td>
			<td>Bar</td>
			<td>83</td>
		</tr>
	</tbody>
</table>
```

### XLSX
```go
xlsx, _ := ds.XLSX()
fmt.Println(xlsx)
// >>>
// binary content
xlsx.WriteTo(...)
```

### ASCII

#### Grid format
```go
ascii := ds.Tabular("grid" /* tablib.TabularGrid */)
fmt.Println(ascii)
```

Will output:
```
+--------------+---------------+--------+
|    firstName |      lastName |    age |
+==============+===============+========+
|       George |    Washington |     90 |
+--------------+---------------+--------+
|        Henry |          Ford |     67 |
+--------------+---------------+--------+
|          Foo |           Bar |     83 |
+--------------+---------------+--------+
```

#### Simple format
```go
ascii := ds.Tabular("simple" /* tablib.TabularSimple */)
fmt.Println(ascii)
```

Will output:
```
--------------  ---------------  --------
    firstName         lastName       age
--------------  ---------------  --------
       George       Washington        90

        Henry             Ford        67

          Foo              Bar        83
--------------  ---------------  --------
```

#### Condensed format
```go
ascii := ds.Tabular("condensed" /* tablib.TabularCondensed */)
fmt.Println(ascii)
```

Similar to simple but with less line feed:
```
--------------  ---------------  --------
    firstName         lastName       age
--------------  ---------------  --------
       George       Washington        90
        Henry             Ford        67
          Foo              Bar        83
--------------  ---------------  --------
```

### Markdown

Markdown tables are similar to the Tabular condensed format, except that they have
pipe characters separating columns.

```go
mkd := ds.Markdown() // or
mkd := ds.Tabular("markdown" /* tablib.TabularMarkdown */)
fmt.Println(mkd)
```

Will output:
```
|     firstName   |       lastName    |    gpa  |
| --------------  | ---------------   | ------- |
|          John   |          Adams    |     90  |
|        George   |     Washington    |     67  |
|        Thomas   |      Jefferson    |     50  |
```

Which equals to the following when rendered as HTML:

|     firstName   |       lastName    |    gpa  |
| --------------  | ---------------   | ------- |
|          John   |          Adams    |     90  |
|        George   |     Washington    |     67  |
|        Thomas   |      Jefferson    |     50  |

### MySQL
```go
sql := ds.MySQL()
fmt.Println(sql)
```

Will output:
```sql
CREATE TABLE IF NOT EXISTS presidents
(
	id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
	firstName VARCHAR(9),
	lastName VARCHAR(8),
	gpa DOUBLE
);

INSERT INTO presidents VALUES(1, 'Jacques', 'Chirac', 88);
INSERT INTO presidents VALUES(2, 'Nicolas', 'Sarkozy', 98);
INSERT INTO presidents VALUES(3, 'François', 'Hollande', 34);

COMMIT;
```

Numeric (`uint`, `int`, `float`, ...) are stored as `DOUBLE`, `string`s as `VARCHAR` with width set to the length of the longest string in the column, and `time.Time`s are stored as `TIMESTAMP`.

### Postgres
```go
sql := ds.Postgres()
fmt.Println(sql)
```

Will output:
```sql
CREATE TABLE IF NOT EXISTS presidents
(
	id SERIAL PRIMARY KEY,
	firstName TEXT,
	lastName TEXT,
	gpa NUMERIC
);

INSERT INTO presidents VALUES(1, 'Jacques', 'Chirac', 88);
INSERT INTO presidents VALUES(2, 'Nicolas', 'Sarkozy', 98);
INSERT INTO presidents VALUES(3, 'François', 'Hollande', 34);

COMMIT;
```

Numerics (`uint`, `int`, `float`, ...) are stored as `NUMERIC`, `string`s as `TEXT` and `time.Time`s are stored as `TIMESTAMP`.

## Databooks

This is an example of how to use Databooks.

```go
db := NewDatabook()
// or loading a JSON content
db, err := LoadDatabookJSON([]byte(`...`))
// or a YAML content
db, err := LoadDatabookYAML([]byte(`...`))

// a dataset of presidents
presidents, _ := LoadJSON([]byte(`[
  {"Age":90,"First name":"John","Last name":"Adams"},
  {"Age":67,"First name":"George","Last name":"Washington"},
  {"Age":83,"First name":"Henry","Last name":"Ford"}
]`))

// a dataset of cars
cars := NewDataset([]string{"Maker", "Model", "Year"})
cars.AppendValues("Porsche", "991", 2012)
cars.AppendValues("Skoda", "Octavia", 2011)
cars.AppendValues("Ferrari", "458", 2009)
cars.AppendValues("Citroen", "Picasso II", 2013)
cars.AppendValues("Bentley", "Continental GT", 2003)

// add the sheets to the Databook
db.AddSheet("Cars", cars.Sort("Year"))
db.AddSheet("Presidents", presidents.SortReverse("Age"))

fmt.Println(db.JSON())
```

Will output the following JSON representation of the Databook:
```json
[
  {
    "title": "Cars",
    "data": [
      {"Maker":"Bentley","Model":"Continental GT","Year":2003},
      {"Maker":"Ferrari","Model":"458","Year":2009},
      {"Maker":"Skoda","Model":"Octavia","Year":2011},
      {"Maker":"Porsche","Model":"991","Year":2012},
      {"Maker":"Citroen","Model":"Picasso II","Year":2013}
    ]
  },
  {
    "title": "Presidents",
    "data": [
      {"Age":90,"First name":"John","Last name":"Adams"},
      {"Age":83,"First name":"Henry","Last name":"Ford"},
      {"Age":67,"First name":"George","Last name":"Washington"}
    ]
  }
]
```

## Installation

```bash
go get github.com/agrison/go-tablib
```

For those wanting the v1 version where export methods returned a `string` and not an `Exportable`:
```bash
go get gopkg.in/agrison/go-tablib.v1
```

## TODO

* Loading in more formats
* Support more formats: DBF, XLS, LATEX, ...

## Contribute

It is a work in progress, so it may exist some bugs and edge cases not covered by the test suite.

But we're on Github and this is Open Source, pull requests are more than welcomed, come and have some fun :)

## Acknowledgement

Thanks to kennethreitz for the first implementation in Python, [`github.com/bndr/gotabulate`](https://github.com/bndr/gotabulate), [`github.com/clbanning/mxj`](https://github.com/clbanning/mxj), [`github.com/tealeg/xlsx`](https://github.com/tealeg/xlsx), [`gopkg.in/yaml.v2`](https://gopkg.in/yaml.v2)
