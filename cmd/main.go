package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"jot/display"
	jot "jot/model"
	"jot/settings"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

func main() {
	// Setup flags and arguments
	var fTitle bool
	var fAll bool
	var fHeaders bool
	var fPopout bool
	var fHelp bool

	flag.BoolVar(&fTitle, "t", false, "Reference note by title instead of id.")
	flag.BoolVar(&fAll, "a", false, "Show all notes.")
	flag.BoolVar(&fHeaders, "h", false, "Show only note headers.")
	flag.BoolVar(&fPopout, "p", false, "Enter input via text editor.")
	flag.BoolVar(&fHelp, "help", false, "Show Help.")
	flag.Parse()

	command := flag.Arg(0)

	// setup paths
	exePath, err := os.Executable()
	check(err)
	dataPath := filepath.Join(exePath, "../data/")

	switch {

	// Help, -h, --help, help, or no args
	case command == "help" || fHelp || command == "":
		// TODO help
		flag.PrintDefaults()

	// List, ls
	case command == "ls":
		switch {
		case fAll && fHeaders:
			display.DisplayAllNoteHeaders()
		case fAll:
			display.DisplayAllNotes()
		case fTitle && fHeaders:
			display.DisplayNoteHeaderByTitle(os.Args[len(os.Args)-1])
		case fTitle:
			display.DisplayNoteByTitle(os.Args[len(os.Args)-1])
		case flag.Arg(1) != "" && fHeaders:
			display.DisplayNoteHeaderById(os.Args[len(os.Args)-1])
		case flag.Arg(1) != "":
			display.DisplayNoteById(os.Args[len(os.Args)-1])
		default:
			// TODO: What should default ls do?
			display.DisplayLastNote()
		}
	// Search keywords
	case command == "search":
		if fHeaders {
			display.DisplayNotesHeadersBySearch(strings.Join(flag.Args()[1:], " "))
		} else {
			display.DisplayNotesBySearch(strings.Join(flag.Args()[1:], " "))
		}

	// New Note
	case command == "new":
		title := flag.Arg(1)

		var note string
		success := true
		switch {
		case !fPopout:
			note = readNoteFromConsole(title)
		case fPopout:
			note, success = readNoteFromTextEditor(dataPath, title)
		}

		if success {
			newNoteId := jot.NewNote(note)
			fmt.Printf("New note created with id: %s", newNoteId)
			fmt.Println()
			// Here we could get away with "DisplayLastNote" but its probably more
			// reliable to display by ID.
			display.DisplayNoteById(newNoteId)
		} else {
			fmt.Printf("Cannot locate text editor. Check your settings.")
			fmt.Println()
		}

	// Delete a note
	case command == "rm" || command == "del":
		switch {
		case fTitle:
			title := flag.Arg(1)
			id, found := jot.DeleteNoteByTitle(title)
			if found {
				fmt.Printf("Note deleted with id: %s", id)
				fmt.Println()
			} else {
				fmt.Printf("No note found with title: %s", title)
				fmt.Println()
			}

		default:
			id := flag.Arg(1)
			title, found := jot.DeleteNote(id)
			if found {
				fmt.Printf("Note deleted with title: %s", title)
				fmt.Println()
			} else {
				fmt.Printf("No note found with id: %s", id)
				fmt.Println()
			}
		}

	// Checking an item on the to-do / check list
	case command == "check":
		nString := flag.Arg(2)
		n, err := strconv.Atoi(nString)

		if err != nil || n < 0 {
			fmt.Printf("'%v' is not an non-negative integer.", nString)
			return
		}

		switch {
		// Reference note by title
		case fTitle:
			title := flag.Arg(1)
			item, success := jot.CheckItemByNoteTitle(title, n)

			if success {
				fmt.Printf("Checked item: '%s' from note with title: '%s'", item, title)
				fmt.Println()
				display.DisplayNoteByTitle(title)
			} else {
				fmt.Printf("Cannot find item number: '%d' from note with title: '%s'", n, title)
				fmt.Println()
			}

		// No options
		default:
			id := flag.Arg(1)
			item, success := jot.CheckItem(id, n)

			if success {
				fmt.Printf("Checked item: '%s' from note with id: '%s'", item, id)
				fmt.Println()
				display.DisplayNoteById(id)
			} else {
				fmt.Printf("Cannot find item number: '%d' from note with id: '%s'", n, id)
				fmt.Println()
			}
		}

	// Unchecking an item on the to-do / check list
	case command == "uncheck":
		nString := flag.Arg(2)
		n, err := strconv.Atoi(nString)

		if err != nil || n < 0 {
			fmt.Printf("'%v' is not an non-negative integer.", nString)
			return
		}

		switch {
		// Reference note by title
		case fTitle:
			title := flag.Arg(1)
			item, success := jot.UnCheckItemByNoteTitle(title, n)

			if success {
				fmt.Printf("Unchecked item: '%s' from note with title: '%s'", item, title)
				fmt.Println()
				display.DisplayNoteByTitle(title)
			} else {
				fmt.Printf("Cannot find item number: '%d' from note with title: '%s'", n, title)
				fmt.Println()
			}

		// No options
		default:
			id := flag.Arg(1)
			item, success := jot.UnCheckItem(id, n)

			if success {
				fmt.Printf("Unchecked item: '%s' from note with id: '%s'", item, id)
				fmt.Println()
				display.DisplayNoteById(id)
			} else {
				fmt.Printf("Cannot find item number: '%d' from note with id: '%s'", n, id)
				fmt.Println()
			}
		}

	// Add an item to the to-do / check list
	case command == "add":
		item := flag.Arg(2)

		switch {
		// Reference note by title
		case fTitle:
			title := flag.Arg(1)
			success := jot.AddItemByNoteTitle(title, item)

			if success {
				fmt.Printf("Added item: '%s' to note with title: '%s'", item, title)
				fmt.Println()
				display.DisplayNoteByTitle(title)
			} else {
				fmt.Printf("Cannot find note with title: '%s'", title)
				fmt.Println()
			}

		// No options
		default:
			id := flag.Arg(1)
			success := jot.AddItem(id, item)

			if success {
				fmt.Printf("Checked item: '%s' from note with id: '%s'", item, id)
				fmt.Println()
				display.DisplayNoteById(id)
			} else {
				fmt.Printf("Cannot find note with id: '%s'", id)
				fmt.Println()
			}
		}

	// Remove an item from the to-do / check list
	case command == "scratch":
		nString := flag.Arg(2)
		n, err := strconv.Atoi(nString)

		if err != nil || n < 0 {
			fmt.Printf("'%v' is not an non-negative integer.", nString)
			return
		}

		switch {
		// Reference note by title
		case fTitle:
			title := flag.Arg(1)
			item, success := jot.RemoveItemByNoteTitle(title, n)

			if success {
				fmt.Printf("Removed item: '%s' from note with title: '%s'", item, title)
				fmt.Println()
				display.DisplayNoteByTitle(title)
			} else {
				fmt.Printf("Cannot find item number: '%d' from note with title: '%s'", n, title)
				fmt.Println()
			}

		// No options
		default:
			id := flag.Arg(1)
			item, success := jot.RemoveItem(id, n)

			if success {
				fmt.Printf("Removed item: '%s' from note with id: '%s'", item, id)
				fmt.Println()
				display.DisplayNoteById(id)
			} else {
				fmt.Printf("Cannot find item number: '%d' from note with id: '%s'", n, id)
				fmt.Println()
			}
		}

	// edit note
	case command == "edit":
		switch {
		case fTitle:
			title := flag.Arg(1)
			id, found := jot.GetIdFromTitle(title)
			oldText, found := jot.GetNoteString(id)

			// After getting user input, edit the note
			if found {
				written, success := readNoteFromTextEditor(dataPath, oldText)
				if success {
					if jot.EditNote(id, written) {
						fmt.Println("Success, note changed:")
						display.DisplayNoteById(id)
					} else {
						fmt.Println("Failure, note not changed.")
					}
				} else {
					fmt.Printf("Cannot locate text editor. Check your settings.")
					fmt.Println()
				}
			} else {
				fmt.Printf("No note found with title: %s", title)
				fmt.Println()
			}

		default:
			id := flag.Arg(1)
			oldText, found := jot.GetNoteString(id)

			// After getting user input, edit the note
			if found {
				written, success := readNoteFromTextEditor(dataPath, oldText)
				if success {
					if jot.EditNote(id, written) {
						fmt.Println("Success, note changed:")
						display.DisplayNoteById(id)
					} else {
						fmt.Println("Failure, note not changed.")
					}
				} else {
					fmt.Printf("Cannot locate text editor. Check your settings.")
					fmt.Println()
				}
			} else {
				fmt.Printf("No note found with id: %s", id)
				fmt.Println()
			}
		}

	// amend, edit a list item
	case command == "amend":
		n, err := strconv.Atoi(flag.Arg(2))
		check(err)
		var id string
		newItem := flag.Arg(3)

		switch {
		case fTitle:
			title := os.Args[len(os.Args)-3]
			var success bool
			id, success = jot.GetIdFromTitle(title)
			if success && jot.EditListItem(id, n, newItem) {
				fmt.Println("Success: ")
				display.DisplayNoteById(id)
			} else {
				fmt.Println("Failure: Cannot amend item.")
			}

		default:
			id = os.Args[len(os.Args)-3]
			if jot.EditListItem(id, n, newItem) {
				fmt.Println("Success: ")
				display.DisplayNoteById(id)
			} else {
				fmt.Println("Failure: Cannot amend item.")
			}
		}

	default:
		fmt.Printf("Unrecognized command: '%s'. use 'jot -help' to see a list of available commands.", command)
		fmt.Println()
	}
}

