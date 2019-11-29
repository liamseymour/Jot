package jot

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"
	"github.com/rs/xid"
)

// Reading and writting

/* An object representing a single note. */
type Note struct {
	Id    string   `json:"id"`
	Title string   `json:"title"`
	Time  int64  `json:"time"`
	Lines []string `json:"lines"`
	Todo  []string `json:"to-do"`
	Done  []string `json:"done"`
}

/* An object representing a collection of notes. */
type Notes struct {
	Notes []Note `json:"notes"`
}

/* Reads Json data from path and returns the Notes object. */
func fetchNotes(path string) Notes {
	file, err := os.Open(path)
	if err != nil {
		panic(err.Error())
	}
	defer file.Close()

	var notes Notes
	bytes, _ := ioutil.ReadAll(file)
	err = json.Unmarshal(bytes, &notes)
	if err != nil {
		panic(err.Error())
	}

	return notes
}

/* Writes the given notes object to the given json file. */
func writeNotes(notes Notes, path string) {
	bytes, err := json.MarshalIndent(notes, "", "    ")
	if err != nil {
		panic(err.Error())
	}
	err = ioutil.WriteFile(path, bytes, 0644)
	if err != nil {
		panic(err.Error())
	}
}

// Display

/* Displays the given note to std out. */
func displayNote(note Note) {
	// Header
	fmt.Println()
	time := time.Unix(int64(note.Time), 0).Format("Jan _2 3:04:05 2006")
	fmt.Println(note.Title, "- Taken:", time)
	fmt.Println("ID: ", note.Id)

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
}

/* Displays the given notes to std out. */
func displayNotes(notes Notes) {
	for i := 0; i < len(notes.Notes); i++ {
		displayNote(notes.Notes[i])
	}
}

/* Displays the stored notes to std out. */
func DisplayAllNotes() {
	displayNotes(fetchNotes("data/notes.json"))
}

/* Displays the last note taken to std out. */
func DisplayLastNote() {
	notes := fetchNotes("data/notes.json")
	displayNote(notes.Notes[len(notes.Notes)-1])
}

/* Displays notes with any of the keywords in the title to std out. */
func DisplayNotesBySearch(search string) {
	notes := fetchNotes("data/notes.json")
	var filtered Notes
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


// Management

/* Given a string, make a new note and record it. Return the id of the new note */
func NewNote(text string) string {
	note := parseNote(text)
	notes := fetchNotes("data/notes.json")
	notes.Notes = append(notes.Notes, note)
	writeNotes(notes, "data/notes.json")
	return note.Id
}
/* Given an id, delete the note with this id and return its contents */
//func DeleteNote(id string) string {;}
//func CheckItem(id string, listItem int) {;}
//func UncheckItem(id string, listItem int) {;}
//func AddItem(id string, listItem int) {;}
//func RemoveItem(id string, listItem int) string {;}


// Helper

/* Parses a string into a note, assuming the first line is a title and lines
 * that begin with " - " are checklist items. */
func parseNote(text string) Note {

	lines := strings.Split(text, "\n")
	if (lines[len(lines)-1] == "") {
		lines = lines[0:len(lines)-1]
	}
	var note Note
	note.Id = xid.New().String()
	note.Title = lines[0]
	lines = lines[1:] // pop title
	note.Time = time.Now().Unix()
	note.Lines = []string{}
	note.Todo = []string{}
	note.Done = []string{}
	for _, line := range lines {
		if strings.HasPrefix(line, " - ") {
			note.Todo = append(note.Todo, line[3:])
		} else {
			note.Lines = append(note.Lines, line)
		}
	}
	return note
} 
