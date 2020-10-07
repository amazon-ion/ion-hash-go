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

var nullTests = []struct {
	ionType     ion.Type
	expectedSum []byte
	String      string
}{
	{ion.NoType, []byte{0x0b, 0x0f, 0x0e}, "null"},
	{ion.NullType, []byte{0x0b, 0x0f, 0x0e}, "null.null"},
	{ion.BoolType, []byte{0x0b, 0x1f, 0x0e}, "null.bool"},
	{ion.IntType, []byte{0x0b, 0x2f, 0x0e}, "null.int"},
	{ion.FloatType, []byte{0x0b, 0x4f, 0x0e}, "null.float"},
	{ion.DecimalType, []byte{0x0b, 0x5f, 0x0e}, "null.decimal"},
	{ion.TimestampType, []byte{0x0b, 0x6f, 0x0e}, "null.timestamp"},
	{ion.SymbolType, []byte{0x0b, 0x7f, 0x0e}, "null.symbol"},
	{ion.StringType, []byte{0x0b, 0x8f, 0x0e}, "null.string"},
	{ion.ClobType, []byte{0x0b, 0x9f, 0x0e}, "null.clob"},
	{ion.BlobType, []byte{0x0b, 0xaf, 0x0e}, "null.blob"},
	{ion.ListType, []byte{0x0b, 0xbf, 0x0e}, "null.list"},
	{ion.SexpType, []byte{0x0b, 0xcf, 0x0e}, "null.sexp"},
	{ion.StructType, []byte{0x0b, 0xdf, 0x0e}, "null.struct"},
}

var scalarTests = []struct {
	ionType     ion.Type
	value       interface{}
	expectedSum []byte
	String      string
}{
	{
		ion.NoType,
		nil,
		[]byte{},
		"",
	},
	{
		ion.BoolType,
		true,
		[]byte{0x0b, 0x11, 0x0e},
		"true",
	},
	{
		ion.IntType,
		uint64(5),
		[]byte{0x0b, 0x20, 0x05, 0x0e},
		"5",
	},
	{
		ion.IntType,
		int64(-5),
		[]byte{0x0b, 0x30, 0x05, 0x0e},
		"-5",
	},
	{
		ion.IntType,
		int64(123456789),
		[]byte{0xb, 0x20, 0x7, 0x5b, 0xcd, 0x15, 0xe},
		"123456789",
	},
	{
		ion.FloatType,
		3.14,
		[]byte{0x0b, 0x40, 0x40, 0x09, 0x1e, 0xb8, 0x51, 0xeb, 0x85, 0x1f, 0x0e},
		"3.14e+0",
	},
	{
		ion.DecimalType,
		"1234.56789",
		[]byte{0x0b, 0x50, 0xc5, 0x07, 0x5b, 0xcd, 0x15, 0x0e},
		"1234.56789",
	},
	{
		ion.TimestampType,
		ion.NewTimestamp(time.Date(1941, time.December, 7, 18, 0, 0, 0, time.UTC),
			ion.TimestampPrecisionSecond, ion.TimezoneUTC),
		[]byte{0x0b, 0x60, 0x80, 0x0f, 0x95, 0x8c, 0x87, 0x92, 0x80, 0x80, 0x0e},
		"1941-12-07T18:00:00Z",
	},
	{
		ion.SymbolType,
		ion.NewSimpleSymbolToken("symbol"),
		[]byte{0x0b, 0x70, 0x73, 0x79, 0x6d, 0x62, 0x6f, 0x6c, 0x0e},
		"symbol",
	},
	{
		ion.StringType,
		"string",
		[]byte{0x0b, 0x80, 0x73, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x0e},
		"\"string\"",
	},
	{
		ion.ClobType,
		[]byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f},
		[]byte{0x0b, 0x90, 0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09,
			0x0a, 0x0c, 0x0b, 0x0c, 0x0c, 0x0d, 0x0c, 0x0e, 0x0f, 0x0e},
		"{{\"\\0\\x01\\x02\\x03\\x04\\x05\\x06\\a\\b\\t\\n\\v\\f\\r\\x0E\\x0F\"}}",
	},
	{
		ion.BlobType,
		[]byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f},
		[]byte{0x0b, 0xa0, 0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09,
			0x0a, 0x0c, 0x0b, 0x0c, 0x0c, 0x0d, 0x0c, 0x0e, 0x0f, 0x0e},
		"{{AAECAwQFBgcICQoLDA0ODw==}}",
	},
}

