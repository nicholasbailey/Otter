package collections

// Helper functions to make it easier to implement the Iterable
// iterface with consistent logic

func forEachHelper(iterable Iterable, iterFn func(interface{})) {
	iterator := iterable.Iterator()
	for iterator.MoveNext() {
		iterFn(iterator.Current())
	}
}

func mapHelper(iterable Iterable, mapFn func(interface{}) interface{}) Iterable {
	oldIterator := iterable.Iterator()
	newIterator := &MapIterator{
		baseIterator: oldIterator,
		mapFn:        mapFn,
	}
	return &Stream{
		iterator: newIterator,
	}
}

func filterHelper(iterable Iterable, filterFn func(interface{}) bool) Iterable {
	oldIterator := iterable.Iterator()
	newIterator := &FilterIterator{
		baseIterator: oldIterator,
		filterFn:     filterFn,
	}
	return &Stream{
		iterator: newIterator,
	}
}

func foldHelper(iterable Iterable, initialValue interface{}, reducerFn func(interface{}, interface{}) interface{}) interface{} {
	val := initialValue
	iterator := iterable.Iterator()
	for iterator.MoveNext() {
		val = reducerFn(val, iterator.Current())
	}
	return val
}

func toSliceHelper(iterable Iterable) []interface{} {
	slice := []interface{}{}
	iterator := iterable.Iterator()
	for iterator.MoveNext() {
		slice = append(slice, iterator.Current())
	}
	return slice
}
