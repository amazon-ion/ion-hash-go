/*
 * Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License").
 * You may not use this file except in compliance with the License.
 * A copy of the License is located at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * or in the "license" file accompanying this file. This file is distributed
 * on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
 * express or implied. See the License for the specific language governing
 * permissions and limitations under the License.
 */

package main

import (
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/amzn/ion-go/ion"
	"github.com/amzn/ion-hash-go"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	if len(os.Args) < 4 {
		fmt.Println("Utility that prints the Ion Hash of the top-level values in a file.")
		fmt.Println()
		fmt.Println("Usage:")
		fmt.Println("  ion-hash [algorithm] [filename]")
		fmt.Println()
		fmt.Println("where [algorithm] is a hash function such as sha256")
		fmt.Println()
		os.Exit(1)
	}

	algorithm := os.Args[1]
	fileName := os.Args[2]

	fmt.Println(algorithm)
	fmt.Println(fileName)

	data, err := ioutil.ReadFile(fileName)
	check(err)
	fmt.Println(data)


	ionReader := ion.NewReaderBytes(data)
	hashReader, err := ionhash.NewHashReader(ionReader, ionhash.NewCryptoHasherProvider(ionhash.Algorithm(strings.ToUpper(algorithm))))
	check(err)

	for hashReader.Next() {
		digest, err := hashReader.Sum(nil)
		if err != nil {
			fmt.Printf(`[unable to digest:%v]`, err)
		} else {
			fmt.Printf("%s", hex.Dump(digest))
		}
	}
}
