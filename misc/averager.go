package misc

import (
	"fmt"

	"github.com/quillaja/goutil/functional"
)

type Averager struct {
	Name     string
	Values   []float64
	insIndex int
}

func NewAverager(name string, len int) *Averager {
	return &Averager{
		Name:   name,
		Values: make([]float64, len),
	}
}

func (a *Averager) Add(val float64) {
	if a.insIndex == len(a.Values) {
		fmt.Printf("%s: %g\n", a.Name, functional.ReduceF(a.Values, functional.Sum)/float64(len(a.Values)))
		a.insIndex = 0
	}
	a.Values[a.insIndex] = val
	a.insIndex++
}
