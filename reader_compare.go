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
	"reflect"

	"github.com/amzn/ion-go/ion"
)

func compareReaders(reader1 ion.Reader, reader2 ion.Reader) (bool, error) {
	for true {
		next, err := hasNext(reader1, reader2)
		if err != nil {
			return false, err
		}

		if !next {
			break
		}

		ionType1 := reader1.Type()
		ionType2 := reader2.Type()

		if ionType1 != ionType2 {
			return false, fmt.Errorf(
				"expected ion types to match; ionType1: %v, ionType2: %v", ionType1, ionType2)
		}

		ionHashReader, ok := reader2.(*hashReader)
		if ok {
			if ionHashReader.isInStruct() {
				compare, err := compareFieldNames(reader1, reader2)
				if !compare || err != nil {
					return compare, err
				}
			}
		} else {
			return false, fmt.Errorf("expected reader2 to be of type hashReader")
		}

		compare, err := compareAnnotations(reader1, reader2)
		if !compare || err != nil {
			return compare, err
		}

		compare, err = compareAnnotationSymbols(reader1, reader2)
		if !compare || err != nil {
			return compare, err
		}

		compare, err = compareHasAnnotations(reader1, reader2)
		if !compare || err != nil {
			return compare, err
		}

		isNull1 := reader1.IsNull()
		isNull2 := reader2.IsNull()
		if isNull1 != isNull2 {
			return false, fmt.Errorf(
				"expected readers to have same value for IsNull(); isNull1: %v, isNull2: %v", isNull1, isNull2)
		}

		switch ionType1 {
		case ion.NullType:
			if !isNull1 {
				return false, fmt.Errorf("expected ionType1 to be null")
			}
			if !isNull2 {
				return false, fmt.Errorf("expected ionType2 to be null")
			}
			break
		case ion.BoolType, ion.IntType, ion.FloatType, ion.DecimalType, ion.TimestampType,
			ion.StringType, ion.SymbolType, ion.BlobType, ion.ClobType:

			compare, err := compareScalars(reader1, reader2)
			if !compare || err != nil {
				return compare, err
			}
			break
		case ion.StructType, ion.ListType, ion.SexpType:
			err := reader1.StepIn()
			if err != nil {
				return false, err
			}

			err = reader2.StepIn()
			if err != nil {
				return false, err
			}

			compare, err := compareReaders(reader1, reader2)
			if !compare || err != nil {
				return compare, err
			}

			err = reader1.StepOut()
			if err != nil {
				return false, err
			}

			err = reader2.StepOut()
			if err != nil {
				return false, err
			}
			break
		default:
			return false, &InvalidIonTypeError{ionType1}
		}
	}

	next, err := hasNext(reader1, reader2)
	if err != nil {
		return false, err
	}

	if next {
		return false, fmt.Errorf("expected hasNext() to return false")
	}

	return true, nil
}

// hasNext() checks that the readers have a Next value
func hasNext(reader1 ion.Reader, reader2 ion.Reader) (bool, error) {
	next1 := reader1.Next()
	next2 := reader2.Next()

	if next1 != next2 {
		return false, fmt.Errorf(
			"next results don't match; reader1.Next(): %v, reader2.Next(): %v", next1, next2)
	}

	if !next1 {
		err := reader1.Err()
		if err != nil {
			return false, fmt.Errorf("expected reader1.next() not to error; %s", err.Error())
		}

		err = reader2.Err()
		if err != nil {
			return false, fmt.Errorf("expected reader2.next() not to error; %s", err.Error())
		}
	}

	return next1, nil
}

func compareFieldNames(reader1 ion.Reader, reader2 ion.Reader) (bool, error) {
	// TODO: Add SymbolToken logic here once SymbolTokens are available

	return true, nil
}

func compareNonNullStrings(str1, str2 string, message string) (bool, error) {
	if str1 == "" || str2 == "" || str1 != str2 {
		return false, fmt.Errorf(message)
	}

	return true, nil
}

func compareAnnotations(reader1 ion.Reader, reader2 ion.Reader) (bool, error) {
	// TODO: Add SymbolToken logic here once SymbolTokens are available

	if !reflect.DeepEqual(reader1.Annotations(), reader2.Annotations()) {
		return false, fmt.Errorf("symbol sequences don't match")
	}

	return true, nil
}

func compareAnnotationSymbols(reader1, reader2 ion.Reader) (bool, error) {
	// TODO: Add SymbolToken logic here once SymbolTokens are available

	return true, nil
}

