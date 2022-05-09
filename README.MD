# structsort #
[![GoDoc](https://godoc.org/github.com/zafnz/structsort?status.svg)](https://godoc.org/github.com/zafnz/structsort)

structsort is a Go library for sorting structs by atribitrary field name in a slice. (eg you don't know what the
field name could be before compile time!)

It exists because I could not find code to do so, and when I tried to write it, it turned out to be 
non-trivial. 

## Installation

Just like every other Go module.
```bash
go get github.com/zafnz/structsort
```

or using import
```go
import "github.com/zafnz/structsort"
```
and then run `go get` without any parameters

## Usage
### Sort
The `Sort(list, field)` function sorts a slice(list) of structs much like the `sort.Slice()` function from the `sort` library. However here you only need to provide a field name and it will do the rest, including doing sensible sorting on strings, ints, and floats.

```go
import "github.com/zafnz/structsort"

type Person struct {
    Name string
    Age int
}
// Create some basic records
var records = []Person {
    Person {
        Name: "Charlie",
        Age: 22,
    },
    Person {
        Name: "Bob",
        Age: 33,
    },
    Person {
        Name: "Alice",
        Age: 44,
    },
}
// Now for the magic
err := structsort.Sort(records, "Name")
for _, r := range records {
    fmt.Printf("%s: %d", r.Name, r.Age)
}
// Output is Alice, Bob, and Charlie
err = structsort.Sort(records, "Age")
for _, r := range records {
    fmt.Printf("%s: %d", r.Name, r.Age)
}
// Now output is Charlie, Bob, Alice
```

### SortByTag
```go
structsort.SortByTag(list, tagName, fieldName)
```

Sometimes your structs have tags, such as 
```go
type MyRecord {
    GivenName string `json:given_name"`
}
```
Using `SortByTag(records, "json", "given_name")` you can now sort by the GivenName 
field by refering to it's json name. 

If it's not clear why this is useful, imagine a web API where the user
can supply the field name to order the results. 
```go
type DatabaseRecord {
    Id int `json:id`
    GivenName string `json:given_name"`
    Surname string `json:surname"`
}

func MyWebApi(w http.ResponseWriter, r *http.Request) {
    query := r.URL.Query()
    sortField := query["sort"]
    records = someLargeApiCall()
    structsort.SortByTag(records, sortField)
    ... Insert rest of code ...
}
```
### Sort fields Compare method
When sorting by a field, if the field itself implements a method `T.Compare(T) bool`
then this method will be called. Note: This is not the struct itself, if you have
a `Compare` method for the struct that you're sorting, you should be using the `sort`
library instead!

Example:
```go
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
```

## Functions
### Sort
```go
func Sort(list interface{}, field string) error 
```
Sorts a slice/list of structs based on the values in the specified field. Fields that are 
pointers will be deferenced correctly. Nil values are sorted to the end of the list.

If a field implements a String() method, then it will be sorted as a simple string.

Unknown fields, unknown types, invalid pointers, etc all return err. This function
should not panic unless you try really hard.

### SortByTag
```go
func SortByTag(list interface{}, tag string, field string) error
```
Sorts a slice/list of structs, similar to the `Sort` function, however it will first try to find
the correct field using the struct tags. (Eg ``json:"field"``)

# Contributing

All contributions are welcome via the usual way.

# License
