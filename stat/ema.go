package stat

type EMA struct {
	alpha  float64
	score  float64
	nodata bool
}

func NewEMA(a float64) *EMA {
	return &EMA{
		alpha:  a,
		score:  0,
		nodata: true,
	}
}

func (st *EMA) Add(r Record) {
	if st.nodata {
		st.nodata = false
		st.score = r.Value
	} else {
		a := st.alpha
		st.score = a * r.Value + (1-a) * st.score
	}
}

func (st *EMA) Score() float64 {
	return st.score
}

func (st *EMA) Inc(x float64) {
  st.score += x
}

func (st *EMA) Set(x float64) {
	st.score = x
	st.nodata = false
}

func (st *EMA) Reset() {
	st.score = 0
	st.nodata = true
}
