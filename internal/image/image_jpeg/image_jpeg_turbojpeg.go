//go:build pdfium_use_turbojpeg

package image_jpeg

import (
	"bufio"
	"image"
	"image/jpeg"
	"io"

	"github.com/bmharper/turbo"
)

func Encode(w io.Writer, m *image.RGBA, o *jpeg.Options) error {
	imageWriter := bufio.NewWriter(w)

	// Clip quality to [1, 100].
	quality := jpeg.DefaultQuality
	if o != nil {
		quality = o.Quality
		if quality < 1 {
			quality = 1
		} else if quality > 100 {
			quality = 100
		}
	}

	dimensions := m.Bounds().Size()

	raw := turbo.Image{
		Width:  dimensions.X,
		Height: dimensions.Y,
		Stride: m.Stride,
		Pixels: m.Pix,
	}

	params := turbo.MakeCompressParams(turbo.PixelFormatRGBA, turbo.Sampling420, quality, 0)
	jpg, err := turbo.Compress(&raw, params)
	if err != nil {
		return err
	}

	_, err = imageWriter.Write(jpg)
	if err != nil {
		return err
	}

	err = imageWriter.Flush()
	if err != nil {
		return err
	}

	return nil
}
