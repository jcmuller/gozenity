// Package gozenity is a simple wrapper for zenity) in Go.
package gozenity

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
	"time"
)

// Gozenity holds the structure of this thing
type Gozenity struct {
	command   string
	arguments []string
}

const (
	zenity = "zenity"
)

// New returns an instance of a Gozenity
func New(prompt string, arguments ...string) *Gozenity {
	titles := []string{`--title`, prompt, `--text`, prompt}
	arguments = append(arguments, titles...)

	program, err := exec.LookPath(zenity)

	if err != nil {
		log.Fatalf("%s not found", zenity)
	}

	return &Gozenity{program, arguments}
}

// EmptySelectionError is returned if there is no selection
type EmptySelectionError struct{}

func (e *EmptySelectionError) Error() string {
	return "Nothing selected"
}

// List pops up the menu
func List(prompt string, options ...string) (selection string, err error) {
	args := []string{`--list`, `--column`, prompt}
	args = append(args, options...)
	g := New(prompt, args...)
	selection, err = g.execute()
	return
}

// Entry asks for input
func Entry(prompt, placeholder string) (entry string, err error) {
	g := New(prompt, `--entry`, `--entry-text`, placeholder)
	entry, err = g.execute()

	return
}

// Calendar picks a date
func Calendar(prompt string, defaultDate time.Time) (date string, err error) {
	g := New(
		prompt,
		`--calendar`,
		fmt.Sprintf("--day=%d", defaultDate.Day()),
		fmt.Sprintf("--month=%d", defaultDate.Month()),
		fmt.Sprintf("--year=%d", defaultDate.Year()),
		"--date-format", `%m/%d/%Y`,
	)

	date, err = g.execute()

	return
}

// Error errors errors
func Error(prompt string) (err error) {
	g := New(prompt, `--error`, `--ellipsize`)
	_, err = g.execute()
	return
}

// Info informs information
func Info(prompt string) (err error) {
	g := New(prompt, `--info`, `--ellipsize`)

	_, err = g.execute()
	return
}

// FileSelection opens a file selector
func FileSelection(prompt string, filtersMap map[string][]string) (files []string, err error) {
	args := []string{`--file-selection`, `--multiple`}
	filters := buildFileFilter(filtersMap)
	args = append(args, filters...)

	g := New(prompt, args...)
	result, err := g.execute()
	files = strings.Split(result, `|`)
	return
}

// DirectorySelection opens a file selector
func DirectorySelection(prompt string) (files []string, err error) {
	g := New(prompt, `--file-selection`, `--multiple`, `--directory`)
	result, err := g.execute()
	files = strings.Split(result, `|`)
	return
}

// Notification notifies notifiees
func Notification(prompt string) (err error) {
	g := New(prompt, `--notification`, `--listen`)
	_, err = g.execute()
	return
}

// Question asks for answer
func Question(prompt string) (answer bool, err error) {
	g := New(prompt, `--question`)

	cmd := exec.Command(g.command, g.arguments...)

	if err := cmd.Start(); err != nil {
		fmt.Printf("cmd.Start: %v\n", err)
	}

	if err := cmd.Wait(); err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				if status.ExitStatus() == 1 {
					return false, nil
				}
			}
		}
		return false, err
	}

	return true, nil
}

// Warning warns about warnings
func Warning(prompt string) (err error) {
	g := New(prompt, `--warning`, `--ellipsize`)

	_, err = g.execute()
	return
}

// ScaleArgs are the options for Scale
type ScaleArgs struct {
	Initial int
	Step    int
	Min     int
	Max     int
}

// Scale shows a nice scale
func Scale(prompt string, args *ScaleArgs) (answer int, err error) {
	g := New(
		prompt,
		`--scale`,
		fmt.Sprintf("--value=%d", args.Initial),
		fmt.Sprintf("--min-value=%d", args.Min),
		fmt.Sprintf("--max-value=%d", args.Max),
		fmt.Sprintf("--step=%d", args.Step),
	)

	ans, err := g.execute()

	if ans == "" {
		return -1, err
	}

	answer, nerr := strconv.Atoi(ans)

	if nerr != nil {
		panic(nerr)
	}

	return
}

// TextInfoArgs represents options for text info
type TextInfoArgs struct {
	Checkbox string
	Editable bool
	Filename string
	Text     string
	URL      string
}

// Parse returns a slice of strings usable for a text info
func (tia *TextInfoArgs) Parse() (args []string, err error) {
	if tia.Checkbox != "" {
		args = append(args, `--checkbox`, tia.Checkbox)
	}

	if tia.Editable {
		args = append(args, `--editable`)
	}

	if (tia.Filename != "" && tia.Text != "") ||
		(tia.Filename != "" && tia.URL != "") ||
		(tia.Text != "" && tia.URL != "") {
		return []string{}, errors.New("Only one of Filename, Text and URL can be supplied")
	}

	if tia.Filename == "" && tia.Text == "" && tia.URL == "" {
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

// TextInfo shows a webview
func TextInfo(prompt string, args *TextInfoArgs) (text string, err error) {
	pArgs := []string{`--text-info`}
	parsedArgs, err := args.Parse()

	if err != nil {
		return
	}

	pArgs = append(pArgs, parsedArgs...)

	g := New(prompt, pArgs...)

	cmd := exec.Command(g.command, g.arguments...)

	if args.Text != "" {
		stdin, err := cmd.StdinPipe()

		if err != nil {
			log.Fatalf("Error getting pipe: %s", err)
		}

		go func(stdin io.WriteCloser) {
			defer stdin.Close()
			io.WriteString(stdin, args.Text)
		}(stdin)
	}

	byteOut, err := cmd.Output()

	if err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				if status.ExitStatus() == 1 {
					return "", nil
				}
			}
		}
		return "", err
	}

	text = string(byteOut)

	return
}

// ColorSelection pops up a color selection dialog
func ColorSelection(prompt, initial string, showPalette bool) (color string, err error) {
	args := []string{`--color-selection`, `--color`, initial}
	if showPalette {
		args = append(args, `--show-palette`)
	}
	g := New(prompt, args...)
	color, err = g.execute()

	return
}

// Password asks for a password
func Password(prompt string) (password string, err error) {
	g := New(prompt, `--password`)
	password, err = g.execute()

	return
}

// UsernameAndPassword asks for a username and password
func UsernameAndPassword(prompt string) (password, username string, err error) {
	g := New(prompt, `--password`, `--username`)
	string, err := g.execute()

	str := strings.Split(string, "|")
	username = str[0]
	password = str[1]

	return
}

func buildFileFilter(filtersMap map[string][]string) (args []string) {
	if len(filtersMap) > 0 {
		for name, patterns := range filtersMap {
			args = append(args, `--file-filter`)
			filter := fmt.Sprintf(`%s|%s`, name, strings.Join(patterns, ` `))
			args = append(args, filter)
		}
	}

	return
}

func (g *Gozenity) execute() (response string, err error) {
	cmd := exec.Command(g.command, g.arguments...)

	byteOut, err := cmd.Output()

	// Cast and trim
	response = strings.TrimSpace(string(byteOut))

	return
}
