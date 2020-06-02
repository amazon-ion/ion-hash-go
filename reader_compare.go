package ionhash

import (
	"fmt"
	"math"
	"reflect"

	"github.com/amzn/ion-go/ion"
)

func Compare(reader1 ion.Reader, reader2 ion.Reader) error {
	for true {
		next, err := HasNext(reader1, reader2)
		if err != nil {
			return err
		}

		if !next {
			break
		}

		ionType1 := reader1.Type()
		ionType2 := reader2.Type()

		if ionType1 != ionType2 {
			return fmt.Errorf("expected ion types to match")
		}

		ionHashReader, ok := reader1.(*hashReader)
		if ok && ionHashReader.isInStruct() {
			err := CompareFieldNames(reader1, reader2)
			if err != nil {
				return err
			}
		}

		err = CompareAnnotations(reader1, reader2)
		if err != nil {
			return err
		}

		err = CompareAnnotationSymbols(reader1, reader2)
		if err != nil {
			return err
		}

		err = CompareHasAnnotations(reader1, reader2)
		if err != nil {
			return err
		}

		isNull1 := reader1.IsNull()
		isNull2 := reader2.IsNull()
		if isNull1 != isNull2 {
			return fmt.Errorf("expected readers to have same value for IsNull()")
		}

		switch ionType1 {
		case ion.NullType:
			if !isNull1 {
				return fmt.Errorf("expected ionType1 to be null")
			}
			if !isNull2 {
				return fmt.Errorf("expected ionType2 to be null")
			}
			break
		case ion.BoolType, ion.IntType, ion.FloatType, ion.DecimalType, ion.TimestampType,
			ion.StringType, ion.SymbolType, ion.BlobType, ion.ClobType:
			err := CompareScalars(ionType1, isNull1, reader1, reader2)
			if err != nil {
				return err
			}
			break
		case ion.StructType, ion.ListType, ion.SexpType:
			err := reader1.StepIn()
			if err != nil {
				return err
			}

			err = reader2.StepIn()
			if err != nil {
				return err
			}

			err = Compare(reader1, reader2)
			if err != nil {
				return err
			}

			err = reader1.StepOut()
			if err != nil {
				return err
			}

			err = reader2.StepOut()
			if err != nil {
				return err
			}
			break
		default:
			return &InvalidIonTypeError{ionType1}
		}
	}

	next, err := HasNext(reader1, reader2)
	if err != nil {
		return err
	}

	if next {
		return fmt.Errorf("expected HasNext() to return false")
	}

	return nil
}

func HasNext(reader1 ion.Reader, reader2 ion.Reader) (bool, error) {
	more := reader1.Next()

	if more != reader2.Next() {
		return false, fmt.Errorf("next results don't match")
	}

	if !more {
		if !reader1.Next() {
			return false, fmt.Errorf("expected reader1.Next() to return true")
		}
		if !reader2.Next() {
			return false, fmt.Errorf("expected reader2.Next() to return true")
		}
	}

	return more, nil
}

func CompareFieldNames(reader1 ion.Reader, reader2 ion.Reader) error {
	// TODO: Rework this once SymbolTokens are available
	/*token1 := reader1.GetFieldNameSymbol()
	token2 := reader2.GetFieldNameSymbol()

	fieldName := token1.Text
	if fieldName != token2.Text {
		return fmt.Errorf("tokens don't match")
	}

	if fieldName != "" {
		field1 := reader1.FieldName()
		field2 := reader2.FieldName()
		return CompareNonNullStrings(field1, field2, "field names don't match")
	}*/

	return nil
}

func CompareNonNullStrings(str1, str2 string, message string) error {
	if str1 == "" || str2 == "" || str1 != str2 {
		return fmt.Errorf(message)
	}

	return nil
}

func CompareAnnotations(reader1 ion.Reader, reader2 ion.Reader) error {
	// TODO: Rework this once SymbolTokens are available
	/*symbols := reader1.GetTypeAnnotationSymbols()

	//Skip comparison if any annotation is zero symbol
	for _, symbol := range symbols {
		if symbol.Text == nil && symbol.Sid == 0 {
			return nil
		}
	}*/

	if !reflect.DeepEqual(reader1.Annotations(), reader2.Annotations()) {
		return fmt.Errorf("symbol sequences don't match")
	}

	return nil
}

func CompareAnnotationSymbols(reader1, reader2 ion.Reader) error {
	// TODO: Rework this once SymbolTokens are available
	/*if !reflect.DeepEqual(reader1.GetTypeAnnotationSymbols(), reader2.GetTypeAnnotationSymbols()) {
		return fmt.Errorf("expected type annotation symbols to match")
	}*/

	return nil
}

