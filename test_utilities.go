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
	"testing"

	"github.com/amzn/ion-go/ion"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func writeToWriter(t *testing.T, reader ion.Reader, textWriter ion.Writer, binaryWriter ion.Writer) {
	ionType := reader.Type()

	if reader.Annotations() != nil {
		err := textWriter.Annotations(reader.Annotations()...)
		require.NoError(t, err)
		err = binaryWriter.Annotations(reader.Annotations()...)
		require.NoError(t, err)
	}

	switch ionType {
	case ion.NullType:
		if reader.IsNull() {
			err := textWriter.WriteNullType(ion.NullType)
			require.NoError(t, err)
			err = binaryWriter.WriteNullType(ion.NullType)
			require.NoError(t, err)
		} else {
			err := textWriter.WriteNull()
			require.NoError(t, err)
			err = binaryWriter.WriteNull()
			require.NoError(t, err)
		}

	case ion.BoolType:
		if reader.IsNull() {
			err := textWriter.WriteNullType(ion.BoolType)
			require.NoError(t, err)
			err = binaryWriter.WriteNullType(ion.BoolType)
			require.NoError(t, err)
		} else {
			boolValue, err := reader.BoolValue()
			require.NoError(t, err)
			err = textWriter.WriteBool(boolValue)
			require.NoError(t, err)
			err = binaryWriter.WriteBool(boolValue)
			require.NoError(t, err)
		}

	case ion.BlobType:
		if reader.IsNull() {
			err := textWriter.WriteNullType(ion.BlobType)
			require.NoError(t, err)
			err = binaryWriter.WriteNullType(ion.BlobType)
			require.NoError(t, err)
		} else {
			byteValue, err := reader.ByteValue()
			require.NoError(t, err)
			err = textWriter.WriteBlob(byteValue)
			require.NoError(t, err)
			err = binaryWriter.WriteBlob(byteValue)
			require.NoError(t, err)
		}

	case ion.ClobType:
		if reader.IsNull() {
			err := textWriter.WriteNullType(ion.ClobType)
			require.NoError(t, err)
			err = binaryWriter.WriteNullType(ion.ClobType)
			require.NoError(t, err)
		} else {
			byteValue, err := reader.ByteValue()
			require.NoError(t, err)
			err = textWriter.WriteClob(byteValue)
			require.NoError(t, err)
			err = binaryWriter.WriteClob(byteValue)
			require.NoError(t, err)
		}

	case ion.DecimalType:
		if reader.IsNull() {
			err := textWriter.WriteNullType(ion.DecimalType)
			require.NoError(t, err)
			err = binaryWriter.WriteNullType(ion.DecimalType)
			require.NoError(t, err)
		} else {
			decimalValue, err := reader.DecimalValue()
			require.NoError(t, err)
			err = textWriter.WriteDecimal(decimalValue)
			require.NoError(t, err)
			err = binaryWriter.WriteDecimal(decimalValue)
			require.NoError(t, err)
		}

	case ion.FloatType:
		if reader.IsNull() {
			err := textWriter.WriteNullType(ion.FloatType)
			require.NoError(t, err)
			err = binaryWriter.WriteNullType(ion.FloatType)
			require.NoError(t, err)
		} else {
			floatValue, err := reader.FloatValue()
			require.NoError(t, err)
			err = textWriter.WriteFloat(floatValue)
			require.NoError(t, err)
			err = binaryWriter.WriteFloat(floatValue)
			require.NoError(t, err)
		}

	case ion.IntType:
		intSize, err := reader.IntSize()
		require.NoError(t, err)

		if reader.IsNull() {
			err = textWriter.WriteNullType(ion.IntType)
			require.NoError(t, err)
			err = binaryWriter.WriteNullType(ion.IntType)
			require.NoError(t, err)
		} else {
			switch intSize {
			case ion.Int32:

				intValue, err := reader.IntValue()
				require.NoError(t, err)
				err = textWriter.WriteInt(int64(intValue))
				require.NoError(t, err)
				err = binaryWriter.WriteInt(int64(intValue))
				require.NoError(t, err)

			case ion.Int64:
				intValue, err := reader.Int64Value()
				require.NoError(t, err)
				err = textWriter.WriteInt(intValue)
				require.NoError(t, err)
				err = binaryWriter.WriteInt(intValue)
				require.NoError(t, err)

			case ion.Uint64:
				intValue, err := reader.Uint64Value()
				require.NoError(t, err)
				err = textWriter.WriteUint(intValue)
				require.NoError(t, err)
				err = binaryWriter.WriteUint(intValue)
				require.NoError(t, err)

			case ion.BigInt:
				intValue, err := reader.BigIntValue()
				require.NoError(t, err)
				err = textWriter.WriteBigInt(intValue)
				require.NoError(t, err)
				err = binaryWriter.WriteBigInt(intValue)
				require.NoError(t, err)

			default:
				t.Error("Expected intSize to be one of Int32, Int64, Uint64, or BigInt")
			}
		}

	case ion.StringType:
		if reader.IsNull() {
			err := textWriter.WriteNullType(ion.StringType)
			require.NoError(t, err)
			err = binaryWriter.WriteNullType(ion.StringType)
			require.NoError(t, err)
		} else {
			stringValue, err := reader.StringValue()
			require.NoError(t, err)
			err = textWriter.WriteString(stringValue)
			require.NoError(t, err)
			err = binaryWriter.WriteString(stringValue)
			require.NoError(t, err)
		}

	case ion.SymbolType:
		if reader.IsNull() {
			err := textWriter.WriteNullType(ion.SymbolType)
			require.NoError(t, err)
			err = binaryWriter.WriteNullType(ion.SymbolType)
			require.NoError(t, err)
		} else {
			stringValue, err := reader.StringValue()
			require.NoError(t, err)
			err = textWriter.WriteSymbol(stringValue)
			require.NoError(t, err)
			err = binaryWriter.WriteSymbol(stringValue)
			require.NoError(t, err)
		}

	case ion.TimestampType:
		if reader.IsNull() {
			err := textWriter.WriteNullType(ion.TimestampType)
			require.NoError(t, err)
			err = binaryWriter.WriteNullType(ion.TimestampType)
			require.NoError(t, err)
		} else {
			timeValue, err := reader.TimeValue()
			require.NoError(t, err)
			err = textWriter.WriteTimestamp(timeValue)
			require.NoError(t, err)
			err = binaryWriter.WriteTimestamp(timeValue)
			require.NoError(t, err)
		}

	case ion.SexpType:
		if reader.IsNull() {
			err := textWriter.WriteNullType(ion.SexpType)
			require.NoError(t, err)
			err = binaryWriter.WriteNullType(ion.SexpType)
			require.NoError(t, err)
		} else {
			err := reader.StepIn()
			require.NoError(t, err)
			err = textWriter.BeginSexp()
			require.NoError(t, err)
			err = binaryWriter.BeginSexp()
			require.NoError(t, err)
			if reader.FieldName() != "" {
				err = textWriter.FieldName(reader.FieldName())
				require.NoError(t, err)
				err = binaryWriter.FieldName(reader.FieldName())
				require.NoError(t, err)
			}
			for reader.Next() {
				writeToWriter(t, reader, textWriter, binaryWriter)
			}

			err = reader.StepOut()
			require.NoError(t, err)
			err = textWriter.EndSexp()
			require.NoError(t, err)
			err = binaryWriter.EndSexp()
			require.NoError(t, err)
		}

	case ion.ListType:
		if reader.IsNull() {
			err := textWriter.WriteNullType(ion.ListType)
			require.NoError(t, err)
			err = binaryWriter.WriteNullType(ion.ListType)
			require.NoError(t, err)
		} else {
			err := reader.StepIn()
			require.NoError(t, err)
			err = textWriter.BeginList()
			require.NoError(t, err)
			err = binaryWriter.BeginList()
			require.NoError(t, err)
			if reader.FieldName() != "" {
				err = textWriter.FieldName(reader.FieldName())
				require.NoError(t, err)
				err = binaryWriter.FieldName(reader.FieldName())
				require.NoError(t, err)
			}
			for reader.Next() {
				writeToWriter(t, reader, textWriter, binaryWriter)
			}

			err = reader.StepOut()
			require.NoError(t, err)
			err = textWriter.EndList()
			require.NoError(t, err)
			err = binaryWriter.EndList()
			require.NoError(t, err)
		}

	case ion.StructType:
		if reader.IsNull() {
			textWriter.WriteNullType(ion.StructType)
			err := textWriter.WriteNullType(ion.StructType)
			require.NoError(t, err)
			err = binaryWriter.WriteNullType(ion.StructType)
			require.NoError(t, err)
		} else {
			err := reader.StepIn()
			require.NoError(t, err)
			err = textWriter.BeginStruct()
			require.NoError(t, err)
			err = binaryWriter.BeginStruct()
			require.NoError(t, err)
			if reader.FieldName() != "" {
				err = textWriter.FieldName(reader.FieldName())
				require.NoError(t, err)
				err = binaryWriter.FieldName(reader.FieldName())
				require.NoError(t, err)
			}

			for reader.Next() {
				writeToWriter(t, reader, textWriter, binaryWriter)
			}

			err = reader.StepOut()
			require.NoError(t, err)
			err = textWriter.EndStruct()
			require.NoError(t, err)
			err = binaryWriter.EndStruct()
			require.NoError(t, err)
		}

	default:
		t.Fatal(InvalidIonTypeError{ionType})
	}
}

func readSexpAndAppendToList(t *testing.T, reader ion.Reader) []byte {
	err := reader.StepIn()
	require.NoError(t, err)
	updateBytes := []byte{}
	for reader.Next() {
		intValue, err := reader.Int64Value()
		assert.NoError(t, err, "Something went wrong executing reader.Int64Value()")
		updateBytes = append(updateBytes, byte(intValue))
	}
	err = reader.StepOut()
	require.NoError(t, err)
	return updateBytes
}
