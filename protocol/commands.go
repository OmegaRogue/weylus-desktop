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

//go:generate go-enum --marshal --names --values
package protocol

import (
	"reflect"

	"github.com/pkg/errors"
)

// WeylusCommand contains the possible commands supported by weylus
/*
ENUM(
TryGetFrame
GetCapturableList
Config
KeyboardEvent
PointerEvent
WheelEvent
)
*/
type WeylusCommand string

// WeylusResponse contains the possible commands supported by weylus
/*
ENUM(
NewVideo
CapturableList
ConfigOk
ConfigError
Error
)
*/
type WeylusResponse string

var commandResponse = map[WeylusCommand]WeylusResponse{
	WeylusCommandGetCapturableList: WeylusResponseCapturableList,
	WeylusCommandConfig:            WeylusResponseConfigOk,
	"":                             WeylusResponseError,
}

func CommandFromOutboundContent[T MessageOutboundContent](content T) (WeylusCommand, error) {
	switch any(content).(type) {
	case PointerEvent:
		return WeylusCommandPointerEvent, nil
	case WheelEvent:
		return WeylusCommandWheelEvent, nil
	case KeyboardEvent:
		return WeylusCommandKeyboardEvent, nil
	case Config:
		return WeylusCommandConfig, nil
	default:
		if b := reflect.ValueOf(content); b.Kind() == reflect.String {
			return WeylusCommand(b.String()), nil
		} else {
			return "", errors.New("Invalid Outbound Content")
		}
	}
}

func ResponseFromOutboundContent[T MessageOutboundContent](content T) (WeylusResponse, error) {
	cmd, err := CommandFromOutboundContent(content)
	if err != nil {
		return "", errors.Wrap(err, "Can't get Response")
	}
	return commandResponse[cmd], nil
}
