package structsort

import (
	"fmt"
	"reflect"
	"sort"
)

type Stringer interface {
	String() string
}

// Sort sorts a slice of structs, by the named field.
// If useJsonTags then look for json field names specified and if the field name is found there,
// then it will sort with that field.
// Pointers to structs and pointers to fields are dereferenced when appropriate.
// Nil pointers are sorted last.
func Sort(list interface{}, field string) error {
	return sortInternal(list, field, "")
}

// SortByTag works just like Sort, except it will check for `json:"field"` for
// the field to sort by. If no tag exists with that field name, it will try to
// use the struct field's name just like Sort.
func SortByTag(list interface{}, tag string, field string) error {
	return sortInternal(list, field, tag)
}

// sortInternal does all the heavy lifting, it supports all the options
func sortInternal(list interface{}, field string, tagName string) error {
	slice := reflect.Indirect(reflect.ValueOf(list))
	if slice.Kind() != reflect.Slice {
		return fmt.Errorf("Sort expects list parameter to be a slice")
	}
	var err error

	if slice.Len() <= 1 {
		// Easy, it's sorted.
		return nil
	}
	structType := slice.Index(0).Type()

	if tagName != "" {
		field = tagToField(structType, field)
	}

	sort.Slice(list, func(a, b int) bool {
		aStruct := deref(slice.Index(a))
		bStruct := deref(slice.Index(b))
		if aStruct == nil || bStruct == nil {
			return false
		}

		aValue := aStruct.FieldByName(field)
		if !aValue.IsValid() {
			err = fmt.Errorf("struct does not have field %s", field)
			return false
		}
		bValue := bStruct.FieldByName(field)
		if aValue.Type().Kind() != bValue.Type().Kind() {
			// It's moderately hard to get to this issue. You're really having to fuck around
			// with lists using reflection.
			panic(fmt.Sprintf("list items aren't of same type %s vs %s", aValue.Type().Name(), bValue.Type().Name()))
		}
		var aNil, bNil bool
		if aValue.Kind() == reflect.Ptr {
			if aValue.IsNil() {
				aNil = true
			} else {
				aValue = aValue.Elem()
			}
		}
		if bValue.Kind() == reflect.Ptr {
			if bValue.IsNil() {
				bNil = true
			} else {
				bValue = bValue.Elem()
			}
		}
		if aNil && bNil {
			return false
		} else if aNil && !bNil {
			return false
		} else if !aNil && bNil {
			return true
		}
		switch aValue.Type().Kind() {
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return aValue.Uint() < bValue.Uint()
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return aValue.Int() < bValue.Int()
		case reflect.Float32, reflect.Float64:
			return aValue.Float() < bValue.Float()
		case reflect.String:
			return aValue.String() < bValue.String()
		case reflect.Bool:
			return !aValue.Bool()
		default:
			stringerType := reflect.TypeOf((*Stringer)(nil)).Elem()

			if aValue.Type().Implements(stringerType) {
				a := aValue.Interface().(Stringer).String()
				b := bValue.Interface().(Stringer).String()
				return a < b
			}
			err = fmt.Errorf("unknown field type %s", aValue.Type().Name())
			return false
		}
	})
	return err
}

func deref(val reflect.Value) *reflect.Value {
	if val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return nil
		}
		val = val.Elem()
		return &val
	}
	return &val
}

// tagToField returns the name of a field with the provided json tag
func tagToField(structType reflect.Type, tag string) (field string) {
	numFields := structType.NumField()
	field = tag
	for i := 0; i < numFields; i++ {
		f := structType.Field(i)

		if f.Tag.Get("json") == field {
			field = f.Name
		}
	}
	return
}
