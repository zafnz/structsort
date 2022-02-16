package structsort

import (
	"fmt"
	"sort"
	"testing"
)

type embeddedStruct struct {
	val string
}

func (e embeddedStruct) String() string {
	return e.val
}

type testStruct struct {
	String     string  `json:"field_str"`
	Int        int     `json:"field_int"`
	StrPtr     *string `json:"field_strp"`
	IntPtr     *int
	StrSlice   []string
	Float      float32
	Bool       bool
	Int8       int8
	Int16      int16
	Int32      int32
	Int64      int64
	Uint       uint
	Uint8      uint8
	Uint16     uint16
	Uint32     uint32
	Uint64     uint64
	Uintptr    *uint
	Float32    float32
	Float64    float64
	Complex64  complex64
	Complex128 complex128
	Struct     embeddedStruct
}

func TestBasicSort(t *testing.T) {
	sillyStrings := []string{"banana", "zebra", "Mike", "john", "apples", "zebra", "zebra"}
	originalLen := len(sillyStrings)
	list := make([]testStruct, len(sillyStrings))
	for i := 0; i < len(sillyStrings); i++ {
		list[i].String = sillyStrings[i]
		list[i].Int = len(sillyStrings) - i
	}
	Sort(list, "String")
	if originalLen != len(list) {
		t.Fatalf("list has changed length, was %d now %d", originalLen, len(list))
	}
	sort.Strings(sillyStrings)
	for i := 0; i < len(sillyStrings); i++ {
		if list[i].String != sillyStrings[i] {
			t.Errorf("list did not sort by strings")
			break
		}
	}
	Sort(list, "Int")
	for i := 0; i < len(sillyStrings); i++ {
		if list[i].Int != i+1 {
			t.Errorf("list did not sort by ints")
			break
		}
	}
}

func TestEmptyLists(t *testing.T) {
	err := Sort(make([]testStruct, 0), "StrPtr")
	if err != nil {
		t.Errorf("Sort returned error when sorting empty list: %s", err.Error())
	}
	err = Sort(make([]testStruct, 1), "StrPtr")
	if err != nil {
		t.Errorf("Sort returned error when sorting single item: %s", err.Error())
	}

	err = Sort(nil, "blah")
	if err == nil {
		t.Errorf("Sort did not report error when trying to sort nil")
	}
}

func TestCaseSort(t *testing.T) {
	sillyStrings := []string{"banana", "Zebra", "zebra", "john", "John"}
	list := make([]testStruct, len(sillyStrings))
	for i := 0; i < len(sillyStrings); i++ {
		list[i].String = sillyStrings[i]
	}
	Sort(list, "String")
	sort.Strings(sillyStrings)
	for i := 0; i < len(sillyStrings); i++ {
		if list[i].String != sillyStrings[i] {
			t.Errorf("list sorted case insensitive")
			break
		}
	}
	return
}

func TestJsonFields(t *testing.T) {
	sillyStrings := []string{"banana", "zebra", "Mike", "john", "apples", "zebra", "zebra"}
	list := make([]testStruct, len(sillyStrings))
	for i := 0; i < len(sillyStrings); i++ {
		list[i].String = sillyStrings[i]
	}
	err := SortByTag(list, "json", "field_str")
	if err != nil {
		t.Fatalf("got error: %s", err)
	}
	sort.Strings(sillyStrings)
	for i := 0; i < len(sillyStrings); i++ {
		if list[i].String != sillyStrings[i] {
			t.Errorf("list not sorted by json tag")
			break
		}
	}
	return
}

func TestValPtrs(t *testing.T) {
	sillyStrings := []string{"banana", "zebra", "Mike", "john", "apples", "zebra", "zebra"}

	list := make([]testStruct, len(sillyStrings))
	for i := 0; i < len(sillyStrings); i++ {
		s := sillyStrings[i]
		list[i].StrPtr = &s
	}
	err := Sort(list, "StrPtr")
	if err != nil {
		t.Fatalf("got error: %s", err)
	}
	sort.Strings(sillyStrings)
	for i := 0; i < len(sillyStrings); i++ {
		if *list[i].StrPtr != sillyStrings[i] {
			t.Errorf("list not sorted by json tag")
			break
		}
	}
}
func TestNilPtrSort(t *testing.T) {
	sillyStrings := []string{"banana", "zebra", "mike", "john", "apples", "zebra", "zebra", "cowboy"}

	list := make([]testStruct, len(sillyStrings))

	for i := 0; i < len(sillyStrings); i++ {
		s := sillyStrings[i]
		list[i].StrPtr = &s
	}
	list[2].StrPtr = nil

	err := Sort(list, "StrPtr")
	if err != nil {
		t.Errorf("got error: %s", err)
		return
	}
	if list[len(sillyStrings)-1].StrPtr != nil {
		t.Errorf("Nil entry not sorted to end of list")
	}
	return
}

func TestStructPtrs(t *testing.T) {
	sillyStrings := []string{"banana", "zebra", "Mike", "john", "apples", "zebra", "zebra"}

	list := make([]*testStruct, len(sillyStrings))
	for i := 0; i < len(sillyStrings); i++ {
		list[i] = new(testStruct)
		list[i].String = sillyStrings[i]
	}
	err := Sort(list, "String")
	if err != nil {
		t.Fatalf("got error: %s", err)
	}
	sort.Strings(sillyStrings)
	for i := 0; i < len(sillyStrings); i++ {
		if list[i].String != sillyStrings[i] {
			t.Errorf("list not sorted by string")
			break
		}
	}
	return
}
func TestInvalidField(t *testing.T) {
	list := make([]testStruct, 5)
	err := Sort(list, "NotValidField")
	if err == nil {
		t.Errorf("did not report errors when trying to sort by a type it can't handle")
	}
}

func TestUnknownType(t *testing.T) {
	sillyStrings := []string{"banana", "Zebra", "zebra", "john", "John", "d", "c", "b", "a"}
	list := make([]testStruct, len(sillyStrings))
	for i := 0; i < len(sillyStrings); i++ {
		list[i].Struct.val = sillyStrings[i]
	}
	err := Sort(list, "StrSlice")
	if err == nil {
		t.Errorf("did not report errors when trying to sort by a type it can't handle")
	}

	err = Sort(list, "Struct")
	if err != nil {
		t.Errorf("reported error sorting Struct that implements String(): %s", err.Error())
	}
	sort.Strings(sillyStrings)
	for i := 0; i < len(sillyStrings); i++ {
		if list[i].Struct.val != sillyStrings[i] {
			t.Errorf("list not sorted by String()")
			break
		}
	}

}

/*
	Float      float32
	Bool       bool
	Int8       int8
	Int16      int16
	Int32      int32
	Int64      int64
	Uint       uint
	Uint8      uint8
	Uint16     uint16
	Uint32     uint32
	Uint64     uint64
	Uintptr    *uint
	Float32    float32
	Float64    float64
	Complex64  complex64
	Complex128 complex128
	Struct     embeddedStruct
*/
func TestAllTheTypes(t *testing.T) {
	//TODO No idea how to do this without a lot of typing

}

func TestExample(t *testing.T) {
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
	fmt.Printf("Sort by Name\n")

	// Now for the magic
	Sort(records, "Name")
	for _, r := range records {
		fmt.Printf("%s: %d\n", r.Name, r.Age)
	}
	// Output is Alice, Bob, and Charlie
	fmt.Printf("Now sort by Age\n")
	Sort(records, "Age")
	for _, r := range records {
		fmt.Printf("%s: %d\n", r.Name, r.Age)
	}
	// Now output is Charlie, Bob, Alice
}
