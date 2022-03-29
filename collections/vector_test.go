package collections

import (
	"errors"
	"fmt"
	"testing"
)

// TODO - break up these tests into little tests

func shouldPanic(t *testing.T, f func(), expectedError error) {
	t.Helper()
	defer func() {
		err := recover()
		if err == nil {
			t.Fatalf("should have paniced")
		}
		if !errors.Is(err.(error), expectedError) {
			t.Fatalf("Failed with incorrect error. Expected %v, got %v", expectedError, err)
		}
	}()
	f()
}

func TestEmptyVectorLength(t *testing.T) {
	vector := EmptyVector()
	if vector.Size() != 0 {
		t.Fatal("Vector length not equal to zero")
	}
}

func buildVectorSlice() []Vector {
	vectors := []Vector{}
	vector := EmptyVector()
	vectors = append(vectors, vector)
	for i := 0; i < MaxVectorSize; i++ {
		vector = vector.Append(i)
		vectors = append(vectors, vector)
	}
	return vectors
}

func TestVectorAppend(t *testing.T) {
	vectors := buildVectorSlice()
	for i := 0; i < len(vectors); i++ {
		testVector := vectors[i]
		if testVector.Size() != i {
			t.Fatalf("expected vector %v to have length %v, but got %v", testVector, i, testVector.Size())
		}
		for j := 0; j < i; j++ {
			val := testVector.Get(j)
			if val != j {
				t.Fatalf("expected %v at index %v of vector %v, got %v", j, j, testVector, val)
			}
		}
	}
}

func TestVectorUpdate(t *testing.T) {
	vectors := buildVectorSlice()

	for i := 0; i < len(vectors); i++ {
		originalVector := vectors[i]
		updatedVector := originalVector
		for j := 0; j < i; j++ {
			updatedVector = updatedVector.Update(j, j*2)

		}
		for j := 0; j < i; j++ {
			updatedVal := updatedVector.Get(j)

			if updatedVal != j*2 {
				t.Fatalf("expected vector to be updated to %v at index %v but was %v", j*2, j, updatedVal)
			}

			originalVal := originalVector.Get(j)

			if originalVal != j {
				t.Fatalf("expected vector not to be mutated, but was different at index %v, expected %v, got %v", j, j, originalVal)
			}
		}
	}
}

func TestVectorIndexOutOfRange(t *testing.T) {
	vectors := buildVectorSlice()

	for i := 0; i < len(vectors); i++ {
		vector := vectors[i]
		shouldPanic(t, func() { vector.Get(-1) }, ErrVectorIndexOutOfRange)

		shouldPanic(t, func() { vector.Get(i) }, ErrVectorIndexOutOfRange)

	}
}

func TestVector1(t *testing.T) {
	vector := EmptyVector()
	for i := 0; i < 32; i++ {
		vector = vector.Append(i)
		for j := 0; j <= i; j++ {
			val := vector.Get(j)

			if val != j {
				t.Fatalf("Unexpected vector behavior in vector of length %v at index %v, got %v", i, j, val)
			}
		}
	}
	for i := 0; i < 32; i++ {
		newVector := vector.Update(i, i*2)

		newVal := newVector.Get(i)

		if newVal != i*2 {
			t.Fatalf("Expected %v at index %v after update, got %v", i*2, i, newVal)
		}
		oldVal := vector.Get(i)

		if oldVal != i {
			t.Fatalf("Immutable vector was secretly updated")
		}
	}
}

func TestVector2(t *testing.T) {
	var vector Vector = EmptyVector()
	for i := 0; i < 1024; i++ {
		vector = vector.Append(i)

		for j := 0; j <= i; j++ {
			val := vector.Get(j)

			if val != j {
				fmt.Printf("%v\n", vector)
				t.Fatalf("Unexpected vector behavior in vector %v of length %v at index %v, got %v", vector, i, j, val)
			}
		}
	}

	for i := 0; i < 1024; i++ {
		newVector := vector.Update(i, i*2)

		newVal := newVector.Get(i)

		if newVal != i*2 {
			t.Fatalf("Expected %v at index %v after update, got %v", i*2, i, newVal)
		}
		oldVal := vector.Get(i)

		if oldVal != i {
			t.Fatalf("Immutable vector was secretly updated")
		}
	}

	shouldPanic(t, func() { vector.Append(1024) }, ErrVectorTooLarge)
}
