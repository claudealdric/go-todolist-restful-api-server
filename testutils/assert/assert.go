package assert

import (
	"reflect"
	"slices"
	"testing"
)

func AssertCalls(t testing.TB, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("incorrect number of calls; got %d, want %d", got, want)
	}
}

func AssertContains[T comparable](t testing.TB, slice []T, element T) {
	t.Helper()
	if !slices.Contains(slice, element) {
		t.Errorf("slice should contain %v but doesn't", element)
	}
}

func AssertContentType(t testing.TB, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("response did not have content-type of %q, got %q", want, got)
	}
}

func AssertDoesNotContain[T comparable](t testing.TB, slice []T, element T) {
	t.Helper()
	if slices.Contains(slice, element) {
		t.Errorf("slice should not contain %v but does", element)
	}
}

func AssertNoError(t testing.TB, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("didn't expect an error but got one, %v", err)
	}

}

func AssertEquals[T any](t testing.TB, got, want T) {
	t.Helper()
	switch v := any(got).(type) {
	case string, int, int64, float64, bool:
		if v != any(want) {
			t.Errorf("got %v, want %v", got, want)
		}
	default:
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	}
}

func AssertStatus(t testing.TB, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("got %d, want %d", got, want)
	}
}
