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
