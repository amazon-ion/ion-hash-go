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

import ion "ion-go"

//TODO add docstrings
type Serializer interface {

	scalar(ionValue interface{})

	stepIn(ionValue interface{})

	stepOut()

	digest() []byte

	handleFieldName(ionValue interface{})

	update(bytes []byte)

	beginMarker()

	endMarker()

	handleAnnotationsBegin(ionValue interface{}, isContainer bool)

	handleAnnotationsEnd(ionValue interface{}, isContainer bool)

	writeSymbol(token string)

	getBytes(ionType ion.Type, ionValue interface{}, isNull bool) []byte

	getLengthLength(bytes []byte) int

// SymbolToken is currently not available
// scalarOrNullSplitParts(ionType ion.Type, symbolToken ion.SymbolToken, isNull bool, bytes byte[]) (byte, []byte)
}

// Holds the commonalities between struct and non-struct serializers.
type serializer struct {
	hashFunction IonHasher
	depth int
	hasContainerAnnotation bool
}

func (serializer *serializer)stepIn(ionValue interface{}) {
	panic("implement me")
}

func (serializer *serializer)digest() []byte{
	panic("implement me")
}

func (serializer *serializer)handleFieldName(ionValue interface{}) {
	panic("implement me")
}

func (serializer *serializer)update(bytes []byte) {
	panic("implement me")
}

func (serializer *serializer)beginMarker() {
	panic("implement me")
}

func (serializer *serializer)endMarker() {
	panic("implement me")
}

func (serializer *serializer)handleAnnotationsBegin(ionValue interface{}, isContainer bool) {
	panic("implement me")
}

func (serializer *serializer)handleAnnotationsEnd(ionValue interface{}, isContainer bool) {
	panic("implement me")
}

func (serializer *serializer)writeSymbol(token string) {
	panic("implement me")
}

func (serializer *serializer)getBytes(ionType ion.Type, ionValue interface{}, isNull bool) []byte {
	panic("implement me")
}

func (serializer *serializer)getLengthLength(bytes []byte) int {
	panic("implement me")
}

// SymbolToken is currently not available
//func (serializer *serializer)scalarOrNullSplitParts(ionType ion.Type, symbolToken ion.SymbolToken, isNull bool, bytes byte[]) (byte, []byte) {
//	panic("implement me")
//}

func escape(bytes []byte) []byte{
	panic("implement me")
}

func serializers(ionType ion.Type, ionValue interface{}, writer HashWriter) {
	panic("implement me")
}

func tq(ionValue interface{}) []byte {
	panic("implement me")
}
