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

type scalarSerializer struct {
	baseSerializer
}

func newScalarSerializer(hashFunction IonHasher, depth int) serializer {
	return &scalarSerializer{baseSerializer{hashFunction: hashFunction, depth: depth}}
}

func (scalarSerializer scalarSerializer) scalar(ionValue interface{}) {
	panic("implement me")
}

func (scalarSerializer scalarSerializer) stepOut() {
	panic("implement me")
}

func (scalarSerializer scalarSerializer) stepIn(ionValue interface{}) {
	scalarSerializer.baseSerializer.stepIn(ionValue.(hashValue))
}

func (scalarSerializer scalarSerializer) sum(b []byte) []byte {
	return scalarSerializer.baseSerializer.sum(b)
}

func (scalarSerializer scalarSerializer) handleFieldName(ionValue interface{}) {
	scalarSerializer.baseSerializer.handleFieldName(ionValue.(hashValue))
}

func (scalarSerializer scalarSerializer) update(bytes []byte) {
	scalarSerializer.baseSerializer.update(bytes)
}

func (scalarSerializer scalarSerializer) beginMarker() {
	scalarSerializer.baseSerializer.beginMarker()
}

func (scalarSerializer scalarSerializer) endMarker() {
	scalarSerializer.baseSerializer.endMarker()
}

func (scalarSerializer scalarSerializer) handleAnnotationsBegin(ionValue interface{}) {
	scalarSerializer.baseSerializer.handleAnnotationsBegin(ionValue.(hashValue))
}

func (scalarSerializer scalarSerializer) handleAnnotationsEnd(ionValue interface{}, isContainer bool) {
	scalarSerializer.baseSerializer.handleAnnotationsEnd(ionValue.(hashValue), isContainer)
}

func (scalarSerializer scalarSerializer) writeSymbol(token string) {
	scalarSerializer.baseSerializer.writeSymbol(token)
}

func (scalarSerializer scalarSerializer) getBytes(ionType ion.Type, ionValue interface{}, isNull bool) []byte {
	bytes, _ := scalarSerializer.baseSerializer.getBytes(ionType, ionValue, isNull)
	return bytes
}

func (scalarSerializer scalarSerializer) getLengthLength(bytes []byte) int {
	length, _ := scalarSerializer.baseSerializer.getLengthLength(bytes)
	return length
}
