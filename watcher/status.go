package watcher

import "fmt"

type Status struct {
	Type    StatusType
	File    string
	Message string
}

type StatusType int

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
