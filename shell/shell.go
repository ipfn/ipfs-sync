package shell

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"
)

var once sync.Once

// Exec - Executes IPFS command.
func Exec(args ...string) (hash string, err error) {
	log.Printf("Exec: ipfs %s", strings.Join(args, " "))
	cmd := exec.Command("ipfs", args...)
	var (
		stdout bytes.Buffer
		stderr bytes.Buffer
	)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		err = fmt.Errorf("%v: %q", err, strings.TrimSpace(string(stderr.Bytes())))
		return
	}
	return strings.TrimSpace(string(stdout.Bytes())), nil
}

// AddOptions - IPFS shell `add` command options.
type AddOptions struct {
	// Ignore - List of paths to ignore.
	Ignore []string
	// IgnoreRulesPath - Path of `.gitignore` file or similar.
	IgnoreRulesPath string
}

// Add - Adds directory to IPFS and returns its hash.
func Add(opts *AddOptions, path string) (string, error) {
	args := []string{"add", "-Q", "-H", "-r"}
	if len(opts.IgnoreRulesPath) != 0 {
		args = append(args, fmt.Sprintf("--ignore-rules-path=%s", opts.IgnoreRulesPath))
	}
	if len(opts.Ignore) != 0 {
		for _, arg := range opts.Ignore {
			args = append(args, fmt.Sprintf("--ignore=%s", arg))
		}
	}
	return Exec(append(args, path)...)
}

// RmLink - Removes link from IPFS object and returns new hash.
func RmLink(last, path string) (string, error) {
	return Exec("object", "patch", "rm-link", last, path)
}

// AddLink - Creates link from IPFS object and returns new hash.
func AddLink(last, path, hash string) (string, error) {
	return Exec("object", "patch", "add-link", last, path, hash)
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
func Publish(key string) chan string {
	var ctx context.Context
	var cancel context.CancelFunc
	if key != "" {
		if !keyExists(key) {
			fmt.Printf("Publish error: key %s deoesn't exist\n", key)
			os.Exit(1)
		} else {
			log.Printf("Publishing to key: %s", key)
		}
	}

	ch := make(chan string)
	go func() {
		for {
			if hash := <-ch; hash != "" {
				if key != "" {
					if cancel != nil {
						cancel()
					}
					ctx, cancel = context.WithCancel(context.Background())
					defer cancel()
					err := exec.CommandContext(ctx, "ipfs", "name", "publish", fmt.Sprintf("--key=%s", key), hash).Start()
					if err != nil {
						log.Fatalf("Publish error: %s", err.Error())
					}
				} else {
					fmt.Println(hash)
				}
			}

		}
	}()
	return ch
}
