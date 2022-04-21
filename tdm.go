// =============================================================================
// Auth: Alex Celani
// File: tdm.go
// Revn: 04-20-2022  6.0
// 
// Func: display and manage progress of a litany of items to be done,
//       with an organization scheme similar to Trello. It's CLI
//       Trello
//
// TODO: add second object to hold command information
//       make subsubtasks for greater depth
// M        or rooms to contain boards
// M     remember last args
//       move most (read: all) helper functions to external files
// M     add priority to tasks
//       attach text files???
//       contemplate adding XOR of data in file, altho who cares
//       contemplate adding battery of tests
//       contemplate moving functions to external file
//       side-by-side printing
//       colored printing?
//       keep track of finished items
//       periodic tasks? something to do with the time package
//       convert file delimiters to (customizable?) variables
//       do some more legit testing
// =============================================================================
//                                 JUMP TAGS
//    for the lucky few who use vim, press * on a tag to jump to that function
//                    CALC CATCH INPUT JOIN FILEIN FILEOUT
//                HELP DELETE FIND СДЕЛАЙ MARK PARSE SHOW MAIN
// =============================================================================
// CHANGE LOG
// -----------------------------------------------------------------------------
// 10-04-2021: init
//             wrote bones of string parsing into data structure
// 10-05-2021: commented
// 10-06-2021: entirely rewrote text format and text processor
// 10-09-2021: added reading from file
//             wrote bones to while loop
//             created empty functions for commands
//             wrote, but did not test, getting user input
// 10-10-2021: fixed user input
//             wrote help()
//             wrote second draft of show(), fix immediately
// 02-21-2022: gutted main loop
//             standardized method to pass commands to subfunctions
//             commented main()
// 02-22-2022: made help() change log and comment block
//             commented help()
//             added XXX tags to "if flag {" statements
//             added searchable structure to ease defining searches
//                  with multiple different types of function calls
//             began restructuring show() to focus on sanitizing
//                  inputs before show() gets called
// 02-23-2022: changed searchable to query
// 02-24-2022: reduced scope of sanitize() to work with a more rigid
//                  set of inputs, just to get a stable release
//             finished sanitize()
//             tested sanitize()
//             commented sanitize()
//             wrote the /all/ portion of show()
// 02-25-2022: wrote rest of show()
//             rewrote /show/ portion of help()
//             connected sanitize() to show(), added show() to main()
// 02-27-2022: commented /all/ portion of show()
// 03-01-2022: planning for find() and delete()
//             wrote and tested find(), checks out
//             changed delete() to del() to avoid overloading builtin
//             wrote and tested first draft of del()
// 03-02-2022: commented /board/, /task/, and /subtask/ portions of
//                  show()
//             wrote mark()
//             made show() accept commands of len greater than 3
//             commented first draft of mark()
// 03-03-2022: commented first draft of find()
//             commented del()
// 03-04-2022: improved error condition in mark()
// 03-06-2022: wrote first draft of parse() to replace sanitize()
// 03-07-2022: wrote first draft of сделай()
// 03-08-2022: tested следай()
//             wrote and tested fileOut()
//             watched Taylor Tomlinson's Look At You
//             updated fileIn() to work with files that don't exist
//             updated work functions to take queries instead of
//                  []strings
// 03-09-2022: updated function comment blocks
// 03-10-2022: commented fileOut()
//             began commenting parse()
//             began updating help()
//*03-11-2022: updated help()
//             commented /move/ and /mark/ of parse()
//             super fixed subtask propogation bug in fileIn()
// 03-12-2022: rewrote /delete/, /show/, and /mark/
// 03-13-2022: added missing call to lower() in input()
//             rewrote /make/
//             began testing /delete/
// 03-14-2022: finished testing /delete/
//             tested /show/, /mark/, and /make/
// 03-15-2022: integrated parse.go as parse()
//             added jump tags :)
//             commented parse()
// 03-20-2022: first draft rewrite of find()
//*03-27-2022: fixed bug where tasks can't be made under boards with
//                  name of length 1
//*03-29-2022: added, tested, and commented function for catching
//                  keyboard interrupts
// 04-07-2022: fixed bug where task fill did not persist when task
//                  contained subtasks
// 04-10-2022: changed type of fill in task/subtask from int to string
//*04-11-2022: wrote, tested, and integrated calc() into show() and
//                  сделай()
//*04-12-2022: calls to os.Args in main() to work with cmd line args
//             fixed bug where fileIn() crashes on empty board
// 04-14-2022: moved data file to home directory with os.UserHomeDir()
//                  changed arguments for fileIn() and fileOut()
//             added quick little subcommands in show() and parse()
//                  /show/ to print boards or tasks
// 04-15-2022: added new function, Help() for when functions fail
//             added /about/ section to help()
//*04-20-2022: fixed bug where keyboard interrupts would fail to write
//                  data out to file
//
// =============================================================================


package main

import ( 
    "fmt"           // used for Println
    "strings"       // used for Split, ToLower, Join
    "strconv"       // used for Atoi
    "os"            // used for Exit, Stdin, Args, UserHomeDir
    "os/signal"     // used for Notify
    "syscall"       // used for SIGTERM
    "io/ioutil"     // used for ReadFile, WriteFile
    "bufio"         // used for NewScanner
)


// aliases for common functions
var p      = fmt.Print
var print  = fmt.Println
var printf = fmt.Printf
var split  = strings.Split
var lower  = strings.ToLower
var atoi   = strconv.Atoi
var itoa   = strconv.Itoa
var exit   = os.Exit
var dir    = os.UserHomeDir
var read   = ioutil.ReadFile
var write  = ioutil.WriteFile


// CALC
// quick little function that calculates the average fill of the
// subtasks contained inside of a task
func calc( t task ) string {

    // init counter
    var count int = 0
    // grab length to determine if count is necessary
    var tlen int = len( t.subt )

    // simply don't iterate if []subtask is empty
    if tlen != 0 {
        // otherwise, iterate over subtasks
        for _, sub := range t.subt {
            // convert to int for addition
            fill, err := atoi( sub.fill )
            if err != nil {     // if conversion returns and error...
                fill = 0        // just default to 0
            }                   // sux, just enter a valid fill, idiot
            count = count + fill    // continue counting up
        }
        count = count / tlen        // actually make it an average
    }
    return itoa( count )            // return conversion to string
}


// CATCH    finally something fun
// asynchronous function that catches keyboard interrupts and exits
// gracefully
func catch() {
    // blocking channel of type signal
    c := make( chan os.Signal )

    // I didn't read what this does, but it was in an example and it
    // worked, so
    // from what I can ascertain, when SIGTERM is detected, it binds
    // Interrupt to channel c
    signal.Notify( c, os.Interrupt, syscall.SIGTERM )

    // asynchronous goroutine
    go func() {
        <-c     // immediately halts in blocking read for signal

        if changed {    // if there was a change to the map...
            fileOut()   // rewrite data file
        }

        print( "\r" )   // this helps with the missing newline

        exit( 0 )   // and exit gracefully
    }()
}


