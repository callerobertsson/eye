package main

import (
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

func main() {
	if len(os.Args) < 3 {
		panic("Too few args")
	}
	eye(os.Args[1], os.Args[2])
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
	for _, fi := range fis {
		if r.MatchString(fi.Name()) {
			entries = append(entries, Entry{fi.Name(), fi.ModTime()})
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
		fmt.Printf("WARNING: Could not execute command %q\n", cmd)
		return
	}

	fmt.Println(string(out))
}
