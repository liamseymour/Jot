package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"jot/display"
	"jot/jot"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	path, err := os.Executable()
	check(err)

	dataPath := filepath.Join(path, "../data/")

	// Create string to run regex on, exlude first arg
	// as it is always jot
	commandString := strings.Join(os.Args[1:], " ")
	switch {

	// Help, -h, --help, help, or no args
	case MatchStringAndCheck("^(-h|--help|help|)( |$)", commandString):
		// TODO help
		fmt.Println("Help unimplemented - you're out of luck.")

	// List, ls
	case MatchStringAndCheck("^ls( |$)", commandString):
		displayByTitle := false
		displayAll := false
		hasArgument := false

		// Parse command
		// Display using title
		if MatchStringAndCheck("^ls( -.+)* -t .+", commandString) {
			displayByTitle = true
		}
		// Display all
		if MatchStringAndCheck("^ls( -.+)* -a( )?", commandString) {
			displayAll = true
		}
		// Some non-option argument is passed
		if MatchStringAndCheck("^ls( -.+)* [[:word:]]", commandString) {
			hasArgument = true
		}

		// Execute
		switch {
		case displayAll:
			display.DisplayAllNotes(dataPath)
		case displayByTitle:
			display.DisplayNoteByTitle(dataPath, os.Args[len(os.Args)-1])
		case hasArgument:
			display.DisplayNoteById(dataPath, os.Args[len(os.Args)-1])
		default:
			// Todo What should default ls do?
			display.DisplayLastNote(dataPath)
		}
	// Search keywords
	case MatchStringAndCheck("^search( |$)", commandString):
		if len(os.Args) < 3 {
			fmt.Println("Insufficient arguments. Use \"jot help search\" for usage.")
			break
		}
		display.DisplayNotesBySearch(dataPath, strings.Join(os.Args[2:], " "))

	// New Note
	case MatchStringAndCheck("^new( |$)", commandString):
		popOut := false
		title := "[TITLE]"

		// Parse command
		// Popout
		if MatchStringAndCheck("^new( -.+)* -p", commandString) {
			popOut = true
		}
		// Title passed
		if MatchStringAndCheck("^new( -.+)* [[:word:]]+", commandString) {
			title = os.Args[len(os.Args)-1]
		}

		// Execute
		var note string
		switch {
		case !popOut:
			note = readNoteFromConsole(title)
		case popOut:
			note = readNoteFromTextEditor(dataPath, title)
		}

		// Always run
		newNoteId := jot.NewNote(dataPath, note)
		fmt.Printf("New note created with id: %s", newNoteId)
		fmt.Println()
		// Here we could get away with "DisplayLastNote" but its probably more
		// reliable to display by ID.
		display.DisplayNoteById(dataPath, newNoteId)

	// Delete a note
	case MatchStringAndCheck("^(del|delete|rm|remove)( |$)", commandString):
		useTitle := false
		// Parse command
		// Bad call
		if !MatchStringAndCheck("^(del|delete|rm|remove)( -.+)* [[:word:]]*", commandString) {
			fmt.Printf("Not a recognized use of %s. Use \"jot help %s\" for usage.", os.Args[1], os.Args[1])
			fmt.Println()
			return
		}
		// Delete by title
		if MatchStringAndCheck("^(del|delete|rm|remove)( -.+)* -t", commandString) {
			useTitle = true
		}

		switch {
		case useTitle:
			title := os.Args[len(os.Args)-1]
			id, found := jot.DeleteNoteByTitle(dataPath, title)
			if found {
				fmt.Printf("Note deleted with id: %s", id)
				fmt.Println()
			} else {
				fmt.Printf("No note found with title: %s", title)
				fmt.Println()
			}

		default:
			id := os.Args[len(os.Args)-1]
			title, found := jot.DeleteNote(dataPath, id)
			if found {
				fmt.Printf("Note deleted with title: %s", title)
				fmt.Println()
			} else {
				fmt.Printf("No note found with id: %s", id)
				fmt.Println()
			}
		}

	// Checking an item on the to-do / check list
	case MatchStringAndCheck("^(check)( |$)", commandString):
		// zero or one arguments given to check
		if MatchStringAndCheck("^check$", commandString) {
			fmt.Printf("Not a recognized use of %s. Use \"jot help %s\" for usage.", os.Args[1], os.Args[1])
			fmt.Println()
			return
		}
		useTitle := false

		// Reference note by title
		if MatchStringAndCheck("^(check)( -[[:word:]]*)* -t [[:word:]]+ [[:word:]]+", commandString) {
			useTitle = true
		}

		nString := os.Args[len(os.Args)-1]
		n, err := strconv.Atoi(nString)

		if err != nil || n < 0 {
			fmt.Printf("'%v' is not an non-negative integer.", nString)
			return
		}

		switch {
		// Reference note by title
		case useTitle:
			title := os.Args[len(os.Args)-2]
			item, success := jot.CheckItemByNoteTitle(dataPath, title, n)

			if success {
				fmt.Printf("Checked item: '%s' from note with title: '%s'", item, title)
				fmt.Println()
				display.DisplayNoteByTitle(dataPath, title)
			} else {
				fmt.Printf("Cannot find item number: '%d' from note with title: '%s'", n, title)
				fmt.Println()
			}

		// No options
		default:
			id := os.Args[len(os.Args)-2]
			item, success := jot.CheckItem(dataPath, id, n)

			if success {
				fmt.Printf("Checked item: '%s' from note with id: '%s'", item, id)
				fmt.Println()
				display.DisplayNoteById(dataPath, id)
			} else {
				fmt.Printf("Cannot find item number: '%d' from note with id: '%s'", n, id)
				fmt.Println()
			}
		}

	// Unchecking an item on the to-do / check list
	case MatchStringAndCheck("^(uncheck)( |$)", commandString):
		// zero or one arguments given to uncheck
		if MatchStringAndCheck("^uncheck$", commandString) {
			fmt.Printf("Not a recognized use of %s. Use \"jot help %s\" for usage.", os.Args[1], os.Args[1])
			fmt.Println()
			return
		}
		useTitle := false

		// Reference note by title
		if MatchStringAndCheck("^(uncheck)( -[[:word:]]*)* -t [[:word:]]+ [[:word:]]+", commandString) {
			useTitle = true
		}

		nString := os.Args[len(os.Args)-1]
		n, err := strconv.Atoi(nString)

		if err != nil || n < 0 {
			fmt.Printf("'%v' is not an non-negative integer.", nString)
			return
		}

		switch {
		// Reference note by title
		case useTitle:
			title := os.Args[len(os.Args)-2]
			item, success := jot.UnCheckItemByNoteTitle(dataPath, title, n)

			if success {
				fmt.Printf("Unchecked item: '%s' from note with title: '%s'", item, title)
				fmt.Println()
				display.DisplayNoteByTitle(dataPath, title)
			} else {
				fmt.Printf("Cannot find item number: '%d' from note with title: '%s'", n, title)
				fmt.Println()
			}

		// No options
		default:
			id := os.Args[len(os.Args)-2]
			item, success := jot.UnCheckItem(dataPath, id, n)

			if success {
				fmt.Printf("Unchecked item: '%s' from note with id: '%s'", item, id)
				fmt.Println()
				display.DisplayNoteById(dataPath, id)
			} else {
				fmt.Printf("Cannot find item number: '%d' from note with id: '%s'", n, id)
				fmt.Println()
			}
		}

	// Add an item to the to-do / check list
	case MatchStringAndCheck("^(add)( |$)", commandString):
		// zero or one arguments given to add
		if MatchStringAndCheck("^add$", commandString) {
			fmt.Printf("Not a recognized use of %s. Use \"jot help %s\" for usage.", os.Args[1], os.Args[1])
			fmt.Println()
			return
		}

		useTitle := false
		// Reference note by title
		if MatchStringAndCheck("^(add)( -[[:word:]]*)* -t [[:word:]]+ [[:word:]]+", commandString) {
			useTitle = true
		}

		item := os.Args[len(os.Args)-1]

		switch {
		// Reference note by title
		case useTitle:
			title := os.Args[len(os.Args)-2]
			success := jot.AddItemByNoteTitle(dataPath, title, item)

			if success {
				fmt.Printf("Added item: '%s' to note with title: '%s'", item, title)
				fmt.Println()
				display.DisplayNoteByTitle(dataPath, title)
			} else {
				fmt.Printf("Cannot find note with title: '%s'", title)
				fmt.Println()
			}

		// No options
		default:
			id := os.Args[len(os.Args)-2]
			success := jot.AddItem(dataPath, id, item)

			if success {
				fmt.Printf("Checked item: '%s' from note with id: '%s'", item, id)
				fmt.Println()
				display.DisplayNoteById(dataPath, id)
			} else {
				fmt.Printf("Cannot find note with id: '%s'", id)
				fmt.Println()
			}
		}

	// Remove an item from the to-do / check list
	case MatchStringAndCheck("^(scratch)( |$)", commandString):
		// zero or one arguments given to scratch
		if MatchStringAndCheck("^scratch$", commandString) {
			fmt.Printf("Not a recognized use of %s. Use \"jot help %s\" for usage.", os.Args[1], os.Args[1])
			fmt.Println()
			return
		}
		useTitle := false

		// Reference note by title
		if MatchStringAndCheck("^(scratch)( -[[:word:]]*)* -t [[:word:]]+ [[:word:]]+", commandString) {
			useTitle = true
		}

		nString := os.Args[len(os.Args)-1]
		n, err := strconv.Atoi(nString)

		if err != nil || n < 0 {
			fmt.Printf("'%v' is not an non-negative integer.", nString)
			return
		}

		switch {
		// Reference note by title
		case useTitle:
			title := os.Args[len(os.Args)-2]
			item, success := jot.RemoveItemByNoteTitle(dataPath, title, n)

			if success {
				fmt.Printf("Removed item: '%s' from note with title: '%s'", item, title)
				fmt.Println()
				display.DisplayNoteByTitle(dataPath, title)
			} else {
				fmt.Printf("Cannot find item number: '%d' from note with title: '%s'", n, title)
				fmt.Println()
			}

		// No options
		default:
			id := os.Args[len(os.Args)-2]
			item, success := jot.RemoveItem(dataPath, id, n)

			if success {
				fmt.Printf("Removed item: '%s' from note with id: '%s'", item, id)
				fmt.Println()
				display.DisplayNoteById(dataPath, id)
			} else {
				fmt.Printf("Cannot find item number: '%d' from note with id: '%s'", n, id)
				fmt.Println()
			}
		}
	}
}

