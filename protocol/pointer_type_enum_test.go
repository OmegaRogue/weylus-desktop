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

func TestPointerTypeNames(t *testing.T) {
	names := PointerTypeNames()
	for _, name := range names {
		if i := lo.IndexOf(_PointerTypeNames, name); i < 0 {
			t.Fatalf("value %v not in list _PointerTypeNames", name)
		}
	}
	for _, name := range _PointerTypeNames {
		if i := lo.IndexOf(names, name); i < 0 {
			t.Fatalf("value %v not returned", name)
		}
	}
}

func TestPointerTypeValues(t *testing.T) {
	values := PointerTypeValues()
	for _, value := range values {
		if _, ok := lo.FindKey(_PointerTypeValue, value); !ok {
			t.Fatalf("value %v not in map _PointerTypeValue", value)
		}
	}
	for _, value := range _PointerTypeValue {
		if i := lo.IndexOf(values, value); i < 0 {
			t.Fatalf("value %v not returned", value)
		}
	}
}

func TestPointerType_String(t *testing.T) {
	for s, command := range _PointerTypeValue {
		if command.String() != s {
			t.Fatalf("String returned invalid result %s for value %v", command.String(), s)
		}
	}
}

func TestPointerType_IsValid(t *testing.T) {
	for _, command := range _PointerTypeValue {
		if !command.IsValid() {
			t.Fatalf("value %v is invalid", command)
		}
	}
}

func TestPointerType_MarshalText(t *testing.T) {
	for s, command := range _PointerTypeValue {
		if b, _ := command.MarshalText(); string(b) != s {
			t.Fatalf("Marshal %v returned invalid value %s", command, string(b))
		}
	}
}

func TestPointerType_UnmarshalText_Correct(t *testing.T) {
	var foo PointerType
	for s, command := range _PointerTypeValue {
		if err := foo.UnmarshalText([]byte(s)); err != nil {
			if err.Error() != fmt.Errorf("%s is %w", s, ErrInvalidPointerType).Error() {
				t.Fatalf("invalid error on unmarshal %s: %v", s, err)
			}
		} else if foo != command {
			t.Fatalf("Unmarshal %s returned invalid value %s", s, foo)
		}
	}
}

func TestPointerType_UnmarshalText_Invalid(t *testing.T) {
	var foo PointerType
	for _, s := range []string{"0"} {
		if err := foo.UnmarshalText([]byte(s)); err != nil {
			if err.Error() != fmt.Errorf("%s is %w", s, ErrInvalidPointerType).Error() {
				t.Fatalf("invalid error on unmarshal %s: %v", s, err)
			}
		}
	}
}

func FuzzPointerType_UnmarshalText(f *testing.F) {
	for _, seed := range PointerTypeValues() {
		b, _ := seed.MarshalText()
		f.Add(b)
	}
	f.Fuzz(func(t *testing.T, in []byte) {
		var res PointerType
		err := res.UnmarshalText(in)
		if err != nil {
			if err.Error() != fmt.Errorf("%s is %w", string(in), ErrInvalidPointerType).Error() {
				t.Fatalf("invalid error on unmarshal %s: %v", string(in), err)
			}
		}
	})
}
