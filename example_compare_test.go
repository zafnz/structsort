package structsort_test

import (
	"fmt"

	"github.com/zafnz/structsort"
)

type Point struct {
	x int
	y int
}

func (d Point) Compare(other Point) bool { // Compare how far from origin
	return (d.x + d.y) < (other.x + other.y)
}

type Thing struct {
	Name     string
	Location Point
}

func Example_customCompare() {
	// This example is for the specific custom Compare function for a type
	// See the more general example for the typical usecase
	list := []Thing{
		{"Orange", Point{2, 3}},
		{"Lemon", Point{1, 2}},
		{"Apple", Point{1, 1}},
	}
	structsort.Sort(list, "Location")
	fmt.Printf("%v\n", list)
	// Output:
	// [{Apple {1 1}} {Lemon {1 2}} {Orange {2 3}}]
}
