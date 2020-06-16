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
	"fmt"
	"io/ioutil"
	"reflect"
	"strings"
	"testing"

	"github.com/amzn/ion-go/ion"
)

func TestMiscMethods(t *testing.T) {
	checkTestToSkip(t)

	str := strings.Builder{}
	writer, err := NewHashWriter(ion.NewTextWriter(&str), newIdentityHasherProvider())
	if err != nil {
		t.Fatalf("expected NewHashWriter() to successfully create a HashWriter; %s", err.Error())
	}

	ionHashWriter, ok := writer.(*hashWriter)
	if !ok {
		t.Fatal("expected ionHashWriter to be of type hashWriter")
	}

	sum, err := ionHashWriter.Sum(nil)
	if err != nil {
		t.Errorf("expected Sum(nil) to execute without errors; %s", err.Error())
	}

	if !reflect.DeepEqual(sum, []byte{}) {
		t.Errorf("expected sum to be %v instead of %v", []byte{}, sum)
	}

	err = ionHashWriter.WriteNull()
	if err != nil {
		t.Errorf("expected ionHashWriter.WriteNull() to execute without errors; %s", err.Error())
	}

	sum, err = ionHashWriter.Sum(nil)
	if err != nil {
		t.Errorf("expected Sum(nil) to execute without errors; %s", err.Error())
	}

	expectedSum := []byte{0x0b, 0x0f, 0x0e}

	if !reflect.DeepEqual(sum, expectedSum) {
		t.Errorf("expected sum to be %v instead of %v", expectedSum, sum)
	}

	err = ionHashWriter.stepIn(ion.ListType)
	if err != nil {
		t.Errorf("expected ionHashWriter.stepIn(ion.ListType) to execute without errors; %s", err.Error())
	}

	sum, err = ionHashWriter.Sum(nil)
	if err != nil {
		_, ok := err.(*InvalidOperationError)
		if !ok {
			t.Errorf("expected Sum(nil) to return an InvalidOperationError; %s", err.Error())
		}
	} else {
		t.Error("expected Sum(nil) to return an error")
	}

	err = ionHashWriter.WriteInt(5)
	if err != nil {
		t.Errorf("expected ionHashWriter.WriteInt(5) to execute without errors; %s", err.Error())
	}

	sum, err = ionHashWriter.Sum(nil)
	if err != nil {
		_, ok := err.(*InvalidOperationError)
		if !ok {
			t.Errorf("expected Sum(nil) to return an InvalidOperationError; %s", err.Error())
		}
	} else {
		t.Error("expected Sum(nil) to return an error")
	}

	err = ionHashWriter.stepOut()
	if err != nil {
		t.Errorf("expected ionHashWriter.stepOut() to execute without errors; %s", err.Error())
	}

	sum, err = ionHashWriter.Sum(nil)
	if err != nil {
		t.Errorf("expected Sum(nil) to execute without errors; %s", err.Error())
	}

	expectedSum = []byte{0x0b, 0xb0, 0x0b, 0x20, 0x05, 0x0e, 0x0e}

	if !reflect.DeepEqual(sum, expectedSum) {
		t.Errorf("expected sum to be %v instead of %v", expectedSum, sum)
	}

	err = ionHashWriter.WriteNull()
	if err != nil {
		t.Errorf("expected ionHashWriter.WriteNull() to execute without errors; %s", err.Error())
	}

	sum, err = ionHashWriter.Sum(nil)
	if err != nil {
		t.Errorf("expected Sum(nil) to execute without errors; %s", err.Error())
	}

	expectedSum = []byte{0x0b, 0x0f, 0x0e}

	if !reflect.DeepEqual(sum, expectedSum) {
		t.Errorf("expected sum to be %v instead of %v", expectedSum, sum)
	}

	if ionHashWriter.isInStruct() {
		t.Error("expected ionHashWriter.isInStruct() to return false")
	}

	err = ionHashWriter.stepIn(ion.StructType)
	if err != nil {
		t.Errorf("expected ionHashWriter.stepIn(ion.StructType) to execute without errors; %s", err.Error())
	}

	if !ionHashWriter.isInStruct() {
		t.Error("expected ionHashWriter.isInStruct() to return true")
	}

	ionHashWriter.setFieldName("hello")

	err = ionHashWriter.Annotation("ion")
	if err != nil {
		t.Errorf("expected ionHashWriter.Annotation(\"ion\") to execute without errors; %s", err.Error())
	}

	err = ionHashWriter.Annotation("hash")
	if err != nil {
		t.Errorf("expected ionHashWriter.Annotation(\"hash\") to execute without errors; %s", err.Error())
	}

	err = ionHashWriter.WriteSymbol("world")
	if err != nil {
		t.Errorf("expected ionHashWriter.WriteSymbol(\"world\") to execute without errors; %s", err.Error())
	}

	err = ionHashWriter.stepOut()
	if err != nil {
		t.Errorf("expected ionHashWriter.stepOut() to execute without errors; %s", err.Error())
	}

	if ionHashWriter.isInStruct() {
		t.Error("expected ionHashWriter.isInStruct() to return false")
	}

	sum, err = ionHashWriter.Sum(nil)
	if err != nil {
		t.Errorf("expected Sum(nil) to execute without errors; %s", err.Error())
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
		t.Errorf("expected sum to be %v instead of %v", expectedSum, sum)
	}

	err = ionHashWriter.Finish()
	if err != nil {
		t.Errorf("expected ionHashWriter.Finish() to execute without errors; %s", err.Error())
	}

	expectedStr := "null [5] null {hello:ion::hash::world}"

	if str.String() != expectedStr {
		t.Errorf("expected str.String() to return \"%s\" instead of %s", expectedStr, str.String())
	}
}

