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

func TestPointerEventTypeNames(t *testing.T) {
	names := PointerEventTypeNames()
	for _, name := range names {
		if i := lo.IndexOf(_PointerEventTypeNames, name); i < 0 {
			t.Fatalf("value %v not in list _PointerEventTypeNames", name)
		}
	}
	for _, name := range _PointerEventTypeNames {
		if i := lo.IndexOf(names, name); i < 0 {
			t.Fatalf("value %v not returned", name)
		}
	}
}

func TestPointerEventTypeValues(t *testing.T) {
	values := PointerEventTypeValues()
	for _, value := range values {
		if _, ok := lo.FindKey(_PointerEventTypeValue, value); !ok {
			t.Fatalf("value %v not in map _PointerEventTypeValue", value)
		}
	}
	for _, value := range _PointerEventTypeValue {
		if i := lo.IndexOf(values, value); i < 0 {
			t.Fatalf("value %v not returned", value)
		}
	}
}

func TestPointerEventType_String(t *testing.T) {
	for s, command := range _PointerEventTypeValue {
		if command.String() != s {
			t.Fatalf("String returned invalid result %s for value %v", command.String(), s)
		}
	}
}

func TestPointerEventType_IsValid(t *testing.T) {
	for _, command := range _PointerEventTypeValue {
		if !command.IsValid() {
			t.Fatalf("value %v is invalid", command)
		}
	}
}

func TestPointerEventType_MarshalText(t *testing.T) {
	for s, command := range _PointerEventTypeValue {
		if b, _ := command.MarshalText(); string(b) != s {
			t.Fatalf("Marshal %v returned invalid value %s", command, string(b))
		}
	}
}

func TestPointerEventType_UnmarshalText(t *testing.T) {
	var foo PointerEventType
	for s, command := range _PointerEventTypeValue {
		if err := foo.UnmarshalText([]byte(s)); err != nil {
			t.Fatalf("Unmarshal %s returned error %v", s, err)
		} else if foo != command {
			t.Fatalf("Unmarshal %s returned invalid value %s", s, foo)
		}
	}
}

func FuzzPointerEventType_UnmarshalText(f *testing.F) {
	for _, seed := range PointerEventTypeValues() {
		b, _ := seed.MarshalText()
		f.Add(b)
	}
	f.Fuzz(func(t *testing.T, in []byte) {
		var res PointerEventType
		err := res.UnmarshalText(in)
		if err != nil {
			if err.Error() != fmt.Errorf("%s is %w", string(in), ErrInvalidPointerEventType).Error() {
				t.Fatalf("invalid error on unmarshal %s: %v", string(in), err)
			}
		}
	})
}
