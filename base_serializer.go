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
	"fmt"
	"time"

	"github.com/amzn/ion-go/ion"
)

// Holds the commonalities between scalar and struct serializers.
type baseSerializer struct {
	hashFunction           IonHasher
	depth                  int
	hasContainerAnnotation bool
}

func (baseSerializer *baseSerializer) stepIn(ionValue hashValue) error {
	err := baseSerializer.handleFieldName(ionValue)
	if err != nil {
		return err
	}

	err = baseSerializer.handleAnnotationsBegin(ionValue, true)
	if err != nil {
		return err
	}

	baseSerializer.beginMarker()

	tq := tq(ionValue)
	if ionValue.isNull() {
		tq = tq | 0x0F
	}

	baseSerializer.update([]byte{tq})

	return nil
}

func (baseSerializer *baseSerializer) Sum(b []byte) []byte {
	return baseSerializer.hashFunction.Sum(b)
}

func (baseSerializer *baseSerializer) handleFieldName(ionValue hashValue) error {
	if baseSerializer.depth > 0 && ionValue.isInStruct() {
		// SymbolTokens are not available right now.
		/*if ionValue.fieldNameSymbol().Text == null && ionValue.fieldNameSymbol().Sid != 0 {
			return &UnknownSymbolError{ionValue.fieldNameSymbol().Sid}
		}
		return baseSerializer.writeSymbol(ionValue.fieldNameSymbol().Text)*/
	}

	return nil
}

func (baseSerializer *baseSerializer) update(bytes []byte) {
	baseSerializer.hashFunction.Update(bytes)
}

func (baseSerializer *baseSerializer) beginMarker() {
	baseSerializer.hashFunction.Update([]byte{BeginMarkerByte})
}

func (baseSerializer *baseSerializer) endMarker() {
	baseSerializer.hashFunction.Update([]byte{EndMarkerByte})
}

func (baseSerializer *baseSerializer) handleAnnotationsBegin(ionValue hashValue, isContainer bool) error {
	annotations := ionValue.getAnnotations()
	if len(annotations) > 0 {
		baseSerializer.beginMarker()
		baseSerializer.update([]byte{TqValue})
		for _, annotation := range annotations {
			err := baseSerializer.writeSymbol(annotation)
			if err != nil {
				return err
			}
		}

		if isContainer {
			baseSerializer.hasContainerAnnotation = true
		}
	}

	return nil
}

func (baseSerializer *baseSerializer) handleAnnotationsEnd(ionValue hashValue, isContainer bool) {
	if (ionValue != nil && len(ionValue.getAnnotations()) > 0) ||
		(isContainer && baseSerializer.hasContainerAnnotation) {
		baseSerializer.endMarker()

		if isContainer {
			baseSerializer.hasContainerAnnotation = false
		}
	}
}

func (baseSerializer *baseSerializer) writeSymbol(token string) error {
	baseSerializer.beginMarker()

	// SymbolTokens are not available right now.
	/*var sid int
	if token == "" {
		sid = 0
	} else {
		sid = ion.SymbolToken.UnknownSid
	}

	symbolToken := &ion.SymbolToken{token, sid}
	scalarBytes, err := baseSerializer.getBytes(ion.SymbolType, symbolToken, false);
	if err != nil {
		return nil
	}

	tq, representation := baseSerializer.scalarOrNullSplitParts(ion.SymbolType, symbolToken, false, scalarBytes)

	baseSerializer.update([]byte{tq})
	if len(representation) > 0 {
		baseSerializer.update(escape(representation))
	}*/

	baseSerializer.endMarker()

	return nil
}

func (baseSerializer *baseSerializer) getBytes(ionType ion.Type, ionValue interface{}, isNull bool) ([]byte, error) {
	if isNull {
		typeCode := byte(ionType)
		return []byte{(typeCode << 4) | 0x0F}, nil
	} else if ionType == ion.FloatType && ionValue == 0 && int64(ionValue.(float64)) >= 0 {
		// value is 0.0, not -0.0
		return []byte{0x40}, nil
	} else {
		buf := bytes.Buffer{}
		writer := NewHashWriter(ion.NewBinaryWriter(&buf), NewCryptoHasherProvider(SHA256))

		err := serializers(ionType, ionValue, writer)
		if err != nil {
			return nil, err
		}

		err = writer.Finish()
		if err != nil {
			return nil, err
		}

		return buf.Bytes()[4:], nil
	}
}

func (baseSerializer *baseSerializer) getLengthLength(bytes []byte) (int, error) {
	if (bytes[0] & 0x0F) == 0x0E {
		// read subsequent byte(s) as the "length" field
		for i := 1; i < len(bytes); i++ {
			if (bytes[i] & 0x80) != 0 {
				return i, nil
			}
		}

		return 0, fmt.Errorf("problem while reading VarUInt")
	}

	return 0, nil
}

// SymbolToken is currently not available
//func (baseSerializer *baseSerializer)scalarOrNullSplitParts(ionType ion.Type, symbolToken ion.SymbolToken, isNull bool, bytes byte[]) (byte, []byte) {
//	panic("implement me")
//}

func escape(bytes []byte) []byte {
	if bytes == nil {
		return nil
	}

	for i := 0; i < len(bytes); i++ {
		b := bytes[i]
		if b == BeginMarkerByte || b == EndMarkerByte || b == EscapeByte {
			// found a byte that needs to be escaped; build a new byte array that
			// escapes that byte as well as any others
			var escapedBytes []byte

			for j := 0; j < len(bytes); j++ {
				c := bytes[j]
				if c == BeginMarkerByte || c == EndMarkerByte || c == EscapeByte {
					escapedBytes = append(escapedBytes, EscapeByte)
				}

				escapedBytes = append(escapedBytes, c)
			}

			return escapedBytes
		}
	}

	return bytes
}

func serializers(ionType ion.Type, ionValue interface{}, writer HashWriter) error {
	switch ionType {
	case ion.BoolType:
		return writer.WriteBool(ionValue.(bool))
	case ion.BlobType:
		return writer.WriteBlob(ionValue.([]byte))
	case ion.ClobType:
		return writer.WriteClob(ionValue.([]byte))
	case ion.DecimalType:
		return writer.WriteDecimal(ionValue.(*ion.Decimal))
	case ion.FloatType:
		return writer.WriteFloat(ionValue.(float64))
	case ion.IntType:
		return writer.WriteInt(ionValue.(int64))
	case ion.StringType:
		return writer.WriteString(ionValue.(string))
	case ion.SymbolType:
		return writer.WriteSymbol(ionValue.(string))
	case ion.TimestampType:
		return writer.WriteTimestamp(ionValue.(time.Time))
	case ion.NullType:
		return writer.WriteNull()
	}

	return &InvalidIonTypeError{ionType}
}

func tq(ionValue hashValue) byte {
	typeCode := byte(ionValue.ionType())
	return typeCode << 4
}