func CompareHasAnnotations(reader1, reader2 ion.Reader) error {
	// TODO: Rework this once SymbolTokens are available
	/*symbols := reader1.GetTypeAnnotationSymbols()

	//Skip comparison if any annotation is zero symbol
	for _, symbol := range symbols {
		if symbol.Text == nil && symbol.Sid == 0 {
			return nil
		}
	}

	annotations1 := reader1.Annotations()
	annotations2 := reader2.Annotations()

	if len(annotations1) != len(annotations2) {
		return fmt.Errorf("expected annotation sequences to have the same length")
	}

	for i := 0; i < len(annotations1); i++ {
		if reader1.HasAnnotation(annotations2[i]) {
			return fmt.Errorf("expected reader1 to have reader2's annotation")
		}
		if reader2.HasAnnotation(annotations1[i]) {
			return fmt.Errorf("expected reader2 to have reader1's annotation")
		}
	}*/

	return nil
}

func CompareScalars(ionType ion.Type, isNull bool, reader1, reader2 ion.Reader) error {
	switch ionType {
	case ion.BoolType:
		if !isNull && reader1.BoolValue() != reader2.BoolValue() {
			return fmt.Errorf("expected readers to have matching bool values")
		}
		break
	case ion.IntType:
		if !isNull && reader1.BigIntValue() != reader2.BigIntValue() {
			return fmt.Errorf("expected readers to have matching big int values")
		}
		break
	case ion.FloatType:
		if !isNull {
			v1, err := reader1.FloatValue()
			if err != nil {
				return err
			}

			v2, err := reader2.FloatValue()
			if err != nil {
				return err
			}

			if math.IsNaN(v1) && math.IsNaN(v2) {
				if v1 != v2 {
					return fmt.Errorf("expected readers to have matching float values")
				}
			} else if math.IsNaN(v1) || math.IsNaN(v2) {
				if v1 == v2 {
					return fmt.Errorf("expected readers to have different float values")
				}
			} else if v1 != v2 {
				return fmt.Errorf("expected readers to have matching float values")
			}
		}
		break
	case ion.DecimalType:
		if !isNull {
			decimal1, err := reader1.DecimalValue()
			if err != nil {
				return err
			}

			decimal2, err := reader2.DecimalValue()
			if err != nil {
				return err
			}

			err = CheckPreciselyEquals(decimal1, decimal2)
			if err != nil {
				return err
			}
		}
		break
	case ion.TimestampType:
		if !isNull {
			timestamp1, err := reader1.TimeValue()
			if err != nil {
				return err
			}

			timestamp2, err := reader2.TimeValue()
			if err != nil {
				return err
			}

			if timestamp1 != timestamp2 {
				return fmt.Errorf("expected readers to have matching timestamp values")
			}
		}
		break
	case ion.StringType:
		if reader1.StringValue() != reader2.StringValue() {
			return fmt.Errorf("expected readers to have matching string values")
		}
		break
	case ion.SymbolType:
		// TODO: Rework this once SymbolTokens are available
		/*token1 := reader1.SymbolValue()
		token2 := reader2.SymbolValue()

		if isNull {
			if token1.Text != nil {
				return fmt.Errorf("expected token to be null")
			}
			if token2.Text != nil {
				return fmt.Errorf("expected token to be null")
			}
		} else if token1.Text == nil || token2.Text == nil {
			if token1.Sid != token2.Sid {
				return fmt.Errorf("expected token SIDs to match")
			}
		} else if token1.Text != token2.Text {
			return fmt.Errorf("expected tokens to match")
		}*/
		break
	case ion.BlobType, ion.ClobType:
		if !isNull {
			b1, err := reader1.ByteValue()
			if err != nil {
				return err
			}

			b2, err := reader1.ByteValue()
			if err != nil {
				return err
			}

			if b1 == nil || b2 == nil {
				return fmt.Errorf("expected byte arrays to be not null")
			}

			if len(b1) != len(b2) {
				return fmt.Errorf("expected byte arrays to have same length")
			}

			for i := 0; i < len(b1); i++ {
				if b1[i] != b2[i] {
					return fmt.Errorf("expected byte arrays to have matching values")
				}
			}
		}
		break
	default:
		return &InvalidIonTypeError{ionType}
	}

	return nil
}

func CheckPreciselyEquals(decimal1, decimal2 *ion.Decimal) error {
	if decimal1 != decimal2 {
		return fmt.Errorf("expected decimals to match")
	}

	zeroDecimal := ion.NewDecimalInt(0)

	expectedNegativeZero := decimal1.Equal(zeroDecimal) && decimal1.Sign() < 0
	actualNegativeZero := decimal2.Equal(zeroDecimal) && decimal2.Sign() < 0

	if expectedNegativeZero != actualNegativeZero {
		return fmt.Errorf("expected decimal values to be both negative zero or both not negatve zero")
	}

	if !decimal1.Equal(decimal2) {
		return fmt.Errorf("expected decimal Equal() to return true for given decimal values")
	}

	return nil
}
