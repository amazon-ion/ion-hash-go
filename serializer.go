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

import ion "ion-go"

//TODO add docstrings
type serializer interface {
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
