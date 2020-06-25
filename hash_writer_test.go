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
	"bytes"
	"io/ioutil"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/amzn/ion-go/ion"
)

func TestWriteNull(t *testing.T) {
	str := strings.Builder{}
	hw, err := NewHashWriter(ion.NewTextWriter(&str), newIdentityHasherProvider())
	if err != nil {
		t.Fatalf("Expected NewHashWriter() to successfully create a HashWriter; %s", err.Error())
	}

	ionHashWriter, ok := hw.(*hashWriter)
	if !ok {
		t.Fatal("Expected hw to be of type hashWriter")
	}

	err = ionHashWriter.WriteNull()
	if err != nil {
		t.Errorf("Something went wrong executing ionHashWriter.WriteNull(); %s", err.Error())
	}

	sum, err := ionHashWriter.Sum(nil)
	if err != nil {
		t.Fatalf("Something went wrong executing ionHashWriter.Sum(nil); %s", err.Error())
	}

	expectedSum := []byte{0x0b, 0x0f, 0x0e}

	if !reflect.DeepEqual(sum, expectedSum) {
		t.Errorf("sum did not match expectation;\n"+
			"Expected sum: %v\n"+
			"Actual sum:   %v",
			expectedSum, sum)
	}

	err = ionHashWriter.WriteNullType(ion.FloatType)
	if err != nil {
		t.Errorf("Something went wrong executing ionHashWriter.WriteNullType(ion.FloatType); %s", err.Error())
	}

	sum, err = ionHashWriter.Sum(nil)
	if err != nil {
		t.Fatalf("Something went wrong executing ionHashWriter.Sum(nil); %s", err.Error())
	}

	expectedSum = []byte{0x0b, 0x4f, 0x0e}

	if !reflect.DeepEqual(sum, expectedSum) {
		t.Errorf("sum did not match expectation;\n"+
			"Expected sum: %v\n"+
			"Actual sum:   %v",
			expectedSum, sum)
	}

	err = ionHashWriter.WriteNullType(ion.BlobType)
	if err != nil {
		t.Errorf("Something went wrong executing ionHashWriter.WriteNullType(ion.BlobType); %s", err.Error())
	}

	sum, err = ionHashWriter.Sum(nil)
	if err != nil {
		t.Fatalf("Something went wrong executing ionHashWriter.Sum(nil); %s", err.Error())
	}

	expectedSum = []byte{0x0b, 0xaf, 0x0e}

	if !reflect.DeepEqual(sum, expectedSum) {
		t.Errorf("sum did not match expectation;\n"+
			"Expected sum: %v\n"+
			"Actual sum:   %v",
			expectedSum, sum)
	}

	err = ionHashWriter.WriteNullType(ion.StructType)
	if err != nil {
		t.Errorf("Something went wrong executing ionHashWriter.WriteNullType(ion.StructType); %s", err.Error())
	}

	sum, err = ionHashWriter.Sum(nil)
	if err != nil {
		t.Fatalf("Something went wrong executing ionHashWriter.Sum(nil); %s", err.Error())
	}

	expectedSum = []byte{0x0b, 0xdf, 0x0e}

	if !reflect.DeepEqual(sum, expectedSum) {
		t.Errorf("sum did not match expectation;\n"+
			"Expected sum: %v\n"+
			"Actual sum:   %v",
			expectedSum, sum)
	}

	err = ionHashWriter.Finish()
	if err != nil {
		t.Errorf("Something went wrong executing ionHashWriter.Finish(); %s", err.Error())
	}

	// We're comparing splits because str.String() uses a cumbersome '\n' separator
	expected := strings.Split("null null.float null.blob null.struct ", " ")
	actual := strings.Split(str.String(), "\n")

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("str.String() did not match expectation;\n"+
			"Expected str.String(): %v"+
			"Actual str.string():   %v",
			expected, actual)
	}
}

