package sync

import (
	"errors"
	"path/filepath"
	"strings"

	"github.com/howeyc/fsnotify"
	"github.com/ipfn/ipfs-sync/shell"
)

// Ops - Sync ops.
type Ops struct {
	base string
}

// Handle - Handles file system notification.
func (ops *Ops) Handle(last string, event *fsnotify.FileEvent) (hash string, err error) {
	path := cleanPath(ops.base, event.Name)
	if event.IsAttrib() {
		return
	}
	if event.IsCreate() {
		return ops.Create(last, path)
	}
	if event.IsDelete() {
		return ops.Delete(last, path)
	}
	if event.IsModify() {
		return ops.Modify(last, path)
	}
	if event.IsRename() {
		return ops.Rename(last, path)
	}
	err = errors.New("unknown error type")
	return
}

// Create - Create event operation.
func (ops *Ops) Create(last string, path string) (hash string, err error) {
	item, err := shell.Add(filepath.Join(ops.base, path))
	if err != nil {
		return
	}
	return shell.AddLink(last, path, item)
}

// Delete - Delete event operation.
func (ops *Ops) Delete(last string, path string) (hash string, err error) {
	return shell.RmLink(last, path)
}

// Modify - Modify event operation.
func (ops *Ops) Modify(last string, path string) (hash string, err error) {
	item, err := shell.Add(filepath.Join(ops.base, path))
	if err != nil {
		return
	}
	return shell.AddLink(last, path, item)
}

// Rename - Rename event operation.
func (ops *Ops) Rename(last string, path string) (hash string, err error) {
	return shell.RmLink(last, path)
}

func cleanPath(base, path string) string {
	path = strings.TrimPrefix(path, base)
	return strings.TrimPrefix(path, "/")
}
