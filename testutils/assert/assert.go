package assert

import (
	"errors"
	"reflect"
	"slices"
	"testing"
)

func Calls(t testing.TB, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("incorrect number of calls; got %d, want %d", got, want)
	}
}

func Contains[T comparable](t testing.TB, slice []T, element T) {
	t.Helper()
	if !slices.Contains(slice, element) {
		t.Errorf("slice should contain %v but doesn't", element)
	}
}

func ContentType(t testing.TB, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("response did not have content-type of %q, got %q", want, got)
	}
}

func DoesNotContain[T comparable](t testing.TB, slice []T, element T) {
	t.Helper()
	if slices.Contains(slice, element) {
		t.Errorf("slice should not contain %v but does", element)
	}
}

func ErrorContains(t testing.TB, got, want error) {
	t.Helper()
	if !errors.Is(got, want) {
		t.Errorf("wanted error of type '%v', got '%v'", want, got)
	}
}

func HasError(t testing.TB, err error) {
	t.Helper()
	if err == nil {
		t.Fatal("expected an error but didn't get one")
	}
}

func HasNoError(t testing.TB, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("didn't expect an error but got one, %v", err)
	}
}

func Equals[T any](t testing.TB, got, want T) {
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

func Status(t testing.TB, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("got %d, want %d", got, want)
	}
}