func TestWriteScalars(t *testing.T) {
	t.Skip() // Skipping test until final str.String() check passes

	str := strings.Builder{}
	hw, err := NewHashWriter(ion.NewTextWriter(&str), newIdentityHasherProvider())
	if err != nil {
		t.Fatalf("Expected NewHashWriter() to successfully create a HashWriter; %s", err.Error())
	}

	ionHashWriter, ok := hw.(*hashWriter)
	if !ok {
		t.Fatal("Expected hw to be of type hashWriter")
	}

	sum, err := ionHashWriter.Sum(nil)
	if err != nil {
		t.Fatalf("Something went wrong executing ionHashWriter.Sum(nil); %s", err.Error())
	}

	if !reflect.DeepEqual(sum, []byte{}) {
		t.Errorf("sum did not match expectation;\n"+
			"Expected sum: %v\n"+
			"Actual sum:   %v",
			[]byte{}, sum)
	}

	err = ionHashWriter.WriteInt(5)
	if err != nil {
		t.Errorf("Something went wrong executing ionHashWriter.WriteInt(5); %s", err.Error())
	}

	sum, err = ionHashWriter.Sum(nil)
	if err != nil {
		t.Fatalf("Something went wrong executing ionHashWriter.Sum(nil); %s", err.Error())
	}

	expectedSum := []byte{0x0b, 0x20, 0x05, 0x0e}

	if !reflect.DeepEqual(sum, expectedSum) {
		t.Errorf("sum did not match expectation;\n"+
			"Expected sum: %v\n"+
			"Actual sum:   %v",
			expectedSum, sum)
	}

	err = ionHashWriter.WriteFloat(3.14)
	if err != nil {
		t.Errorf("Something went wrong executing ionHashWriter.WriteInt(5); %s", err.Error())
	}

	sum, err = ionHashWriter.Sum(nil)
	if err != nil {
		t.Fatalf("Something went wrong executing ionHashWriter.Sum(nil); %s", err.Error())
	}

	expectedSum = []byte{0x0b, 0x40, 0x40, 0x09, 0x1e, 0xb8, 0x51, 0xeb, 0x85, 0x1f, 0x0e}

	if !reflect.DeepEqual(sum, expectedSum) {
		t.Errorf("sum did not match expectation;\n"+
			"Expected sum: %v\n"+
			"Actual sum:   %v",
			expectedSum, sum)
	}

	err = ionHashWriter.WriteTimestamp(time.Date(1941, time.December, 7, 18, 0, 0, 0, time.UTC))
	if err != nil {
		t.Errorf("Something went wrong executing ionHashWriter.WriteTimestamp(time.Date(...)); %s", err.Error())
	}

	sum, err = ionHashWriter.Sum(nil)
	if err != nil {
		t.Fatalf("Something went wrong executing ionHashWriter.Sum(nil); %s", err.Error())
	}

	expectedSum = []byte{0x0b, 0x60, 0x80, 0x0f, 0x95, 0x8c, 0x87, 0x92, 0x80, 0x80, 0x0e}

	if !reflect.DeepEqual(sum, expectedSum) {
		t.Errorf("sum did not match expectation;\n"+
			"Expected sum: %v\n"+
			"Actual sum:   %v",
			expectedSum, sum)
	}

	err = ionHashWriter.WriteBlob(
		[]byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f})
	if err != nil {
		t.Errorf("Something went wrong executing ionHashWriter.WriteBlob(...); %s", err.Error())
	}

	sum, err = ionHashWriter.Sum(nil)
	if err != nil {
		t.Fatalf("Something went wrong executing ionHashWriter.Sum(nil); %s", err.Error())
	}

	expectedSum = []byte{0x0b, 0xa0, 0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08,
		0x09, 0x0a, 0x0c, 0x0b, 0x0c, 0x0c, 0x0d, 0x0c, 0x0e, 0x0f, 0x0e}

	if !reflect.DeepEqual(sum, expectedSum) {
		t.Errorf("sum did not match expectation;\n"+
			"Expected sum: %v\n"+
			"Actual sum:   %v",
			expectedSum, sum)
	}

	err = ionHashWriter.Finish()
	if err != nil {
		t.Errorf("Something went wrong executing ionHashWriter.Finish(); %s", err.Error())
	}

	// We're comparing splits because str.String() uses a cumbersome '\n' separator
	expected := strings.Split("5 3.14e0 1941-12-07T18:00:00.0000000-00:00 {{AAECAwQFBgcICQoLDA0ODw==}} ", " ")
	actual := strings.Split(str.String(), "\n")

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("str.String() did not match expectation;\n"+
			"Expected str.String(): %v\n"+
			"Actual str.string():   %v",
			expected, actual)
	}
}

