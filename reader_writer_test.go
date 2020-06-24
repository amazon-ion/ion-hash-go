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
	"reflect"
	"testing"

	"github.com/amzn/ion-go/ion"
)

// This test writes a nested struct {a: {b:1}} where the Ion Writer writes the outer struct
// and the HashWriter writes the inner struct.
// We then read the struct back following a similar pattern where the Ion Reader reads the outer struct
// and the HashReader reads the inner struct.
// We then confirm that the HashReader reads the same hash written by the HashWriter.
func TestFieldNameAsymmetry(t *testing.T) {
	//t.Skip() // Skipping test until reader's IsInStruct logic matches dot net

	buf := bytes.Buffer{}
	writer := ion.NewBinaryWriter(&buf)

	hw, err := NewHashWriter(writer, newIdentityHasherProvider())
	if err != nil {
		t.Fatalf("Expected NewHashWriter() to successfully create a HashWriter; %s", err.Error())
	}

	ionHashWriter, ok := hw.(*hashWriter)
	if !ok {
		t.Fatal("Expected hw to be of type hashWriter")
	}

	// Writing a nested struct: {a:{b:1}}
	// We use the ion writer to write the outer struct (ie. {a:_})
	err = writer.BeginStruct()
	if err != nil {
		t.Errorf("Something went wrong executing writer.BeginStruct(); %s", err.Error())
	}

	err = writer.FieldName("a")
	if err != nil {
		t.Errorf("Something went wrong executing writer.FieldName(\"a\"); %s", err.Error())
	}

	// We use the ion hash writer to write the inner struct (ie. {b:1} inside {a:{b:1}})
	err = ionHashWriter.BeginStruct()
	if err != nil {
		t.Errorf("Something went wrong executing ionHashWriter.BeginStruct(); %s", err.Error())
	}

	err = ionHashWriter.FieldName("b")
	if err != nil {
		t.Errorf("Something went wrong executing ionHashWriter.FieldName(\"b\"); %s", err.Error())
	}

	err = ionHashWriter.WriteInt(1)
	if err != nil {
		t.Errorf("Something went wrong executing ionHashWriter.WriteInt(1); %s", err.Error())
	}

	err = ionHashWriter.EndStruct()
	if err != nil {
		t.Errorf("Something went wrong executing ionHashWriter.EndStruct(); %s", err.Error())
	}

	err = writer.EndStruct()
	if err != nil {
		t.Errorf("Something went wrong executing writer.EndStruct(); %s", err.Error())
	}

	writeHash, err := ionHashWriter.Sum(nil)
	if err != nil {
		t.Fatalf("Something went wrong executing ionHashWriter.Sum(nil); %s", err.Error())
	}

	err = ionHashWriter.Finish()
	if err != nil {
		t.Errorf("Something went wrong executing ionHashWriter.Finish(); %s", err.Error())
	}

	err = writer.Finish()
	if err != nil {
		t.Errorf("Something went wrong executing writer.Finish(); %s", err.Error())
	}

	reader := ion.NewReaderBytes(buf.Bytes())

	hr, err := NewHashReader(reader, newIdentityHasherProvider())
	if err != nil {
		t.Fatalf("Expected NewHashReader() to successfully create a HashReader; %s", err.Error())
	}

	ionHashReader, ok := hr.(*hashReader)
	if !ok {
		t.Fatalf("Expected hr to be of type hashReader")
	}

	// We are reading the nested struct that we just wrote: {a:{b:1}}
	// We use the ion reader to read the outer struct (ie. {a:_})
	if !reader.Next() {
		err = reader.Err()
		if err != nil {
			t.Errorf("Something went wrong executing reader.Next(); %s", err.Error())
		}
	}

	err = reader.StepIn()
	if err != nil {
		t.Errorf("Something went wrong executing reader.StepIn(); %s", err.Error())
	}

	if !reader.Next() {
		err = reader.Err()
		if err != nil {
			t.Errorf("Something went wrong executing reader.Next(); %s", err.Error())
		}
	}

	// We use the ion hash reader to read the inner struct (ie. {b:1} inside {a:{b:1}})
	err = ionHashReader.StepIn()
	if err != nil {
		t.Errorf("Something went wrong executing ionHashReader.StepIn(); %s", err.Error())
	}

	if !ionHashReader.Next() {
		err = ionHashReader.Err()
		if err != nil {
			t.Errorf("Something went wrong executing ionHashReader.Next(); %s", err.Error())
		}
	}

	err = ionHashReader.StepOut()
	if err != nil {
		t.Errorf("Something went wrong executing ionHashReader.StepOut(); %s", err.Error())
	}

	err = reader.StepOut()
	if err != nil {
		t.Errorf("Something went wrong executing reader.StepOut(); %s", err.Error())
	}

	sum, err := ionHashReader.Sum(nil)
	if err != nil {
		t.Fatalf("Something went wrong executing ionHashReader.Sum(nil); %s", err.Error())
	}

	if !reflect.DeepEqual(sum, writeHash) {
		t.Errorf("sum did not match expectation;\n"+
			"Expected sum: %v\n"+
			"Actual sum:   %v",
			writeHash, sum)
	}
}

