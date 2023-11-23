package shell

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os/exec"
	"strings"
	"sync"
)

var once sync.Once

// Exec - Executes IPFS command.
func Exec(args ...string) (_ string, err error) {
	log.Printf("Exec: ipfs %s", strings.Join(args, " "))
	cmd := exec.Command("ipfs", args...)
	var (
		stdout bytes.Buffer
		stderr bytes.Buffer
	)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err = cmd.Run(); err != nil {
		err = fmt.Errorf("%v: %q", err, strings.TrimSpace(string(stderr.Bytes())))
		return
	}
	return strings.TrimSpace(string(stdout.Bytes())), nil
}

// AddOptions - IPFS shell `add` command options.
type AddOptions struct {
	// IgnorePaths - List of paths to ignore.
	IgnorePaths []string
	// IgnoreRulesPath - Path of `.gitignore` file or similar.
	IgnoreRulesPath string
	// Hidden - Include hidden files.
	Hidden bool
}

// Add - Adds directory to IPFS and returns its hash.
func Add(opts *AddOptions, path string) (string, error) {
	args := []string{"add", "-Q", "-r"}
	if opts.Hidden {
		args = append(args, "-H")
	}
	if len(opts.IgnoreRulesPath) != 0 {
		args = append(args, fmt.Sprintf("--ignore-rules-path=%s", opts.IgnoreRulesPath))
	}
	if len(opts.IgnorePaths) != 0 {
		for _, arg := range opts.IgnorePaths {
			args = append(args, fmt.Sprintf("--ignore=%s", arg))
		}
	}
	return Exec(append(args, path)...)
}

// RmLink - Removes link from IPFS object and returns new hash.
func RmLink(last, path string) (string, error) {
	return Exec("object", "patch", "rm-link", "--allow-big-block", last, path)
}

// AddLink - Creates link from IPFS object and returns new hash.
func AddLink(last, path, hash string) (string, error) {
	return Exec("object", "patch", "add-link", "--allow-big-block", last, path, hash)
}

func keyExists(key string) bool {
	keys, err := Exec("key", "list", "-l")
	if err != nil {
		return false
	}
	for _, line := range strings.Split(keys, "\n") {
		for _, keyName := range strings.Split(line, " ") {
			if keyName != "" && keyName == key {
				return true
			}
		}
	}
	return false
}

// Publish - Publishes the ipfs hash to the provided ipns key
func Publish(key string, hashChan <-chan string) (err error) {
	if !keyExists(key) {
		return fmt.Errorf("key %s doesn't exist", key)
	}
	go func() {
		var cancel context.CancelFunc
		for {
			hash := <-hashChan
			if cancel != nil {
				cancel()
			}
			var ctx context.Context
			ctx, cancel = context.WithCancel(context.Background())
			args := []string{"name", "publish", fmt.Sprintf("--key=%s", key), hash}
			log.Printf("Exec: ipfs %s", strings.Join(args, " "))
			cmd := exec.CommandContext(ctx, "ipfs", args...)
			if err := cmd.Start(); err != nil {
				log.Printf("Publish start error: %v", err)
			}
		}
	}()
	return
}
