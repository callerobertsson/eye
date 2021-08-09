// Package watcher Status data type
package watcher

import "fmt"

// Status data type hold info about one watched file
type Status struct {
	Type    StatusType
	File    string
	Message string
}

// StatusType is the status of a watched file
type StatusType int

// Statuses of StatusType
const (
	StatusNone StatusType = iota
	StatusInfo
	StatusModified
	StatusDeleted
	StatusAdded
)

// Create Status
func newStatus(s StatusType, f string, mf string, a ...interface{}) Status {
	return Status{s, f, fmt.Sprintf(mf, a...)}
}