// INPUT
// wrapper function to print input prompt, take user input
func input( ps1 string ) string {
    p( ps1 )                                // print prompt
    // create scanner on stdin
    scanner := bufio.NewScanner( os.Stdin )
    scanner.Scan()                          // read input
    return lower( scanner.Text() )          // return text
}


// JOIN
// quick function to join an array of strings, with delimiting spaces
func join( subcommands []string ) string {
    
    var ret string  // declare returnable string

    // iterate over contained strings inside array
    for _, word := range subcommands {
        // add word and delimiting space to returnable
        ret = ret + word + " "
    }

    // return everything but last character, a trailing space
    return ret[:len( ret ) - 1]

}


// FILEOUT
// wrapper function to contain all the file writing
func fileOut() {

    // initialize data to be written as empty string
    var data string = ""
    // keep track of how many boards have been written, so the
    // delimiter isn't written after the last board
    var boardCount int = 0

    // iterate over all boards
    for bname, bcontents := range boards {

        // increment number of boards written
        // actually caused a medium bug, because it's incremented
        // before the check against len(), so if != len() - 1 placed
        // delimiters in the wrong places
        boardCount++

        data = data + bname         // add board name to writestring
        data = data + "|"           // XXX task delimiter

        // iterate over tasks
        for tind, tcontents := range bcontents {
            data = data + tcontents.name    // add task name
            data = data + ","               // XXX task delimiter

            data = data + tcontents.fill    // add task fill
            // does the task contain subtasks?
            if len( tcontents.subt ) != 0 {
                // XXX task's subtask delimiter
                data = data + ":"
                // iterate over a task's subtasks
                for sind, scontents := range tcontents.subt {
                // add subtask name

                    // XXX
                    if flag {   // if debug set, print incoming
                                // subtask and data before...
                        print( "subt.name: ", scontents.name )
                        print( "pre data: ", data )
                    }

                    // add subtask name
                    data = data + scontents.name

                    // XXX
                    if flag {   // and after
                        print( "post data: ", data )
                    }

                    data = data + ">"       // XXX subtask fill delim

                    // add subtask fill
                    data = data + scontents.fill

                    // check if NOT last subtask
                    if sind != len( tcontents.subt ) - 1 {
                        data = data + "."   // XXX subtask delimiter
                    }
                }
            }
            // check if NOT last task
            if tind != len( bcontents ) - 1 {
                data = data + ";"           // XXX task delimiter
            }
        }
        // check if NOT last board
        // slightly different than before, explained above
        if boardCount != len( boards ) {
            data = data + "$"               // XXX board delimiter
        }
    }

    // get home directory from user
    home, err := dir()

    if err != nil {     // check if home dir failed
        panic( err )    // panic if so
    }

    // XXX
    if flag {   // if debug set, print home directory
        print( "home (fO): ", home )
    }

    // cast writestring to byte array, write to file named "data"
    // if a file doesn't exist, make with 0x644 permissions
    // take error as err, if fails
    err = write( home + "/data", []byte( data ), 0644 )

    // XXX
    if flag {   // if debug set, print data written
        print( data )
    }

    if err != nil {     // if write() failed, tell user
        print( "Error in fileOut(): file not properly written" )
        print( err )
    }

}


// FILEIN
// warpper function to contain all the file reading
func fileIn() {
    
    // test string to use in lieu of a file
    // file := 
    // "b1|tA,:sI>30.sII>30;tB,25;tC,40$b2|tA,5;tB,"
    //  b1|tA,:sI>30.sII>30;tB,25;tC,40               ~ first board
    //                                 $              ~ board delimiter
    //                                  b2|tA,5;tB,   ~ second board
    //  b1                              b2            ~ board names
    //    |tA,:sI>30.sII>30;tB,25;tC,40   |tA,5;tB,   ~ board contents
    //    |                               |           ~ task marker
    //     ta               tB    tC       tA   tB    ~ task names
    //                     ;     ;             ;      ~ task delimiter
    //       ,:sI>30.sII>30   ,25   ,40      ,5   ,   ~ task contents
    //       ,                ,25   ,40      ,5   ,   ~ task fill
    //       ,                ,     ,        ,    ,   ~ task fill delim
    //        :sI>30.sII>30                           ~ task subtasks
    //        :                                       ~ subtask marker
    //         sI    sII                              ~ subtask names
    //              .                                 ~ subtask delim
    //           >30    >30                           ~ subtask fill
    //           >      >                             ~ subtask fill delim

    // get home directory from user
    home, err := dir()

    if err != nil {     // check if home dir failed
        panic( err )    // panic if so
    }

    // XXX
    if flag {   // if debug set, print home directory
        print( "home (fI): ", home )
    }

    file, err := read( home + "/data" )         // open file, failure set err

    // declare board, because if statement gets mad if I don't
    var board []string

    // XXX
    if flag {   // if debug set, print raw file
        print( "file: ", string( file ) )
    }

    // no file -> stop trying to build object, let boards be empty
    if err != nil {
        // notify user and leave
        print( "No file found" )
        return
    } else {    // if there is a file...
        // $ -> separates different boards
        board = split( string( file ), "$" )
    }

    // iterate over all boards
    for a := 0; a < len( board ); a++ {
        // split board by name and tasks
        bContents := split( board[a], "|" )
        bName := bContents[0]       // assign name of board to var
        bTasks := bContents[1]      // assign tasks of board to var

        // XXX
        if flag {   // if debug set, print board name
            print( "board name -> ", bName )
        }

        // if there is no text right of the task marker
        if len( bTasks ) == 0 {
            // literally just initialize a new empty array of tasks
            // and add to the map with the board name
            boards[bName] = []task{}
            continue        // jump to the next board because
                            // there's no contents to process
        }

        // separate tasks from each other
        tTasks := split( bTasks, ";" )
        var tk task     // make empty task

        // iterate over all tasks
        for b := 0; b < len( tTasks ); b++ {

            // XXX
            if flag {   // if debug set, print raw tasks
                print( "tTasks -> ", tTasks[b] )
            }

            // split subtasks out of tasks
            subs := split( tTasks[b], ":" )
            // len is not one, there are subtasks
            // i.e. no reason to process subtasks
            if len( subs ) != 1 {
                // split list of subtasks
                sTasks := split( subs[1], "." )
                // iterate over subtasks
                for c := 0; c < len( sTasks ); c++ {
                    // split subtasks name from fill
                    sContents := split( sTasks[c], ">" )
                    // create subtask object with this extracted info
                    stk := subtask { name: sContents[0],
                                     fill: sContents[1] }
                    // add subtask object to array of such
                    subtasks = append( subtasks, stk )
                }
                // add full subtask array to task object
                tk.subt = subtasks
            }
            // separate task from fill
            tContents := split( tTasks[b], "," )

            // if there is no fill, interpret this as 0 fill
            if len( tContents[1] ) == 0 {
                // inferred fill is a string, so it can be processed
                // like all the other fills
                tContents[1] = "x"
            }
            tk.name = tContents[0]          // set name of task
            // convert fill to int
            // tContents[1] contains ALL contents of task
            // splitting over : separates fill from subtasks
            // split()[0] is fill
            tk.fill = split( tContents[1], ":" )[0]
            tasks = append( tasks, tk )     // add task to task array

            // XXX
            if flag {   // if debug set, print array of tasks
                print( "[]task -> ", tasks )
            }

            // these lines are so stupid
            tk.subt = nil
            subtasks = nil
            // if tk.subt is not set to nil, then the same array would
            // go to every other task subtask array that follows
        }

        // make a key/value pair out of the board name and task array
        // add it to hashmap
        boards[bName] = tasks
        tasks = nil     // nil out task array, same reason with subt
    }

    // XXX
    if flag {   // if debug set, print entire structure
        print( "\n\nboards -> ", boards, "\n\n" )
    }
}


