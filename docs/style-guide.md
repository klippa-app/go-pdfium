# Style Guide

In this project, we follow the style guide explained in 
https://github.com/golang-standards/project-layout. Below is a style guide 
for each directory in the project. In short, the project structure should 
look like this: 
```
/go-pdfium
    /docs
        {documentation}
    /internal
        {internal files/directories}
    /pkg
        {public files/directories}
    /vendor
        {external packages}
.gitignore
go.mod
README.md
```

## `/go-pdfium`
This is the main directory, this will hold all subdirectories (listed below) 
and other files that can not be placed in a subdirectory, like `.gitignore` 
and `go.mod`.

### `/docs`
The directory containing all documents belonging to the project. No code is 
allowed in this directory.

### `/internal`
Code used by the project, but not suitable for use outside the project. When 
the project is imported, this code can not be accessed.

### `/pkg`
This directory should contain the main code. Everything in here is 
accessible and importable by outside users. 

If necessary, this package can also have an internal directory 
(`/pkg/internal`).

### `/vendor`
This directory contains all external packages imported by the project and is 
generated automatically by the command `go mod vendor`.