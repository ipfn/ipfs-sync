package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/ipfn/ipfs-sync/shell"
	"github.com/ipfn/ipfs-sync/sync"
)

var (
	verbose = flag.Bool("verbose", false, "Print logs to stderr")
	nodeURL = flag.String("node-addr", "/ip4/127.0.0.1/tcp/5001/", "IPFS node URL")
	publish = flag.Bool("publish", false, "Option to publish to IPNS")
	key     = flag.String("key", "", "key to publish to publish to IPNS")

	// Ignore .git and .gitignore files by default
	git = flag.Bool("git", true, "Ignores files from .gitignore and .git directory itself.")

	// IPFS ignore rules
	ignore          stringList
	ignoreRulesPath = flag.String("ignore-rules-path", "", "Ignores files from .gitignore.")
)

func init() {
	flag.Var(&ignore, "ignore", "List of paths to ignore.")
}

func main() {
	flag.Parse()
	if *verbose {
		log.SetOutput(os.Stderr)
	} else {
		log.SetOutput(ioutil.Discard)
	}

	path := flag.Arg(0)
	if *key == "" {
		*key = "self"
	}
	if path == "" {
		log.Fatal("Usage: ipfs-sync --node-addr=multiaddr <directory>")
	}

	if path == "." {
		var err error
		if path, err = os.Getwd(); err != nil {
			log.Fatalf("getwd: %v", err)
		}
	} else {
		var err error
		if path, err = filepath.Abs(path); err != nil {
			log.Fatalf("path error: %v", err)
		}
	}

	log.Printf("Starting in %s", path)

	if _, err := exec.LookPath("ipfs"); err != nil {
		log.Fatal("Error: ipfs was not found in $PATH")
	}

	_, err := os.Stat(".gitignore")
	if *git && !os.IsNotExist(err) {
		// Since --git is true by default we only respect --ignore-rules-path flag.
		if len(*ignoreRulesPath) == 0 {
			*ignoreRulesPath = ".gitignore"
		}
		ignore = append(ignore, ".git")
	}

	snc, err := sync.Watch(*nodeURL, path, shell.AddOptions{
		Ignore:          ignore,
		IgnoreRulesPath: *ignoreRulesPath,
	})
	if err != nil {
		log.Fatalf("watch error: %v", err)
	}

	var cmd *exec.Cmd
	cmd = checkPublish(cmd, snc.Hash(), true)

	for hash := range snc.Events() {
		cmd = checkPublish(cmd, hash, false)
	}
}

func checkPublish(cmd *exec.Cmd, hash string, printLogs bool) *exec.Cmd {
	var err error
	if *publish && hash != "" {
		cmd, err = shell.Publish(cmd, hash, *key)
		if err != nil {
			fmt.Printf("Publish error: %s\n", err.Error())
			os.Exit(1)
		} else if printLogs {
			fmt.Printf("Publishing to key %s\n", *key)
		}
	} else if !*publish {
		fmt.Println(hash)
	}
	return cmd
}

type stringList []string

func (i *stringList) String() string {
	return strings.Join(*i, ", ")
}

func (i *stringList) Set(value string) error {
	*i = append(*i, value)
	return nil
}
