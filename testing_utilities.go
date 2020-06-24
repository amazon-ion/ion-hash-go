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
	"reflect"
	"testing"

	"github.com/amzn/ion-go/ion"
)

func compareReaders(t *testing.T, reader1 ion.Reader, reader2 ion.Reader) {
	for true {
		if !hasNext(t, reader1, reader2) {
			break
		}

		ionType1 := reader1.Type()
		ionType2 := reader2.Type()

		if ionType1 != ionType2 {
			t.Errorf("Ion types do not match;\n"+
				"Type #1: %s\n"+
				"Type #2: %s",
				ionType1.String(), ionType2.String())
		}

		ionHashReader, ok := reader2.(*hashReader)
		if ok {
			if ionHashReader.isInStruct() {
				compareFieldNames(t, reader1, reader2)
			}
		} else {
			t.Errorf("Expected reader2 to be of type hashReader")
		}

		compareAnnotations(t, reader1, reader2)

		compareAnnotationSymbols(t, reader1, reader2)

		compareHasAnnotations(t, reader1, reader2)

		isNull1 := reader1.IsNull()
		isNull2 := reader2.IsNull()
		if isNull1 != isNull2 {
			t.Errorf("Expected readers to have matching IsNull() values;\n"+
				"isNull1: %v,\n"+
				"isNull2: %v",
				isNull1, isNull2)
		}

		switch ionType1 {
		case ion.NullType:
			if !isNull1 {
				t.Errorf("Expected ionType1 to be null")
			}
			if !isNull2 {
				t.Errorf("Expected ionType2 to be null")
			}
		case ion.BoolType, ion.IntType, ion.FloatType, ion.DecimalType, ion.TimestampType,
			ion.StringType, ion.SymbolType, ion.BlobType, ion.ClobType:

			compareScalars(t, reader1, reader2)
		case ion.StructType, ion.ListType, ion.SexpType:
			err := reader1.StepIn()
			if err != nil {
				t.Errorf("Something went wrong executing reader1.StepIn(); %s", err.Error())
			}

			err = reader2.StepIn()
			if err != nil {
				t.Errorf("Something went wrong executing reader2.StepIn(); %s", err.Error())
			}

			compareReaders(t, reader1, reader2)

			err = reader1.StepOut()
			if err != nil {
				t.Errorf("Something went wrong executing reader1.StepOut(); %s", err.Error())
			}

			err = reader2.StepOut()
			if err != nil {
				t.Errorf("Something went wrong executing reader1.StepOut(); %s", err.Error())
			}
		default:
			t.Error(&InvalidIonTypeError{ionType1})
		}
	}

	if hasNext(t, reader1, reader2) {
		t.Errorf("Expected hasNext() to return false")
	}
}

// hasNext() checks that the readers have a Next value
func hasNext(t *testing.T, reader1 ion.Reader, reader2 ion.Reader) bool {
	next1 := reader1.Next()
	next2 := reader2.Next()

	if next1 != next2 {
		t.Errorf("next results don't match;\n" +
			"next1: %v,\n" +
			"next2: %v",
			next1, next2)
	}

	if !next1 {
		err := reader1.Err()
		if err != nil {
			t.Errorf("Something went wrong executing reader1.next(); %s", err.Error())
		}
	}

	if !next2 {
		err := reader2.Err()
		if err != nil {
			t.Errorf("Something went wrong executing reader2.next(); %s", err.Error())
		}
	}

	return next1 && next2
}

func compareFieldNames(t *testing.T, reader1 ion.Reader, reader2 ion.Reader) {
	// TODO: Add SymbolToken logic here once SymbolTokens are available
}

func compareNonNullStrings(t *testing.T, str1, str2 string, message string) bool {
	if str1 == "" || str2 == "" || str1 != str2 {
		t.Error(message)
	}

	return true
}

