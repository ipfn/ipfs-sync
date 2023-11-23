package sync

import (
	"fmt"
	"log"
	"strings"

	ignore "github.com/crackcomm/go-gitignore"
	"github.com/crackcomm/rfsnotify"

	"github.com/ipfn/ipfs-sync/shell"
)

// Synchronizer - IPFS directory synchronizer.
type Synchronizer struct {
	path string
	hash string

	ignoreRules *ignore.GitIgnore

	ops    *Ops
	watch  *rfsnotify.RWatcher
	events chan string
}

// Watch - Constructs new IPFS synchronizer for a directory.
func Watch(path string, opts shell.AddOptions) (sync *Synchronizer, err error) {
	sync = &Synchronizer{
		path: path,
		ops: &Ops{
			base: strings.ReplaceAll(path, "\\", "/"),
			opts: opts,
		},
	}

	if opts.IgnoreRulesPath == "" {
		sync.ignoreRules, err = ignore.CompileIgnoreLines(opts.IgnorePaths...)
	} else {
		sync.ignoreRules, err = ignore.CompileIgnoreFileAndLines(opts.IgnoreRulesPath, opts.IgnorePaths...)
	}

	if sync.ignoreRules == nil {
		sync.watch, err = rfsnotify.NewWatcher()
	} else {
		sync.watch, err = rfsnotify.NewWatcherWithIgnore(sync.ignoreRules)
	}
	if err != nil {
		return nil, fmt.Errorf("rfsnotify: %v", err)
	}

	sync.hash, err = shell.Add(&sync.ops.opts, path)
	if err != nil {
		return nil, fmt.Errorf("ipfs add %q: %v", path, err)
	}

	go sync.watchForEvents()

	err = sync.watch.AddRecursive(path)
	if err != nil {
		return nil, fmt.Errorf("watch %q: %v", path, err)
	}

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
		case ev := <-sync.watch.Events:
			hash, err := sync.ops.Handle(sync.hash, ev)
			if err != nil {
				log.Println("error:", err)
			}
			if hash == "" {
				continue
			}
			if sync.events != nil && sync.hash != hash && hash != "" {
				sync.events <- hash
			}
			if hash != "" {
				sync.hash = hash
			}
		case err := <-sync.watch.Errors:
			log.Println("error:", err)
		}
	}
}
