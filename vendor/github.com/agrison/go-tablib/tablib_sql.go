package tablib

import (
	"bytes"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	typePostgres = "postgres"
	typeMySQL    = "mysql"
	defaults     = map[string]string{"various." + typePostgres: "TEXT",
		"various." + typeMySQL: "VARCHAR(100)", "numeric." + typePostgres: "NUMERIC",
		"numeric." + typeMySQL: "DOUBLE"}
)

// columnSQLType determines the type of a column
// if throughout the whole column values have the same type then this type is
// returned, otherwise the VARCHAR/TEXT type is returned.
// numeric types are coerced into DOUBLE/NUMERIC
func (d *Dataset) columnSQLType(header, dbType string) (string, []interface{}) {
	types := 0
	currentType := ""
	maxString := 0
	values := d.Column(header)
	for _, c := range values {
		switch c.(type) {
		case uint, uint8, uint16, uint32, uint64,
			int, int8, int16, int32, int64,
			float32, float64:
			if currentType != "numeric" {
				currentType = "numeric"
				types++
			}
		case time.Time:
			if currentType != "time" {
				currentType = "time"
				types++
			}
		case string:
			if currentType != "string" {
				currentType = "string"
				types++
			}
			if len(c.(string)) > maxString {
				maxString = len(c.(string))
			}
		}
	}

	if types > 1 {
		return defaults["various."+dbType], values
	}
	switch currentType {
	case "numeric":
		return defaults["numeric."+dbType], values
	case "time":
		return "TIMESTAMP", values
	default:
		if dbType == typePostgres {
			return "TEXT", values
		}
		return "VARCHAR(" + strconv.Itoa(maxString) + ")", values
	}
}

// isStringColumn returns whether a column is VARCHAR/TEXT
func isStringColumn(c string) bool {
	return strings.HasPrefix(c, "VARCHAR") || strings.HasPrefix(c, "TEXT")
}

// MySQL returns a string representing a suite of MySQL commands
// recreating the Dataset into a table.
func (d *Dataset) MySQL(table string) *Exportable {
	return d.sql(table, typeMySQL)
}

// Postgres returns a string representing a suite of Postgres commands
// recreating the Dataset into a table.
func (d *Dataset) Postgres(table string) *Exportable {
	return d.sql(table, typePostgres)
}

// sql returns a string representing a suite of SQL commands
// recreating the Dataset into a table.
func (d *Dataset) sql(table, dbType string) *Exportable {
	b := newBuffer()

	tableSQL, columnTypes, columnValues := d.createTable(table, dbType)
	b.WriteString(tableSQL)

	reg, _ := regexp.Compile("[']")
	// inserts
	for i := range d.data {
		b.WriteString("INSERT INTO " + table + " VALUES(" + strconv.Itoa(i+1) + ", ")
		for j, col := range d.headers {
			asStr := d.asString(columnValues[col][i])
			if isStringColumn(columnTypes[col]) {
				b.WriteString("'" + reg.ReplaceAllString(asStr, "''") + "'")
			} else if strings.HasPrefix(columnTypes[col], "TIMESTAMP") {
				if dbType == typeMySQL {
					b.WriteString("CONVERT_TZ('" + asStr[:10] + " " + asStr[11:19] + "', '" + asStr[len(asStr)-6:] + "', 'SYSTEM')")
				} else {
					b.WriteString("'" + asStr + "'") // simpler with Postgres
				}
			} else {
				b.WriteString(asStr)
			}
			if j < len(d.headers)-1 {
				b.WriteString(", ")
			}
		}
		b.WriteString(");\n")
	}
	b.WriteString("\nCOMMIT;\n")

	return newExportable(b)
}

func (d *Dataset) createTable(table, dbType string) (string, map[string]string, map[string][]interface{}) {
	var b bytes.Buffer
	columnValues := make(map[string][]interface{})
	columnTypes := make(map[string]string)

	// create table
	b.WriteString("CREATE TABLE IF NOT EXISTS " + table)
	if dbType == typePostgres {
		b.WriteString("\n(\n\tid SERIAL PRIMARY KEY,\n")
	} else {
		b.WriteString("\n(\n\tid INT NOT NULL AUTO_INCREMENT PRIMARY KEY,\n")
	}
	for i, h := range d.headers {
		b.WriteString("\t" + h)
		t, v := d.columnSQLType(h, dbType)
		columnValues[h] = v
		columnTypes[h] = t
		b.WriteString(" " + t)
		if i < len(d.headers)-1 {
			b.WriteString(",")
		}
		b.WriteString("\n")
	}

	b.WriteString(");\n\n")

	return b.String(), columnTypes, columnValues
}
