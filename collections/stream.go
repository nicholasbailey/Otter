package collections

type Stream struct {
	iterator Iterator
}

func NewStream(iterator Iterator) *Stream {
	return &Stream{
		iterator: iterator,
	}
}

func (iterable *Stream) Iterator() Iterator {
	return iterable.iterator
}

func (iterable *Stream) ForEach(iterFn func(interface{})) {
	forEachHelper(iterable, iterFn)
}

func (iterable *Stream) Map(mapFn func(interface{}) interface{}) Iterable {
	return mapHelper(iterable, mapFn)
}

func (iterable *Stream) Filter(filterFn func(interface{}) bool) Iterable {
	return filterHelper(iterable, filterFn)
}

func (iterable *Stream) Fold(initialValue interface{}, reducerFn func(interface{}, interface{}) interface{}) interface{} {
	return foldHelper(iterable, initialValue, reducerFn)
}

func (iterable *Stream) ToSlice() []interface{} {
	return toSliceHelper(iterable)
}
