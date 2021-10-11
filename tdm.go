// =============================================================================
// Auth: Alex Celani
// File: tdm.go
// Revn: 10-10-2021  0.5
// 
// Func: display and manage progress of a litany of items to be done,
//       with an organization scheme similar to Trello. It's CLI
//       Trello
//
// TODO: COMMENT
//       file out
//       fix data storage issue
//       contemplate adding XOR of data in file, altho who cares
//       contemplate adding battery of tests
//       contemplate moving functions to external file
//       side-by-side printing
//       colored printing?
//       keep track of finished items
//       periodic tasks? something to do with the time package
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
//*10-10-2021: fixed user input
//             wrote help()
//             wrote second draft of show(), fix immediately
//
// =============================================================================

package main

import ( 
    "fmt"           // used for Println
    "strings"       // used for Split
    "strconv"       // used for Atoi
    "io/ioutil"     // used for ReadFile
    "os"            // used for Exit, Stdin
    "bufio"         // used for NewScanner
)


// aliases for common functions
var p      = fmt.Print
var print  = fmt.Println
var printf = fmt.Printf
var split  = strings.Split
var lower  = strings.ToLower
var atoi   = strconv.Atoi
var exit   = os.Exit
var read   = ioutil.ReadFile


func input( ps1 string ) string {
    p( ps1 )
    scanner := bufio.NewScanner( os.Stdin )
    scanner.Scan()
    return scanner.Text()
}



func choose( commands []string ) {
    
}


func help( c string ) {
    var helpFlag bool = c == "all" || c == ""
    if helpFlag || c == "choose" {
        print( " choose <board>" )
        print( "   \\-- switch focus to selected board" )
    } else if helpFlag || c == "make" {
        print( " make <board = b>" )
        print( "   \\-- create new board, <b>" )
        print( "   \\  <board> <task = t>" )
        print( "   \\-- create new task <t> under board <board>" )
        print( "   \\  <board> <task> <subtask = s>" )
        print( "   \\-- create new subtask <s> under task <board.task>" )    
    } else if helpFlag || c == "mark" {
        print( " mark <task> [<subtask>] <status = 100>" )
        print( "   \\                  \\-- status -> 0 - 100" )
        print( "   \\-- set progress of <task> or <task.subtask> to <status>" )
    } else if helpFlag || c == "move" {
        print( " move <board1> <task> <board2>" )
        print( "   \\-- change scope of <task> from <board1> to <board2>" )    
    } else if helpFlag || c == "show" {
        print( " show [all|<board> [<task> [<subtask>]]]" )
        print( "   \\-- display contents of <board>, <board.task>, or <board.task.subtask>" )
        print( "   \\-- or, with [all] set, display everything" )    
    } else if helpFlag || c == "exit" || c == "quit" {
        print( " exit" )
        print( " quit" )
        print( "   \\-- exit program" )    
    }
}


func mark( commands []string ) {
    
}


func move( commands []string ) {
    
}


func show( commands []string ) {

    //commands = append( commands, "all" )
    switch commands[0] {
        case "all":
            for b, c := range boards {
                print( b )
                for _, task := range c {
                    printf( "%-20s", task.name )
                    printf( "%20d\n", task.fill )
                    for _, subt := range task.subt {
                        printf( "%5" )
                        printf( "%-15s", subt.name )
                        printf( "%20d\n", subt.fill )
                    }
                }
            }
        case "board":
            _, prs := boards[commands[1]]
            if prs {
                for _, task := range boards[commands[1]] {
                    printf( "%-20s", task.name )
                    printf( "%20d\n", task.fill )
                    for _, subt := range task.subt {
                        printf( "%5" )
                        printf( "%-15s", subt.name )
                        printf( "%20d\n", subt.fill )
                    }
                }
            } else {
                print( "board <", commands[1], "> does not exist" )
            }
        case "task":
            var tFlag bool = false
            // make sure commands[1] is in boards
            for _, c := range boards {
                for _, task := range c {
                    if task.name == commands[1] {
                        tFlag = true
                        for _, subt := range task.subt {
                            printf( "%5" )
                            printf( "%-15s", subt.name )
                            printf( "%20d\n", subt.fill )
                        }
                    }
                }
            }
            if !tFlag {
                print( "task <", commands[1], "> does not exist" )
            }
        case "subtask": 
            var sFlag bool = false
            // make sure commands[1] is in boards
            for _, c := range boards {
                for _, task := range c {
                    for _, subt := range task.subt {
                        if subt.name == commands[1] {
                            sFlag = true
                            printf( "%5" )
                            printf( "%-15s", subt.name )
                            printf( "%20d\n", subt.fill )
                        }
                    }
                }
            }
            if !sFlag {
                print( "subtask <", commands[1], "> does not exist" )
            }
        default:
            //var str []string
            //str = append( str, "show" )
            //help( str )
            help( "show" )
    }
}


var flag bool = false

// struct used for defining what constitutes a 'subtask'
type subtask struct {
    name string         // subtasks are named
    fill int            // subtasks contain progress
    //ssbt subsubtask     // subtasks can contain a subsubtask
}


// struct used for defining what constitutes a 'task'
type task struct {
    name string         // tasks are named
    fill int            // tasks contain progress
    subt []subtask      // tasks con contain a subtask
}


// boards is a hashmap      map[      ]
// that relates strings         string
// to an array of tasks                []task
var boards = make( map[string][]task )
var tasks []task    // the array of tasks referenced above
var subtasks []subtask


func main() {
 

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

    file, err := read( "data" )
    if err != nil {
        panic( "file read error" )
    }

    board := split( string( file ), "$" )     // $ -> separates different boards
    
    for a := 0; a < len( board ); a++ {
        bContents := split( board[a], "|" )
        bName := bContents[0]
        bTasks := bContents[1]
        if flag {
            print( "board name -> ", bName )
        }

        tTasks := split( bTasks, ";" )
        var tk task
        for b := 0; b < len( tTasks ); b++ {
            if flag {
                print( "tTasks -> ", tTasks[b] )
            }
            subs := split( tTasks[b], ":" )
            if len( subs ) != 1 {
                sTasks := split( subs[1], "." )
                for c := 0; c < len( sTasks ); c++ {
                    sContents := split( sTasks[c], ">" )
                    fil, _ := atoi( sContents[1] )
                    stk := subtask { name: sContents[0], fill: fil }
                    subtasks = append( subtasks, stk )
                }
                tk.subt = subtasks
            }
            tContents := split( tTasks[b], "," )
            if len( tContents[1] ) == 0 {
                tContents[1] = "0"
            }
            tk.name = tContents[0]
            fil, _ := atoi( tContents[1] )
            tk.fill = fil
            tasks = append( tasks, tk )
            if flag {
                print( "[]task -> ", tasks )
            }
            tk.subt = nil
        }

        boards[bName] = tasks
        tasks = nil
    }

    if flag {
        print( "\n\nboards -> ", boards, "\n\n" )
    }

    for {

        ui := input( "-> " )
        userIn := split( ui, " " )
        command := userIn[0]

        switch command {
            case "choose":
                choose( userIn[1:] )
            case "help":
                help( userIn[1] )
            case "mark":
                mark( userIn[1:] )
            case "move":
                move( userIn[1:] )
            case "show":
                show( userIn[1:] )
            case "exit":
                os.Exit( 0 )
            case "quit":
                os.Exit( 0 )
            default:
                continue
        }
    }

}
