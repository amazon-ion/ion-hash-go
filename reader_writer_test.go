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
	t.Skip() // Skipping test until reader's IsInStruct logic matches dot net

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

	// Writing a nested struct: {a:{b:1}}
	// We use the ion writer to write the outer struct (ie. {a:_})
	err = writer.BeginStruct()
	if err != nil {
		t.Errorf("expected writer.BeginStruct() to execute without errors; %s", err.Error())
	}

	err = writer.FieldName("a")
	if err != nil {
		t.Errorf("expected writer.FieldName(\"a\") to execute without errors; %s", err.Error())
	}

	// We use the ion hash writer to write the inner struct (ie. {b:1} inside {a:{b:1}})
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

	err = writer.EndStruct()
	if err != nil {
		t.Errorf("expected writer.EndStruct() to execute without errors; %s", err.Error())
	}

	writeHash, err := ionHashWriter.Sum(nil)
	if err != nil {
		t.Errorf("expected ionHashWriter.Sum(nil) to execute without errors; %s", err.Error())
	}

	err = ionHashWriter.Finish()
	if err != nil {
		t.Errorf("expected ionHashWriter.Finish() to execute without errors; %s", err.Error())
	}

	err = writer.Finish()
	if err != nil {
		t.Errorf("expected writer.Finish() to execute without errors; %s", err.Error())
	}

	reader := ion.NewReaderBytes(buf.Bytes())

	hr, err := NewHashReader(reader, newIdentityHasherProvider())
	if err != nil {
		t.Fatalf("expected NewHashReader() to successfully create a HashReader; %s", err.Error())
	}

	ionHashReader, ok := hr.(*hashReader)
	if !ok {
		t.Fatalf("expected hr to be of type hashReader")
	}

	// We are reading the nested struct that we just wrote: {a:{b:1}}
	// We use the ion reader to read the outer struct (ie. {a:_})
	if !reader.Next() {
		err = reader.Err()
		if err != nil {
			t.Errorf("expected reader.Next() to execute without errors; %s", err.Error())
		}
	}

	err = reader.StepIn()
	if err != nil {
		t.Errorf("expected reader.StepIn() to execute without errors; %s", err.Error())
	}

	if !reader.Next() {
		err = reader.Err()
		if err != nil {
			t.Errorf("expected reader.Next() to execute without errors; %s", err.Error())
		}
	}

	// We use the ion hash reader to read the inner struct (ie. {b:1} inside {a:{b:1}})
	err = ionHashReader.StepIn()
	if err != nil {
		t.Errorf("expected ionHashReader.StepIn() to execute without errors; %s", err.Error())
	}

	if !ionHashReader.Next() {
		err = ionHashReader.Err()
		if err != nil {
			t.Errorf("expected ionHashReader.Next() to execute without errors; %s", err.Error())
		}
	}

	err = ionHashReader.StepOut()
	if err != nil {
		t.Errorf("expected ionHashReader.StepOut() to execute without errors; %s", err.Error())
	}

	err = reader.StepOut()
	if err != nil {
		t.Errorf("expected reader.StepOut() to execute without errors; %s", err.Error())
	}

	sum, err := ionHashReader.Sum(nil)
	if err != nil {
		t.Fatalf("expected ionHashReader.Sum(nil) to execute without errors; %s", err.Error())
	}

	if !reflect.DeepEqual(sum, writeHash) {
		t.Errorf("expected sum to be %v instead of %v", writeHash, sum)
	}
}

