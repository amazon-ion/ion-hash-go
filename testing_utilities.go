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
	"math"
	"testing"

	"github.com/amzn/ion-go/ion"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func compareReaders(t *testing.T, reader1 ion.Reader, reader2 ion.Reader) {
	for hasNext(t, reader1, reader2) {
		require.Equal(t, reader1.Type().String(), reader2.Type().String(), "Ion Types did not match")

		ionHashReader, ok := reader2.(*hashReader)
		require.True(t, ok, "Expected reader2 to be of type hashReader")

		if ionHashReader.isInStruct() {
			compareFieldNames(t, reader1, reader2)
		}

		compareAnnotations(t, reader1, reader2)

		compareAnnotationSymbols(t, reader1, reader2)

		compareHasAnnotations(t, reader1, reader2)

		require.Equal(t, reader1.IsNull(), reader2.IsNull(), "Expected readers to have matching IsNull() values")

		switch reader1.Type() {
		case ion.NullType:
			assert.True(t, reader1.IsNull(), "Expected reader1.IsNull() to return true")
			assert.True(t, reader2.IsNull(), "Expected reader2.IsNull() to return true")
		case ion.BoolType, ion.IntType, ion.FloatType, ion.DecimalType, ion.TimestampType,
			ion.StringType, ion.SymbolType, ion.BlobType, ion.ClobType:

			compareScalars(t, reader1.Type(), reader1, reader2)
		case ion.StructType, ion.ListType, ion.SexpType:
			assert.NoError(t, reader1.StepIn(), "Something went wrong executing reader1.StepIn()")

			assert.NoError(t, reader2.StepIn(), "Something went wrong executing reader2.StepIn()")

			compareReaders(t, reader1, reader2)

			assert.NoError(t, reader1.StepOut(), "Something went wrong executing reader1.StepOut()")

			assert.NoError(t, reader2.StepOut(), "Something went wrong executing reader2.StepOut()")
		default:
			t.Error(&InvalidIonTypeError{reader1.Type()})
		}
	}

	assert.False(t, hasNext(t, reader1, reader2), "Expected hasNext() to return false")
}

