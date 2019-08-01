package stat

type ActiveAtK struct {
	k      int
	active bool
	recs   []Record
	stat   *EMA
}

func NewActiveAtK(k int, a float64) *ActiveAtK {
	return &ActiveAtK{
		k:    k,
		recs: []Record{},
		stat: NewEMA(a),
	}
}

func (alk *ActiveAtK) Add(r Record) {
	alk.recs = append(alk.recs, r)
  if alk.active {
	  alk.stat.Add(r)
  } else if len(alk.recs) >= alk.k {
    alk.active = true
    for _, r := range alk.recs {
      alk.stat.Add(r)
    }
  }
}

func(alk *ActiveAtK) Score() float64 {
  return alk.stat.Score()
}


func (alk *ActiveAtK) RmBefore(t uint64) {
	for len(alk.recs) > 0 {
		if alk.recs[0].Time >= t {
			break
		}
		alk.recs = alk.recs[1:]
	}
	if alk.active && len(alk.recs) == 0 {
		alk.active = false
		alk.stat.Reset()
	}
}

func (alk *ActiveAtK) ScoreCount() (float64, int) {
	if alk.active {
		return alk.stat.Score(), len(alk.recs)
	}
	return 0, 0
}

