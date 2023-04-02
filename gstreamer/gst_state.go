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

//go:generate go-enum --names --values
package gstreamer

// GstState represents the possible states an element can be in.
// ENUM(VoidPending,Null,Ready,Paused,Playing)
type GstState int

// GstStateChangeReturn represents the possible return values from a state change function. Only Failure is a real failure.
// ENUM(Failure,Success,Async,NoPreRoll)
type GstStateChangeReturn int
