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

package bmp

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"io"

	"github.com/pkg/errors"
)

// ErrUnsupported means that the input BMP image uses a valid but unsupported
// feature.
var ErrUnsupported = errors.New("bmp: unsupported BMP image")

func readUint16(b []byte) uint16 {
	return uint16(b[0]) | uint16(b[1])<<8
}

func readUint32(b []byte) uint32 {
	return uint32(b[0]) | uint32(b[1])<<8 | uint32(b[2])<<16 | uint32(b[3])<<24
}

// NBGRA holds an alpha-premultiplied 32-bit BGRA image.
type NBGRA struct {
	Pix    []uint8
	Stride int
	Rect   image.Rectangle
}

// NewNBGRA returns a new NBGRA image with the given bounds.
func NewNBGRA(r image.Rectangle) *NBGRA {
	w, h := r.Dx(), r.Dy()
	buf := make([]uint8, 4*w*h)
	return &NBGRA{buf, 4 * w, r}
}

// BGRADecoder is used to decode an BGRA BMP image of the same size as the one
// used for construction.
type BGRADecoder struct {
	cfg *Config
	hdr int64
}

// NewBGRADecoder returns a new BGRADecoder.
func NewBGRADecoder() *BGRADecoder {
	return &BGRADecoder{}
}

// Decode reads an BGRA BMP image from r and write it to dst (if not nil). If
// dst is nil, then a new one is allocated. If topDown is false, the image rows
// will be read bottom-up.
func (dec *BGRADecoder) Decode(src []byte, dst *NBGRA) (*NBGRA, error) {
	var c Config

	if dec.cfg != nil {
		c = *dec.cfg
		src = src[dec.hdr:]
	} else {
		buf := bytes.NewBuffer(src)

		var err error
		c, err = DecodeConfig(buf)
		if err != nil {
			return nil, fmt.Errorf("bmp: failed to decode config: %w", err)
		}

		dec.cfg = &c
		dec.hdr = int64(len(src) - buf.Len())
		src = buf.Bytes() // swap to unread bytes
	}

	if dst != nil {
		if c.Width != dst.Rect.Dx() || c.Height != dst.Rect.Dy() {
			return nil, fmt.Errorf(
				"bmp: image size mismatch: %dx%d != %dx%d",
				c.Width, c.Height, dst.Rect.Dx(), dst.Rect.Dy())
		}
	} else {
		dst = NewNBGRA(image.Rect(0, 0, c.Width, c.Height))
	}

	if c.TopDown {
		// This is rarely the case.
		copy(dst.Pix, src)
	} else {
		y1 := 0
		y2 := c.Height - 1
		for y1 < c.Height {
			b := src[y1*dst.Stride : y1*dst.Stride+c.Width*4]
			p := dst.Pix[y2*dst.Stride : y2*dst.Stride+c.Width*4]
			copy(p, b)
			y1++
			y2--
		}
	}

	return dst, nil
}

// Config extends image.Config to add BMP image config.
type Config struct {
	image.Config
	BPP        int
	TopDown    bool
	AllowAlpha bool
}

// DecodeConfig returns the color model and dimensions of a BMP image without
// decoding the entire image.
// Limitation: The file must be 8, 24 or 32 bits per pixel.
func DecodeConfig(r io.Reader) (Config, error) {
	config, bpp, topDown, allowAlpha, err := decodeConfig(r)
	return Config{config, bpp, topDown, allowAlpha}, err
}

