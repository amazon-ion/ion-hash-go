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

type HashWriter interface {
	hashValue
	// Embed interface of Ion writer.
	ion.Writer

	// Sum appends the current hash to b and returns the resulting slice.
	// It does not change the underlying hash state.
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

func NewHashWriter(ionWriter ion.Writer, hasherProvider IonHasherProvider) (HashWriter, error) {
	newHasher, err := newHasher(hasherProvider)
	if err != nil {
		return nil, err
	}

	return &hashWriter{ionWriter: ionWriter, hasher: *newHasher}, nil
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
	err := hashWriter.hashScalar(ion.NullType, nil)
	if err != nil {
		return err
	}
	return hashWriter.ionWriter.WriteNull()
}

func (hashWriter *hashWriter) WriteNullType(ionType ion.Type) error {
	err := hashWriter.hashScalar(ionType, nil)
	if err != nil {
		return err
	}
	return hashWriter.ionWriter.WriteNullType(ionType)
}

func (hashWriter *hashWriter) WriteBool(val bool) error {
	err := hashWriter.hashScalar(ion.BoolType, val)
	if err != nil {
		return err
	}
	return hashWriter.ionWriter.WriteBool(val)
}

func (hashWriter *hashWriter) WriteInt(val int64) error {
	err := hashWriter.hashScalar(ion.IntType, val)
	if err != nil {
		return err
	}
	return hashWriter.ionWriter.WriteInt(val)
}

func (hashWriter *hashWriter) WriteUint(val uint64) error {
	err := hashWriter.hashScalar(ion.IntType, val)
	if err != nil {
		return err
	}
	return hashWriter.ionWriter.WriteUint(val)
}

func (hashWriter *hashWriter) WriteBigInt(val *big.Int) error {
	err := hashWriter.hashScalar(ion.IntType, val)
	if err != nil {
		return err
	}
	return hashWriter.ionWriter.WriteBigInt(val)
}

func (hashWriter *hashWriter) WriteFloat(val float64) error {
	err := hashWriter.hashScalar(ion.FloatType, val)
	if err != nil {
		return err
	}
	return hashWriter.ionWriter.WriteFloat(val)
}

func (hashWriter *hashWriter) WriteDecimal(val *ion.Decimal) error {
	err := hashWriter.hashScalar(ion.DecimalType, val)
	if err != nil {
		return err
	}
	return hashWriter.ionWriter.WriteDecimal(val)
}

func (hashWriter *hashWriter) WriteTimestamp(val time.Time) error {
	err := hashWriter.hashScalar(ion.TimestampType, val)
	if err != nil {
		return err
	}
	return hashWriter.ionWriter.WriteTimestamp(val)
}

func (hashWriter *hashWriter) WriteSymbol(val string) error {
	err := hashWriter.hashScalar(ion.SymbolType, val)
	if err != nil {
		return err
	}
	return hashWriter.ionWriter.WriteSymbol(val)
}

func (hashWriter *hashWriter) WriteString(val string) error {
	err := hashWriter.hashScalar(ion.StringType, val)
	if err != nil {
		return err
	}
	return hashWriter.ionWriter.WriteString(val)
}

func (hashWriter *hashWriter) WriteClob(val []byte) error {
	err := hashWriter.hashScalar(ion.ClobType, val)
	if err != nil {
		return err
	}
	return hashWriter.ionWriter.WriteClob(val)
}

func (hashWriter *hashWriter) WriteBlob(val []byte) error {
	err := hashWriter.hashScalar(ion.BlobType, val)
	if err != nil {
		return err
	}
	return hashWriter.ionWriter.WriteBlob(val)
}

func (hashWriter *hashWriter) BeginList() error {
	err := hashWriter.stepIn(ion.ListType)
	if err != nil {
		return err
	}

	return hashWriter.ionWriter.BeginList()
}

func (hashWriter *hashWriter) EndList() error {
	err := hashWriter.hasher.stepOut()
	if err != nil {
		return err
	}

	return hashWriter.ionWriter.EndList()
}

func (hashWriter *hashWriter) BeginSexp() error {
	err := hashWriter.stepIn(ion.SexpType)
	if err != nil {
		return err
	}

	return hashWriter.ionWriter.BeginSexp()
}

func (hashWriter *hashWriter) EndSexp() error {
	err := hashWriter.hasher.stepOut()
	if err != nil {
		return err
	}

	return hashWriter.ionWriter.EndSexp()
}

func (hashWriter *hashWriter) BeginStruct() error {
	err := hashWriter.stepIn(ion.StructType)
	if err != nil {
		return err
	}

	return hashWriter.ionWriter.BeginStruct()
}

func (hashWriter *hashWriter) EndStruct() error {
	err := hashWriter.hasher.stepOut()
	if err != nil {
		return err
	}

	return hashWriter.ionWriter.EndStruct()
}

func (hashWriter *hashWriter) Finish() error {
	return hashWriter.ionWriter.Finish()
}

func (hashWriter *hashWriter) Sum(b []byte) ([]byte, error) {
	return hashWriter.hasher.sum(b)
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

func (hashWriter *hashWriter) value() (interface{}, error) {
	return hashWriter.currentValue, nil
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
