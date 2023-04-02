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
	"strings"
	"testing"

	"github.com/OmegaRogue/weylus-desktop/utils"
)

func TestWrapMessage(t *testing.T) {
	t.Run("PointerEvent", func(t *testing.T) {
		out := WrapMessage(PointerEvent{})
		if res, ok := out.(map[WeylusCommand]PointerEvent); !ok {
			t.Errorf("invalid result: %v", out)
		} else if _, ok := res[WeylusCommandPointerEvent]; !ok {
			t.Errorf("invalid key in map: %v", res)
		}
	})

	t.Run("WheelEvent", func(t *testing.T) {
		out := WrapMessage(WheelEvent{})
		if res, ok := out.(map[WeylusCommand]WheelEvent); !ok {
			t.Errorf("invalid result: %v", out)
		} else if _, ok := res[WeylusCommandWheelEvent]; !ok {
			t.Errorf("invalid key in map: %v", res)
		}
	})
	t.Run("KeyboardEvent", func(t *testing.T) {
		out := WrapMessage(KeyboardEvent{})
		if res, ok := out.(map[WeylusCommand]KeyboardEvent); !ok {
			t.Errorf("invalid result: %v", out)
		} else if _, ok := res[WeylusCommandKeyboardEvent]; !ok {
			t.Errorf("invalid key in map: %v", res)
		}
	})
	t.Run("Config", func(t *testing.T) {
		out := WrapMessage(Config{})
		if res, ok := out.(map[WeylusCommand]Config); !ok {
			t.Errorf("invalid result: %v", out)
		} else if _, ok := res[WeylusCommandConfig]; !ok {
			t.Errorf("invalid key in map: %v", res)
		}
	})
	t.Run("WeylusCommand", func(t *testing.T) {
		out := WrapMessage(WeylusCommandTryGetFrame)
		if res, ok := out.(WeylusCommand); !ok {
			t.Errorf("invalid result: %v", out)
		} else if res != "TryGetFrame" {
			t.Errorf("invalid value: %v", res)
		}
	})
	t.Run("string", func(t *testing.T) {
		out := WrapMessage("test")
		if res, ok := out.(string); !ok {
			t.Errorf("invalid result: %v", out)
		} else if res != "test" {
			t.Errorf("invalid value: %v", res)
		}
	})
	t.Run("UnderlyingString", func(t *testing.T) {
		out := WrapMessage(utils.UnderlyingString("test"))
		if res, ok := out.(string); !ok {
			t.Errorf("invalid result: %v", out)
		} else if res != "test" {
			t.Errorf("invalid value: %v", res)
		}
	})
}

func TestParseMessage(t *testing.T) {
	t.Run("CapturableList", func(t *testing.T) {
		_, err := ParseMessage([]byte("{\"CapturableList\":[\"Desktop\",\"Monitor: DP-4\",\"Weylus - 0.11.4\",\"Desktop (autopilot)\"]}"))
		if err != nil {
			t.Error(err)
		}
	})
	t.Run("ConfigError", func(t *testing.T) {
		_, err := ParseMessage([]byte("{\"ConfigError\":\"test\"}"))
		if err != nil {
			t.Error(err)
		}
	})
	t.Run("Error", func(t *testing.T) {
		_, err := ParseMessage([]byte("{\"Error\":\"test\"}"))
		if err != nil {
			t.Error(err)
		}
	})
	t.Run("V", func(t *testing.T) {
		_, err := ParseMessage([]byte("{\"Error\":\"test\"}"))
		if err != nil {
			t.Error(err)
		}
	})

}

func FuzzParseMessage(f *testing.F) {
	f.Fuzz(func(t *testing.T, in []byte) {
		out, err := ParseMessage(in)
		if err != nil {
			if !strings.HasPrefix(err.Error(), "failed unmarshaling data") {
				t.Errorf("invalid error: %v", err)
			}
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
