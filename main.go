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
	ipnsKey = flag.String("ipns-key", "", "IPNS publish key or name")

	// Ignore .git and .gitignore files by default
	git = flag.Bool("git", true, "Ignores files from .gitignore and .git directory itself")

	// IPFS ignore rules
	ignore          stringList
	ignoreRulesPath = flag.String("ignore-rules-path", "", "Ignores files from .gitignore")
)

func init() {
	flag.Var(&ignore, "ignore", "List of paths to ignore")
}

func fatal(msg string, v ...interface{}) {
	fmt.Printf(msg, v...)
	os.Exit(1)
}

func main() {
	flag.Parse()
	if *verbose {
		log.SetOutput(os.Stderr)
	} else {
		log.SetOutput(ioutil.Discard)
	}

	path := flag.Arg(0)
	if path == "" {
		fatal("Usage: ipfs-sync --node-addr=multiaddr <directory>")
	}

	if path == "." {
		var err error
		if path, err = os.Getwd(); err != nil {
			fatal("Error: getwd: %v", err)
		}
	} else {
		var err error
		if path, err = filepath.Abs(path); err != nil {
			fatal("Error: filepath: %v", err)
		}
	}

	log.Printf("Starting in %s", path)

	if _, err := exec.LookPath("ipfs"); err != nil {
		fatal("Error: ipfs was not found in $PATH")
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
		fatal("Error: watch: %v", err)
	}
	fmt.Println(snc.Hash())
	var pubChan chan string

	if *ipnsKey != "" {
		pubChan = make(chan string, 1)
		pubChan <- snc.Hash()
		if err := shell.Publish(*ipnsKey, pubChan); err != nil {
			fatal("Publish error: %s\n", err)
		}
		log.Printf("Publishing to key: %s", *ipnsKey)
	}

	for hash := range snc.Events() {
		fmt.Println(hash)
		if pubChan != nil {
			pubChan <- hash
		}
	}
}

type stringList []string

func (i *stringList) String() string {
	return strings.Join(*i, ", ")
}

func (i *stringList) Set(value string) error {
	*i = append(*i, value)
	return nil
}
