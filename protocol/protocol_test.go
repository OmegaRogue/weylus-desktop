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

	"github.com/pkg/errors"
)

func TestWrapMessage(t *testing.T) {

	//case PointerEvent:
	//	wrapper[WeylusCommandPointerEvent] = a
	//	case WheelEvent:
	//	wrapper[WeylusCommandWheelEvent] = a
	//	case KeyboardEvent:
	//	wrapper[WeylusCommandKeyboardEvent] = a
	//	case Config:
	//	wrapper[WeylusCommandConfig] = a
	//	case string, WeylusCommand:
	t.Run("PointerEvent", func(t *testing.T) {
		out := WrapMessage(PointerEvent{})
		t.Log(out)
	})

	t.Run("WheelEvent", func(t *testing.T) {
		out := WrapMessage(WheelEvent{})
		t.Log(out)
	})
	t.Run("KeyboardEvent", func(t *testing.T) {
		out := WrapMessage(KeyboardEvent{})
		t.Log(out)
	})
	t.Run("Config", func(t *testing.T) {
		out := WrapMessage(Config{})
		t.Log(out)
	})
	t.Run("WeylusCommand", func(t *testing.T) {
		out := WrapMessage(WeylusCommandTryGetFrame)
		t.Log(out)
	})
	t.Run("string", func(t *testing.T) {
		out := WrapMessage("")
		t.Log(out)
	})
	t.Run("underlyingString", func(t *testing.T) {
		out := WrapMessage(underlyingString("test"))
		t.Log(out)
	})
}

func TestParseMessage(t *testing.T) {
}

func FuzzParseMessage(f *testing.F) {
	f.Fuzz(func(t *testing.T, in []byte) {
		out, err := ParseMessage(in)
		if err != nil {
			t.Error(errors.Cause(err))
		}
		t.Log(out)
	})
}

func TestWeylusError_Error(t *testing.T) {
	err := WeylusError{
		ErrorMessage: "test",
	}
	if err.Error() != fmt.Sprintf("WeylusError: %s", err.ErrorMessage) {
		t.Errorf("invalid error message: %v", err)
	}
}

func TestWeylusConfigError_Error(t *testing.T) {
	err := WeylusConfigError{
		ErrorMessage: "test",
	}
	if err.Error() != fmt.Sprintf("WeylusConfigError: %s", err.ErrorMessage) {
		t.Errorf("invalid error message: %v", err)
	}
}
