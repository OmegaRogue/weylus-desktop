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

// Package protocol
package protocol

type ButtonFlags byte

const (
	// ButtonNone is a ButtonFlags of type None.
	ButtonNone ButtonFlags = 0
	// ButtonPrimary is a ButtonFlags of type Primary. Usually the left button
	ButtonPrimary ButtonFlags = 1 << (iota - 1)
	// ButtonSecondary is a ButtonFlags of type Secondary. Usually the right button
	ButtonSecondary
	// ButtonAuxiliary is a ButtonFlags of type Auxiliary. Usually the wheel button or the middle button (if present)
	ButtonAuxiliary
	// ButtonFourth is a ButtonFlags of type Fourth. Typically the Browser Back button
	ButtonFourth
	// ButtonFifth is a ButtonFlags of type Fifth. Typically the Browser Forward button
	ButtonFifth
	// ButtonEraser is a ButtonFlags of type Eraser.
	ButtonEraser
)
