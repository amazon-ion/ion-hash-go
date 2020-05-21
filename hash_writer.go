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

	"ion-go"
)

type HashWriter interface {
	hashValue
	// Embed interface of Ion writer.
	ion.Writer

	// Sum appends the current hash to b and returns the resulting slice.
	// It does not change the underlying hash state.
	Sum(b []byte) []byte
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

func NewHashWriter(ionWriter ion.Writer, hasherProvider IonHasherProvider) HashWriter {
	hasher := newHasher(hasherProvider)
	hashWriter := &hashWriter{ionWriter: ionWriter, hasher: *hasher}

	return hashWriter
}

func (hashWriter *hashWriter) FieldName(val string) error {
	hashWriter.currentFieldName = val
	return hashWriter.ionWriter.FieldName(val)
}

func (hashWriter *hashWriter) Annotation(val string) error {
	hashWriter.annotations = append(hashWriter.annotations, val)
	return hashWriter.ionWriter.Annotations(val)
}

func (hashWriter *hashWriter) Annotations(vals ...string) error {
	hashWriter.annotations = append(hashWriter.annotations, vals...)
	return hashWriter.ionWriter.Annotations(vals...)
}

func (hashWriter *hashWriter) WriteNull() error {
	error := hashWriter.hashScalar(ion.NullType, nil)
	if error != nil {
		return error
	}
	return hashWriter.ionWriter.WriteNull()
}

func (hashWriter *hashWriter) WriteNullType(ionType ion.Type) error {
	error := hashWriter.hashScalar(ionType, nil)
	if error != nil {
		return error
	}
	return hashWriter.ionWriter.WriteNullType(ionType)
}

func (hashWriter *hashWriter) WriteBool(val bool) error {
	error := hashWriter.hashScalar(ion.BoolType, val)
	if error != nil {
		return error
	}
	return hashWriter.ionWriter.WriteBool(val)
}

func (hashWriter *hashWriter) WriteInt(val int64) error {
	error := hashWriter.hashScalar(ion.IntType, val)
	if error != nil {
		return error
	}
	return hashWriter.ionWriter.WriteInt(val)
}

func (hashWriter *hashWriter) WriteUint(val uint64) error {
	error := hashWriter.hashScalar(ion.IntType, val)
	if error != nil {
		return error
	}
	return hashWriter.ionWriter.WriteUint(val)
}

func (hashWriter *hashWriter) WriteBigInt(val *big.Int) error {
	error := hashWriter.hashScalar(ion.IntType, val)
	if error != nil {
		return error
	}
	return hashWriter.ionWriter.WriteBigInt(val)
}

func (hashWriter *hashWriter) WriteFloat(val float64) error {
	error := hashWriter.hashScalar(ion.FloatType, val)
	if error != nil {
		return error
	}
	return hashWriter.ionWriter.WriteFloat(val)
}

func (hashWriter *hashWriter) WriteDecimal(val *ion.Decimal) error {
	error := hashWriter.hashScalar(ion.DecimalType, val)
	if error != nil {
		return error
	}
	return hashWriter.ionWriter.WriteDecimal(val)
}

func (hashWriter *hashWriter) WriteTimestamp(val time.Time) error {
	error := hashWriter.hashScalar(ion.TimestampType, val)
	if error != nil {
		return error
	}
	return hashWriter.ionWriter.WriteTimestamp(val)
}

func (hashWriter *hashWriter) WriteSymbol(val string) error {
	error := hashWriter.hashScalar(ion.SymbolType, val)
	if error != nil {
		return error
	}
	return hashWriter.ionWriter.WriteSymbol(val)
}

func (hashWriter *hashWriter) WriteString(val string) error {
	error := hashWriter.hashScalar(ion.StringType, val)
	if error != nil {
		return error
	}
	return hashWriter.ionWriter.WriteString(val)
}

func (hashWriter *hashWriter) WriteClob(val []byte) error {
	error := hashWriter.hashScalar(ion.ClobType, val)
	if error != nil {
		return error
	}
	return hashWriter.ionWriter.WriteClob(val)
}

func (hashWriter *hashWriter) WriteBlob(val []byte) error {
	error := hashWriter.hashScalar(ion.BlobType, val)
	if error != nil {
		return error
	}
	return hashWriter.ionWriter.WriteBlob(val)
}

func (hashWriter *hashWriter) BeginList() error {
	error := hashWriter.stepIn(ion.ListType)
	if error != nil {
		return error
	}

	hashWriter.ionWriter.BeginList()
	panic("implement me")
}

func (hashWriter *hashWriter) EndList() error {
	error := hashWriter.hasher.stepOut()
	if error != nil {
		return error
	}

	return hashWriter.ionWriter.EndList()
}

func (hashWriter *hashWriter) BeginSexp() error {
	error := hashWriter.stepIn(ion.SexpType)
	if error != nil {
		return error
	}

	hashWriter.ionWriter.BeginSexp()
	panic("implement me")
}

func (hashWriter *hashWriter) EndSexp() error {
	error := hashWriter.hasher.stepOut()
	if error != nil {
		return error
	}

	return hashWriter.ionWriter.EndSexp()
}

func (hashWriter *hashWriter) BeginStruct() error {
	error := hashWriter.stepIn(ion.StructType)
	if error != nil {
		return error
	}

	hashWriter.ionWriter.BeginStruct()
	panic("implement me")
}

func (hashWriter *hashWriter) EndStruct() error {
	error := hashWriter.hasher.stepOut()
	if error != nil {
		return error
	}

	return hashWriter.ionWriter.EndStruct()
}

func (hashWriter *hashWriter) Finish() error {
	return hashWriter.ionWriter.Finish()
}

func (hashWriter *hashWriter) Sum(b []byte) []byte {
	return hashWriter.hasher.sum()
}

// The following implements HashValue interface.

func (hashWriter *hashWriter) getFieldName() string {
	return hashWriter.currentFieldName
}

func (hashWriter *hashWriter) getAnnotations() []string {
	return hashWriter.annotations
}

func (hashWriter *hashWriter) isNull() bool {
	return hashWriter.currentIsNull
}

func (hashWriter *hashWriter) ionType() ion.Type {
	return hashWriter.currentType
}

func (hashWriter *hashWriter) value() interface{} {
	return hashWriter.currentValue
}

func (hashWriter *hashWriter) isInStruct() bool {
	return hashWriter.currentType == ion.StructType
}

func (hashWriter *hashWriter) hashScalar(ionType ion.Type, value interface{}) error {
	hashWriter.currentType = ionType
	hashWriter.currentValue = value
	hashWriter.currentIsNull = value == nil
	hashWriter.currentFieldName = ""
	hashWriter.annotations = nil

	return hashWriter.hasher.scalar(hashWriter)
}

func (hashWriter *hashWriter) stepIn(ionType ion.Type) error {
	hashWriter.currentType = ionType
	hashWriter.currentValue = nil
	hashWriter.currentIsNull = false
	hashWriter.currentFieldName = ""
	hashWriter.annotations = nil

	return hashWriter.hasher.stepIn(hashWriter)
}