func TestNoFieldNameInCurrentHash(t *testing.T) {
	//t.Skip() // Skipping test until reader's IsInStruct logic matches dot net

	AssertNoFieldnameInCurrentHash(t, "null", []byte{0x0b, 0x0f, 0x0e})
	AssertNoFieldnameInCurrentHash(t, "false", []byte{0x0b, 0x10, 0x0e})
	AssertNoFieldnameInCurrentHash(t, "5", []byte{0x0b, 0x20, 0x05, 0x0e})
	AssertNoFieldnameInCurrentHash(
		t,
		"2e0",
		[]byte{0x0b, 0x40, 0x40, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0e})
	AssertNoFieldnameInCurrentHash(t, "1234.500", []byte{0x0b, 0x50, 0xc3, 0x12, 0xd6, 0x44, 0x0e})
	AssertNoFieldnameInCurrentHash(t, "hi", []byte{0x0b, 0x70, 0x68, 0x69, 0x0e})
	AssertNoFieldnameInCurrentHash(t, "\"hi\"", []byte{0x0b, 0x80, 0x68, 0x69, 0x0e})
	AssertNoFieldnameInCurrentHash(t, "{{\"hi\"}}", []byte{0x0b, 0x90, 0x68, 0x69, 0x0e})
	AssertNoFieldnameInCurrentHash(t, "{{aGVsbG8=}}", []byte{0x0b, 0xa0, 0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x0e})
	AssertNoFieldnameInCurrentHash(
		t,
		"[1,2,3]",
		[]byte{0x0b, 0xb0, 0x0b, 0x20, 0x01, 0x0e, 0x0b, 0x20, 0x02, 0x0e, 0x0b, 0x20, 0x03, 0x0e, 0x0e})
	AssertNoFieldnameInCurrentHash(
		t,
		"(1 2 3)",
		[]byte{0x0b, 0xc0, 0x0b, 0x20, 0x01, 0x0e, 0x0b, 0x20, 0x02, 0x0e, 0x0b, 0x20, 0x03, 0x0e, 0x0e})
	AssertNoFieldnameInCurrentHash(
		t,
		"{a:1,b:2,c:3}",
		[]byte{
			0x0b, 0xd0, 0x0c, 0x0b, 0x70, 0x61, 0x0c, 0x0e, 0x0c, 0x0b, 0x20, 0x01, 0x0c,
			0x0e, 0x0c, 0x0b, 0x70, 0x62, 0x0c, 0x0e, 0x0c, 0x0b, 0x20, 0x02, 0x0c, 0x0e,
			0x0c, 0x0b, 0x70, 0x63, 0x0c, 0x0e, 0x0c, 0x0b, 0x20, 0x03, 0x0c, 0x0e, 0x0e})
	AssertNoFieldnameInCurrentHash(
		t,
		"hi::7",
		[]byte{0x0b, 0xe0, 0x0b, 0x70, 0x68, 0x69, 0x0e, 0x0b, 0x20, 0x07, 0x0e, 0x0e})
}

