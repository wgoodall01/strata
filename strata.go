package main

import (
	"math/rand"
	"time"
	"flag"
	"fmt"
	"math"
)


// Start up a RNG
var Rng = rand.New(rand.NewSource(time.Now().UnixNano()));

func main(){
	// argv parsing to config
	iterations := flag.Int("iterations", 1000000, "Number of iterations, default 1 million.")

	realAvg := Sample(ALL)

	//Sampling method average declarations
	errorAverages := map[Strategy]*Average{
		SIMPLE_RANDOM: &Average{},
		CONVENIENCE: &Average{},
		STRATIFIED_ROW: &Average{},
		STRATIFIED_COL: &Average{},
	}

	for i := 0; i < *iterations; i++ {
		for strat := SIMPLE_RANDOM; strat <= STRATIFIED_COL; strat++ {
			//Get sample avg
			sampleAvg := Sample(strat)

			//Calculate margin of error
			errorMargin := 1.0 - realAvg.Value() / sampleAvg.Value()

			//Add the absolute value of the error to the error average
			absError := math.Abs(errorMargin)
			errorAverages[strat].Include(absError)
		}
	}

	fmt.Printf("Real          : %f\n", realAvg.Value())
	fmt.Printf("Simple random : %f%% error\n", errorAverages[SIMPLE_RANDOM].Value() * 100)
	fmt.Printf("Convenience   : %f%% error\n", errorAverages[CONVENIENCE].Value() * 100)
	fmt.Printf("Stratified row: %f%% error\n", errorAverages[STRATIFIED_ROW].Value() * 100)
	fmt.Printf("Stratified col: %f%% error\n", errorAverages[STRATIFIED_COL].Value() * 100)

}


