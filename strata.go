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

var outputting = false

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

	//Open output file or stdout and create buffer
	var file *os.File
	var writer *bufio.Writer
	if *arg_output != "" {
		var err error
		file, err = os.Create(*arg_output)
		if err != nil {
			panic(err)
		}
		defer file.Close()
		writer = bufio.NewWriter(file)
		defer writer.Flush()
		outputting = true
	}

	msgRaw("Taking %d samples... ", *arg_iterations)
	start := time.Now()

	//Sampling method average declarations
	errorAverages := map[Strategy]*Average{
		SIMPLE_RANDOM: &Average{},
		CONVENIENCE: &Average{},
		STRATIFIED_ROW: &Average{},
		STRATIFIED_COL: &Average{},
	}

	sampleAverages := map[Strategy]*Average{
		SIMPLE_RANDOM: &Average{},
		CONVENIENCE: &Average{},
		STRATIFIED_ROW: &Average{},
		STRATIFIED_COL: &Average{},
	}

	//Get the real average, and save it to the output
	realAvg := Sample(ALL)
	writeRecord(writer, ALL, realAvg.Value())

	convenienceAvg := Sample(CONVENIENCE)
	writeRecord(writer, CONVENIENCE, convenienceAvg.Value())

	for i := 0; i < *arg_iterations; i++ {
		for strat := SIMPLE_RANDOM; strat <= STRATIFIED_COL; strat++ {
			//Get sample avg
			sampleAvg := Sample(strat)

			//Add the sample avg to the sample average
			sampleAverages[strat].Merge(sampleAvg)

			//Calculate margin of error
			errorMargin := 1.0 - realAvg.Value() / sampleAvg.Value()

			//Add the absolute value of the error to the error average
			absError := math.Abs(errorMargin)
			errorAverages[strat].Include(absError)

			//Add the absolute value of the average to the output buffer.
			writeRecord(writer, strat, sampleAvg.Value())
		}
	}

	msg("[DONE in %dms]", time.Since(start).Nanoseconds() / 1000000);
	msg("")

	msg("  --- Result ---")
	msg("Population    : %f avg", realAvg.Value())
	msg("Convenience   : %f avg : %f%% error", convenienceAvg.Value(), math.Abs(1 - convenienceAvg.Value() / realAvg.Value()) * 100)
	msg("Simple random : %f avg : %f%% error", sampleAverages[SIMPLE_RANDOM].Value(), errorAverages[SIMPLE_RANDOM].Value() * 100)
	msg("Stratified row: %f avg : %f%% error", sampleAverages[STRATIFIED_ROW].Value(), errorAverages[STRATIFIED_ROW].Value() * 100)
	msg("Stratified col: %f avg : %f%% error", sampleAverages[STRATIFIED_COL].Value(), errorAverages[STRATIFIED_COL].Value() * 100)
}

func msg(msg string, a ...interface{}){
	msgRaw(msg + "\n", a...)
}

func msgRaw(msg string, a ...interface{}){
	//Writes a formatted message to stderr
	fmt.Fprintf(os.Stderr, msg, a...)
}

func writeRecord(writer *bufio.Writer, strat Strategy, average float64) (int, error){
	if(outputting) {
		return writer.WriteString(fmt.Sprintf("%d,%f\n", int(strat), average))
	}
	return 0, nil
}


