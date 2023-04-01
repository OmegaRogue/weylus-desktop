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
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

package protocol

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/samber/lo"
)

func TestKeyboardLocationNames(t *testing.T) {
	names := KeyboardLocationNames()
	for _, name := range names {
		if i := lo.IndexOf(_KeyboardLocationNames, name); i < 0 {
			t.Fatalf("value %v not in list _KeyboardLocationNames", name)
		}
	}
	for _, name := range _KeyboardLocationNames {
		if i := lo.IndexOf(names, name); i < 0 {
			t.Fatalf("value %v not returned", name)
		}
	}
}

func TestKeyboardLocationValues(t *testing.T) {
	values := KeyboardLocationValues()
	for _, value := range values {
		if _, ok := lo.FindKey(_KeyboardLocationValue, value); !ok {
			t.Fatalf("value %v not in map _KeyboardLocationValue", value)
		}
	}
	for _, value := range _KeyboardLocationValue {
		if i := lo.IndexOf(values, value); i < 0 {
			t.Fatalf("value %v not returned", value)
		}
	}
}

func TestKeyboardLocation_String(t *testing.T) {
	for s, command := range _KeyboardLocationValue {
		if command.String() != s {
			t.Fatalf("String returned invalid result %s for value %v", command.String(), s)
		}
	}
}

func FuzzKeyboardLocation_String(f *testing.F) {
	for _, seed := range append(KeyboardLocationValues(), 4, -1) {
		f.Add(int(seed))
	}
	f.Fuzz(func(t *testing.T, in int) {
		val := KeyboardLocation(in)
		out := val.String()

		if lo.Contains(KeyboardLocationValues(), val) {
			if !lo.Contains(KeyboardLocationNames(), out) {
				t.Errorf("invalid string %s on valid value %v", out, _KeyboardLocationMap[val])
			}
		} else {
			if s := fmt.Sprintf("KeyboardLocation(%d)", in); out != s {
				t.Errorf("invalid string %s on unknown value %d, should be %s", out, in, s)
			}
		}
	})
}

func TestKeyboardLocation_MarshalText(t *testing.T) {
	for s, command := range _KeyboardLocationValue {
		if b, _ := command.MarshalText(); string(b) != s {
			t.Fatalf("Marshal %v returned invalid value %s", command, string(b))
		}
	}
}

func TestKeyboardLocation_UnmarshalText(t *testing.T) {
	var foo KeyboardLocation
	for s, command := range _KeyboardLocationValue {
		if err := foo.UnmarshalText([]byte(s)); err != nil {
			t.Fatalf("Unmarshal %s returned error %v", s, err)
		} else if foo != command {
			t.Fatalf("Unmarshal %s returned invalid value %s", s, foo)
		}
	}
}

func FuzzParseKeyboardLocation(f *testing.F) {
	for _, seed := range append(KeyboardLocationNames(), "fdsaghjkfcgkjhsdgvf") {
		f.Add(seed)
	}
	f.Fuzz(func(t *testing.T, in string) {
		res, err := ParseKeyboardLocation(in)
		if lo.Contains(_KeyboardLocationNames, in) {
			if err != nil {
				t.Errorf("Got error on correct value %s: %v", in, err)
			}
		} else {
			if err == nil {
				t.Errorf("Got no error on invalid value %s, result: %v", in, res)
			}
		}
	})
}

func FuzzKeyboardLocation_UnmarshalText(f *testing.F) {
	for _, seed := range KeyboardLocationValues() {
		b, _ := seed.MarshalText()
		f.Add(b)
	}
	f.Fuzz(func(t *testing.T, in []byte) {
		var res KeyboardLocation
		err := res.UnmarshalText(in)
		if err != nil {
			if res2, err2 := res.MarshalText(); err2 != nil {
				t.Fatalf("error on remarshal %v: %v", res, err2)
			} else if string(in) != string(res2) && string(in) != strconv.Itoa(int(res)) && res != KeyboardLocationStandard {
				t.Errorf("Values dont match after remarshal: %v != %v", string(in), string(res2))
			}
			if err.Error() != fmt.Errorf("%s is %w", string(in), ErrInvalidKeyboardLocation).Error() {
				t.Fatalf("invalid error on unmarshal %s: %v", string(in), err)
			}
		}
	})
}
