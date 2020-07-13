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
	newHasher, err := hashFunctionProvider.NewHasher()
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

func (structSerializer *structSerializer) handleFieldName(ionValue interface{}) error {
	return structSerializer.baseSerializer.handleFieldName(ionValue.(hashValue))
}

func (structSerializer *structSerializer) handleAnnotationsBegin(ionValue interface{}) error {
	return structSerializer.baseSerializer.handleAnnotationsBegin(ionValue.(hashValue), false)
}

func (structSerializer *structSerializer) handleAnnotationsEnd(ionValue interface{}, isContainer bool) error {
	return structSerializer.baseSerializer.handleAnnotationsEnd(ionValue.(hashValue), isContainer)
}

func (structSerializer *structSerializer) appendFieldHash(sum []byte) {
	structSerializer.fieldHashes = append(structSerializer.fieldHashes, sum)
}

func (structSerializer *structSerializer) scalarOrNullSplitParts(
	ionType ion.Type, isNull bool, bytes []byte) (byte, []byte, error) {

	return structSerializer.baseSerializer.scalarOrNullSplitParts(ionType, isNull, bytes)
}
