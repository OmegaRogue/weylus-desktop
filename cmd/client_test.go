/*
 * Copyright © 2023 omegarogue
 * SPDX-License-Identifier: GPL-3.0-or-later
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program. If not, see <http://www.gnu.org/licenses/>.
 */

package cmd

import (
	"testing"

	"github.com/OmegaRogue/weylus-desktop/internal/logger"
)

func TestActivate(t *testing.T) {
	logger.SetupLogger()
	cmd := GetRootCmd()
	cmd.SetArgs([]string{"client"})
	err := cmd.Execute()
	if err != nil {
		t.Fatalf("error: %v", err)
	}
}
