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

package ionhash

import (
	"bufio"
	"encoding/base64"
	"os"
	"strings"
	"testing"

	"github.com/amzn/ion-go/ion"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNaughtyStrings(t *testing.T) {
	t.Skip() // Skipping test until ion text reader SymbolTable() is implemented

	file, err := os.Open("ion-hash-test/big_list_of_naughty_strings.txt")
	require.NoError(t, err, "Something went wrong loading big_list_of_naughty_strings.txt")

	defer func() {
		assert.NoError(t, file.Close(), "Something went wrong executing file.Close()")
	}()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		if strings.Trim(line, " ") == "" || strings.HasPrefix(line, "#") {
			continue
		}

		tv := newTestValue(line)

		NaughtyStrings(t, tv, tv.asSymbol())
		NaughtyStrings(t, tv, tv.asString())
		NaughtyStrings(t, tv, tv.asLongString())
		NaughtyStrings(t, tv, tv.asClob())
		NaughtyStrings(t, tv, tv.asBlob())

		NaughtyStrings(t, tv, tv.asSymbol()+"::"+tv.asSymbol())
		NaughtyStrings(t, tv, tv.asSymbol()+"::"+tv.asString())
		NaughtyStrings(t, tv, tv.asSymbol()+"::"+tv.asLongString())
		NaughtyStrings(t, tv, tv.asSymbol()+"::"+tv.asClob())
		NaughtyStrings(t, tv, tv.asSymbol()+"::"+tv.asBlob())

		NaughtyStrings(t, tv, tv.asSymbol()+"::{"+tv.asSymbol()+":"+tv.asSymbol()+"}")
		NaughtyStrings(t, tv, tv.asSymbol()+"::{"+tv.asSymbol()+":"+tv.asString()+"}")
		NaughtyStrings(t, tv, tv.asSymbol()+"::{"+tv.asSymbol()+":"+tv.asLongString()+"}")
		NaughtyStrings(t, tv, tv.asSymbol()+"::{"+tv.asSymbol()+":"+tv.asClob()+"}")
		NaughtyStrings(t, tv, tv.asSymbol()+"::{"+tv.asSymbol()+":"+tv.asBlob()+"}")

		if tv.isValidIon() {
			NaughtyStrings(t, tv, tv.asIon())
			NaughtyStrings(t, tv, tv.asSymbol()+"::"+tv.asIon())
			NaughtyStrings(t, tv, tv.asSymbol()+"::{"+tv.asSymbol()+":"+tv.asIon()+"}")
			NaughtyStrings(t, tv, tv.asSymbol()+"::{"+tv.asSymbol()+":"+tv.asSymbol()+"::"+tv.asIon()+"}")
		}

		list := tv.asSymbol() + "::[" + tv.asSymbol() + ", " + tv.asString() + ", " + tv.asLongString() + ", " + tv.asClob() + ", " + tv.asBlob() + ", "
		if tv.isValidIon() {
			list += tv.asIon()
		}
		list += "]"

		NaughtyStrings(t, tv, list)

		sexp := tv.asSymbol() + "::(" + tv.asSymbol() + " " + tv.asString() + " " + tv.asLongString() + " " + tv.asClob() + " " + tv.asBlob() + " "
		if tv.isValidIon() {
			sexp += tv.asIon()
		}
		sexp += ")"

		NaughtyStrings(t, tv, sexp)

		// multiple annotations
		NaughtyStrings(t, tv, tv.asSymbol()+"::"+tv.asSymbol()+"::"+tv.asSymbol()+"::"+tv.asString())
	}

	assert.NoError(t, scanner.Err(), "Expected scanner to scan without errors")
}

func NaughtyStrings(t *testing.T, tv testValue, s string) {
	hasherProvider := NewCryptoHasherProvider(SHA256)

	str := strings.Builder{}
	hw, err := NewHashWriter(ion.NewTextWriter(&str), hasherProvider)
	require.NoError(t, err, "Expected NewHashWriter() to successfully create a HashWriter")

	ionHashWriter, ok := hw.(*hashWriter)
	require.True(t, ok, "Expected hw to be of type hashWriter")

	writeFromReaderToWriter(t, ion.NewReaderString(s), ionHashWriter)

	hr, err := NewHashReader(ion.NewReaderString(s), hasherProvider)
	require.NoError(t, err, "Expected NewHashReader() to successfully create a HashReader")

	ionHashReader, ok := hr.(*hashReader)
	require.True(t, ok, "Expected hr to be of type hashReader")

	if !ionHashReader.Next() {
		assert.NoError(t, ionHashReader.Err(), "Something went wrong executing ionHashReader.Next()")
	}

	if !ionHashReader.Next() {
		assert.NoError(t, ionHashReader.Err(), "Something went wrong executing ionHashReader.Next()")
	}

	if tv.isValidIon() {
		writerSum, err := ionHashWriter.Sum(nil)
		require.NoError(t, err, "Something went wrong executing ionHashWriter.Sum(nil)")

		readerSum, err := ionHashReader.Sum(nil)
		require.NoError(t, err, "Something went wrong executing ionHashReader.Sum(nil)")

		assert.Equalf(t, writerSum, readerSum, "Expected reader/writer sums for \"%s\" to match", s)
	}
}

type testValue struct {
	ionPrefix        string
	invalidIonPrefix string

	ion      string
	validIon bool
}

func newTestValue(ion string) testValue {
	prefix := "ion::"
	invalidPrefix := "invalid_ion::"
	validIon := false

	if strings.HasPrefix(ion, prefix) {
		validIon = true
		ion = ion[len(prefix):]
	} else if strings.HasPrefix(ion, invalidPrefix) {
		ion = ion[len(invalidPrefix):]
	}

	return testValue{ionPrefix: prefix, invalidIonPrefix: invalidPrefix, ion: ion, validIon: validIon}
}

func (tv *testValue) asIon() string {
	return tv.ion
}

func (tv *testValue) asSymbol() string {
	s := tv.ion
	s = strings.Replace(s, "\\", "\\\\", -1)
	s = strings.Replace(s, "'", "\\'", -1)
	s = "'" + s + "'"

	return s
}

func (tv *testValue) asString() string {
	s := tv.ion
	s = strings.Replace(s, "\\", "\\\\", -1)
	s = strings.Replace(s, "\"", "\\\"", -1)
	s = "\"" + s + "\""

	return s
}

func (tv *testValue) asLongString() string {
	s := tv.ion
	s = strings.Replace(s, "\\", "\\\\", -1)
	s = strings.Replace(s, "'", "\\'", -1)
	s = "'''" + s + "'''"

	return s
}

func (tv *testValue) asClob() string {
	s := ""

	bytes := []byte(tv.asString())

	for _, b := range bytes {
		c := b & 0xFF
		if c >= 128 {
			s += "\\x" + string(c)
		} else {
			s += string(c)
		}
	}

	return s
}

func (tv *testValue) asBlob() string {
	bytes := []byte(tv.asIon())

	return "{{" + base64.StdEncoding.EncodeToString(bytes) + "}}"
}

func (tv *testValue) isValidIon() bool {
	return tv.validIon
}
