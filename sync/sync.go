package sync

import (
	"log"

	"github.com/howeyc/fsnotify"

	"github.com/ipfn/ipfs-sync/shell"
)

// Synchronizer - IPFS directory synchronizer.
type Synchronizer struct {
	path string
	hash string

	ops    *Ops
	watch  *fsnotify.Watcher
	events chan string
}

// Watch - Constructs new IPFS synchronizer for a directory.
func Watch(url, path string) (sync *Synchronizer, err error) {
	watch, err := fsnotify.NewWatcher()
	if err != nil {
		return
	}

	sync = &Synchronizer{
		path: path,
		ops: &Ops{
			base: path,
		},
		watch: watch,
	}

	err = sync.watch.Watch(path)
	if err != nil {
		return
	}

	sync.hash, err = shell.Add(path)
	if err != nil {
		return
	}

	go sync.watchForEvents()
	return
}

// Events - Sends new hashes.
func (sync *Synchronizer) Events() <-chan string {
	if sync.events == nil {
		sync.events = make(chan string, 1)
	}
	return sync.events
}

// Hash - Current hash.
func (sync *Synchronizer) Hash() string {
	return sync.hash
}

// Close - Closes synchronizer.
func (sync *Synchronizer) Close() (err error) {
	return sync.watch.Close()
}

func (sync *Synchronizer) watchForEvents() {
	for {
		select {
		case ev := <-sync.watch.Event:
			hash, err := sync.ops.Handle(sync.hash, ev)
			if err != nil {
				log.Println("error:", err)
			}
			if hash != "" {
				sync.hash = hash
			}
			if sync.events != nil {
				sync.events <- hash
			}
		case err := <-sync.watch.Error:
			log.Println("error:", err)
		}
	}
}
