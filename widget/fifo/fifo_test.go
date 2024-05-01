package fifo_test

import (
	"fmt"
	"testing"

	"github.com/jeromelesaux/fyne-io/widget/fifo"
	"github.com/stretchr/testify/assert"
)

func TestFifo(t *testing.T) {

	t.Run("push", func(t *testing.T) {
		f := fifo.NewFifo()
		for i := 0; i <= 100; i++ {
			f.Push(i)
		}
	})

	t.Run("pop", func(t *testing.T) {
		f := fifo.NewFifo()
		for i := 0; i <= 100; i++ {
			f.Push(i)
		}
		for i := 0; i <= 100; i++ {
			f.Pop()
		}
	})
	t.Run("overload_pop", func(t *testing.T) {
		f := fifo.NewFifo()
		for i := 0; i <= 10; i++ {
			f.Push(i)
		}
		for i := 0; i <= 10; i++ {
			f.Pop()
		}

		for i := 0; i <= 5; i++ {
			res := f.Pop()
			assert.Equal(t, 0, res)
		}
	})

	t.Run("stack_structure", func(t *testing.T) {
		f := fifo.NewFifo()
		type bean struct {
			v int
			s string
		}

		for i := 0; i <= 100; i++ {
			b := bean{v: i, s: fmt.Sprintf("%d", i)}
			f.Push(b)
		}
		for i := 100; i >= 0; i-- {
			b := f.Pop()
			assert.Equal(t, i, b.(bean).v)
			assert.Equal(t, fmt.Sprintf("%d", i), b.(bean).s)
		}
	})
}