func compareAnnotations(t *testing.T, reader1 ion.Reader, reader2 ion.Reader) {
	// TODO: Add SymbolToken logic here once SymbolTokens are available

	annotations1 := reader1.Annotations()
	annotations2 := reader2.Annotations()

	if !reflect.DeepEqual(annotations1, annotations2) {
		t.Errorf("Expected symbol sequences to match;\n" +
			"Annotations #1: %v\n" +
			"Annotations #2; %v",
			annotations1, annotations2)
	}
}

func compareAnnotationSymbols(t *testing.T, reader1, reader2 ion.Reader) {
	// TODO: Add SymbolToken logic here once SymbolTokens are available
}

func compareHasAnnotations(t *testing.T, reader1, reader2 ion.Reader) {
	// TODO: Add SymbolToken logic here once SymbolTokens are available
}

func compareScalars(t *testing.T, reader1, reader2 ion.Reader) {
	ionType := reader1.Type()
	isNull := reader1.IsNull()

	switch ionType {
	case ion.BoolType:
		if !isNull {
			value1, err := reader1.BoolValue()
			if err != nil {
				t.Errorf("Something went wrong executing reader1.BoolValue(); %s", err.Error())
			}

			value2, err := reader2.BoolValue()
			if err != nil {
				t.Errorf("Something went wrong executing reader2.BoolValue(); %s", err.Error())
			}

			if value1 != value2 {
				t.Errorf("Expected bool values to match;\n" +
					"Bool #1: %v\n" +
					"Bool #2: %v",
					value1, value2)
			}
		}
	case ion.IntType:
		if !isNull {
			intSize, err := reader1.IntSize()
			if err != nil {
				t.Fatalf("Something went wrong executing reader1.IntSize(); %s", err.Error())
			}

			switch intSize {
			case ion.Int32, ion.Int64:
				int1, err := reader1.Int64Value()
				if err != nil {
					t.Errorf("Something went wrong executing reader1.Int64Value(); %s", err.Error())
				}

				int2, err := reader2.Int64Value()
				if err != nil {
					t.Errorf("Something went wrong executing reader2.Int64Value(); %s", err.Error())
				}

				if int1 != int2 {
					t.Errorf("Expected int values to match;\n"+
						"Int #1: %v\n"+
						"Int #2: %v",
						int1, int2)
				}
			case ion.Uint64:
				uint1, err := reader1.Uint64Value()
				if err != nil {
					t.Errorf("Something went wrong executing reader1.Uint64Value(); %s", err.Error())
				}

				uint2, err := reader2.Uint64Value()
				if err != nil {
					t.Errorf("Something went wrong executing reader2.Uint64Value(); %s", err.Error())
				}

				if uint1 != uint2 {
					t.Errorf("Expected uint values to match;\n"+
						"UInt #1: %v\n"+
						"UInt #2: %v",
						uint1, uint2)
				}
			case ion.BigInt:
				bigInt1, err := reader1.BigIntValue()
				if err != nil {
					t.Errorf("Something went wrong executing reader1.BigIntValue(); %s", err.Error())
				}

				bigInt2, err := reader2.BigIntValue()
				if err != nil {
					t.Errorf("Something went wrong executing reader2.BigIntValue(); %s", err.Error())
				}

				if bigInt1 != bigInt2 {
					t.Errorf("Expected big int values to match;\n" +
						"Big Int #1: %v\n" +
						"Big Int #2: %v",
						bigInt1, bigInt2)
				}
			default:
				t.Error("Expected intSize to be one of Int32, Int64, Uint64, or BigInt")
			}
		}
	case ion.FloatType:
		if !isNull {
			float1, err := reader1.FloatValue()
			if err != nil {
				t.Errorf("Something went wrong executing reader1.FloatValue(); %s", err.Error())
			}

			float2, err := reader2.FloatValue()
			if err != nil {
				t.Errorf("Something went wrong executing reader2.FloatValue(); %s", err.Error())
			}

			if math.IsNaN(float1) && math.IsNaN(float2) {
				if float1 != float2 {
					t.Errorf("Expected NaN float values to match")
				}
			} else if math.IsNaN(float1) || math.IsNaN(float2) {
				if float1 == float2 {
					t.Errorf("Expected float values to differ")
				}
			} else if float1 != float2 {
				t.Errorf("Expected float values to match;\n" +
					"Float #1: %v\n" +
					"Float #2: %v",
					float1, float2)
			}
		}
	case ion.DecimalType:
		if !isNull {
			decimal1, err := reader1.DecimalValue()
			if err != nil {
				t.Errorf("Something went wrong executing reader1.DecimalValue(); %s", err.Error())
			}

			decimal2, err := reader2.DecimalValue()
			if err != nil {
				t.Errorf("Something went wrong executing reader2.DecimalValue(); %s", err.Error())
			}

			decimalStrictEquals(t, decimal1, decimal2)
		}
	case ion.TimestampType:
		if !isNull {
			timestamp1, err := reader1.TimeValue()
			if err != nil {
				t.Errorf("Something went wrong executing reader1.TimeValue(); %s", err.Error())
			}

			timestamp2, err := reader2.TimeValue()
			if err != nil {
				t.Errorf("Something went wrong executing reader2.TimeValue(); %s", err.Error())
			}

			if timestamp1 != timestamp2 {
				t.Errorf("Expected timestamp values to match;\n" +
					"Timestamp #1: %v\n" +
					"Timestamp #2: %v",
					timestamp1, timestamp2)
			}
		}
	case ion.StringType:
		str1, err := reader1.StringValue()
		if err != nil {
			t.Errorf("Something went wrong executing reader1.StringValue(); %s", err.Error())
		}

		str2, err := reader2.StringValue()
		if err != nil {
			t.Errorf("Something went wrong executing reader2.StringValue(); %s", err.Error())
		}

		if str1 != str2 {
			t.Errorf("Expected string values to match;\n" +
				"String #1: %s\n" +
				"String #2: %s",
				str1, str2)
		}
	case ion.SymbolType:
		// TODO: Add SymbolToken logic here once SymbolTokens are available
	case ion.BlobType, ion.ClobType:
		if !isNull {
			b1, err := reader1.ByteValue()
			if err != nil {
				t.Errorf("Something went wrong executing reader1.ByteValue(); %s", err.Error())
			}

			b2, err := reader2.ByteValue()
			if err != nil {
				t.Errorf("Something went wrong executing reader2.ByteValue(); %s", err.Error())
			}

			if b1 == nil || b2 == nil {
				t.Errorf("Expected byte arrays to be non-null")
			}

			if len(b1) != len(b2) {
				t.Errorf("Expected byte arrays to have same length")
			}

			for i := 0; i < len(b1); i++ {
				t.Errorf("Expected byte arrays to match;\n" +
					"Array #1: %v\n" +
					"Array #2: %v",
					b1, b2)
			}
		}
	default:
		t.Error(InvalidIonTypeError{ionType})
	}
}

