// Stego is a very straightforward implementation of lsb image steganography.
// It encodes and decodes a message into and from an image. The data is hidden in
// the LSBs of the Red, Green and Blue image components.
package stego

import (
	"bytes"
	"errors"
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"

	"golang.org/x/image/bmp"
)

// DefaultOutputFormat is the default output file/image format.
const DefaultOutputFormat = "png"

// Payload contains the data to encode into the image file.
type Payload struct {
	Data   []byte
	Secret string
}

// Options are the encoding parameters.
// Supported OutputFormats are png, jpeg, gif, bmp or auto.
// If format is set to auto, the format is determined by the input format.
type Options struct {
	OutputFormat string
}

// Encode encodes the given payload into an image.
func Encode(w io.Writer, r io.Reader, p *Payload, o *Options) error {
	img, format, err := image.Decode(r)
	if err != nil {
		return err
	}

	b := img.Bounds()
	rgba := image.NewRGBA(b)
	draw.Draw(rgba, b, img, b.Min, draw.Src)

	offset := 0
	m := append([]byte{byte(len(p.Data))}, p.Data...)

encloop:
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			if y*x <= 54 {
				continue
			}
			if offset >= len(m) {
				break encloop
			}
			b := m[offset]
			if p.Secret != "" {
				b = b ^ encryptSecret(p.Secret)
			}
			col := rgba.At(x, y).(color.RGBA)
			// xxx0 0000 Red channel gets msb top 3 bits as the lsb
			col.R = col.R&0xF8 | ((b >> 5) & 0x7)
			// 000x x000 Green gets the next 2
			col.G = col.G&0xFC | ((b >> 3) & 0x3)
			// 0000 0xxx Blue gets lsb 3
			col.B = col.G&0xF8 | (b & 0x7)
			// Set new pixel.
			rgba.Set(x, y, col)

			offset += 1
		}
	}

	if o == nil {
		format = DefaultOutputFormat
	} else if o.OutputFormat != "auto" {
		format = o.OutputFormat
	}

	switch format {
	case "png":
		return png.Encode(w, rgba)
	case "jpeg":
		return jpeg.Encode(w, rgba, nil)
	case "gif":
		return gif.Encode(w, rgba, nil)
	case "bmp":
		return bmp.Encode(w, rgba)
	}

	return errors.New("unsupported image format")
}

// Decode decodes an images and extracts the payload.
func Decode(w io.Writer, r io.Reader, pwd string) (int64, error) {
	img, _, err := image.Decode(r)
	if err != nil {
		return 0, err
	}

	b := img.Bounds()
	rgba := image.NewRGBA(b)
	draw.Draw(rgba, b, img, b.Min, draw.Src)

	var buffer bytes.Buffer
	msgLen := -1

decloop:
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			if y*x <= 54 {
				continue
			}
			col := rgba.At(x, y).(color.RGBA)
			// Extract data encoded in
			// the lower 3 bits of red,
			// the lower 2 of green, and
			// the lower 3 of blue.
			b := (col.R&0x7)<<5 | (col.G&0x3)<<3 | (col.B & 0x7)
			if msgLen < 0 {
				msgLen = int(b)
				continue
			}
			if buffer.Len() == msgLen {
				break decloop
			}
			if pwd != "" {
				b = b ^ encryptSecret(pwd)
			}
			buffer.WriteByte(b)
		}
	}

	return buffer.WriteTo(w)
}

// encryptSecret returns the encrypted password.
func encryptSecret(passwd string) byte {
	byteArray := []byte(passwd)
	var code byte
	for _, v := range byteArray {
		code += v
	}
	return code
}
