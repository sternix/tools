/*
small utility for generating raw image for wayland
you can use this file as:
mmap the out file
create pool for the file
create buffer with given width and height
attach the buffer to surface
and commit
you dont need to paint the surface
*/
package main

import (
	"bufio"
	"flag"
	"fmt"
	"image"
	"image/color"
	"log"
	"os"
)

import (
	_ "golang.org/x/image/bmp"
	_ "golang.org/x/image/tiff"
	_ "golang.org/x/image/webp"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
)

var (
	inFile  = flag.String("in", "", "Input image file")
	outFile = flag.String("out", "", "Raw output file")
)

func init() {
	flag.Parse()
	log.SetFlags(0)
}

func ImageFromFile(fileName string) (image.Image, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	br := bufio.NewReader(f)
	img, _, err := image.Decode(br)
	if err != nil {
		return nil, err
	}
	return img, nil
}

func main() {
	if *inFile == "" {
		flag.Usage()
		os.Exit(1)
	}

	if *outFile == "" {
		flag.Usage()
		os.Exit(1)
	}

	img, err := ImageFromFile(*inFile)
	if err != nil {
		log.Fatal(err)
	}

	out, err := os.Create(*outFile)
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	bwOut := bufio.NewWriter(out)
	var buf [4]byte
	// https://blog.golang.org/go-image-package
	b := img.Bounds()
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			r = r >> 8
			g = g >> 8
			b = b >> 8
			a = a >> 8
			// in order b , g , r , a
			buf[0] = byte(b * a / 255)
			buf[1] = byte(g * a / 255)
			buf[2] = byte(r * a / 255)
			buf[3] = byte(a)
			bwOut.Write(buf[:])
		}
	}

	bwOut.Flush()

	fmt.Printf("Width:%d , Height:%d\n", b.Dx(), b.Dy())
}
