package structsort_test

import (
	"log"
	"math"
	"sort"
	"testing"
	"time"

	"github.com/zafnz/structsort"
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
	Time       time.Time
}

func TestBasicSort(t *testing.T) {
	sillyStrings := []string{"banana", "zebra", "Mike", "john", "apples", "zebra", "zebra"}
	originalLen := len(sillyStrings)
	list := make([]testStruct, len(sillyStrings))
	for i := 0; i < len(sillyStrings); i++ {
		list[i].String = sillyStrings[i]
		list[i].Int = len(sillyStrings) - i
	}
	structsort.Sort(list, "String")
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
	structsort.Sort(list, "Int")
	for i := 0; i < len(sillyStrings); i++ {
		if list[i].Int != i+1 {
			t.Errorf("list did not sort by ints")
			break
		}
	}
}

func TestEmptyLists(t *testing.T) {
	err := structsort.Sort(make([]testStruct, 0), "StrPtr")
	if err != nil {
		t.Errorf("Sort returned error when sorting empty list: %s", err.Error())
	}
	err = structsort.Sort(make([]testStruct, 1), "StrPtr")
	if err != nil {
		t.Errorf("Sort returned error when sorting single item: %s", err.Error())
	}

	err = structsort.Sort(nil, "blah")
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
	structsort.Sort(list, "String")
	sort.Strings(sillyStrings)
	for i := 0; i < len(sillyStrings); i++ {
		if list[i].String != sillyStrings[i] {
			t.Errorf("list sorted case insensitive")
			break
		}
	}
}

func TestJsonFields(t *testing.T) {
	sillyStrings := []string{"banana", "zebra", "Mike", "john", "apples", "zebra", "zebra"}
	list := make([]testStruct, len(sillyStrings))
	for i := 0; i < len(sillyStrings); i++ {
		list[i].String = sillyStrings[i]
	}
	err := structsort.SortByTag(list, "json", "field_str")
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
}

func TestTimeSort(t *testing.T) {
	list := make([]testStruct, 10)
	list[0].Time = time.Unix(3000, 0)
	list[1].Time = time.Unix(9000, 0)
	list[2].Time = time.Unix(9000, 0)
	list[3].Time = time.Unix(1000, 0)
	list[4].Time = time.Unix(2511, 1)

	err := structsort.Sort(list, "Time")
	if err != nil {
		t.Fatalf("got error: %s", err)
	}
	for i := 1; i < len(list); i++ {
		if list[i-1].Time.After(list[i].Time) {
			t.Errorf("Time sort has %s after %s", list[i-1].Time.String(), list[i].Time.String())
		}
		//fmt.Println("times: ", list[i].Time.String())
	}
}

func TestValPtrs(t *testing.T) {
	sillyStrings := []string{"banana", "zebra", "Mike", "john", "apples", "zebra", "zebra"}

	list := make([]testStruct, len(sillyStrings))
	for i := 0; i < len(sillyStrings); i++ {
		s := sillyStrings[i]
		list[i].StrPtr = &s
	}
	err := structsort.Sort(list, "StrPtr")
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

	err := structsort.Sort(list, "StrPtr")
	if err != nil {
		t.Errorf("got error: %s", err)
		return
	}
	if list[len(sillyStrings)-1].StrPtr != nil {
		t.Errorf("Nil entry not sorted to end of list")
	}
}

func TestStructPtrs(t *testing.T) {
	sillyStrings := []string{"banana", "zebra", "Mike", "john", "apples", "zebra", "zebra"}

	list := make([]*testStruct, len(sillyStrings))
	for i := 0; i < len(sillyStrings); i++ {
		list[i] = new(testStruct)
		list[i].String = sillyStrings[i]
	}
	err := structsort.Sort(list, "String")
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
}
func TestInvalidField(t *testing.T) {
	list := make([]testStruct, 5)
	err := structsort.Sort(list, "NotValidField")
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
	err := structsort.Sort(list, "StrSlice")
	if err == nil {
		t.Errorf("did not report errors when trying to sort by a type it can't handle")
	}

	err = structsort.Sort(list, "Struct")
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

type MyType struct {
	x int
}

func (t MyType) Compare(t2 MyType) bool {
	return t.x <= t2.x
}
func TestCustomComparison(t *testing.T) {
	list := []struct {
		Thing MyType
	}{
		{Thing: MyType{1}},
		{Thing: MyType{5}},
		{Thing: MyType{2}},
		{Thing: MyType{4}},
		{Thing: MyType{3}},
	}
	err := structsort.Sort(list, "Thing")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(list)

	if list[0].Thing.x != 1 || list[4].Thing.x != 5 {
		t.Error("Out of order")
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
	list := []testStruct{
		{Float64: 3.1, Uint16: 5},
		{Float64: -0.9, Uint16: 2},
		{Float64: math.Log(-1.0), Uint16: 200},
	}
	err := structsort.Sort(list, "Float64")
	if err != nil {
		t.Fatal(err)
	}
	log.Print(list)
	if list[0].Float64 != -0.9 {
		t.Error("Output in wrong order")
	}
	if !math.IsNaN(list[2].Float64) {
		t.Error("NaN isn't last")
	}
	err = structsort.Sort(list, "Uint16")
	if err != nil {
		t.Fatal(err)
	}
	if list[0].Uint16 != 2 || list[2].Uint16 != 200 {
		t.Error("Uint16 sort order incorrect")
	}
}