// -----------------------------------------------------------------------------
//   HELP( commands []string )
//
//   Revn: 04-14-2022
//   Args: commands - array of strings
//                    list of commands, should begin with "help" here
//                    any subcommands would specify with which
//                    commands that help is needed
//   Func: Prints help messages to the user. If the user requests help about
//         a specific topic, function will print only that topic, otherwise
//         print everything
//   TODO: FIXME rewrite
//         add /about/ subcommand
//
//         add fail response? change structure to a switch/case
//         check that commands[0] is actualy "help"
//         make more targeted help prompts for when help has args
// -----------------------------------------------------------------------------
//   CHANGE LOG
// - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
//   ??-??-????: wrote bones to
//   02-22-2022: made change log
//               commented
//   02-25-2022: rewrote /show/ to reflect new scope
//   03-10-2022: rewrote /delete/ to reflect new scope
//               removed /move/ and /choose/ portions, because those
//                  commands were removed
//   04-15-2022: added /about/ section
//
// -----------------------------------------------------------------------------
func help( commands []string ) {


    var helpFlag bool = true                // init flag to view all
    var c string                            // init var for command[1]

    if len( commands ) > 1 {                // check for plural commands

        c = commands[1]                     // capture first command
        // recall that commands[0] is the word "help"

        // set flag to see every entry on user request
        helpFlag = c == "all" || c == ""
    }

    // XXX
    if flag {   // if debug set, print commands[0] and helpFlag
        print( "\n\nc ->", c )
        print( "helpFlag ->", helpFlag )
        print( "\n" )

    }

    print()     // print extra line to avoid cluttered look

    // help text really should speak for itself
    // important things to note: helpFlag is set true by default
    // helpFlag is set false if user asks for specific entry
    // c is set to the specific entry
    // if the specific entry is not present, nothing is printed
    // this really should be changed to a switch/case tbh
    if helpFlag || c == "delete" {
        print( " delete [board] <b>" )
        print( "   \\-- delete board <b>" )
        print( "   \\    [board|task|subtask] <b|t|s>" )
        print( "   \\-- delete board <b>, task <t>, or subtask <s>, provided argument exists " )
    }
    if helpFlag || c == "make" {
        print( " make [board] <b>" )
        print( "   \\-- create new board, <b>" )
        print( "   \\  [board] <b> task <t>" )
        print( "   \\-- create new task <t> under board <b>" )
        print( "   \\  task <t> subtask <s>" )
        print( "   \\-- create new subtask <s> under task <t>" )    
    }
    if helpFlag || c == "mark" {
        print( " mark task <t> <fill>" )
        print( "   \\        \\-- fill -> 0 - 100" )
        print( "   \\-- set progress of <t> to <fill>" )
        print( " mark subtask <s> <fill>" )
        print( "   \\           \\-- fill -> 0 - 100" )
        print( "   \\-- set progress of <s> to <fill>" )
    }
    if helpFlag || c == "show" {
        print( " show [all]" )
        print( "   \\-- display everything" )
        print( " show [board] <b>" )
        print( "   \\-- display contents of board <b>" )
        print( " show task <t>" )
        print( "   \\-- display contents of task <t>" )
        print( " show subtask <s>" )
        print( "   \\-- display contents of subtask <s>" )
    }
    if helpFlag || c == "exit" || c == "quit" {
        print( " exit" )
        print( " quit" )
        print( "   \\-- exit program" )    
    }
    if helpFlag || c == "save" {
        print( " save" )
        print( "   \\-- save changes to file without exiting" )
    }
    if c == "about" {
        print( " about" )
        print( "   \\-- kanban board in the terminal" )
    }
    print()     // print extra line to avoid cluttered look

}


