package collections

type EmptyIterator struct {
}

func (empty *EmptyIterator) MoveNext() bool {
	return false
}

func (empty *EmptyIterator) Current() interface{} {
	panic(ErrIterationOutOfRange)
}

type SequenceIterator struct {
	index    int
	sequence Sequence
}

func NewSequenceIterator(sequence Sequence) Iterator {
	return &SequenceIterator{
		sequence: sequence,
		index:    -1,
	}
}

func (iterator *SequenceIterator) MoveNext() bool {
	iterator.index += 1
	return iterator.index < iterator.sequence.Size()
}

func (iterator *SequenceIterator) Current() interface{} {
	if iterator.index < iterator.sequence.Size() {
		panic(ErrIterationOutOfRange)
	}
	return iterator.sequence.Get(iterator.index)
}

type MapIterator struct {
	baseIterator Iterator
	mapFn        func(interface{}) interface{}
}

func (iterator *MapIterator) MoveNext() bool {
	return iterator.baseIterator.MoveNext()
}

func (iterator *MapIterator) Current() interface{} {
	return iterator.mapFn(iterator.baseIterator.Current())
}

type FilterIterator struct {
	baseIterator Iterator
	filterFn     func(interface{}) bool
}

func (iterator *FilterIterator) MoveNext() bool {
	base := iterator.baseIterator
	for base.MoveNext() {
		current := base.Current()
		if iterator.filterFn(current) {
			return true
		}
	}
	return false
}

func (iterator *FilterIterator) Current() interface{} {
	return iterator.baseIterator.Current()
}