func TestWriteNull(t *testing.T) {
	str := strings.Builder{}
	tihp := newTestIonHasherProvider("identity")
	hw, err := NewHashWriter(ion.NewTextWriter(&str), tihp.getInstance())
	require.NoError(t, err, "Expected NewHashWriter() to successfully create a HashWriter")

	ionHashWriter, ok := hw.(*hashWriter)
	require.True(t, ok, "Expected hw to be of type hashWriter")

	expectedStr := ""
	for _, test := range nullTests {
		if test.ionType == ion.NoType {
			assert.NoError(t, ionHashWriter.WriteNull(), "Something went wrong executing ionHashWriter.WriteNull()")
		} else {
			assert.NoErrorf(t, ionHashWriter.WriteNullType(test.ionType),
				"Something went wrong executing ionHashWriter.WriteNullType() for type %s", test.ionType.String())
		}

		sum, err := ionHashWriter.Sum(nil)
		require.NoError(t, err, "Something went wrong executing ionHashWriter.Sum(nil)")

		assert.Equal(t, test.expectedSum, sum, "sum did not match expectation")

		expectedStr += test.String + "\n"
	}

	assert.NoError(t, ionHashWriter.Finish(), "Something went wrong executing ionHashWriter.Finish()")

	assert.Equal(t, expectedStr, str.String(), "str.String() did not match expectation")
}

func TestWriteScalars(t *testing.T) {
	str := strings.Builder{}
	tihp := newTestIonHasherProvider("identity")
	hw, err := NewHashWriter(ion.NewTextWriter(&str), tihp.getInstance())
	require.NoError(t, err, "Expected NewHashWriter() to successfully create a HashWriter")

	ionHashWriter, ok := hw.(*hashWriter)
	require.True(t, ok, "Expected hw to be of type hashWriter")

	expectedStr := ""
	for _, test := range scalarTests {
		switch test.ionType {
		case ion.BoolType:
			assert.NoErrorf(t, ionHashWriter.WriteBool(test.value.(bool)),
				"Something went wrong executing ionHashWriter.WriteBool(%s)", test.String)
		case ion.IntType:
			switch test.value {
			case uint64(5):
				assert.NoErrorf(t, ionHashWriter.WriteUint(test.value.(uint64)),
					"Something went wrong executing ionHashWriter.WriteUint(%s)", test.String)
			case int64(-5):
				assert.NoErrorf(t, ionHashWriter.WriteInt(test.value.(int64)),
					"Something went wrong executing ionHashWriter.WriteInt(%s)", test.String)
			case int64(123456789):
				bigInt := big.NewInt(test.value.(int64))
				assert.NoError(t, ionHashWriter.WriteBigInt(bigInt),
					"Something went wrong executing ionHashWriter.WriteBigInt(bigInt)")
			}
		case ion.FloatType:
			assert.NoErrorf(t, ionHashWriter.WriteFloat(test.value.(float64)),
				"Something went wrong executing ionHashWriter.WriteFloat(%s)", test.String)
		case ion.DecimalType:
			dec, err := ion.ParseDecimal(test.value.(string))
			assert.NoErrorf(t, err, "Something went wrong executing ion.ParseDecimal(\"%s\")", test.String)
			assert.NoError(t, ionHashWriter.WriteDecimal(dec),
				"Something went wrong executing ionHashWriter.WriteDecimal(dec)")
		case ion.TimestampType:
			assert.NoError(t, ionHashWriter.WriteTimestamp(test.value.(ion.Timestamp)),
				"Something went wrong executing ionHashWriter.WriteTimestamp(...)")
		case ion.SymbolType:
			assert.NoErrorf(t, ionHashWriter.WriteSymbol(test.value.(ion.SymbolToken)),
				"Something went wrong executing ionHashWriter.WriteSymbol(\"%s\")", test.String)
		case ion.StringType:
			assert.NoErrorf(t, ionHashWriter.WriteString(test.value.(string)),
				"Something went wrong executing ionHashWriter.WriteString(\"%s\")", test.String)
		case ion.ClobType:
			err = ionHashWriter.WriteClob(test.value.([]byte))
			assert.NoError(t, err, "Something went wrong executing ionHashWriter.WriteClob(...)")
		case ion.BlobType:
			err = ionHashWriter.WriteBlob(test.value.([]byte))
			assert.NoError(t, err, "Something went wrong executing ionHashWriter.WriteBlob(...)")
		}

		sum, err := ionHashWriter.Sum(nil)
		require.NoError(t, err, "Something went wrong executing ionHashWriter.Sum(nil)")

		assert.Equal(t, test.expectedSum, sum, "sum did not match expectation")

		if test.String != "" {
			expectedStr += test.String + "\n"
		}
	}

	assert.NoError(t, ionHashWriter.Finish(), "Something went wrong executing ionHashWriter.Finish()")

	assert.Equal(t, expectedStr, str.String(), "str.String() did not match expectation")
}

