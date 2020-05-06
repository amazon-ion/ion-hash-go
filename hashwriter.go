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

type HashWriter interface {
	// Extend interface of Ion writer.
	ion.Writer

	// Sum appends the current hash to b and returns the resulting slice.
	// It does not change the underlying hash state.
	Sum(b []byte) []byte
}

type hashWriter struct {
	ionWriter ion.Writer
	hasherProvider IonHashProvider
}

func NewHashWriter(ionWriter ion.Writer, hasherProvider IonHashProvider) HashWriter {
	hw := &hashWriter{ionWriter, hasherProvider}

	return hw
}

// FieldName sets the field name for the next value written.
func (hw *hashWriter) FieldName(val string) error {
	hw.FieldName(val)
	panic("implement me")
}

// Annotation adds a single annotation to the next value written.
func (hw *hashWriter) Annotation(val string) error {
	hw.Annotations(val)
	panic("implement me")
}

// Annotations adds multiple annotations to the next value written.
func (hw *hashWriter) Annotations(vals ...string) error {
	hw.Annotations(vals...)
	panic("implement me")
}

// WriteNull writes an untyped null value.
func (hw *hashWriter) WriteNull() error {
	hw.WriteNull()
	panic("implement me")
}

// WriteNullType writes a null value with a type qualifier, e.g. null.bool.
func (hw *hashWriter) WriteNullType(t ion.Type) error {
	hw.WriteNullType(t)
	panic("implement me")
}

// WriteBool writes a boolean value.
func (hw *hashWriter) WriteBool(val bool) error {
	hw.WriteBool(val)
	panic("implement me")
}

// WriteInt writes an integer value.
func (hw *hashWriter) WriteInt(val int64) error {
	hw.WriteInt(val)
	panic("implement me")
}

// WriteUint writes an unsigned integer value.
func (hw *hashWriter) WriteUint(val uint64) error {
	hw.WriteUint(val)
	panic("implement me")
}

// WriteBigInt writes a big integer value.
func (hw *hashWriter) WriteBigInt(val *big.Int) error {
	hw.WriteBigInt(val)
	panic("implement me")
}

// WriteFloat writes a floating-point value.
func (hw *hashWriter) WriteFloat(val float64) error {
	hw.WriteFloat(val)
	panic("implement me")
}

// WriteDecimal writes an arbitrary-precision decimal value.
func (hw *hashWriter) WriteDecimal(val *ion.Decimal) error {
	hw.WriteDecimal(val)
	panic("implement me")
}

// WriteTimestamp writes a timestamp value.
func (hw *hashWriter) WriteTimestamp(val time.Time) error {
	hw.WriteTimestamp(val)
	panic("implement me")
}

// WriteSymbol writes a symbol value.
func (hw *hashWriter) WriteSymbol(val string) error {
	hw.WriteSymbol(val)
	panic("implement me")
}

// WriteString writes a string value.
func (hw *hashWriter) WriteString(val string) error {
	hw.WriteString(val)
	panic("implement me")
}

// WriteClob writes a clob value.
func (hw *hashWriter) WriteClob(val []byte) error {
	hw.WriteClob(val)
	panic("implement me")
}

// WriteBlob writes a blob value.
func (hw *hashWriter) WriteBlob(val []byte) error {
	hw.WriteBlob(val)
	panic("implement me")
}

// BeginList begins writing a list value.
func (hw *hashWriter) BeginList() error {
	hw.BeginList()
	panic("implement me")
}

// EndList finishes writing a list value.
func (hw *hashWriter) EndList() error {
	hw.EndList()
	panic("implement me")
}

// BeginSexp begins writing an s-expression value.
func (hw *hashWriter) BeginSexp() error {
	hw.BeginSexp()
	panic("implement me")
}

// EndSexp finishes writing an s-expression value.
func (hw *hashWriter) EndSexp() error {
	hw.EndSexp()
	panic("implement me")
}

// BeginStruct begins writing a struct value.
func (hw *hashWriter) BeginStruct() error {
	hw.BeginStruct()
	panic("implement me")
}

// EndStruct finishes writing a struct value.
func (hw *hashWriter) EndStruct() error {
	hw.EndStruct()
	panic("implement me")
}

// Finish finishes writing values and flushes any buffered data.
func (hw *hashWriter) Finish() error {
	hw.Finish()
	panic("implement me")
}

// Sum appends the current hash to b and returns the resulting slice.
// It does not change the underlying hash state.
func (hw *hashWriter) Sum(b []byte) []byte {
	hw.Sum(b)
	panic("implement me")
}
