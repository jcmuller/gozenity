package gozenity_test

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"github.com/jcmuller/gozenity"
)

func ExampleList() {

	output, err := gozenity.List(
		"Choose an option:",
		"One word",
		"Two",
		"Three things",
	)

	if err != nil {
		if err, ok := err.(*gozenity.EmptySelectionError); !ok {
			fmt.Println(fmt.Errorf("Error getting output: %s", err))
		}
	}

	fmt.Println(output)
	// Output: Two
}

func ExampleCalendar() {
	args := time.Now()
	entry, err := gozenity.Calendar("Please select 4/12/2018", args)

	if err != nil {
		panic(err)
	}

	fmt.Println(entry)
	// Output: 04/12/2018
}

func ExampleEntry() {
	entry, err := gozenity.Entry("Please type an answer (expecting 'Foo'):", "Placeholder")

	if err != nil {
		panic(err)
	}

	fmt.Println(entry)
	// Output: Foo
}

func ExampleError() {
	err := gozenity.Error("Something turrible happened :(.")
	if err != nil {
		panic(err)
	}
	// Output:
}

func ExampleInfo() {
	gozenity.Info("This thing happened.")
	// Output:
}
func ExampleFileSelection() {
	files, err := gozenity.FileSelection("Choose that guy", nil)

	if err != nil {
		panic(err)
	}

	fmt.Println(strings.Join(files, " "))
	// Output: file1 file2
}

func ExampleFileSelection_second() {
	filters := map[string][]string{
		`Go files`:       {`*.go`},
		`Markdown files`: {`*.md`, `*.markdown`},
	}

	files, err := gozenity.FileSelection("Choose that guy", filters)

	if err != nil {
		panic(err)
	}

	fmt.Println(strings.Join(files, " "))
	// Output: file1 file2
}

func ExampleDirectorySelection() {
	dirs, err := gozenity.DirectorySelection("Choose that guy")

	if err != nil {
		panic(err)
	}

	fmt.Println(strings.Join(dirs, " "))
	// Output: dir1 dir2
}

// This is broken
func _ExampleNotification() {
	err := gozenity.Notification("This happened. Deal.")

	if err != nil {
		panic(err)
	}

	// Output:
}

func ExampleProgress() {
	progress := make(chan int)

	go gozenity.Progress("We're doing that", progress)

	for i := 0; i <= 100; i++ {
		progress <- i
		time.Sleep(time.Millisecond * 50)
	}

	// Output: 2
}

// This example ties a scale to a progress bar.
func ExampleProgress_second() {
	progress := make(chan int)

	go gozenity.Progress("Here's the thing", progress)

	args := &gozenity.ScaleArgs{
		Max:     100,
		Partial: true,
		Stream:  progress,
	}

	_, err := gozenity.Scale("Select a value", args)

	if err != nil {
		log.Fatalf("Error getting scale: %q", err)
	}

	// Output: foo
}

func ExampleQuestion() {
	answer, err := gozenity.Question("Who? Answer 'her'.")

	if err != nil {
		panic(err)
	}

	fmt.Println(answer)
	// Output: true
}

func ExampleWarning() {
	err := gozenity.Warning("This thing happened. Bad.")

	if err != nil {
		panic(err)
	}

	// Output:
}

func ExampleScale() {
	val, err := gozenity.Scale(
		"Select a value",
		&gozenity.ScaleArgs{Initial: 30, Max: 100},
	)

	if err != nil {
		panic(err)
	}

	fmt.Println(val)
	// Output: 23
}

func ExampleTextInfo() {
	args := &gozenity.TextInfoArgs{
		Editable: true,
	}

	_, err := gozenity.TextInfo("Do something here:", args)

	if err != nil {
		fmt.Println(err)
	}

	// Output: One of Filename, Text or URL need to be supplied
}

func ExampleTextInfo_second() {
	args := &gozenity.TextInfoArgs{
		Editable: true,
		Text: `Hello, worldly
worlded world.`,
	}

	string, err := gozenity.TextInfo("Do something here:", args)

	if err != nil {
		panic(err)
	}

	fmt.Println(string)

	// Output: Hello, worldly
	// worlded world.
}

func ExampleTextInfo_third() {
	args := &gozenity.TextInfoArgs{
		Checkbox: "Agree?",
		Editable: true,
		Text: `Hello, worldly
worlded world.`,
	}

	string, err := gozenity.TextInfo("Do something here:", args)

	if err != nil {
		panic(err)
	}

	fmt.Println(string)

	// Output: Hello, worldly
	// worlded world.
}

func ExampleTextInfo_fourth() {
	args := &gozenity.TextInfoArgs{
		Editable: true,
		URL:      `https://google.com`,
	}

	string, err := gozenity.TextInfo("Do something here:", args)

	if err != nil {
		panic(err)
	}

	fmt.Println(string)
	// Output:
}

func ExampleTextInfo_fifth() {
	tmpfile, err := ioutil.TempFile("", "filenametext")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(tmpfile.Name()) // clean up

	if _, err = tmpfile.WriteString("Hello, world!"); err != nil {
		log.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		log.Fatal(err)
	}

	args := &gozenity.TextInfoArgs{
		Editable: true,
		Filename: tmpfile.Name(),
	}

	string, err := gozenity.TextInfo("Do something here:", args)

	if err != nil {
		panic(err)
	}

	fmt.Println(string)
	// Output: Hello, world!
}

func ExampleTextInfo_sixth() {
	args := &gozenity.TextInfoArgs{
		Editable: true,
		Filename: "/tmp/foobar.txt",
	}

	_, err := gozenity.TextInfo("Do something here:", args)

	fmt.Println(err)

	// Output: stat /tmp/foobar.txt: no such file or directory
}

func ExampleTextInfo_seventh() {
	args := &gozenity.TextInfoArgs{
		Editable: true,
		Text: `Hello, worldly
worlded world.`,
		URL: "https://goobar.com",
	}

	_, err := gozenity.TextInfo("Do something here:", args)

	if err != nil {
		fmt.Println(err)
	}

	// Output: Only one of Filename, Text and URL can be supplied
}

func ExampleColorSelection() {
	color, err := gozenity.ColorSelection("Choose green:", "green", false)

	if err != nil {
		panic(err)
	}

	fmt.Println(color)
	// Output: rgb(0,128,0)
}

func ExampleColorSelection_second() {
	color, err := gozenity.ColorSelection("Choose green:", "green", true)

	if err != nil {
		panic(err)
	}

	fmt.Println(color)
	// Output: rgb(0,128,0)
}

func ExamplePassword() {
	password, err := gozenity.Password("Enter password:")

	if err != nil {
		panic(err)
	}

	fmt.Println(password)
	// Output: hunter2
}

func ExampleUsernameAndPassword() {
	password, username, err := gozenity.UsernameAndPassword("Enter password:")

	if err != nil {
		panic(err)
	}

	fmt.Printf("%s\n%s\n", username, password)
	// Output: user
	// hunter2
}
