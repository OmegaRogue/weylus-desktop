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

func FuzzParseWeylusCommand(f *testing.F) {
	for _, seed := range append(WeylusCommandNames(), "fdsaghjkfcgkjhsdgvf") {
		f.Add(seed)
	}
	f.Fuzz(func(t *testing.T, in string) {
		_, err := ParseWeylusCommand(in)
		if err != nil {
			if err.Error() != fmt.Errorf("%s is %w", in, ErrInvalidWeylusCommand).Error() {
				t.Fatalf("invalid error on ParseWeylusCommand %s: %v", in, err)
			}
		}
	})
}

func FuzzParseWeylusResponse(f *testing.F) {
	for _, seed := range append(WeylusResponseNames(), "fdsaghjkfcgkjhsdgvf") {
		f.Add(seed)
	}
	f.Fuzz(func(t *testing.T, in string) {
		_, err := ParseWeylusResponse(in)
		if err != nil {
			if err.Error() != fmt.Errorf("%s is %w", in, ErrInvalidWeylusResponse).Error() {
				t.Fatalf("invalid error on ParseWeylusResponse %s: %v", in, err)
			}
		}
	})
}

func TestWeylusCommandNames(t *testing.T) {
	names := WeylusCommandNames()
	for _, name := range names {
		if i := lo.IndexOf(_WeylusCommandNames, name); i < 0 {
			t.Fatalf("value %v not in list _WeylusCommandNames", name)
		}
	}
	for _, name := range _WeylusCommandNames {
		if i := lo.IndexOf(names, name); i < 0 {
			t.Fatalf("value %v not returned", name)
		}
	}
}

func TestWeylusCommandValues(t *testing.T) {
	values := WeylusCommandValues()
	for _, value := range values {
		if _, ok := lo.FindKey(_WeylusCommandValue, value); !ok {
			t.Fatalf("value %v not in map _WeylusCommandValue", value)
		}
	}
	for _, value := range _WeylusCommandValue {
		if i := lo.IndexOf(values, value); i < 0 {
			t.Fatalf("value %v not returned", value)
		}
	}
}

func TestWeylusCommand_String(t *testing.T) {
	for s, command := range _WeylusCommandValue {
		if command.String() != s {
			t.Fatalf("String returned invalid result %s for value %v", command.String(), s)
		}
	}
}

func TestWeylusCommand_IsValid(t *testing.T) {
	for _, command := range _WeylusCommandValue {
		if !command.IsValid() {
			t.Fatalf("value %v is invalid", command)
		}
	}
}

func TestWeylusCommand_MarshalText(t *testing.T) {
	for s, command := range _WeylusCommandValue {
		if b, _ := command.MarshalText(); string(b) != s {
			t.Fatalf("Marshal %v returned invalid value %s", command, string(b))
		}
	}
}

func TestWeylusCommand_UnmarshalText(t *testing.T) {
	var foo WeylusCommand
	for s, command := range _WeylusCommandValue {
		if err := foo.UnmarshalText([]byte(s)); err != nil {
			t.Fatalf("Unmarshal %s returned error %v", s, err)
		} else if foo != command {
			t.Fatalf("Unmarshal %s returned invalid value %s", s, foo)
		}
	}
}

func FuzzWeylusCommand_UnmarshalText(f *testing.F) {
	for _, seed := range WeylusCommandValues() {
		b, _ := seed.MarshalText()
		f.Add(b)
	}
	f.Fuzz(func(t *testing.T, in []byte) {
		var res WeylusCommand
		err := res.UnmarshalText(in)
		if err != nil {
			if err.Error() != fmt.Errorf("%s is %w", string(in), ErrInvalidWeylusCommand).Error() {
				t.Fatalf("invalid error on unmarshal %s: %v", string(in), err)
			}
		}
	})
}

func TestWeylusResponseNames(t *testing.T) {
	names := WeylusResponseNames()
	for _, name := range names {
		if i := lo.IndexOf(_WeylusResponseNames, name); i < 0 {
			t.Fatalf("value %v not in list _WeylusResponseNames", name)
		}
	}
	for _, name := range _WeylusResponseNames {
		if i := lo.IndexOf(names, name); i < 0 {
			t.Fatalf("value %v not returned", name)
		}
	}
}

func TestWeylusResponseValues(t *testing.T) {
	values := WeylusResponseValues()
	for _, value := range values {
		if _, ok := lo.FindKey(_WeylusResponseValue, value); !ok {
			t.Fatalf("value %v not in map _WeylusResponseValue", value)
		}
	}
	for _, value := range _WeylusResponseValue {
		if i := lo.IndexOf(values, value); i < 0 {
			t.Fatalf("value %v not returned", value)
		}
	}
}

func TestWeylusResponse_String(t *testing.T) {
	for s, command := range _WeylusResponseValue {
		if command.String() != s {
			t.Fatalf("String returned invalid result %s for value %v", command.String(), s)
		}
	}
}

func TestWeylusResponse_IsValid(t *testing.T) {
	for _, command := range _WeylusResponseValue {
		if !command.IsValid() {
			t.Fatalf("value %v is invalid", command)
		}
	}
}

func TestWeylusResponse_MarshalText(t *testing.T) {
	for s, command := range _WeylusResponseValue {
		if b, _ := command.MarshalText(); string(b) != s {
			t.Fatalf("Marshal %v returned invalid value %s", command, string(b))
		}
	}
}

func TestWeylusResponse_UnmarshalText(t *testing.T) {
	var foo WeylusResponse
	for s, command := range _WeylusResponseValue {
		if err := foo.UnmarshalText([]byte(s)); err != nil {
			t.Fatalf("Unmarshal %s returned error %v", s, err)
		} else if foo != command {
			t.Fatalf("Unmarshal %s returned invalid value %s", s, foo)
		}
	}
}

func FuzzWeylusResponse_UnmarshalText(f *testing.F) {
	for _, seed := range WeylusResponseValues() {
		b, _ := seed.MarshalText()
		f.Add(b)
	}
	f.Fuzz(func(t *testing.T, in []byte) {
		var res WeylusResponse
		err := res.UnmarshalText(in)
		if err != nil {
			if err.Error() != fmt.Errorf("%s is %w", string(in), ErrInvalidWeylusResponse).Error() {
				t.Fatalf("invalid error on unmarshal %s: %v", string(in), err)
			}
		}
	})
}