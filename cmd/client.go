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
	"context"
	"fmt"
	"net"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/OmegaRogue/weylus-desktop/client"
	"github.com/OmegaRogue/weylus-desktop/internal/event"
	"github.com/OmegaRogue/weylus-desktop/protocol"
	"github.com/diamondburned/gotk4/pkg/cairo"
	"github.com/diamondburned/gotk4/pkg/gio/v2"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// clientCmd represents the client command
var clientCmd = NewClientCmd()

// NewClientCmd creates a new client command
func NewClientCmd() *cobra.Command {
	var clientCmd = &cobra.Command{
		Use:   "client",
		Short: "Start the weylus-desktop client",
		Long:  `Start the weylus-desktop client`,
		Run: func(cmd *cobra.Command, args []string) {
			app := gtk.NewApplication("codes.omegavoid.weylus-desktop", gio.ApplicationHandlesCommandLine)
			app.ConnectCommandLine(func(commandLine *gio.ApplicationCommandLine) (gint int) {
				app.Activate()
				return 0
			})
			app.ConnectActivate(func() { activate(app) })
			if code := app.Run(os.Args); code > 0 {
				os.Exit(code)
			}
		},
	}
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// clientCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	clientCmd.Flags().StringP("hostname", "", "localhost", "Hostname to connect to")
	if err := viper.BindPFlag("hostname", clientCmd.Flags().Lookup("hostname")); err != nil {
		log.Fatal().Err(err).Msg("failed binding flag hostname")
	}
	return clientCmd
}

func init() {
	rootCmd.AddCommand(clientCmd)
}

type State struct {
	X float64
	Y float64
}

func activate(app *gtk.Application) {
	window := gtk.NewApplicationWindow(app)
	window.SetTitle("weylus-client")
	drawArea := gtk.NewDrawingArea()
	drawArea.SetVExpand(true)
	drawArea.SetDrawFunc(func(draw *gtk.DrawingArea, cr *cairo.Context, w, h int) {
		// Draw a red rectangle at the X and Y location.
		cr.SetSourceRGB(255, 0, 0)
		cr.Fill()
	})

	stylusLabel := gtk.NewLabel("")
	stylusLabel.SetHAlign(gtk.AlignStart)
	clickLabel := gtk.NewLabel("")
	clickLabel.SetHAlign(gtk.AlignStart)
	touchLabel := gtk.NewLabel("")
	touchLabel.SetHAlign(gtk.AlignStart)
	keyLabel := gtk.NewLabel("")
	keyLabel.SetHAlign(gtk.AlignStart)
	scrollLabel := gtk.NewLabel("")
	scrollLabel.SetHAlign(gtk.AlignStart)

	layout := gtk.NewGrid()
	layout.Attach(stylusLabel, 0, 0, 1, 1)
	layout.Attach(clickLabel, 0, 1, 1, 1)
	layout.Attach(touchLabel, 0, 2, 1, 1)
	layout.Attach(keyLabel, 0, 3, 1, 1)
	layout.Attach(scrollLabel, 0, 4, 1, 1)

	manager := event.NewControllerManager()
	manager.AddCallback(func(m *event.ControllerManager) {
		stylusLabel.SetMarkup(fmt.Sprintf("<span font_desc=\"mono\">%v</span>", m.StylusState))
		clickLabel.SetMarkup(fmt.Sprintf("<span font_desc=\"mono\">%v</span>", m.MouseState))
		touchLabel.SetMarkup(fmt.Sprintf("<span font_desc=\"mono\">%v</span>", m.TouchState))
		keyLabel.SetMarkup(fmt.Sprintf("<span font_desc=\"mono\">%v</span>", m.KeyState))
		scrollLabel.SetMarkup(fmt.Sprintf("<span font_desc=\"mono\">%v</span>", m.ScrollState))
	})

	overlay := gtk.NewOverlay()
	overlay.SetVExpand(true)
	overlay.SetHExpand(true)

	manager.ConnectControllers(overlay)
	window.AddController(manager.Key)
	window.AddController(manager.Scroll)
	overlay.SetChild(layout)
	overlay.AddOverlay(drawArea)
	window.SetChild(overlay)
	window.SetDefaultSize(400, 300)
	ctx, _ := context.WithTimeout(context.Background(), time.Minute*10)
	weylusClient := client.NewWeylusClient(ctx, 30)

	address := url.URL{
		Scheme: "ws",
		Host:   net.JoinHostPort(viper.GetString("hostname"), strconv.FormatUint(uint64(viper.GetUint16("websocket-port")), 10)),
	}

	if err := weylusClient.Dial(address.String()); err != nil {
		log.Err(err).Msg("dial weylusClient")
	}
	go weylusClient.Listen()
	go weylusClient.Run()
	go weylusClient.RunVideo()

	capturables, err := weylusClient.GetCapturableList()
	if err != nil {
		log.Fatal().Err(err).Msg("get capturables")
	} else {
		log.Debug().Strs("capturables", capturables.CapturableList).Msg("get capturables")
	}

	if _, err := weylusClient.Config(protocol.Config{
		UInputSupport: true,
		CapturableID:  1,
		CaptureCursor: true,
		MaxWidth:      640,
		MaxHeight:     480,
		ClientName:    "weylus-desktop",
	}); err != nil {
		log.Err(err).Msg("send Config")
	}
	window.Show()
}
