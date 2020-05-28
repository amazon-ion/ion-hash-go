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

import "github.com/amzn/ion-go/ion"

type structSerializer struct {
	baseSerializer

	scalarSerializer serializer
	fieldHashes      [][]byte
}

func newStructSerializer(hashFunction IonHasher, depth int, hashFunctionProvider IonHasherProvider) serializer {
	return &structSerializer{
		baseSerializer:   baseSerializer{hashFunction: hashFunction, depth: depth},
		scalarSerializer: newScalarSerializer(hashFunctionProvider.newHasher(), depth+1)}
}

func (structSerializer structSerializer) scalar(ionValue interface{}) {
	panic("implement me")
}

func (structSerializer structSerializer) stepOut() {
	panic("implement me")
}

func (structSerializer structSerializer) stepIn(ionValue interface{}) error {
	return structSerializer.baseSerializer.stepIn(ionValue.(hashValue))
}

func (structSerializer structSerializer) sum(b []byte) []byte {
	return structSerializer.baseSerializer.sum(b)
}

func (structSerializer structSerializer) handleFieldName(ionValue interface{}) error {
	return structSerializer.baseSerializer.handleFieldName(ionValue.(hashValue))
}

func (structSerializer structSerializer) update(bytes []byte) error {
	return structSerializer.baseSerializer.update(bytes)
}

func (structSerializer structSerializer) beginMarker() error {
	return structSerializer.baseSerializer.beginMarker()
}

func (structSerializer structSerializer) endMarker() error {
	return structSerializer.baseSerializer.endMarker()
}

func (structSerializer structSerializer) handleAnnotationsBegin(ionValue interface{}) error {
	return structSerializer.baseSerializer.handleAnnotationsBegin(ionValue.(hashValue))
}

func (structSerializer structSerializer) handleAnnotationsEnd(ionValue interface{}, isContainer bool) error {
	return structSerializer.baseSerializer.handleAnnotationsEnd(ionValue.(hashValue), isContainer)
}

func (structSerializer structSerializer) writeSymbol(token string) error {
	return structSerializer.baseSerializer.writeSymbol(token)
}

func (structSerializer structSerializer) getBytes(ionType ion.Type, ionValue interface{}, isNull bool) []byte {
	bytes, _ := structSerializer.baseSerializer.getBytes(ionType, ionValue, isNull)
	return bytes
}

func (structSerializer structSerializer) getLengthFieldLength(bytes []byte) (int, error) {
	return structSerializer.baseSerializer.getLengthFieldLength(bytes)
}

func (structSerializer *structSerializer) appendFieldHash(sum []byte) {
	panic("implement me")
}

func compareBytes(bs1, bs2 []byte) []int16 {
	panic("implement me")
}
