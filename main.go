package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/ipfn/ipfs-sync/sync"
)

var (
	verbose = flag.Bool("verbose", false, "Print logs to stderr")
	nodeURL = flag.String("node-addr", "/ip4/127.0.0.1/tcp/5001/", "IPFS node URL")
)

func main() {
	flag.Parse()
	if *verbose {
		log.SetOutput(os.Stderr)
	} else {
		log.SetOutput(ioutil.Discard)
	}

	path := flag.Arg(0)
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

	snc, err := sync.Watch(*nodeURL, path)
	if err != nil {
		log.Fatalf("watch error: %v", err)
	}

	fmt.Println(snc.Hash())

	for hash := range snc.Events() {
		fmt.Println(hash)
	}
}
