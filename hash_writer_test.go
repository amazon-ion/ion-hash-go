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
	"math/big"
	"strings"
	"testing"
	"time"

	"github.com/amzn/ion-go/ion"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWriteNull(t *testing.T) {
	str := strings.Builder{}
	hw, err := NewHashWriter(ion.NewTextWriter(&str), newIdentityHasherProvider())
	require.NoError(t, err, "Expected NewHashWriter() to successfully create a HashWriter")

	ionHashWriter, ok := hw.(*hashWriter)
	require.True(t, ok, "Expected hw to be of type hashWriter")

	// Null
	assert.NoError(t, ionHashWriter.WriteNull(), "Something went wrong executing ionHashWriter.WriteNull()")

	sum, err := ionHashWriter.Sum(nil)
	require.NoError(t, err, "Something went wrong executing ionHashWriter.Sum(nil)")

	assert.Equal(t, []byte{0x0b, 0x0f, 0x0e}, sum, "sum did not match expectation")

	// Null Type
	assert.NoError(t, ionHashWriter.WriteNullType(ion.NullType),
		"Something went wrong executing ionHashWriter.WriteNullType(ion.NullType)")

	sum, err = ionHashWriter.Sum(nil)
	require.NoError(t, err, "Something went wrong executing ionHashWriter.Sum(nil)")

	assert.Equal(t, []byte{0x0b, 0x0f, 0x0e}, sum, "sum did not match expectation")

	// Bool
	assert.NoError(t, ionHashWriter.WriteNullType(ion.BoolType),
		"Something went wrong executing ionHashWriter.WriteNullType(ion.BoolType)")

	sum, err = ionHashWriter.Sum(nil)
	require.NoError(t, err, "Something went wrong executing ionHashWriter.Sum(nil)")

	assert.Equal(t, []byte{0x0b, 0x1f, 0x0e}, sum, "sum did not match expectation")

	// Int
	assert.NoError(t, ionHashWriter.WriteNullType(ion.IntType),
		"Something went wrong executing ionHashWriter.WriteNullType(ion.IntType)")

	sum, err = ionHashWriter.Sum(nil)
	require.NoError(t, err, "Something went wrong executing ionHashWriter.Sum(nil)")

	assert.Equal(t, []byte{0x0b, 0x2f, 0x0e}, sum, "sum did not match expectation")

	// Float
	assert.NoError(t, ionHashWriter.WriteNullType(ion.FloatType),
		"Something went wrong executing ionHashWriter.WriteNullType(ion.FloatType)")

	sum, err = ionHashWriter.Sum(nil)
	require.NoError(t, err, "Something went wrong executing ionHashWriter.Sum(nil)")

	assert.Equal(t, []byte{0x0b, 0x4f, 0x0e}, sum, "sum did not match expectation")

	// Decimal
	assert.NoError(t, ionHashWriter.WriteNullType(ion.DecimalType),
		"Something went wrong executing ionHashWriter.WriteNullType(ion.DecimalType)")

	sum, err = ionHashWriter.Sum(nil)
	require.NoError(t, err, "Something went wrong executing ionHashWriter.Sum(nil)")

	assert.Equal(t, []byte{0x0b, 0x5f, 0x0e}, sum, "sum did not match expectation")

	// Timestamp
	assert.NoError(t, ionHashWriter.WriteNullType(ion.TimestampType),
		"Something went wrong executing ionHashWriter.WriteNullType(ion.TimestampType)")

	sum, err = ionHashWriter.Sum(nil)
	require.NoError(t, err, "Something went wrong executing ionHashWriter.Sum(nil)")

	assert.Equal(t, []byte{0x0b, 0x6f, 0x0e}, sum, "sum did not match expectation")

	// Symbol
	assert.NoError(t, ionHashWriter.WriteNullType(ion.SymbolType),
		"Something went wrong executing ionHashWriter.WriteNullType(ion.SymbolType)")

	sum, err = ionHashWriter.Sum(nil)
	require.NoError(t, err, "Something went wrong executing ionHashWriter.Sum(nil)")

	assert.Equal(t, []byte{0x0b, 0x7f, 0x0e}, sum, "sum did not match expectation")

	// String
	assert.NoError(t, ionHashWriter.WriteNullType(ion.StringType),
		"Something went wrong executing ionHashWriter.WriteNullType(ion.StringType)")

	sum, err = ionHashWriter.Sum(nil)
	require.NoError(t, err, "Something went wrong executing ionHashWriter.Sum(nil)")

	assert.Equal(t, []byte{0x0b, 0x8f, 0x0e}, sum, "sum did not match expectation")

	// Clob
	assert.NoError(t, ionHashWriter.WriteNullType(ion.ClobType),
		"Something went wrong executing ionHashWriter.WriteNullType(ion.ClobType)")

	sum, err = ionHashWriter.Sum(nil)
	require.NoError(t, err, "Something went wrong executing ionHashWriter.Sum(nil)")

	assert.Equal(t, []byte{0x0b, 0x9f, 0x0e}, sum, "sum did not match expectation")

	// Blob
	assert.NoError(t, ionHashWriter.WriteNullType(ion.BlobType),
		"Something went wrong executing ionHashWriter.WriteNullType(ion.BlobType)")

	sum, err = ionHashWriter.Sum(nil)
	require.NoError(t, err, "Something went wrong executing ionHashWriter.Sum(nil)")

	assert.Equal(t, []byte{0x0b, 0xaf, 0x0e}, sum, "sum did not match expectation")

	// List
	assert.NoError(t, ionHashWriter.WriteNullType(ion.ListType),
		"Something went wrong executing ionHashWriter.WriteNullType(ion.ListType)")

	sum, err = ionHashWriter.Sum(nil)
	require.NoError(t, err, "Something went wrong executing ionHashWriter.Sum(nil)")

	assert.Equal(t, []byte{0x0b, 0xbf, 0x0e}, sum, "sum did not match expectation")

	// Sexp
	assert.NoError(t, ionHashWriter.WriteNullType(ion.SexpType),
		"Something went wrong executing ionHashWriter.WriteNullType(ion.SexpType)")

	sum, err = ionHashWriter.Sum(nil)
	require.NoError(t, err, "Something went wrong executing ionHashWriter.Sum(nil)")

	assert.Equal(t, []byte{0x0b, 0xcf, 0x0e}, sum, "sum did not match expectation")

	// Struct
	assert.NoError(t, ionHashWriter.WriteNullType(ion.StructType),
		"Something went wrong executing ionHashWriter.WriteNullType(ion.StructType)")

	sum, err = ionHashWriter.Sum(nil)
	require.NoError(t, err, "Something went wrong executing ionHashWriter.Sum(nil)")

	assert.Equal(t, []byte{0x0b, 0xdf, 0x0e}, sum, "sum did not match expectation")

	assert.NoError(t, ionHashWriter.Finish(), "Something went wrong executing ionHashWriter.Finish()")

	// We're comparing splits because str.String() uses a cumbersome '\n' separator
	expected := strings.Split("null null.null null.bool null.int null.float null.decimal null.timestamp "+
		"null.symbol null.string null.clob null.blob null.list null.sexp null.struct ", " ")
	actual := strings.Split(str.String(), "\n")

	assert.Equal(t, expected, actual, "str.String() did not match expectation")
}

