// =============================================================================
// Auth: Alex Celani
// File: tdm.go
// Revn: 03-11-2022  1.0
// 
// Func: display and manage progress of a litany of items to be done,
//       with an organization scheme similar to Trello. It's CLI
//       Trello
//
// TODO: add debug statements
//
//       literally how do I handle parsing multi-word input to the
//          program???
//
//       add second object to hold command information
//       make subsubtasks for greater depth
// M        or rooms to contain boards
//       handle command line args
// M     remember last args
//       research git branches
//       handle ^C gracefully
//       move most (read: all) helper functions to external files
// M     add priority to tasks
//       attach text files???
//       fix data storage issue
//       contemplate adding XOR of data in file, altho who cares
//       contemplate adding battery of tests
//       contemplate moving functions to external file
//       side-by-side printing
//       colored printing?
//       keep track of finished items
//       periodic tasks? something to do with the time package
//       convert file delimiters to (customizable?) variables
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
//
// =============================================================================

package main

import ( 
    "fmt"           // used for Println
    "strings"       // used for Split
    "strconv"       // used for Atoi
    "os"            // used for Exit, Stdin
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
var read   = ioutil.ReadFile
var write  = ioutil.WriteFile


// wrapper function to print input prompt, take user input
func input( ps1 string ) string {
    p( ps1 )                                // print prompt
    // create scanner on stdin
    scanner := bufio.NewScanner( os.Stdin )
    scanner.Scan()                          // read input
    return scanner.Text()                   // return text
}


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

            // convert fill back to ascii
            // no error, because I guess all ints can be converted to
            // ascii, but not all ascii can be an int?
            // TODO think about keeping fill as ascii
            // maybe check the conversion to see if it worked (i.e.
            // user entered a number), but the ints are never
            // _actually_ used
            fill := itoa( tcontents.fill )

            data = data + fill              // add task fill
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


                    data = data + scontents.name

                    // XXX
                    if flag {   // and after
                        print( "post data: ", data )
                    }

                    data = data + ">"       // XXX subtask fill delim

                    // convert fill back to ascii
                    // TODO as above
                    fill = itoa( scontents.fill )

                    data = data + fill      // add subtask fill

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

    // cast writestring to byte array, write to file named "data"
    // if a file doesn't exist, make with 0x644 permissions
    // take error as err, if fails
    err := write( "data", []byte( data ), 0644 )

    if flag {   // if debug set, print data written
        print( data )
    }

    if err != nil {     // if write() failed, tell user
        print( "Error in fileOut(): file not properly written" )
        print( err )
    }

}

// warpper function to contain all the file reading
// literally yanked from main() and pasted here
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

    file, err := read( "data" )         // open file, failure set err

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
                    // convert fill from string to int, discard err
                    fil, _ := atoi( sContents[1] )
                    // create subtask object with this extracted info
                    stk := subtask { name: sContents[0], fill: fil }
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
                tContents[1] = "0"
            }
            tk.name = tContents[0]          // set name of task
            fil, _ := atoi( tContents[1] )  // convert fill to int
            tk.fill = fil                   // set fill of task
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
//   help( commands []string )
//
//   Revn: 03-10-2022
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

    print()     // print extra line to avoid cluttered look

}


// -----------------------------------------------------------------------------
//   delete( search query )
//
//   Revn: 03-08-2022
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
//
// -----------------------------------------------------------------------------
func del( search query ) {

    // determine what needs searching

    // check to see if query returned empty
    // not entirely sure if that's possible, but check for it
    if search == ( query{} ) {
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
        print( "Error in delete(): board does not exist or not specified" )
    }
}


