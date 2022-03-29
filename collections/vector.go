package collections

// This is a (more or less) straight port of the Vector implementation used
// in Scala's standard library
// [https://github.com/scala/scala/blob/2.13.x/src/library/scala/collection/immutable/Vector.scala]
// It's not guarenteed to be nearly as optimized, because Go and the JVM have different
// runtime characteristics, but you should be able to expect good performance for:
// Head and Tail access (amoritized O(1))
// Random access and functional update (O(log n))
// Append and prepend (amortirzed O(1))
//
// Along with excellent memory sharing.

const bits = 5
const width = 32
const bits2 = 10
const width2 = 1024 // 1 << 19
const bits3 = 15
const width3 = 32768 // 1 << 15
const bits4 = 20
const width4 = 1048576
const bits5 = 25
const width5 = 33554432
const bits6 = 30
const width6 = 1073741824
const mask = 31
const lastWidth = 64
const log2ConcatFactor = 5

const MaxVectorSize = width2

type Arr1 []interface{}

type Arr2 [][]interface{}

type Vector interface {
	Sequence
}

type Vector0 struct{}

func EmptyVector() Vector {
	return &Vector0{}
}

func (vector0 *Vector0) Get(index int) interface{} {
	panic(ErrVectorIndexOutOfRange)
}

func (vector0 *Vector0) Append(value interface{}) Sequence {
	return &Vector1{
		data: Arr1{value},
	}
}

func (vector0 *Vector0) Update(index int, value interface{}) Sequence {
	panic(ErrVectorIndexOutOfRange)
}

func (vector0 *Vector0) String() string {
	return "<>"
}

func (vector0 *Vector0) Size() int {
	return 0
}

func (vector *Vector0) Iterator() Iterator {
	return &EmptyIterator{}
}

func (vector0 *Vector0) Map(mapFn func(interface{}) interface{}) Iterable {
	return mapHelper(vector0, mapFn)
}

func (vector0 *Vector0) Filter(filterFn func(interface{}) bool) Iterable {
	return filterHelper(vector0, filterFn)
}

func (vector0 *Vector0) Fold(initialValue interface{}, reducerFn func(interface{}, interface{}) interface{}) interface{} {
	return foldHelper(vector0, initialValue, reducerFn)
}

func (vector0 *Vector0) ForEach(iterFn func(interface{})) {
	forEachHelper(vector0, iterFn)
}

func (vector *Vector0) ToSlice() []interface{} {
	return toSliceHelper(vector)
}

type Vector1 struct {
	data Arr1
}

func (vector1 *Vector1) Get(index int) interface{} {
	if index < 0 || index >= len(vector1.data) {
		panic(ErrVectorIndexOutOfRange)
	}
	return vector1.data[index]
}

func (vector1 *Vector1) Append(value interface{}) Sequence {
	length := len(vector1.data)
	if length < width {
		return &Vector1{
			data: append(vector1.data, value),
		}
	}

	return &Vector2{
		prefix1:       vector1.data,
		data2:         Arr2{},
		suffix1:       Arr1{value},
		length:        width + 1,
		prefix1Length: width,
	}
}

func (vector1 *Vector1) Size() int {
	return len(vector1.data)
}

func (vector1 *Vector1) Update(index int, value interface{}) Sequence {
	length := len(vector1.data)
	if index < 0 || index >= length {
		panic(ErrVectorIndexOutOfRange)
	}
	newSlice := make([]interface{}, length)
	copy(newSlice, vector1.data)
	newSlice[index] = value
	return &Vector1{
		data: newSlice,
	}
}

func (vector *Vector1) Iterator() Iterator {
	return NewSequenceIterator(vector)
}

func (vector *Vector1) Map(mapFn func(interface{}) interface{}) Iterable {
	return mapHelper(vector, mapFn)
}

func (vector *Vector1) Filter(filterFn func(interface{}) bool) Iterable {
	return filterHelper(vector, filterFn)
}

func (vector *Vector1) Fold(initialValue interface{}, reducerFn func(interface{}, interface{}) interface{}) interface{} {
	return foldHelper(vector, initialValue, reducerFn)
}

func (vector *Vector1) ForEach(iterFn func(interface{})) {
	forEachHelper(vector, iterFn)
}

