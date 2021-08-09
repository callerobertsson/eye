// Package resources implements the Resources data type and functions
package resources

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"
)

// Resources holds the map of settings as strings and the resource file used
type Resources struct {
	Map  map[string]string
	File string
}

// Regular expressions for matching lines
var (
	regexpEmptyLine   = regexp.MustCompile(`^\s*$`)
	regexpCommentLine = regexp.MustCompile(`^\s*!.*$`)
	regexpSettingLine = regexp.MustCompile(`^\s*([^\s]+)\s*:\s*(.*)\s*$`)
)

// New tries to read resource files in fs.
// Returns error if a file has error
// NOTE: No error if no files are found, but File will be ""
func New(fs []string) (r Resources, err error) {
	// Try resource files in order, skip if not found
	for _, f := range fs {
		// Try to open file or continue
		bs, err := ioutil.ReadFile(f)
		if err != nil {
			continue
		}

		// Parse resource data
		err = r.parseResource(string(bs))
		r.File = f

		// Exit after successfully reading first file found
		return r, err
	}

	// Not an error, if no file is found
	return r, nil
}

// Loop lines in data and build map of settings as strings, returns error on illegal lines
func (r *Resources) parseResource(data string) error {
	r.Map = make(map[string]string)
	ls := strings.Split(data, "\n")

	for i, l := range ls {
		// Skip empty lines and comments
		if regexpEmptyLine.MatchString(l) || regexpCommentLine.MatchString(l) {
			continue
		}

		// If correct formatted resource setting line
		if regexpSettingLine.MatchString(l) {
			ms := regexpSettingLine.FindStringSubmatch(l)
			r.Map[ms[1]] = ms[2]
			continue
		}

		// Something unmatched is not allowed
		return fmt.Errorf("syntax error on line %d: %q", i, l)
	}

	return nil
}
