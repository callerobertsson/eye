// Package params defines the Params data type and functions
package params

import (
	"flag"
	"fmt"
	"strconv"

	"github.com/callerobertsson/eye/resources"
)

// Params holds the settings collected from .eyerc and overidden by arguments
type Params struct {
	Help           bool
	Pattern        string
	Command        string
	Recursive      bool
	IntervalMillis int
	ResourceFile   string
}

// AddParamsFromResourceFile reads the fs file and sets the values
func (ps *Params) AddParamsFromResourceFile(fs []string) error {

	// Create Resource
	rcs, err := resources.New(fs)
	if err != nil {
		return err
	}

	// Get applicable values from Resource

	if v, ok := rcs.Map["pattern"]; ok {
		ps.Pattern = v
	}
	if v, ok := rcs.Map["command"]; ok {
		ps.Command = v
	}
	if v, ok := rcs.Map["recursive"]; ok {
		recursive, err := strconv.ParseBool(v)
		if err != nil {
			return fmt.Errorf("the recursive setting can only be true or false")
		}
		ps.Recursive = recursive
	}
	if v, ok := rcs.Map["interval-millis"]; ok {
		millis, err := strconv.Atoi(v)
		if err != nil {
			return fmt.Errorf("the interval-millis setting must be a numerical value")
		}
		ps.IntervalMillis = millis
	}

	ps.ResourceFile = rcs.File

	return nil
}

// AddParamsFromCommandLine initializes flags and options
func (ps *Params) AddParamsFromCommandLine() error {

	// Command line flags and options
	var (
		helpFlag             bool
		patternOption        string
		commandOption        string
		noRecursionFlag      bool
		intervalMillisOption int
	)

	// Read resource file
	flag.BoolVar(&helpFlag, "h", false, "Show usage information and exit")
	flag.BoolVar(&helpFlag, "help", false, "Show usage information and exit")
	flag.BoolVar(&noRecursionFlag, "R", false, "Do *not* recurse sub directories")
	flag.StringVar(&patternOption, "p", "", "Matching files regex pattern")
	flag.StringVar(&commandOption, "c", "", "Command to execute on changes")
	flag.IntVar(&intervalMillisOption, "i", 1000, "Interval between checks in millis, defalt 1000 ms")
	flag.Parse()

	// Set values if not default
	ps.Help = helpFlag
	if noRecursionFlag {
		ps.Recursive = false
	}
	if patternOption != "" {
		ps.Pattern = patternOption
	}
	if commandOption != "" {
		ps.Command = commandOption
	}
	if intervalMillisOption != 1000 {
		ps.IntervalMillis = intervalMillisOption
	}

	return nil
}
