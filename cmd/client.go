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
	"io"
	"io/fs"
	"net"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/OmegaRogue/weylus-desktop/bmp"
	"github.com/OmegaRogue/weylus-desktop/client"
	"github.com/OmegaRogue/weylus-desktop/internal/event"
	"github.com/OmegaRogue/weylus-desktop/protocol"
	"github.com/OmegaRogue/weylus-desktop/utils"
	"github.com/diamondburned/gotk4/pkg/cairo"
	"github.com/diamondburned/gotk4/pkg/gdk/v4"
	"github.com/diamondburned/gotk4/pkg/gio/v2"
	"github.com/diamondburned/gotk4/pkg/glib/v2"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	"github.com/edsrzf/mmap-go"
	"github.com/pkg/errors"
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
	var wg sync.WaitGroup

	tmpdir, err := os.MkdirTemp(fmt.Sprintf("/run/user/%d", os.Getuid()), "giotest-")
	if err != nil {
		log.Fatal().Err(err).Msg("make tempdir")
	}

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
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*10)

	app.ConnectShutdown(func() {
		cancel()
		wg.Wait()

		log.Info().Msgf("cleaning up %s", tmpdir)
		if err := os.RemoveAll(tmpdir); err != nil {
			log.Warn().Err(err).Msg("warning: failed to remove temp dir")
		}
	})

	wg.Add(1)
	go func() {
		defer wg.Done()

		<-ctx.Done()
		app.Quit()
	}()

	errCh := make(chan error, 1)

	weylusClient := client.NewWeylusClient(ctx, 30)

	weylusClient.BufPipe = utils.NewBufPipe()

	wg.Add(1)
	go func() {
		defer wg.Done()

		// Atomic writing is where the magic happens!
		// Also, "-pix_fmt rgba" actually gives us BGRA.
		err := sh(ctx, fmt.Sprintf(`
			ffmpeg -y \
				-hide_banner -loglevel error \
				-f mp4 -re -i - \
				-c:v bmp -pix_fmt rgba -update 1 -atomic_writing 1 %s/screen.bmp
		`, tmpdir), weylusClient.BufPipe)
		// err := sh(ctx, fmt.Sprintf(`
		// 	ffmpeg -y \
		// 		-hide_banner -loglevel error \
		// 		-stream_loop -1 -re -i ~/Videos/LLOGE.mp4 \
		// 		-c:v bmp -update 1 -atomic_writing 1 %s/screen.bmp
		// `, tmpdir))
		if err != nil {
			select {
			case errCh <- err:
			default:
			}
		}
	}()

	bmpr := newBMPReader(filepath.Join(tmpdir, "screen.bmp"), int(weylusClient.Framerate))

	wg.Add(1)
	go func() {
		defer wg.Done()

		err := bmpr.start(ctx)
		if err != nil {
			select {
			case errCh <- err:
			default:
			}
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		select {
		case <-ctx.Done():
			return
		case err := <-errCh:
			log.Err(err).Msg("error occurred")
			cancel()
		}
	}()

	screen := gtk.NewPicture()
	screen.SetKeepAspectRatio(true)
	screen.SetHExpand(true)
	screen.AddTickCallback(func(_ gtk.Widgetter, clock gdk.FrameClocker) bool {
		bmpr.acquire(func(txt *gdk.MemoryTexture) { screen.SetPaintable(txt) })
		return true
	})

	glib.TimeoutAdd(1000/weylusClient.Framerate, func() bool {
		screen.QueueDraw()
		return ctx.Err() == nil
	})

	layout.Attach(screen, 0, 5, 1, 1)

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
		MaxWidth:      5120 / 2,
		MaxHeight:     1440,
		ClientName:    "weylus-desktop",
	}); err != nil {
		log.Err(err).Msg("send Config")
	}
	window.Show()
}

type bmpReader struct {
	path string
	freq time.Duration
	dec  *bmp.BGRADecoder

	bmp  *bmp.NBGRA
	txtv atomic.Value // *gdk.MemoryTexture
}

func newBMPReader(path string, fps int) *bmpReader {
	return &bmpReader{
		path: path,
		freq: time.Second / time.Duration(fps),
		dec:  bmp.NewBGRADecoder(),
	}
}

func (r *bmpReader) start(ctx context.Context) error {
	clock := time.NewTicker(r.freq)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-clock.C:
			// ok
		}

		if err := r.update(ctx); err != nil && !errors.Is(err, fs.ErrNotExist) {
			return err
		}
	}
}

func (r *bmpReader) update(ctx context.Context) error {
	f, err := os.Open(r.path)
	if err != nil {
		return errors.Wrap(err, "failed to open bmp snapshot")
	}
	defer f.Close()

	buf, err := mmap.Map(f, mmap.RDONLY, 0)
	if err != nil {
		return errors.Wrap(err, "failed to mmap bmp snapshot")
	}
	defer buf.Unmap()

	// TODO: figure out double buffering to avoid locking for too long.
	r.bmp, err = r.dec.Decode(buf, r.bmp)
	if err != nil {
		return errors.Wrap(err, "failed to decode bmp snapshot")
	}

	// This unfortunately makes a newly-allocated copy of the image every call.
	// It's probably the slowest part of this code and why you should write your
	// own Paintable.
	newTexture := gdk.NewMemoryTexture(
		r.bmp.Rect.Dx(),
		r.bmp.Rect.Dy(),
		// I'm not sure if this is the format that GTK uses. They might be
		// swizzling this on their own in the code which adds cost.
		gdk.MemoryB8G8R8A8Premultiplied,
		glib.NewBytesWithGo(r.bmp.Pix),
		uint(r.bmp.Stride),
	)

	r.txtv.Store(newTexture)
	return nil
}

func (r *bmpReader) acquire(f func(*gdk.MemoryTexture)) {
	txt, _ := r.txtv.Swap((*gdk.MemoryTexture)(nil)).(*gdk.MemoryTexture)
	if txt != nil {
		f(txt)
	}
}

func sh(ctx context.Context, shcmd string, reader io.Reader) error {
	cmd := exec.CommandContext(ctx, "sh", "-c", shcmd)
	cmd.Stdout = nil
	cmd.Stderr = os.Stderr
	in, err := cmd.StdinPipe()
	if err != nil {
		return err
	}
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				_, err := io.Copy(in, reader)
				if err != nil {
					log.Err(err).Msg("error on copy")
				}
			}
		}
	}()

	return cmd.Run()
}
