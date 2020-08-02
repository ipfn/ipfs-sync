# ipfs-sync

[![IPFN project](https://img.shields.io/badge/project-IPFN-blue.svg?style=flat-square)](//github.com/ipfn)
[![IPFS project](https://img.shields.io/badge/project-IPFS-blue.svg?style=flat-square)](//github.com/ipfs)
[![GoDoc](https://godoc.org/github.com/ipfn/ipfs-sync/sync?status.svg)](https://godoc.org/github.com/ipfn/ipfs-sync/sync)
[![Travis CI](https://travis-ci.org/ipfn/ipfs-sync.svg?branch=master)](https://travis-ci.org/ipfn/ipfs-sync)

Atomically syncs changes in directory on IPFS.

## Install

Installation requires Go:

```console
$ go get -u github.com/ipfn/ipfs-sync
```

## Usage

```console
$ ipfs-sync --node-addr=multiaddr <directory>
```

Publish to IPNS
```console
$ ipfs-sync --node-addr=multiaddr --publish=true --key=QmdZjqVMbATh1YVdbtZA9j9aVkGKUXjM4aCHAT6yGiPY9R <directory>
```

You can also specify name of the key. 
```console
$ ipfs-sync --node-addr=multiaddr --publish=true --key=testKey <directory>
```

If key is not specified, self key will be used by default. 
```console
$ ipfs-sync --node-addr=multiaddr --publish=true <directory>
```
## License

                                 Apache License
                           Version 2.0, January 2004
                        http://www.apache.org/licenses/
