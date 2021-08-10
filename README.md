# `go-pdfium` A package which allows for multithreaded pdf rendering

## Created by [Klippa](https://www.klippa.com/)

# Download C library

Download a pdfium binary here: https://github.com/bblanchon/pdfium-binaries/releases
Extract it somewhere

# Configure pkg-config

Create/edit file `/usr/lib/pkgconfig/pdfium.pc`

```
prefix={path}
libdir={path}/lib
includedir={path}/include

Name: pdfium
Description: pdfium
Version: 4320
Requires:

Libs: -L${libdir} -lpdfium
Cflags: -I${includedir}
```

Replace `{path}` with the path you extracted pdfium in.

Make sure you extend your library path when running:

`export LD_LIBRARY_PATH={path}/lib`

You can do this globally or just in your editor.

this can globally be done on ubuntu by editing `~/.profile`
and adding the line in this file.
reloading for bash can be done by relogging or running `source ~/.profile` can be used to test the change for a terminal
