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
	"testing"

	"github.com/OmegaRogue/weylus-desktop/utils"
)

//nolint:funlen,gocognit
func TestResponseFromOutboundContent(t *testing.T) {
	var tests = []struct {
		name  string
		input any
		want  WeylusResponse
	}{
		{"CapturableList", WeylusCommandGetCapturableList, WeylusResponseCapturableList},
		{"CapturableList", WeylusCommandGetCapturableList.String(), WeylusResponseCapturableList},
		{"ConfigOk", Config{}, WeylusResponseConfigOk},
		{"Error", "", WeylusResponseError},
		{"KeyboardEvent", KeyboardEvent{}, ""},
		{"PointerEvent", PointerEvent{}, ""},
		{"WheelEvent", WheelEvent{}, ""},
	}
	//nolint:dupl
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var res WeylusResponse
			switch val := tt.input.(type) {
			case PointerEvent:
				res1 := ResponseFromOutboundContent(val)
				res = res1
			case WheelEvent:
				res1 := ResponseFromOutboundContent(val)
				res = res1
			case KeyboardEvent:
				res1 := ResponseFromOutboundContent(val)
				res = res1
			case Config:
				res1 := ResponseFromOutboundContent(val)
				res = res1
			case string:
				res1 := ResponseFromOutboundContent(val)
				res = res1
			default:
				v, err := utils.GetUnderlyingString(tt.input)
				if err != nil {
					t.Error(err)
				}
				res1 := ResponseFromOutboundContent(utils.UnderlyingString(v))
				res = res1
			}
			if res != tt.want {
				t.Errorf("got %s, want %s", res, tt.want)
			}
		})
	}
}
