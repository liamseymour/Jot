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
	TitleColor string `json:"title-color"`
	DateColor  string `json:"date-color"`
	IdColor    string `json:"id-color"`
}

/* Get the color from string representation */
func color_of_string(colorString string) color.Color {
	switch colorString {
	case "Black":
		return color.Black
	case "White":
		return color.White
	case "Gray":
		return color.Gray
	case "Red":
		return color.Red
	case "Green":
		return color.Green
	case "Yellow":
		return color.Yellow
	case "Blue":
		return color.Blue
	case "Magenta":
		return color.Magenta
	case "Cyan":
		return color.Cyan
	default:
		return color.Black
	}
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
	titleStyle := color_of_string(style.TitleColor)
	dateStyle := color_of_string(style.DateColor)
	idStyle := color_of_string(style.IdColor)

	// Header
	fmt.Println()
	time := time.Unix(int64(note.Time), 0).Format("Jan 2 3:04 2006")
	titleStyle.Printf("%s\n", note.Title)
	fmt.Print("Taken: ")
	dateStyle.Printf("%v\n", time)
	fmt.Print("ID: ")
	idStyle.Printf("%s", note.Id)

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
