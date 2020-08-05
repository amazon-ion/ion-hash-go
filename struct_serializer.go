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

func (ss *structSerializer) scalar(ionValue interface{}) error {
	err := ss.scalarSerializer.handleFieldName(ionValue)
	if err != nil {
		return err
	}

	err = ss.scalarSerializer.scalar(ionValue)
	if err != nil {
		return err
	}

	sum := ss.scalarSerializer.sum(nil)
	ss.appendFieldHash(sum)

	return nil
}

func (ss *structSerializer) stepOut() error {
	// Sort fieldHashes using the sortableBytes sorting interface.
	sort.Sort(sortableBytes(ss.fieldHashes))

	for _, digest := range ss.fieldHashes {
		err := ss.write(escape(digest))
		if err != nil {
			return err
		}
	}

	return ss.baseSerializer.stepOut()
}

func (ss *structSerializer) stepIn(ionValue interface{}) error {
	return ss.baseSerializer.stepIn(ionValue.(hashValue))
}

func (ss *structSerializer) handleFieldName(ionValue interface{}) error {
	return ss.baseSerializer.handleFieldName(ionValue.(hashValue))
}

func (ss *structSerializer) handleAnnotationsBegin(ionValue interface{}) error {
	return ss.baseSerializer.handleAnnotationsBegin(ionValue.(hashValue), false)
}

func (ss *structSerializer) handleAnnotationsEnd(ionValue interface{}, isContainer bool) error {
	return ss.baseSerializer.handleAnnotationsEnd(ionValue.(hashValue), isContainer)
}

func (ss *structSerializer) appendFieldHash(sum []byte) {
	ss.fieldHashes = append(ss.fieldHashes, sum)
}

func (ss *structSerializer) scalarOrNullSplitParts(
	ionType ion.Type, isNull bool, bytes []byte) (byte, []byte, error) {

	return ss.baseSerializer.scalarOrNullSplitParts(ionType, isNull, bytes)
}
