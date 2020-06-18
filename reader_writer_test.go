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
	"reflect"
	"testing"

	"github.com/amzn/ion-go/ion"
)

func TestNoFieldNameInCurrentHash(t *testing.T) {
	t.Skip()

	err := AssertNoFieldnameInCurrentHash("null", []byte{0x0b, 0x0f, 0x0e})
	if err != nil {
		t.Errorf(err.Error())
	}

	err = AssertNoFieldnameInCurrentHash("false", []byte{0x0b, 0x10, 0x0e})
	if err != nil {
		t.Errorf(err.Error())
	}

	err = AssertNoFieldnameInCurrentHash("5", []byte{0x0b, 0x20, 0x05, 0x0e})
	if err != nil {
		t.Errorf(err.Error())
	}

	err = AssertNoFieldnameInCurrentHash(
		"2e0",
		[]byte{0x0b, 0x40, 0x40, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0e})
	if err != nil {
		t.Errorf(err.Error())
	}

	err = AssertNoFieldnameInCurrentHash("1234.500", []byte{0x0b, 0x50, 0xc3, 0x12, 0xd6, 0x44, 0x0e})
	if err != nil {
		t.Errorf(err.Error())
	}

	err = AssertNoFieldnameInCurrentHash("hi", []byte{0x0b, 0x70, 0x68, 0x69, 0x0e})
	if err != nil {
		t.Errorf(err.Error())
	}

	err = AssertNoFieldnameInCurrentHash("\"hi\"", []byte{0x0b, 0x80, 0x68, 0x69, 0x0e})
	if err != nil {
		t.Errorf(err.Error())
	}

	err = AssertNoFieldnameInCurrentHash("{{\"hi\"}}", []byte{0x0b, 0x90, 0x68, 0x69, 0x0e})
	if err != nil {
		t.Errorf(err.Error())
	}

	err = AssertNoFieldnameInCurrentHash("{{aGVsbG8=}}", []byte{0x0b, 0xa0, 0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x0e})
	if err != nil {
		t.Errorf(err.Error())
	}

	err = AssertNoFieldnameInCurrentHash(
		"[1,2,3]",
		[]byte{0x0b, 0xb0, 0x0b, 0x20, 0x01, 0x0e, 0x0b, 0x20, 0x02, 0x0e, 0x0b, 0x20, 0x03, 0x0e, 0x0e})
	if err != nil {
		t.Errorf(err.Error())
	}

	err = AssertNoFieldnameInCurrentHash(
		"(1 2 3)",
		[]byte{0x0b, 0xc0, 0x0b, 0x20, 0x01, 0x0e, 0x0b, 0x20, 0x02, 0x0e, 0x0b, 0x20, 0x03, 0x0e, 0x0e})
	if err != nil {
		t.Errorf(err.Error())
	}

	err = AssertNoFieldnameInCurrentHash(
		"{a:1,b:2,c:3}",
		[]byte{
			0x0b, 0xd0, 0x0c, 0x0b, 0x70, 0x61, 0x0c, 0x0e, 0x0c, 0x0b, 0x20, 0x01, 0x0c,
			0x0e, 0x0c, 0x0b, 0x70, 0x62, 0x0c, 0x0e, 0x0c, 0x0b, 0x20, 0x02, 0x0c, 0x0e,
			0x0c, 0x0b, 0x70, 0x63, 0x0c, 0x0e, 0x0c, 0x0b, 0x20, 0x03, 0x0c, 0x0e, 0x0e})
	if err != nil {
		t.Errorf(err.Error())
	}

	err = AssertNoFieldnameInCurrentHash(
		"hi::7",
		[]byte{0x0b, 0xe0, 0x0b, 0x70, 0x68, 0x69, 0x0e, 0x0b, 0x20, 0x07, 0x0e, 0x0e})
	if err != nil {
		t.Errorf(err.Error())
	}
}