// -----------------------------------------------------------------------------
//   find( search query ) ( string, int, int )
//
//   Revn: 03-03-2022
//   Args: search  - query
//                   object that outlines what needs finding
//   Retn: x, y, z - string, int, int
//                   "coordinates" that outline where sought-after
//                      objects are
//   Func: iterate over all everything until exact positions of
//         sought-after object is found
//   TODO: incorporate with sanitize?
//         error messages
// -----------------------------------------------------------------------------
//   CHANGE LOG
// - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
//   03-01-2022: essentially stole from show()
//   03-03-2022: commented
//
// -----------------------------------------------------------------------------
func find( search query ) ( string, int, int ) {
    
    var rboard string       // init string as empty
    var rtask int = -1      // init ints as -1, the forbidden index
    var rsubt int = -1

    //subcom := search.subcom
    brd := search.board
    tsk := search.task
    sbt := search.subtask

    // if user is looking for a board
    if brd != "" {
        _, prs := boards[brd]   // quick check to see if board exists
        if prs == true {        // if exists...
            rboard = brd        // set returnable board name
        //} else {    // requested board doesn't exist
            // print error about it
        //    print( "Board /" + brd + "/ does not exist" )
        }
    } else if tsk != "" {       // if user is looking for task
        // iterate over all boards
        for bd, contents := range boards {
            // iterate over tasks in given board
            for tindex, ts := range contents {
                // if found task matches sought task
                if ts.name == tsk {
                    // keep track of board and task index
                    rboard = bd
                    rtask = tindex
                }
            }
        }
    } else if sbt != "" {       // if user is looking for subtask
        // iterate over all boards
        for bd, contents := range boards {
            // iterate over tasks in given board
            for tindex, ts := range contents {
                // iterate over subtask in given board
                for sindex, st := range ts.subt {
                    // if found subtask matches sought subtask
                    if st.name == sbt {
                        // keep track of board, task index, and
                        // subtask index
                        rboard = bd
                        rtask = tindex
                        rsubt = sindex
                    }
                }
            }
        }
    } else {    // should be unreachable
                // implies that user wasnt looking for board, task, or
                // subtask
        //print( "what is happen" )
    }

    // return returnable 3-tuple
    return rboard, rtask, rsubt

}


// -----------------------------------------------------------------------------
//   сделай( search query )
//
//   Revn: 03-09-2022
//   Args: search - query
//                  object that summarizes what needs making
//   Func: create new boards, tasks, and subtasks
//   Meth: find parent object, create new child object, attach
//   TODO: write second draft
//         add debug
// -----------------------------------------------------------------------------
//   CHANGE LOG
// - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
//   03-07-2022: wrote first draft
//   03-10-2022: commented
//
// -----------------------------------------------------------------------------
func сделай( search query ) {

    if flag {   // if debug set, print contents of argument
        print( search )
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
                makeFailure := []string{ "help", "make" }
                help( makeFailure )
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
                makeFailure := []string{ "help", "make" }
                help( makeFailure )
                return
            } else {                    // board is specified
                exists = true           // signal parent did exist
                // create task with user-given name
                tk := task{ name: search.task }
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
            makeFailure := []string{ "help", "make" }
            help( makeFailure )
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
        makeFailure := []string{ "help", "make" }
        help( makeFailure )
    }

}

// -----------------------------------------------------------------------------
//   mark( search query )
//
//   Revn: 03-11-2022
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
//
// -----------------------------------------------------------------------------
func mark( search query ) {

    // check to see if query returned empty
    // not entirely sure if that's possible, but check for it
    if search == ( query{} ) {
        return
    }

    // determine "coordinates" of object that needs finding
    // x -> ( string ) name of board, shouldn't ever be blank
    // y -> ( int ) index of task array for a given board
    // z -> ( int ) index of subtask array for a given task
    x, y, z := find( search )
    
    // index of last item           len( commands ) - 1
    // last item in array  commands[len( commands ) - 1]
    // convert string to int
    //               atoi( commands[len( commands ) - 1] )
    // markFill -> converted value
    // err      -> contains value of atoi() error
    //             would happen if passed value isn't ascii
    markFill, err := atoi( search.subcom )
    
    // on error, print failure, call help about it, return
    if err != nil {
        print( "Fill value /" + search.subcom + "/ is not a number" )
        markFailure := []string{ "help", "mark" }
        help( markFailure )
        return
    }

    // round fill down to 100 max, specifically allow negative fill
    if markFill > 100 {
        markFill = 100
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
                boards[x][y].fill = markFill
                changed = true  // signal that file needs rewriting
            } else {       // yes subtask
                // boards                    -> map[string][]task
                // boards[x]                 -> []task
                // boards[x][y]              -> task
                // boards[x][y].subt         -> []subtask
                // boards[x][y].subt[z]      -> subtask
                // boards[x][y].subt[z].fill -> subtask.fill
                boards[x][y].subt[z].fill = markFill
                changed = true  // signal that file needs rewriting
            }
        }
    } else {    // reachable if the user gives input that DNE
        // notify user requested object did not exist
        p( "Error in mark(): " )
        if search.board != "" {
            p( "board /", search.board )
        } else if search.task != "" {
            p( "task /", search.task )
        } else if search.subtask != "" {
            p( "subtask /", search.subtask )
        } else {
            p( "object /", search.board )
        }
        print( "/ does not exist" )
    }
}