func TestWriteContainers(t *testing.T) {
	str := strings.Builder{}
	tihp := newTestIonHasherProvider("identity")
	hw, err := NewHashWriter(ion.NewTextWriter(&str), tihp.getInstance())
	require.NoError(t, err, "Expected NewHashWriter() to successfully create a HashWriter")

	ionHashWriter, ok := hw.(*hashWriter)
	require.True(t, ok, "Expected hw to be of type hashWriter")

	sum, err := ionHashWriter.Sum(nil)
	require.NoError(t, err, "Something went wrong executing ionHashWriter.Sum(nil)")

	assert.Equal(t, []byte{}, sum, "sum did not match expectation")

	assert.NoError(t, ionHashWriter.BeginList(), "Something went wrong executing ionHashWriter.BeginList()")

	sum, err = ionHashWriter.Sum(nil)
	assert.Error(t, err, "Expected ionHashWriter.Sum(nil) to return an error")
	assert.IsType(t, &InvalidOperationError{}, err,
		"Expected ionHashWriter.Sum(nil) to return an InvalidOperationError")

	assert.NoError(t, ionHashWriter.WriteBool(true), "Something went wrong executing ionHashWriter.WriteBool(true)")

	sum, err = ionHashWriter.Sum(nil)
	assert.Error(t, err, "Expected ionHashWriter.Sum(nil) to return an error")
	assert.IsType(t, &InvalidOperationError{}, err,
		"Expected ionHashWriter.Sum(nil) to return an InvalidOperationError")

	assert.NoError(t, ionHashWriter.EndList(), "Something went wrong executing ionHashWriter.EndList()")

	sum, err = ionHashWriter.Sum(nil)
	require.NoError(t, err, "Something went wrong executing ionHashWriter.Sum(nil)")

	assert.Equal(t, []byte{0x0b, 0xb0, 0x0b, 0x11, 0x0e, 0x0e}, sum, "sum did not match expectation")

	assert.False(t, ionHashWriter.IsInStruct())

	assert.NoError(t, ionHashWriter.BeginStruct(), "Something went wrong executing ionHashWriter.BeginStruct()")

	assert.True(t, ionHashWriter.IsInStruct())

	assert.NoError(t, ionHashWriter.FieldName(ion.NewSimpleSymbolToken("hello")),
		"Something went wrong executing ionHashWriter.FieldName(...)")

	assert.NoError(t, ionHashWriter.Annotation(ion.NewSimpleSymbolToken("ion")),
		"Something went wrong executing ionHashWriter.Annotation(...)")

	assert.NoError(t, ionHashWriter.Annotation(ion.NewSimpleSymbolToken("hash")),
		"Something went wrong executing ionHashWriter.Annotation(...)")

	assert.NoError(t, ionHashWriter.WriteSymbolFromString("world"),
		"Something went wrong executing ionHashWriter.WriteSymbolFromString(\"world\")")

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

	assert.Equal(t, "[true]\n{hello:ion::hash::world}\n", str.String(), "str.String() did not match expectation")
}

