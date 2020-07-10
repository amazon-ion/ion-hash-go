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
	"testing"

	"github.com/amzn/ion-go/ion"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// This test writes a nested struct {a:{b:1}} where the Ion Writer writes the outer struct
// and the HashWriter writes the inner struct.
// We then read the struct back following a similar pattern where the Ion Reader reads the outer struct
// and the HashReader reads the inner struct.
// We then confirm that the HashReader reads the same hash written by the HashWriter.
func TestFieldNameAsymmetry(t *testing.T) {
	buf := bytes.Buffer{}
	writer := ion.NewBinaryWriter(&buf)

	hw, err := NewHashWriter(writer, newIdentityHasherProvider())
	require.NoError(t, err, "Expected NewHashWriter() to successfully create a HashWriter")

	ionHashWriter, ok := hw.(*hashWriter)
	require.True(t, ok, "Expected hw to be of type hashWriter")

	// Writing a nested struct: {a:{b:1}}
	// We use the ion writer to write the outer struct (ie. {a:_})
	assert.NoError(t, writer.BeginStruct(), "Something went wrong executing writer.BeginStruct()")

	assert.NoError(t, writer.FieldName("a"), "Something went wrong executing writer.FieldName(\"a\")")

	// We use the ion hash writer to write the inner struct (ie. {b:1} inside {a:{b:1}})
	assert.NoError(t, ionHashWriter.BeginStruct(), "Something went wrong executing ionHashWriter.BeginStruct()")

	assert.NoError(t, ionHashWriter.FieldName("b"), "Something went wrong executing ionHashWriter.FieldName(\"b\")")

	assert.NoError(t, ionHashWriter.WriteInt(1), "Something went wrong executing ionHashWriter.WriteInt(1)")

	assert.NoError(t, ionHashWriter.EndStruct(), "Something went wrong executing ionHashWriter.EndStruct()")

	assert.NoError(t, writer.EndStruct(), "Something went wrong executing writer.EndStruct()")

	writeHash, err := ionHashWriter.Sum(nil)
	require.NoError(t, err, "Something went wrong executing ionHashWriter.Sum(nil)")

	assert.NoError(t, ionHashWriter.Finish(), "Something went wrong executing ionHashWriter.Finish()")

	assert.NoError(t, writer.Finish(), "Something went wrong executing writer.Finish()")

	reader := ion.NewReaderBytes(buf.Bytes())

	hr, err := NewHashReader(reader, newIdentityHasherProvider())
	require.NoError(t, err, "Expected NewHashReader() to successfully create a HashReader")

	ionHashReader, ok := hr.(*hashReader)
	require.True(t, ok, "Expected hr to be of type hashReader")

	// We are reading the nested struct that we just wrote: {a:{b:1}}
	// We use the ion reader to read the outer struct (ie. {a:_})
	if !reader.Next() {
		assert.NoError(t, reader.Err(), "Something went wrong executing reader.Next()")
	}

	assert.NoError(t, reader.StepIn(), "Something went wrong executing reader.StepIn()")

	if !reader.Next() {
		assert.NoError(t, reader.Err(), "Something went wrong executing reader.Next()")
	}

	// We use the ion hash reader to read the inner struct (ie. {b:1} inside {a:{b:1}})
	assert.NoError(t, ionHashReader.StepIn(), "Something went wrong executing ionHashReader.StepIn()")

	if !ionHashReader.Next() {
		assert.NoError(t, ionHashReader.Err(), "Something went wrong executing ionHashReader.Next()")
	}

	assert.NoError(t, ionHashReader.StepOut(), "Something went wrong executing ionHashReader.StepOut()")

	assert.NoError(t, reader.StepOut(), "Something went wrong executing reader.StepOut()")

	sum, err := ionHashReader.Sum(nil)
	require.NoError(t, err, "Something went wrong executing ionHashReader.Sum(nil)")

	assert.Equal(t, writeHash, sum, "sum did not match expectation")
}

func TestNoFieldNameInCurrentHash(t *testing.T) {
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

	assert.NoError(t, writer.BeginStruct(), "Something went wrong executing writer.BeginStruct()")

	hw, err := NewHashWriter(writer, newIdentityHasherProvider())
	require.NoError(t, err, "Expected NewHashWriter() to successfully create a HashWriter")

	ionHashWriter, ok := hw.(*hashWriter)
	require.True(t, ok, "Expected hw to be of type hashWriter")

	assert.NoError(t, ionHashWriter.FieldName("field_name"),
		"Something went wrong executing ionHashWriter.FieldName(\"field_name\")")

	writeFromReaderToWriter(t, reader, ionHashWriter)

	actual, err := ionHashWriter.Sum(nil)
	require.NoError(t, err, "Something went wrong executing ionHashWriter.Sum(nil)")

	assert.Equal(t, actual, expectedBytes, "sum did not match expectation")

	assert.NoError(t, writer.EndStruct(), "Something went wrong executing writer.EndStruct()")

	assert.NoError(t, ionHashWriter.Finish(), "Something went wrong executing ionHashWriter.Finish()")

	assert.NoError(t, writer.Finish(), "Something went wrong executing writer.Finish()")

	reader = ion.NewReaderBytes(buf.Bytes())

	if !reader.Next() {
		assert.NoError(t, reader.Err(), "Something went wrong executing reader.Next()")
	}

	assert.NoError(t, reader.StepIn(), "Something went wrong executing reader.StepIn()")

	hr, err := NewHashReader(reader, newIdentityHasherProvider())
	require.NoError(t, err, "Expected NewHashReader() to successfully create a HashReader")

	ionHashReader := hr.(*hashReader)
	require.True(t, ok, "Expected hr to be of type hashReader")

	// List
	if !ionHashReader.Next() {
		assert.NoError(t, ionHashReader.Err(), "Something went wrong executing ionHashReader.Next()")
	}

	// None
	if !ionHashReader.Next() {
		assert.NoError(t, ionHashReader.Err(), "Something went wrong executing ionHashReader.Next()")
	}

	actualBytes, err := ionHashReader.Sum(nil)
	require.NoError(t, err, "Something went wrong executing ionHashReader.Sum(nil)")

	assert.Equal(t, expectedBytes, actualBytes, "sum did not match expectation")
}
