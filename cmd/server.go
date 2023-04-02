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

package cmd

import (
	"fmt"
	"net"
	"os"

	"github.com/OmegaRogue/weylus-desktop/web"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// serverCmd represents the client command
var serverCmd = NewServerCmd()

// NewServerCmd creates a new server command
func NewServerCmd() *cobra.Command {
	var serverCmd = &cobra.Command{
		Use:   "server",
		Short: "Start the weylus-desktop server",
		Long:  `Start the weylus-desktop server`,
		Run:   runServerCommand,
	}
	serverCmd.Flags().BoolP("auto-start", "", false, "Start Weylus server immediately on program start.")
	serverCmd.Flags().BoolP("no-gui", "", false, "Run Weylus without gui and start immediately.")
	serverCmd.Flags().BoolP("print-access-html", "", false, "Print access.html served by Weylus.")
	serverCmd.Flags().BoolP("print-index-html", "", false, "Print template of index.html served by Weylus.")
	serverCmd.Flags().BoolP("print-lib-js", "", false, "Print lib.js served by Weylus.")
	serverCmd.Flags().BoolP("print-style-css", "", false, "Print style.css served by Weylus.")
	serverCmd.Flags().IPP("bind-address", "", net.IPv4(0, 0, 0, 0), "Bind address")
	serverCmd.Flags().StringP("custom-access-html", "", "", "Use custom access.html to be served by Weylus.")
	serverCmd.Flags().Uint16P("websocket-port", "", 9001, "Websocket port")
	serverCmd.Flags().StringP("access-code", "", "", "Access code")
	serverCmd.Flags().StringP("custom-index-html", "", "", "Use custom template of index.html to be served by Weylus.")
	serverCmd.Flags().StringP("custom-lib-js", "", "", "Use custom lib.js to be served by Weylus.")
	serverCmd.Flags().StringP("custom-style-css", "", "", "Use custom style.css to be served by Weylus.")
	serverCmd.Flags().Uint16P("web-port", "", 1701, "Web port")

	if err := serverCmd.MarkFlagFilename("custom-access-html", "html"); err != nil {
		log.Fatal().Err(err).Msg("failed mark flag custom-access-html as filename")
	}
	if err := serverCmd.MarkFlagFilename("custom-index-html", "html", "gohtml"); err != nil {
		log.Fatal().Err(err).Msg("failed mark flag custom-index-html as filename")
	}
	if err := serverCmd.MarkFlagFilename("custom-lib-js", "js"); err != nil {
		log.Fatal().Err(err).Msg("failed mark flag custom-lib-js as filename")
	}
	if err := serverCmd.MarkFlagFilename("custom-style-css", "css"); err != nil {
		log.Fatal().Err(err).Msg("failed mark flag custom-style-css as filename")
	}
	serverFlagsOSSpecific(serverCmd)
	serverCmd.Flags().VisitAll(func(flag *pflag.Flag) {
		if err := viper.BindPFlag(flag.Name, flag); err != nil {
			log.Fatal().Err(err).Msgf("failed binding flag %s", flag.Name)
		}
	})

	return serverCmd
}

func runServerCommand(_ *cobra.Command, _ []string) {
	if accessPath := viper.GetString("custom-access-html"); accessPath != "" {
		data, err := os.ReadFile(accessPath)
		if err != nil {
			log.Fatal().Err(err).Str("path", accessPath).Msg("failed reading file")
		}
		web.AccessHTML = string(data)
	}
	if indexPath := viper.GetString("custom-index-html"); indexPath != "" {
		data, err := os.ReadFile(indexPath)
		if err != nil {
			log.Fatal().Err(err).Str("path", indexPath).Msg("failed reading file")
		}
		web.IndexHTML = string(data)
	}
	if libPath := viper.GetString("custom-lib-js"); libPath != "" {
		data, err := os.ReadFile(libPath)
		if err != nil {
			log.Fatal().Err(err).Str("path", libPath).Msg("failed reading file")
		}
		web.LibJS = string(data)
	}
	if stylePath := viper.GetString("custom-style-css"); stylePath != "" {
		data, err := os.ReadFile(stylePath)
		if err != nil {
			log.Fatal().Err(err).Str("path", stylePath).Msg("failed reading file")
		}
		web.StyleCSS = string(data)
	}
	switch {
	case viper.GetBool("print-access-html"):
		fmt.Println(web.AccessHTML)
		return
	case viper.GetBool("print-index-html"):
		fmt.Println(web.IndexHTML)
		return
	case viper.GetBool("print-lib-js"):
		fmt.Println(web.LibJS)
		return
	case viper.GetBool("print-style-css"):
		fmt.Println(web.StyleCSS)
		return
	}
	// TODO serverLogger := log.With().Str("component", "server").Logger()
	// TODO server.NewWeylusServer()
	// server.WeylusWeb(serverLogger)
}

func init() {
	rootCmd.AddCommand(serverCmd)
}
