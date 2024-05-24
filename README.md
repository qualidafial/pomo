# pomo

`pomo` is a terminal-based pomodoro timer with task tracking, inspired by the
commercial [Pomotodo](pomotodo.com) service.

Unlike Pomotodo, `pomo` is locally hosted. All data is stored in the `~/.pomo`
directory in human-readable YAML files.

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
  * User can start or cancel breaks.
  * Plays an alarm and display a notification when the break is over.
  * Doesn't start the next pomodoro until the user starts it.
  * Resume any running pomodoro or break timers if the user exits `pomo` and
    starts it again later.
* Saves as you go: every pomodoro action or task change is saved to disk.

## Installation

`pomo` requires Go 1.22 or later to install:

```shell
go install github.com/qualidafial/pomo/cmd/pomo
```

You can also build the application from the repository root:

```shell
make
```

## Configuration

`pomo` may be configured by modifying `~/.pomo/config.yaml`, which is
automatically generated the first time it runs:

```yaml
pomo:
    daily-goal: 8
timer:
    break: 5m
    long-break: 15m
    pomodoro: 25m
```
