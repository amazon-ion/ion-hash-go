** This package is considered beta. While the API is relatively stable it is still subject to change. **

# Amazon Ion Hash Go

An implementation of [Amazon Ion Hash](http://amzn.github.io/ion-hash) in Go.

[![build](https://github.com/amzn/ion-hash-go/workflows/Build/badge.svg)](https://github.com/amzn/ion-hash-go/actions?query=workflow%3ABuild)
[![license](https://img.shields.io/hexpm/l/plug.svg)](https://github.com/amzn/ion-hash-go/blob/master/LICENSE)
[![docs](https://img.shields.io/badge/docs-api-green.svg?style=flat-square)](https://pkg.go.dev/github.com/amzn/ion-hash-go?tab=doc)

## Getting Started

You can start using ion-hash-go by simply importing it, e.g.,

```Go

import (
	ionhash "github.com/amzn/ion-hash-go"
)

```

## Generating a hash while reading

```Go

// Create a hasher provider, using MD5
hasherProvider := ionhash.NewCryptoHasherProvider("MD5")

// Create an Ion reader over the input [1,2,3]
ionReader := ion.NewReaderString("[1,2,3]")

// Create a hash reader
hashReader, err := ionhash.NewHashReader(ionReader, hasherProvider)
if err != nil {
	panic(err)
}

// Read over the top level value and calculate its hash
hashReader.Next()
hashReader.Next()


// Get the hash value
res, err := hashReader.Sum(nil)
if err != nil {
	panic(err)
}

// Print out the hash in Hex
fmt.Printf("Digest = %x\n", res) // prints: Digest = 8f3bf4b1935cf469c9c10c31524b2625

```

## Generating a hash while writing

```Go

// Create a hasher provider, using MD5
hasherProvider := ionhash.NewCryptoHasherProvider("MD5")

// Create an Ion writer
ionWriter := ion.NewTextWriter(new(bytes.Buffer))

// Create a hash writer
hashWriter, err := ionhash.NewHashWriter(ionWriter, hasherProvider)
if err != nil {
	panic(err)
}

// Write the list [1,2,3]
hashWriter.BeginList()
hashWriter.WriteInt(1)
hashWriter.WriteInt(2)
hashWriter.WriteInt(3)
hashWriter.EndList()

// Get the hash value
res, err := hashWriter.Sum(nil)
if err != nil {
	panic(err)
}

// Print out the hash in Hex
fmt.Printf("Digest = %x\n", res) // prints: Digest = 8f3bf4b1935cf469c9c10c31524b2625/

```

## Development

This package uses [Go Modules](https://github.com/golang/go/wiki/Modules) to model
its dependencies.

Assuming the `go` command is in your path, building the module can be done as:

```
$ go build -v ./...
```

Running all the tests can be executed with:

```
$ go test -v ./...
```

We use [`goimports`](https://pkg.go.dev/golang.org/x/tools/cmd/goimports?tab=doc) to format
our imports and files in general.  Running this before commit is advised:

```
$ goimports -w .
```

It is recommended that you hook this in your favorite IDE (`Tools` > `File Watchers` in Goland, for example).

## Known Issues

Any tests commented out in
[ion_hash_tests.ion](https://github.com/amzn/ion-hash-go/blob/master/ion_hash_tests.ion)
are not expected to work at this time.

## License

This library is licensed under the Apache-2.0 License.
