package types

import "time"

type Entry struct {
	Url          string
	Title        string
	Text         string
	Date         time.Time
	ChildrenURLs map[string]string
	Children     map[string]*Entry
}

type Storage interface {
	Init(url string)
	Write(entry *Entry) error
	ResolveConnections() error
	GetLeafs() (es []string, e error)
	Shutdown()
}