# Simple Lisp

A simple `Lisp` interpreter written in Go with minimal dependencies, even to 
the standard library.

It requires the following types to be present:

* `os.File`
* `unsafe.Pointer`

It requires the following constants / variables to be present:

* `os.Args`
* `os.Stdout`

It requires the following functions / methods to be present:

* `fmt.Sprintf`
* `math.Hypot`
* `math.RoundToEven`
* `os.(*File).Close`
* `os.(*File).Read`
* `os.(*File).Write`
* `os.OpenFile`
* `strconv.FormatFloat`
* `strconv.Itoa`
* `strconv.ParseComplex`
* `strconv.ParseFloat`
* `strconv.ParseInt`
* `strconv.Quote`
* `strconv.Unquote`
* `strings.ContainsRune`
* `strings.HasPrefix`
* `strings.Join`
* `strings.Split`

Command to run the Mandelbrot Set example program:

```bash
$ go run . mandelbrot.scm
```

It should give you this image as output:

![Mandelbrot Set](mandelbrot.png)
