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

type hasher struct {
	hasherProvider IonHasherProvider
	currentHasher  serializer
	hasherStack    stack
}

func newHasher(hasherProvider IonHasherProvider) (*hasher, error) {
	newHasher, err := hasherProvider.newHasher()
	if err != nil {
		return nil, err
	}

	currentHasher := newScalarSerializer(newHasher, 0)

	var hasherStack stack
	hasherStack.push(currentHasher)

	return &hasher{hasherProvider, currentHasher, hasherStack}, nil
}

func (h *hasher) scalar(ionValue hashValue) error {
	return h.currentHasher.scalar(ionValue)
}

func (h *hasher) stepIn(ionValue hashValue) error {
	var hashFunction IonHasher

	_, ok := h.currentHasher.(*structSerializer)
	if ok {
		newHasher, err := h.hasherProvider.newHasher()
		if err != nil {
			return err
		}

		hashFunction = newHasher
	} else {
		hashFunction = h.currentHasher.(*scalarSerializer).hashFunction
	}

	if ionValue.ionType() == ion.StructType {
		newStructSerializer, err := newStructSerializer(hashFunction, 0, h.hasherProvider)
		if err != nil {
			return err
		}

		h.currentHasher = newStructSerializer
	} else {
		h.currentHasher = newScalarSerializer(hashFunction, 0)
	}

	h.hasherStack.push(h.currentHasher)
	return h.currentHasher.stepIn(ionValue)
}

func (h *hasher) stepOut() error {
	if h.depth() == 0 {
		return &InvalidOperationError{"hasher", "stepOut", "Depth is zero. Hasher cannot step out any further"}
	}

	err := h.currentHasher.stepOut()
	if err != nil {
		return err
	}

	poppedHasher, err := h.hasherStack.pop()
	if err != nil {
		return err
	}
	peekedHasher, err := h.hasherStack.peek()
	if err != nil {
		return err
	}

	h.currentHasher = peekedHasher.(serializer)

	structHasher, ok := h.currentHasher.(*structSerializer)
	if ok {
		sum := poppedHasher.(serializer).sum(nil)
		structHasher.appendFieldHash(sum)
	}

	return nil
}

func (h *hasher) sum(b []byte) ([]byte, error) {
	if h.depth() != 0 {
		return nil, &InvalidOperationError{
			"hasher", "sum", "A sum may only be provided at the same depth hashing started"}
	}

	return h.currentHasher.sum(b), nil
}

func (h *hasher) depth() int {
	return h.hasherStack.size() - 1
}
