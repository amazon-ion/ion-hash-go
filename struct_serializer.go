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

func newStructSerializer(hashFunction IonHasher, depth int, hashFunctionProvider IonHasherProvider) (serializer, error) {
	newHasher, err := hashFunctionProvider.newHasher()
	if err != nil {
		return nil, err
	}

	return &structSerializer{
		baseSerializer:   baseSerializer{hashFunction: hashFunction, depth: depth},
		scalarSerializer: newScalarSerializer(newHasher, depth+1)}, nil
}

func (structSerializer *structSerializer) scalar(ionValue interface{}) error {
	err := structSerializer.scalarSerializer.handleFieldName(ionValue)
	if err != nil {
		return err
	}

	err = structSerializer.scalarSerializer.scalar(ionValue)
	if err != nil {
		return err
	}

	sum := structSerializer.scalarSerializer.sum(nil)
	structSerializer.appendFieldHash(sum)

	return nil
}

func (structSerializer *structSerializer) stepOut() error {
	// Sort fieldHashes using the sortableBytes sorting interface
	sort.Sort(sortableBytes(structSerializer.fieldHashes))

	for _, digest := range structSerializer.fieldHashes {
		err := structSerializer.write(escape(digest))
		if err != nil {
			return err
		}
	}

	return structSerializer.baseSerializer.stepOut()
}

func (structSerializer *structSerializer) stepIn(ionValue interface{}) error {
	return structSerializer.baseSerializer.stepIn(ionValue.(hashValue))
}

func (structSerializer *structSerializer) sum(b []byte) []byte {
	return structSerializer.baseSerializer.sum(b)
}

func (structSerializer *structSerializer) handleFieldName(ionValue interface{}) error {
	return structSerializer.baseSerializer.handleFieldName(ionValue.(hashValue))
}

func (structSerializer *structSerializer) write(bytes []byte) error {
	return structSerializer.baseSerializer.write(bytes)
}

func (structSerializer *structSerializer) beginMarker() error {
	return structSerializer.baseSerializer.beginMarker()
}

func (structSerializer *structSerializer) endMarker() error {
	return structSerializer.baseSerializer.endMarker()
}

func (structSerializer *structSerializer) handleAnnotationsBegin(ionValue interface{}) error {
	return structSerializer.baseSerializer.handleAnnotationsBegin(ionValue.(hashValue), false)
}

func (structSerializer *structSerializer) handleAnnotationsEnd(ionValue interface{}, isContainer bool) error {
	return structSerializer.baseSerializer.handleAnnotationsEnd(ionValue.(hashValue), isContainer)
}

func (structSerializer *structSerializer) writeSymbol(token string) error {
	return structSerializer.baseSerializer.writeSymbol(token)
}

func (structSerializer *structSerializer) getBytes(ionType ion.Type, ionValue interface{}, isNull bool) ([]byte, error) {
	return structSerializer.baseSerializer.getBytes(ionType, ionValue.(hashValue), isNull)
}

func (structSerializer *structSerializer) getLengthFieldLength(bytes []byte) (int, error) {
	return structSerializer.baseSerializer.getLengthFieldLength(bytes)
}

func (structSerializer *structSerializer) appendFieldHash(sum []byte) {
	structSerializer.fieldHashes = append(structSerializer.fieldHashes, sum)
}

func (structSerializer *structSerializer) scalarOrNullSplitParts(
	ionType ion.Type, isNull bool, bytes []byte) (byte, []byte, error) {

	return structSerializer.baseSerializer.scalarOrNullSplitParts(ionType, isNull, bytes)
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
