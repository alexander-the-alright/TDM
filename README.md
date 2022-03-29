# ToDo Manager
## _Command-Line Todo List Manager_

It's really just a quick program that I wrote to help myself keep
track of things I needed to do, and I like using the terminal a whole
lot. I wrote this in Go, because Python felt like cheating, and I was
having a hard time getting Rust to be happy about variable-sized maps.

## Installation

[GoLang](https://go.dev/doc/install) needs to be installed already.

## Usage

Run TDM with

```sh
./tdm
```

Once inside, you can run the following forms of commands:

### Delete
Delete the requested argument.
```
-> delete board|task|subtask <b|t|s>
```

### Make
Create a new object. Exists in three footprints
- Make a board
```
-> make [board] <b>
```
- Make a task
```
-> make [board] <b> task <t>
```
- Make a subtask
```
-> make task <t> subtask <s>
```

### Mark
Mark an object as completed. Exists in two footprints
- Mark board as completed
```
-> mark [board] <b> [100]
```
- Mark task or subtask. Fill can be negative, and caps at 100
```
-> mark task|subtask <t|s> <fill>
```

### Show
Print contents of requested arguments.
```
-> show [all|[board] <b>|task <t>|subtask <s>]
```

