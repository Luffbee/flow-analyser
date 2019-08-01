package stat

import (
	"github.com/montanaflynn/stats"
)

type Median struct {
	size  int
	count int
	recs  []float64
}

func NewMedian(sz int) *Median {
	return &Median{
		size:  sz,
		count: 0,
		recs:  make([]float64, sz),
	}
}

func (st *Median) Add(r Record) {
	st.recs[st.count%st.size] = r.Value
	st.count++
}

func (st *Median) Score() float64 {
	r, err := stats.Median(st.recs)
	if err != nil {
		return 0
	}
	return r
}
