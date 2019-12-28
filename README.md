# Jot
Lightweight command line, note taking &amp; list keeping application designed with work flow and productivity in mind.

The two primary advantages jot offers over traditional or mainstream note taking applications are speed and list taking functionality. Jot attempts to strip out unnecessary overhead commonly associated with modern software. This means that jot is not as pretty but still offers some "common sense" customization and formating. List taking features is a primary goal of jot; ability to maintain to-do lists is built into the software.

# Install
In this early stage of development, installing jot is not straight forward.
Jot must be compiled (using golang) from source. 

This of course means the dependencies must be first installed.

cd into the top most directory: "jot". From this point build cmd/main.go, e.g. on windows:

`go build -o jot.exe .\cmd\main.go`

non-windows:

`go build -o jot /cmd/main.go`

Then all files in jot/data that have a `_sample` prefix should be renamed (duplicated and renamed if you plan on contributing). e.g. `notes_sample.json` -> `notes.json`. This is to avoid committing personal files.
Change settings.json to fit your needs, especially give the path to your preferred text editor.
Then jot should be in working order.

# Commands
- `help [command]`, gets help on command
- `ls [id]`, display notes
- `new [title]`, create a new note with title if provided
- `add [id] [item]`, add item to note with [id]
- `check [id] [n]` check the nth item on note with [id]
- `uncheck [id] [n]` uncheck the nth item on note with [id]
- `scratch [id] [n]` remove the nth item on note with [id]
- `edit [id]`, edit the note in preferred text editor
- `amend [id] [n] [s]`, amend the nth item of note with [id] to be [s]

# Titles vs. Ids
Be default commands take note ids instead of the user supplied titles. The rational behind this is that titles are not necessarily unique and Ids are. Requiring the user to state that they want to use a title prevents unexpected behavior. E.g. if the user has two notes with titles foo and runs "jot rm foo" then jot will remove the first note with foo as the title.

Any command that takes an id can instead take a title when the "-t" option is passed.

Jot is still in an infantile stage and may change this to be more user friendly (and quicker to use). It may be a good idea to use titles by default but warn the user if more than one note has the same title.

# Quick Tour of jot
If you just installed, running `jot ls` should display the global jot to-do list, this is because `jot ls` with no parameters simply lists the most recent note and on install you should only have one note. If you don't care about the global jot to-do list `jot rm -t jot` will delete it; we see `rm` (alternatively `remove`, `del`, and `delete`) is used to delete an entire note. Additional the option `-t` is used to refer to the note by title, we could also use `jot rm bngre9ku76li6v1ts97g`. Be aware that jot will use the **first note it finds** with the supplied title. This is not a problem if you don't have any notes with the same title.

## Making a Note
Lets take a note: `jot new` has a few forms, `jot new "foo"` starts the note with the title "foo" and prompts for the rest of the note, line by line. `jot new` is the same but will ask for a title first. Usually you will want to use the `-p` (popout) option, which takes input from an external text editor. 

Let's run `jot new -p`. If you specified a text editor in settings.json, you should be looking at a nearly blank text file. The first line is the title and all the following lines are treated as normal. If one of these lines starts with " - " then it will be treated as a list item and added to the to-do list. Otherwise it is a standard line and will have normal formatting. An example note is below:

```
foobar
this is just a normal line,
	I can include formatting such as tabs.
 - this is a list item
list items are added to the to-do list
 - this is also a list item
```

After saving and hitting enter in the terminal to confirm our note. We should see that jot has parsed and added our note. To check, we can run the command `jot ls -a` which will display all notes you have taken. A compressed form is available with `jot ls -a -h` which will show all of the headers of your notes.

## Mutating a Note
With our newly created note, lets check an item off of the list. `jot check -t foobar 0` will check the 0th to-do item from the foobar note, after running you can see that it has been added to the "done" list. To uncheck this item, use the command `jot uncheck -t foobar 0` and the change will be reverted. 

Say we realized that we have something else to do, we can add a to-do item with `jot add -t foobar "Just one more thing"`. At the same time we realized that the second item on our list is not necessary, it can be removed entirely with `jot scratch -t foobar 1`.

I realized that I want my lists items to use proper grammar, so lets change "this is a list item" to "This is a list item." with `jot amend -t foobar 0 "This is a list item."`

If many changes are to be made it is best to use `jot edit -t foobar`. This will allow for editing in a text editor. If there are any completed list items, they will be preceded by " X ".
