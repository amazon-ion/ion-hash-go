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

func (structSerializer structSerializer) stepIn(ionValue interface{}) {
	structSerializer.baseSerializer.stepIn(ionValue.(hashValue))
}

// TODO: Remove digest() once we've fully sorted out how Sum(b []bytes) can replace all instances of digest()
func (structSerializer structSerializer) digest() []byte {
	return structSerializer.baseSerializer.digest()
}

func (structSerializer structSerializer) handleFieldName(ionValue interface{}) {
	structSerializer.baseSerializer.handleFieldName(ionValue.(hashValue))
}

func (structSerializer structSerializer) update(bytes []byte) {
	structSerializer.baseSerializer.update(bytes)
}

func (structSerializer structSerializer) beginMarker() {
	structSerializer.baseSerializer.beginMarker()
}

func (structSerializer structSerializer) endMarker() {
	structSerializer.baseSerializer.endMarker()
}

func (structSerializer structSerializer) handleAnnotationsBegin(ionValue interface{}, isContainer bool) {
	structSerializer.baseSerializer.handleAnnotationsBegin(ionValue.(hashValue), isContainer)
}

func (structSerializer structSerializer) handleAnnotationsEnd(ionValue interface{}, isContainer bool) {
	structSerializer.baseSerializer.handleAnnotationsEnd(ionValue.(hashValue), isContainer)
}

func (structSerializer structSerializer) writeSymbol(token string) {
	structSerializer.baseSerializer.writeSymbol(token)
}

func (structSerializer structSerializer) getBytes(ionType ion.Type, ionValue interface{}, isNull bool) []byte {
	bytes, _ := structSerializer.baseSerializer.getBytes(ionType, ionValue, isNull)
	return bytes
}

func (structSerializer structSerializer) getLengthLength(bytes []byte) int {
	length, _ := structSerializer.baseSerializer.getLengthLength(bytes)
	return length
}

func (structSerializer *structSerializer) appendFieldHash(digest []byte) {
	panic("implement me")
}

func compareBytes(bs1, bs2 []byte) []int16 {
	panic("implement me")
}
