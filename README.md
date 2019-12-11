# Jot
Lightweight note taking &amp; list keeping application.

# Install
Currently, Jot can be installed by building main.go to the top level directory.

# Commands
- help [command], gets help on command
  - no arguments, general help
- ls [id], display notes
  - -t, by note title instead of id
  - -a, all notes
  - no arguments, most recent note
- new [title], create a new note with title if provided
  - -p, popout into sublime
- add [id] [item], add item to note with [id]
  - -t, by note title instead of id
- check [id] [n] check the nth item on note with [id]
  - -t, by note title instead of id
- uncheck [id] [n] uncheck the nth item on note with [id]
  - -t, by note title instead of id
- scratch [id] [n] remove the nth item on note with [id]
  - -t, by note title instead of id
