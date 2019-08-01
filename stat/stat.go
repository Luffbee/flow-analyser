package stat

type Record struct {
	Time  uint64
	Value float64
}

type Stat interface {
	Add(r Record)
	Score() float64
}
