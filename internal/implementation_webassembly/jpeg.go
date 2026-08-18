package implementation_webassembly

import (
	"errors"
	"image"
	"io"

	"github.com/klippa-app/go-pdfium/internal/image/image_jpeg"
)

// Input pixel formats of pdfium_jpeg_encode.
const (
	jpegEncodeFormatRGB  = 0
	jpegEncodeFormatRGBA = 1
	jpegEncodeFormatBGRA = 2
	jpegEncodeFormatGray = 3
)

// encodeJPEG encodes m to JPEG (baseline, or progressive when
// opt.Progressive is set). When the loaded pdfium module exports
// libjpeg-turbo's encoder (pdfium_jpeg_encode) and the image is a type it
// can consume directly, the encode runs inside the guest (with the SIMD
// kernels); otherwise it falls back to image_jpeg.Encode.
func (p *PdfiumImplementation) encodeJPEG(w io.Writer, m image.Image, opt image_jpeg.Options) error {
	encode := p.Fn("pdfium_jpeg_encode")
	// Guard against custom wasm binaries with an older/newer shim signature:
	// (data, width, height, stride, format, quality, progressive, out_buf,
	// out_size).
	if encode == nil || len(encode.Definition().ParamTypes()) != 9 {
		return image_jpeg.Encode(w, m, opt)
	}

	var pixels []byte
	var stride, format int
	switch img := m.(type) {
	case *image.RGBA:
		pixels = img.Pix
		stride = img.Stride
		format = jpegEncodeFormatRGBA
	case *image.Gray:
		pixels = img.Pix
		stride = img.Stride
		format = jpegEncodeFormatGray
	default:
		return image_jpeg.Encode(w, m, opt)
	}

	dimensions := m.Bounds().Size()
	if dimensions.X == 0 || dimensions.Y == 0 || len(pixels) == 0 {
		return image_jpeg.Encode(w, m, opt)
	}

	quality := 75
	if opt.Options != nil {
		quality = opt.Options.Quality
		if quality < 1 {
			quality = 1
		} else if quality > 100 {
			quality = 100
		}
	}

	inPtr, err := p.MallocNoZero(uint64(len(pixels)))
	if err != nil {
		return err
	}
	defer p.Free(inPtr)

	// Two output parameters: the buffer pointer and its size.
	outParams, err := p.Malloc(16)
	if err != nil {
		return err
	}
	defer p.Free(outParams)

	if !p.Module.Memory().Write(uint32(inPtr), pixels) {
		return errors.New("could not write pixel data to guest memory")
	}

	progressive := uint64(0)
	if opt.Progressive {
		progressive = 1
	}

	res, err := p.callFn(encode, inPtr, uint64(dimensions.X),
		uint64(dimensions.Y), uint64(stride), uint64(format),
		uint64(quality), progressive, outParams, outParams+8)
	if err != nil {
		return err
	}
	if res[0] != 1 {
		return errors.New("in-module JPEG encode failed")
	}

	bufPtr, ok := p.Module.Memory().ReadUint32Le(uint32(outParams))
	if !ok {
		return errors.New("could not read JPEG output pointer")
	}
	defer p.call("pdfium_jpeg_free", uint64(bufPtr))

	size, ok := p.Module.Memory().ReadUint32Le(uint32(outParams + 8))
	if !ok {
		return errors.New("could not read JPEG output size")
	}

	jpegData, ok := p.Module.Memory().Read(bufPtr, size)
	if !ok {
		return errors.New("could not read JPEG output data")
	}

	_, err = w.Write(jpegData)
	return err
}
