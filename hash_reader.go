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
	"math/big"

	"github.com/amzn/ion-go/ion"
)

// A HashReader reads a stream of Ion values and calculates its hash.
//
// The HashReader has a logical position within the stream of values, influencing the
// values returned from its methods. Initially, the HashReader is positioned before the
// first value in the stream. A call to Next advances the HashReader to the first value
// in the stream, with subsequent calls advancing to subsequent values. When a call to
// Next moves the HashReader to the position after the final value in the stream, it returns
// false, making it easy to loop through the values in a stream. e.g.,
//
// 	   var r HashReader
// 	   for r.Next() {
// 		   // ...
// 	   }
//
// Next also returns false in case of error. This can be distinguished from a legitimate
// end-of-stream by calling HashReader.Err after exiting the loop.
//
// When positioned on an Ion value, the type of the value can be retrieved by calling
// HashReader.Type. If it has an associated field name (inside a struct) or annotations, they can
// be read by calling HashReader.FieldName and HashReader.Annotations respectively.
//
// For atomic values, an appropriate XxxValue method can be called to read the value.
// For lists, sexps, and structs, you should instead call HashReader.StepIn to move the HashReader in
// to the contained sequence of values. The HashReader will initially be positioned before
// the first value in the container. Calling HashReader.Next without calling HashReader.StepIn will skip over
// the composite value and return the next value in the outer value stream.
//
// At any point while reading through a composite value, including when HashReader.Next returns false
// to indicate the end of the contained values, you may call HashReader.StepOut to move back to the
// outer sequence of values. The HashReader will be positioned at the end of the composite value,
// such that a call to HashReader.Next will move to the immediately-following value (if any).
//
// HashReader.Sum will return the hash of the entire stream of Ion values that have been seen thus far, e.g.,
//
//     cryptoHasherProvider := NewCryptoHasherProvider(SHA256)
// 	   r := NewTextReaderStr("[foo, bar] [baz]")
//     hr := NewHashReader(r, cryptoHasherProvider)
// 	   for hr.Next() {
// 		   if err := hr.StepIn(); err != nil {
// 			   return err
// 		   }
// 		   for hr.Next() {
// 			   fmt.Println(hr.StringValue())
// 		   }
// 		   if err := hr.StepOut(); err != nil {
// 			   return err
// 		   }
// 	   }
// 	   if err := hr.Err(); err != nil {
// 		   return err
// 	   }
//
//     fmt.Printf("%v", hr.Sum(nil))
//
type HashReader interface {
	// Embed interface of Ion reader.
	ion.Reader

	// hashValue methods.
	getFieldName() (*ion.SymbolToken, error)
	getAnnotations() ([]ion.SymbolToken, error)
	value() (interface{}, error)

	// Sum appends the current hash to b and returns the resulting slice.
	// It resets the Hash to its initial state.
	Sum(b []byte) ([]byte, error)
}

type hashReader struct {
	ionReader   ion.Reader
	hasher      hasher
	currentType ion.Type
	err         error
}

// NewHashReader takes an Ion reader and a hash provider and returns a new HashReader.
func NewHashReader(ionReader ion.Reader, hasherProvider IonHasherProvider) (HashReader, error) {
	newHasher, err := newHasher(hasherProvider)
	if err != nil {
		return nil, err
	}

	return &hashReader{ionReader: ionReader, hasher: *newHasher}, nil
}

// SymbolTable returns the current symbol table, or nil if there isn't one.
// Text Readers do not, generally speaking, have an associated symbol table.
// Binary Readers do.
func (hr *hashReader) SymbolTable() ion.SymbolTable {
	return hr.ionReader.SymbolTable()
}

// Next advances the Reader to the next position in the current value stream.
// It returns true if this is the position of an Ion value, and false if it
// is not. On error, it returns false and sets Err.
func (hr *hashReader) Next() bool {
	hr.err = nil

	if hr.currentType != ion.NoType {
		if ion.IsScalar(hr.currentType) || hr.IsNull() {
			hr.err = hr.hasher.scalar(hr)
			if hr.err != nil {
				return false
			}
		} else {
			hr.err = hr.StepIn()
			if hr.err != nil {
				return false
			}

			hr.err = hr.traverse()
			if hr.err != nil {
				return false
			}

			hr.err = hr.StepOut()
			if hr.err != nil {
				return false
			}
		}
	}

	next := hr.ionReader.Next()
	if !next {
		hr.err = hr.ionReader.Err()
	}

	hr.currentType = hr.ionReader.Type()

	return next
}

// Err returns an error if a previous call to Next failed.
func (hr *hashReader) Err() error {
	return hr.err
}

// Type returns the type of the Ion value the Reader is currently positioned on.
// It returns NoType if the Reader is positioned before or after a value.
func (hr *hashReader) Type() ion.Type {
	return hr.ionReader.Type()
}

// IsNull returns true if the current value is an explicit null. This may be true
// even if the Type is not NullType (for example, null.struct has type Struct).
func (hr *hashReader) IsNull() bool {
	return hr.ionReader.IsNull()
}

// FieldName returns the field name associated with the current value as a pointer. It returns
// nil if there is no current value or the current value has no field name.
func (hr *hashReader) FieldName() (*ion.SymbolToken, error) {
	return hr.ionReader.FieldName()
}

// Annotations returns the set of annotations associated with the current value.
// It returns nil if there is no current value or the current value has no annotations.
func (hr *hashReader) Annotations() ([]ion.SymbolToken, error) {
	return hr.ionReader.Annotations()
}

