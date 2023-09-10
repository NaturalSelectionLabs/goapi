package ff_test

import (
	"reflect"
	"testing"

	ff "github.com/NaturalSelectionLabs/goapi/lib/flat-fields"
	"github.com/ysmood/got"
)

type Animal struct {
	Name string
}

type Dog struct {
	Animal
	Breed string
}

func TestBasic(t *testing.T) {
	g := got.T(t)

	hound := Dog{Animal: Animal{Name: "Spotty"}, Breed: "Dalmatian"}
	val := reflect.ValueOf(&hound)

	s := ff.Parse(reflect.TypeOf(hound))

	g.Eq(s.Fields[0].Get(val).Interface(), "Spotty")

	for _, f := range s.Fields {
		f.Set(val, reflect.ValueOf("test"))
	}

	g.Eq(val.Interface(), &Dog{Animal: Animal{Name: "test"}, Breed: "test"})
}

type Cat struct {
	*Animal
	Breed string
}

func TestPointer(t *testing.T) {
	g := got.T(t)

	cat := Cat{Animal: &Animal{Name: "Kitty"}, Breed: "Bobtail"}
	val := reflect.ValueOf(&cat)

	s := ff.Parse(reflect.TypeOf(cat))

	g.Eq(s.Fields[0].Get(val).Interface(), "Kitty")

	for _, f := range s.Fields {
		f.Set(val, reflect.ValueOf("test"))
	}

	g.Eq(val.Interface(), &Cat{Animal: &Animal{Name: "test"}, Breed: "test"})
}
