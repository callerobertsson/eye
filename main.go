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

type Entry struct {
	path    string
	changed time.Time
}

var (
	helpFlag      bool
	recursiveFlag bool
	patternOption string
	commandOption string
)

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

func main() {
	eye(patternOption, commandOption)
}

func usage() {
	fmt.Println(`Usage:
	eye -p <PATTERN> -c <COMMAND>

	PATTERN - a regex pattern for matching files to watch
	COMMAND - the command to execute on changes`)
}

func eye(pattern, command string) {
	fmt.Printf("Eyeing pattern %q for command %q\n", pattern, command)

	r := regexp.MustCompile(pattern)

	entries := getMatchingEntries(".", r)
	for {
		time.Sleep(1 * time.Second)
		newEntries := getMatchingEntries(".", r)
		if isDifferent(entries, newEntries) {
			runCommand(command)
		}
		entries = newEntries
	}
}

func getMatchingEntries(dir string, r *regexp.Regexp) []Entry {
	fis, err := ioutil.ReadDir(dir)
	if err != nil {
		panic("Error reading cwd!")
	}

	entries := []Entry{}
	dirs := []string{}

	for _, fi := range fis {
		if strings.Index(fi.Name(), ".") != 0 && fi.IsDir() {
			dirs = append(dirs, dir+"/"+fi.Name())
			continue
		}
		if r.MatchString(fi.Name()) {
			entries = append(entries, Entry{dir + "/" + fi.Name(), fi.ModTime()})
		}
	}

	if recursiveFlag {
		for _, d := range dirs {
			entries = append(entries, getMatchingEntries(d, r)...)
		}
	}

	return entries
}

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

func runCommand(cmd string) {
	fmt.Printf("Running %q\n", cmd)
	args := strings.Fields(cmd)

	out, err := exec.Command(args[0], args[1:]...).Output()
	if err != nil {
		fmt.Printf("WARNING: Could not execute command %q: %v\n", cmd, err.Error())
	}

	if out != nil {
		fmt.Println(string(out))
	}
}