// StepIn steps in to the current value if it is a container. It returns an error if there
// is no current value or if the current value is not a container. On success, the Reader is
// positioned before the first value in the container.
func (hr *hashReader) StepIn() error {
	err := hr.hasher.stepIn(hr)
	if err != nil {
		return err
	}

	err = hr.ionReader.StepIn()
	if err != nil {
		return err
	}

	hr.currentType = ion.NoType

	return nil
}

// StepOut steps out of the current container value being read. It returns an error if
// this Reader is not currently positioned inside a container. On success, the Reader is
// positioned after the end of the container, but before any subsequent value(s) in the
// stream.
func (hr *hashReader) StepOut() error {
	err := hr.traverse()
	if err != nil {
		return err
	}

	err = hr.ionReader.StepOut()
	if err != nil {
		return err
	}

	err = hr.hasher.stepOut()
	if err != nil {
		return err
	}

	return nil
}

// BoolValue returns the current value as a boolean if the current value is an Ion boolean.
// It returns an error if the current value is not an Ion bool.
func (hr *hashReader) BoolValue() (*bool, error) {
	return hr.ionReader.BoolValue()
}

// IntSize returns the size of integer needed to losslessly represent the current value.
// It returns an error if the current value is not an Ion int.
func (hr *hashReader) IntSize() (ion.IntSize, error) {
	return hr.ionReader.IntSize()
}

// IntValue returns the current value as a 32-bit integer.
// It returns an error if the current value is not an Ion integer or requires more than
// 32 bits to represent the Ion integer.
func (hr *hashReader) IntValue() (*int, error) {
	return hr.ionReader.IntValue()
}

// Int64Value returns the current value as a 64-bit integer.
// It returns an error if the current value is not an Ion integer or requires more than
// 64 bits to represent the Ion integer.
func (hr *hashReader) Int64Value() (*int64, error) {
	return hr.ionReader.Int64Value()
}

// BigIntValue returns the current value as a big.Integer.
// It returns an error if the current value is not an Ion integer.
func (hr *hashReader) BigIntValue() (*big.Int, error) {
	return hr.ionReader.BigIntValue()
}

// FloatValue returns the current value as a 64-bit floating point number.
// It returns an error if the current value is not an Ion float.
func (hr *hashReader) FloatValue() (*float64, error) {
	return hr.ionReader.FloatValue()
}

// DecimalValue returns the current value as an arbitrary-precision Decimal.
// It returns an error if the current value is not an Ion decimal.
func (hr *hashReader) DecimalValue() (*ion.Decimal, error) {
	return hr.ionReader.DecimalValue()
}

// TimeValue returns the current value as a timestamp.
// It returns an error if the current value is not an Ion timestamp.
func (hr *hashReader) TimestampValue() (*ion.Timestamp, error) {
	return hr.ionReader.TimestampValue()
}

// StringValue returns the current value as a string.
// It returns an error if the current value is not an Ion symbol or an Ion string.
func (hr *hashReader) StringValue() (*string, error) {
	return hr.ionReader.StringValue()
}

// SymbolValue returns the current value as a symbol token.
// It returns an error if the current value is not an Ion symbol or an Ion string.
func (hr *hashReader) SymbolValue() (*ion.SymbolToken, error) {
	return hr.ionReader.SymbolValue()
}

// ByteValue returns the current value as a byte slice.
// It returns an error if the current value is not an Ion clob or an Ion blob.
func (hr *hashReader) ByteValue() ([]byte, error) {
	return hr.ionReader.ByteValue()
}

// Sum appends the current hash to b and returns the resulting slice.
// It resets the Hash to its initial state.
func (hr *hashReader) Sum(b []byte) ([]byte, error) {
	return hr.hasher.sum(b)
}

func (hr *hashReader) traverse() error {
	for hr.Next() {
		if ion.IsContainer(hr.currentType) && !hr.IsNull() {
			err := hr.StepIn()
			if err != nil {
				return err
			}

			err = hr.traverse()
			if err != nil {
				return err
			}

			err = hr.StepOut()
			if err != nil {
				return err
			}
		}
	}

	return hr.Err()
}

// The following implements hashValue interface.

func (hr *hashReader) getFieldName() (*ion.SymbolToken, error) {
	return hr.FieldName()
}

func (hr *hashReader) getAnnotations() ([]ion.SymbolToken, error) {
	return hr.ionReader.Annotations()
}

func (hr *hashReader) value() (interface{}, error) {
	switch hr.currentType {
	case ion.BoolType:
		return hr.BoolValue()
	case ion.BlobType:
		return hr.ionReader.ByteValue()
	case ion.ClobType:
		return hr.ionReader.ByteValue()
	case ion.DecimalType:
		return hr.DecimalValue()
	case ion.FloatType:
		return hr.FloatValue()
	case ion.IntType:
		intSize, err := hr.IntSize()
		if err != nil {
			return nil, err
		}

		switch intSize {
		case ion.Int32:
			return hr.IntValue()
		case ion.Int64:
			return hr.Int64Value()
		case ion.BigInt:
			return hr.BigIntValue()
		default:
			return nil, &InvalidOperationError{
				"hashReader", "value", "Expected intSize to be one of Int32, Int64, Uint64, or BigInt"}
		}
	case ion.StringType:
		return hr.StringValue()
	case ion.SymbolType:
		return hr.SymbolValue()
	case ion.TimestampType:
		return hr.TimestampValue()
	case ion.NoType:
		return ion.NoType, nil
	}

	return nil, &InvalidIonTypeError{hr.currentType}
}

// IsInStruct indicates if the reader is currently positioned inside a struct.
func (hr *hashReader) IsInStruct() bool {
	return hr.ionReader.IsInStruct()
}
