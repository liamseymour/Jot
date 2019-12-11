package jot

import (
	"encoding/json"
	"fmt"
	"github.com/gookit/color"
	"github.com/rs/xid"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

// Reading and writting

/* An object representing a single note. */
type Note struct {
	Id    string   `json:"id"`
	Title string   `json:"title"`
	Time  int64    `json:"time"`
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
	note, _ := getNoteById(path, id)
	displayNote(note)
}

func DisplayNoteByTitle(path string, title string) {
	id, found := getIdFromTitle(path, title)
	if found {
		DisplayNoteById(path, id)
	}
}

/* Displays the given notes to std out. */
func displayNotes(notes Notes) {
	for i := 0; i < len(notes.Notes); i++ {
		displayNote(notes.Notes[i])
	}
}

/* Displays the stored notes to std out. */
func DisplayAllNotes(path string) {
	displayNotes(fetchNotes(path))
}

/* Displays the last note taken to std out. */
func DisplayLastNote(path string) {
	notes := fetchNotes(path)
	displayNote(notes.Notes[len(notes.Notes)-1])
}

/* Displays notes with any of the keywords in the title to std out. */
func DisplayNotesBySearch(path string, search string) {
	notes := fetchNotes(path)
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
func NewNote(path string, text string) string {
	note := parseNote(text)
	notes := fetchNotes(path)
	notes.Notes = append(notes.Notes, note)
	writeNotes(notes, path)
	return note.Id
}

/* Given an id, delete the note with this id and return its title */
func DeleteNote(path string, id string) (deletedTitle string, found bool) {
	notes := fetchNotes(path)
	found = false
	deletedTitle = ""
	for i := 0; i < len(notes.Notes); i++ {
		if notes.Notes[i].Id == id {
			found = true
			deletedTitle = notes.Notes[i].Title
			notes.Notes = append(notes.Notes[:i], notes.Notes[i+1:]...)
			break
		}
	}
	writeNotes(notes, path)
	return
}

/* Given a title, delete the first note that has the same title.
 * Return the id of the deleted note */
func DeleteNoteByTitle(path string, title string) (deletedId string, found bool) {
	deletedId, found = getIdFromTitle(path, title)
	if found {
		DeleteNote(path, deletedId)
	}
	return
}

/* Given the id of the note, check the nth item.
 * return the item and if the operation was successful. */
func CheckItem(path string, id string, n int) (item string, success bool) {
	note, foundNote := getNoteById(path, id)
	foundItem := false
	item = ""
	if n < len(note.Todo) && n >= 0 && foundNote {
		foundItem = true
		item = note.Todo[n]
		note.Todo = append(note.Todo[:n], note.Todo[n+1:]...)
		note.Done = append(note.Done, item)
	}

	success = false
	if foundItem {
		var notes Notes
		notes, success = replaceNote(path, id, note)
		if foundNote && foundItem && success {
			writeNotes(notes, path)
		}
	}

	return
}

func CheckItemByNoteTitle(path string, title string, n int) (item string, success bool) {
	id, found := getIdFromTitle(path, title)
	if found {
		return CheckItem(path, id, n)
	}
	return "", found
}

/* Given the id of the note, uncheck the nth item.
 * return the item and if the operation was successful. */
func UnCheckItem(path string, id string, n int) (item string, success bool) {
	note, foundNote := getNoteById(path, id)
	foundItem := false
	item = ""
	if n < len(note.Done) && n >= 0 && foundNote {
		foundItem = true
		item = note.Done[n]
		note.Done = append(note.Done[:n], note.Done[n+1:]...)
		note.Todo = append(note.Todo, item)
	}

	success = false
	if foundItem {
		var notes Notes
		notes, success = replaceNote(path, id, note)
		if foundNote && foundItem && success {
			writeNotes(notes, path)
		}
	}
	return
}

func UnCheckItemByNoteTitle(path string, title string, n int) (item string, success bool) {
	id, found := getIdFromTitle(path, title)
	if found {
		return UnCheckItem(path, id, n)
	}
	return "", found
}

/* Given the id of the note, uncheck the nth item.
 * return the item and if the operation was successful. */
func RemoveItem(path string, id string, n int) (item string, success bool) {
	note, foundNote := getNoteById(path, id)
	foundItem := false
	item = ""
	if n < len(note.Todo) && n >= 0 && foundNote {
		foundItem = true
		item = note.Todo[n]
		note.Todo = append(note.Todo[:n], note.Todo[n+1:]...)
	}

	success = false
	if foundItem {
		var notes Notes
		notes, success = replaceNote(path, id, note)
		if foundNote && foundItem && success {
			writeNotes(notes, path)
		}
	}
	return
}

func RemoveItemByNoteTitle(path string, title string, n int) (item string, success bool) {
	id, found := getIdFromTitle(path, title)
	if found {
		return RemoveItem(path, id, n)
	}
	return "", found
}

/* Given the id of the note, uncheck the nth item.
 * return the item and if the operation was successful. */
func AddItem(path string, id string, item string) (success bool) {
	note, foundNote := getNoteById(path, id)
	note.Todo = append(note.Todo, item)

	success = false
	if foundNote {
		var notes Notes
		notes, success = replaceNote(path, id, note)
		if foundNote && success {
			writeNotes(notes, path)
		}
	}
	return
}

func AddItemByNoteTitle(path string, title string, item string) (success bool) {
	id, found := getIdFromTitle(path, title)
	if found {
		return AddItem(path, id, item)
	}
	return found
}

//func AddItem(id string, listItem int) {;}
//func RemoveItem(id string, listItem int) string {;}

// Helper

/* Parses a string into a note, assuming the first line is a title and lines
 * that begin with " - " are checklist items. */
func parseNote(text string) Note {

	lines := strings.Split(text, "\n")
	if lines[len(lines)-1] == "" {
		lines = lines[0 : len(lines)-1]
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
		} else if strings.HasPrefix(line, " X ") {
			note.Done = append(note.Done, line[3:])
		} else {
			note.Lines = append(note.Lines, line)
		}
	}
	return note
}

func noteToString(note Note) string {
	s := note.Title + "\n"
	for _, line := range note.Lines {
		s += line + "\n"
	}
	for _, line := range note.Todo {
		s += " - " + line + "\n"
	}
	for _, line := range note.Done {
		s += " X " + line + "\n"
	}
	return s
}

func getNoteById(path string, id string) (note Note, found bool) {
	notes := fetchNotes(path)
	found = false
	for i := 0; i < len(notes.Notes); i++ {
		if notes.Notes[i].Id == id {
			found = true
			note = notes.Notes[i]
			break
		}
	}
	return
}

func getIdFromTitle(path string, title string) (id string, found bool) {
	notes := fetchNotes(path)
	found = false
	id = ""
	for i := 0; i < len(notes.Notes); i++ {
		if notes.Notes[i].Title == title {
			found = true
			id = notes.Notes[i].Id
		}
	}
	return
}

func replaceNote(path string, id string, newNote Note) (notes Notes, success bool) {
	success = false
	notes = fetchNotes(path)

	for i := 0; i < len(notes.Notes); i++ {
		// find note to replace
		if notes.Notes[i].Id == id {
			notes.Notes[i] = newNote
			success = true
			break
		}
	}

	return
}
