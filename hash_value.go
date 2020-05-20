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

type hashValue interface {
	getFieldName() string

	getAnnotations() []string

	// SymbolTokens are not available right now.
	//fieldNameSymbol() ion.FieldNameSymbol

	isNull() bool

	ionType() ion.Type

	value() (interface{}, error)

	isInStruct() bool
}
