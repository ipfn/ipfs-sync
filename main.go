package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/ipfn/ipfs-sync/sync"
)

var (
	nodeURL = flag.String("node-addr", "/ip4/127.0.0.1/tcp/5001/", "IPFS node URL")
)

func main() {
	flag.Parse()

	path := flag.Arg(0)
	if path == "" {
		log.Fatal("Usage: ipfs-sync --node-addr=multiaddr <directory>")
	}

	snc, err := sync.Watch(*nodeURL, path)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(snc.Hash())

	for hash := range snc.Events() {
		fmt.Println(hash)
	}
}