/*************Helper Functions*************/
func check(err error) {
	if err != nil {
		panic(err.Error())
	}
}

func readNoteFromConsole(title string) string {
	s := ""
	if title == "" {
		fmt.Print("New Note Title: ")
	} else {
		fmt.Printf("%v: ", title)
		s += title + "\n"
	}
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		s += scanner.Text()
		s += "\n"
	}

	return s
}

func readNoteFromTextEditor(path, seedText string) (written string, success bool) {
	// create text file
	fp := filepath.Join(path, "input.txt")
	file, err := os.Create(fp)
	check(err)
	seedTextBytes := []byte(seedText)
	_, err = file.Write(seedTextBytes)
	check(err)
	file.Close()

	// open in text editor
	success = true
	editorSettings := settings.GetTextEditor()
	editorPath := editorSettings.TextEditorPath
	editorArgs := editorSettings.TextEditorArgs

	// prepend filepath into args
	editorArgs = append([]string{fp}, editorArgs...)
	cmd := exec.Command(editorPath, editorArgs...)
	err = cmd.Run()
	if err != nil {
		success = false
	}

	// If a text editor is located
	if success {
		// when the user says so, read it
		fmt.Println("Press enter to continue.")
		reader := bufio.NewReader(os.Stdin)
		_, _, _ = reader.ReadRune()

		bytes, err := ioutil.ReadFile(fp)
		check(err)

		// delete text file
		err = os.Remove(fp)
		check(err)

		written = string(bytes)
	} else {
		written = ""
	}

	return
}
