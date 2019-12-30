package help

import (
	"fmt"
	"reflect"
)

/* Struct representing entire Help file */
type Help struct {
	Commands string `help:"COMMANDS"`
	Help     string `help:"HELP"`
	Ls       string `help:"LS"`
	New      string `help:"NEW"`
	Add      string `help:"ADD"`
	Check    string `help:"CHECK"`
	Uncheck  string `help:"UNCHECK"`
	Scratch  string `help:"SCRATCH"`
	Edit     string `help:"EDIT"`
	Amend    string `help:"AMEND"`
	Search   string `help:"SEARCH"`
}

// Header of each section in help document
var sections []string{}

/* Read and parse help file */
func init() {
	sections = []string{"COMMANDS", "HELP", "LS", "NEW", "ADD", "CHECK", "UNCHECK", "SCRATCH", "EDIT", "AMEND", "SEARCH"}
	
	// open and close settings
	exePath, err := os.Executable()
	if err != nil {
		panic(err.Error())
	}

	file, err := os.Open(filepath.Join(exePath, "../data/help"))
	if err != nil {
		panic("The settings file is missing: jot/data/settings.json")
	}
	defer file.Close()

	
}

/* Display the relavent section to console */
func PrintHelp(section string) (found bool) {
	fmt.Println()
	return true
}

/*  */
func parse(s string) {
	
}
