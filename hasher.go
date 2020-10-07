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
	"github.com/amzn/ion-go/ion"
	"github.com/amzn/ion-hash-go/internal"
)

type hasher struct {
	hasherProvider IonHasherProvider
	currentHasher  serializer
	hasherStack    internal.Stack
}

func newHasher(hasherProvider IonHasherProvider) (*hasher, error) {
	newHasher, err := hasherProvider.NewHasher()
	if err != nil {
		return nil, err
	}

	currentHasher := newScalarSerializer(newHasher, 0)

	var hasherStack internal.Stack
	hasherStack.Push(currentHasher)

	return &hasher{hasherProvider, currentHasher, hasherStack}, nil
}

func (h *hasher) scalar(ionValue hashValue) error {
	return h.currentHasher.scalar(ionValue)
}

func (h *hasher) stepIn(ionValue hashValue) error {
	var hashFunction IonHasher

	if _, ok := h.currentHasher.(*structSerializer); ok {
		newHasher, err := h.hasherProvider.NewHasher()
		if err != nil {
			return err
		}

		hashFunction = newHasher
	} else {
		hashFunction = h.currentHasher.(*scalarSerializer).hashFunction
	}

	if ionValue.Type() == ion.StructType {
		newStructSerializer, err := newStructSerializer(hashFunction, h.depth(), h.hasherProvider)
		if err != nil {
			return err
		}

		h.currentHasher = newStructSerializer
	} else {
		h.currentHasher = newScalarSerializer(hashFunction, h.depth())
	}

	h.hasherStack.Push(h.currentHasher)
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

	poppedHasher, err := h.hasherStack.Pop()
	if err != nil {
		return &InvalidOperationError{
			"hasher",
			"stepOut",
			err.Error(),
		}
	}
	peekedHasher, err := h.hasherStack.Peek()
	if err != nil {
		return &InvalidOperationError{
			"hasher",
			"stepOut",
			err.Error(),
		}
	}

	h.currentHasher = peekedHasher.(serializer)

	if structHasher, ok := h.currentHasher.(*structSerializer); ok {
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
	return h.hasherStack.Size() - 1
}
