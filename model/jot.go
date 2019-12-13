package jot

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/rs/xid"
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

/* Reads Json data from path/notes.json and returns the Notes object. */
func FetchNotes(path string) Notes {
	file, err := os.Open(filepath.Join(path, "/notes.json"))
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

/* Writes the given notes object to the given path/notes.json. */
func writeNotes(notes Notes, path string) {
	bytes, err := json.MarshalIndent(notes, "", "    ")
	if err != nil {
		panic(err.Error())
	}
	err = ioutil.WriteFile(filepath.Join(path, "/notes.json"), bytes, 0644)
	if err != nil {
		panic(err.Error())
	}
}

// Management

/* Given a string, make a new note and record it. Return the id of the new note */
func NewNote(path string, text string) string {
	note := parseNote(text)
	notes := FetchNotes(path)
	notes.Notes = append(notes.Notes, note)
	writeNotes(notes, path)
	return note.Id
}

/* Given an id, delete the note with this id and return its title */
func DeleteNote(path string, id string) (deletedTitle string, found bool) {
	notes := FetchNotes(path)
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
	deletedId, found = GetIdFromTitle(path, title)
	if found {
		DeleteNote(path, deletedId)
	}
	return
}

/* Given the id of the note, check the nth item.
 * return the item and if the operation was successful. */
func CheckItem(path string, id string, n int) (item string, success bool) {
	note, foundNote := GetNoteById(path, id)
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
	id, found := GetIdFromTitle(path, title)
	if found {
		return CheckItem(path, id, n)
	}
	return "", found
}

/* Given the id of the note, uncheck the nth item.
 * return the item and if the operation was successful. */
func UnCheckItem(path string, id string, n int) (item string, success bool) {
	note, foundNote := GetNoteById(path, id)
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
	id, found := GetIdFromTitle(path, title)
	if found {
		return UnCheckItem(path, id, n)
	}
	return "", found
}

/* Given the id of the note, uncheck the nth item.
 * return the item and if the operation was successful. */
func RemoveItem(path string, id string, n int) (item string, success bool) {
	note, foundNote := GetNoteById(path, id)
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
	id, found := GetIdFromTitle(path, title)
	if found {
		return RemoveItem(path, id, n)
	}
	return "", found
}

/* Given the id of the note, uncheck the nth item.
 * return the item and if the operation was successful. */
func AddItem(path string, id string, item string) (success bool) {
	note, foundNote := GetNoteById(path, id)
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
	id, found := GetIdFromTitle(path, title)
	if found {
		return AddItem(path, id, item)
	}
	return found
}

// Helper
/* Parses a string into a note, assuming the first line is a title and lines
 * that begin with " - " are checklist items. */
func parseNote(text string) Note {

	lines := strings.Split(text, "\n")
	for i := 0; i < len(lines); i++ {
		lines[i] = strings.Trim(lines[i], "\r")
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

func GetNoteById(path string, id string) (note Note, found bool) {
	notes := FetchNotes(path)
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

func GetIdFromTitle(path string, title string) (id string, found bool) {
	notes := FetchNotes(path)
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
	notes = FetchNotes(path)

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
