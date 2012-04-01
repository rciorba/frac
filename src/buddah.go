package main

import (
	"fmt"
)

const size = 128
const max_iter = 1

func square(x int, y int) (int, int) {
	return (x*x - y*y), (2*x*y)
}

func mandelbrot(x int, y int) int {
	// c := complex(float64(x), float64(y))
	o_x := x
	o_y := y
	i:=0
	for (i<max_iter && (x*x)+(y*y)<4) {
		i++
		x, y = square(x, y)
		x += o_x
		y += o_y
	}
	return i
}

func main() {
	var brot [size][size]bool
	for x:=0; x<size; x++ {
		for y:=0; y<size; y++ {
			p := mandelbrot(x, y)
			if p>0 {
				brot[x][y] = true
				fmt.Print(".")
			} else {
				fmt.Print(" ")
			}
		}
		fmt.Print("\n")
	}
}