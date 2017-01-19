// Eye - A simple file change command executioner
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

// Entry represents the path to a file to watch and the modification time
type Entry struct {
	path    string
	changed time.Time
}

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

	if recursiveFlag {
		fmt.Println("Warning: recursing sub directories not supported in this version. Ignoring.")
	}

	if patternOption == "" || commandOption == "" {
		usage()
		os.Exit(1)
	}
}

// Main function
func main() {
	eye(patternOption, commandOption)
}

// Print usage information
func usage() {
	fmt.Println(`Usage:
	eye -p <PATTERN> -c <COMMAND>

	PATTERN - a regex pattern for matching files to watch
	COMMAND - the command to execute on changes`)
}

// Primus motor
func eye(pattern, command string) {
	fmt.Printf("Eyeing pattern %q for command %q\n", pattern, command)

	r := regexp.MustCompile(pattern)

	t := time.Now()
	entries, total := getMatchingEntries(".", r)
	fmt.Printf("Watching %d files of %d [%v]\n", len(entries), total, time.Now().Sub(t))

	for {
		time.Sleep(1 * time.Second)
		newEntries, _ := getMatchingEntries(".", r)
		if isDifferent(entries, newEntries) {
			runCommand(command)
		}
		entries = newEntries
	}
}

// Get matching entries in dir, and if recursive, all subdirs
// Returns list of matches and the total count of files in tree
func getMatchingEntries(dir string, r *regexp.Regexp) ([]Entry, int) {
	fis, err := ioutil.ReadDir(dir)
	if err != nil {
		return []Entry{}, 0
	}

	entries := []Entry{}
	dirs := []string{}
	total := 0

	for _, fi := range fis {
		if strings.Index(fi.Name(), ".") != 0 && fi.IsDir() {
			dirs = append(dirs, dir+"/"+fi.Name())
			continue
		}
		total++
		if r.MatchString(fi.Name()) {
			entries = append(entries, Entry{dir + "/" + fi.Name(), fi.ModTime()})
		}
	}

	if recursiveFlag {
		for _, d := range dirs {
			newEntries, newTotal := getMatchingEntries(d, r)
			total += newTotal
			entries = append(entries, newEntries...)
		}
	}

	return entries, total
}

// Compare if two lists of entries are equal
func isDifferent(old, new []Entry) bool {

	if len(old) > len(new) {
		fmt.Printf("%d file(s) was removed\n", len(old)-len(new))
		return true
	}
	if len(old) < len(new) {
		fmt.Printf("%d file(s) was added\n", len(new)-len(old))
		return true
	}

	for i := range old {
		if old[i] != new[i] {
			fmt.Printf("The file %q was modified\n", old[i].path)
			return true
		}
	}

	return false
}

// Executes system command
func runCommand(cmd string) {
	fmt.Printf("Running %q\n", cmd)
	args := strings.Fields(cmd)

	out, err := exec.Command(args[0], args[1:]...).CombinedOutput()
	if err != nil {
		fmt.Printf("WARNING: Could not execute command %q: %v\n", cmd, err.Error())
	}

	if out != nil {
		fmt.Println(string(out))
	}
}