// hasNext() checks that the readers have a Next value
func hasNext(t *testing.T, reader1 ion.Reader, reader2 ion.Reader) bool {
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

func compareFieldNames(t *testing.T, reader1 ion.Reader, reader2 ion.Reader) {
	// TODO: Add SymbolToken logic here once SymbolTokens are available
}

func compareNonNullStrings(t *testing.T, str1, str2 string) {
	assert.NotNil(t, str1, "Expected str1 to be not null")
	assert.NotNil(t, str2, "Expected str2 to be not null")
	assert.Equal(t, str1, str2, "Expected strings to match")
}

func compareAnnotations(t *testing.T, reader1 ion.Reader, reader2 ion.Reader) {
	// TODO: Add SymbolToken logic here once SymbolTokens are available

	assert.Equal(t, reader1.Annotations(), reader2.Annotations(), "Expected symbol sequences to match")
}

func compareAnnotationSymbols(t *testing.T, reader1, reader2 ion.Reader) {
	// TODO: Add SymbolToken logic here once SymbolTokens are available
}

func compareHasAnnotations(t *testing.T, reader1, reader2 ion.Reader) {
	// TODO: Add SymbolToken logic here once SymbolTokens are available
}

func compareScalars(t *testing.T, ionType ion.Type, reader1 ion.Reader, reader2 ion.Reader) {
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
		case ion.Int32, ion.Int64:
			int1, err := reader1.Int64Value()
			assert.NoError(t, err, "Something went wrong executing reader1.Int64Value()")

			int2, err := reader2.Int64Value()
			assert.NoError(t, err, "Something went wrong executing reader2.Int64Value()")

			assert.Equal(t, int1, int2, "Expected int values to match")
		case ion.Uint64:
			uint1, err := reader1.Uint64Value()
			assert.NoError(t, err, "Something went wrong executing reader1.Uint64Value()")

			uint2, err := reader2.Uint64Value()
			assert.NoError(t, err, "Something went wrong executing reader2.Uint64Value()")

			assert.Equal(t, uint1, uint2, "Expected uint values to match")
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

		if math.IsNaN(float1) && math.IsNaN(float2) {
			assert.Equal(t, float1, float2, "Expected NaN float values to match")
		} else if math.IsNaN(float1) || math.IsNaN(float2) {
			assert.NotEqual(t, float1, float2, "Expected IsNaN float value to differ from a non-IsNaN float value")
		} else {
			assert.Equal(t, float1, float2, "Expected float values to match")
		}
	case ion.DecimalType:
		decimal1, err := reader1.DecimalValue()
		assert.NoError(t, err, "Something went wrong executing reader1.DecimalValue()")

		decimal2, err := reader2.DecimalValue()
		assert.NoError(t, err, "Something went wrong executing reader2.DecimalValue()")

		decimalStrictEquals(t, decimal1, decimal2)
	case ion.TimestampType:
		timestamp1, err := reader1.TimeValue()
		assert.NoError(t, err, "Something went wrong executing reader1.TimeValue()")

		timestamp2, err := reader2.TimeValue()
		assert.NoError(t, err, "Something went wrong executing reader2.TimeValue()")

		assert.Equal(t, timestamp1, timestamp2, "Expected timestamp values to match")
	case ion.StringType:
		str1, err := reader1.StringValue()
		assert.NoError(t, err, "Something went wrong executing reader1.StringValue()")

		str2, err := reader2.StringValue()
		assert.NoError(t, err, "Something went wrong executing reader2.StringValue()")

		assert.Equal(t, str1, str2, "Expected string values to match")
	case ion.SymbolType:
		// TODO: Add SymbolToken logic here once SymbolTokens are available
		t.Fail()
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

// Read all the values in the reader and write them in the writer
func writeFromReaderToWriter(t *testing.T, reader ion.Reader, writer ion.Writer) {
	for reader.Next() {
		name := reader.FieldName()
		if name != "" {
			require.NoError(t, writer.FieldName(name), "Something went wrong executing writer.FieldName(name)")
		}

		an := reader.Annotations()
		if len(an) > 0 {
			require.NoError(t, writer.Annotations(an...), "Something went wrong executing writer.Annotations(an...)")
		}

		currentType := reader.Type()
		if reader.IsNull() {
			require.NoError(t, writer.WriteNullType(currentType),
				"Something went wrong executing writer.WriteNullType(currentType)")
			continue
		}

		switch currentType {
		case ion.BoolType:
			val, err := reader.BoolValue()
			assert.NoError(t, err, "Something went wrong when reading Boolean value")

			assert.NoError(t, writer.WriteBool(val), "Something went wrong when writing Boolean value")
		case ion.IntType:
			intSize, err := reader.IntSize()
			require.NoError(t, err, "Something went wrong when retrieving the Int size")

			switch intSize {
			case ion.Int32, ion.Int64:
				val, err := reader.Int64Value()
				assert.NoError(t, err, "Something went wrong when reading Int value")

				assert.NoError(t, writer.WriteInt(val), "Something went wrong when writing Int value")
			case ion.Uint64:
				val, err := reader.Uint64Value()
				assert.NoError(t, err, "Something went wrong when reading UInt value")

				assert.NoError(t, writer.WriteUint(val), "Something went wrong when writing UInt value")
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

			assert.NoError(t, writer.WriteFloat(val), "Something went wrong when writing Float value")
		case ion.DecimalType:
			val, err := reader.DecimalValue()
			assert.NoError(t, err, "Something went wrong when reading Decimal value")

			assert.NoError(t, writer.WriteDecimal(val), "Something went wrong when writing Decimal value")
		case ion.TimestampType:
			val, err := reader.TimeValue()
			assert.NoError(t, err, "Something went wrong when reading Timestamp value")

			assert.NoError(t, writer.WriteTimestamp(val), "Something went wrong when writing Timestamp value")
		case ion.SymbolType:
			val, err := reader.StringValue()
			assert.NoError(t, err, "Something went wrong when reading Symbol value")

			assert.NoError(t, writer.WriteSymbol(val), "Something went wrong when writing Symbol value")
		case ion.StringType:
			val, err := reader.StringValue()
			assert.NoError(t, err, "Something went wrong when reading String value")

			assert.NoError(t, writer.WriteString(val), "Something went wrong when writing String value")
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

			writeFromReaderToWriter(t, reader, writer)

			require.NoError(t, reader.StepOut(), "Something went wrong executing reader.StepOut()")

			require.NoError(t, writer.EndSexp(), "Something went wrong executing writer.EndSexp()")
		case ion.ListType:
			require.NoError(t, reader.StepIn(), "Something went wrong executing reader.StepIn()")

			require.NoError(t, writer.BeginList(), "Something went wrong executing writer.BeginList()")

			writeFromReaderToWriter(t, reader, writer)

			require.NoError(t, reader.StepOut(), "Something went wrong executing reader.StepOut()")

			require.NoError(t, writer.EndList(), "Something went wrong executing writer.EndList()")
		case ion.StructType:
			require.NoError(t, reader.StepIn(), "Something went wrong executing reader.StepIn()")

			require.NoError(t, writer.BeginStruct(), "Something went wrong executing writer.BeginStruct()")

			writeFromReaderToWriter(t, reader, writer)

			require.NoError(t, reader.StepOut(), "Something went wrong executing reader.StepOut()")

			require.NoError(t, writer.EndStruct(), "Something went wrong executing writer.EndStruct()")
		}
	}

	assert.NoError(t, reader.Err(), "Something went wrong executing reader.Next()")
}
