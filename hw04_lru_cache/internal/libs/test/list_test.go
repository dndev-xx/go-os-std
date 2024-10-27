package hw04lrucache

import (
	"testing"

	"github.com/dndev-xx/go-os-std/hw04_lru_cache/internal/libs" //nolint:depguard
	"github.com/stretchr/testify/require"                        //nolint:depguard
)

func TestList(t *testing.T) {
	t.Run("empty list", func(t *testing.T) {
		l := libs.NewList()

		require.Equal(t, 0, l.Len())
		require.Nil(t, l.Front())
		require.Nil(t, l.Back())
	})

	t.Run("push only front", func(t *testing.T) {
		l := libs.NewList()
		l.PushFront(10)
		l.PushFront(20)
		l.PushFront(30)
		require.Equal(t, 3, l.Len())
		require.Equal(t, 30, l.Front().Value)
	})

	t.Run("push front and back", func(t *testing.T) {
		l := libs.NewList()
		l.PushFront(10)
		l.PushBack(20)
		l.PushBack(30)
		l.PushBack(40)
		l.PushFront(50)
		require.Equal(t, 5, l.Len())
		require.Equal(t, 50, l.Front().Value)
		require.Equal(t, 40, l.Back().Value)
		require.Equal(t, 30, l.Back().Prev.Value)
		require.Nil(t, l.Front().Prev)
		require.Nil(t, l.Back().Next)
	})

	t.Run("remove", func(t *testing.T) {
		l := libs.NewList()
		l.PushFront(10)
		l.PushBack(20)
		l.PushBack(30)
		l.PushBack(40)
		l.PushFront(50)
		remove := l.Back().Prev
		l.Remove(remove)
		require.Equal(t, 4, l.Len())
		require.Equal(t, 20, l.Back().Prev.Value)
	})

	t.Run("remove last", func(t *testing.T) {
		l := libs.NewList()
		l.PushFront(10)
		l.PushBack(20)
		l.PushBack(30)
		l.PushBack(40)
		l.PushFront(50)
		remove := l.Back()
		l.Remove(remove)
		require.Equal(t, 4, l.Len())
		require.Equal(t, 30, l.Back().Value)
	})

	t.Run("complex", func(t *testing.T) {
		l := libs.NewList()

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
	})
}
