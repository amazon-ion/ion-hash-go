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
	"time"

	"github.com/amzn/ion-go/ion"
)

// A HashWriter writes a stream of Ion values and calculates its hash.
//
// The various Write methods write atomic values to the current output stream. Methods
// prefixed with Begin start writing a list, sexp, or struct respectively. Subsequent
// calls to Write will write values inside of the container until a matching
// End method is called, e.g.,
//
// 	   var hw HashWriter
// 	   hw.BeginSexp()
// 	   {
// 		   hw.WriteInt(1)
// 		   hw.WriteSymbol("+")
// 		   hw.WriteInt(1)
// 	   }
// 	   hw.EndSexp()
//
// When writing values inside a struct, the FieldName method must be called before
// each value to set the value's field name. The Annotation method may likewise
// be called before writing any value to add an annotation to the value.
//
// 	   var hw HashWriter
// 	   hw.Annotation("user")
// 	   hw.BeginStruct()
// 	   {
// 		   hw.FieldName("id")
// 		   hw.WriteString("foo")
// 		   hw.FieldName("name")
// 		   hw.WriteString("bar")
// 	   }
// 	   hw.EndStruct()
//
// When you're done writing values, you should call Finish to ensure everything has
// been flushed from in-memory buffers. While individual methods all return an error
// on failure, implementations will remember any errors, no-op subsequent calls, and
// return the previous error. This lets you keep code a bit cleaner by only checking
// the return value of the final method call (generally Finish).
//
// Sum will return the hash of the entire stream of Ion values that have been written thus far.
//
// 	   var hw HashWriter
// 	   writeSomeStuff(hw)
// 	   if err := hw.Finish(); err != nil {
// 		   return err
// 	   }
//     fmt.Printf("%v", hw.Sum(nil))
//
type HashWriter interface {
	// Embed interface of Ion writer.
	ion.Writer

	// Remaining hashValue methods.
	getFieldName() *string
	getAnnotations() []string
	IsNull() bool
	Type() ion.Type
	value() (interface{}, error)

	// Sum appends the current hash to b and returns the resulting slice.
	// It resets the Hash to its initial state.
	Sum(b []byte) ([]byte, error)
}

type hashWriter struct {
	ionWriter ion.Writer
	hasher    hasher

	currentFieldName string
	currentType      ion.Type
	currentValue     interface{}
	currentIsNull    bool
	annotations      []string
}

// NewHashWriter takes an Ion Writer and a hash provider and returns a new HashWriter.
func NewHashWriter(ionWriter ion.Writer, hasherProvider IonHasherProvider) (HashWriter, error) {
	newHasher, err := newHasher(hasherProvider)
	if err != nil {
		return nil, err
	}

	return &hashWriter{ionWriter: ionWriter, hasher: *newHasher}, nil
}

// FieldName sets the field name for the next value written.
// It may only be called while writing a struct.
func (hw *hashWriter) FieldName(val string) error {
	hw.currentFieldName = val
	return hw.ionWriter.FieldName(val)
}

// Annotation adds an annotation to the next value written.
func (hw *hashWriter) Annotation(val string) error {
	hw.annotations = append(hw.annotations, val)
	return hw.ionWriter.Annotations(val)
}

// Annotations adds one or more annotations to the next value written.
func (hw *hashWriter) Annotations(vals ...string) error {
	hw.annotations = append(hw.annotations, vals...)
	return hw.ionWriter.Annotations(vals...)
}

// WriteNull writes an untyped null value.
func (hw *hashWriter) WriteNull() error {
	err := hw.hashScalar(ion.NullType, nil)
	if err != nil {
		return err
	}
	return hw.ionWriter.WriteNull()
}

// WriteNullType writes a null value with a type qualifier, e.g. null.bool.
func (hw *hashWriter) WriteNullType(ionType ion.Type) error {
	err := hw.hashScalar(ionType, nil)
	if err != nil {
		return err
	}
	return hw.ionWriter.WriteNullType(ionType)
}

// WriteBool writes a boolean value.
func (hw *hashWriter) WriteBool(val bool) error {
	err := hw.hashScalar(ion.BoolType, val)
	if err != nil {
		return err
	}
	return hw.ionWriter.WriteBool(val)
}

// WriteInt writes an integer value.
func (hw *hashWriter) WriteInt(val int64) error {
	err := hw.hashScalar(ion.IntType, val)
	if err != nil {
		return err
	}
	return hw.ionWriter.WriteInt(val)
}

// WriteUint writes an unsigned integer value.
func (hw *hashWriter) WriteUint(val uint64) error {
	err := hw.hashScalar(ion.IntType, val)
	if err != nil {
		return err
	}
	return hw.ionWriter.WriteUint(val)
}

// WriteBigInt writes a big integer value.
func (hw *hashWriter) WriteBigInt(val *big.Int) error {
	err := hw.hashScalar(ion.IntType, val)
	if err != nil {
		return err
	}
	return hw.ionWriter.WriteBigInt(val)
}

