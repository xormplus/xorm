## History

### 2016-02-26
- Added support for Markdown tables export

### 2016-02-25

- Constrained columns
  - `Dataset.ValidSubset()`
  - `Dataset.InvalidSubset()`
- Tagging a specific row after it was already created
- Loading Databooks
  - JSON
  - YAML
- Loading Datasets
  - CSV
  - TSV
  - XML
- Unit test coverage

### 2016-02-24

- Constrained columns
- Support for `time.Time` in `Dataset.MySQL()` and `Dataset.Postgres()` export.
- Source files refactoring
- Added on travis-ci
- Retrieving specific rows
  - `Dataset.Row(int)`
  - `Dataset.Rows(int...)`
  - `Dataset.Slice(int, int)`

### 2016-02-23

- First release with support for:
  - Loading YAML, JSON
  - Exporting YAML, JSON, CSV, TSV, XLS, XML, ASCII
  - Filtering + Tagging
  - Sorting
  - ...
