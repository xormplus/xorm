package tablib

import "errors"

var (
	// ErrInvalidDimensions is returned when trying to append/insert too much
	// or not enough values to a row or column
	ErrInvalidDimensions = errors.New("tablib: Invalid dimension")
	// ErrInvalidColumnIndex is returned when trying to insert a column at an
	// invalid index
	ErrInvalidColumnIndex = errors.New("tablib: Invalid column index")
	// ErrInvalidRowIndex is returned when trying to insert a row at an
	// invalid index
	ErrInvalidRowIndex = errors.New("tablib: Invalid row index")
	// ErrInvalidDataset is returned when trying to validate a Dataset against
	// the constraints that have been set on its columns.
	ErrInvalidDataset = errors.New("tablib: Invalid dataset")
	// ErrInvalidTag is returned when trying to add a tag which is not a string.
	ErrInvalidTag = errors.New("tablib: A tag must be a string")
)