func TestWriteContainers(t *testing.T) {
	t.Skip() // Skipping test until reader's IsInStruct logic matches dot net

	str := strings.Builder{}
	hw, err := NewHashWriter(ion.NewTextWriter(&str), newIdentityHasherProvider())
	if err != nil {
		t.Fatalf("Something went wrong executing NewHashWriter(); %s", err.Error())
	}

	ionHashWriter, ok := hw.(*hashWriter)
	if !ok {
		t.Fatal("Expected hw to be of type hashWriter")
	}

	sum, err := ionHashWriter.Sum(nil)
	if err != nil {
		t.Fatalf("Something went wrong executing ionHashWriter.Sum(nil); %s", err.Error())
	}

	if !reflect.DeepEqual(sum, []byte{}) {
		t.Errorf("sum did not match expectation;\n"+
			"Expected sum: %v\n"+
			"Actual sum:   %v",
			[]byte{}, sum)
	}

	err = ionHashWriter.BeginList()
	if err != nil {
		t.Errorf("Something went wrong executing ionHashWriter.BeginList(); %s", err.Error())
	}

	sum, err = ionHashWriter.Sum(nil)
	if err != nil {
		_, ok := err.(*InvalidOperationError)
		if !ok {
			t.Errorf("Expected ionHashWriter.Sum(nil) to return an InvalidOperationError; %s", err.Error())
		}
	} else {
		t.Error("Expected ionHashWriter.Sum(nil) to return an error")
	}

	err = ionHashWriter.WriteBool(true)
	if err != nil {
		t.Errorf("Something went wrong executing ionHashWriter.WriteBool(true); %s", err.Error())
	}

	sum, err = ionHashWriter.Sum(nil)
	if err != nil {
		_, ok := err.(*InvalidOperationError)
		if !ok {
			t.Errorf("Expected ionHashWriter.Sum(nil) to return an InvalidOperationError; %s", err.Error())
		}
	} else {
		t.Error("Expected ionHashWriter.Sum(nil) to return an error")
	}

	err = ionHashWriter.EndList()
	if err != nil {
		t.Errorf("Something went wrong executing ionHashWriter.EndList(); %s", err.Error())
	}

	sum, err = ionHashWriter.Sum(nil)
	if err != nil {
		t.Fatalf("Something went wrong executing ionHashWriter.Sum(nil); %s", err.Error())
	}

	expectedSum := []byte{0x0b, 0xb0, 0x0b, 0x11, 0x0e, 0x0e}

	if !reflect.DeepEqual(sum, expectedSum) {
		t.Errorf("sum did not match expectation;\n"+
			"Expected sum: %v\n"+
			"Actual sum:   %v",
			expectedSum, sum)
	}

	if ionHashWriter.isInStruct() {
		t.Error("Expected ionHashWriter.isInStruct() to return false")
	}

	err = ionHashWriter.BeginStruct()
	if err != nil {
		t.Errorf("Something went wrong executing ionHashWriter.BeginStruct(); %s", err.Error())
	}

	if !ionHashWriter.isInStruct() {
		t.Error("Expected ionHashWriter.isInStruct() to return true")
	}

	err = ionHashWriter.FieldName("hello")
	if err != nil {
		t.Errorf("Something went wrong executing ionHashWriter.FieldName(\"hello\"); %s", err.Error())
	}

	err = ionHashWriter.Annotation("ion")
	if err != nil {
		t.Errorf("Something went wrong executing ionHashWriter.Annotation(\"ion\"); %s", err.Error())
	}

	err = ionHashWriter.Annotation("hash")
	if err != nil {
		t.Errorf("Something went wrong executing ionHashWriter.Annotation(\"hash\"); %s", err.Error())
	}

	err = ionHashWriter.WriteSymbol("world")
	if err != nil {
		t.Errorf("Something went wrong executing ionHashWriter.WriteSymbol(\"world\"); %s", err.Error())
	}

	err = ionHashWriter.EndStruct()
	if err != nil {
		t.Errorf("Something went wrong executing ionHashWriter.EndStruct(); %s", err.Error())
	}

	if ionHashWriter.isInStruct() {
		t.Error("Expected ionHashWriter.isInStruct() to return false")
	}

	sum, err = ionHashWriter.Sum(nil)
	if err != nil {
		t.Fatalf("Something went wrong executing ionHashWriter.Sum(nil); %s", err.Error())
	}

	expectedSum = []byte{0x0b, 0xd0,
		0x0c, 0x0b, 0x70, 0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x0c, 0x0e, // hello:
		0x0c, 0x0b, 0xe0,
		0x0c, 0x0b, 0x70, 0x69, 0x6f, 0x6e, 0x0c, 0x0e, // ion::
		0x0c, 0x0b, 0x70, 0x68, 0x61, 0x73, 0x68, 0x0c, 0x0e, // hash::
		0x0c, 0x0b, 0x70, 0x77, 0x6f, 0x72, 0x6c, 0x64, 0x0c, 0x0e, // world
		0x0c, 0x0e,
		0x0e}

	if !reflect.DeepEqual(sum, expectedSum) {
		t.Errorf("sum did not match expectation;\n"+
			"Expected sum: %v\n"+
			"Actual sum:   %v",
			expectedSum, sum)
	}

	err = ionHashWriter.Finish()
	if err != nil {
		t.Errorf("Something went wrong executing ionHashWriter.Finish(); %s", err.Error())
	}

	// We're comparing splits because str.String() uses a cumbersome '\n' separator
	expected := strings.Split("[true] {hello:ion::hash::world} ", " ")
	actual := strings.Split(str.String(), "\n")

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("str.String() did not match expectation;\n"+
			"Expected str.String(): %v"+
			"Actual str.string():   %v",
			expected, actual)
	}
}

