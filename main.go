package main

import (
	"flow-analyser/stat"
	"fmt"
	"os"
	"time"
)

// low load: [0, 20)
// middle load: [20, 70)
// high load: [70, 100]

var examples []Data = []Data{
	// startup noise
	DataFromSlice("startup", []float64{100, 1, 2, 3, 4, 5, 6, 1, 1, 0}),
	// noise
	DataFromSlice("noise", []float64{1, 1, 0, 1, 0, 100, 6, 4, 4, 3, 5, 1, 2, 3, 0}),
	// pulse high load with window < 5
	DataFromSlice("pulse high", []float64{0, 1, 0, 1, 100, 4, 5, 100, 10, 1, 2, 100, 3, 3, 3, 100, 1, 1, 1, 1, 1, 0}),
	// continuous high load
	DataFromSlice("cont high", []float64{0, 0, 0, 0, 0, 100, 98, 85, 85, 75, 85, 99, 100, 100, 1, 1, 2, 3, 4, 1, 0}),
	// continuous middle load
	DataFromSlice("cont mid", []float64{1, 0, 0, 1, 0, 25, 30, 17, 21, 44, 55, 61, 30, 40, 18, 20, 1, 2, 3, 1, 1, 2, 0}),
	// continuous middle load with pulse high load
	DataFromSlice("c-m-p-h", []float64{60, 75, 100, 44, 33, 20, 99, 50, 32, 89, 19, 12, 30, 99, 30, 15, 7.5, 2, 1, 1, 0.5, 0, 0}),
}

type Data struct {
	Name string
	Records []stat.Record
}

func DataFromSlice(name string, data []float64) Data {
	recs := make([]stat.Record, len(data))
	for i, v := range data {
		recs[i] = stat.Record{Time: uint64(i), Value: v}
	}
	return Data{name, recs}
}

func ShowData(data Data, prefix string, endl bool) {
	fmt.Printf("%s %10s ", prefix, data.Name)
	for _, v := range data.Records {
		fmt.Printf("%6.2f ", v.Value)
	}
  if endl {
    fmt.Println("")
  }
}

func RunStat(st stat.Stat, data Data) Data {
  out := make([]stat.Record, len(data.Records))
  for i, v := range data.Records {
    st.Add(v)
    out[i].Time = v.Time
    out[i].Value = st.Score()
  }
  return Data{Records: out}
}

func RunDispStat(name string, st stat.Stat, data Data) {
  res := RunStat(st, data)
  res.Name = name
  ShowData(res, "Stat", false)

	var sum float64 = 0.0
  dis := 1 + data.Records[len(data.Records)-1].Time - data.Records[0].Time

	tim := time.Now()
	for i := 0; i < 10; i++ {
		for _, v := range data.Records {
			v.Time += uint64(i + 1) * dis
			st.Add(v)
			sum += st.Score()
		}
	}
	dur := time.Since(tim)

	sum = sum / float64(len(data.Records)*10)

	fmt.Printf("Time %v, %.2f\n", dur, sum)
}

func runExamples(_ []string) {
	for _, data := range examples {
		ShowData(data, "Data", true)
		RunDispStat("Median", stat.NewMedian(5), data)
		RunDispStat("EMA", stat.NewEMA(2.0/(5.0+1.0)), data)
		RunDispStat("3 Level", stat.NewThreeLevel(5, 20, 70, 2, 1, 2), data)
		//RunStat("EEMA", stat.NewEEMA(5, 20, 70), data)
		fmt.Println("")
	}
}

func runStdin(args []string) {
  if len(args) < 1 {
    os.Stderr.WriteString("which stat?")
  }

  var n int
  fmt.Scanf("%d", &n)
  data := make([]stat.Record, n)
  for i := 0; i < n; i++ {
    fmt.Scanf("%d %f", &data[i].Time, &data[i].Value)
  }

  stat := stat.NewEEMA(60, 5, 5 * 1024 * 1024)
  //stat := stat.NewEMA(2.0/6.0)
  res := RunStat(stat, Data{Records: data})
  for _, r := range res.Records {
    fmt.Printf("%v %v\n", r.Time, r.Value)
  }
}

type RunFunc func([]string)

var run map[string]RunFunc = map[string]RunFunc {
  "examples": runExamples,
  "stdin": runStdin,
}

func main() {
	if len(os.Args) < 2 {
		os.Stderr.WriteString("which function?\n")
    for s, _ := range run {
      os.Stderr.WriteString("  " + s + "\n")
    }
		os.Exit(1)
	}
  if f, ok := run[os.Args[1]]; ok {
    f(os.Args[2:])
  } else {
    os.Stderr.WriteString("no such function, current functions:\n")
    for s, _ := range run {
      os.Stderr.WriteString("  " + s + "\n")
    }
  }
}
