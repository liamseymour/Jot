package display

import (
	"fmt"
	jot "jot/model"
	"jot/settings"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/gookit/color"
	"golang.org/x/crypto/ssh/terminal"
)

/* 			   	 Display 			   */
/* Displays the given note to std out using style settings from path/settings.json. */
func displayNote(note jot.Note) {
	// Load style settings
	style := settings.GetStyle()

	// Setup styles
	defaultStyle := color.New(color.FgColors["default"], color.BgColors["default"])

	contentStyle := color.New(color.FgColors[style.ContentColor], color.BgColors[style.ContentBackground])
	titleStyle := color.New(color.FgColors[style.TitleColor], color.BgColors[style.TitleBackground])
	dateStyle := color.New(color.FgColors[style.DateColor], color.BgColors[style.DateBackground])
	idStyle := color.New(color.FgColors[style.IdColor], color.BgColors[style.IdBackground])
	todoHeadStyle := color.New(color.FgColors[style.TodoHeadColor], color.BgColors[style.TodoHeadBackground])
	todoBulletStyle := color.New(color.FgColors[style.TodoBulletColor], color.BgColors[style.TodoBulletBackground])
	todoItemStyle := color.New(color.FgColors[style.TodoItemColor], color.BgColors[style.TodoItemBackground])
	doneHeadStyle := color.New(color.FgColors[style.DoneHeadColor], color.BgColors[style.DoneHeadBackground])
	doneBulletStyle := color.New(color.FgColors[style.DoneBulletColor], color.BgColors[style.DoneBulletBackground])
	doneItemStyle := color.New(color.FgColors[style.DoneItemColor], color.BgColors[style.DoneItemBackground])

	indent := ""
	for i := style.IndentWidth; i > 0; i-- {
		indent += " "
	}

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
		SplitPrintln(indent, note.Lines[i], defaultStyle, contentStyle)
	}

	// 'to-do'
	if len(note.Todo) != 0 {
		fmt.Println()
		fmt.Print(indent)
		todoHeadStyle.Printf("To-do:")
		fmt.Println()
	}
	for i := 0; i < len(note.Todo); i++ {
		prefix := fmt.Sprintf(indent+"%3d) ", i)
		SplitPrintln(prefix, note.Todo[i], todoBulletStyle, todoItemStyle)
	}

	// Done
	if len(note.Done) != 0 {
		fmt.Println()
		fmt.Print(indent)
		doneHeadStyle.Printf("Done:")
		fmt.Println()
	}
	for i := 0; i < len(note.Done); i++ {
		prefix := fmt.Sprintf(indent+"%3d) ", i)
		SplitPrintln(prefix, note.Done[i], doneBulletStyle, doneItemStyle)
	}

	fmt.Println()
}

func DisplayNoteById(id string) {
	note, found := jot.GetNoteById(id)
	if found {
		displayNote(note)
	} else {
		fmt.Println("No note found.")
	}
}

func DisplayNoteByTitle(title string) {
	id, found := jot.GetIdFromTitle(title)
	if found {
		DisplayNoteById(id)
	}
}

/* Displays the given notes to std out. */
func displayNotes(notes jot.Notes) {
	for i := 0; i < len(notes.Notes); i++ {
		displayNote(notes.Notes[i])
	}
}

/* Displays the stored notes to std out. */
func DisplayAllNotes() {
	displayNotes(jot.GetNotes())
}

/* Displays the last note taken to std out. */
func DisplayLastNote() {
	notes := jot.GetNotes()
	displayNote(notes.Notes[len(notes.Notes)-1])
}

/* Displays notes with any of the keywords in the title to std out. */
func DisplayNotesBySearch(search string) {
	notes := jot.GetNotes()
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

// Helper functions

/* Splits the string with respect to terminal width and indents based on the prefix width.
Prints all out to console. Will try to split on word breaks.*/
func SplitPrintln(prefix, str string, prefixStyle, strStyle color.Style) {
	// Terminal dimentions
	width := GetConsoleWidth()

	// no formating needed
	if len(prefix)+len(str) <= width {
		prefixStyle.Printf(prefix)
		strStyle.Printf(str)
		fmt.Println()
		return
	}

	// determine tabbing
	tab := ""
	for i := len(prefix); i > 0; i-- {
		tab += " "
	}

	// print with prefix
	breakIndex := findLastBreak(str, width-len(prefix))
	if breakIndex == -1 {
		breakIndex = width - len(prefix)
	}

	head := str[:breakIndex] // portion of str to print
	str = str[breakIndex+1:]
	prefixStyle.Printf(prefix)
	strStyle.Printf(head)
	fmt.Println()

	for len(prefix)+len(str) > width {
		breakIndex := findLastBreak(str, width-len(prefix))
		if breakIndex == -1 {
			breakIndex = width - len(prefix)
		}
		head = str[:breakIndex]
		str = str[breakIndex+1:]
		fmt.Printf(tab)
		strStyle.Printf(head)
		fmt.Println()
	}

	if len(str) > 0 {
		fmt.Printf(tab)
		strStyle.Printf(str)
		fmt.Println()
	}
}

/* find the last white space with respect to pos */
func findLastBreak(str string, pos int) int {
	for i := pos; i >= 0; i-- {
		if str[i] == ' ' || str[i] == '\t' || str[i] == '\n' {
			return i
		}
	}
	return -1
}

func GetConsoleWidth() int {
	var fd int
	if runtime.GOOS == "windows" {
		// windows needs to use stdout or will throw an error
		fd = int(os.Stdout.Fd())
	} else {
		fd = int(os.Stdin.Fd())
	}
	termWidth, _, err := terminal.GetSize(fd)
	if err != nil {
		panic(err.Error)
	}
	return termWidth
}
