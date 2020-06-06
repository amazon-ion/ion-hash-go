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

func Compare(reader1 ion.Reader, reader2 ion.Reader) (bool, error) {
	for true {
		next, err := HasNext(reader1, reader2)
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

		ionHashReader, ok := reader1.(*hashReader)
		if ok && ionHashReader.isInStruct() {
			compare, err := CompareFieldNames(reader1, reader2)
			if !compare || err != nil {
				return compare, err
			}
		}

		compare, err := CompareAnnotations(reader1, reader2)
		if !compare || err != nil {
			return compare, err
		}

		compare, err = CompareAnnotationSymbols(reader1, reader2)
		if !compare || err != nil {
			return compare, err
		}

		compare, err = CompareHasAnnotations(reader1, reader2)
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
			compare, err := CompareScalars(ionType1, isNull1, reader1, reader2)
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

			compare, err := Compare(reader1, reader2)
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

	next, err := HasNext(reader1, reader2)
	if err != nil {
		return false, err
	}

	if next {
		return false, fmt.Errorf("expected HasNext() to return false")
	}

	return true, nil
}

func HasNext(reader1 ion.Reader, reader2 ion.Reader) (bool, error) {
	next1 := reader1.Next()
	next2 := reader2.Next()

	if next1 != next2 {
		return false, fmt.Errorf("next results don't match; reader1.Next(): %v, reader2.Next(): %v", next1, next2)
	}

	if !next1 {
		err := reader1.Err()
		if err != nil {
			return false, fmt.Errorf("expected reader1.next() not to error; %s", err.Error())
		}

		err = reader2.Err()
		if err != nil {
			return false, fmt.Errorf("expected reader1.next() not to error; %s", err.Error())
		}
	}

	return next1, nil
}

func CompareFieldNames(reader1 ion.Reader, reader2 ion.Reader) (bool, error) {
	// TODO: Rework this once SymbolTokens are available
	/*token1 := reader1.GetFieldNameSymbol()
	token2 := reader2.GetFieldNameSymbol()

	tokenText1 := token1.Text
	tokenText2 := token2.Text
	if tokenText1 != tokenText2 {
		return false, fmt.Errorf("tokens don't match; token1: %s, token2: %s", tokenText1, tokenText2)
	}

	if tokenText1 != "" {
		field1 := reader1.FieldName()
		field2 := reader2.FieldName()
		return CompareNonNullStrings(field1, field2,
			fmt.Sprintf("field names don't match; field1: %s, field2: %s", field1, field2))
	}*/

	return true, nil
}

func CompareNonNullStrings(str1, str2 string, message string) (bool, error) {
	if str1 == "" || str2 == "" || str1 != str2 {
		return false, fmt.Errorf(message)
	}

	return true, nil
}

func CompareAnnotations(reader1 ion.Reader, reader2 ion.Reader) (bool, error) {
	// TODO: Rework this once SymbolTokens are available
	/*symbols := reader1.GetTypeAnnotationSymbols()

	//Skip comparison if any annotation is zero symbol
	for _, symbol := range symbols {
		if symbol.Text == nil && symbol.Sid == 0 {
			return true, nil
		}
	}*/

	if !reflect.DeepEqual(reader1.Annotations(), reader2.Annotations()) {
		return false, fmt.Errorf("symbol sequences don't match")
	}

	return true, nil
}

func CompareAnnotationSymbols(reader1, reader2 ion.Reader) (bool, error) {
	// TODO: Rework this once SymbolTokens are available
	/*if !reflect.DeepEqual(reader1.GetTypeAnnotationSymbols(), reader2.GetTypeAnnotationSymbols()) {
		return false, fmt.Errorf("expected type annotation symbols to match")
	}*/

	return true, nil
}

func CompareHasAnnotations(reader1, reader2 ion.Reader) (bool, error) {
	// TODO: Rework this once SymbolTokens are available
	/*symbols := reader1.GetTypeAnnotationSymbols()

	//Skip comparison if any annotation is zero symbol
	for _, symbol := range symbols {
		if symbol.Text == nil && symbol.Sid == 0 {
			return true, nil
		}
	}

	annotations1 := reader1.Annotations()
	annotations2 := reader2.Annotations()

	if len(annotations1) != len(annotations2) {
		return false, fmt.Errorf("expected annotation sequences to have the same length")
	}

	for i := 0; i < len(annotations1); i++ {
		if reader1.HasAnnotation(annotations2[i]) {
			return false, fmt.Errorf("expected reader1 to have reader2's annotation")
		}
		if reader2.HasAnnotation(annotations1[i]) {
			return false, fmt.Errorf("expected reader2 to have reader1's annotation")
		}
	}*/

	return true, nil
}

func CompareScalars(ionType ion.Type, isNull bool, reader1, reader2 ion.Reader) (bool, error) {
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

			check, err := CheckPreciselyEquals(decimal1, decimal2)
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
		// TODO: Rework this once SymbolTokens are available
		/*token1 := reader1.SymbolValue()
		token2 := reader2.SymbolValue()

		if isNull {
			if token1.Text != nil {
				return false, fmt.Errorf("expected token to be null")
			}
			if token2.Text != nil {
				return false, fmt.Errorf("expected token to be null")
			}
		} else if token1.Text == nil || token2.Text == nil {
			if token1.Sid != token2.Sid {
				return false, fmt.Errorf("expected token SIDs to match")
			}
		} else if token1.Text != token2.Text {
			return false, fmt.Errorf("expected tokens to match")
		}*/
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

func CheckPreciselyEquals(decimal1, decimal2 *ion.Decimal) (bool, error) {
	if decimal1 != decimal2 {
		return false, fmt.Errorf("expected decimals to match")
	}

	zeroDecimal := ion.NewDecimalInt(0)

	expectedNegativeZero := decimal1.Equal(zeroDecimal) && decimal1.Sign() < 0
	actualNegativeZero := decimal2.Equal(zeroDecimal) && decimal2.Sign() < 0

	if expectedNegativeZero != actualNegativeZero {
		return false, fmt.Errorf("expected decimal values to be both negative zero or both not negative zero")
	}

	if !decimal1.Equal(decimal2) {
		return false, fmt.Errorf("expected decimal Equal() to return true for given decimal values")
	}

	return true, nil
}