func TestExtraStepOut(t *testing.T) {
	str := strings.Builder{}
	writer, err := NewHashWriter(ion.NewTextWriter(&str), newIdentityHasherProvider())
	if err != nil {
		t.Fatalf("expected NewHashWriter() to successfully create a HashWriter; %s", err.Error())
	}

	ionHashWriter, ok := writer.(*hashWriter)
	if !ok {
		t.Fatal("expected ionHashWriter to be of type hashWriter")
	}

	err = ionHashWriter.stepOut()
	if err != nil {
		_, ok := err.(*InvalidOperationError)
		if !ok {
			t.Errorf("expected ionHashWriter.stepOut() to return an InvalidOperationError; %s", err.Error())
		}
	} else {
		t.Error("expected ionHashWriter.stepOut() to return an error")
	}
}

func TestIonWriterContractWriteValue(t *testing.T) {
	checkTestToSkip(t)

	file, err := ioutil.ReadFile("ion-hash-test/ion_hash_tests.ion")
	if err != nil {
		t.Fatal("expected ion_hash_tests.ion to load properly")
	}

	reader := ion.NewReaderBytes(file)

	expected, err := ExerciseWriter(reader, false, nextWriteValue)
	if err != nil {
		t.Fatal(err.Error())
	}

	actual, err := ExerciseWriter(reader, true, nextWriteValue)
	if err != nil {
		t.Fatal(err.Error())
	}

	if len(expected) <= 10 {
		t.Error("expected the ion writer to write more than 10 bytes")
	}

	if len(actual) <= 10 {
		t.Error("expected the hash writer to write more than 10 bytes")
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("expected the hash writer to write %v instead of %v", expected, actual)
	}
}

func TestIonWriterContractWriteValues(t *testing.T) {
	checkTestToSkip(t)

	file, err := ioutil.ReadFile("ion-hash-test/ion_hash_tests.ion")
	if err != nil {
		t.Fatal("expected ion_hash_tests.ion to load properly")
	}

	reader := ion.NewReaderBytes(file)

	expected, err := ExerciseWriter(reader, false, writeValues)
	if err != nil {
		t.Fatal(err.Error())
	}

	actual, err := ExerciseWriter(reader, true, writeValues)
	if err != nil {
		t.Fatal(err.Error())
	}

	if len(expected) <= 1000 {
		t.Error("expected the ion writer to write more than 1000 bytes")
	}

	if len(actual) <= 1000 {
		t.Error("expected the hash writer to write more than 1000 bytes")
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("expected the hash writer to write %v instead of %v", expected, actual)
	}
}

func TestWriterUnresolvedSid(t *testing.T) {
	checkTestToSkip(t)

	// TODO: Implement test once SymbolToken is available
}

func ExerciseWriter(reader ion.Reader, useHashWriter bool, function writeFunction) ([]byte, error) {
	var err error

	buf := bytes.Buffer{}
	writer := ion.NewBinaryWriter(&buf)

	if useHashWriter {
		writer, err = NewHashWriter(writer, newIdentityHasherProvider())
		if err != nil {
			return nil, fmt.Errorf("expected NewHashWriter() to successfully create a HashWriter; %s", err.Error())
		}
	}

	err = function(reader, writer)
	if err != nil {
		return nil, err
	}

	err = writer.Finish()
	if err != nil {
		return nil, fmt.Errorf("expected writer.Finish() to execute without errors; %s", err.Error())
	}

	return buf.Bytes(), nil
}

type writeFunction func(reader ion.Reader, writer ion.Writer) error

func nextWriteValue(reader ion.Reader, writer ion.Writer) error {
	next := reader.Next()
	if !next {
		err := reader.Err()
		if err != nil {
			return err
		}
	}

	// TODO: Implement WriteValue logic once writer.WriteValue() is available

	return nil
}

func writeValues(reader ion.Reader, writer ion.Writer) error {
	// TODO: Implement WriteValue logic once writer.WriteValue() is available

	return nil
}
