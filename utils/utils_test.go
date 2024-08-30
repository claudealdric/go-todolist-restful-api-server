package utils

import (
	"testing"

	"github.com/claudealdric/go-todolist-restful-api-server/testutils/assert"
)

func TestConvertToJSON(t *testing.T) {
	t.Run("converts a struct to a JSON object", func(t *testing.T) {
		type node struct {
			Id int `json:"id"`
		}
		n := node{1}
		json, err := ConvertToJSON(n)
		want := `{"id":1}`

		assert.HasNoError(t, err)
		assert.Equals(t, string(json), want)
	})

	t.Run("converts a slice of structs to an array of JSON objects", func(t *testing.T) {
		type node struct {
			Id int `json:"id"`
		}
		s := []node{node{1}, node{2}}
		json, err := ConvertToJSON(s)
		want := `[{"id":1},{"id":2}]`

		assert.HasNoError(t, err)
		assert.Equals(t, string(json), want)
	})

	t.Run("returns a nil value and an error with an invalid JSON string", func(t *testing.T) {
		invalidJsonInput := func() {}
		json, err := ConvertToJSON(invalidJsonInput)
		assert.HasError(t, err)
		assert.Equals(t, json, nil)
	})
}

func TestSliceFind(t *testing.T) {
	t.Run("returns the searched element if it exists", func(t *testing.T) {
		s := []int{1, 2, 3, 4, 5}
		want := s[2]
		got, ok := SliceFind(s, func(i int) bool {
			return i == want
		})
		assert.Equals(t, ok, true)
		assert.Equals(t, got, want)
	})

	t.Run("returns the zero value if it does not exist", func(t *testing.T) {
		t.Run("primitives", func(t *testing.T) {
			ints := []int{1, 2, 3, 4, 5}
			got, ok := SliceFind(ints, func(i int) bool {
				return i == -1
			})
			assert.Equals(t, ok, false)
			assert.Equals(t, got, 0)
		})

		t.Run("structs", func(t *testing.T) {
			type s struct {
				Id int
			}
			structs := []s{
				{Id: 1},
			}
			got, ok := SliceFind(structs, func(s s) bool {
				return s.Id == -1
			})
			var zeroVal s
			assert.Equals(t, ok, false)
			assert.Equals(t, got, zeroVal)
		})
	})
}
