// structsort sorts a supplied list of structs by an arbitrary supplied field. You do not need
// to know the name of the field ahead of time, it's type, or even if it exists in the struct.
//  list := []struct {val: int, str:string} {
//    {2,"x"},
//    {1,"y"},
//  }
//  structsort.Sort(list, "var") // sorts by int
//  structsort.Sort(list, "str") // sorts by string
//
// It can sort by anything that is either a native type, implements a Compare() method (see
// customCompare example, or implements a String() function.
package structsort

import (
	"errors"
	"fmt"
	"reflect"
	"sort"
)

type stringer interface {
	String() string
}

func sortInternal(list interface{}, field string, tagName string) error {
	s := genericSort{}
	s.slice = reflect.Indirect(reflect.ValueOf(list))
	if s.slice.Kind() != reflect.Slice {
		return fmt.Errorf("sort expects list parameter to be a slice")
	}

	if s.slice.Len() <= 1 {
		// Easy, it's sorted.
		return nil
	}
	s.err = new(errHolder)
	first := s.slice.Index(0)
	s.rawType = first.Type()
	s.fieldIdx = fieldIndex(ref(first).Type(), tagName, field)
	if s.fieldIdx == -1 {
		return fmt.Errorf("no such field %s", field)
	}
	s.fieldType = ref(ref(first).Field(s.fieldIdx)).Type()
	var err error
	s.hasCompare, err = s.checkHasCompare(s.fieldType)
	if err != nil {
		fmt.Printf("struct's sort field has a Compare method with issues: %s\n", err.Error())
	}
	sort.Sort(s)
	return s.err.err
}

type errHolder struct {
	err error
}

type genericSort struct {
	slice      reflect.Value
	rawType    reflect.Type // The type of the struct itself
	fieldIdx   int
	err        *errHolder
	hasCompare bool         // Whether the field we are sorting by has a Compare method
	fieldType  reflect.Type // The type of the field we are sorting by
}

func (s genericSort) Len() int {
	return s.slice.Len()
}
func (s genericSort) Swap(i, j int) {
	a := s.slice.Index(i)
	b := s.slice.Index(j)
	tmp := reflect.New(s.rawType).Elem()
	tmp.Set(a)
	a.Set(b)
	b.Set(tmp)
}
func (s genericSort) Less(i, j int) bool {
	aStruct := ref(s.slice.Index(i))
	bStruct := ref(s.slice.Index(j))
	if aStruct == nil || bStruct == nil {
		return false
	}
	aVal := ref(aStruct.Field(s.fieldIdx))
	bVal := ref(bStruct.Field(s.fieldIdx))
	if aVal == nil && bVal == nil {
		return false
	} else if aVal == nil {
		return false
	} else if bVal == nil {
		return true
	}
	result, err := s.compare(*aVal, *bVal)
	if err != nil {
		s.err.err = err
	}
	return result
}
func (s genericSort) checkHasCompare(t reflect.Type) (bool, error) {
	// t is the type we are sorting on, that may have a method called "Compare"
	m, found := t.MethodByName("Compare")
	if !found {
		return false, nil
	}
	if m.Type.NumIn() != 2 { // Method's signature includes it's own type as first parameter
		return false, errors.New("compare method for type does not have one parameter")
	}
	if m.Type.NumOut() != 1 {
		return false, errors.New("compare method for type does not return only one value")
	}
	if m.Type.In(0) != t || m.Type.In(1) != t {
		return false, errors.New("compare method's input is not of it's own type")
	}
	if m.Type.Out(0).Kind() != reflect.Bool {
		return false, errors.New("compare method's output is not of type bool")

	}
	return true, nil
}
func (s genericSort) useCompare(a, b reflect.Value) (bool, error) {
	compareFunc := a.MethodByName("Compare")
	result := compareFunc.Call([]reflect.Value{b})
	if len(result) != 1 {
		return false, fmt.Errorf("compare method did not return only one argument")
	}
	if result[0].Type().Kind() != reflect.Bool {
		return false, errors.New("compare method's result was not bool")
	}
	return result[0].Bool(), nil
}

func (s genericSort) compare(a, b reflect.Value) (bool, error) {
	if s.hasCompare {
		return s.useCompare(a, b)
	}
	switch a.Type().Kind() {
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return a.Uint() < b.Uint(), nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return a.Int() < b.Int(), nil
	case reflect.Float32, reflect.Float64:
		return a.Float() < b.Float(), nil
	case reflect.String:
		return a.String() < b.String(), nil
	case reflect.Bool:
		return !a.Bool(), nil
	default:
		stringerType := reflect.TypeOf((*stringer)(nil)).Elem()

		if a.Type().Implements(stringerType) {
			aStr := a.Interface().(stringer).String()
			bStr := b.Interface().(stringer).String()
			return aStr < bStr, nil
		}
		err := fmt.Errorf("unsupported field type '%s'. Add a Compare(%s) bool method", a.Type(), a.Type())
		s.err.err = err
		return false, err
	}
}

func fieldIndex(structType reflect.Type, tag string, name string) int {
	numFields := structType.NumField()
	namedIdx := -1
	tagIdx := -1
	for i := 0; i < numFields; i++ {
		f := structType.Field(i)
		if f.Name == name {
			namedIdx = i
		}
		if tag != "" && f.Tag.Get(tag) == name {
			tagIdx = i
		}
	}
	if tagIdx == -1 {
		return namedIdx
	}
	return tagIdx
}

// Takes a reflect.Value, if it's a Ptr and IsNil, return Nil, otherwise return a pointer to that reflect.Value.
func ref(val reflect.Value) *reflect.Value {
	if val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return nil
		}
		val = val.Elem()
		return &val
	}
	return &val
}

// Sort sorts a slice of structs, by the named field.
// Pointers to structs and pointers to fields are dereferenced when appropriate.
// Nil pointers are sorted last.
//
// Structsort has built in sorting for strings, ints floats.
// If the sort field has a method called `Compare`` and the method has a signature
// that it takes it's own type, and returns a bool, then this method will be called
// when doing comparisons.
//
// If the sort field does not have a built in type, but does have a `String()` method,
// then it will be used in sorting
//
// NOTE: See examples, they demonstrate this quite well.
func Sort(list interface{}, field string) error {
	return sortInternal(list, field, "")
}

// SortByTag works just like Sort, except it will check for `json:"field"` for
// the field to sort by. If no tag exists with that field name, it will try to
// use the struct field's name just like Sort.
func SortByTag(list interface{}, tag string, field string) error {
	return sortInternal(list, field, tag)
}
