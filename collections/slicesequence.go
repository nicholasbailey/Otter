package collections

type SliceSequence struct {
	slice []interface{}
}

func NewSliceSequence(slice []interface{}) *SliceSequence {
	return &SliceSequence{
		slice: slice,
	}
}

type SliceIterator struct {
	slice []interface{}
	index int
}

func NewSliceIterator(slice []interface{}) *SliceIterator {
	return &SliceIterator{
		slice: slice,
		index: -1,
	}
}

func (iterator *SliceIterator) MoveNext() bool {
	iterator.index += 1
	return iterator.index < len(iterator.slice)
}

func (iterator *SliceIterator) Current() interface{} {
	if iterator.index >= len(iterator.slice) {
		panic(ErrIterationOutOfRange)
	}
	return iterator.slice[iterator.index]
}
