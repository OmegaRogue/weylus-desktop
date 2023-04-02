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

func ResponseFromOutboundContent[T MessageOutboundContent](content T) WeylusResponse {
	cmd := CommandFromOutboundContent(content)
	return commandResponse[cmd]
}
