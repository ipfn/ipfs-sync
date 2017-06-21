package shell

import (
	"os/exec"
	"strings"
)

// Exec - Executes IPFS command.
func Exec(args ...string) (hash string, err error) {
	output, err := exec.Command("ipfs", args...).Output()
	if err != nil {
		return
	}
	return strings.TrimSpace(string(output)), nil
}

// Add - Adds directory to IPFS and returns its hash.
func Add(path string) (string, error) {
	return Exec("add", "-Q", "-r", path)
}

// RmLink - Removes link from IPFS object and returns new hash.
func RmLink(last, path string) (string, error) {
	return Exec("object", "patch", "rm-link", last, path)
}

// AddLink - Creates link from IPFS object and returns new hash.
func AddLink(last, path, hash string) (string, error) {
	return Exec("object", "patch", "add-link", last, path, hash)
}
