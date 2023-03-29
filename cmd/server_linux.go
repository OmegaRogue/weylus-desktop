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
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func serverFlagsOSSpecific(cmd *cobra.Command) {
	cmd.Flags().BoolP("wayland-support", "", false, "Wayland/PipeWire Support.")
	cmd.Flags().BoolP("try-vaapi", "", false, "Try to use hardware acceleration through the Video Acceleration API.")
	cmd.Flags().BoolP("try-nvenc", "", false, "Try to use Nvidia's NVENC to encode the video via GPU.")
	if err := viper.BindPFlag("try-nvenc", cmd.Flags().Lookup("try-nvenc")); err != nil {
		log.Fatal().Err(err).Msg("failed binding flag try-nvenc")
	}
	if err := viper.BindPFlag("try-vaapi", cmd.Flags().Lookup("try-vaapi")); err != nil {
		log.Fatal().Err(err).Msg("failed binding flag try-vaapi")
	}
	if err := viper.BindPFlag("wayland-support", cmd.Flags().Lookup("wayland-support")); err != nil {
		log.Fatal().Err(err).Msg("failed binding flag wayland-support")
	}
}
