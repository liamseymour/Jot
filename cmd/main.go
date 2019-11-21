package main

import (
	"jot"
	"os"
)

func main() {
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
			jot.DisplayLastNote()
		} else {
			switch os.Args[2] { // list with an optional
			case "-a", "--all": // display all notes
				jot.DisplayAllNotes()
			}
		}
	case "search": // search for a note
		if len(os.Args) == 3 {
			jot.DisplayNotesBySearch(os.Args[2])
		}
	case "new": // add a new note
		// TODO
	case "remove": // remove a note
		// TODO
	case "check": // check an item as complete
		// TODO
	case "uncheck": // uncheck an item
		// TODO
	case "add": // add a item to the checklist
		// TODO
	case "scratch": // discard an item from the checklist
		// TODO
	
	// TEMP
	case "test":
		jot.Write()
	} 
}