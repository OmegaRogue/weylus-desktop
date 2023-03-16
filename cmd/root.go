/*
Copyright © 2023 OmegaRogue

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/

package cmd

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/diamondburned/gotk4/pkg/cairo"
	"github.com/diamondburned/gotk4/pkg/gio/v2"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
	"weylus-surface/internal/event"
	"weylus-surface/protocol"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "weylus-surface",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {

		data, _ := json.Marshal(protocol.WeylusError{ErrorMessage: "test"})
		log.Info().RawJSON("data", data).Msg("")
		app := gtk.NewApplication("codes.omegavoid.weylus-client", gio.ApplicationFlagsNone)
		app.ConnectActivate(func() { activate(app) })

		if code := app.Run(os.Args); code > 0 {
			os.Exit(code)
		}

	},
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
		// Draw a red rectagle at the X and Y location.
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
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*10)
	go func() {

		c, _, err := websocket.Dial(ctx, "ws://192.168.0.49:9001", nil)
		if err != nil {
			log.Err(err).Msg("dial websocket")
		}

		if err := wsjson.Write(ctx, c, "GetCapturableList"); err != nil {
			log.Err(err).Msg("send GetCapturableList")
		}
		msg, err := readSocket(ctx, c)
		log.Err(err).RawJSON("msg", []byte(msg)).Msg("read socket")

		config := protocol.WrapMessage(protocol.Config{
			UInputSupport: true,
			CapturableID:  0,
			CaptureCursor: false,
			MaxWidth:      1920,
			MaxHeight:     1080,
			ClientName:    "weylus-surface",
		})
		if err := wsjson.Write(ctx, c, config); err != nil {
			log.Err(err).Msg("send Config")
		}
		msg, err = readSocket(ctx, c)
		log.Err(err).Stack().RawJSON("msg", []byte(msg)).Msg("read socket")

		for i := 0; i < 3; i++ {
			eventMessage := protocol.WrapMessage(protocol.WheelEvent{
				Timestamp: uint64(11968000), Dy: 60,
			})

			if err := wsjson.Write(ctx, c, eventMessage); err != nil {
				log.Err(err).Msg("send WheelEvent")
			}

			for ctx.Err() == nil {
				msg, err := readSocket(ctx, c)
				log.Err(err).RawJSON("msg", []byte(msg)).Msg("read socket")
			}
		}

		if err := c.Close(websocket.StatusNormalClosure, ""); err != nil {
			log.Err(err).Msg("close websocket")
		}
		defer cancel()

	}()

	//window.Show()
}

func readSocket(ctx context.Context, c *websocket.Conn) (string, error) {
	_, data, err := c.Read(ctx)
	if err != nil {
		return "", errors.Wrap(err, "read websocket")
	}
	text := strings.TrimSpace(string(data))
	if text != "" {
		if strings.Contains(text, "ConfigError") {
			var err2 protocol.WeylusConfigError
			err := json.Unmarshal(data, &err2)
			if err != nil {
				return "", errors.New(fmt.Sprintf("failed unmarshaling error: %s", text))
			}
			return "", &err2
		} else if strings.Contains(text, "Error") {
			var err2 protocol.WeylusError
			err := json.Unmarshal(data, &err2)
			if err != nil {
				return "", errors.New(fmt.Sprintf("failed unmarshaling error: %s", text))
			}
			return "", &err2
		}
		return text, nil
	}
	log.Warn().Msg("received empty")
	return "", nil
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout}).With().Caller().Stack().Logger()
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.weylus-surface.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
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

		// Search config in home directory with name ".weylus-surface" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".weylus-surface")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
