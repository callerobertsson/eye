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
	helpFlag      bool
	recursiveFlag bool
	patternOption string
	commandOption string
)

// Initialize flags and options
func init() {
	flag.BoolVar(&helpFlag, "h", false, "Show usage information and exit")
	flag.BoolVar(&helpFlag, "help", false, "Show usage information and exit")
	flag.BoolVar(&recursiveFlag, "r", false, "Recurse sub directories")
	flag.StringVar(&patternOption, "p", "", "Matching files regex pattern")
	flag.StringVar(&commandOption, "c", "", "Command to execute on changes")
	flag.Parse()

	if helpFlag {
		usage()
		os.Exit(0)
	}

	if patternOption == "" || commandOption == "" {
		usage()
		os.Exit(1)
	}
}

// Main function
func main() {
	// Compile the input pattern
	p, err := regexp.Compile(patternOption) // TODO: Handle errors in patterns
	if err != nil {
		fmt.Printf("-p %q is not a valid regular expression\n", patternOption)
		os.Exit(1)
	}

	// Create the command function
	c := func() {
		watcher.RunSystemCommand(commandOption)
	}

	// Create Watcher
	ss := make(chan watcher.Status)
	w := watcher.New(p, c, recursiveFlag, 1*time.Second, ss)

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
	fmt.Println("Usage:\n\teye [-r] -p <PATTERN> -c <COMMAND>\n")
	flag.PrintDefaults()
}
