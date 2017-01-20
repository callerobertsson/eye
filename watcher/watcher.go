// Watches for changes matching a pattern in a file structure on an interval
package watcher

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

type Watcher struct {
	p *regexp.Regexp // regex?
	c CommandFunc    // executed on change
	r bool           // recurse sub directories
	i time.Duration  // interval between change scan
}

type CommandFunc func()

// Entry represents the path to a file to watch and the modification time
type Entry struct {
	path    string
	changed time.Time
}

// Create a new Watcher
func New(p *regexp.Regexp, c CommandFunc, recursive bool, i time.Duration) Watcher {
	return Watcher{p, c, recursive, i}
}

// Watch for changes and execute command
func (w Watcher) Watch() {
	//	fmt.Printf("Eyeing pattern %q for command %q\n", pattern, command)

	entries, _ := getMatchingEntries(".", w.p, w.r)
	//fmt.Printf("Watching %d files of %d [%v]\n", len(entries), total, time.Now().Sub(t))

	for {
		time.Sleep(w.i)
		newEntries, _ := getMatchingEntries(".", w.p, w.r)
		if hasChanged(entries, newEntries) {
			w.c()
		}
		entries = newEntries
	}
}

// Executes system command
func RunSystemCommand(cmd string) {
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

// Get matching entries in dir, and if recursive, all subdirs
// Returns list of matches and the total count of files in tree
func getMatchingEntries(dir string, r *regexp.Regexp, recursive bool) ([]Entry, int) {

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

	if recursive {
		for _, d := range dirs {
			newEntries, newTotal := getMatchingEntries(d, r, recursive)
			total += newTotal
			entries = append(entries, newEntries...)
		}
	}

	return entries, total
}

// Compare if two lists of entries are equal
func hasChanged(old, new []Entry) bool {

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