func TestNoFieldNameInCurrentHash(t *testing.T) {
	t.Skip() // Skipping test until reader's IsInStruct logic matches dot net

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
		t.Errorf("expected writer.BeginStruct() to execute without errors; %s", err.Error())
	}

	hw, err := NewHashWriter(writer, newIdentityHasherProvider())
	if err != nil {
		t.Errorf("expected NewHashWriter() to successfully create a HashWriter; %s", err.Error())
	}

	ionHashWriter, ok := hw.(*hashWriter)
	if !ok {
		t.Error("expected ionHashWriter to be of type hashWriter")
	}

	err = ionHashWriter.FieldName("field_name")
	if err != nil {
		t.Errorf("expected ionHashWriter.FieldName(\"field_name\") to execute without errors; %s", err.Error())
	}

	writeToWriterFromReader(t, reader, ionHashWriter)

	actual, err := ionHashWriter.Sum(nil)
	if err != nil {
		t.Errorf("expected ionHashWriter.Sum(nil) to execute without errors; %s", err.Error())
	}

	if !reflect.DeepEqual(actual, expectedBytes) {
		t.Errorf("expected sum to be %v instead of %v", expectedBytes, actual)
	}

	err = writer.EndStruct()
	if err != nil {
		t.Errorf("expected writer.EndStruct() to execute without errors; %s", err.Error())
	}

	err = ionHashWriter.Finish()
	if err != nil {
		t.Errorf("expected ionHashWriter.Finish() to execute without errors; %s", err.Error())
	}

	err = writer.Finish()
	if err != nil {
		t.Errorf("expected writer.Finish() to execute without errors; %s", err.Error())
	}

	reader = ion.NewReaderBytes(buf.Bytes())

	if !reader.Next() {
		err = reader.Err()
		if err != nil {
			t.Errorf("expected reader.Next() to execute without errors; %s", err.Error())
		}
	}

	err = reader.StepIn()
	if err != nil {
		t.Errorf("expected reader.StepIn() to execute without errors; %s", err.Error())
	}

	hr, err := NewHashReader(reader, newIdentityHasherProvider())
	if err != nil {
		t.Errorf("expected NewHashReader() to successfully create a HashReader; %s", err.Error())
	}

	ionHashReader := hr.(*hashReader)
	if !ok {
		t.Fatal("expected hr to be of type hashReader")
	}

	// List
	if !ionHashReader.Next() {
		err = ionHashReader.Err()
		if err != nil {
			t.Errorf("expected ionHashReader.Next() to execute without errors; %s", err.Error())
		}
	}

	// None
	if !ionHashReader.Next() {
		err = ionHashReader.Err()
		if err != nil {
			t.Errorf("expected ionHashReader.Next() to execute without errors; %s", err.Error())
		}
	}

	actualBytes, err := ionHashReader.Sum(nil)
	if err != nil {
		t.Errorf("expected ionHashReader.Sum(nil) to execute without errors; %s", err.Error())
	}

	if !reflect.DeepEqual(expectedBytes, actualBytes) {
		t.Errorf("expected sum to be %v instead of %v", expectedBytes, actualBytes)
	}
}

