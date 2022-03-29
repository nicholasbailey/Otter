package collections

// An Iterable is a value that supports basic functional iteration
// operations. Iterables do not, in general, guarentee iteration order
// or that iteration will ever terminate.
type Iterable interface {
	ForEach(iterFn func(interface{}))
	Map(mapFn func(interface{}) interface{}) Iterable
	Filter(filterFn func(interface{}) bool) Iterable
	Fold(initialValue interface{}, reducerFn func(interface{}, interface{}) interface{}) interface{}
	Iterator() Iterator
	ToSlice() []interface{}
}

type Iterator interface {
	MoveNext() bool
	Current() interface{}
}

type FiniteIterable interface {
	Iterable
	Size() int
}

// A Sequence is an immutable iterable with a fixed length and order.
// Examples of Sequence types include LinkedLists, SliceSequences
// and Vectors.
type Sequence interface {
	//Iterable
	// Returns the number of elements in the sequence
	Size() int

	Get(index int) interface{}
	Update(index int, value interface{}) Sequence
	Append(value interface{}) Sequence
	//Prepend(value interface{}) Sequence
	//Slice(start int, end int) Sequence
}

// A Set is an immutable iterable with no duplicates. Sets support standard
// set operations. Sets do not in general guarentee iteration
// order, though many implementations do.
type Set interface {
	Iterable
	Size() int
	Contains(value interface{}) bool
	SubsetOf(other Set) bool
	Add(value interface{}) Set
	Remove(value interface{}) Set
	Intersect(other Set) Set
	Union(other Set) Set
	Difference(other Set) Set
}

// A Map is an immutable iterable of Key-Value pairs supporting
// Key-Value store semantics. Maps do not in general guarentee
// iteration order, though many implementations do
type Map interface {
	Iterable
	Contains(key interface{}) bool
	Get(key interface{}) (interface{}, bool)
	Set(key interface{}, value interface{}) Map
	Remove(key interface{}) Map
	Merge(other Map) Map
	Keys() Iterable
	Values() Iterable
	KeySet() Set
}
