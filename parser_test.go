package main

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParser(t *testing.T) {
	src, err := ioutil.ReadFile("mandelbrot.scm")
	require.NoError(t, err)
	ret := CreateParser(string(src)).Parse()
	println(ret.String())
}
