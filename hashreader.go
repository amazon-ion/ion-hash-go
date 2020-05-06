/*
 * Copyright 2020 Amazon.com, Inc. or its affiliates. All Rights Reserved.
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
	"ion-go"
	"math/big"
	"time"
)

type HashReader interface {
	// Extend interface of Ion reader.
	ion.Reader

	// Sum appends the current hash to b and returns the resulting slice.
	// It does not change the underlying hash state.
	Sum(b []byte) []byte
}

type hashReader struct {
	ionReader ion.Reader
	hasherProvider IonHashProvider
}

func NewHashReader(ionReader ion.Reader, hasherProvider IonHashProvider) HashReader {
	hr := &hashReader{ionReader, hasherProvider}

	return hr
}

// SymbolTable returns the current symbol table, or nil if there isn't one.
// Text Readers do not, generally speaking, have an associated symbol table.
// Binary Readers do.
func (hr *hashReader) SymbolTable() ion.SymbolTable {
	panic("implement me")
	return hr.SymbolTable()
}

// Next advances the Reader to the next position in the current value stream.
// It returns true if this is the position of an Ion value, and false if it
// is not. On error, it returns false and sets Err.
func (hr *hashReader) Next() bool {
	panic("implement me")
	return hr.Next()
}

// Err returns an error if a previous call call to Next has failed.
func (hr *hashReader) Err() error {
	panic("implement me")
	return hr.Err()
}

// Type returns the type of the Ion value the Reader is currently positioned on.
// It returns NoType if the Reader is positioned before or after a value.
func (hr *hashReader) Type() ion.Type {
	panic("implement me")
	return hr.Type()
}

// IsNull returns true if the current value is an explicit null. This may be true
// even if the Type is not NullType (for example, null.struct has type Struct). Yes,
// that's a bit confusing.
func (hr *hashReader) IsNull() bool {
	panic("implement me")
	return hr.IsNull()
}

// FieldName returns the field name associated with the current value. It returns
// the empty string if there is no current value or the current value has no field
// name.
func (hr *hashReader) FieldName() string {
	panic("implement me")
	return hr.FieldName()
}

// Annotations returns the set of annotations associated with the current value.
// It returns nil if there is no current value or the current value has no annotations.
func (hr *hashReader) Annotations() []string {
	panic("implement me")
	return hr.Annotations()
}

// StepIn steps in to the current value if it is a container. It returns an error if there
// is no current value or if the value is not a container. On success, the Reader is
// positioned before the first value in the container.
func (hr *hashReader) StepIn() error {
	panic("implement me")
	return hr.StepIn()
}

// StepOut steps out of the current container value being read. It returns an error if
// this Reader is not currently stepped in to a container. On success, the Reader is
// positioned after the end of the container, but before any subsequent values in the
// stream.
func (hr *hashReader) StepOut() error {
	panic("implement me")
	return hr.StepOut()
}

// BoolValue returns the current value as a boolean (if that makes sense). It returns
// an error if the current value is not an Ion bool.
func (hr *hashReader) BoolValue() (bool, error) {
	panic("implement me")
	return hr.BoolValue()
}

// IntSize returns the size of integer needed to losslessly represent the current value
// (if that makes sense). It returns an error if the current value is not an Ion int.
func (hr *hashReader) IntSize() (ion.IntSize, error) {
	panic("implement me")
	return hr.IntSize()
}

// IntValue returns the current value as a 32-bit integer (if that makes sense). It
// returns an error if the current value is not an Ion integer or requires more than
// 32 bits to represent losslessly.
func (hr *hashReader) IntValue() (int, error) {
	panic("implement me")
	return hr.IntValue()
}

// Int64Value returns the current value as a 64-bit integer (if that makes sense). It
// returns an error if the current value is not an Ion integer or requires more than
// 64 bits to represent losslessly.
func (hr *hashReader) Int64Value() (int64, error) {
	panic("implement me")
	return hr.Int64Value()
}

// Uint64Value returns the current value as an unsigned 64-bit integer (if that makes
// sense). It returns an error if the current value is not an Ion integer, is negative,
// or requires more than 64 bits to represent losslessly.
func (hr *hashReader) Uint64Value() (uint64, error) {
	panic("implement me")
	return hr.Uint64Value()
}

// BigIntValue returns the current value as a big.Integer (if that makes sense). It
// returns an error if the current value is not an Ion integer.
func (hr *hashReader) BigIntValue() (*big.Int, error) {
	panic("implement me")
	return hr.BigIntValue()
}

// FloatValue returns the current value as a 64-bit floating point number (if that
// makes sense). It returns an error if the current value is not an Ion float.
func (hr *hashReader) FloatValue() (float64, error) {
	panic("implement me")
	return hr.FloatValue()
}

// DecimalValue returns the current value as an arbitrary-precision Decimal (if that
// makes sense). It returns an error if the current value is not an Ion decimal.
func (hr *hashReader) DecimalValue() (*ion.Decimal, error) {
	panic("implement me")
	return hr.DecimalValue()
}

// TimeValue returns the current value as a timestamp (if that makes sense). It returns
// an error if the current value is not an Ion timestamp.
func (hr *hashReader) TimeValue() (time.Time, error) {
	return hr.TimeValue()
}

// StringValue returns the current value as a string (if that makes sense). It returns
// an error if the current value is not an Ion symbol or an Ion string.
func (hr *hashReader) StringValue() (string, error) {
	panic("implement me")
	return hr.StringValue()
}

// ByteValue returns the current value as a byte slice (if that makes sense). It returns
// an error if the current value is not an Ion clob or an Ion blob.
func (hr *hashReader) ByteValue() ([]byte, error) {
	panic("implement me")
	return hr.ByteValue()
}

// Sum appends the current hash to b and returns the resulting slice.
// It does not change the underlying hash state.
func (hr *hashReader) Sum(b []byte) []byte {
	panic("implement me")
}
