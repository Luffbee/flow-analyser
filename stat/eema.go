package stat

import "math"
// Y = (1-beta^t) * X + beta^t * Y
type EEMA struct {
	interval float64
	neglect  float64
	beta     float64
	score    float64
  lastTime uint64
}

func NewEEMA(T float64, N int, neg float64) *EEMA {
	beta := math.Pow(0.001, 1/(float64(N)*T))
	return &EEMA{
		interval: T,
		neglect:  neg,
		beta:     beta,
		score:    0,
    lastTime: 0,
	}
}

func (st *EEMA) getUpperBound(t float64) float64 {
  lin := st.score + (st.neglect * t / st.interval)
  exp := st.score * math.Pow(2.0, t / st.interval)
  return math.Max(lin, exp)
}

func (st *EEMA) getLowerBound(t float64) float64 {
  lin := st.score - (st.neglect * t / st.interval)
  exp := st.score * math.Pow(0.5, t / st.interval)
  return math.Min(lin, exp)
}

func (st *EEMA) Add(r Record) {
  if st.lastTime == 0 {
    if r.Value > st.neglect {
      st.score = st.neglect
    } else {
      st.score = r.Value
    }
    st.lastTime = r.Time
  } else if r.Time > st.lastTime {
    t := float64(r.Time - st.lastTime)
    b := math.Pow(st.beta, t)
    y := (1 - b) * r.Value + b * st.score
    if y > st.score {
      st.score = math.Min(st.getUpperBound(t), y)
    } else if y < st.score {
      st.score = math.Max(st.getLowerBound(t), y)
    }
    st.lastTime = r.Time
  }
}

func (st *EEMA) Score() float64 {
	return st.score
}
