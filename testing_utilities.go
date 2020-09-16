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
	"fmt"
	"math"
	"testing"

	"github.com/amzn/ion-go/ion"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newString(value string) *string {
	return &value
}

// SymbolTokenString returns a SymbolToken's text value or "$<sid>"
func SymbolTokenString(st ion.SymbolToken) string {
	if st.Text != nil {
		return *st.Text
	}

	if st.LocalSID != ion.SymbolIDUnknown {
		return fmt.Sprintf("$%v", st.LocalSID)
	}

	return ""
}

func compareReaders(t *testing.T, reader1, reader2 ion.Reader) {
	for hasNext(t, reader1, reader2) {
		type1 := reader1.Type()
		type2 := reader2.Type()
		require.Equal(t, type1.String(), type2.String(), "Ion Types did not match")

		if type1 == ion.NoType {
			break
		}

		ionHashReader, ok := reader2.(*hashReader)
		require.True(t, ok, "Expected reader2 to be of type hashReader")

		if ionHashReader.IsInStruct() {
			compareFieldNames(t, reader1, reader2)
		}

		compareAnnotations(t, reader1, reader2)

		isNull1 := reader1.IsNull()
		isNull2 := reader2.IsNull()
		require.Equal(t, isNull1, isNull2, "Expected readers to have matching IsNull() values")

		if type1 == ion.NullType {
			assert.True(t, isNull1, "Expected reader1.IsNull() to return true")
			assert.True(t, isNull2, "Expected reader2.IsNull() to return true")
		} else if ion.IsScalar(type1) {
			compareScalars(t, type1, reader1, reader2)
		} else if ion.IsContainer(type1) {
			if !isNull1 {
				assert.NoError(t, reader1.StepIn(), "Something went wrong executing reader1.StepIn()")

				assert.NoError(t, reader2.StepIn(), "Something went wrong executing reader2.StepIn()")

				compareReaders(t, reader1, reader2)

				assert.NoError(t, reader1.StepOut(), "Something went wrong executing reader1.StepOut()")

				assert.NoError(t, reader2.StepOut(), "Something went wrong executing reader2.StepOut()")
			}
		} else {
			t.Error(&InvalidIonTypeError{type1})
		}
	}

	assert.False(t, hasNext(t, reader1, reader2), "Expected hasNext() to return false")
}

// hasNext() checks that the readers have a Next value.
func hasNext(t *testing.T, reader1, reader2 ion.Reader) bool {
	next1 := reader1.Next()
	next2 := reader2.Next()
	assert.Equal(t, next1, next2, "next results don't match")

	if !next1 {
		assert.NoError(t, reader1.Err(), "Something went wrong executing reader1.next()")
	}

	if !next2 {
		assert.NoError(t, reader2.Err(), "Something went wrong executing reader2.next()")
	}

	return next1 && next2
}

func compareFieldNames(t *testing.T, reader1, reader2 ion.Reader) {
	token1, err := reader1.FieldName()
	require.NoError(t, err, "Something went wrong executing reader1.FieldName()")
	token2, err := reader2.FieldName()
	require.NoError(t, err, "Something went wrong executing reader2.FieldName()")

	require.True(t, token1.Equal(token2), "Expected field names to match")
}

func compareAnnotations(t *testing.T, reader1, reader2 ion.Reader) {
	an1, err := reader1.Annotations()
	require.NoError(t, err, "Something went wrong executing reader1.Annotations()")

	an2, err := reader2.Annotations()
	require.NoError(t, err, "Something went wrong executing reader2.Annotations()")

	require.Equal(t, len(an1), len(an2), "Expected readers to have same number of annotations")

	for i := 0; i < len(an1); i++ {
		assert.True(t, an1[i].Equal(&an2[i]))
	}
}

