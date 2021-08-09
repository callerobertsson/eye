// Package main implements Eye, a simple file change command executioner
package main

import (
	"flag"
	"fmt"
	"os"
	"os/user"
	"regexp"
	"time"

	"github.com/callerobertsson/eye/params"
	"github.com/callerobertsson/eye/watcher"
)

// Main function
func main() {
	fmt.Printf("Eye - the file change watcher\n")

	// Read params from resource file and command line
	ps, err := initParams()
	if err != nil {
		fmt.Printf("Error in arguments or resource file: %v\n", err)
		os.Exit(1)
	}

	// Print usage and exit if help flag present
	if ps.Help {
		usage()
		os.Exit(0)
	}

	// Compile the input pattern
	p, err := regexp.Compile(ps.Pattern)
	if err != nil {
		fmt.Printf("-p %q is not a valid regular expression\n", ps.Pattern)
		os.Exit(1)
	}

	// Show params
	fmt.Printf("  %-14v %v\n", "Pattern:", ps.Pattern)
	fmt.Printf("  %-14v %v\n", "Command:", ps.Command)
	fmt.Printf("  %-14v %v\n", "Recursive:", ps.Recursive)
	fmt.Printf("  %-14v %v millis\n", "Interval:", ps.IntervalMillis)
	if ps.ResourceFile != "" {
		fmt.Printf("  %-14v %v\n", "Resource file:", ps.ResourceFile)
	}

	// Create the command function
	c := func() {
		err := watcher.RunSystemCommand(ps.Command)
		if err != nil {
			fmt.Println(err)
			os.Exit(2)
		}
	}

	// Create status reporting channel
	ss := make(chan watcher.Status)

	// Create Watcher
	i := time.Duration(ps.IntervalMillis) * time.Millisecond
	w := watcher.New(p, c, ps.Recursive, i, ss)

	// Watch
	go w.Watch()

	// Read statuses until end of time
	for {
		s := <-ss
		printStatus(s)
	}
}

func initParams() (ps params.Params, err error) {
	// Set default values for Params
	ps.Recursive = true
	ps.IntervalMillis = 3000

	usr, err := user.Current()
	if err != nil {
		fmt.Printf("Could not get current user: %v", err)
		os.Exit(1)
	}

	// Read parameters from resources, if present
	err = ps.AddParamsFromResourceFile([]string{"./.eyerc", usr.HomeDir + "/.eyerc"})
	if err != nil {
		fmt.Printf("Error reading resource file: %v\n", err)
		os.Exit(1)
	}

	// Command line arguments will override, if present
	err = ps.AddParamsFromCommandLine()
	if err != nil {
		fmt.Printf("Error in arguments: %v\n", err)
		os.Exit(1)
	}

	// Pattern and Command are mandatory
	if ps.Pattern == "" {
		return ps, fmt.Errorf("pattern is empty")
	}
	if ps.Command == "" {
		return ps, fmt.Errorf("command is empty")
	}

	return ps, nil
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
	fmt.Print("Resources read from ./.eyerc or ~/.eyerc will be overridden by command line parameters\n\n")
	flag.PrintDefaults()
}
