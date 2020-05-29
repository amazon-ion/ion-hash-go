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

type serializer interface {
	scalar(ionValue interface{}) error

	stepIn(ionValue interface{}) error

	stepOut() error

	sum(b []byte) []byte

	handleFieldName(ionValue interface{}) error

	update(bytes []byte) error

	beginMarker() error

	endMarker() error

	handleAnnotationsBegin(ionValue interface{}) error

	handleAnnotationsEnd(ionValue interface{}, isContainer bool) error

	writeSymbol(token string) error

	getBytes(ionType ion.Type, ionValue interface{}, isNull bool) ([]byte, error)

	getLengthFieldLength(bytes []byte) (int, error)

	// SymbolToken is currently not available
	// scalarOrNullSplitParts(ionType ion.Type, symbolToken ion.SymbolToken, isNull bool, bytes byte[]) (byte, []byte)
}