/*************Helper Functions*************/

func MatchStringAndCheck(pattern string, string string) bool {
	match, err := regexp.MatchString(pattern, string)
	if err != nil {
		fmt.Println("Command parsing error.")
		panic(err.Error())
	}
	return match
}

func check(err error) {
	if err != nil {
		panic(err.Error())
	}
}

func readNoteFromConsole(title string) string {
	if title == "" {
		fmt.Print("New Note Title: ")
	} else {
		fmt.Printf("%v: ", title)
	}
	scanner := bufio.NewScanner(os.Stdin)
	s := ""

	for scanner.Scan() {
		s += scanner.Text()
		s += ""
	}

	return s
}

func readNoteFromTextEditor(path, title string) string {
	// To-do generalize text editor, pull path from a settings file
	// find the available programs etc...

	// create text file
	fp := filepath.Join(path, "input.txt")
	file, err := os.Create(fp)
	check(err)
	titleBytes := []byte(title)
	_, err = file.Write(titleBytes)
	check(err)
	file.Close()

	// open in sublime
	sublPath, err := exec.LookPath("subl")
	check(err)
	cmd := exec.Command(sublPath, "-n", fp)
	err = cmd.Run()
	check(err)

	// when the user says so, read it
	fmt.Println("Press enter to continue.")
	reader := bufio.NewReader(os.Stdin)
	_, _, _ = reader.ReadRune()

	bytes, err := ioutil.ReadFile(fp)
	check(err)

	// delete text file
	err = os.Remove(fp)
	check(err)

	return string(bytes)
}
