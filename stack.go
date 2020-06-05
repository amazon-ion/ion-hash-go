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

type stack []interface{}

// Check if stack is empty
func (s *stack) isEmpty() bool {
	return len(*s) == 0
}

// Push element into stack
func (s *stack) push(element interface{}) {
	*s = append(*s, element)
}

// Remove and return top element of stack. Return error if stack is empty.
func (s *stack) pop() (interface{}, error) {
	if s.isEmpty() {
		return nil, &InvalidOperationError{"stack", "pop", "stack is empty"}
	}

	index := len(*s) - 1   // Get the index of the top most element.
	element := (*s)[index] // Index into the slice and obtain the element.
	*s = (*s)[:index]      // Remove it from the stack by slicing it off.

	return element, nil
}

// Return top element of stack. Return error if stack is empty.
func (s *stack) peek() (interface{}, error) {
	if s.isEmpty() {
		return nil, &InvalidOperationError{"stack", "pop", "stack is empty"}
	}

	index := len(*s) - 1   // Get the index of the top most element.
	element := (*s)[index] // Index into the slice and obtain the element.

	return element, nil
}

// Return number of elements in stack
func (s *stack) size() int {
	return len(*s)
}