func compareHasAnnotations(reader1, reader2 ion.Reader) (bool, error) {
	// TODO: Add SymbolToken logic here once SymbolTokens are available

	return true, nil
}

func compareScalars(reader1, reader2 ion.Reader) (bool, error) {
	ionType := reader1.Type()
	isNull := reader1.IsNull()

	switch ionType {
	case ion.BoolType:
		if !isNull {
			value1, err := reader1.BoolValue()
			if err != nil {
				return false, err
			}

			value2, err := reader2.BoolValue()
			if err != nil {
				return false, err
			}

			if value1 != value2 {
				return false, fmt.Errorf("expected readers to have matching bool values")
			}
		}
		break
	case ion.IntType:
		if !isNull {
			value1, err := reader1.BigIntValue()
			if err != nil {
				return false, err
			}

			value2, err := reader2.BigIntValue()
			if err != nil {
				return false, err
			}

			if value1 != value2 {
				return false, fmt.Errorf("expected readers to have matching big int values")
			}
		}
		break
	case ion.FloatType:
		if !isNull {
			v1, err := reader1.FloatValue()
			if err != nil {
				return false, err
			}

			v2, err := reader2.FloatValue()
			if err != nil {
				return false, err
			}

			if math.IsNaN(v1) && math.IsNaN(v2) {
				if v1 != v2 {
					return false, fmt.Errorf("expected readers to have matching float values")
				}
			} else if math.IsNaN(v1) || math.IsNaN(v2) {
				if v1 == v2 {
					return false, fmt.Errorf("expected readers to have different float values")
				}
			} else if v1 != v2 {
				return false, fmt.Errorf("expected readers to have matching float values")
			}
		}
		break
	case ion.DecimalType:
		if !isNull {
			decimal1, err := reader1.DecimalValue()
			if err != nil {
				return false, err
			}

			decimal2, err := reader2.DecimalValue()
			if err != nil {
				return false, err
			}

			check, err := decimalStrictEquals(decimal1, decimal2)
			if !check || err != nil {
				return check, err
			}
		}
		break
	case ion.TimestampType:
		if !isNull {
			timestamp1, err := reader1.TimeValue()
			if err != nil {
				return false, err
			}

			timestamp2, err := reader2.TimeValue()
			if err != nil {
				return false, err
			}

			if timestamp1 != timestamp2 {
				return false, fmt.Errorf("expected readers to have matching timestamp values")
			}
		}
		break
	case ion.StringType:
		value1, err := reader1.StringValue()
		if err != nil {
			return false, err
		}

		value2, err := reader2.StringValue()
		if err != nil {
			return false, err
		}

		if value1 != value2 {
			return false, fmt.Errorf("expected readers to have matching string values")
		}
		break
	case ion.SymbolType:
		// TODO: Add SymbolToken logic here once SymbolTokens are available
		break
	case ion.BlobType, ion.ClobType:
		if !isNull {
			b1, err := reader1.ByteValue()
			if err != nil {
				return false, err
			}

			b2, err := reader1.ByteValue()
			if err != nil {
				return false, err
			}

			if b1 == nil || b2 == nil {
				return false, fmt.Errorf("expected byte arrays to be not null")
			}

			if len(b1) != len(b2) {
				return false, fmt.Errorf("expected byte arrays to have same length")
			}

			for i := 0; i < len(b1); i++ {
				if b1[i] != b2[i] {
					return false, fmt.Errorf("expected byte arrays to have matching values")
				}
			}
		}
		break
	default:
		return false, &InvalidIonTypeError{ionType}
	}

	return true, nil
}

// decimalStrictEquals() compares two Ion Decimal values by equality and negative zero.
func decimalStrictEquals(decimal1, decimal2 *ion.Decimal) (bool, error) {
	if decimal1 != decimal2 {
		return false, fmt.Errorf("expected decimals to match")
	}

	zeroDecimal := ion.NewDecimalInt(0)

	negativeZero1 := decimal1.Equal(zeroDecimal) && decimal1.Sign() < 0
	negativeZero2 := decimal2.Equal(zeroDecimal) && decimal2.Sign() < 0

	if negativeZero1 != negativeZero2 {
		return false, fmt.Errorf("expected decimal values to be both negative zero or both not negative zero")
	}

	if !decimal1.Equal(decimal2) {
		return false, fmt.Errorf("expected decimal Equal() to return true for given decimal values")
	}

	return true, nil
}
