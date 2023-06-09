/*
 * Copyright © 2023 omegarogue
 * SPDX-License-Identifier: AGPL-3.0-or-later
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as
 * published by the Free Software Foundation, either version 3 of the
 * License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program. If not, see <https://www.gnu.org/licenses/>.
 */

// Code generated by go-enum DO NOT EDIT.
// Version:
// Revision:
// Build Date:
// Built By:

package protocol

import (
	"fmt"
	"strings"
)

const (
	// PointerTypeUnknown is a PointerType of type unknown.
	PointerTypeUnknown PointerType = ""
	// PointerTypeMouse is a PointerType of type mouse.
	PointerTypeMouse PointerType = "mouse"
	// PointerTypePen is a PointerType of type pen.
	PointerTypePen PointerType = "pen"
	// PointerTypeTouch is a PointerType of type touch.
	PointerTypeTouch PointerType = "touch"
)

var ErrInvalidPointerType = fmt.Errorf("not a valid PointerType, try [%s]", strings.Join(_PointerTypeNames, ", "))

var _PointerTypeNames = []string{
	string(PointerTypeUnknown),
	string(PointerTypeMouse),
	string(PointerTypePen),
	string(PointerTypeTouch),
}

// PointerTypeNames returns a list of possible string values of PointerType.
func PointerTypeNames() []string {
	tmp := make([]string, len(_PointerTypeNames))
	copy(tmp, _PointerTypeNames)
	return tmp
}

// PointerTypeValues returns a list of the values for PointerType
func PointerTypeValues() []PointerType {
	return []PointerType{
		PointerTypeUnknown,
		PointerTypeMouse,
		PointerTypePen,
		PointerTypeTouch,
	}
}

// String implements the Stringer interface.
func (x PointerType) String() string {
	return string(x)
}

// String implements the Stringer interface.
func (x PointerType) IsValid() bool {
	_, err := ParsePointerType(string(x))
	return err == nil
}

var _PointerTypeValue = map[string]PointerType{
	"":      PointerTypeUnknown,
	"mouse": PointerTypeMouse,
	"pen":   PointerTypePen,
	"touch": PointerTypeTouch,
}

// ParsePointerType attempts to convert a string to a PointerType.
func ParsePointerType(name string) (PointerType, error) {
	if x, ok := _PointerTypeValue[name]; ok {
		return x, nil
	}
	return PointerType(""), fmt.Errorf("%s is %w", name, ErrInvalidPointerType)
}

// MarshalText implements the text marshaller method.
func (x PointerType) MarshalText() ([]byte, error) {
	return []byte(string(x)), nil
}

// UnmarshalText implements the text unmarshaller method.
func (x *PointerType) UnmarshalText(text []byte) error {
	tmp, err := ParsePointerType(string(text))
	if err != nil {
		return err
	}
	*x = tmp
	return nil
}
