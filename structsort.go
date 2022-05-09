// structsort sorts a supplied list of structs by an arbitrary supplied field. You do not need
// to know the name of the field ahead of time, it's type, or even if it exists in the struct.
//  list := []struct {val: int, str:string} {
//    {2,"x"},
//    {1,"y"},
//  }
// structsort.Sort(list, "var") // sorts by int
// structsort.Sort(list, "str") // sorts by string
//
// It can sort by anything that is either a native type, or implements a String() function.
package structsort

import (
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
	s.rawType = s.slice.Index(0).Type()
	s.fieldIdx = fieldIndex(ref(s.slice.Index(0)).Type(), tagName, field)
	if s.fieldIdx == -1 {
		return fmt.Errorf("no such field %s", field)
	}

	sort.Sort(s)
	return s.err.err
}

type errHolder struct {
	err error
}

type genericSort struct {
	slice    reflect.Value
	rawType  reflect.Type
	fieldIdx int
	err      *errHolder
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

func (s genericSort) compare(a, b reflect.Value) (bool, error) {
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
		err := fmt.Errorf("unknown field type '%s'", a.Type())
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
// Structsort has built in sorting for strings, ints floats and if the underlying
// type supports it, by .String()
func Sort(list interface{}, field string) error {
	return sortInternal(list, field, "")
}

// SortByTag works just like Sort, except it will check for `json:"field"` for
// the field to sort by. If no tag exists with that field name, it will try to
// use the struct field's name just like Sort.
func SortByTag(list interface{}, tag string, field string) error {
	return sortInternal(list, field, tag)
}
