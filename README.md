# pomo

A terminal-based pomodoro timer with task tracking. Inspired by pomotodo.com

## Goals

Provide a personal, locally hosted, low ceremony pomodoro timer with task
management. Pomodoro history includes which tasks were worked on (both Doing and
Done) to facilitate time tracking.

Store all task and pomodoro data in human-readable text files. Probably YAML.

Generate reports based on pomodoro history to facilitate time/task estimates,
e.g. for post hoc cap ex reports, or for itemizing hourly billing.

## Features

* Task management
  * Manage a list of tags that may be attached to tasks.
  * Create, edit, and delete tasks, including the task title, status (to do,
    doing, or done), tags, and optional notes.
  * Ask the user if they want to start a pomodoro if they move a task into Doing
    or Done when there is not an active pomodoro. Completed pomodoros that are
    not yet reported are considered active.
  * Present tasks in a Kanban board with columns for each status.
  * Use arrow keys to navigate tasks, and shift+arrow to move tasks around.
* Pomodoro timer
  * User can start, ~~pause, resume,~~ abort, or complete pomodoros.
  * Pomodoro timers count down automatically.
  * Service alerts the user when the pomodoro is over.
  * When a Pomodoro is done, prompt the user to update their tasks to reflect
    what they worked on, and the status of each task at the end of the pomodoro.
  * After the user updates their tasks and saves the final pomodoro state to
    history, remove all completed tasks from the Done column, and start the
    break.
  * Automatically select a short break (5 minutes) or long break (15 minutes)
    based on how many .
  * User can pause, resume, and abort breaks.
  * Service alerts the user when the break is over.
  * User explicitly starts the next pomodoro 
* Pomodoro history
  * User can browse past pomodoros in the app. They can also view them in plain
    text in the pomo data folder (probably in `~/.pomo/data/`).

Nonfunctional requirements:

* Save as you go, so no data is lost.
* Pomodoro and break timers run in a background task, so that a terminal can be
  popped 
* Guard against exiting when a pomodoro is in progress.
* Resume a running pomodoro if the user exits and later runs the application.
* Guard against multiple pomodoro processes at the same time--e.g. use a lock
  file to prevent concurrent use, or use file watchers to automatically reload
  if the current pomodoro is changed by another process.

## TODOs

* [x] save on every modify action
* [x] coalesce board and kanban packages
* [x] faint border around inactive columns (instead of invisible border)
* [x] move left goes to top, move right goes to bottom
* [x] debounce saves and show save status
* [x] auto save on timer start
* [x] resume in progress timer on start
* [x] finalize pomo when timer goes off
* [x] save completed pomos in ~/.pomo/history/
* [x] play sound when pomodoro or break ends
* [ ] show alert when pomodoro or break ends
* [ ] load pomo history on start
* [ ] show number of pomodoros completed today
* [ ] swap in current.yaml instead of saving history of changes
* [ ] guard against concurrent modification (through lock files or file watchers)
* [ ] project tags for tasks
* [ ] responsive kanban layout when the screen is too small for 3 columns
* [ ] prompt start pomodoro on task changes (e.g. when status is idle or break ended)
* [ ] show progress bar when pomodoro timer is running
* [ ] make pomodoro section more prominent when user action is wanted (idle, pomo ended, break ended)