// decimalStrictEquals() compares two Ion Decimal values by equality and negative zero.
func decimalStrictEquals(t *testing.T, decimal1, decimal2 *ion.Decimal) {
	if decimal1 != decimal2 {
		t.Errorf("Expected decimal values to match;\n" +
			"Decimal #1: %v\n" +
			"Decimal #2: %v",
			decimal1, decimal2)
	}

	zeroDecimal := ion.NewDecimalInt(0)

	negativeZero1 := decimal1.Equal(zeroDecimal) && decimal1.Sign() < 0
	negativeZero2 := decimal2.Equal(zeroDecimal) && decimal2.Sign() < 0

	if negativeZero1 != negativeZero2 {
		t.Errorf("Expected decimal values to be both negative zero or both not negative zero")
	}

	if !decimal1.Equal(decimal2) {
		t.Errorf("Expected decimal Equal() to return true for given decimal values")
	}
}

// Read all the values in the reader and write them in the writer
func writeFromReaderToWriter(t *testing.T, reader ion.Reader, writer ion.Writer) {
	for reader.Next() {
		name := reader.FieldName()
		if name != "" {
			err := writer.FieldName(name)
			if err != nil {
				t.Fatalf("Something went wrong executing writer.FieldName(name); %s", err.Error())
			}
		}

		an := reader.Annotations()
		if len(an) > 0 {
			err := writer.Annotations(an...)
			if err != nil {
				t.Fatalf("Something went wrong executing writer.Annotations(an...); %s", err.Error())
			}
		}

		currentType := reader.Type()
		if reader.IsNull() {
			err := writer.WriteNullType(currentType)
			if err != nil {
				t.Fatalf("Something went wrong executing writer.WriteNullType(currentType); %s", err.Error())
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
			intSize, err := reader.IntSize()
			if err != nil {
				t.Fatalf("Something went wrong when retrieving the Int size; %s", err.Error())
			}

			switch intSize {
			case ion.Int32, ion.Int64:
				val, err := reader.Int64Value()
				if err != nil {
					t.Errorf("Something went wrong when reading Int value; %s", err.Error())
				}

				err = writer.WriteInt(val)
				if err != nil {
					t.Errorf("Something went wrong when writing Int value; %s", err.Error())
				}
			case ion.Uint64:
				val, err := reader.Uint64Value()
				if err != nil {
					t.Errorf("Something went wrong when reading UInt value; %s", err.Error())
				}

				err = writer.WriteUint(val)
				if err != nil {
					t.Errorf("Something went wrong when writing UInt value; %s", err.Error())
				}
			case ion.BigInt:
				val, err := reader.BigIntValue()
				if err != nil {
					t.Errorf("Something went wrong when reading Big Int value; %s", err.Error())
				}

				err = writer.WriteBigInt(val)
				if err != nil {
					t.Errorf("Something went wrong when writing Big Int value; %s", err.Error())
				}
			default:
				t.Error("Expected intSize to be one of Int32, Int64, Uint64, or BigInt")
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
				t.Fatalf("Something went wrong executing reader.StepIn(); %s", err.Error())
			}

			err = writer.BeginSexp()
			if err != nil {
				t.Fatalf("Something went wrong executing writer.BeginSexp(); %s", err.Error())
			}

			writeFromReaderToWriter(t, reader, writer)

			err = reader.StepOut()
			if err != nil {
				t.Fatalf("Something went wrong executing reader.StepOut(); %s", err.Error())
			}

			err = writer.EndSexp()
			if err != nil {
				t.Fatalf("Something went wrong executing writer.EndSexp(); %s", err.Error())
			}
		case ion.ListType:
			err := reader.StepIn()
			if err != nil {
				t.Fatalf("Something went wrong executing reader.StepIn(); %s", err.Error())
			}

			err = writer.BeginList()
			if err != nil {
				t.Fatalf("Something went wrong executing writer.BeginList(); %s", err.Error())
			}

			writeFromReaderToWriter(t, reader, writer)

			err = reader.StepOut()
			if err != nil {
				t.Fatalf("Something went wrong executing reader.StepOut(); %s", err.Error())
			}

			err = writer.EndList()
			if err != nil {
				t.Fatalf("Something went wrong executing writer.EndList(); %s", err.Error())
			}
		case ion.StructType:
			err := reader.StepIn()
			if err != nil {
				t.Fatalf("Something went wrong executing reader.StepIn(); %s", err.Error())
			}

			err = writer.BeginStruct()
			if err != nil {
				t.Fatalf("Something went wrong executing writer.BeginStruct(); %s", err.Error())
			}

			writeFromReaderToWriter(t, reader, writer)

			err = reader.StepOut()
			if err != nil {
				t.Fatalf("Something went wrong executing reader.StepOut(); %s", err.Error())
			}

			err = writer.EndStruct()
			if err != nil {
				t.Fatalf("Something went wrong executing writer.EndStruct(); %s", err.Error())
			}
		}
	}

	if reader.Err() != nil {
		t.Errorf("Something went wrong executing reader.Next(); %s", reader.Err().Error())
	}
}
