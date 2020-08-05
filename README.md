:warning: **This package is considered experimental, under active early development, and the API is subject to change.** :warning:

# Amazon Ion Hash Go

An implementation of [Amazon Ion Hash](http://amzn.github.io/ion-hash) in Go.

[![build](https://github.com/amzn/ion-hash-go/workflows/Build/badge.svg)](https://github.com/amzn/ion-hash-go/actions?query=workflow%3ABuild)
[![license](https://img.shields.io/hexpm/l/plug.svg)](https://github.com/amzn/ion-hash-go/blob/master/LICENSE)
[![docs](https://img.shields.io/badge/docs-api-green.svg?style=flat-square)](https://amzn.github.io/ion-hash-go/api)

## Getting Started

The following example code illustrates how to use it:

```go
func main() {
	chp := ionhash.NewCryptoHasherProvider(ionhash.SHA256)

	// Write a simple Ion struct and compute the hash.
	str := strings.Builder{}
	w := ion.NewTextWriter(&str)
	hw, _ := ionhash.NewHashWriter(w, chp)

	fmt.Println("writer")
	hw.BeginStruct()
	hw.FieldName("first_name")
	hw.WriteString("Amanda")
	hw.FieldName("middle_name")
	hw.WriteString("Amanda")
	hw.FieldName("last_name")
	hw.WriteString("Smith")
	hw.EndStruct()
	hw.Finish()
	digest, _ := hw.Sum(nil)
	fmt.Println(bytesToHex(digest))

	ionData := fmt.Sprintf("Ion data: %v", str.String())
	fmt.Println(ionData)

	// Read the struct and compute the hash.
	r := ion.NewReaderString(str.String())
	hr, _ := ionhash.NewHashReader(r, chp)

	fmt.Println("reader")
	hr.Next() // Position reader at the first value.
	hr.Next() // Position reader just after the struct.
	digest, _ = hr.Sum(nil)
	fmt.Println(bytesToHex(digest))
}


func bytesToHex(b []byte) string {
	s := hex.EncodeToString(b)

	buffer := bytes.Buffer{}
	for i, rune := range s {
		buffer.WriteRune(rune)
		if i%2 == 1 {
			buffer.WriteRune(' ')
		}
	}
	return buffer.String()
}
```

Upon execution, the above code produces the following output:

```
writer
37 82 6e 71 92 a1 e4 e1 24 aa 73 f9 85 0f f1 0f 1c b5 cc ca f2 07 b0 9e 65 af 42 56 ae 8c 80 55 
Ion data: {first_name:"Amanda",middle_name:"Amanda",last_name:"Smith"}

reader
37 82 6e 71 92 a1 e4 e1 24 aa 73 f9 85 0f f1 0f 1c b5 cc ca f2 07 b0 9e 65 af 42 56 ae 8c 80 55 
```

## Development

This repository contains a [git submodule](https://git-scm.com/docs/git-submodule)
called `ion-hash-test`, which holds test data used by `ion-hash-go`'s unit tests.

The easiest way to clone the `ion-hash-go` repository and initialize its `ion-hash-test`
submodule is to run the following command:

```
$ git clone --recursive https://github.com/amzn/ion-hash-go.git ion-hash-go
```

Alternatively, the submodule may be initialized independently from the clone
by running the following commands:

```
$ git submodule init
$ git submodule update
```

## Known Issues

Any tests commented out in [ion_hash_tests.ion](https://github.com/amzn/ion-hash-go/blob/master/ion_hash_tests.ion)
are not expected to work at this time.

## License

This library is licensed under the Apache-2.0 License.