// Read all the values in the reader and write them in the writer
func writeToWriterFromReader(t *testing.T, reader ion.Reader, writer ion.Writer) {
	for reader.Next() {
		name := reader.FieldName()
		if name != "" {
			err := writer.FieldName(name)
			if err != nil {
				t.Fatalf("expected writer.FieldName(name) to execute without errors; %s", err.Error())
			}
		}

		an := reader.Annotations()
		if len(an) > 0 {
			err := writer.Annotations(an...)
			if err != nil {
				t.Fatalf("expected writer.Annotations(an...) to execute without errors; %s", err.Error())
			}
		}

		currentType := reader.Type()
		if reader.IsNull() {
			err := writer.WriteNullType(currentType)
			if err != nil {
				t.Fatalf("expected writer.WriteNullType(currentType) to execute without errors; %s", err.Error())
			}
			return
		}

		switch currentType {
		case ion.BoolType:
			val, err := reader.BoolValue()
			if err != nil {
				t.Errorf("Something went wrong when reading Boolean value; %s", err.Error())
			}
			err = writer.WriteBool(val)
			if err != nil {
				t.Errorf("Something went wrong when writing Boolean value; %s", err.Error())
			}

		case ion.IntType:
			val, err := reader.Int64Value()
			if err != nil {
				t.Errorf("Something went wrong when reading Int value; %s", err.Error())
			}
			err = writer.WriteInt(val)
			if err != nil {
				t.Errorf("Something went wrong when writing Int value; %s", err.Error())
			}

		case ion.FloatType:
			val, err := reader.FloatValue()
			if err != nil {
				t.Errorf("Something went wrong when reading Float value; %s", err.Error())
			}
			err = writer.WriteFloat(val)
			if err != nil {
				t.Errorf("Something went wrong when writing Float value; %s", err.Error())
			}

		case ion.DecimalType:
			val, err := reader.DecimalValue()
			if err != nil {
				t.Errorf("Something went wrong when reading Decimal value; %s", err.Error())
			}
			err = writer.WriteDecimal(val)
			if err != nil {
				t.Errorf("Something went wrong when writing Decimal value; %s", err.Error())
			}

		case ion.TimestampType:
			val, err := reader.TimeValue()
			if err != nil {
				t.Errorf("Something went wrong when reading Timestamp value; %s", err.Error())
			}
			err = writer.WriteTimestamp(val)
			if err != nil {
				t.Errorf("Something went wrong when writing Timestamp value; %s", err.Error())
			}

		case ion.SymbolType:
			val, err := reader.StringValue()
			if err != nil {
				t.Errorf("Something went wrong when reading Symbol value; %s", err.Error())
			}
			err = writer.WriteSymbol(val)
			if err != nil {
				t.Errorf("Something went wrong when writing Symbol value; %s", err.Error())
			}

		case ion.StringType:
			val, err := reader.StringValue()
			if err != nil {
				t.Errorf("Something went wrong when reading String value; %s", err.Error())
			}
			err = writer.WriteString(val)
			if err != nil {
				t.Errorf("Something went wrong when writing String value; %s", err.Error())
			}

		case ion.ClobType:
			val, err := reader.ByteValue()
			if err != nil {
				t.Errorf("Something went wrong when reading Clob value; %s", err.Error())
			}
			err = writer.WriteClob(val)
			if err != nil {
				t.Errorf("Something went wrong when writing Clob value; %s", err.Error())
			}

		case ion.BlobType:
			val, err := reader.ByteValue()
			if err != nil {
				t.Errorf("Something went wrong when reading Blob value; %s", err.Error())
			}
			err = writer.WriteBlob(val)
			if err != nil {
				t.Errorf("Something went wrong when writing Blob value; %s", err.Error())
			}

		case ion.SexpType:
			err := reader.StepIn()
			if err != nil {
				t.Fatalf("expected reader.StepIn() to execute without errors; %s", err.Error())
			}
			err = writer.BeginSexp()
			if err != nil {
				t.Fatalf("expected writer.BeginSexp() to execute without errors; %s", err.Error())
			}
			writeToWriterFromReader(t, reader, writer)
			err = reader.StepOut()
			if err != nil {
				t.Fatalf("expected reader.StepOut() to execute without errors; %s", err.Error())
			}
			err = writer.EndSexp()
			if err != nil {
				t.Fatalf("expected writer.EndSexp() to execute without errors; %s", err.Error())
			}

		case ion.ListType:
			err := reader.StepIn()
			if err != nil {
				t.Fatalf("expected reader.StepIn() to execute without errors; %s", err.Error())
			}
			err = writer.BeginList()
			if err != nil {
				t.Fatalf("expected writer.BeginList() to execute without errors; %s", err.Error())
			}
			writeToWriterFromReader(t, reader, writer)
			err = reader.StepOut()
			if err != nil {
				t.Fatalf("expected reader.StepOut() to execute without errors; %s", err.Error())
			}
			err = writer.EndList()
			if err != nil {
				t.Fatalf("expected writer.EndList() to execute without errors; %s", err.Error())
			}

		case ion.StructType:
			err := reader.StepIn()
			if err != nil {
				t.Fatalf("expected reader.StepIn() to execute without errors; %s", err.Error())
			}
			err = writer.BeginStruct()
			if err != nil {
				t.Fatalf("expected writer.BeginStruct() to execute without errors; %s", err.Error())
			}
			writeToWriterFromReader(t, reader, writer)
			err = reader.StepOut()
			if err != nil {
				t.Fatalf("expected reader.StepOut() to execute without errors; %s", err.Error())
			}
			err = writer.EndStruct()
			if err != nil {
				t.Fatalf("expected writer.EndStruct() to execute without errors; %s", err.Error())
			}
		}
	}

	if reader.Err() != nil {
		t.Errorf("expected reader.Next() to execute without errors; %s", reader.Err().Error())
	}
}
