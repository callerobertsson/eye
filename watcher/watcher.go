// Package watcher defines the type that watches the files for changes matching a pattern in a file structure on an interval
package watcher

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

// Watcher data type
type Watcher struct {
	p *regexp.Regexp // regex?
	c func()         // executed on change
	r bool           // recurse sub directories
	i time.Duration  // interval between change scan
	s chan Status    // status reporting channel
}

// entry represents the path to a file to watch and the modification time
type entry struct {
	path    string
	changed time.Time
}

// New creates a new Watcher
func New(p *regexp.Regexp, c func(), recursive bool, i time.Duration, ch chan Status) Watcher {
	return Watcher{p, c, recursive, i, ch}
}

// NewDefault creates an new Watcher with default values
func NewDefault(p *regexp.Regexp, c func()) Watcher {
	return New(p, c, true, 1*time.Second, nil)
}

// Watch watches for changes and execute command
func (w Watcher) Watch() {
	entries, total := getMatchingEntries(".", w.p, w.r)
	w.reportStatus(newStatus(StatusInfo, "", "Watching %d files of %d\n", len(entries), total))

	for {
		time.Sleep(w.i)

		newEntries, _ := getMatchingEntries(".", w.p, w.r)
		s := w.getChangeStatus(entries, newEntries)
		if s.Type == StatusNone {
			continue
		}

		w.reportStatus(s)

		// Execute cammond
		w.c()

		entries = newEntries
		w.reportStatus(newStatus(StatusInfo, "", "Watching %d files of %d\n", len(entries), total))
	}
}

// Report status to channel, if channel is defined
func (w Watcher) reportStatus(s Status) {
	if w.s != nil {
		w.s <- s
	}
}

// RunSystemCommand runs the command and display the output during the execution
// Can be used by client when wrapped in anonymous func
func RunSystemCommand(c string) error {
	args := strings.Fields(c)

	cmd := exec.Command(args[0], args[1:]...)
	outPipe, _ := cmd.StdoutPipe()
	errPipe, _ := cmd.StderrPipe()

	if err := cmd.Start(); err != nil {
		return err
	}
	reader := bufio.NewReader(outPipe)
	line, err := reader.ReadString('\n')
	for err == nil {
		fmt.Printf("%v", line)
		line, err = reader.ReadString('\n')
	}
	if err != nil && err != io.EOF {
		return err
	}
	reader = bufio.NewReader(errPipe)
	line, err = reader.ReadString('\n')
	for err == nil {
		fmt.Printf("%v", line)
		line, err = reader.ReadString('\n')
	}
	if err != nil && err != io.EOF {
		return err
	}

	return cmd.Wait()
}

// Get matching entries in dir, and if recursive, all subdirs
// Returns list of matches and the total count of files in tree
func getMatchingEntries(dir string, r *regexp.Regexp, recursive bool) ([]entry, int) {

	fis, err := ioutil.ReadDir(dir)
	if err != nil {
		return []entry{}, 0
	}

	entries := []entry{}
	dirs := []string{}
	total := 0

	for _, fi := range fis {
		if strings.Index(fi.Name(), ".") != 0 && fi.IsDir() {
			dirs = append(dirs, dir+"/"+fi.Name())
			continue
		}
		total++
		if r.MatchString(fi.Name()) {
			entries = append(entries, entry{dir + "/" + fi.Name(), fi.ModTime()})
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
func (w Watcher) getChangeStatus(old, new []entry) Status {

	if len(old) > len(new) {
		return Status{StatusDeleted, "", fmt.Sprintf("%d file(s) was removed", len(old)-len(new))}
	}

	if len(old) < len(new) {
		return Status{StatusAdded, "", fmt.Sprintf("%d file(s) was added", len(new)-len(old))}
	}

	for i := range old {
		if old[i] != new[i] {
			f := old[i].path
			return Status{StatusModified, f, fmt.Sprintf("The file %q was modified\n", f)}
		}
	}

	return Status{StatusNone, "", ""}
}
