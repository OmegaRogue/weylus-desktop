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

package protocol

import (
	"fmt"
	"testing"

	"github.com/samber/lo"
)

func TestKeyboardEventTypeNames(t *testing.T) {
	names := KeyboardEventTypeNames()
	for _, name := range names {
		if i := lo.IndexOf(_KeyboardEventTypeNames, name); i < 0 {
			t.Fatalf("value %v not in list _KeyboardEventTypeNames", name)
		}
	}
	for _, name := range _KeyboardEventTypeNames {
		if i := lo.IndexOf(names, name); i < 0 {
			t.Fatalf("value %v not returned", name)
		}
	}
}

func TestKeyboardEventTypeValues(t *testing.T) {
	values := KeyboardEventTypeValues()
	for _, value := range values {
		if _, ok := lo.FindKey(_KeyboardEventTypeValue, value); !ok {
			t.Fatalf("value %v not in map _KeyboardEventTypeValue", value)
		}
	}
	for _, value := range _KeyboardEventTypeValue {
		if i := lo.IndexOf(values, value); i < 0 {
			t.Fatalf("value %v not returned", value)
		}
	}
}

func TestKeyboardEventType_String(t *testing.T) {
	for s, command := range _KeyboardEventTypeValue {
		if command.String() != s {
			t.Fatalf("String returned invalid result %s for value %v", command.String(), s)
		}
	}
}

func TestKeyboardEventType_IsValid(t *testing.T) {
	for _, command := range _KeyboardEventTypeValue {
		if !command.IsValid() {
			t.Fatalf("value %v is invalid", command)
		}
	}
}

func TestKeyboardEventType_MarshalText(t *testing.T) {
	for s, command := range _KeyboardEventTypeValue {
		if b, _ := command.MarshalText(); string(b) != s {
			t.Fatalf("Marshal %v returned invalid value %s", command, string(b))
		}
	}
}

func TestKeyboardEventType_UnmarshalText_Correct(t *testing.T) {
	var foo KeyboardEventType
	for s, command := range _KeyboardEventTypeValue {
		if err := foo.UnmarshalText([]byte(s)); err != nil {
			if err.Error() != fmt.Errorf("%s is %w", s, ErrInvalidKeyboardEventType).Error() {
				t.Fatalf("invalid error on unmarshal %s: %v", s, err)
			}
		} else if foo != command {
			t.Fatalf("Unmarshal %s returned invalid value %s", s, foo)
		}
	}
}

func TestKeyboardEventType_UnmarshalText_Invalid(t *testing.T) {
	var foo KeyboardEventType
	for _, s := range []string{"0"} {
		if err := foo.UnmarshalText([]byte(s)); err != nil {
			if err.Error() != fmt.Errorf("%s is %w", s, ErrInvalidKeyboardEventType).Error() {
				t.Fatalf("invalid error on unmarshal %s: %v", s, err)
			}
		}
	}
}

func FuzzKeyboardEventType_UnmarshalText(f *testing.F) {
	for _, seed := range KeyboardEventTypeValues() {
		b, _ := seed.MarshalText()
		f.Add(b)
	}
	f.Fuzz(func(t *testing.T, in []byte) {
		var res KeyboardEventType
		err := res.UnmarshalText(in)
		if err != nil {
			if err.Error() != fmt.Errorf("%s is %w", string(in), ErrInvalidKeyboardEventType).Error() {
				t.Fatalf("invalid error on unmarshal %s: %v", string(in), err)
			}
		}
	})
}
