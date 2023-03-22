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
	"os"
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

var (
	cfgFile string
	rootCmd = NewRootCmd()
)

// NewRootCmd returns a new command which
// represents the base command when called without any subcommands
// TODO add description
func NewRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "weylus-desktop",
		Short: "A brief description of your application",
		Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
		Run: func(cmd *cobra.Command, args []string) {
			app := gtk.NewApplication("codes.omegavoid.weylus-client", gio.ApplicationFlagsNone)
			app.ConnectActivate(func() { activate(app) })
			if code := app.Run(os.Args); code > 0 {
				os.Exit(code)
			}
		},
	}

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.weylus-desktop.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	return rootCmd
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

	vfs := gio.VFSGetLocal()
	file := vfs.FileForURI("rtmp://localhost:1935/live/app")
	log.Info().Msg(file.URI())
	video := gtk.NewVideo()
	video.SetFile(file)
	video.SetVExpand(true)
	video.SetHExpand(true)
	go func() {
		time.Sleep(time.Second * 20)
		video.SetAutoplay(true)
		log.Info().Msg("start")
	}()

	layout := gtk.NewGrid()
	layout.Attach(stylusLabel, 0, 0, 1, 1)
	layout.Attach(clickLabel, 0, 1, 1, 1)
	layout.Attach(touchLabel, 0, 2, 1, 1)
	layout.Attach(keyLabel, 0, 3, 1, 1)
	layout.Attach(scrollLabel, 0, 4, 1, 1)
	layout.Attach(video, 0, 5, 1, 1)

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

	if err := weylusClient.Dial("ws://192.168.0.49:9001"); err != nil {
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

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		log.Fatal().Err(err).Msg("exit with error")
	}
}

func init() {
	cobra.OnInitialize(initConfig)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".weylus-desktop" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".weylus-desktop")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		_, err := fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
		if err != nil {
			log.Err(err).Msg("print config file location")
		}
	}
}
