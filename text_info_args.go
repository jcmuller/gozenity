package gozenity

import (
	"errors"
	"os"
)

// TextInfoArgs represents options for text info
type TextInfoArgs struct {
	Checkbox string
	Editable bool
	Filename string
	Text     string
	URL      string
}

func (tia *TextInfoArgs) tooManyArguments() bool {
	return (tia.Filename != "" && tia.Text != "") ||
		(tia.Filename != "" && tia.URL != "") ||
		(tia.Text != "" && tia.URL != "")
}

func (tia *TextInfoArgs) notEnoughArguments() bool {
	return tia.Filename == "" && tia.Text == "" && tia.URL == ""
}

// Parse returns a slice of strings usable for a text info
func (tia *TextInfoArgs) Parse() (args []string, err error) {
	if tia.Checkbox != "" {
		args = append(args, `--checkbox`, tia.Checkbox)
	}

	if tia.Editable {
		args = append(args, `--editable`)
	}

	if tia.tooManyArguments() {
		return []string{}, errors.New("Only one of Filename, Text and URL can be supplied")
	}

	if tia.notEnoughArguments() {
		return []string{}, errors.New("One of Filename, Text or URL need to be supplied")
	}

	if tia.Filename != "" {
		if _, err = os.Stat(tia.Filename); err != nil {
			return
		}

		args = append(args, `--filename`, tia.Filename)
	}

	if tia.URL != "" {
		args = append(args, `--html`, `--url`, tia.URL)
	}

	return
}
