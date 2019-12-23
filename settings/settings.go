package settings

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
)

/* Whole (or top level) settings file */
type Settings struct {
	Style      Style      `json:"style"`
	TextEditor TextEditor `json:"text-editor"`
}

/* Style section of settings file */
type Style struct {
	IndentWidth          int    `json:"indent-width"`
	TitleColor           string `json:"title-color"`
	TitleBackground      string `json:"title-background"`
	DateColor            string `json:"date-color"`
	DateBackground       string `json:"date-background"`
	ContentColor         string `json:"content-color"`
	ContentBackground    string `json:"content-background"`
	IdColor              string `json:"id-color"`
	IdBackground         string `json:"id-background"`
	TodoHeadColor        string `json:"todo-head-color"`
	TodoHeadBackground   string `json:"todo-head-background"`
	TodoBulletColor      string `json:"todo-bullet-color"`
	TodoBulletBackground string `json:"todo-bullet-background"`
	TodoItemColor        string `json:"todo-item-color"`
	TodoItemBackground   string `json:"todo-item-background"`
	DoneHeadColor        string `json:"done-head-color"`
	DoneHeadBackground   string `json:"done-head-background"`
	DoneBulletColor      string `json:"done-bullet-color"`
	DoneBulletBackground string `json:"done-bullet-background"`
	DoneItemColor        string `json:"done-item-color"`
	DoneItemBackground   string `json:"done-item-background"`
}

/* Settings regarding the text editor used with jot */
type TextEditor struct {
	TextEditorPath string   `json:"prefered-text-editor-path"`
	TextEditorArgs []string `json:"text-editor-args"`
}

var settings Settings

/* Setup settings */
func init() {

	// open and close settings
	exePath, err := os.Executable()
	if err != nil {
		panic(err.Error())
	}

	file, err := os.Open(filepath.Join(exePath, "../data/settings.json"))
	if err != nil {
		panic("The settings file is missing: jot/data/settings.json")
	}
	defer file.Close()

	// load settings
	bytes, _ := ioutil.ReadAll(file)
	err = json.Unmarshal(bytes, &settings)
	if err != nil {
		panic("The settings file is corrupted: jot/data/settings.json")
	}
}

/* Returns the whole settings file */
func GetSettings() Settings {
	return settings
}

/* Returns the style settings */
func GetStyle() Style {
	return settings.Style
}

/* Returns the settings for the text editor used with jot */
func GetTextEditor() TextEditor {
	return settings.TextEditor
}