// -----------------------------------------------------------------------------
//   HELP( command string )
//
//   Revn: 04-15-2022
//   Args: commands - string
//   Func: more in depth help messages than the normal function, to be
//         called from a different function
//   TODO: rewrite?
// -----------------------------------------------------------------------------
//   CHANGE LOG
// - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
//   04-15-2022: init
//
// -----------------------------------------------------------------------------
func Help( command string ) {
    switch command {
        case "delete":
            print( " usage - /delete/" )
            print( "   \\-- permanently remove an object" )
            print( "   |" )
            print( "   \\-> delete [board] <b>")
            print( "   |" )
            print( "   \\-- deletes a board named /b/" )
            print( "   \\   if the flag /board/ is not present, it will be assumed" )
            print( "   \\   if /b/ doesn't exist, /delete/ will alert user" )
            print( "   |" )
            print( "   \\-> delete task <t>" )
            print( "   |" )
            print( "   \\-- deletes a task named /t/" )
            print( "   \\   will fail without /task/ flag" )
            print( "   \\   if /t/ doesn't exist, /delete/ will alert user" )
            print( "   |" )
            print( "   \\-> delete subtask <s>" )
            print( "   |" )
            print( "   \\-- deletes a subtask named /s/" )
            print( "   \\   will fail without /subtask/ flag" )
            print( "   \\   if /s/ doesn't exist, /delete/ will alert user" )
            print()
        case "make":
            print( " usage - /make/" )
            print( "   \\-- create a new object" )
            print( "   |" )
            print( "   \\-> make [board] <b>" )
            print( "   |" )
            print( "   \\-- creates a new board named /b/" )
            print( "   \\   if the flag /board/ is not present, it will be assumed" )
            print( "   \\   board names cannot contain the word \"task\"" )
            print( "   |" )
            print( "   \\-> make [board] <b> task <t>" )
            print( "   |" )
            print( "   \\-- creates a new task named <t> under board <b>" )
            print( "   \\   will fail without /task/ flag" )
            print( "   \\   if the flag /board/ is not present, it will be assumed" )
            print( "   \\   neither board nor task names may contain the word \"task\"" )
            print( "   |" )
            print( "   \\-> make task <t> subtask <s>" )
            print( "   |" )
            print( "   \\-- creates a new subtask named <s> under task <t>" )
            print( "   \\   will fail without /task/ and /subtask/ flags" )
            print( "   \\   neither task nor subtask names may contain the words" )
            print( "   \\   \"task\" or \"subtask\"" )
            print()
        case "mark":
            print( " usage - /mark/" )
            print( "   \\-- set a fill of a task or subtask" )
            print( "   |" )
            print( "   \\-> mark task <t> <fill>" )
            print( "   |" )
            print( "   \\-- set the status of task /t/ to /fill/" )
            print( "   \\   /fill/" )
            print( "        \\-- must be an integer" )
            print( "        \\-- negative numbers are allowed (and lowkey encouraged)" )
            print( "        \\-- numbers above 100 are truncated" )
            print( "   \\   will fail without /task/ flag" )
            print( "   \\   task names cannot contain the word \"task\"" )
            print( "   |" )
            print( "   \\-> mark subtask <s> <fill>" )
            print( "   |" )
            print( "   \\-- set the status of subtask /t/ to /fill/" )
            print( "   \\   /fill/" )
            print( "        \\-- must be an integer" )
            print( "        \\-- negative numbers are allowed (and lowkey encouraged)" )
            print( "        \\-- numbers above 100 are truncated" )
            print( "   \\   will fail without /subtask/ flag" )
            print( "   \\   subtask names cannot contain the words" )
            print( "   \\   \"task\" or \"subtask\"" )
            print()
        case "show":
            print( " usage - /show/" )
            print( "   \\-- print user-requested objects" )
            print( "   |" )
            print( "   \\-> show boards" )
            print( "   |" )
            print( "   \\-- print a quick list of all boards" )
            print( "   |" )
            print( "   \\-> show tasks" )
            print( "   |" )
            print( "   \\-- print a quick list of all tasks (under their respective boards)" )
            print( "   |" )
            print( "   \\-> show [all]" )
            print( "   |" )
            print( "   \\-- print everything" )
            print( "   |" )
            print( "   \\-> show [board] <b>" )
            print( "   |" )
            print( "   \\-- print everything inside board /b/" )
            print( "   \\   if the flag /board/ is not present, it will be assumed" )
            print( "   |" )
            print( "   \\-> show task <t>" )
            print( "   |" )
            print( "   \\-- print parent board of task /t/, task /t/, and its contents" )
            print( "   \\   will fail without /task/ flag" )
            print( "   |" )
            print( "   \\-> show subtask <s>" )
            print( "   |" )
            print( "   \\-- show the parent board and task of subtask /s/" )
            print( "   \\   will fail without /subtask/ flag" )
            print()
        default:
            print( "what? how?" )
    }
}


// -----------------------------------------------------------------------------
//   DELETE( search query )
//
//   Revn: 03-15-2022
//   Args: search - query
//                  object that summarizes what needs deleting
//   Func: delete entries
//   Meth: call find() to get position of wanted entry, replace with
//         last item in list, replace list with [:len()-1] slice
//   TODO:
// -----------------------------------------------------------------------------
//   CHANGE LOG
// - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
//   03-01-2022: major planning for deleting functionality
//               plans for moving search functionality to new find()
//                  function
//               wrote bones, not working
//   03-03-2022: commented
//   03-08-2022: changed arguments from commands []string to search query
//   03-15-2022: added quick check for search.comm = fail
//
// -----------------------------------------------------------------------------
func del( search query ) {

    // determine what needs searching

    // check to see if query returned empty
    // not entirely sure if that's possible, but check for it
    if search == ( query{} ) {
        return
    }

    // if fail is set, quit
    if search.command == "fail" {
        Help( "delete" )
        return
    }

    // determine "coordinates" of object that needs finding
    // x -> ( string ) name of board, shouldn't ever be blank
    // y -> ( int ) index of task array for a given board
    // z -> ( int ) index of subtask array for a given task
    x, y, z := find( search )

    if x != "" {        // should be a given
        changed = true
        if y == -1 {    // not task, should be board
            // maps have easy delete function
            // delete, from map 'boards', key 'x'
            delete( boards, x )
        } else {            // yes task
            if z == -1 {    // not subtask, should be task
            // XXX it really needs to be mentioned that everything
            // here has to be referenced absolutely, because if you
            // make some other variable about it, like
            //      b := boards[x][y]
            // when you go to reassign the array, it won't get rid of
            // the last task, of which there are now two

            //  boards          -> map[string][]task
            //  boards[x]       -> []task
            //  boards[x][y]    -> task

            //  int                      len( boards[x] ) - 1
            //  task           boards[x][len( boards[x] ) - 1]
                // take task at the end of the array, place it over
                // unwanted task
                boards[x][y] = boards[x][len( boards[x] ) - 1]
                // take slice of everything but the last task (
                // because there's two of them in the array ), assign
                // it to array to replace array with two copies
                boards[x] = boards[x][:len( boards[x] ) - 1]
            } else {                    // yes subtask
            // oh boy
            //  boards                  -> map[string][]task
            //  boards[x]               -> []task
            //  boards[x][y]            -> task
            //  boards[x][y].subt       -> []subtask
            //  boards[x][y].subt[z]    -> subtask

            //  int                                      len( boards[x][y].subt ) - 1
            //  subtask                boards[x][y].subt[len( boards[x][y].subt ) - 1]
                // take subtask at the end of the array, place it over
                // unwanted subtask
                boards[x][y].subt[z] = boards[x][y].subt[len(
                boards[x][y].subt ) - 1]
                // take slice of everything but the last subtask (
                // because there's two of them in the array ), assign
                // it to array to replace array with two copies
                boards[x][y].subt = boards[x][y].subt[:len( boards[x][y].subt ) - 1]
            }
        }
    } else {    // should be unreachable
        print( "Error in delete(): object does not exist or not specified" )
    }
}