// WriteFloat writes a floating-point value.
func (hw *hashWriter) WriteFloat(val float64) error {
	err := hw.hashScalar(ion.FloatType, val)
	if err != nil {
		return err
	}
	return hw.ionWriter.WriteFloat(val)
}

// WriteDecimal writes an arbitrary-precision decimal value.
func (hw *hashWriter) WriteDecimal(val *ion.Decimal) error {
	err := hw.hashScalar(ion.DecimalType, val)
	if err != nil {
		return err
	}
	return hw.ionWriter.WriteDecimal(val)
}

// WriteTimestamp writes a timestamp value.
func (hw *hashWriter) WriteTimestamp(val time.Time) error {
	err := hw.hashScalar(ion.TimestampType, val)
	if err != nil {
		return err
	}
	return hw.ionWriter.WriteTimestamp(val)
}

// WriteSymbol writes a symbol value.
func (hw *hashWriter) WriteSymbol(val string) error {
	err := hw.hashScalar(ion.SymbolType, val)
	if err != nil {
		return err
	}
	return hw.ionWriter.WriteSymbol(val)
}

// WriteString writes a string value.
func (hw *hashWriter) WriteString(val string) error {
	err := hw.hashScalar(ion.StringType, val)
	if err != nil {
		return err
	}
	return hw.ionWriter.WriteString(val)
}

// WriteClob writes a clob value.
func (hw *hashWriter) WriteClob(val []byte) error {
	err := hw.hashScalar(ion.ClobType, val)
	if err != nil {
		return err
	}
	return hw.ionWriter.WriteClob(val)
}

// WriteBlob writes a blob value.
func (hw *hashWriter) WriteBlob(val []byte) error {
	err := hw.hashScalar(ion.BlobType, val)
	if err != nil {
		return err
	}
	return hw.ionWriter.WriteBlob(val)
}

// BeginList begins writing a list value.
func (hw *hashWriter) BeginList() error {
	err := hw.stepIn(ion.ListType)
	if err != nil {
		return err
	}

	return hw.ionWriter.BeginList()
}

// EndList finishes writing a list value.
func (hw *hashWriter) EndList() error {
	err := hw.hasher.stepOut()
	if err != nil {
		return err
	}

	return hw.ionWriter.EndList()
}

// BeginSexp begins writing an s-expression value.
func (hw *hashWriter) BeginSexp() error {
	err := hw.stepIn(ion.SexpType)
	if err != nil {
		return err
	}

	return hw.ionWriter.BeginSexp()
}

// EndSexp finishes writing an s-expression value.
func (hw *hashWriter) EndSexp() error {
	err := hw.hasher.stepOut()
	if err != nil {
		return err
	}

	return hw.ionWriter.EndSexp()
}

// BeginStruct begins writing a struct value.
func (hw *hashWriter) BeginStruct() error {
	err := hw.stepIn(ion.StructType)
	if err != nil {
		return err
	}

	return hw.ionWriter.BeginStruct()
}

// EndStruct finishes writing a struct value.
func (hw *hashWriter) EndStruct() error {
	err := hw.hasher.stepOut()
	if err != nil {
		return err
	}

	return hw.ionWriter.EndStruct()
}

// Finish finishes writing values and flushes any buffered data.
func (hw *hashWriter) Finish() error {
	return hw.ionWriter.Finish()
}

// Sum appends the current hash to b and returns the resulting slice.
// It resets the Hash to its initial state.
func (hw *hashWriter) Sum(b []byte) ([]byte, error) {
	return hw.hasher.sum(b)
}

// The following implements hashValue interface.

func (hw *hashWriter) getFieldName() *string {
	return &hw.currentFieldName
}

func (hw *hashWriter) getAnnotations() []string {
	return hw.annotations
}

// IsNull returns true if the current value is an explicit null. This may be true
// even if the Type is not NullType (for example, null.struct has type Struct).
func (hw *hashWriter) IsNull() bool {
	return hw.currentIsNull
}

// Type returns the type of the Ion value the hashWriter is currently positioned on.
// It returns NoType if the hashWriter is positioned before or after a value.
func (hw *hashWriter) Type() ion.Type {
	return hw.currentType
}

func (hw *hashWriter) value() (interface{}, error) {
	return hw.currentValue, nil
}

// IsInStruct indicates if the writer is currently positioned inside a struct.
func (hw *hashWriter) IsInStruct() bool {
	return hw.ionWriter.IsInStruct()
}

func (hw *hashWriter) hashScalar(ionType ion.Type, value interface{}) error {
	hw.currentType = ionType
	hw.currentValue = value
	hw.currentIsNull = value == nil

	err := hw.hasher.scalar(hw)
	if err != nil {
		return err
	}

	hw.currentFieldName = ""
	hw.annotations = nil

	return nil
}

func (hw *hashWriter) stepIn(ionType ion.Type) error {
	hw.currentType = ionType
	hw.currentValue = nil
	hw.currentIsNull = false

	err := hw.hasher.stepIn(hw)
	if err != nil {
		return err
	}

	hw.currentFieldName = ""
	hw.annotations = nil

	return nil
}
