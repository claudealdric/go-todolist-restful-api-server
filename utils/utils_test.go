package utils

import (
	"testing"

	"github.com/claudealdric/go-todolist-restful-api-server/testutils/assert"
)

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
