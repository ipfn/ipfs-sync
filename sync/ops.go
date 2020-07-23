package sync

import (
	"errors"
	"fmt"
	"log"
	"path/filepath"
	"strings"

	fsnotify "gopkg.in/fsnotify.v1"

	"github.com/ipfn/ipfs-sync/shell"
)

// Ops - Sync ops.
type Ops struct {
	base string
	opts shell.AddOptions
}

// Handle - Handles file system notification.
func (ops *Ops) Handle(last string, event fsnotify.Event) (hash string, err error) {
	path := cleanPath(ops.base, event.Name)
	for _, ignored := range ops.opts.Ignore {
		dirpath := fmt.Sprintf("%s/", ignored)
		if path == ignored || path == dirpath || strings.HasPrefix(path, dirpath) {
			return last, nil
		}
	}
	log.Printf("File event: %s;%s (%s)", last, path, eventType(event))
	if event.Op&fsnotify.Create == fsnotify.Create {
		return ops.Create(last, path)
	}
	if event.Op&fsnotify.Remove == fsnotify.Remove {
		return ops.Remove(last, path)
	}
	if event.Op&fsnotify.Write == fsnotify.Write {
		return ops.Write(last, path)
	}
	if event.Op&fsnotify.Rename == fsnotify.Rename {
		return ops.Rename(last, path)
	}
	err = errors.New("unknown error type")
	return
}

// Create - Create event operation.
func (ops *Ops) Create(last string, path string) (hash string, err error) {
	item, err := shell.Add(&ops.opts, filepath.Join(ops.base, path))
	if err != nil {
		return
	}
	return shell.AddLink(last, path, item)
}

// Remove - Remove event operation.
func (ops *Ops) Remove(last string, path string) (hash string, err error) {
	return shell.RmLink(last, path)
}

// Write - Write event operation.
func (ops *Ops) Write(last string, path string) (hash string, err error) {
	item, err := shell.Add(&ops.opts, filepath.Join(ops.base, path))
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
	path = strings.ReplaceAll(path, "\\", "/")
	path = strings.TrimPrefix(path, base)
	return strings.TrimPrefix(path, "/")
}

func eventType(event fsnotify.Event) string {
	if event.Op&fsnotify.Create == fsnotify.Create {
		return "Create"
	}
	if event.Op&fsnotify.Remove == fsnotify.Remove {
		return "Remove"
	}
	if event.Op&fsnotify.Write == fsnotify.Write {
		return "Write"
	}
	if event.Op&fsnotify.Rename == fsnotify.Rename {
		return "Rename"
	}
	return "Unknown"
}
