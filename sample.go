package main

type Strategy int
const (
	ALL Strategy = iota
	CONVENIENCE Strategy = iota
	SIMPLE_RANDOM Strategy = iota
	STRATIFIED_ROW Strategy = iota
	STRATIFIED_COL Strategy = iota
)

func Sample(strat Strategy) Average{
	avg := Average{}

	addPlot := func(x int, y int){
		avg.Include(float64(Farm[x][y]))
	}

	switch strat{
	case ALL:
		for i := range Farm {
			for j := range Farm[i]{
				addPlot(i, j);
			}
		}

	case SIMPLE_RANDOM:
		for i := 0; i < 10; i++ {
			x := Rng.Intn(10)
			y := Rng.Intn(10)
			addPlot(x, y)
		}

	case CONVENIENCE:
		for i := 0; i < 10; i++ {
			addPlot(i, 0)
		}

	case STRATIFIED_ROW:
		for i := 0; i < 10; i++ {
			x := Rng.Intn(10)
			addPlot(i, x)
		}

	case STRATIFIED_COL:
		for i := 0; i < 10; i++ {
			x := Rng.Intn(10)
			addPlot(x, i)
		}
	}

	//Return average
	return avg
}
