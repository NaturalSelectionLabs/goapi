package ff

import (
	"reflect"
)

type FlattenedStruct struct {
	Type   reflect.Type
	Fields []*FlattenedField
}

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

type FlattenedField struct {
	Path  []int
	Field reflect.StructField
}

// read is to read the embedding struct field.
func (f *FlattenedField) Get(target reflect.Value) reflect.Value {
	path := f.Path

	var value reflect.Value = target

	for _, index := range path {
		if value.Kind() == reflect.Ptr && value.Elem().Kind() == reflect.Struct {
			value = value.Elem()
		}

		value = value.Field(index)
	}

	return value
}

// set is used to set the value of the embedding struct field.
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

		if i == len(path)-1 {
			break
		} else {
			fieldType = fieldType.Field(path[i]).Type
		}
	}

	return fieldType.Field(path[i])
}