// -----------------------------------------------------------------------------
//   parse( commands []string )
//
//   Revn: 03-10-2022 
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
//   TODO: comments for
//              /make/
//              /mark/
//
//         contemplate help case
//         move to external file?
// -----------------------------------------------------------------------------
//   CHANGE LOG
// - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
//   03-06-2022: planned
//               wrote everything but help case
//               tested pass cases
//   03-10-2022: commented /delete/ portion
//               commented /show/ portion
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
    //search.command = comm

    // big switch statement over command
    switch comm {
        // these don't really need processing, break out immediately
        case "set", "quit", "exit":
            return search
        case "delete":
            // there are only so many types of commands that work here
            //      -> delete <board>
            //      -> delete board <board>
            //      -> delete task <task>
            //      -> delete subtask <subtask>
            // i.e. subcommand can have either length 1 or 2
            switch subcomlen {
                case 1:
                    // -> delete <board>
                    // requested item is subcommands[0], so set
                    search.board = subcommands[0]
                case 2:
                    // -> delete board|task|subtask <object>
                    // can be any of three, so switch over subcom[0]
                    // kind of wish I could do like a reverse set
                    // like subcommands[1] = board or task or subtask
                    switch subcommands[0] {
                        case "board":
                            search.board = subcommands[1]
                        case "task":
                            search.task = subcommands[1]
                        case "subtask":
                            search.subtask = subcommands[1]
                        default:    // user didn't input valid object
                            // print error, call help about it
                            print( "Incorrect object specified" )
                            p( "Object /", subcommands[0] )
                            print( "/ is not valid" )
                            delFailure := []string{ "help", "delete" }
                            help( delFailure )
                    }
                default:    // user had the wrong number of args
                    // print error, call help about it
                    print( "Incorrect number of arguments" )
                    delFailure := []string{ "help", "delete" }
                    help( delFailure )
            }
        case "help":
//            help( subcommands )
        case "make":
            //      -> make <board>
            //      -> make board <b>
            //      -> make <board> task <t>
            //      -> make board <b> task <t>
            //      -> make task <t> subtask <s>
            // if subcommand length is 1 or 2, make a board
            // if subcommand length is 3 or 4, make a task/subtask
            switch subcomlen {
                case 1:
                    // -> make <board>
                    // only subcommand is argument, so set
                    search.board = subcommands[0]
                case 2:
                    // -> make board <b>
                    // double triple check that user asked for board,
                    // then set second subcommand
                    if subcommands[0] == "board" {
                        search.board = subcommands[1]
                    } else {    // user didn't seem to want to make a
                                // board
                        // print error, call help about it
                        print( "Incorrect argument type:" )
                        print( "Object /", subcommands[0], "/ does not match the following types:" )
                        print( "/board/" )
                        makeFailure := []string{ "help", "make" }
                        help( makeFailure )
                    }
                case 3:
                    // -> make <board> task <t>
                    // if task is actually task, set <board> to board
                    // and <t> to task
                    if subcommands[1] == "task" {
                        search.board = subcommands[0]
                        search.task = subcommands[2]
                    } else {    // user didn't want to make a task
                        // print error, call help about it
                        print( "Incorrect argument type:" )
                        print( "Object /", subcommands[1], "/ does not match the following types:" )
                        makeFailure := []string{ "help", "make" }
                        help( makeFailure )
                    }
                case 4:
                    // -> make board <b> task <t>
                    // -> make task <t> subtask <s>
                    // look for a signature that matches the two above
                    // set appropriate fields
                    if subcommands[0] == "board" && subcommands[2] == "task" {
                        search.board = subcommands[1]
                        search.task = subcommands[3]
                    } else if subcommands[0] == "task" && subcommands[2] == "subtask" {
                        search.task = subcommands[1]
                        search.subtask = subcommands[3]
                    } else {    // user input didn't conform to scope
                        // print error, call help about it
                        print( "Unknown call signature" )
                        makeFailure := []string{ "help", "make" }
                        help( makeFailure )
                    }
                default:        // subcommand length > 4
                    // print error, call help about it
                    print( "Incorrect number of arguments" )
                    makeFailure := []string{ "help", "make" }
                    help( makeFailure )
            }
        case "mark":
            // NYI  -> mark <board> [100]
            // NYI  -> mark board <b> [100]
            //      -> mark task <t> fill
            //      -> mark subtask <s> fill
            switch subcomlen {
                case 1, 2:
                    // -> mark <board>
                    // -> mark board <b>
                    // -> mark <board> 100
                    // marking boards isn't support yet
//                    if subcomlen == 2 {
//                        if subcommands[0] == "board" {
//                            search.board = subcommands[1]
//                        } else {
//                            // fail
//                            // unrecognized argument object
//                            // /[subcommands[0]]/
//                            // help mark
//                            print( "fail h" )
//                        }
//                    } else {
//                        search.board = subcommands[0]
//                    }
//                    search.subcom = "100"
                    print( "Incorrect number of arguments" )
                    markFailure := []string{ "help", "mark" }
                    help( markFailure )
                case 3:
                    // -> mark task <t> fill
                    // -> mark subtask <s> fill
                    // -> mark board <b> 100    NYI
                    // pretty simple, double check first subcommand to
                    // see what object is being marked, set name as
                    // needed, set fill to subcommand
                    switch subcommands[0] {
                        case "task":
                            search.task = subcommands[1]
                        case "subtask":
                            search.subtask = subcommands[1]
                        default:    // user didnt want task or subtask
                            // print error, call help about it
                            p( "Object /", subcommands[0] )
                            print( "/ not recognized" )
                            markFailure := []string{ "help", "mark" }
                            help( markFailure )
                    }
                    search.subcom = subcommands[2]
                default:        // subcommand length != 3
                    // print error, call help about it
                    print( "Incorrect number of arguments" )
                    markFailure := []string{ "help", "mark" }
                    help( markFailure )
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
            switch subcomlen {
                case 0:
                    // show
                    // equivalent to show all, so set to all
                    search.subcom = "all"
                case 1:
                    // show all
                    // show <board>
                    // if user asks for all, give all
                    if subcommands[0] == "all" {
                        search.subcom = "all"
                    } else {    // otherwise, user is asking for board
                        search.board = subcommands[0]
                    }
                case 2:
                    // show [board|task|subtask] <b|t|s>
                    // switch to see what was specified, set the
                    // appropriate field in search, easy peasy
                    switch subcommands[0] {
                        case "board":
                            search.board = subcommands[1]
                        case "task":
                            search.task = subcommands[1]
                        case "subtask":
                            search.subtask = subcommands[1]
                        default:    // user didn't input valid object
                            // print error, call help about it
                            print( "Incorrect object specified" )
                            p( "Object /", subcommands[0] )
                            print( "/ is not valid" )
                            showFailure := []string{ "help", "show" }
                            help( showFailure )
                    }
                default:    // user had the wrong number of args
                    // print error, call help about it
                    print( "Incorrect number of arguments" )
                    showFailure := []string{ "help", "show" }
                    help( showFailure )
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
//   show( search query )
//
//   Revn: 03-09-2022
//   Args: search - query
//                  object that summarizes what needs showing
//   Func: display stored information to the user
//   Meth: iterate through everything, look for mathcing info, print
//   TODO: stylize print in some way
//         add debug flag
//         add call to find()
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
//
// -----------------------------------------------------------------------------
func show( search query ) {

    // check to see if query returned empty
    // not entirely sure if that's possible, but check for it
    if search == ( query{} ) {
        return
    }

    // ease of use variables from search struct
    subcom := search.subcom
    sbd := search.board
    stk := search.task
    sstk := search.subtask

    if subcom != "" {    // show all
        // might not actually contain "all", double check
        if subcom == "all" {
            // iterate over boards' contents
            for showBoard, contents := range boards {
                // show all -> print board name
                print( showBoard )
                // iterate over board's tasks
                for _, showTask := range contents {
                    // print indent with no trailing newline
                    // indent has len 4
                    p( "    " )
                    // print task name, right justified (with respect
                    // to the indent)
                    printf( "%-16s", showTask.name )
                    // print task fill, left justified
                    // indent len 4 + 16 = 20
                    printf( "%16d\n", showTask.fill )
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
                        printf( "%12d\n", showSubtask.fill )
                    }
                }
            }
        } else {    // if command isn't "all"
                    // print error, call help about it
            print( "Error in show()" )
            print( "Subcommand /" + subcom + "/ does not exist" )
            showFailure := []string{ "help", "show" }
            help( showFailure )
            return
        }
    } else if sbd != "" {   // show board <>
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
                // print task fill, left justified
                // indent len 4 + 16 = 20
                printf( "%16d\n", showTask.fill )
                for _, showSubtask := range showTask.subt {
                    // double indent
                    // double indent has len 8
                    p( "        " )
                    // print subtask name, right justified (with
                    // respect to indent again)
                    printf( "%-12s", showSubtask.name )
                    // print subtask fill, left justified
                    // indent len 8 + 12 = 20
                    printf( "%16d\n", showSubtask.fill )
                }
            }
        } else {    // if board DNE
                    // print error, call help about it
            print( "Board /" + sbd + "/ does not exist" )
            showFailure := []string{ "help", "show" }
            help( showFailure )
            return
        }
    } else if stk != "" {   // show task <>
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
                    // print task fill with right justification
                    printf( "%20d\n", showTask.fill )
                    // iterate over contained subtasks
                    for _, showSubtask := range showTask.subt {
                        // print indent
                        p( "    " )
                        // print subtask name left justified
                        printf( "%-16s", showSubtask.name )
                        // print subtask fill right justified
                        printf( "%16d\n", showSubtask.fill )
                    }
                }
            }
        }
        if !exists {        // if task not found
                            // print error
                            // call help about it
            print( "Task /" + stk + "/ does not exist" )
            showFailure := []string{ "help", "show" }
            help( showFailure )
            return
        }
    } else if sstk != "" {  // show subtask <>
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
                        printf( "%20d\n", showSubtask.fill )
                    }
                }
            }
        }

        if !exists {        // if subtask not found
                            // print error
                            // call help about it
            print( "Subtask /" + sstk + "/ does not exist" )
            showFailure := []string{ "help", "show" }
            help( showFailure )
            return
        }

    } else {                // should be undefined behavior
        // print contents of the query (human-readable)
        print( "Error with input query!" )
        print( "  query.board      -> " + sbd )
        print( "  query.task       -> " + stk )
        print( "  query.subtask    -> " + sstk )
        print( "  query.subcommand -> " + subcom )

        // call help about it
        showFailure := []string{ "help", "show" }
        help( showFailure )
        return
    }
}



