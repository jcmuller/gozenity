// Package gozenity is a simple wrapper for zenity) in Go.
package gozenity

import (
	"fmt"
	"io"
	"log"
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
	arguments = append(titles, arguments...)

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
	args := []string{`--list`, `--hide-header`, `--column`, prompt}
	args = append(args, options...)
	g := New(prompt, args...)
	selection, err = g.Execute()
	return
}

// Entry asks for input
func Entry(prompt, placeholder string) (entry string, err error) {
	g := New(prompt, `--entry`, `--entry-text`, placeholder)
	entry, err = g.Execute()

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

	date, err = g.Execute()

	return
}

// Error errors errors
func Error(prompt string) (err error) {
	g := New(prompt, `--error`, `--ellipsize`)
	_, err = g.Execute()
	return
}

// Info informs information
func Info(prompt string) (err error) {
	g := New(prompt, `--info`, `--ellipsize`)

	_, err = g.Execute()
	return
}

// FileSelection opens a file selector
func FileSelection(prompt string, filtersMap map[string][]string) (files []string, err error) {
	args := []string{`--file-selection`, `--multiple`}
	filters := buildFileFilter(filtersMap)
	args = append(args, filters...)

	g := New(prompt, args...)
	result, err := g.Execute()
	files = strings.Split(result, `|`)
	return
}

// DirectorySelection opens a file selector
func DirectorySelection(prompt string) (files []string, err error) {
	g := New(prompt, `--file-selection`, `--multiple`, `--directory`)
	result, err := g.Execute()
	files = strings.Split(result, `|`)
	return
}

// Notification notifies notifiees
func Notification(prompt string) (err error) {
	g := New(prompt, `--notification`, `--listen`)
	_, err = g.Execute()
	return
}

// Progress shows a pogress bar modal
func Progress(prompt string, progress <-chan int) (err error) {
	g := New(prompt, "--progress", "--auto-close", "--auto-kill")

	cmd := exec.Command(g.command, g.arguments...)

	stdin, err := cmd.StdinPipe()
	if err != nil {
		log.Fatalf("Error getting pipe: %s", err)
	}

	go func(stdin io.WriteCloser) {
		defer stdin.Close()

		for {
			select {
			case p, ok := <-progress:
				if !ok {
					return
				}

				io.WriteString(stdin, fmt.Sprintf("%d\n", p))
			}
		}
	}(stdin)

	err = cmd.Run()

	return
}

type ChecklistOptions struct {
	Question              string
	CheckColumnName       string
	DescriptionColumnName string
	Checks                []bool
	Descriptions          []string
	Width                 uint
	Height                uint
}

func Checklist(options ChecklistOptions) ([]string, error) {
	if len(options.Question) == 0 {
		return nil, fmt.Errorf("you should be asking a question")
	}
	if len(options.Descriptions) != len(options.Checks) {
		return nil, fmt.Errorf("the amount of description entries and check entries should be equal")
	}
	if options.Width == 0 {
		options.Width = 800
	}
	if options.Height == 0 {
		options.Height = 400
	}
	if len(options.CheckColumnName) == 0 {
		options.CheckColumnName = "check"
	}
	if len(options.DescriptionColumnName) == 0 {
		options.DescriptionColumnName = "description"
	}

	checklistSlice := []string{
		"--list",
		"--checklist",
		fmt.Sprintf("--width=%d", options.Width),
		fmt.Sprintf("--height=%d", options.Height),
		fmt.Sprintf("--column=%s", options.CheckColumnName),
		fmt.Sprintf("--column=%s", options.DescriptionColumnName),
	}

	for i := range options.Checks {
		state := "FALSE"
		if options.Checks[i] {
			state = "TRUE"
		}

		checklistSlice = append(checklistSlice, state)
		checklistSlice = append(checklistSlice, options.Descriptions[i])
	}

	g := New(options.Question, checklistSlice...)

	selection, err := g.Execute()
	if err != nil {
		return nil, err
	}

	if len(selection) == 0 {
		return make([]string, 0), nil
	}

	selections := strings.Split(selection, "|")

	return selections, nil
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

	_, err = g.Execute()
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

	ans, err := g.Execute()

	if ans == "" {
		return -1, err
	}

	answer, nerr := strconv.Atoi(ans)

	if nerr != nil {
		log.Fatalf("Error converting to int: %s", nerr)
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
	color, err = g.Execute()

	return
}

// Password asks for a password
func Password(prompt string) (password string, err error) {
	g := New(prompt, `--password`)
	password, err = g.Execute()

	return
}

// UsernameAndPassword asks for a username and password
func UsernameAndPassword(prompt string) (password, username string, err error) {
	g := New(prompt, `--password`, `--username`)
	string, err := g.Execute()

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

func (g *Gozenity) Execute() (response string, err error) {
	cmd := exec.Command(g.command, g.arguments...)

	byteOut, err := cmd.Output()

	// Cast and trim
	response = strings.TrimSpace(string(byteOut))

	return
}
