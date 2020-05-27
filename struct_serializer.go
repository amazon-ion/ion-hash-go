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
	"sort"

	"github.com/amzn/ion-go/ion"
)

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
	structSerializer.scalarSerializer.handleFieldName(ionValue)
	structSerializer.scalarSerializer.scalar(ionValue)

	digest := structSerializer.scalarSerializer.digest()
	structSerializer.appendFieldHash(digest)
}

func (structSerializer structSerializer) stepOut() {
	// Sort fieldHashes using the sortableBytes sorting interface
	sort.Sort(sortableBytes(structSerializer.fieldHashes))

	for _, digest := range structSerializer.fieldHashes {
		structSerializer.update(escape(digest))
	}

	structSerializer.baseSerializer.stepOut()
}

func (structSerializer structSerializer) stepIn(ionValue interface{}) {
	structSerializer.baseSerializer.stepIn(ionValue)
}

func (structSerializer structSerializer) digest() []byte {
	return structSerializer.baseSerializer.digest()
}

func (structSerializer structSerializer) handleFieldName(ionValue interface{}) {
	structSerializer.baseSerializer.handleFieldName(ionValue)
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
	structSerializer.baseSerializer.handleAnnotationsBegin(ionValue, isContainer)
}

func (structSerializer structSerializer) handleAnnotationsEnd(ionValue interface{}, isContainer bool) {
	structSerializer.baseSerializer.handleAnnotationsEnd(ionValue, isContainer)
}

func (structSerializer structSerializer) writeSymbol(token string) {
	structSerializer.baseSerializer.writeSymbol(token)
}

func (structSerializer structSerializer) getBytes(ionType ion.Type, ionValue interface{}, isNull bool) []byte {
	return structSerializer.baseSerializer.getBytes(ionType, ionValue, isNull)
}

func (structSerializer structSerializer) getLengthLength(bytes []byte) int {
	return structSerializer.baseSerializer.getLengthLength(bytes)
}

func (structSerializer *structSerializer) appendFieldHash(digest []byte) {
	structSerializer.fieldHashes = append(structSerializer.fieldHashes, digest)
}

func compareBytes(bytes1, bytes2 []byte) int {
	for i := 0; i < len(bytes1) && i < len(bytes2); i++ {
		byte1 := bytes1[i]
		byte2 := bytes2[i]
		if byte1 != byte2 {
			return int(byte1 - byte2)
		}
	}

	return len(bytes1) - len(bytes2)
}

// sortableBytes implements the sort.Interface so we can sort fieldHashes in stepOut()
type sortableBytes [][]byte

func (sb sortableBytes) Len() int {
	return len(sb)
}

func (sb sortableBytes) Less(i, j int) bool {
	bytes1 := sb[i]
	bytes2 := sb[j]

	return compareBytes(bytes1, bytes2) < 0
}

func (sb sortableBytes) Swap(i, j int) {
	sb[i], sb[j] = sb[j], sb[i]
}