func compareScalars(t *testing.T, ionType ion.Type, reader1, reader2 ion.Reader) {
	isNull1 := reader1.IsNull()
	isNull2 := reader2.IsNull()

	require.Equal(t, isNull1, isNull2, "Expected readers to be both null or both non-null")
	if isNull1 {
		return
	}

	switch ionType {
	case ion.BoolType:
		value1, err := reader1.BoolValue()
		assert.NoError(t, err, "Something went wrong executing reader1.BoolValue()")

		value2, err := reader2.BoolValue()
		assert.NoError(t, err, "Something went wrong executing reader2.BoolValue()")

		assert.Equal(t, value1, value2, "Expected bool values to match")
	case ion.IntType:
		intSize, err := reader1.IntSize()
		assert.NoError(t, err, "Something went wrong executing reader1.IntSize()")

		switch intSize {
		case ion.Int32:
			int1, err := reader1.IntValue()
			assert.NoError(t, err, "Something went wrong executing reader1.IntValue()")

			int2, err := reader2.IntValue()
			assert.NoError(t, err, "Something went wrong executing reader2.IntValue()")

			assert.Equal(t, int1, int2, "Expected int values to match")
		case ion.Int64:
			int1, err := reader1.Int64Value()
			assert.NoError(t, err, "Something went wrong executing reader1.Int64Value()")

			int2, err := reader2.Int64Value()
			assert.NoError(t, err, "Something went wrong executing reader2.Int64Value()")

			assert.Equal(t, int1, int2, "Expected int values to match")
		case ion.BigInt:
			bigInt1, err := reader1.BigIntValue()
			assert.NoError(t, err, "Something went wrong executing reader1.BigIntValue()")

			bigInt2, err := reader2.BigIntValue()
			assert.NoError(t, err, "Something went wrong executing reader2.BigIntValue()")

			assert.Equal(t, bigInt1, bigInt2, "Expected big int values to match")
		default:
			t.Error("Expected intSize to be one of Int32, Int64, Uint64, or BigInt")
		}
	case ion.FloatType:
		float1, err := reader1.FloatValue()
		assert.NoError(t, err, "Something went wrong executing reader1.FloatValue()")

		float2, err := reader2.FloatValue()
		assert.NoError(t, err, "Something went wrong executing reader2.FloatValue()")

		require.True(t, (float1 == nil) == (float2 == nil),
			"Expected float values to be either both null or both not null")

		if float1 != nil {
			if !math.IsNaN(*float1) && !math.IsNaN(*float2) {
				assert.Equal(t, float1, float2, "Expected float values to match")
			} else if !math.IsNaN(*float1) || !math.IsNaN(*float2) {
				assert.NotEqual(t, float1, float2, "Expected IsNaN float value to differ from a non-IsNaN float value")
			}
		}
	case ion.DecimalType:
		decimal1, err := reader1.DecimalValue()
		assert.NoError(t, err, "Something went wrong executing reader1.DecimalValue()")

		decimal2, err := reader2.DecimalValue()
		assert.NoError(t, err, "Something went wrong executing reader2.DecimalValue()")

		decimalStrictEquals(t, decimal1, decimal2)
	case ion.TimestampType:
		timestamp1, err := reader1.TimestampValue()
		assert.NoError(t, err, "Something went wrong executing reader1.TimestampValue()")

		timestamp2, err := reader2.TimestampValue()
		assert.NoError(t, err, "Something went wrong executing reader2.TimestampValue()")

		assert.Equal(t, timestamp1, timestamp2, "Expected timestamp values to match")
	case ion.StringType:
		str1, err := reader1.StringValue()
		assert.NoError(t, err, "Something went wrong executing reader1.StringValue()")

		str2, err := reader2.StringValue()
		assert.NoError(t, err, "Something went wrong executing reader2.StringValue()")

		assert.Equal(t, str1, str2, "Expected string values to match")
	case ion.SymbolType:
		token1, err := reader1.SymbolValue()
		require.NoError(t, err, "Something went wrong executing reader1.SymbolValue()")

		token2, err := reader2.SymbolValue()
		require.NoError(t, err, "Something went wrong executing reader2.SymbolValue()")

		if isNull1 {
			assert.Nil(t, token1.Text, "Expected token1 to have null text")
			assert.Nil(t, token2.Text, "Expected token2 to have null text")
		} else {
			require.Equal(t, token1.Text == nil, token2.Text == nil,
				"Expected the text of both tokens to be null or both not null")

			if token1.Text == nil {
				assert.Equal(t, token1.LocalSID, token2.LocalSID, "Expected token SIDs to match")
			} else {
				assert.Equal(t, token1.Text, token2.Text, "Expected token to have matching text")
			}
		}
	case ion.BlobType, ion.ClobType:
		b1, err := reader1.ByteValue()
		assert.NoError(t, err, "Something went wrong executing reader1.ByteValue()")

		b2, err := reader2.ByteValue()
		assert.NoError(t, err, "Something went wrong executing reader2.ByteValue()")

		assert.True(t, b1 != nil && b2 != nil, "Expected byte arrays to be non-null")

		assert.Equal(t, len(b1), len(b2), "Expected byte arrays to have same length")

		assert.Equal(t, b1, b2, "Expected byte arrays to match")
	default:
		t.Error(InvalidIonTypeError{ionType})
	}
}

