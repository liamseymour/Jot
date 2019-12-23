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

var notes Notes
var path string

/* Load notes data and initialize path. */
func init() {
	// open and close settings
	exePath, err := os.Executable()
	if err != nil {
		panic(err.Error())
	}

	path = filepath.Join(exePath, "../data/notes.json")
	file, err := os.Open(path)
	if err != nil {
		panic(err.Error())
	}
	defer file.Close()

	bytes, _ := ioutil.ReadAll(file)
	err = json.Unmarshal(bytes, &notes)
	if err != nil {
		panic(err.Error())
	}
}

/* Writes notes to path. */
func writeNotes() {
	bytes, err := json.MarshalIndent(notes, "", "    ")
	if err != nil {
		panic(err.Error())
	}
	err = ioutil.WriteFile(filepath.Join(path), bytes, 0644)
	if err != nil {
		panic(err.Error())
	}
}

func GetNotes() Notes {
	return notes
}

// Management

/* Given a string, make a new note and record it. Return the id of the new note */
func NewNote(text string) string {
	note := parseNote(text)
	notes.Notes = append(notes.Notes, note)
	writeNotes()
	return note.Id
}

/* Given an id, delete the note with this id and return its title */
func DeleteNote(id string) (deletedTitle string, found bool) {
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
	writeNotes()
	return
}

/* Given a title, delete the first note that has the same title.
 * Return the id of the deleted note */
func DeleteNoteByTitle(title string) (deletedId string, found bool) {
	deletedId, found = GetIdFromTitle(title)
	if found {
		DeleteNote(deletedId)
	}
	return
}

/* Given the id of the note, check the nth item.
 * return the item and if the operation was successful. */
func CheckItem(id string, n int) (item string, success bool) {
	note, foundNote := GetNoteById(id)
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
		success = replaceNote(id, note)
		if foundNote && foundItem && success {
			writeNotes()
		}
	}

	return
}

func CheckItemByNoteTitle(title string, n int) (item string, success bool) {
	id, found := GetIdFromTitle(title)
	if found {
		return CheckItem(id, n)
	}
	return "", found
}

/* Given the id of the note, uncheck the nth item.
 * return the item and if the operation was successful. */
func UnCheckItem(id string, n int) (item string, success bool) {
	note, foundNote := GetNoteById(id)
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
		success = replaceNote(id, note)
		if foundNote && foundItem && success {
			writeNotes()
		}
	}
	return
}

func UnCheckItemByNoteTitle(title string, n int) (item string, success bool) {
	id, found := GetIdFromTitle(title)
	if found {
		return UnCheckItem(id, n)
	}
	return "", found
}

/* Given the id of the note, uncheck the nth item.
 * return the item and if the operation was successful. */
func RemoveItem(id string, n int) (item string, success bool) {
	note, foundNote := GetNoteById(id)
	foundItem := false
	item = ""
	if n < len(note.Todo) && n >= 0 && foundNote {
		foundItem = true
		item = note.Todo[n]
		note.Todo = append(note.Todo[:n], note.Todo[n+1:]...)
	}

	success = false
	if foundItem {
		success = replaceNote(id, note)
		if foundNote && foundItem && success {
			writeNotes()
		}
	}
	return
}

func RemoveItemByNoteTitle(title string, n int) (item string, success bool) {
	id, found := GetIdFromTitle(title)
	if found {
		return RemoveItem(id, n)
	}
	return "", found
}

/* Given the id of the note, uncheck the nth item.
 * return the item and if the operation was successful. */
func AddItem(id string, item string) (success bool) {
	note, foundNote := GetNoteById(id)
	note.Todo = append(note.Todo, item)

	success = false
	if foundNote {
		success = replaceNote(id, note)
		if foundNote && success {
			writeNotes()
		}
	}
	return
}

func AddItemByNoteTitle(title string, item string) (success bool) {
	id, found := GetIdFromTitle(title)
	if found {
		return AddItem(id, item)
	}
	return found
}

/* Return the string representation of a Note */
func GetNoteString(id string) (noteString string, success bool) {
	var note Note
	note, success = GetNoteById(id)
	if success {
		noteString = noteToString(note)
	} else {
		noteString = ""
	}
	return
}

/* Return the string representation of a Note with specified title */
func GetNoteStringByTitle(title string) (noteString string, success bool) {
	id, found := GetIdFromTitle(title)
	if found {
		return GetNoteString(id)
	}
	return "", found
}

/* Given an id and a string representation of a note, overwrite the note with id with the newNoteString */
func EditNote(id, newNoteString string) bool {
	// Create edited version of note
	newNote := parseNote(newNoteString)
	oldNote, found := GetNoteById(id)
	if found {
		newNote.Id = oldNote.Id
		newNote.Time = oldNote.Time

		// write it
		success := replaceNote(id, newNote)
		writeNotes()
		return success
	} else {
		return found
	}

}

/* Given a path to the folder containing notes.json, an id to a note,
and a item number to change, replace that list item with newItem. Return if
the operation was succesful or not. */
func EditListItem(id string, listItem int, newItem string) bool {
	// find note
	var noteToEdit Note
	for _, note := range notes.Notes {
		if note.Id == id {
			noteToEdit = note
			break
		}
	}

	// replace list item
	success := false
	if listItem < len(noteToEdit.Todo) {
		noteToEdit.Todo[listItem] = newItem
		success = true
		writeNotes()
	}
	return success
}

// Helper
/* Parses a string into a note, assuming the first line is a title and lines
 * that begin with " - " are checklist items. */
func parseNote(text string) Note {

	lines := strings.Split(text, "\n")
	for i := 0; i < len(lines); i++ {
		lines[i] = strings.Trim(lines[i], "\r")
	}
	// delete any trailing empty string from splitting on newlines
	if len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
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

func GetNoteById(id string) (note Note, found bool) {
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

func GetIdFromTitle(title string) (id string, found bool) {
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

/* Replace a note with id with the given note.
This does not write the notes to file but simply mutates the global notes variable. */
func replaceNote(id string, newNote Note) (success bool) {
	success = false

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