// -----------------------------------------------------------------------------
//   FIND( search query ) ( string, int, int )
//
//   Revn: 03-27-2022
//   Args: search  - query
//                   object that outlines what needs finding
//   Retn: x, y, z - string, int, int
//                   "coordinates" that outline where sought-after
//                      objects are
//   Func: iterate over all everything until exact positions of
//         sought-after object is found
//   TODO: error messages
//         make returnables into struct ( "coords"? )
// -----------------------------------------------------------------------------
//   CHANGE LOG
// - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
//   03-01-2022: essentially stole from show()
//   03-03-2022: commented
//   03-20-2022: rewrote for updated query definition
//   03-27-2022: commented
//
// -----------------------------------------------------------------------------
func find( search query ) ( string, int, int ) {
    
    var rboard string       // init string as empty
    var rtask int = -1      // init ints as -1, the forbidden index
    var rsubt int = -1


    // look at subcommand A specifically
    switch search.subA {
        case "board":
            // prs is true if search.board is in boards
            _, prs := boards[search.board]
            if prs {    // if prs _is_ true, keep track of that board
                rboard = search.board
            }
        case "task":
            // iterate over boards in boards
            // board -> names of boards ( key in map )
            // contents -> []task ( value in map )
            for board, contents := range boards {
                // iterate through task array
                // ind -> index
                // task -> actual task object
                for ind, task := range contents {
                    // if name of found task is requested task name
                    if task.name == search.task {
                        // keep track of board and index
                        rboard = board
                        rtask = ind
                    }
                }
            }
        case "subtask":
            // iterate over boards in boards
            // board -> names of boards ( key in map )
            // contents -> []task ( value in map )
            for board, contents := range boards {
                // iterate through task array
                // tind -> task index
                // task -> actual task object
                for tind, task := range contents {
                    // iterate through subtask array
                    // sind -> subtask index
                    // subtask -> actual subtask object
                    for sind, subtask := range task.subt {
                        // if name of found subtask is requested
                        // subtask name
                        if subtask.name == search.subtask {
                            // keep track of board, task index, and
                            // subtask index
                            rboard = board
                            rtask = tind
                            rsubt = sind
                        }
                    }
                }
            }
        default:
            // if the subcommand A is not a known object, error out
            // I'm like 90% sure that this is unreachable, but
            // whatever
            print( "Error in find()" )
            print( "Object /", search.subA, "/ not recognized" )
    }

    // return returnable 3-tuple
    return rboard, rtask, rsubt

}


// -----------------------------------------------------------------------------
//   СДЕЛАЙ( search query )
//
//   Revn: 04-11-2022
//   Args: search - query
//                  object that summarizes what needs making
//   Func: create new boards, tasks, and subtasks
//   Meth: find parent object, create new child object, attach
//   TODO: add debug
// -----------------------------------------------------------------------------
//   CHANGE LOG
// - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
//   03-07-2022: wrote first draft
//   03-10-2022: commented
//   03-16-2022: rewrote to work with updated query{} definition
//   03-27-2022: added check for command == "fail"
//   04-11-2022: made auto-update the default for task fill
//
// -----------------------------------------------------------------------------
func сделай( search query ) {

    if flag {   // if debug set, print contents of argument
        print( search )
    }

    // if fail is set, quit
    if search.command == "fail" {
        Help( "make" )
        return
    }

    // init variable to mark if parent object exists
    var exists bool = false
    // declare variable to determine what caused error
    var err string

    if search.subtask == "" {           // if subtask not specified
        if search.task == "" {          // if task not specified
            if search.board == "" {     // if board also not specified
                // failure, make called without argument
                // print error, call help about it
                print( "Error in make()" )
                print( "No object specified" )
                Help( "make" )
                return
            } else {                    // yes board
                exists = true   // signal that parent did exist
                // declare empty task array
                var t []task
                // map board name to empty task array
                boards[search.board] = t
                changed = true  // signal that file needs rewriting
            }
        } else {                        // yes task
            if search.board == "" {     // if no parent board
                err = "board"
                // failure, called to make a task belonging to no
                // board
                // print error, call help about it
                print( "Error in make()" )
                print( "Cannot make task without specifying board" )
                Help( "make" )
                return
            } else {                    // board is specified
                exists = true           // signal parent did exist
                // create task with user-given name
                // and default to auto-update task fill
                tk := task{ name: search.task, 
                            fill: "x" }
                // grab []task that given board maps to
                // i literally don't want to type search.board out two
                // more times, v clunky and causes line overrun
                preTask := boards[search.board]
                if len( preTask ) == 0 {    // if empty
                    // declare empty task array
                    var t []task
                    t = append( t, tk )     // place inside array
                    // map task array to new board
                    boards[search.board] = t
                } else {                    // array not empty
                    // simply append to pre-existing task array
                    boards[search.board] = append( preTask, tk )
                }
                changed = true          // signal file needs rewriting
            }
        }
    } else {                                // subtask specified
        if search.task == "" {              // no task tho
            err = "task"
            // failure, called to make subtask with a parent task
            // print error, call help about it
            print( "Error in make()" )
            print( "Cannot make subtask with specifying task" )
            Help( "make" )
            return
        } else {                            // yes task
            // iterate over boards to look for requested task
            for bname, bcontents := range boards {
                // iterate over tasks in boards
                for ind, t := range bcontents {
                    // if task name is given parent task name
                    if t.name == search.task {
                        exists = true   // signal parent exists
                        // create subtask with user-given name
                        stk := subtask{ name: search.subtask }
                        // if task's subtask array is empty
                        if len( t.subt ) == 0 {
                            // declare empty subtask array
                            var s []subtask
                            // place created subtask inside array
                            s = append( s, stk )
                            // map subtask array to task
                            boards[bname][ind].subt = s
                        } else {            // array not empty
                            // simply append to pre-existing subtask
                            // array
                            boards[bname][ind].subt = append( t.subt, stk )
                        }
                        // signal file needs rewriting
                        changed = true
                    }
                }
            }
        }
    }

    if !exists {        // nothing was changed, because parent object
                        // doesn't exist
        // print error, call help about it
        print( "Error in make()" )
        p( "Parent ", err )
        if err == "board" {
            p( " /", search.board, "/ " )
        } else if err == "task" {
            p( " /", search.task, "/ " )
        } else {
            p( " object " )
        }
        print( "does not exist" )
        Help( "make" )
    }

}