func TestExtraEndContainer(t *testing.T) {
	str := strings.Builder{}
	tihp := newTestIonHasherProvider("identity")
	hw, err := NewHashWriter(ion.NewTextWriter(&str), tihp.getInstance())
	require.NoError(t, err, "Expected NewHashWriter() to successfully create a HashWriter")

	ionHashWriter, ok := hw.(*hashWriter)
	require.True(t, ok, "Expected hw to be of type hashWriter")

	err = ionHashWriter.EndList()
	assert.Error(t, err, "Expected ionHashWriter.EndList() to return an error")
	assert.IsType(t, &InvalidOperationError{}, err,
		"Expected ionHashWriter.EndList() to return an InvalidOperationError")

	err = ionHashWriter.EndSexp()
	assert.Error(t, err, "Expected ionHashWriter.EndSexp() to return an error")
	assert.IsType(t, &InvalidOperationError{}, err,
		"Expected ionHashWriter.EndSexp() to return an InvalidOperationError")

	err = ionHashWriter.EndStruct()
	assert.Error(t, err, "Expected ionHashWriter.EndStruct() to return an error")
	assert.IsType(t, &InvalidOperationError{}, err,
		"Expected ionHashWriter.EndStruct() to return an InvalidOperationError")
}

func TestIonWriterContractWriteValue(t *testing.T) {
	file, err := ioutil.ReadFile("ion-hash-test/ion_hash_tests.ion")
	require.NoError(t, err, "Something went wrong loading ion_hash_tests.ion")

	expected := ExerciseWriter(t, ion.NewReaderBytes(file), false, writeFromReaderToWriterAfterNext)

	actual := ExerciseWriter(t, ion.NewReaderBytes(file), true, writeFromReaderToWriterAfterNext)

	assert.Greater(t, len(expected), 10, "Expected the ion writer to write more than 10 bytes")

	assert.Greater(t, len(actual), 10, "Expected the ion writer to write more than 10 bytes")

	assert.Equal(t, expected, actual, "sum did not match expectation")
}

func TestIonWriterContractWriteValues(t *testing.T) {
	file, err := ioutil.ReadFile("ion-hash-test/ion_hash_tests.ion")
	require.NoError(t, err, "Something went wrong loading ion_hash_tests.ion")

	expected := ExerciseWriter(t, ion.NewReaderBytes(file), false, writeFromReaderToWriter)

	actual := ExerciseWriter(t, ion.NewReaderBytes(file), true, writeFromReaderToWriter)

	assert.Greater(t, len(expected), 1000, "Expected the ion writer to write more than 1000 bytes")

	assert.Greater(t, len(actual), 1000, "Expected the ion writer to write more than 1000 bytes")

	assert.Equal(t, expected, actual, "sum did not match expectation")
}

func ExerciseWriter(t *testing.T, reader ion.Reader, useHashWriter bool, fn func(*testing.T, ion.Reader, ion.Writer, bool)) []byte {
	var err error

	buf := bytes.Buffer{}
	writer := ion.NewBinaryWriter(&buf)

	if useHashWriter {
		tihp := newTestIonHasherProvider("identity")
		writer, err = NewHashWriter(writer, tihp.getInstance())
		require.NoError(t, err, "Expected NewHashWriter() to successfully create a HashWriter")
	}

	fn(t, reader, writer, false)

	assert.NoError(t, writer.Finish(), "Something went wrong executing writer.Finish()")

	return buf.Bytes()
}

func writeFromReaderToWriterAfterNext(t *testing.T, reader ion.Reader, writer ion.Writer, errExpected bool) {
	require.True(t, reader.Next())

	writeFromReaderToWriter(t, reader, writer, errExpected)
}
