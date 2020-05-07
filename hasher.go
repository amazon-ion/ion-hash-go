/*
 * Copyright 2020 Amazon.com, Inc. or its affiliates. All Rights Reserved.
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

type hasher struct {
	hasherProvider IonHashProvider
	currentHasher Serializer
	hasherStack stack
}

func newHasher (hasherProvider IonHashProvider) *hasher{
	currentHasher := newScalarSerializer(hasherProvider.newHasher(), 0)

	var hasherStack stack
	hasherStack.push(currentHasher)

	return &hasher{hasherProvider, currentHasher, hasherStack}
}

func (h *hasher) scalar() {
	panic("implement me")
}

func (h *hasher) stepIn() {
	panic("implement me")
}

func (h *hasher) stepOut() {
	panic("implement me")
}

func (h *hasher) digest() []byte {
	panic("implement me")
}

func (h *hasher) depth() int {
	panic("implement me")
}
