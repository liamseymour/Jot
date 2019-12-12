package display

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"jot/jot"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gookit/color"
)

// Whole settings file
type Settings struct {
	Style Style `json:"style"`
}

// Style section of settings file
type Style struct {
	TitleColor      string `json:"title-color"`
	TitleBackground string `json:"title-background"`
	DateColor       string `json:"date-color"`
	DateBackground  string `json:"date-background"`
	IdColor         string `json:"id-color"`
	IdBackground    string `json:"id-background"`
}

/* 			   	 Display 			   */
/* Displays the given note to std out using style settings from path/settings.json. */
func displayNote(dataPath string, note jot.Note) {
	// Load style settings
	file, err := os.Open(filepath.Join(dataPath, "/settings.json"))
	if err != nil {
		panic(err.Error())
	}
	defer file.Close()

	var settings Settings
	bytes, _ := ioutil.ReadAll(file)
	err = json.Unmarshal(bytes, &settings)
	if err != nil {
		panic(err.Error())
	}
	style := settings.Style

	// Setup styles
	titleStyle := color.New(color.FgColors[style.TitleColor], color.BgColors[style.TitleBackground])
	dateStyle := color.New(color.FgColors[style.DateColor], color.BgColors[style.DateBackground])
	idStyle := color.New(color.FgColors[style.IdColor], color.BgColors[style.IdBackground])

	// Header
	fmt.Println()
	time := time.Unix(int64(note.Time), 0).Format("Jan 2 3:04 2006")
	titleStyle.Printf(note.Title)
	fmt.Println()
	fmt.Print("Taken: ")
	dateStyle.Printf("%v", time)
	fmt.Println()
	fmt.Print("ID: ")
	idStyle.Printf(note.Id)

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
	displayNote(path, note)
}

func DisplayNoteByTitle(path string, title string) {
	id, found := jot.GetIdFromTitle(path, title)
	if found {
		DisplayNoteById(path, id)
	}
}

/* Displays the given notes to std out. */
func displayNotes(path string, notes jot.Notes) {
	for i := 0; i < len(notes.Notes); i++ {
		displayNote(path, notes.Notes[i])
	}
}

/* Displays the stored notes to std out. */
func DisplayAllNotes(path string) {
	displayNotes(path, jot.FetchNotes(path))
}

/* Displays the last note taken to std out. */
func DisplayLastNote(path string) {
	notes := jot.FetchNotes(path)
	displayNote(path, notes.Notes[len(notes.Notes)-1])
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
	displayNotes(path, filtered)
}