func decodeConfig(r io.Reader) (config image.Config, bitsPerPixel int, topDown bool, allowAlpha bool, err error) {
	// We only support those BMP images with one of the following DIB headers:
	// - BITMAPINFOHEADER (40 bytes)
	// - BITMAPV4HEADER (108 bytes)
	// - BITMAPV5HEADER (124 bytes)
	const (
		fileHeaderLen   = 14
		infoHeaderLen   = 40
		v4InfoHeaderLen = 108
		v5InfoHeaderLen = 124
	)
	var b [1024]byte
	if _, err := io.ReadFull(r, b[:fileHeaderLen+4]); err != nil {
		if err == io.EOF {
			err = io.ErrUnexpectedEOF
		}
		return image.Config{}, 0, false, false, err
	}
	if string(b[:2]) != "BM" {
		return image.Config{}, 0, false, false, errors.New("bmp: invalid format")
	}
	offset := readUint32(b[10:14])
	infoLen := readUint32(b[14:18])
	if infoLen != infoHeaderLen && infoLen != v4InfoHeaderLen && infoLen != v5InfoHeaderLen {
		return image.Config{}, 0, false, false, ErrUnsupported
	}
	if _, err := io.ReadFull(r, b[fileHeaderLen+4:fileHeaderLen+infoLen]); err != nil {
		if err == io.EOF {
			err = io.ErrUnexpectedEOF
		}
		return image.Config{}, 0, false, false, err
	}
	width := int(int32(readUint32(b[18:22])))
	height := int(int32(readUint32(b[22:26])))
	if height < 0 {
		height, topDown = -height, true
	}
	if width < 0 || height < 0 {
		return image.Config{}, 0, false, false, ErrUnsupported
	}
	// We only support 1 plane and 8, 24 or 32 bits per pixel and no
	// compression.
	planes, bpp, compression := readUint16(b[26:28]), readUint16(b[28:30]), readUint32(b[30:34])
	// if compression is set to BI_BITFIELDS, but the bitmask is set to the default bitmask
	// that would be used if compression was set to 0, we can continue as if compression was 0
	if compression == 3 && infoLen > infoHeaderLen &&
		readUint32(b[54:58]) == 0xff0000 && readUint32(b[58:62]) == 0xff00 &&
		readUint32(b[62:66]) == 0xff && readUint32(b[66:70]) == 0xff000000 {
		compression = 0
	}
	if planes != 1 || compression != 0 {
		return image.Config{}, 0, false, false, ErrUnsupported
	}
	switch bpp {
	case 8:
		if offset != fileHeaderLen+infoLen+256*4 {
			return image.Config{}, 0, false, false, ErrUnsupported
		}
		_, err = io.ReadFull(r, b[:256*4])
		if err != nil {
			return image.Config{}, 0, false, false, err
		}
		pcm := make(color.Palette, 256)
		for i := range pcm {
			// BMP images are stored in BGR order rather than RGB order.
			// Every 4th byte is padding.
			pcm[i] = color.RGBA{b[4*i+2], b[4*i+1], b[4*i+0], 0xFF}
		}
		return image.Config{ColorModel: pcm, Width: width, Height: height}, 8, topDown, false, nil
	case 24:
		if offset != fileHeaderLen+infoLen {
			return image.Config{}, 0, false, false, ErrUnsupported
		}
		return image.Config{ColorModel: color.RGBAModel, Width: width, Height: height}, 24, topDown, false, nil
	case 32:
		if offset != fileHeaderLen+infoLen {
			return image.Config{}, 0, false, false, ErrUnsupported
		}
		// 32 bits per pixel is possibly RGBX (X is padding) or RGBA (A is
		// alpha transparency). However, for BMP images, "Alpha is a
		// poorly-documented and inconsistently-used feature" says
		// https://source.chromium.org/chromium/chromium/src/+/bc0a792d7ebc587190d1a62ccddba10abeea274b:third_party/blink/renderer/platform/image-decoders/bmp/bmp_image_reader.cc;l=621
		//
		// That goes on to say "BITMAPV3HEADER+ have an alpha bitmask in the
		// info header... so we respect it at all times... [For earlier
		// (smaller) headers we] ignore alpha in Windows V3 BMPs except inside
		// ICO files".
		//
		// "Ignore" means to always set alpha to 0xFF (fully opaque):
		// https://source.chromium.org/chromium/chromium/src/+/bc0a792d7ebc587190d1a62ccddba10abeea274b:third_party/blink/renderer/platform/image-decoders/bmp/bmp_image_reader.h;l=272
		//
		// Confusingly, "Windows V3" does not correspond to BITMAPV3HEADER, but
		// instead corresponds to the earlier (smaller) BITMAPINFOHEADER:
		// https://source.chromium.org/chromium/chromium/src/+/bc0a792d7ebc587190d1a62ccddba10abeea274b:third_party/blink/renderer/platform/image-decoders/bmp/bmp_image_reader.cc;l=258
		//
		// This Go package does not support ICO files and the (infoLen >
		// infoHeaderLen) condition distinguishes BITMAPINFOHEADER (40 bytes)
		// vs later (larger) headers.
		allowAlpha = infoLen > infoHeaderLen
		return image.Config{ColorModel: color.RGBAModel, Width: width, Height: height}, 32, topDown, allowAlpha, nil
	}
	return image.Config{}, 0, false, false, ErrUnsupported
}
