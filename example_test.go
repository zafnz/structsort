package structsort_test

import (
	"fmt"

	"github.com/zafnz/structsort"
)

func Example() {
	type Person struct {
		Name string
		Age  int
	}
	// Create some basic records
	var records = []Person{
		{
			Name: "Charlie",
			Age:  22,
		},
		{
			Name: "Bob",
			Age:  33,
		},
		{
			Name: "Alice",
			Age:  44,
		},
	}

	// Now for the magic
	fmt.Println("Sort by Name")
	structsort.Sort(records, "Name")
	for _, r := range records {
		fmt.Printf("%s: %d\n", r.Name, r.Age)
	}
	fmt.Println("Now sort by Age")
	structsort.Sort(records, "Age")
	for _, r := range records {
		fmt.Printf("%s: %d\n", r.Name, r.Age)
	}
	// Output:
	// Sort by Name
	// Alice: 44
	// Bob: 33
	// Charlie: 22
	// Now sort by Age
	// Charlie: 22
	// Bob: 33
	// Alice: 44
}