func TestWriteScalars(t *testing.T) {
	str := strings.Builder{}
	hw, err := NewHashWriter(ion.NewTextWriter(&str), newIdentityHasherProvider())
	require.NoError(t, err, "Expected NewHashWriter() to successfully create a HashWriter")

	ionHashWriter, ok := hw.(*hashWriter)
	require.True(t, ok, "Expected hw to be of type hashWriter")

	sum, err := ionHashWriter.Sum(nil)
	require.NoError(t, err, "Something went wrong executing ionHashWriter.Sum(nil)")

	assert.Equal(t, []byte{}, sum, "sum did not match expectation")

	// Bool
	assert.NoError(t, ionHashWriter.WriteBool(true), "Something went wrong executing ionHashWriter.WriteBool(true)")

	sum, err = ionHashWriter.Sum(nil)
	require.NoError(t, err, "Something went wrong executing ionHashWriter.Sum(nil)")

	assert.Equal(t, []byte{0x0b, 0x11, 0x0e}, sum, "sum did not match expectation")

	// Uint
	assert.NoError(t, ionHashWriter.WriteUint(5), "Something went wrong executing ionHashWriter.WriteUint(5)")

	sum, err = ionHashWriter.Sum(nil)
	require.NoError(t, err, "Something went wrong executing ionHashWriter.Sum(nil)")

	assert.Equal(t, []byte{0x0b, 0x20, 0x05, 0x0e}, sum, "sum did not match expectation")

	// Int
	assert.NoError(t, ionHashWriter.WriteInt(-5), "Something went wrong executing ionHashWriter.WriteInt(5)")

	sum, err = ionHashWriter.Sum(nil)
	require.NoError(t, err, "Something went wrong executing ionHashWriter.Sum(nil)")

	assert.Equal(t, []byte{0x0b, 0x30, 0x05, 0x0e}, sum, "sum did not match expectation")

	// Big Int
	bigInt := big.NewInt(123456789)
	assert.NoError(t, ionHashWriter.WriteBigInt(bigInt), "Something went wrong executing ionHashWriter.WriteBigInt(bigInt)")

	sum, err = ionHashWriter.Sum(nil)
	require.NoError(t, err, "Something went wrong executing ionHashWriter.Sum(nil)")

	assert.Equal(t, []byte{0xb, 0x20, 0x7, 0x5b, 0xcd, 0x15, 0xe}, sum, "sum did not match expectation")

	// Float
	assert.NoError(t, ionHashWriter.WriteFloat(3.14), "Something went wrong executing ionHashWriter.WriteFloat(3.14)")

	sum, err = ionHashWriter.Sum(nil)
	require.NoError(t, err, "Something went wrong executing ionHashWriter.Sum(nil)")

	assert.Equal(t, []byte{0x0b, 0x40, 0x40, 0x09, 0x1e, 0xb8, 0x51, 0xeb, 0x85, 0x1f, 0x0e}, sum,
		"sum did not match expectation")

	// Decimal
	dec, err := ion.ParseDecimal("1234.56789")
	assert.NoError(t, err, "Something went wrong executing ion.ParseDecimal(\"1234.56789\")")
	assert.NoError(t, ionHashWriter.WriteDecimal(dec), "Something went wrong executing ionHashWriter.WriteDecimal(dec)")

	sum, err = ionHashWriter.Sum(nil)
	require.NoError(t, err, "Something went wrong executing ionHashWriter.Sum(nil)")

	assert.Equal(t, []byte{0x0b, 0x50, 0xc5, 0x07, 0x5b, 0xcd, 0x15, 0x0e}, sum, "sum did not match expectation")

	// Timestamp
	assert.NoError(t, ionHashWriter.WriteTimestamp(time.Date(1941, time.December, 7, 18, 0, 0, 0, time.UTC)),
		"Something went wrong executing ionHashWriter.WriteTimestamp(time.Date(...))")

	sum, err = ionHashWriter.Sum(nil)
	require.NoError(t, err, "Something went wrong executing ionHashWriter.Sum(nil)")

	assert.Equal(t, []byte{0x0b, 0x60, 0x80, 0x0f, 0x95, 0x8c, 0x87, 0x92, 0x80, 0x80, 0x0e}, sum,
		"sum did not match expectation")

	// Symbol
	assert.NoError(t, ionHashWriter.WriteSymbol("symbol"),
		"Something went wrong executing ionHashWriter.WriteSymbol(\"symbol\")")

	sum, err = ionHashWriter.Sum(nil)
	require.NoError(t, err, "Something went wrong executing ionHashWriter.Sum(nil)")

	assert.Equal(t, []byte{0x0b, 0x70, 0x73, 0x79, 0x6d, 0x62, 0x6f, 0x6c, 0x0e}, sum, "sum did not match expectation")

	// String
	assert.NoError(t, ionHashWriter.WriteString("string"),
		"Something went wrong executing ionHashWriter.WriteString(\"string\")")

	sum, err = ionHashWriter.Sum(nil)
	require.NoError(t, err, "Something went wrong executing ionHashWriter.Sum(nil)")

	assert.Equal(t, []byte{0x0b, 0x80, 0x73, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x0e}, sum, "sum did not match expectation")

	// Clob
	err = ionHashWriter.WriteClob(
		[]byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f})
	assert.NoError(t, err, "Something went wrong executing ionHashWriter.WriteClob(...)")

	sum, err = ionHashWriter.Sum(nil)
	require.NoError(t, err, "Something went wrong executing ionHashWriter.Sum(nil)")

	assert.Equal(t, []byte{0x0b, 0x90, 0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08,
		0x09, 0x0a, 0x0c, 0x0b, 0x0c, 0x0c, 0x0d, 0x0c, 0x0e, 0x0f, 0x0e}, sum, "sum did not match expectation")

	// Blob
	err = ionHashWriter.WriteBlob(
		[]byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f})
	assert.NoError(t, err, "Something went wrong executing ionHashWriter.WriteBlob(...)")

	sum, err = ionHashWriter.Sum(nil)
	require.NoError(t, err, "Something went wrong executing ionHashWriter.Sum(nil)")

	expectedSum := []byte{0x0b, 0xa0, 0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08,
		0x09, 0x0a, 0x0c, 0x0b, 0x0c, 0x0c, 0x0d, 0x0c, 0x0e, 0x0f, 0x0e}

	assert.Equal(t, expectedSum, sum, "sum did not match expectation")

	assert.NoError(t, ionHashWriter.Finish(), "Something went wrong executing ionHashWriter.Finish()")

	// We're comparing splits because str.String() uses a cumbersome '\n' separator
	expected := strings.Split("true 5 -5 123456789 3.14e+0 1234.56789 1941-12-07T18:00:00Z symbol \"string\" "+
		"{{\"\\0\\x01\\x02\\x03\\x04\\x05\\x06\\a\\b\\t\\n\\v\\f\\r\\x0E\\x0F\"}}"+
		" {{AAECAwQFBgcICQoLDA0ODw==}} ", " ")
	actual := strings.Split(str.String(), "\n")

	assert.Equal(t, expected, actual, "str.String() did not match expectation")
}

