// Eye - A simple file change command executioner
package main

import (
	"flag"
	"fmt"
	"os"
	"regexp"
	"time"

	"github.com/callerobertsson/eye/watcher"
)

// Command line flags and options
var (
	helpFlag        bool
	patternOption   string
	commandOption   string
	noRecursionFlag bool
)

// Initialize flags and options
func init() {
	flag.BoolVar(&helpFlag, "h", false, "Show usage information and exit")
	flag.BoolVar(&helpFlag, "help", false, "Show usage information and exit")
	flag.BoolVar(&noRecursionFlag, "R", false, "Do *not* recurse sub directories")
	flag.StringVar(&patternOption, "p", "", "Matching files regex pattern")
	flag.StringVar(&commandOption, "c", "", "Command to execute on changes")
	flag.Parse()

	// Print usage and exit if help flag present
	if helpFlag {
		usage()
		os.Exit(0)
	}

	// Pattern and Command are mandatory
	if patternOption == "" || commandOption == "" {
		usage()
		os.Exit(1)
	}
}

// Main function
func main() {
	// Compile the input pattern
	p, err := regexp.Compile(patternOption)
	if err != nil {
		fmt.Printf("-p %q is not a valid regular expression\n", patternOption)
		os.Exit(1)
	}

	// Create the command function
	c := func() {
		watcher.RunSystemCommand(commandOption)
	}

	// Create status reporting channel
	ss := make(chan watcher.Status)

	// Create Watcher
	w := watcher.New(p, c, !noRecursionFlag, 1*time.Second, ss)

	// Watch
	go w.Watch()

	// Read statuses until end of time
	for {
		s := <-ss
		printStatus(s)
	}
}

// Print status to the console
func printStatus(s watcher.Status) {
	switch s.Type {
	case watcher.StatusModified:
		fmt.Printf("\nMODIFIED: The file %q was changed\n", s.File)
	case watcher.StatusAdded:
		fmt.Printf("\nADDED: %v\n", s.Message)
	case watcher.StatusDeleted:
		fmt.Printf("\nDELETED: %v\n", s.Message)
	default:
		fmt.Printf("%v", s.Message)
	}
}

// Print usage information
func usage() {
	fmt.Printf("Usage:\n\teye [-R] -p <PATTERN> -c <COMMAND>\n\n")
	flag.PrintDefaults()
}
