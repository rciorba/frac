package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"runtime"
)

const size = 128
const max_iter = 10000

func square(x int, y int) (int, int) {
	return (x*x - y*y), (2*x*y)
}


func to_img(brot [size][size]uint) {
	max := uint(0)
	gray := image.NewGray16(image.Rect(0,0,size,size))
	for x:=0; x<size; x++ {
		for y:=0; y<size; y++ {
			if brot[x][y] > max{
				max = brot[x][y]
			}
			if brot[x][y] >= max_iter{
				gray.SetGray16(x, y, color.Black)
			} else {
				gray.SetGray16(
					x, y,
					color.Gray16{uint16(brot[x][y]*200)})
			}
		}
	}
	fmt.Printf("%v", max)
	w, _ := os.OpenFile("./brot.png", os.O_CREATE|os.O_WRONLY, 0666)
	png.Encode(w, gray)
}



func mandelbrot(x int, y int) uint {
	c := complex(float64(x)/size, float64(y)/size)
	z := c
	var i = uint(0)
	for (i<max_iter && real(z)+imag(z)<4) {
		i++
		z = z*z + c
	}
	return i
}

// setup 4 goroutines all consuming from one queue
// all 4 goroutines will send one message when they finish on a second queue

func main() {
	runtime.GOMAXPROCS(4)
	var brot [size][size]uint
	// in := make(chan image.Point)
	// finished := make(chan bool)
	for x:=0; x<size; x++ {
		fmt.Print(".")
		for y:=0; y<size; y++ {
			brot[x][y] = mandelbrot(x, y)
		}
	}
	to_img(brot)
}