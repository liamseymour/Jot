package main

import (
	"fmt"
	"jot/jot"
	"os"
	"bufio"
	"strconv"
	//"path/filepath"
)

func main() {
	path, err := os.Executable()
	check(err)
	//path = filepath.Join(path, "../../data/notes.json")
	/* Debugging */path = "C:/Users/liamg/go/src/jot/data/notes.json"

	// User has entered no arguments
	if len(os.Args) == 1 {
		// TODO help
		return
	} 

	switch os.Args[1] {
	
	case "-h", "help":
		// TODO help

	case "ls": // display notes
		if len(os.Args) == 2 {
			jot.DisplayLastNote(path)
		} else {
			switch os.Args[2] { // list with an optional
			case "-a", "--all": // display all notes
				jot.DisplayAllNotes(path)
			}
		}

	case "search": // search for a note
		if len(os.Args) == 3 {
			jot.DisplayNotesBySearch(path, os.Args[2])
		}

	case "new": // add a new note
		note := ""
		if len(os.Args) == 3 { // If a title is supplied, use it
			note += os.Args[2] + "\n"
			note += readNoteFromConsole(false, os.Args[2])
		} else {
			note = readNoteFromConsole(true, "")
		}

		jot.NewNote(path, note)
		fmt.Println("Note added:")
		jot.DisplayLastNote(path)

	case "del", "delete", "remove": // remove a note
		if len(os.Args) == 2 {
			fmt.Println("Please specify note ID to be deleted.")
			break
		}
		if os.Args[2] == "-t" || os.Args[2] == "--title" {
			// Delete by title
			id, found := jot.DeleteNoteByTitle(path, os.Args[3])
			if (found) {
				fmt.Printf("Note deleted with id: '%s'", id)
				fmt.Println()
			} else {
				fmt.Println("No note found.")
			}
		} else {
			// Delete by ID
			title, found := jot.DeleteNote(path, os.Args[2])
			if (found) {
				fmt.Printf("Note deleted with title: '%s'", title)
				fmt.Println()
			} else {
				fmt.Println("No note found.")
			}
		}

	case "check": // check an item as complete
		n, err := strconv.Atoi(os.Args[2])
		id := os.Args[3]
		if err != nil || n < 0 {
			fmt.Printf("'%v' is not an non-negative integer.", os.Args[2])
			break
		}
		
		item, success := jot.CheckItem(path, id, n)

		if success {
			fmt.Printf("Checked item: '%s' from note with id: '%s'\n", item, id)
		} else {
			fmt.Printf("Cannot find item number: '%d' from note with id: '%s'\n", n, id)
		}

		jot.DisplayNoteById(path, id)

	case "uncheck": // uncheck an item
		// TODO
	case "add": // add a item to the checklist
		// TODO
	case "scratch": // discard an item from the checklist
		// TODO
	default:
		fmt.Printf("Unkown command: '%v'. Use 'jot help' to see a list of available commands.", os.Args[1])
	}
}

func readNoteFromConsole(getTitle bool, title string) string {
	if getTitle {
		fmt.Print("New Note Title: ")
	} else {
		fmt.Printf("%v: \n", title)
	}
	scanner := bufio.NewScanner(os.Stdin)
	s := ""

	for scanner.Scan() {
		s += scanner.Text()
		s += "\n"
	}

	return s
}

func check(err error) {
	if err != nil {
		panic(err.Error())
	}
}