func TestExtraEndContainer(t *testing.T) {
	str := strings.Builder{}
	hw, err := NewHashWriter(ion.NewTextWriter(&str), newIdentityHasherProvider())
	if err != nil {
		t.Fatalf("Expected NewHashWriter() to successfully create a HashWriter; %s", err.Error())
	}

	ionHashWriter, ok := hw.(*hashWriter)
	if !ok {
		t.Fatal("Expected hw to be of type hashWriter")
	}

	err = ionHashWriter.EndList()
	if err != nil {
		_, ok := err.(*InvalidOperationError)
		if !ok {
			t.Errorf("Expected ionHashWriter.EndList() to return an InvalidOperationError; %s", err.Error())
		}
	} else {
		t.Error("Expected ionHashWriter.EndList() to return an error")
	}

	err = ionHashWriter.EndSexp()
	if err != nil {
		_, ok := err.(*InvalidOperationError)
		if !ok {
			t.Errorf("Expected ionHashWriter.EndSexp() to return an InvalidOperationError; %s", err.Error())
		}
	} else {
		t.Error("Expected ionHashWriter.EndSexp() to return an error")
	}

	err = ionHashWriter.EndStruct()
	if err != nil {
		_, ok := err.(*InvalidOperationError)
		if !ok {
			t.Errorf("Expected ionHashWriter.EndStruct() to return an InvalidOperationError; %s", err.Error())
		}
	} else {
		t.Error("Expected ionHashWriter.EndStruct() to return an error")
	}
}

func TestIonWriterContractWriteValue(t *testing.T) {
	t.Skip()

	file, err := ioutil.ReadFile("ion-hash-test/ion_hash_tests.ion")
	if err != nil {
		t.Fatalf("Something went wrong loading ion_hash_tests.ion; %s", err.Error())
	}

	reader := ion.NewReaderBytes(file)

	expected := ExerciseWriter(t, reader, false, writeFromReaderToWriterAfterNext)

	actual := ExerciseWriter(t, reader, true, writeFromReaderToWriterAfterNext)

	if len(expected) <= 10 {
		t.Error("Expected the ion writer to write more than 10 bytes")
	}

	if len(actual) <= 10 {
		t.Error("Expected the hash writer to write more than 10 bytes")
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("sum did not match expectation;\n"+
			"Expected sum: %v\n"+
			"Actual sum:   %v",
			expected, actual)
	}
}

func TestIonWriterContractWriteValues(t *testing.T) {
	t.Skip()

	file, err := ioutil.ReadFile("ion-hash-test/ion_hash_tests.ion")
	if err != nil {
		t.Fatalf("Something went wrong loading ion_hash_tests.ion; %s", err.Error())
	}

	reader := ion.NewReaderBytes(file)

	expected := ExerciseWriter(t, reader, false, writeFromReaderToWriter)

	actual := ExerciseWriter(t, reader, true, writeFromReaderToWriter)

	if len(expected) <= 1000 {
		t.Error("Expected the ion writer to write more than 1000 bytes")
	}

	if len(actual) <= 1000 {
		t.Error("Expected the hash writer to write more than 1000 bytes")
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("sum did not match expectation;\n"+
			"Expected sum: %v\n"+
			"Actual sum:   %v",
			expected, actual)
	}
}

func TestWriterUnresolvedSid(t *testing.T) {
	t.Skip() // Skipping test until test is implemented once SymbolToken is available

	// TODO: Implement test once SymbolToken is available
}

func ExerciseWriter(t *testing.T, reader ion.Reader, useHashWriter bool, function writeFunction) []byte {
	var err error

	buf := bytes.Buffer{}
	writer := ion.NewBinaryWriter(&buf)

	if useHashWriter {
		writer, err = NewHashWriter(writer, newIdentityHasherProvider())
		if err != nil {
			t.Fatalf("Expected NewHashWriter() to successfully create a HashWriter; %s", err.Error())
		}
	}

	function(t, reader, writer)

	err = writer.Finish()
	if err != nil {
		t.Errorf("Something went wrong executing writer.Finish(); %s", err.Error())
	}

	return buf.Bytes()
}

type writeFunction func(*testing.T, ion.Reader, ion.Writer)

func writeFromReaderToWriterAfterNext(t *testing.T, reader ion.Reader, writer ion.Writer) {
	if !reader.Next() {
		err := reader.Err()
		if err != nil {
			t.Errorf("Something went wrong executing reader.Next(); %s", err.Error())
		}
	}

	writeFromReaderToWriter(t, reader, writer)
}
