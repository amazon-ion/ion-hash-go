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
	hashFunction IonHasher
	depth        int
}

func (baseSerializer *baseSerializer) stepOut() error {
	err := baseSerializer.endMarker()
	if err != nil {
		return err
	}

	err = baseSerializer.handleAnnotationsEnd(nil, true)
	if err != nil {
		return err
	}

	return nil
}

func (baseSerializer *baseSerializer) stepIn(ionValue hashValue) error {
	err := baseSerializer.handleFieldName(ionValue)
	if err != nil {
		return err
	}

	err = baseSerializer.handleAnnotationsBegin(ionValue)
	if err != nil {
		return err
	}

	err = baseSerializer.beginMarker()
	if err != nil {
		return err
	}

	tq := typeQualifier(ionValue)
	if ionValue.isNull() {
		tq = tq | 0x0F
	}

	err = baseSerializer.write([]byte{tq})
	if err != nil {
		return err
	}

	return nil
}

func (baseSerializer *baseSerializer) sum(b []byte) []byte {
	return baseSerializer.hashFunction.Sum(b)
}

func (baseSerializer *baseSerializer) handleFieldName(ionValue hashValue) error {
	if baseSerializer.depth > 0 && ionValue.isInStruct() {
		// TODO: Rework this once SymbolTokens become available
		/*if ionValue.fieldNameSymbol().Text == null && ionValue.fieldNameSymbol().Sid != 0 {
			return &UnknownSymbolError{ionValue.fieldNameSymbol().Sid}
		}
		return baseSerializer.writeSymbol(ionValue.fieldNameSymbol().Text)*/
	}

	return nil
}

func (baseSerializer *baseSerializer) write(bytes []byte) error {
	_, err := baseSerializer.hashFunction.Write(bytes)
	return err
}

func (baseSerializer *baseSerializer) beginMarker() error {
	_, err := baseSerializer.hashFunction.Write([]byte{beginMarkerByte})
	return err
}

func (baseSerializer *baseSerializer) endMarker() error {
	_, err := baseSerializer.hashFunction.Write([]byte{endMarkerByte})
	return err
}

func (baseSerializer *baseSerializer) handleAnnotationsBegin(ionValue hashValue) error {
	if ionValue == nil {
		return &InvalidArgumentError{"ionValue", ionValue}
	}

	annotations := ionValue.getAnnotations()
	if len(annotations) > 0 {
		err := baseSerializer.beginMarker()
		if err != nil {
			return err
		}

		err = baseSerializer.write([]byte{tqValue})
		if err != nil {
			return err
		}

		for _, annotation := range annotations {
			err = baseSerializer.writeSymbol(annotation)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (baseSerializer *baseSerializer) handleAnnotationsEnd(ionValue hashValue, isContainer bool) error {
	if (ionValue != nil && len(ionValue.getAnnotations()) > 0) || isContainer {
		err := baseSerializer.endMarker()
		if err != nil {
			return err
		}
	}

	return nil
}

func (baseSerializer *baseSerializer) writeSymbol(token string) error {
	err := baseSerializer.beginMarker()
	if err != nil {
		return err
	}

	// TODO: Rework this once SymbolTokens become available
	/*var sid int
	if token == "" {
		sid = 0
	} else {
		sid = ion.SymbolToken.UnknownSid
	}

	symbolToken := &ion.SymbolToken{token, sid}
	scalarBytes, err := baseSerializer.getBytes(ion.SymbolType, symbolToken, false);
	if err != nil {
		return err
	}

	tq, representation, err := baseSerializer.scalarOrNullSplitParts(ion.SymbolType, symbolToken, false, scalarBytes)
	if err != nil {
		return err
	}

	err = baseSerializer.write([]byte{tq})
	if err != nil {
		return err
	}

	if len(representation) > 0 {
		err = baseSerializer.write(escape(representation))
		if err != nil {
			return err
		}
	}*/

	err = baseSerializer.endMarker()
	if err != nil {
		return err
	}

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
		writer, err := NewHashWriter(ion.NewBinaryWriter(&buf), newCryptoHasherProvider(SHA256))
		if err != nil {
			return nil, err
		}

		err = serializers(ionType, ionValue, writer)
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

func (baseSerializer *baseSerializer) getLengthFieldLength(bytes []byte) (int, error) {
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

// TODO: Rework this once SymbolTokens become available
/*func (baseSerializer *baseSerializer)scalarOrNullSplitParts(
	ionType ion.Type, symbolToken ion.SymbolToken, isNull bool, bytes []byte) (byte, []byte, error) {

	offset, err := baseSerializer.getLengthFieldLength(bytes)
	if err != nil {
		return byte(0), nil, err
	}
	offset++

	if ionType == ion.IntType && len(bytes) > offset {
		// ignore sign byte when the magnitude ends at byte boundary
		if (bytes[offset] & 0xFF) == 0 {
			offset++
		}
	}

	// the representation is everything after TL (first byte) and length
	representation := bytes[offset:]
	tq := bytes[0]

	if ionType == ion.SymbolType {
		// symbols are serialized as strings; use the correct TQ:
		tq = 0x70
		if isNull {
			tq = tq | 0x0F
		} else if symbolToken != nil && symbolToken.Value.Text == nil && symbolToken.Value.Sid == 0 {
			tq = 0x71
		}
	// not a symbol, bool, or null value
	} else if ionType != ion.BoolType && (tq & 0x0F) != 0x0F {
		// zero - out the L nibble
		tq = tq & 0xF0
	}

	return tq, representation, nil
}*/

func needsEscape(b byte) bool {
	switch b {
	case beginMarkerByte, endMarkerByte, escapeByte:
		return true
	}

	return false
}

func escape(bytes []byte) []byte {
	if bytes == nil {
		return nil
	}

	for i := 0; i < len(bytes); i++ {
		b := bytes[i]
		if needsEscape(b) {
			// found a byte that needs to be escaped; build a new byte array that
			// escapes that byte as well as any others
			var escapedBytes []byte

			for j := 0; j < len(bytes); j++ {
				c := bytes[j]
				if needsEscape(c) {
					escapedBytes = append(escapedBytes, escapeByte)
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

func typeQualifier(ionValue hashValue) byte {
	typeCode := byte(ionValue.ionType())
	return typeCode << 4
}
