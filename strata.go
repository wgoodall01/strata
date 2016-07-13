package main

import (
	"math/rand"
	"time"
	"bufio"
	"os"
	"fmt"
	"math"
	"github.com/alecthomas/kingpin"
)


// Start up a RNG
var Rng = rand.New(rand.NewSource(time.Now().UnixNano()));

func main(){
	// argv parsing to config
	arg_iterations := kingpin.Flag("iterations", "Number of iterations").
		Default("10").
		Short('i').
		Int()

	arg_output := kingpin.Flag("output", "Output file location").
		Default(""). // If output is blank, send to stdout.
		Short('o').
		String()
	kingpin.Parse()


	//Sampling method average declarations
	errorAverages := map[Strategy]*Average{
		SIMPLE_RANDOM: &Average{},
		CONVENIENCE: &Average{},
		STRATIFIED_ROW: &Average{},
		STRATIFIED_COL: &Average{},
	}

	msgRaw("Taking %d samples... ", *arg_iterations)
	start := time.Now()

	//Open output file or stdout and create buffer
	file := os.Stdout
	if *arg_output != "" {
		var err error
		file, err = os.Create(*arg_output)
		if err != nil {
			panic(err)
		}
		defer file.Close()
	}
	writer := bufio.NewWriter(file)
	defer writer.Flush()

	//Get the real average, and save it to the output
	realAvg := Sample(ALL)
	writeRecord(writer, ALL, realAvg.Value())

	for i := 0; i < *arg_iterations; i++ {
		for strat := SIMPLE_RANDOM; strat <= STRATIFIED_COL; strat++ {
			//Get sample avg
			sampleAvg := Sample(strat)

			//Calculate margin of error
			errorMargin := 1.0 - realAvg.Value() / sampleAvg.Value()

			//Add the absolute value of the error to the error average
			absError := math.Abs(errorMargin)
			errorAverages[strat].Include(absError)

			//Add the absolute value of the error to the output buffer.
			writeRecord(writer, strat, sampleAvg.Value())
		}
	}

	msg("[DONE in %dms]", time.Since(start).Nanoseconds() / 1000000);
	msg("")

	msg("  --- Error Digest ---")
	msg("Real          : %f avg", realAvg.Value())
	msg("Simple random : %f%% error", errorAverages[SIMPLE_RANDOM].Value() * 100)
	msg("Convenience   : %f%% error", errorAverages[CONVENIENCE].Value() * 100)
	msg("Stratified row: %f%% error", errorAverages[STRATIFIED_ROW].Value() * 100)
	msg("Stratified col: %f%% error", errorAverages[STRATIFIED_COL].Value() * 100)
}

func msg(msg string, a ...interface{}){
	msgRaw(msg + "\n", a...)
}

func msgRaw(msg string, a ...interface{}){
	//Writes a formatted message to stderr
	fmt.Fprintf(os.Stderr, msg, a...)
}

func writeRecord(writer *bufio.Writer, strat Strategy, average float64) (int, error){
	return writer.WriteString(fmt.Sprintf("%d,%f\n", int(strat), average))
}