// -----------------------------------------------------------------------------
//   MARK( search query )
//
//   Revn: 04-10-2022
//   Args: search - query
//                  object that summarizes what needs showing
//   Func: allow user to complete tasks
//   Meth: call to find() to get "coords", set object at those coords
//              to fill
//   TODO:
// -----------------------------------------------------------------------------
//   CHANGE LOG
// - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
//   03-02-2022: planned and wrote first draft, needs much error
//                  checking
//               commented first draft
//   03-04-2022: made error condition show more information
//   03-09-2022: fixed error condition bug where program didn't know
//                  what kind of object user and printed hanging slash
//   03-11-2022: added reference to variable 'change', so fileOut()
//                  will run
//   03-17-2022: added references to updated definition of a query
//   04-10-2022: replaced fill setting with string type instead of int
//
// -----------------------------------------------------------------------------
func mark( search query ) {

    // check to see if query returned empty
    // not entirely sure if that's possible, but check for it
    if search == ( query{} ) {
        return
    }

    // if parse() failed, simply return
    // probably print failure too?
    if search.command == "fail" {
        Help( "mark" )
        return
    }

    // determine "coordinates" of object that needs finding
    // x -> ( string ) name of board, shouldn't ever be blank
    // y -> ( int ) index of task array for a given board
    // z -> ( int ) index of subtask array for a given task
    x, y, z := find( search )
    
    // markFill isn't actually used anymore, but it's good for input
    // sanitization
    markFill, err := atoi( search.subB )
    
    // on error, print failure, call help about it, return
    if err != nil {
        print( "Fill value /" + search.subB + "/ is not a number" )
        Help( "mark" )
        return
    }

    // round fill down to 100 max, specifically allow negative fill
    if markFill > 100 {
        search.subB = "100"
    }

    if x != "" {        // should be true, board is given
        if y == -1 {    // not task, should be board
            // shouldn't be possible, boards don't have a fill
            print( "Error in mark(): task does not exist or not specified" )
        } else {           // yes task
            if z == -1 {   // not subtask, should be task
                // boards            -> map[string][]task
                // boards[x]         -> []task
                // boards[x][y]      -> task
                // boards[x][y].fill -> task.fill, fill field in task
                boards[x][y].fill = search.subB
                changed = true  // signal that file needs rewriting
            } else {       // yes subtask
                // boards                    -> map[string][]task
                // boards[x]                 -> []task
                // boards[x][y]              -> task
                // boards[x][y].subt         -> []subtask
                // boards[x][y].subt[z]      -> subtask
                // boards[x][y].subt[z].fill -> subtask.fill
                boards[x][y].subt[z].fill = search.subB
                changed = true  // signal that file needs rewriting
            }
        }
    } else {    // reachable if the user gives input that DNE
        // notify user requested object did not exist
        p( "Error in mark(): " )
        p( search.subA + " /" )
        switch search.subA {
            case "board":
                p( search.board )
            case "task":
                p( search.task )
            case "subtask":
                p( search.subtask )
            default:
                p( "object" )
        }
        print( "/ does not exist" )
    }
}


// -----------------------------------------------------------------------------
//   PARSE( commands []string )
//
//   Revn: 04-14-2022
//   Args: commands - []string
//                    list of commands of what board, task, or subtask
//                    to make
//   Retn: search - query
//                  object that contains pertinent information
//   Func: parse user input, determine what command to call and about
//         what object is being called about
//   Meth: lots and lots of switch statements
//         switch over first command
//              switch over length of subcommand
//                  switch over contents of subcommand
//   TODO: move to external file?
// -----------------------------------------------------------------------------
//   CHANGE LOG
// - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
//   03-06-2022: planned
//               wrote everything but help case
//               tested pass cases
//   03-10-2022: commented /delete/ portion
//               commented /show/ portion
//   03-12-2022: rewrote /delete/, /show/, and /mark/
//   03-13-2022: rewrote /make/
//               began testing /delete/
//   03-14-2022: finished testing /delete/
//               tested /show/, /mark/, and /make/
//   03-15-2022: commented
//               fixed small bugs related to valid keyword positions
//   03-17-2022: added small reference to /save/ in the bugout case
//   03-27-2022: fixed bug in /make/ where tasks under boards with
//                  names of length 1 failed ( literally changed "> 1"
//                  to "> 0" )
//   04-14-2022: added recognition to /boards/ and /tasks/ subcommand
//                  inside of /show/
//
// -----------------------------------------------------------------------------
func parse( commands []string ) query {

    // declare return object
    var search query

    // separate command from subcommands, take note of amount of
    // subcommands
    comm := commands[0]
    subcommands := commands[1:]
    subcomlen := len( subcommands )

    // set pre-emptively, clear if it's bad
    // unused as of right now
    search.command = comm

    // big switch statement over command
    switch comm {
        // these don't really need processing, break out immediately
        case "save", "set", "quit", "exit", "help", "":
            return search
        case "delete":
            // there are only so many types of commands that work here
            //      -> delete <board>
            //      -> delete board <board>
            //      -> delete task <task>
            //      -> delete subtask <subtask>
            // i.e. subcommand can have either length 1 or 2

            // if user enters JUST delete, set fail and return
            if subcomlen == 0 {
                search.command = "fail"
                return search
            }

            switch subcommands[0] {
                case "task":
                    // delete task xyz
                    // set subcommand ( focus ) of search
                    search.subA = subcommands[0]
                    // join array into string and set as task
                    search.task = join( subcommands[1:] )
                case "subtask":
                    // delete subtask xyz
                    // set subcommand ( focus ) of search
                    search.subA = subcommands[0]
                    // join array into string and set as subtask
                    search.subtask = join( subcommands[1:] )
                default:
                    // delete board xyz
                    // delete xyz
                    // set subcommand ( focus ) of search
                    search.subA = "board"
                    // if keyword board is given, join without keyword
                    if subcommands[0] == "board" && subcomlen > 1 {
                        search.board = join( subcommands[1:] )
                    } else {    // otherwise, join everything
                        search.board = join( subcommands[:] )
                    }
            }
        case "make":
            //      -> make <board>
            //      -> make board <b>
            //      -> make <board> task <t>
            //      -> make board <b> task <t>
            //      -> make task <t> subtask <s>
            // if subcommand length is 1 or 2, make a board
            // if subcommand length is 3 or 4, make a task/subtask

            // if user enters JUST make, set fail and return
            if subcomlen == 0 {
                search.command = "fail"
                return search
            }

            var x int = -1  // keep track of position of subcom

            if subcommands[0] == "task" {
                // make task xyz subtask ijk
                // iterate over subcommands ( not counting task )
                // looking for the subtask keyword
                for ind, _ := range subcommands[1:] {
                    // set x equal to index when subtask is found
                    if subcommands[ind] == "subtask" {
                        x = ind
                    }
                }
                // the absolute minimum index for subtask keyword is 2
                // if its found before 2, thats a problem
                // set fail and return
                if x < 2 {
                    search.command = "fail"
                    return search
                } else {        // valid keyword position
                    // set secondary focus to subtask
                    search.subB = "subtask"
                    // set main focus to parent task
                    search.subA = "task"
                    // join everything between keywords, set to task
                    search.task = join( subcommands[1:x] )
                    // join everything after subtask, set to subtask
                    search.subtask = join( subcommands[x+1:] )
                }
            } else {
                // make xyz
                // make board xyz
                // make xyz task ijk
                // make board xyz task ijk

                // declare string to hold subsubcommand
                var subsubcom []string
                
                // determine what subsubcommand should include
                if subcommands[0] == "board" {
                    subsubcom = subcommands[1:]
                } else {
                    subsubcom = subcommands[:]
                }
                // iterate over subsubcommand looking for task keyword
                for ind, _ := range subsubcom {
                    // set x equal to index when task is found
                    if subsubcom[ind] == "task" {
                        x = ind
                    }
                }

                // if task was never found, making a board
                if x == -1 {
                    // set focus to board
                    search.subA = "board"
                    // join everything and set to board
                    search.board = join( subsubcom )
                // task was found in a valid position
                } else if x > 0 {
                    // set secondary focus to task
                    search.subB = "task"
                    // set main focus to parent board
                    search.subA = "board"
                    // join everything between keywords, set to board
                    search.board = join( subsubcom[:x] )
                    // join everything after task, set to task
                    search.task = join( subsubcom[x+1:] )
                // task was not in valid position
                } else {
                    // set fail and return
                    return search
                }
            }
        case "mark":
            // NYI  -> mark <board>
            // NYI  -> mark board <b>
            //      -> mark task <t> fill
            //      -> mark subtask <s> fill

            // if user enters JUST mark, set fail and return
            if subcomlen == 0 {
                search.command = "fail"
                return search
            }

            switch subcommands[0] {
                case "task":
                    // mark task xyz #
                    if subcomlen < 3 {
                        // absolute minimum length possible is 3
                        // task, name, and fill
                        search.command = "fail"
                        return search
                    }
                    // set focus to task
                    // set fill to absolute last word in list
                    search.subA = subcommands[0]
                    search.subB = subcommands[subcomlen - 1]
                    // join all the other words, set to task
                    search.task = join( subcommands[1:subcomlen - 1] )
                case "subtask":
                    // mark subtask xyz #
                    if subcomlen < 3 {
                        // absolute minimum length possible is 3
                        // subtask, name, and fill
                        search.command = "fail"
                        return search
                    }
                    // set focus to subtask
                    // set fill to absolute last word in list
                    search.subA = subcommands[0]
                    search.subB = subcommands[subcomlen - 1]
                    // join all the other words, set to subtask
                    search.subtask = join( subcommands[1:subcomlen - 1] )
                default:
                    // delete board xyz
                    // delete xyz
                    // set focus to board
                    // set fill to 100
                    search.subA = "board"
                    search.subB = "100"

                    // if keyword board is given, join without keyword
                    if subcommands[0] == "board" && subcomlen > 1 {
                        search.board = join( subcommands[1:] )
                    } else {    // otherwise, join everything
                        search.board = join( subcommands[:] )
                    }
            }
        case "show":
            // similar call signature to delete()
            //      -> show
            //      -> show all
            //      -> show <board>
            //      -> show board <b>
            //      -> show task <t>
            //      -> show subtask <s>
            // subcom lengths of 1, 2, and 3 are all possible
            if subcomlen == 0 {
                // show
                // set subcommand to all
                search.subA = "all"
                return search
            } else {
                switch subcommands[0] {
                    case "boards", "tasks":
                        // show boards
                        // show tasks
                        // iterate through just the names of the
                        // boards/tasks
                        search.subA = subcommands[0]
                    case "task":
                        // show task xyz
                        // set focus to task
                        // join all the other words, set to task
                        search.subA = subcommands[0]
                        search.task = join( subcommands[1:] )
                    case "subtask":
                        // show subtask xyz
                        // set focus to subtask
                        // join all the other words, set to subtask
                        search.subA = subcommands[0]
                        search.subtask = join( subcommands[1:] )
                    case "all":
                        // show all
                        // show <board name beginning with all>
                        if subcomlen == 1 {
                            // if there's only one subcommand, it's
                            // show all...
                            // set subcommand to all
                            search.subA = subcommands[0]
                        } else {    // otherwise, it must be a board
                            // set focus to board
                            // join all the other words, set to board
                            search.subA = "board"
                            search.board = join( subcommands[:] )
                        }
                    default:
                        // show board xyz
                        // show xyz
                        // set focus to board
                        search.subA = "board"
                        // if keyword board is given, join without keyword
                        if subcommands[0] == "board" && subcomlen > 1 {
                            search.board = join( subcommands[1:] )
                        } else {    // otherwise, join everything
                            search.board = join( subcommands[:] )
                        }
                }
            }
        default:    // command not recognized
            // print error, call help about it
            print( "Command /", comm, "/ not recognized" )
            commFailure := []string{ "help" }
            help( commFailure )
    }

    return search

}