func TestWriteContainers(t *testing.T) {
	str := strings.Builder{}
	hw, err := NewHashWriter(ion.NewTextWriter(&str), newIdentityHasherProvider())
	require.NoError(t, err, "Expected NewHashWriter() to successfully create a HashWriter")

	ionHashWriter, ok := hw.(*hashWriter)
	require.True(t, ok, "Expected hw to be of type hashWriter")

	sum, err := ionHashWriter.Sum(nil)
	require.NoError(t, err, "Something went wrong executing ionHashWriter.Sum(nil)")

	assert.Equal(t, []byte{}, sum, "sum did not match expectation")

	assert.NoError(t, ionHashWriter.BeginList(), "Something went wrong executing ionHashWriter.BeginList()")

	sum, err = ionHashWriter.Sum(nil)
	assert.Error(t, err, "Expected ionHashWriter.Sum(nil) to return an error")
	assert.IsType(t, &InvalidOperationError{}, err, "Expected ionHashWriter.Sum(nil) to return an InvalidOperationError")

	assert.NoError(t, ionHashWriter.WriteBool(true), "Something went wrong executing ionHashWriter.WriteBool(true)")

	sum, err = ionHashWriter.Sum(nil)
	assert.Error(t, err, "Expected ionHashWriter.Sum(nil) to return an error")
	assert.IsType(t, &InvalidOperationError{}, err, "Expected ionHashWriter.Sum(nil) to return an InvalidOperationError")

	assert.NoError(t, ionHashWriter.EndList(), "Something went wrong executing ionHashWriter.EndList()")

	sum, err = ionHashWriter.Sum(nil)
	require.NoError(t, err, "Something went wrong executing ionHashWriter.Sum(nil)")

	assert.Equal(t, []byte{0x0b, 0xb0, 0x0b, 0x11, 0x0e, 0x0e}, sum, "sum did not match expectation")

	assert.False(t, ionHashWriter.IsInStruct())

	assert.NoError(t, ionHashWriter.BeginStruct(), "Something went wrong executing ionHashWriter.BeginStruct()")

	assert.True(t, ionHashWriter.IsInStruct())

	assert.NoError(t, ionHashWriter.FieldName("hello"),
		"Something went wrong executing ionHashWriter.FieldName(\"hello\")")

	assert.NoError(t, ionHashWriter.Annotation("ion"),
		"Something went wrong executing ionHashWriter.Annotation(\"ion\")")

	assert.NoError(t, ionHashWriter.Annotation("hash"),
		"Something went wrong executing ionHashWriter.Annotation(\"hash\")")

	assert.NoError(t, ionHashWriter.WriteSymbol("world"),
		"Something went wrong executing ionHashWriter.WriteSymbol(\"world\")")

	assert.NoError(t, ionHashWriter.EndStruct(), "Something went wrong executing ionHashWriter.EndStruct()")

	assert.False(t, ionHashWriter.IsInStruct())

	sum, err = ionHashWriter.Sum(nil)
	require.NoError(t, err, "Something went wrong executing ionHashWriter.Sum(nil)")

	expectedSum := []byte{0x0b, 0xd0,
		0x0c, 0x0b, 0x70, 0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x0c, 0x0e, // hello:
		0x0c, 0x0b, 0xe0,
		0x0c, 0x0b, 0x70, 0x69, 0x6f, 0x6e, 0x0c, 0x0e, // ion::
		0x0c, 0x0b, 0x70, 0x68, 0x61, 0x73, 0x68, 0x0c, 0x0e, // hash::
		0x0c, 0x0b, 0x70, 0x77, 0x6f, 0x72, 0x6c, 0x64, 0x0c, 0x0e, // world
		0x0c, 0x0e,
		0x0e}

	assert.Equal(t, expectedSum, sum, "sum did not match expectation")

	assert.NoError(t, ionHashWriter.Finish(), "Something went wrong executing ionHashWriter.Finish()")

	// We're comparing splits because str.String() uses a cumbersome '\n' separator
	expected := strings.Split("[true] {hello:ion::hash::world} ", " ")
	actual := strings.Split(str.String(), "\n")

	assert.Equal(t, expected, actual, "sum did not match expectation")
}

