package jot

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"
	// "strconv"
)

// Individual note struct used to load json data
type Note struct {
	Id    string   `json:"id"`
	Title string   `json:"title"`
	Time  float64  `json:"time"`
	Lines []string `json:"lines"`
	Todo  []string `json:"to-do"`
	Done  []string `json:"done"`
}

// All of the Notes
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


func writeNotes(notes Notes, path string) {
	bytes, err := json.Marshal(notes)
	if err != nil {
		panic(err.Error())
	}
	err = ioutil.WriteFile(path, bytes, 0644)
	if err != nil {
		panic(err.Error())
	}
}

/* Displays the given note to std out */
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

/* Displays the given notes to std out */
func displayNotes(notes Notes) {
	for i := 0; i < len(notes.Notes); i++ {
		displayNote(notes.Notes[i])
	}
}

/* Displays the stored notes to std out */
func DisplayAllNotes() {
	displayNotes(fetchNotes("data/temp-notes.json"))
}

/* Displays the last note taken to std out */
func DisplayLastNote() {
	notes := fetchNotes("data/temp-notes.json")
	displayNote(notes.Notes[len(notes.Notes)-1])
}

/* Displays notes with any of the keywords in the title to std out */
func DisplayNotesBySearch(search string) {
	notes := fetchNotes("data/temp-notes.json")
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

func Write() {
	writeNotes(fetchNotes("data/temp-notes.json"), "data/notes.json")
}
