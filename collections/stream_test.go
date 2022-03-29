package collections

import (
	"reflect"
	"strings"
	"testing"
)

func buildStream(data []interface{}) *Stream {
	return NewStream(NewSliceIterator(data))
}

func TestMap(t *testing.T) {
	data := []interface{}{"A", "B", "C", "D", "E"}
	stream := buildStream(data)

	mapFn := func(v interface{}) interface{} {
		return strings.ToLower(v.(string))
	}

	expected := []interface{}{"a", "b", "c", "d", "e"}
	actual := stream.Map(mapFn).ToSlice()

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("expected %v got %v", expected, actual)
	}
}

func TestFilter(t *testing.T) {
	data := []interface{}{5, 2, 3, 4, 4}
	stream := buildStream(data)
	filterFn := func(v interface{}) bool {
		return v.(int) > 3
	}
	expected := []interface{}{5, 4, 4}
	actual := stream.Filter(filterFn).ToSlice()
	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("expected %v got %v", expected, actual)
	}
}

func TestFold(t *testing.T) {

}

func TestForEach(t *testing.T) {

}