var flag bool = false       // debug flag
var changed bool = false


// struct used for defining search parameters
type query struct {         // relatively clean way to specify...
    board   string          // a board
    task    string          // a task
    subtask string          // a subtask
    subcom  string          // and\or a subcommand
    command string          // and finally the command
}


/*
type subsubtask struct {
    name string         // subsubtasks are named
    fill int            // subsubtasks contain progress
    //note string         // subsubtasks can contain files
}
*/

// struct used for defining what constitutes a 'subtask'
type subtask struct {
    name string         // subtasks are named
    fill int            // subtasks contain progress
    //ssbt subsubtask     // subtasks can contain subsubtasks
    //note string         // subtasks can contain files
}


// struct used for defining what constitutes a 'task'
type task struct {
    name string         // tasks are named
    fill int            // tasks contain progress
    subt []subtask      // tasks can contain subtasks
    //note string         // tasks can contain files
}


// boards is a hashmap      map[      ]
// that relates strings         string
// to an array of tasks                []task
var boards = make( map[string][]task )
var tasks []task        // array of tasks referenced above
var subtasks []subtask  // array of subtasks belonging to tasks


func main() {

    fileIn()

    for {       // infinite loop, where everything happens

        userIn := input( "-> " )        // take user input
        ui := split( userIn, " " )      // split over spaces
        command := ui[0]                // first entry is command
        // convention: send entire array to subfunction, do not strip
        // any subcommands off

        // switch case over command
        // recognize user input, call appropriate command

        var search query = parse( ui )

        
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
            case "show":
                show( search )
            case "exit", "quit":

                if flag {   // if debug set, print if file was changed
                    print( changed )
                }

                // if file was changed, rewrite to file
                if changed {
                    fileOut()
                }

                os.Exit( 0 )    // exeunt
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
        
    }
}

