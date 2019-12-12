package display

import (
	"fmt"
	"jot/jot"
	"strings"
	"time"

	"github.com/gookit/color"
)

/* 			   	 Display 			   */
/* Displays the given note to std out. */
func displayNote(note jot.Note) {
	// Header
	fmt.Println()
	time := time.Unix(int64(note.Time), 0).Format("Jan 2 3:04 2006")
	color.Yellow.Printf("%s\n", note.Title)
	fmt.Print("Taken: ")
	color.Blue.Printf("%v\n", time)
	fmt.Print("ID: ")
	color.Blue.Printf("%s", note.Id)

	// Lines
	if len(note.Lines) != 0 {
		fmt.Println()
	}
	for i := 0; i < len(note.Lines); i++ {
		fmt.Println(" ", note.Lines[i])
	}

	// to-do
	if len(note.Todo) != 0 {
		fmt.Println()
		fmt.Println("  To-do: ")
	}
	for i := 0; i < len(note.Todo); i++ {
		fmt.Printf("%5v", i)
		fmt.Println(")", note.Todo[i])
	}

	// Done
	if len(note.Done) != 0 {
		fmt.Println()
		fmt.Println("  Done: ")
	}
	for i := 0; i < len(note.Done); i++ {
		fmt.Printf("%5v", i)
		fmt.Println(")", note.Done[i])
	}

	fmt.Println()
}

func DisplayNoteById(path string, id string) {
	note, _ := jot.GetNoteById(path, id)
	displayNote(note)
}

func DisplayNoteByTitle(path string, title string) {
	id, found := jot.GetIdFromTitle(path, title)
	if found {
		DisplayNoteById(path, id)
	}
}

/* Displays the given notes to std out. */
func displayNotes(notes jot.Notes) {
	for i := 0; i < len(notes.Notes); i++ {
		displayNote(notes.Notes[i])
	}
}

/* Displays the stored notes to std out. */
func DisplayAllNotes(path string) {
	displayNotes(jot.FetchNotes(path))
}

/* Displays the last note taken to std out. */
func DisplayLastNote(path string) {
	notes := jot.FetchNotes(path)
	displayNote(notes.Notes[len(notes.Notes)-1])
}

/* Displays notes with any of the keywords in the title to std out. */
func DisplayNotesBySearch(path string, search string) {
	notes := jot.FetchNotes(path)
	var filtered jot.Notes
	keywords := strings.Split(search, " ")

	// First find notes with the keywords in the title
	for i := 0; i < len(notes.Notes); i++ {
		for j, found := 0, false; j < len(keywords) && !found; j++ {
			if strings.Contains(strings.ToLower(notes.Notes[i].Title), strings.ToLower(keywords[j])) {
				filtered.Notes = append(filtered.Notes, notes.Notes[i])
				found = true
			}
		}
	}
	displayNotes(filtered)
}
