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

//go:generate go run . ../../../man ../../../docs
package main

import (
	"os"
	"strings"

	cmd2 "github.com/OmegaRogue/weylus-desktop/cmd"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra/doc"
)

func main() {
	cmd := cmd2.GetRootCmd()
	header := &doc.GenManHeader{
		Title:   "MINE",
		Section: "3",
	}
	manDir := "man"
	if len(os.Args) > 1 {
		manDir = os.Args[1]
	}
	mdDir := "docs"
	if len(os.Args) > 2 {
		mdDir = os.Args[2]
	}

	if err := doc.GenManTree(cmd, header, manDir); err != nil {
		log.Fatal().Err(err).Msg("failed to generate mandoc")
	}

	identity := func(s string) string { return strings.TrimSuffix(s, ".md") }
	emptyStr := func(s string) string { return "" }
	if err := doc.GenMarkdownTreeCustom(cmd, mdDir, emptyStr, identity); err != nil {
		log.Fatal().Err(err).Msg("failed to generate md docs")
	}
}
