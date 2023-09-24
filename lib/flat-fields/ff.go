// Package ff contains a struct flattening utility.
package ff

import (
	"reflect"
)

// FlattenedStruct is a struct that contains the flattened fields.
type FlattenedStruct struct {
	Type   reflect.Type
	Fields []*FlattenedField
}

// Parse is to parse the embedding struct and return the flattened struct.
func Parse(t reflect.Type) *FlattenedStruct {
	paths := parse(t)

	fields := make([]*FlattenedField, len(paths))
	for i, path := range paths {
		fields[i] = &FlattenedField{
			Path:  path,
			Field: get(t, path),
		}
	}

	return &FlattenedStruct{
		Type:   t,
		Fields: fields,
	}
}

// FlattenedField is a struct that contains the path indices and the field.
type FlattenedField struct {
	Path  []int
	Field reflect.StructField
}

// Get the value of embedding struct field.
func (f *FlattenedField) Get(target reflect.Value) reflect.Value {
	path := f.Path

	value := target

	for _, index := range path {
		if value.Kind() == reflect.Ptr && value.Elem().Kind() == reflect.Struct {
			value = value.Elem()
		}

		value = value.Field(index)
	}

	return value
}

// Set the value of embedding struct field.
func (f *FlattenedField) Set(target reflect.Value, val reflect.Value) {
	path := f.Path

	for i := 0; i < len(path); i++ {
		if target.Kind() == reflect.Ptr && target.Elem().Kind() == reflect.Struct {
			target = target.Elem()
		}

		if i == len(path)-1 {
			target.Field(path[i]).Set(val)
		} else {
			target = target.Field(path[i])
		}
	}
}

// parse is to parse embedding struct and return the path indices.
func parse(target reflect.Type) [][]int {
	var paths [][]int

	if target == nil || target.Kind() != reflect.Struct {
		return paths
	}

	for i := 0; i < target.NumField(); i++ {
		field := target.Field(i)
		if field.Anonymous {
			embeddedType := field.Type
			if embeddedType.Kind() == reflect.Ptr {
				embeddedType = embeddedType.Elem()
			}

			embeddedPaths := parse(embeddedType)

			for _, embeddedPath := range embeddedPaths {
				path := append([]int{i}, embeddedPath...)
				paths = append(paths, path)
			}
		} else {
			paths = append(paths, []int{i})
		}
	}

	return paths
}

func get(target reflect.Type, path []int) reflect.StructField {
	fieldType := target
	i := 0

	for ; i < len(path); i++ {
		if fieldType.Kind() == reflect.Ptr {
			fieldType = fieldType.Elem()
		}

		if i != len(path)-1 {
			fieldType = fieldType.Field(path[i]).Type
		} else {
			break
		}
	}

	return fieldType.Field(path[i])
}
