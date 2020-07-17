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
	"math/big"
	"time"

	"github.com/amzn/ion-go/ion"
)

// Holds the commonalities between scalar and struct serializers.
type baseSerializer struct {
	hashFunction           IonHasher
	depth                  int
	hasContainerAnnotation bool
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

	err = baseSerializer.handleAnnotationsBegin(ionValue, true)
	if err != nil {
		return err
	}

	err = baseSerializer.beginMarker()
	if err != nil {
		return err
	}

	tq := typeQualifier(ionValue)
	if ionValue.IsNull() {
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
	if baseSerializer.depth > 0 && ionValue.IsInStruct() {
		fieldName := ionValue.getFieldName()

		// TODO: Add logic returning UnknownSymbolError once SymbolToken is available

		if fieldName != nil {
			return baseSerializer.writeSymbol(*fieldName)
		}
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

func (baseSerializer *baseSerializer) handleAnnotationsBegin(ionValue hashValue, isContainer bool) error {
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

		if isContainer {
			baseSerializer.hasContainerAnnotation = true
		}
	}

	return nil
}

func (baseSerializer *baseSerializer) handleAnnotationsEnd(ionValue hashValue, isContainer bool) error {
	if (ionValue != nil && len(ionValue.getAnnotations()) > 0) ||
		(isContainer && baseSerializer.hasContainerAnnotation) {

		err := baseSerializer.endMarker()
		if err != nil {
			return err
		}

		if isContainer {
			baseSerializer.hasContainerAnnotation = false
		}
	}

	return nil
}

func (baseSerializer *baseSerializer) writeSymbol(token string) error {
	err := baseSerializer.beginMarker()
	if err != nil {
		return err
	}

	// TODO: Add SymbolToken logic here once SymbolTokens are available

	scalarBytes, err := baseSerializer.getBytes(ion.SymbolType, token, false)
	if err != nil {
		return err
	}

	tq, representation, err := baseSerializer.scalarOrNullSplitParts(ion.SymbolType, false, scalarBytes)
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
	}

	err = baseSerializer.endMarker()
	if err != nil {
		return err
	}

	return nil
}

func (baseSerializer *baseSerializer) getBytes(ionType ion.Type, ionValue interface{}, isNull bool) ([]byte, error) {
	if isNull {
		var typeCode byte
		if ionType <= ion.IntType {
			// The Ion binary encodings of NoType, NullType, BoolType, and IntType
			// differ from their enum values by one
			typeCode = byte(ionType - 1)
		} else {
			typeCode = byte(ionType)
		}

		return []byte{(typeCode << 4) | 0x0F}, nil
	} else if ionType == ion.FloatType && ionValue == 0 && int64(ionValue.(float64)) >= 0 {
		// value is 0.0, not -0.0
		return []byte{0x40}, nil
	} else {
		buf := bytes.Buffer{}
		writer := ion.NewBinaryWriter(&buf)

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

func (baseSerializer *baseSerializer) scalarOrNullSplitParts(
	ionType ion.Type, isNull bool, bytes []byte) (byte, []byte, error) {

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
		}

		// TODO: Add SymbolToken logic here once SymbolTokens are available

	} else if ionType != ion.BoolType && (tq&0x0F) != 0x0F {
		// not a symbol, bool, or null value
		// zero - out the L nibble
		tq = tq & 0xF0
	}

	return tq, representation, nil
}

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

func serializers(ionType ion.Type, ionValue interface{}, writer ion.Writer) error {
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
		ionValInt64, ok := ionValue.(int64)
		if ok {
			return writer.WriteInt(ionValInt64)
		}

		ionValUint64, ok := ionValue.(uint64)
		if ok {
			return writer.WriteUint(ionValUint64)
		}

		ionValInt32, ok := ionValue.(int32)
		if ok {
			return writer.WriteInt(int64(ionValInt32))
		}

		ionValUint32, ok := ionValue.(uint32)
		if ok {
			return writer.WriteUint(uint64(ionValUint32))
		}

		ionValInt, ok := ionValue.(int)
		if ok {
			return writer.WriteInt(int64(ionValInt))
		}

		ionValBigInt, ok := ionValue.(*big.Int)
		if ok {
			return writer.WriteBigInt(ionValBigInt)
		}

		return &InvalidArgumentError{"ionValue", ionValue}
	case ion.StringType:
		return writer.WriteString(ionValue.(string))
	case ion.SymbolType:
		ionValueStr, ok := ionValue.(string)
		if ok {
			return writer.WriteString(ionValueStr)
		}

		ionValueSymbol, ok := ionValue.(ion.SymbolTable)
		if ok {
			symbols := ionValueSymbol.Symbols()
			if len(symbols) > 0 {
				return writer.WriteString(symbols[0])
			}
		}

		return &InvalidArgumentError{"ionValue", ionValue}
	case ion.TimestampType:
		return writer.WriteTimestamp(ionValue.(time.Time))
	case ion.NullType:
		return writer.WriteNull()
	}

	return &InvalidIonTypeError{ionType}
}

func typeQualifier(ionValue hashValue) byte {
	typeCode := byte(ionValue.Type())
	return typeCode << 4
}

func compareBytes(bytes1, bytes2 []byte) int {
	for i := 0; i < len(bytes1) && i < len(bytes2); i++ {
		byte1 := bytes1[i]
		byte2 := bytes2[i]
		if byte1 != byte2 {
			return int(byte1) - int(byte2)
		}
	}

	return len(bytes1) - len(bytes2)
}

// sortableBytes implements the sort.Interface so we can sort fieldHashes
type sortableBytes [][]byte

func (sb sortableBytes) Len() int {
	return len(sb)
}

func (sb sortableBytes) Less(i, j int) bool {
	bytes1 := sb[i]
	bytes2 := sb[j]

	return compareBytes(bytes1, bytes2) < 0
}

func (sb sortableBytes) Swap(i, j int) {
	sb[i], sb[j] = sb[j], sb[i]
}
