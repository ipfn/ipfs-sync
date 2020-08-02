package shell

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"os/exec"
	"strings"
)

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
	lines := strings.Split(keys, "\n")
	for _, line := range lines {
		keyNames := strings.Split(line, " ")
		for _, keyName := range keyNames {
			if keyName != "" && keyName == key {
				return true
			}
		}
	}
	return false
}

// Publish - Publishes the ipfs hash to the provided ipns key
func Publish(cmd *exec.Cmd, hash, key string) (*exec.Cmd, error) {
	if !keyExists(key) {
		return nil, errors.New("Key does not exist")
	}
	if cmd != nil {
		cmd.Process.Kill()
	}

	cmd = exec.Command("ipfs", "name", "publish", hash)
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}
	go func() { cmd.Wait() }()
	return cmd, nil
}
