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
	HashValue
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
	hashWriter := &hashWriter{ionWriter, hasherProvider}

	return hashWriter
}

func (hashWriter *hashWriter) FieldName(val string) error {
	hashWriter.FieldName(val)
	panic("implement me")
}

func (hashWriter *hashWriter) Annotation(val string) error {
	hashWriter.Annotations(val)
	panic("implement me")
}

func (hashWriter *hashWriter) Annotations(vals ...string) error {
	hashWriter.Annotations(vals...)
	panic("implement me")
}

func (hashWriter *hashWriter) WriteNull() error {
	hashWriter.WriteNull()
	panic("implement me")
}

func (hashWriter *hashWriter) WriteNullType(t ion.Type) error {
	hashWriter.WriteNullType(t)
	panic("implement me")
}

func (hashWriter *hashWriter) WriteBool(val bool) error {
	hashWriter.WriteBool(val)
	panic("implement me")
}

func (hashWriter *hashWriter) WriteInt(val int64) error {
	hashWriter.WriteInt(val)
	panic("implement me")
}

func (hashWriter *hashWriter) WriteUint(val uint64) error {
	hashWriter.WriteUint(val)
	panic("implement me")
}

func (hashWriter *hashWriter) WriteBigInt(val *big.Int) error {
	hashWriter.WriteBigInt(val)
	panic("implement me")
}

func (hashWriter *hashWriter) WriteFloat(val float64) error {
	hashWriter.WriteFloat(val)
	panic("implement me")
}

func (hashWriter *hashWriter) WriteDecimal(val *ion.Decimal) error {
	hashWriter.WriteDecimal(val)
	panic("implement me")
}

func (hashWriter *hashWriter) WriteTimestamp(val time.Time) error {
	hashWriter.WriteTimestamp(val)
	panic("implement me")
}

func (hashWriter *hashWriter) WriteSymbol(val string) error {
	hashWriter.WriteSymbol(val)
	panic("implement me")
}

func (hashWriter *hashWriter) WriteString(val string) error {
	hashWriter.WriteString(val)
	panic("implement me")
}

func (hashWriter *hashWriter) WriteClob(val []byte) error {
	hashWriter.WriteClob(val)
	panic("implement me")
}

func (hashWriter *hashWriter) WriteBlob(val []byte) error {
	hashWriter.WriteBlob(val)
	panic("implement me")
}

func (hashWriter *hashWriter) BeginList() error {
	hashWriter.BeginList()
	panic("implement me")
}

func (hashWriter *hashWriter) EndList() error {
	hashWriter.EndList()
	panic("implement me")
}

func (hashWriter *hashWriter) BeginSexp() error {
	hashWriter.BeginSexp()
	panic("implement me")
}

func (hashWriter *hashWriter) EndSexp() error {
	hashWriter.EndSexp()
	panic("implement me")
}

func (hashWriter *hashWriter) BeginStruct() error {
	hashWriter.BeginStruct()
	panic("implement me")
}

func (hashWriter *hashWriter) EndStruct() error {
	hashWriter.EndStruct()
	panic("implement me")
}

func (hashWriter *hashWriter) Finish() error {
	hashWriter.Finish()
	panic("implement me")
}

func (hashWriter *hashWriter) Sum(b []byte) []byte {
	hashWriter.Sum(b)
	panic("implement me")
}

// The following implements HashValue interface.

func (hashWriter hashWriter) GetFieldName() string {
	panic("implement me")
}

func (hashWriter hashWriter) GetAnnotations() []string {
	panic("implement me")
}

func (hashWriter hashWriter) IsNull() bool {
	panic("implement me")
}

func (hashWriter hashWriter) Type() ion.Type {
	panic("implement me")
}

func (hashWriter hashWriter) Value() interface{} {
	panic("implement me")
}

func (hashWriter hashWriter) IsInStruct() bool {
	panic("implement me")
}