func TestExtraEndContainer(t *testing.T) {
	str := strings.Builder{}
	hw, err := NewHashWriter(ion.NewTextWriter(&str), newIdentityHasherProvider())
	require.NoError(t, err, "Expected NewHashWriter() to successfully create a HashWriter")

	ionHashWriter, ok := hw.(*hashWriter)
	require.True(t, ok, "Expected hw to be of type hashWriter")

	err = ionHashWriter.EndList()
	assert.Error(t, err, "Expected ionHashWriter.EndList() to return an error")
	assert.IsType(t, &InvalidOperationError{}, err, "Expected ionHashWriter.EndList() to return an InvalidOperationError")

	err = ionHashWriter.EndSexp()
	assert.Error(t, err, "Expected ionHashWriter.EndSexp() to return an error")
	assert.IsType(t, &InvalidOperationError{}, err, "Expected ionHashWriter.EndSexp() to return an InvalidOperationError")

	err = ionHashWriter.EndStruct()
	assert.Error(t, err, "Expected ionHashWriter.EndStruct() to return an error")
	assert.IsType(t, &InvalidOperationError{}, err, "Expected ionHashWriter.EndStruct() to return an InvalidOperationError")
}

func TestIonWriterContractWriteValue(t *testing.T) {
	// Skipping test until FieldNameSymbol logic is available.
	// Test currently fails with empty field name ie. {'':1}
	t.Skip()

	file, err := ioutil.ReadFile("ion-hash-test/ion_hash_tests.ion")
	require.NoError(t, err, "Something went wrong loading ion_hash_tests.ion")

	expected := ExerciseWriter(t, ion.NewReaderBytes(file), false, writeFromReaderToWriterAfterNext)

	actual := ExerciseWriter(t, ion.NewReaderBytes(file), true, writeFromReaderToWriterAfterNext)

	assert.Greater(t, len(expected), 10, "Expected the ion writer to write more than 10 bytes")

	assert.Greater(t, len(actual), 10, "Expected the ion writer to write more than 10 bytes")

	assert.Equal(t, expected, actual, "sum did not match expectation")
}