// -----------------------------------------------------------------------------
//   SHOW( search query )
//
//   Revn: 04-11-2022
//   Args: search - query
//                  object that summarizes what needs showing
//   Func: display stored information to the user
//   Meth: iterate through everything, look for matching info, print
//   TODO: stylize print in some way
//         add debug flag
// -----------------------------------------------------------------------------
//   CHANGE LOG
// - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
//   02-24-2022: wrote bones
//               wrote /all/ portion
//   02-25-2022: implemented /board/, /task/, and /subtask/ portions
//               added calls to help() on error
//               changed input argument from type query to type
//                  []string, so it can call sanitize()
//   02-27-2022: commented /all/ portion
//   03-02-2022: commented /board/, /task/, and /subtask/ portions
//   03-16-2022: updated to switch/case system to work with the new
//                  version of parse()
//   03-27-2022: removed 1.x version
//   04-11-2022: changed task fill in to use calc() in /all/, /board/,
//                  and /task/
//
// -----------------------------------------------------------------------------
func show( search query ) {

    // check to see if query returned empty
    // not entirely sure if that's possible, but check for it
    if search == ( query{} ) {
        return
    }

    // ease of use variables from search struct
    subcom := search.subA
    sbd := search.board
    stk := search.task
    sstk := search.subtask

    switch search.subA {
        case "all":
            // iterate over boards' contents
            for boardName, boardContents := range boards {
                // show all -> print board name
                print( boardName )
                // iterate over board's tasks
                for _, showTask := range boardContents {
                    // print indent with no trailing newline
                    // indent has len 4
                    p( "    " )
                    // print task name, right justified (with respect
                    // to the indent)
                    printf( "%-16s", showTask.name )
                    // init string to contain task fill
                    var showTaskFill string
                    // if task fill is in auto mode
                    if showTask.fill == "x" {
                        // call averaging method over task
                        showTaskFill = calc( showTask )
                    } else {    // otherwise, grab fill as usual
                        showTaskFill = showTask.fill
                    }
                    // print task fill, left justified
                    // indent len 4 + 16 = 20
                    printf( "%16s\n", showTaskFill )
                    // iterate over the task's subtasks
                    for _, showSubtask := range showTask.subt {
                        // double indent
                        // double indent has len 8
                        p( "        " )
                        // print subtask name, right justified (with
                        // respect to indent again)
                        printf( "%-16s", showSubtask.name )
                        // print subtask fill, left justified
                        // indent len 8 + 12 = 20
                        printf( "%12s\n", showSubtask.fill )
                    }
                }
            }
        case "board":
            // grab user requested board from boards
            // prs is set true if board exists
            contents, prs := boards[sbd]
            if prs {    // if user specified board exists
                print( sbd )    // print board name
                // iterate over tasks in board
                for _, showTask := range contents {
                    // print indent with no trailing newline
                    // indent has len 4
                    p( "    " )
                    // print task name, right justified (with respect
                    // to the indent)
                    printf( "%-16s", showTask.name )
                    // init string to contain task fill
                    var showTaskFill string
                    // if task fill is in auto mode
                    if showTask.fill == "x" {
                        // call averaging method over task
                        showTaskFill = calc( showTask )
                    } else {    // otherwise, grab fill as usual
                        showTaskFill = showTask.fill
                    }
                    // print task fill, left justified
                    // indent len 4 + 16 = 20
                    printf( "%16s\n", showTaskFill )
                    for _, showSubtask := range showTask.subt {
                        // double indent
                        // double indent has len 8
                        p( "        " )
                        // print subtask name, right justified (with
                        // respect to indent again)
                        printf( "%-12s", showSubtask.name )
                        // print subtask fill, left justified
                        // indent len 8 + 12 = 20
                        printf( "%16s\n", showSubtask.fill )
                    }
                }
            } else {    // if board DNE
                // print error, call help about it
                print( "Board /" + sbd + "/ does not exist" )
                Help( "show" )
                return
            }
        case "boards":
            // literally just print the names of all the boards
            // just a compact version of "show all"
            for showBoard, _ := range boards {
                print( showBoard )
            }
        case "task":
            // can check if a board is present with _, prs :=
            // but this is a flag to be set if a task exists
            var exists bool

            // iterate over all boards looking for named task
            for showBoard, contents := range boards {
                // iterate through all tasks in boards
                for _, showTask := range contents {
                    // look for matching task name
                    if showTask.name == stk {
                        exists = true           // set exists true
                        // print containing board name
                        print( showBoard )
                        // print task name with same justification
                        printf( "%-16s", showTask.name )
                        // init string to contain task fill
                        var showTaskFill string
                        // if task fill is in auto mode
                        if showTask.fill == "x" {
                            // call averaging method over task
                            showTaskFill = calc( showTask )
                        } else {    // otherwise, grab fill as usual
                            showTaskFill = showTask.fill
                        }
                        // print task fill with right justification
                        printf( "%20s\n", showTaskFill )
                        // iterate over contained subtasks
                        for _, showSubtask := range showTask.subt {
                            // print indent
                            p( "    " )
                            // print subtask name left justified
                            printf( "%-16s", showSubtask.name )
                            // print subtask fill right justified
                            printf( "%16s\n", showSubtask.fill )
                        }
                    }
                }
            }
            if !exists {        // if task not found
                                // print error
                                // call help about it
                print( "Task /" + stk + "/ does not exist" )
                Help( "show" )
                return
            }
        case "tasks":
            // compact version of showing all boards and their tasks
            // no subtasks tho, thats just show all
            // doesn't print any fills or anything like that, simple
            for showBoard, contents := range boards {
                print( showBoard )
                for _, showTask := range contents {
                    print( "\t" + showTask.name )
                }
            }
        case "subtask":
            // can check if a board is present with _, prs :=
            // but this is a flag to be set if a task exists
            var exists bool

            // iterate over all boards looking for named task
            for showBoard, contents := range boards {
                // iterate through all tasks in boards
                for _, showTask := range contents {
                    // iterate over contained subtasks
                    for _, showSubtask := range showTask.subt {
                        // look for matching subtask name
                        if showSubtask.name == sstk {
                            exists = true       // set exists true
                            // print containing board name
                            print( showBoard )
                            // print containing task name
                            print( showTask.name )
                            // print subtask name left justified
                            printf( "%-16s", showSubtask.name )
                            // print subtask fill right justified
                            printf( "%20s\n", showSubtask.fill )
                        }
                    }
                }
            }

            if !exists {        // if subtask not found
                                // print error
                                // call help about it
                print( "Subtask /" + sstk + "/ does not exist" )
                Help( "show" )
                return
            }
        default:    // unrecognized subcommand
            // print error, call help about it
            print( "Error in show()" )
            print( "Subcommand /" + subcom + "/ does not exist" )
            Help( "show" )
            return
    }
}