func AssertNoFieldnameInCurrentHash(t *testing.T, value string, expectedBytes []byte) {
	var err error

	reader := ion.NewReaderStr(value)

	buf := bytes.Buffer{}
	writer := ion.NewBinaryWriter(&buf)

	err = writer.BeginStruct()
	if err != nil {
		t.Errorf("Something went wrong executing writer.BeginStruct(); %s", err.Error())
	}

	hw, err := NewHashWriter(writer, newIdentityHasherProvider())
	if err != nil {
		t.Fatalf("Expected NewHashWriter() to successfully create a HashWriter; %s", err.Error())
	}

	ionHashWriter, ok := hw.(*hashWriter)
	if !ok {
		t.Fatalf("Expected hw to be of type hashWriter")
	}

	err = ionHashWriter.FieldName("field_name")
	if err != nil {
		t.Errorf("Something went wrong executing ionHashWriter.FieldName(\"field_name\"); %s", err.Error())
	}

	writeFromReaderToWriter(t, reader, ionHashWriter)

	actual, err := ionHashWriter.Sum(nil)
	if err != nil {
		t.Fatalf("Something went wrong executing ionHashWriter.Sum(nil); %s", err.Error())
	}

	if !reflect.DeepEqual(actual, expectedBytes) {
		t.Errorf("sum did not match expectation;\n"+
			"Expected sum: %v\n"+
			"Actual sum:   %v",
			expectedBytes, actual)
	}

	err = writer.EndStruct()
	if err != nil {
		t.Errorf("Something went wrong executing writer.EndStruct(); %s", err.Error())
	}

	err = ionHashWriter.Finish()
	if err != nil {
		t.Errorf("Something went wrong executing ionHashWriter.Finish(); %s", err.Error())
	}

	err = writer.Finish()
	if err != nil {
		t.Errorf("Something went wrong executing writer.Finish(); %s", err.Error())
	}

	reader = ion.NewReaderBytes(buf.Bytes())

	if !reader.Next() {
		err = reader.Err()
		if err != nil {
			t.Errorf("Something went wrong executing reader.Next(); %s", err.Error())
		}
	}

	err = reader.StepIn()
	if err != nil {
		t.Errorf("Something went wrong executing reader.StepIn(); %s", err.Error())
	}

	hr, err := NewHashReader(reader, newIdentityHasherProvider())
	if err != nil {
		t.Fatalf("Expected NewHashReader() to successfully create a HashReader; %s", err.Error())
	}

	ionHashReader := hr.(*hashReader)
	if !ok {
		t.Fatal("Expected hr to be of type hashReader")
	}

	// List
	if !ionHashReader.Next() {
		err = ionHashReader.Err()
		if err != nil {
			t.Errorf("Something went wrong executing ionHashReader.Next(); %s", err.Error())
		}
	}

	// None
	if !ionHashReader.Next() {
		err = ionHashReader.Err()
		if err != nil {
			t.Errorf("Something went wrong executing ionHashReader.Next(); %s", err.Error())
		}
	}

	actualBytes, err := ionHashReader.Sum(nil)
	if err != nil {
		t.Fatalf("Something went wrong executing ionHashReader.Sum(nil); %s", err.Error())
	}

	if !reflect.DeepEqual(expectedBytes, actualBytes) {
		t.Errorf("sum did not match expectation;\n"+
			"Expected sum: %v\n"+
			"Actual sum:   %v",
			expectedBytes, actualBytes)
	}
}
