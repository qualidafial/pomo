# pomo

`pomo` is a terminal-based pomodoro timer with task tracking, inspired by the
commercial [Pomotodo](pomotodo.com) service.

Unlike Pomotodo, `pomo` is locally hosted. All data is stored in the `~/.pomo`
directory in human-readable YAML files.

## Installation

`pomo` requires Go 1.22 or later to install:

```shell
go install github.com/qualidafial/pomo/cmd/pomo
```

You can also build the application from the repository root:

```shell
make
```

## Features

* Task management
  * Create, edit, and delete tasks, including the task title, status (to do,
    doing, or done), tags, and optional notes.
  * Present tasks in a Kanban board with columns for each status.
  * Navigate through tasks using arrow keys.
  * Move tasks around using shift+arrow keys. Moving a ticket to the right sends
    it to the bottom of the next list. Moving it left moves it to the top of the
    previous list.
* Pomodoro timer
  * User can start, cancel, or complete pomodoros.
  * Pomodoro and break timers count down automatically, and resume automatically
    when the app is closed and reopened.
  * Plays an alarm and displays a notification when a pomodoro is over.
  * After a Pomodoro is done, prompt the user to update their tasks to reflect
    what they worked on, and the status of each task at the end of the pomodoro.
  * After the user updates their tasks and completes the pomodoro, save the
    pomodoro to history, remove all completed tasks from the Done column, and
    start the break.
  * Automatically select a short break (5 minutes) or long break (15 minutes)
    based on how many pomodoros have been completed today.
  * User can pause, resume, and abort breaks.
  * Plays an alarm and display a notification when the break is over.
  * Doesn't start the next pomodoro until the user starts it.
  * Resume any running pomodoro or break timers if the user exits `pomo` and
    starts it again later.
* Saves as you go: every pomodoro action or task change is saved to disk.

## Planned features

* Task management
  * Manage a list of tags that may be attached to tasks.
  * Ask the user if they want to start a pomodoro if they move a task into Doing
    or Done when there is not an active pomodoro. Completed pomodoros that are
    not yet reported are considered active.
* Pomodoro timer
  * Automatically select a short break (5 minutes) or long break (15 minutes)
    based on how many .
  * User can pause, resume, and abort breaks.
  * Service alerts the user when the break is over.
  * User explicitly starts the next pomodoro
  * Active timers run in a background task, so that a terminal can be popped
    even when the application is closed.
* Pomodoro history
  * User can browse past pomodoros in the app. They can also view them in plain
    text in the pomo data folder (probably in `~/.pomo/data/`).
  * Generate reports based on pomodoro history to facilitate time/task
    estimates, e.g. for post hoc cap ex reports, or for itemizing hourly
    billing.
* Guard against multiple pomodoro processes at the same time--e.g. use a lock
  file to prevent concurrent use, or use file watchers to automatically reload
  if the current pomodoro is changed by another process.

## TODOs

* [ ] move pomodoro header to the bottom row, below the help line
* [ ] show progress bar when pomodoro timer is running
* [ ] make pomodoro section more prominent when user action is wanted (idle, pomo ended, break ended)
* [ ] load pomo history on start
* [ ] show number of pomodoros completed today
* [ ] swap in current.yaml instead of saving history of changes
* [ ] guard against concurrent modification (through lock files or file watchers)
* [ ] project tags for tasks
* [ ] responsive kanban layout when the screen is too small for 3 columns
* [ ] prompt start pomodoro on task changes (e.g. when status is idle or break ended)
* [ ] record time when tasks are moved to done