var flag bool = false       // debug flag
var changed bool = false    // keep track of if file is changed


// struct used for defining search parameters
type query struct {
    command string      // command passed, or fail
    board string        // requested board
    task string         // requested task
    subtask string      // requested subtaks
    subA string         // subcom A, usually object being acted upon
    subB string         // sometimes, two are needed
}


/*
type subsubtask struct {
    name string         // subsubtasks are named
    fill string         // subsubtasks contain progress
    //note string         // subsubtasks can contain files
}
*/

// struct used for defining what constitutes a 'subtask'
type subtask struct {
    name string         // subtasks are named
    fill string         // subtasks contain progress
    //ssbt []subsubtask     // subtasks can contain subsubtasks
    //note string         // subtasks can contain files
}


// struct used for defining what constitutes a 'task'
type task struct {
    name string         // tasks are named
    fill string         // tasks contain progress
    subt []subtask      // tasks can contain subtasks
    //note string         // tasks can contain files
}


// boards is a hashmap      map[      ]
// that relates strings         string
// to an array of tasks                []task
var boards = make( map[string][]task )
var tasks []task        // array of tasks referenced above
var subtasks []subtask  // array of subtasks belonging to tasks
var home string         // declare string to keep track of home dir

// MAIN
func main() {

    catch()     // call function to catch keyboard interrupts

    fileIn()    // open data file, process into structure

    for {       // infinite loop, where everything happens

        // declare array of strings to contain user input
        var ui []string
        // and quick flag to keep track of operation mode:
        // cmd line args or inline?
        var cmd bool = false

        // check to see if there was a given argument
        if len( os.Args[1:] ) > 0 {     // if yes
            ui = os.Args[1:]        // declare ui as list of args

            // XXX
            if flag {   // if debug set, print arguments list
                print( ui )
            }

            cmd = true              // set cmd line args flag
        } else {                        // if no
            userIn := input( "-> " )        // take user input
            ui = split( userIn, " " )       // split over spaces
        }

        command := ui[0]                // first entry is command
        // convention: send entire array to subfunction, do not strip
        // any subcommands off

        // process the user input
        var search query = parse( ui )

        // switch case over command
        // recognize user input, call appropriate command
        switch command {
            case "delete":
                del( search )
            case "help":
                help( ui )
            case "make":
                // make translated into Russian, because make is a
                // GoLang keyword, and Norwegian word was less cool
                // pronounced like 'zdye-lie'
                сделай( search )
            case "mark":
                mark( search )
            case "save":
                fileOut()
                changed = false
                print( "Updates saved" )
            case "show":
                show( search )
            case "exit", "quit":

                // XXX
                if flag {   // if debug set, print if file was changed
                    print( changed )
                }

                // if file was changed, rewrite to file
                if changed {
                    fileOut()
                }

                exit( 0 )    // exeunt
            case "set":
            // XXX
                // turn debug flag on and off
                // conveniently only recognizes first and last word of
                // command; "set true" will work, as will "set flag
                // false" and "set sprite coke pepsi"
                flag = ui[len( ui ) - 1] == "true"
            default:            // if the input is unrecognized
                continue        // try again
        }

        // if cmd line arg mode, only run once, break infinite loop
        if cmd {
            // if file was changed, rewrite to file
            if changed {
                fileOut()
            }
            break
        }
    }
}

