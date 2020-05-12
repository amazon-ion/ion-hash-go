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

type HashReader interface {
	hashValue
	// Embed interface of Ion reader.
	ion.Reader

	// Sum appends the current hash to b and returns the resulting slice.
	// It does not change the underlying hash state.
	Sum(b []byte) []byte
}

type hashReader struct {
	ionReader      ion.Reader
	hasherProvider IonHasherProvider
}

func NewHashReader(ionReader ion.Reader, hasherProvider IonHasherProvider) HashReader {
	hashReader := &hashReader{ionReader, hasherProvider}

	return hashReader
}

func (hashReader *hashReader) SymbolTable() ion.SymbolTable {
	panic("implement me")
	return hashReader.ionReader.SymbolTable()
}

func (hashReader *hashReader) Next() bool {
	panic("implement me")
	return hashReader.ionReader.Next()
}

func (hashReader *hashReader) Err() error {
	panic("implement me")
	return hashReader.ionReader.Err()
}

func (hashReader *hashReader) Type() ion.Type {
	panic("implement me")
	return hashReader.ionReader.Type()
}

func (hashReader *hashReader) IsNull() bool {
	panic("implement me")
	return hashReader.ionReader.IsNull()
}

func (hashReader *hashReader) FieldName() string {
	panic("implement me")
	return hashReader.ionReader.FieldName()
}

func (hashReader *hashReader) Annotations() []string {
	panic("implement me")
	return hashReader.ionReader.Annotations()
}

func (hashReader *hashReader) StepIn() error {
	panic("implement me")
	return hashReader.ionReader.StepIn()
}

func (hashReader *hashReader) StepOut() error {
	panic("implement me")
	return hashReader.ionReader.StepOut()
}

func (hashReader *hashReader) BoolValue() (bool, error) {
	panic("implement me")
	return hashReader.ionReader.BoolValue()
}

func (hashReader *hashReader) IntSize() (ion.IntSize, error) {
	panic("implement me")
	return hashReader.ionReader.IntSize()
}

func (hashReader *hashReader) IntValue() (int, error) {
	panic("implement me")
	return hashReader.ionReader.IntValue()
}

func (hashReader *hashReader) Int64Value() (int64, error) {
	panic("implement me")
	return hashReader.ionReader.Int64Value()
}

func (hashReader *hashReader) Uint64Value() (uint64, error) {
	panic("implement me")
	return hashReader.ionReader.Uint64Value()
}

func (hashReader *hashReader) BigIntValue() (*big.Int, error) {
	panic("implement me")
	return hashReader.ionReader.BigIntValue()
}

func (hashReader *hashReader) FloatValue() (float64, error) {
	panic("implement me")
	return hashReader.ionReader.FloatValue()
}

func (hashReader *hashReader) DecimalValue() (*ion.Decimal, error) {
	panic("implement me")
	return hashReader.ionReader.DecimalValue()
}

func (hashReader *hashReader) TimeValue() (time.Time, error) {
	return hashReader.ionReader.TimeValue()
}

func (hashReader *hashReader) StringValue() (string, error) {
	panic("implement me")
	return hashReader.ionReader.StringValue()
}

func (hashReader *hashReader) ByteValue() ([]byte, error) {
	panic("implement me")
	return hashReader.ionReader.ByteValue()
}

func (hashReader *hashReader) Sum(b []byte) []byte {
	panic("implement me")
}

// The following implements HashValue interface.

func (hashReader hashReader) getFieldName() string {
	panic("implement me")
}

func (hashReader hashReader) getAnnotations() []string {
	panic("implement me")
}

func (hashReader *hashReader) value() interface{} {
	panic("implement me")
}

func (hashReader *hashReader) isInStruct() bool {
	panic("implement me")
}

func (hashReader *hashReader) ionType() ion.Type {
	panic("implement me")
	return hashReader.Type()
}

func (hashReader *hashReader) isNull() bool {
	panic("implement me")
	return hashReader.IsNull()
}
