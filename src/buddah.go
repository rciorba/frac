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
const max_iter = 2000
const procs = 4

func to_img(brot* [size][size]uint, max uint) {
	gray := image.NewGray16(image.Rect(0,0,size,size))
	norm_gray := image.NewGray16(image.Rect(0,0,size,size))
	for x:=0; x<size; x++ {
		for y:=0; y<size; y++ {
			pix := brot[x][y]
			norm := float64(brot[x][y]*4)/float64(max)*65534
			if norm > 65534{
				norm = 65534
			}
			if pix > 65534{
				pix = 65534
			}
			gray.SetGray16(
				x, y,
				color.Gray16{uint16(pix)})
			norm_gray.SetGray16(
				x, y,
				color.Gray16{uint16(norm)})
		}
	}
	w, _ := os.OpenFile("./brot.png", os.O_CREATE|os.O_WRONLY, 0666)
	png.Encode(w, gray)
	n, _ := os.OpenFile("./brot-norm.png", os.O_CREATE|os.O_WRONLY, 0666)
	png.Encode(n, norm_gray)
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
		for orbit == nil && done < procs-1  {
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
				brot[x][y] += 1
				if (brot[x][y] > max){
					max = brot[x][y]
				}
			}
		}
		orbit = next()
	}
	fmt.Printf("max: %v\n", max)
	finished <- max
}

func mandelbrot(c complex128) (uint, []complex128) {
	orbit := make([]complex128, max_iter)
	z := c
	var iterations = uint(0)
	x, y := real(z), imag(z)
	for (iterations<max_iter && x*x+y*y<4) {
		z = z*z + c
		orbit[iterations] = z
		iterations++
		x, y = real(z), imag(z)
	}
	return iterations, orbit[0:iterations]
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
	points := make(chan complex128, 12*procs)
	orbits := make(chan []complex128, 12*procs)
	finished := make(chan uint)
	for i:=0; i<procs; i++ {
		go brot_routine(&brot, points, orbits)
	}
	go buddahbrot_renderer(&brot, orbits, finished)
	exposure := uint64(40000000)
	fe := float64(exposure)
	for i:=uint64(0); i<exposure; i++ {
		points <- complex(-2+rand.Float64()*3, -1.5+rand.Float64()*3)
		if (i%100000==0){
			fmt.Printf("%v %%\n", float64(i)/fe*100)
		}
	}
	close(points)
	fmt.Print("done\n")
	max := <- finished
	to_img(&brot, max)
}