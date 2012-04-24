package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"runtime"
	"math/rand"
)

const size = 1000
const max_iter = 1000
const procs = 4

func to_img(brot* [size][size]uint, max uint) {
	gray := image.NewGray16(image.Rect(0,0,size,size))
	for x:=0; x<size; x++ {
		for y:=0; y<size; y++ {
			norm := float64(brot[x][y]*4)/float64(max)*65534
			if norm > 65534{
				norm = 65534
			}
			gray.SetGray16(
				x, y,
				color.Gray16{uint16(norm)})
		}
	}
	w, _ := os.OpenFile("./brot.png", os.O_CREATE|os.O_WRONLY, 0666)
	png.Encode(w, gray)
}

func cartesian_2_complex(x, y int) complex128 {
	return complex(
		-2 + (float64(x)/size * 3),
		-1.5 + (float64(y)/size * 3))
}

func complex_2_cartesian(point complex128) (float64, float64) {
	x := (real(point)+2)/3 * size
	y := (imag(point)+1.5)/3 * size
	return x, y
}

func buddahbrot_renderer(brot* [size][size]uint, orbits chan []complex128,
	finished chan uint) {
	done := 0
	max := uint(0)
	next := func() ([] complex128){
		orbit := <- orbits
		for orbit == nil && done < 3  {
			fmt.Printf("%v", done)
			done += 1
			orbit = <- orbits
		}
		return orbit
	}
	orbit := next()
	for orbit != nil {
		for _, point := range orbit {
			fx, fy := complex_2_cartesian(point)
			x, y := int(fx), int(fy)
			if (0<x && x<size && 0<y && y<size){
				brot[int(x)][int(y)] += 1
				if (brot[int(x)][int(y)] > max){
					max = brot[int(x)][int(y)]
				}
			}
		}
		orbit = next()
	}
	finished <- max
}

func mandelbrot(c complex128) (uint, []complex128) {
	orbit := make([]complex128, max_iter)
	z := c
	var iterations = uint(0)
	for (iterations<max_iter && real(z)+imag(z)<4) {
		z = z*z + c
		orbit[iterations] = z
		iterations++
	}
	return iterations, orbit
}

func brot_routine(brot *[size][size]uint, points chan complex128,
	orbits chan []complex128) {
	point, open := <- points
	for open {
		iterations, orbit := mandelbrot(point)
		if iterations == max_iter {
			orbits <- orbit
		}
		point, open = <- points
	}
	orbits <- nil
}


func main() {
	runtime.GOMAXPROCS(procs)
	rand.Seed(92388)
	var brot [size][size]uint
	points := make(chan complex128, 12)
	orbits := make(chan []complex128, 12)
	finished := make(chan uint)
	for i:=0; i<procs; i++ {
		go brot_routine(&brot, points, orbits)
	}
	go buddahbrot_renderer(&brot, orbits, finished)
	// for x:=0; x<size; x++ {
	// 	for y:=0; y<size; y++ {
	// 		points <- cartesian_2_complex(x, y)
	// 	}
	// }
	// fmt.Print(".\n")
	for i:=uint64(0); i<40000000; i++ {
		points <- complex(-2+rand.Float64()*3, -1.5+rand.Float64()*3)
	}
	close(points)
	fmt.Print("done\n")
	max := <- finished
	to_img(&brot, max)
}