func (vector *Vector1) ToSlice() []interface{} {
	return toSliceHelper(vector)
}

type Vector2 struct {
	prefix1       Arr1
	data2         Arr2
	suffix1       Arr1
	length        int
	prefix1Length int
}

func (vector2 *Vector2) Get(index int) interface{} {
	if index < 0 || index >= vector2.length {
		panic(ErrVectorIndexOutOfRange)
	}
	if index < vector2.prefix1Length {
		return vector2.prefix1[index]
	}
	i0 := index - vector2.prefix1Length
	i1 := i0 & mask
	i2 := i0 >> bits
	if i2 < len(vector2.data2) {
		return vector2.data2[i2][i1]
	}
	return vector2.suffix1[i1]
}

func (vector2 *Vector2) Append(value interface{}) Sequence {
	if len(vector2.suffix1) < width {
		newSuffix1 := append(vector2.suffix1, value)
		return &Vector2{
			prefix1:       vector2.prefix1,
			data2:         vector2.data2,
			suffix1:       newSuffix1,
			length:        vector2.length + 1,
			prefix1Length: vector2.prefix1Length,
		}
	}
	if len(vector2.data2) < width-2 {
		newData2 := append(vector2.data2, vector2.suffix1)
		newSuffix1 := Arr1{value}
		return &Vector2{
			prefix1:       vector2.prefix1,
			data2:         newData2,
			suffix1:       newSuffix1,
			length:        vector2.length + 1,
			prefix1Length: vector2.prefix1Length,
		}
	}
	panic(ErrVectorTooLarge)
}

func (vector2 *Vector2) Size() int {
	return vector2.length
}

func (vector2 *Vector2) Update(index int, value interface{}) Sequence {
	// Handle index out of range
	if index < 0 || index > vector2.length {
		panic(ErrVectorIndexOutOfRange)
	}

	// Handle index in prefix
	if index < vector2.prefix1Length {
		newPrefix1 := make(Arr1, vector2.prefix1Length)
		copy(newPrefix1, vector2.prefix1)
		newPrefix1[index] = value
		return &Vector2{
			prefix1:       newPrefix1,
			data2:         vector2.data2,
			suffix1:       vector2.suffix1,
			length:        vector2.length,
			prefix1Length: vector2.prefix1Length,
		}
	}

	i0 := index - vector2.prefix1Length
	i1 := i0 & mask
	i2 := i0 >> bits

	data2Length := len(vector2.data2)
	if i2 < data2Length {
		newData2 := make(Arr2, data2Length)
		copy(newData2, vector2.data2)

		sliceToModify := newData2[i2]
		modifiedSlice := make(Arr1, len(sliceToModify))
		copy(modifiedSlice, sliceToModify)

		modifiedSlice[i1] = value
		newData2[i2] = modifiedSlice
		return &Vector2{
			prefix1:       vector2.prefix1,
			data2:         newData2,
			suffix1:       vector2.suffix1,
			length:        vector2.length,
			prefix1Length: vector2.prefix1Length,
		}
	}

	newSuffix1 := make(Arr1, len(vector2.suffix1))
	copy(newSuffix1, vector2.suffix1)
	newSuffix1[i1] = value
	return &Vector2{
		prefix1:       vector2.prefix1,
		data2:         vector2.data2,
		suffix1:       newSuffix1,
		length:        vector2.length,
		prefix1Length: vector2.prefix1Length,
	}
}

func (vector *Vector2) Iterator() Iterator {
	return NewSequenceIterator(vector)
}

func (vector *Vector2) Map(mapFn func(interface{}) interface{}) Iterable {
	return mapHelper(vector, mapFn)
}

func (vector *Vector2) Filter(filterFn func(interface{}) bool) Iterable {
	return filterHelper(vector, filterFn)
}

func (vector *Vector2) Fold(initialValue interface{}, reducerFn func(interface{}, interface{}) interface{}) interface{} {
	return foldHelper(vector, initialValue, reducerFn)
}

func (vector *Vector2) ForEach(iterFn func(interface{})) {
	forEachHelper(vector, iterFn)
}

func (vector *Vector2) ToSlice() []interface{} {
	return toSliceHelper(vector)
}
