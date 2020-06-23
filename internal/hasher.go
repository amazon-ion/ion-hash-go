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

package internal

import (
	"github.com/amzn/ion-go/ion"
	"github.com/amzn/ion-hash-go/ihp"
)

// Hasher struct is responsible for hashing Ion values.
type Hasher struct {
	hasherProvider ihp.IonHasherProvider
	currentHasher  serializer
	hasherStack    stack
}

// NewHasher creates a new Hasher.
func NewHasher(hasherProvider ihp.IonHasherProvider) (*Hasher, error) {
	newHasher, err := hasherProvider.NewHasher()
	if err != nil {
		return nil, err
	}

	currentHasher := newScalarSerializer(newHasher, 0)

	var hasherStack stack
	hasherStack.push(currentHasher)

	return &Hasher{hasherProvider, currentHasher, hasherStack}, nil
}

// Scalar hashes a scalar Ion value.
func (h *Hasher) Scalar(ionValue HashValue) error {
	return h.currentHasher.scalar(ionValue)
}

// StepIn will step in the Ion container.
func (h *Hasher) StepIn(ionValue HashValue) error {
	var hashFunction ihp.IonHasher

	_, ok := h.currentHasher.(*structSerializer)
	if ok {
		newHasher, err := h.hasherProvider.NewHasher()
		if err != nil {
			return err
		}

		hashFunction = newHasher
	} else {
		hashFunction = h.currentHasher.(*scalarSerializer).hashFunction
	}

	if ionValue.IonType() == ion.StructType {
		newStructSerializer, err := newStructSerializer(hashFunction, h.Depth(), h.hasherProvider)
		if err != nil {
			return err
		}

		h.currentHasher = newStructSerializer
	} else {
		h.currentHasher = newScalarSerializer(hashFunction, h.Depth())
	}

	h.hasherStack.push(h.currentHasher)
	return h.currentHasher.stepIn(ionValue)
}

// StepOut will step out of the Ion container.
func (h *Hasher) StepOut() error {
	if h.Depth() == 0 {
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

// Sum appends the current hash to b and returns the resulting slice.
func (h *Hasher) Sum(b []byte) ([]byte, error) {
	if h.Depth() != 0 {
		return nil, &InvalidOperationError{
			"hasher", "sum", "A sum may only be provided at the same depth hashing started"}
	}

	return h.currentHasher.sum(b), nil
}

// Depth returns the current size of the Hasher stack.
func (h *Hasher) Depth() int {
	return h.hasherStack.size() - 1
}
