package main

type Average struct{
	Sum   float64
	Count int
}

func (a *Average) Include(x float64){
	a.Sum += x
	a.Count += 1
}

func (a Average) Value() float64{
	return float64(a.Sum) / float64(a.Count)
}

func (a *Average) Merge(other Average){
	a.Sum += other.Sum
	a.Count += other.Count
}