// decimalStrictEquals() compares two Ion Decimal values by equality and negative zero.
func decimalStrictEquals(t *testing.T, decimal1, decimal2 *ion.Decimal) {
	assert.Equal(t, decimal1, decimal2, "Expected decimal values to match")

	zeroDecimal := ion.NewDecimalInt(0)

	negativeZero1 := decimal1.Equal(zeroDecimal) && decimal1.Sign() < 0
	negativeZero2 := decimal2.Equal(zeroDecimal) && decimal2.Sign() < 0

	assert.Equal(t, negativeZero1, negativeZero2,
		"Expected decimal values to be both negative zero or both not negative zero")

	assert.True(t, decimal1.Equal(decimal2), "Expected decimal1.Equal(decimal2) to return true")
	assert.True(t, decimal2.Equal(decimal1), "Expected decimal2.Equal(decimal1) to return true")
}

// Read all the values in the reader and write them in the writer.
func writeFromReaderToWriter(t *testing.T, reader ion.Reader, writer ion.Writer, errExpected bool) {
	for reader.Next() {
		name, err := reader.FieldName()
		require.NoError(t, err, "Something went wrong executing reader.Annotations()")

		if name != nil && name.Text != nil {
			require.NoError(t, writer.FieldName(*name.Text), "Something went wrong executing writer.FieldName(*name)")
		}

		annotations, err := reader.Annotations()
		require.NoError(t, err, "Something went wrong executing reader.Annotations()")

		if len(annotations) > 0 {
			require.NoError(t, writer.Annotations(annotations...), "Something went wrong executing writer.Annotations(annotations...)")
		}

		currentType := reader.Type()
		if reader.IsNull() {
			require.NoError(t, writer.WriteNullType(currentType),
				"Something went wrong executing writer.WriteNullType(currentType)")
			continue
		}

		switch currentType {
		case ion.NullType:
			assert.NoError(t, writer.WriteNullType(ion.NullType), "Something went wrong while writing a Null value")

		case ion.BoolType:
			val, err := reader.BoolValue()
			assert.NoError(t, err, "Something went wrong when reading Boolean value")

			if val == nil {
				assert.NoError(t, writer.WriteNullType(ion.BoolType))
			} else {
				assert.NoError(t, writer.WriteBool(*val), "Something went wrong while writing a Boolean value")
			}

		case ion.IntType:
			intSize, err := reader.IntSize()
			require.NoError(t, err, "Something went wrong when retrieving the Int size")

			switch intSize {
			case ion.Int32, ion.Int64:
				val, err := reader.Int64Value()
				assert.NoError(t, err, "Something went wrong when reading Int value")

				assert.NoError(t, writer.WriteInt(*val), "Something went wrong when writing Int value")
			case ion.BigInt:
				val, err := reader.BigIntValue()
				assert.NoError(t, err, "Something went wrong when reading Big Int value")

				assert.NoError(t, writer.WriteBigInt(val), "Something went wrong when writing Big Int value")
			default:
				t.Error("Expected intSize to be one of Int32, Int64, Uint64, or BigInt")
			}

		case ion.FloatType:
			val, err := reader.FloatValue()
			assert.NoError(t, err, "Something went wrong when reading Float value")

			assert.NoError(t, writer.WriteFloat(*val), "Something went wrong when writing Float value")

		case ion.DecimalType:
			val, err := reader.DecimalValue()
			assert.NoError(t, err, "Something went wrong when reading Decimal value")

			assert.NoError(t, writer.WriteDecimal(val), "Something went wrong when writing Decimal value")

		case ion.TimestampType:
			val, err := reader.TimestampValue()
			assert.NoError(t, err, "Something went wrong when reading Timestamp value")

			assert.NoError(t, writer.WriteTimestamp(*val), "Something went wrong when writing Timestamp value")

		case ion.SymbolType:
			val, err := reader.SymbolValue()
			assert.NoError(t, err, "Something went wrong when reading Symbol value")

			assert.NoError(t, writer.WriteSymbol(SymbolTokenString(*val)), "Something went wrong when writing Symbol value")

		case ion.StringType:
			val, err := reader.StringValue()
			assert.NoError(t, err, "Something went wrong when reading String value")

			require.NotNil(t, val)
			assert.NoError(t, writer.WriteString(*val), "Something went wrong when writing String value")

		case ion.ClobType:
			val, err := reader.ByteValue()
			assert.NoError(t, err, "Something went wrong when reading Clob value")

			assert.NoError(t, writer.WriteClob(val), "Something went wrong when writing Clob value")

		case ion.BlobType:
			val, err := reader.ByteValue()
			assert.NoError(t, err, "Something went wrong when reading Blob value")

			assert.NoError(t, writer.WriteBlob(val), "Something went wrong when writing Blob value")

		case ion.SexpType:
			require.NoError(t, reader.StepIn(), "Something went wrong executing reader.StepIn()")
			require.NoError(t, writer.BeginSexp(), "Something went wrong executing writer.BeginSexp()")

			writeFromReaderToWriter(t, reader, writer, errExpected)

			err := reader.StepOut()
			if !errExpected {
				require.NoError(t, err, "Something went wrong executing reader.StepOut()")
			}

			require.NoError(t, writer.EndSexp(), "Something went wrong executing writer.EndSexp()")

		case ion.ListType:
			require.NoError(t, reader.StepIn(), "Something went wrong executing reader.StepIn()")
			require.NoError(t, writer.BeginList(), "Something went wrong executing writer.BeginList()")

			writeFromReaderToWriter(t, reader, writer, errExpected)

			err := reader.StepOut()
			if !errExpected {
				require.NoError(t, err, "Something went wrong executing reader.StepOut()")
			}

			require.NoError(t, writer.EndList(), "Something went wrong executing writer.EndList()")

		case ion.StructType:
			require.NoError(t, reader.StepIn(), "Something went wrong executing reader.StepIn()")
			require.NoError(t, writer.BeginStruct(), "Something went wrong executing writer.BeginStruct()")

			writeFromReaderToWriter(t, reader, writer, errExpected)

			err := reader.StepOut()
			if !errExpected {
				require.NoError(t, err, "Something went wrong executing reader.StepOut()")
			}

			require.NoError(t, writer.EndStruct(), "Something went wrong executing writer.EndStruct()")
		}
	}

	if !errExpected {
		assert.NoError(t, reader.Err(), "Something went wrong writing from reader to writer")
	}
}