func TestIonWriterContractWriteValues(t *testing.T) {
	// Skipping test until FieldNameSymbol logic is available.
	// Test currently fails with empty field name ie. {'':1}
	t.Skip()

	file, err := ioutil.ReadFile("ion-hash-test/ion_hash_tests.ion")
	require.NoError(t, err, "Something went wrong loading ion_hash_tests.ion")

	expected := ExerciseWriter(t, ion.NewReaderBytes(file), false, writeFromReaderToWriter)

	actual := ExerciseWriter(t, ion.NewReaderBytes(file), true, writeFromReaderToWriter)

	assert.Greater(t, len(expected), 1000, "Expected the ion writer to write more than 1000 bytes")

	assert.Greater(t, len(actual), 1000, "Expected the ion writer to write more than 1000 bytes")

	assert.Equal(t, expected, actual, "sum did not match expectation")
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
		require.NoError(t, err, "Expected NewHashWriter() to successfully create a HashWriter")
	}

	function(t, reader, writer)

	assert.NoError(t, writer.Finish(), "Something went wrong executing writer.Finish()")

	return buf.Bytes()
}

type writeFunction func(*testing.T, ion.Reader, ion.Writer)

func writeFromReaderToWriterAfterNext(t *testing.T, reader ion.Reader, writer ion.Writer) {
	require.True(t, reader.Next())

	writeFromReaderToWriter(t, reader, writer)
}
