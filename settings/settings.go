package settings

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
)

// Whole (or top level) settings file
type Settings struct {
	Style          Style  `json:"style"`
	TextEditorPath string `json:"prefered-text-editor-path"`
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

/* Load the settings.json file in the given path into a Settings struct.*/
func LoadSettings(path string) Settings {
	file, err := os.Open(filepath.Join(path, "/settings.json"))
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

	return settings
}

func GetStyle(path string) Style {
	return LoadSettings(path).Style
}

func GetTextEditorPath(path string) string {
	return LoadSettings(path).TextEditorPath
}