func writeToWriters(t *testing.T, reader ion.Reader, writers ...ion.Writer) {
	ionType := reader.Type()

	annotations, err := reader.Annotations()
	require.NoError(t, err, "Something went wrong executing reader.Annotations()")

	if len(annotations) > 0 {
		for _, writer := range writers {
			require.NoError(t, writer.Annotations(annotations...),
				"Something went wrong executing writer.Annotations(annotations...)")
		}
	}

	fieldName, err := reader.FieldName()
	if err == nil && fieldName.Text != nil && *fieldName.Text != "ion" && *fieldName.Text != "10n" {
		for _, writer := range writers {
			require.NoError(t, writer.FieldName(*fieldName.Text),
				"Something went wrong executing writer.FieldName(*fieldName)")
		}
	}

	if reader.IsNull() {
		for _, writer := range writers {
			require.NoError(t, writer.WriteNullType(reader.Type()),
				"Something went wrong executing writer.WriteNullType(reader.Type())")
		}

		return
	}

	switch ionType {
	case ion.NullType:
		for _, writer := range writers {
			require.NoError(t, writer.WriteNull(), "Something went wrong executing writer.WriteNull()")
		}
	case ion.BoolType:
		boolValue, err := reader.BoolValue()
		require.NoError(t, err)
		for _, writer := range writers {
			require.NoError(t, writer.WriteBool(*boolValue),
				"Something went wrong executing writer.WriteBool(*boolValue)")
		}
	case ion.BlobType:
		byteValue, err := reader.ByteValue()
		require.NoError(t, err)
		for _, writer := range writers {
			require.NoError(t, writer.WriteBlob(byteValue),
				"Something went wrong executing writer.WriteBlob(byteValue)")
		}
	case ion.ClobType:
		byteValue, err := reader.ByteValue()
		require.NoError(t, err)
		for _, writer := range writers {
			require.NoError(t, writer.WriteClob(byteValue),
				"Something went wrong executing writer.WriteClob(byteValue)")
		}
	case ion.DecimalType:
		decimalValue, err := reader.DecimalValue()
		require.NoError(t, err)
		for _, writer := range writers {
			require.NoError(t, writer.WriteDecimal(decimalValue),
				"Something went wrong executing writer.WriteDecimal(decimalValue)")
		}
	case ion.FloatType:
		floatValue, err := reader.FloatValue()
		require.NoError(t, err)
		for _, writer := range writers {
			require.NoError(t, writer.WriteFloat(*floatValue),
				"Something went wrong executing writer.WriteFloat(*floatValue)")
		}
	case ion.IntType:
		intSize, err := reader.IntSize()
		require.NoError(t, err)

		switch intSize {
		case ion.Int32, ion.Int64:
			intValue, err := reader.Int64Value()
			require.NoError(t, err)
			for _, writer := range writers {
				require.NoError(t, writer.WriteInt(*intValue),
					"Something went wrong executing writer.WriteInt(*intValue)")
			}
		case ion.BigInt:
			bigIntValue, err := reader.BigIntValue()
			require.NoError(t, err)
			for _, writer := range writers {
				require.NoError(t, writer.WriteBigInt(bigIntValue),
					"Something went wrong executing writer.WriteBigInt(bigIntValue)")
			}
		default:
			t.Error("Expected intSize to be one of Int32, Int64, Uint64, or BigInt")
		}

	case ion.StringType:
		stringValue, err := reader.StringValue()
		require.NoError(t, err)
		require.NotNil(t, stringValue)
		for _, writer := range writers {
			require.NoError(t, writer.WriteString(*stringValue),
				"Something went wrong executing writer.WriteString(stringValue)")
		}
	case ion.SymbolType:
		symbolValue, err := reader.SymbolValue()
		require.NoError(t, err)
		for _, writer := range writers {
			require.NoError(t, writer.WriteSymbol(SymbolTokenString(*symbolValue)),
				"Something went wrong executing writer.WriteSymbol(SymbolTokenString(*symbolValue))")
		}
	case ion.TimestampType:
		timestampValue, err := reader.TimestampValue()
		require.NoError(t, err)
		for _, writer := range writers {
			require.NoError(t, writer.WriteTimestamp(*timestampValue),
				"Something went wrong executing writer.WriterTimestamp(*timestampValue)")
		}

	case ion.SexpType:
		require.NoError(t, reader.StepIn())
		for _, writer := range writers {
			require.NoError(t, writer.BeginSexp(), "Something went wrong executing writer.BeginSexp()")
		}
		for reader.Next() {
			writeToWriters(t, reader, writers...)
		}
		require.NoError(t, reader.Err(), "Something went wrong executing reader.Next()")

		require.NoError(t, reader.StepOut())
		for _, writer := range writers {
			require.NoError(t, writer.EndSexp(), "Something went wrong executing writer.EndSexp()")
		}
	case ion.ListType:
		require.NoError(t, reader.StepIn())
		for _, writer := range writers {
			require.NoError(t, writer.BeginList(), "Something went wrong executing writer.BeginList()")
		}

		for reader.Next() {
			writeToWriters(t, reader, writers...)
		}
		require.NoError(t, reader.Err(), "Something went wrong executing reader.Next()")

		require.NoError(t, reader.StepOut())
		for _, writer := range writers {
			require.NoError(t, writer.EndList(), "Something went wrong executing writer.EndList()")
		}
	case ion.StructType:
		require.NoError(t, reader.StepIn())
		for _, writer := range writers {
			require.NoError(t, writer.BeginStruct(), "Something went wrong executing writer.BeginStruct()")
		}

		for reader.Next() {
			writeToWriters(t, reader, writers...)
		}
		require.NoError(t, reader.Err(), "Something went wrong executing reader.Next()")

		require.NoError(t, reader.StepOut())
		for _, writer := range writers {
			require.NoError(t, writer.EndStruct(), "Something went wrong executing writer.EndStruct()")
		}

	default:
		t.Fatal(InvalidIonTypeError{ionType})
	}
}

func readSexpAndAppendToList(t *testing.T, reader ion.Reader) []byte {
	require.NoError(t, reader.StepIn())
	updateBytes := []byte{}
	for reader.Next() {
		intValue, err := reader.Int64Value()
		require.NoError(t, err, "Something went wrong executing reader.Int64Value()")
		updateBytes = append(updateBytes, byte(*intValue))
	}
	require.NoError(t, reader.Err(), "Something went wrong executing reader.Next()")
	require.NoError(t, reader.StepOut(), "Something went wrong executing reader.StepOut()")
	return updateBytes
}
