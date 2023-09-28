package image_jpeg

import "image/jpeg"

type Options struct {
	*jpeg.Options
	Progressive bool // Render in progressive mode, only available with libturbojpeg.
}
