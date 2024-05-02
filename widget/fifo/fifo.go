package fifo

type Fifo struct {
	q []interface{}
}

func NewFifo() *Fifo {
	return &Fifo{q: make([]interface{}, 0)}
}

func (f *Fifo) Reset() {
	f.q = f.q[:0]
}

func (f *Fifo) Push(i interface{}) {
	f.q = append(f.q, i)
}

func (f *Fifo) Pop() interface{} {
	indice := -1
	if len(f.q) >= 1 {
		indice = len(f.q) - 1
	}
	if indice < 0 {
		indice = 0
	}
	res := f.q[indice]
	if indice > 0 {
		f.q = f.q[0:indice]
	}

	return res
}