func TestFieldNameAsymmetry(t *testing.T) {
	t.Skip()

	buf := bytes.Buffer{}
	writer := ion.NewBinaryWriter(&buf)

	hw, err := NewHashWriter(writer, newIdentityHasherProvider())
	if err != nil {
		t.Fatalf("expected NewHashWriter() to successfully create a HashWriter; %s", err.Error())
	}

	ionHashWriter, ok := hw.(*hashWriter)
	if !ok {
		t.Fatal("expected ionHashWriter to be of type hashWriter")
	}

	// A nested struct: {a:{b:1}}
	err = writer.BeginStruct()
	if err != nil {
		t.Errorf("expected writer.BeginStruct() to execute without errors; %s", err.Error())
	}

	err = writer.FieldName("a")
	if err != nil {
		t.Errorf("expected writer.FieldName(\"a\") to execute without errors; %s", err.Error())
	}

	err = ionHashWriter.BeginStruct()
	if err != nil {
		t.Errorf("expected ionHashWriter.BeginStruct() to execute without errors; %s", err.Error())
	}

	err = ionHashWriter.FieldName("b")
	if err != nil {
		t.Errorf("expected ionHashWriter.FieldName(\"b\") to execute without errors; %s", err.Error())
	}

	err = ionHashWriter.WriteInt(1)
	if err != nil {
		t.Errorf("expected ionHashWriter.WriteInt(1) to execute without errors; %s", err.Error())
	}

	err = ionHashWriter.EndStruct()
	if err != nil {
		t.Errorf("expected ionHashWriter.EndStruct() to execute without errors; %s", err.Error())
	}

	writeHash, err := ionHashWriter.Sum(nil)
	if err != nil {
		t.Errorf("expected ionHashWriter.Sum(nil) to execute without errors; %s", err.Error())
	}

	err = ionHashWriter.Finish()
	if err != nil {
		t.Errorf("expected ionHashWriter.stepOut() to execute without errors; %s", err.Error())
	}

	err = writer.EndStruct()
	if err != nil {
		t.Errorf("expected writer.EndStruct() to execute without errors; %s", err.Error())
	}

	err = writer.Finish()
	if err != nil {
		t.Errorf("expected writer.Finish() to execute without errors; %s", err.Error())
	}

	ionHashReader, err := NewHashReader(ion.NewReaderBytes(buf.Bytes()), newIdentityHasherProvider())
	if err != nil {
		t.Fatalf("expected NewHashReader() to successfully create a HashReader; %s", err.Error())
	}

	if !ionHashReader.Next() {
		err = ionHashReader.Err()
		if err != nil {
			t.Errorf("expected ionHashReader.Next() to execute without errors; %s", err.Error())
		}
	}

	if !ionHashReader.Next() {
		err = ionHashReader.Err()
		if err != nil {
			t.Errorf("expected ionHashReader.Next() to execute without errors; %s", err.Error())
		}
	}

	sum, err := ionHashReader.Sum(nil)
	if err != nil {
		t.Fatalf("expected Sum(nil) to execute without errors; %s", err.Error())
	}

	if !reflect.DeepEqual(sum, writeHash) {
		t.Errorf("expected sum to be %v instead of %v", writeHash, sum)
	}
}

func AssertNoFieldnameInCurrentHash(value string, expectedBytes []byte) error {
	var err error

	reader := ion.NewReaderStr(value)

	if !reader.Next() {
		err = reader.Err()
		if err != nil {
			return fmt.Errorf("expected reader.Next() to execute without errors; %s", err.Error())
		}
	}

	buf := bytes.Buffer{}
	writer := ion.NewBinaryWriter(&buf)

	err = writer.BeginStruct()
	if err != nil {
		return fmt.Errorf("expected writer.BeginStruct() to execute without errors; %s", err.Error())
	}

	hw, err := NewHashWriter(writer, newIdentityHasherProvider())
	if err != nil {
		return fmt.Errorf("expected NewHashWriter() to successfully create a HashWriter; %s", err.Error())
	}

	ionHashWriter, ok := hw.(*hashWriter)
	if !ok {
		return fmt.Errorf("expected ionHashWriter to be of type hashWriter")
	}

	err = ionHashWriter.FieldName("field_name")
	if err != nil {
		return fmt.Errorf("expected ionHashWriter.FieldName(\"field_name\") to execute without errors; %s", err.Error())
	}

	//ionHashWriter.writeValue(reader)

	actual, err := ionHashWriter.Sum(nil)
	if err != nil {
		return fmt.Errorf("expected ionHashWriter.Sum(nil) to execute without errors; %s", err.Error())
	}

	if !reflect.DeepEqual(actual, expectedBytes) {
		return fmt.Errorf("expected sum to be %v instead of %v", expectedBytes, actual)
	}

	err = writer.EndStruct()
	if err != nil {
		return fmt.Errorf("expected writer.EndStruct() to execute without errors; %s", err.Error())
	}

	err = ionHashWriter.Finish()
	if err != nil {
		return fmt.Errorf("expected ionHashWriter.Finish() to execute without errors; %s", err.Error())
	}

	err = writer.Finish()
	if err != nil {
		return fmt.Errorf("expected writer.Finish() to execute without errors; %s", err.Error())
	}

	reader = ion.NewReaderBytes(buf.Bytes())

	if !reader.Next() {
		err = reader.Err()
		if err != nil {
			return fmt.Errorf("expected reader.Next() to execute without errors; %s", err.Error())
		}
	}

	err = reader.StepIn()
	if err != nil {
		return fmt.Errorf("expected reader.StepIn() to execute without errors; %s", err.Error())
	}

	ionHashReader, err := NewHashReader(reader, newIdentityHasherProvider())
	if err != nil {
		return fmt.Errorf("expected NewHashReader() to successfully create a HashReader; %s", err.Error())
	}

	// List
	if !ionHashReader.Next() {
		err = ionHashReader.Err()
		if err != nil {
			return fmt.Errorf("expected ionHashReader.Next() to execute without errors; %s", err.Error())
		}
	}

	// None
	if !ionHashReader.Next() {
		err = ionHashReader.Err()
		if err != nil {
			return fmt.Errorf("expected ionHashReader.Next() to execute without errors; %s", err.Error())
		}
	}

	actualBytes, err := ionHashReader.Sum(nil)
	if err != nil {
		return fmt.Errorf("expected ionHashReader.Sum(nil) to execute without errors; %s", err.Error())
	}

	if !reflect.DeepEqual(expectedBytes, actualBytes) {
		return fmt.Errorf("expected sum to be %v instead of %v", expectedBytes, actualBytes)
	}

	return nil
}
