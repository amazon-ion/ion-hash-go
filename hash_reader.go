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

// HashReader inherits functions from Ion's Reader interface and adds
// the function Sum that allows read access to the hash value held by this reader.
type HashReader interface {
	hashValue
	// Embed interface of Ion reader.
	ion.Reader

	// Sum appends the current hash to b and returns the resulting slice.
	// It does not change the underlying hash state.
	Sum(b []byte) ([]byte, error)
}

type hashReader struct {
	ionReader   ion.Reader
	hasher      hasher
	currentType ion.Type
	err         error
}

// NewHashReader takes an Ion reader and a hash provider and returns a new HashReader
func NewHashReader(ionReader ion.Reader, hasherProvider IonHasherProvider) (HashReader, error) {
	newHasher, err := newHasher(hasherProvider)
	if err != nil {
		return nil, err
	}

	return &hashReader{ionReader: ionReader, hasher: *newHasher}, nil
}

func (hashReader *hashReader) SymbolTable() ion.SymbolTable {
	return hashReader.ionReader.SymbolTable()
}

func (hashReader *hashReader) Next() bool {
	hashReader.err = nil

	if hashReader.currentType != ion.NoType {
		if ion.IsScalar(hashReader.currentType) || hashReader.IsNull() {
			err := hashReader.hasher.scalar(hashReader)
			if err != nil {
				hashReader.err = err
				return false
			}
		} else {
			err := hashReader.StepIn()
			if err != nil {
				hashReader.err = err
				return false
			}

			err = hashReader.traverse()
			if err != nil {
				hashReader.err = err
				return false
			}

			err = hashReader.StepOut()
			if err != nil {
				hashReader.err = err
				return false
			}
		}
	}

	next := hashReader.ionReader.Next()
	if !next {
		hashReader.err = hashReader.ionReader.Err()
	}

	hashReader.currentType = hashReader.ionReader.Type()

	return next
}

func (hashReader *hashReader) Err() error {
	return hashReader.err
}

func (hashReader *hashReader) Type() ion.Type {
	return hashReader.ionReader.Type()
}

func (hashReader *hashReader) IsNull() bool {
	return hashReader.ionReader.IsNull()
}

func (hashReader *hashReader) FieldName() string {
	return hashReader.ionReader.FieldName()
}

func (hashReader *hashReader) Annotations() []string {
	return hashReader.ionReader.Annotations()
}

func (hashReader *hashReader) StepIn() error {
	err := hashReader.hasher.stepIn(hashReader)
	if err != nil {
		return err
	}

	err = hashReader.ionReader.StepIn()
	if err != nil {
		return err
	}

	hashReader.currentType = ion.NoType

	return nil
}

func (hashReader *hashReader) StepOut() error {
	err := hashReader.traverse()
	if err != nil {
		return err
	}

	err = hashReader.ionReader.StepOut()
	if err != nil {
		return err
	}

	err = hashReader.hasher.stepOut()
	if err != nil {
		return err
	}

	return nil
}

func (hashReader *hashReader) BoolValue() (bool, error) {
	return hashReader.ionReader.BoolValue()
}

func (hashReader *hashReader) IntSize() (ion.IntSize, error) {
	return hashReader.ionReader.IntSize()
}

func (hashReader *hashReader) IntValue() (int, error) {
	return hashReader.ionReader.IntValue()
}

func (hashReader *hashReader) Int64Value() (int64, error) {
	return hashReader.ionReader.Int64Value()
}

func (hashReader *hashReader) Uint64Value() (uint64, error) {
	return hashReader.ionReader.Uint64Value()
}

func (hashReader *hashReader) BigIntValue() (*big.Int, error) {
	return hashReader.ionReader.BigIntValue()
}

func (hashReader *hashReader) FloatValue() (float64, error) {
	return hashReader.ionReader.FloatValue()
}

func (hashReader *hashReader) DecimalValue() (*ion.Decimal, error) {
	return hashReader.ionReader.DecimalValue()
}

func (hashReader *hashReader) TimeValue() (time.Time, error) {
	return hashReader.ionReader.TimeValue()
}

func (hashReader *hashReader) StringValue() (string, error) {
	return hashReader.ionReader.StringValue()
}

func (hashReader *hashReader) ByteValue() ([]byte, error) {
	return hashReader.ionReader.ByteValue()
}

func (hashReader *hashReader) Sum(b []byte) ([]byte, error) {
	return hashReader.hasher.sum(b)
}

func (hashReader *hashReader) traverse() error {
	for hashReader.Next() {
		if ion.IsContainer(hashReader.currentType) && !hashReader.IsNull() {
			err := hashReader.StepIn()
			if err != nil {
				return err
			}

			err = hashReader.traverse()
			if err != nil {
				return err
			}

			err = hashReader.StepOut()
			if err != nil {
				return err
			}
		}
	}

	return hashReader.Err()
}

// The following implements HashValue interface.

func (hashReader *hashReader) getFieldName() string {
	return hashReader.FieldName()
}

func (hashReader *hashReader) getAnnotations() []string {
	return hashReader.Annotations()
}

func (hashReader *hashReader) value() (interface{}, error) {
	switch hashReader.currentType {
	case ion.BoolType:
		return hashReader.BoolValue()
	case ion.BlobType:
		return hashReader.ionReader.ByteValue()
	case ion.ClobType:
		return hashReader.ionReader.ByteValue()
	case ion.DecimalType:
		return hashReader.DecimalValue()
	case ion.FloatType:
		return hashReader.FloatValue()
	case ion.IntType:
		return hashReader.Int64Value()
	case ion.StringType:
		return hashReader.StringValue()
	case ion.SymbolType:
		return hashReader.SymbolTable(), nil
	case ion.TimestampType:
		return hashReader.TimeValue()
	case ion.NoType:
		return ion.NoType, nil
	}

	return nil, &InvalidIonTypeError{hashReader.currentType}
}

func (hashReader *hashReader) isInStruct() bool {
	return hashReader.currentType == ion.StructType
}

func (hashReader *hashReader) ionType() ion.Type {
	return hashReader.Type()
}

func (hashReader *hashReader) isNull() bool {
	return hashReader.IsNull()
}
