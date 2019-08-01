package stat

type ThreeLevel struct {
	size      int
	lowBound  float64
	highBound float64
	low       *ActiveAtK
	mid       *ActiveAtK
	high      *ActiveAtK
}

func NewThreeLevel(sz int, low, high float64, l, m, h int) *ThreeLevel {
	a := 2.0 / (float64(sz)/2 + 1.0)
	return &ThreeLevel{
		size:      sz,
		lowBound:  low,
		highBound: high,
		low:       NewActiveAtK(l, a),
		mid:       NewActiveAtK(m, a),
		high:      NewActiveAtK(h, a),
	}
}

func (st *ThreeLevel) Add(r Record) {
  t := uint64(0)
	if r.Time > uint64(st.size) {
    t = r.Time - uint64(st.size)
  }
  st.low.RmBefore(t)
  st.mid.RmBefore(t)
  st.high.RmBefore(t)

	switch {
	case r.Value < st.lowBound: st.low.Add(r)
	case r.Value < st.highBound: st.mid.Add(r)
	default: st.high.Add(r)
	}
}

func (st *ThreeLevel) Score() float64 {
	var score float64 = 0.0
	var cnt int = 0

  s, c := st.low.ScoreCount()
  score += s * float64(c)
  cnt += c

  s, c = st.mid.ScoreCount()
  score += s * float64(c)
  cnt += c

  s, c = st.high.ScoreCount()
  score += s * float64(c)
  cnt += c

	if cnt == 0 {
		return 0.0
	}
	return score / float64(cnt)
}
