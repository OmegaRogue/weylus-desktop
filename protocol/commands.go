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
	"github.com/OmegaRogue/weylus-desktop/utils"
	"github.com/rs/zerolog/log"
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

func CommandFromOutboundContent[T MessageOutboundContent](content T) WeylusCommand {
	switch any(content).(type) {
	case PointerEvent:
		return WeylusCommandPointerEvent
	case WheelEvent:
		return WeylusCommandWheelEvent
	case KeyboardEvent:
		return WeylusCommandKeyboardEvent
	case Config:
		return WeylusCommandConfig
	default:
		str, err := utils.GetUnderlyingString(content)
		if err != nil {
			log.Panic().Err(err).Msg("what the fuck did you do? (I'm genuinely curious) This should never happen, immediately report this as an issue")
		}
		return WeylusCommand(str)
	}
}
