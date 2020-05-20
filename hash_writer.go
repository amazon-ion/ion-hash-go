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
	Sum(b []byte) []byte
}

type hashWriter struct {
	ionWriter      ion.Writer
	hasherProvider IonHasherProvider
}

func NewHashWriter(ionWriter ion.Writer, hasherProvider IonHasherProvider) HashWriter {
	hashWriter := &hashWriter{ionWriter, hasherProvider}

	return hashWriter
}

func (hashWriter *hashWriter) FieldName(val string) error {
	hashWriter.ionWriter.FieldName(val)
	panic("implement me")
}

func (hashWriter *hashWriter) Annotation(val string) error {
	hashWriter.ionWriter.Annotations(val)
	panic("implement me")
}

func (hashWriter *hashWriter) Annotations(vals ...string) error {
	hashWriter.ionWriter.Annotations(vals...)
	panic("implement me")
}

func (hashWriter *hashWriter) WriteNull() error {
	hashWriter.ionWriter.WriteNull()
	panic("implement me")
}

func (hashWriter *hashWriter) WriteNullType(t ion.Type) error {
	hashWriter.ionWriter.WriteNullType(t)
	panic("implement me")
}

func (hashWriter *hashWriter) WriteBool(val bool) error {
	hashWriter.ionWriter.WriteBool(val)
	panic("implement me")
}

func (hashWriter *hashWriter) WriteInt(val int64) error {
	hashWriter.ionWriter.WriteInt(val)
	panic("implement me")
}

func (hashWriter *hashWriter) WriteUint(val uint64) error {
	hashWriter.ionWriter.WriteUint(val)
	panic("implement me")
}

func (hashWriter *hashWriter) WriteBigInt(val *big.Int) error {
	hashWriter.ionWriter.WriteBigInt(val)
	panic("implement me")
}

func (hashWriter *hashWriter) WriteFloat(val float64) error {
	hashWriter.ionWriter.WriteFloat(val)
	panic("implement me")
}

func (hashWriter *hashWriter) WriteDecimal(val *ion.Decimal) error {
	hashWriter.ionWriter.WriteDecimal(val)
	panic("implement me")
}

func (hashWriter *hashWriter) WriteTimestamp(val time.Time) error {
	hashWriter.ionWriter.WriteTimestamp(val)
	panic("implement me")
}

func (hashWriter *hashWriter) WriteSymbol(val string) error {
	hashWriter.ionWriter.WriteSymbol(val)
	panic("implement me")
}

func (hashWriter *hashWriter) WriteString(val string) error {
	hashWriter.ionWriter.WriteString(val)
	panic("implement me")
}

func (hashWriter *hashWriter) WriteClob(val []byte) error {
	hashWriter.ionWriter.WriteClob(val)
	panic("implement me")
}

func (hashWriter *hashWriter) WriteBlob(val []byte) error {
	hashWriter.ionWriter.WriteBlob(val)
	panic("implement me")
}

func (hashWriter *hashWriter) BeginList() error {
	hashWriter.ionWriter.BeginList()
	panic("implement me")
}

func (hashWriter *hashWriter) EndList() error {
	hashWriter.ionWriter.EndList()
	panic("implement me")
}

func (hashWriter *hashWriter) BeginSexp() error {
	hashWriter.ionWriter.BeginSexp()
	panic("implement me")
}

func (hashWriter *hashWriter) EndSexp() error {
	hashWriter.ionWriter.EndSexp()
	panic("implement me")
}

func (hashWriter *hashWriter) BeginStruct() error {
	hashWriter.ionWriter.BeginStruct()
	panic("implement me")
}

func (hashWriter *hashWriter) EndStruct() error {
	hashWriter.ionWriter.EndStruct()
	panic("implement me")
}

func (hashWriter *hashWriter) Finish() error {
	hashWriter.ionWriter.Finish()
	panic("implement me")
}

func (hashWriter *hashWriter) Sum(b []byte) []byte {
	panic("implement me")
}

// The following implements HashValue interface.

func (hashWriter hashWriter) getFieldName() string {
	panic("implement me")
}

func (hashWriter hashWriter) getAnnotations() []string {
	panic("implement me")
}

func (hashWriter hashWriter) isNull() bool {
	panic("implement me")
}

func (hashWriter hashWriter) ionType() ion.Type {
	panic("implement me")
}

func (hashWriter hashWriter) value() interface{} {
	panic("implement me")
}

func (hashWriter hashWriter) isInStruct() bool {
	panic("implement me")
}
