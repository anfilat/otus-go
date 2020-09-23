package hw04_lru_cache //nolint:golint,stylecheck

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestList(t *testing.T) {
	t.Run("empty list", func(t *testing.T) {
		l := NewList()

		require.Equal(t, 0, l.Len())
		require.Nil(t, l.Front())
		require.Nil(t, l.Back())
	})

	t.Run("complex", func(t *testing.T) {
		l := NewList()

		l.PushFront(10) // [10]
		l.PushBack(20)  // [10, 20]
		l.PushBack(30)  // [10, 20, 30]
		require.Equal(t, 3, l.Len())

		middle := l.Front().Next // 20
		l.Remove(middle)         // [10, 30]
		require.Equal(t, 2, l.Len())

		for i, v := range [...]int{40, 50, 60, 70, 80} {
			if i%2 == 0 {
				l.PushFront(v)
			} else {
				l.PushBack(v)
			}
		} // [80, 60, 40, 10, 30, 50, 70]

		require.Equal(t, 7, l.Len())
		require.Equal(t, 80, l.Front().Value)
		require.Equal(t, 70, l.Back().Value)

		l.MoveToFront(l.Front()) // [80, 60, 40, 10, 30, 50, 70]
		l.MoveToFront(l.Back())  // [70, 80, 60, 40, 10, 30, 50]

		elems := make([]int, 0, l.Len())
		for i := l.Front(); i != nil; i = i.Next {
			elems = append(elems, i.Value.(int))
		}
		require.Equal(t, []int{70, 80, 60, 40, 10, 30, 50}, elems)
	})

	t.Run("push front to empty list", func(t *testing.T) {
		l := NewList()

		l.PushFront(10)

		require.Equal(t, 1, l.Len())
		require.Equal(t, 10, l.Front().Value)
		require.Equal(t, 10, l.Back().Value)
	})

	t.Run("push back to empty list", func(t *testing.T) {
		l := NewList()

		l.PushBack(10)

		require.Equal(t, 1, l.Len())
		require.Equal(t, 10, l.Front().Value)
		require.Equal(t, 10, l.Back().Value)
	})

	t.Run("remove first item", func(t *testing.T) {
		l := NewList()

		l.PushBack(10)
		l.PushBack(20)
		l.PushBack(30)
		l.PushBack(40)

		require.Equal(t, 10, l.Front().Value)
		l.Remove(l.Front())
		require.Equal(t, 20, l.Front().Value)

		elems := make([]int, 0, l.Len())
		for i := l.Front(); i != nil; i = i.Next {
			elems = append(elems, i.Value.(int))
		}
		require.Equal(t, []int{20, 30, 40}, elems)
	})

	t.Run("remove last item", func(t *testing.T) {
		l := NewList()

		l.PushBack(10)
		l.PushBack(20)
		l.PushBack(30)
		l.PushBack(40)

		require.Equal(t, 40, l.Back().Value)
		l.Remove(l.Back())
		require.Equal(t, 30, l.Back().Value)

		elems := make([]int, 0, l.Len())
		for i := l.Front(); i != nil; i = i.Next {
			elems = append(elems, i.Value.(int))
		}
		require.Equal(t, []int{10, 20, 30}, elems)
	})

	t.Run("remove only item", func(t *testing.T) {
		l := NewList()

		item := l.PushFront(10)
		l.Remove(item)

		require.Equal(t, 0, l.Len())
		require.Nil(t, l.Front())
		require.Nil(t, l.Back())
	})

	t.Run("swap first and second items", func(t *testing.T) {
		l := NewList()

		l.PushBack(10)
		l.PushBack(20)

		require.Equal(t, 20, l.Front().Next.Value)
		require.Equal(t, 10, l.Back().Prev.Value)

		l.MoveToFront(l.Back())

		require.Equal(t, 10, l.Front().Next.Value)
		require.Equal(t, 20, l.Back().Prev.Value)

		elems := make([]int, 0, l.Len())
		for i := l.Front(); i != nil; i = i.Next {
			elems = append(elems, i.Value.(int))
		}
		require.Equal(t, []int{20, 10}, elems)
	})
}
