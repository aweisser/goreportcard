package report

import (
	"strconv"
	"strings"
)

// Check describes what methods various checks (gofmt, go lint, etc.)
// should implement
type Check interface {
	Name() string
	Description() string
	Weight() float64
	// Percentage returns the passing percentage of the check,
	// as well as a map of filename to output
	Percentage() (float64, []FileSummary, error)
}

// Error contains the line number and the reason for
// an error output from a command
type Error struct {
	LineNumber  int    `json:"line_number"`
	ErrorString string `json:"error_string"`
}

// FileSummary contains the filename, location of the file
// on GitHub, and all of the errors related to the file
type FileSummary struct {
	Filename string  `json:"filename"`
	FileURL  string  `json:"file_url"`
	Errors   []Error `json:"errors"`
}

// AddError adds an Error to FileSummary
func (fs *FileSummary) AddError(out string) error {
	s := strings.SplitN(out, ":", 2)
	msg := strings.SplitAfterN(s[1], ":", 3)[2]

	e := Error{ErrorString: msg}
	ls := strings.Split(s[1], ":")
	ln, err := strconv.Atoi(ls[0])
	if err != nil {
		return err
	}
	e.LineNumber = ln

	fs.Errors = append(fs.Errors, e)

	return nil
}
