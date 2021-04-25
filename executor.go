package main

import (
	"flag"
	"fmt"
	"github.com/yaricom/goNEAT/v2/experiment"
	"github.com/yaricom/goNEAT/v2/experiments/pole"
	"github.com/yaricom/goNEAT/v2/experiments/xor"
	"github.com/yaricom/goNEAT/v2/neat"
	"github.com/yaricom/goNEAT/v2/neat/genetics"
	"log"
	"math/rand"
	"os"
	"time"
)

// The experiment runner boilerplate code
func main() {
	var outDirPath = flag.String("out", "./out", "The output directory to store results.")
	var contextPath = flag.String("context", "./data/xor.neat", "The execution context configuration file.")
	var genomePath = flag.String("genome", "./data/xorstartgenes", "The seed genome to start with.")
	var experimentName = flag.String("experiment", "XOR", "The name of experiment to run. [XOR, cart_pole, cart_2pole_markov, cart_2pole_non-markov]")
	var trialsCount = flag.Int("trials", 0, "The numbar of trials for experiment. Overrides the one set in configuration.")
	var logLevel = flag.String("log_level", "", "The logger level to be used. Overrides the one set in configuration.")

	flag.Parse()

	// Seed the random-number generator with current time so that
	// the numbers will be different every time we run.
	rand.Seed(time.Now().Unix())

	// Load neatOptions configuration
	configFile, err := os.Open(*contextPath)
	if err != nil {
		log.Fatal("Failed to open context configuration file: ", err)
	}
	neatOptions, err := neat.LoadNeatOptions(configFile)
	if err != nil {
		log.Fatal("Failed to load NEAT options: ", err)
	}

	// Load Genome
	log.Printf("Loading start genome for %s experiment\n", *experimentName)
	genomeFile, err := os.Open(*genomePath)
	if err != nil {
		log.Fatal("Failed to open genome file: ", err)
	}
	startGenome, err := genetics.ReadGenome(genomeFile, 1)
	if err != nil {
		log.Fatal("Failed to read start genome: ", err)
	}
	fmt.Println(startGenome)

	// Check if output dir exists
	outDir := *outDirPath
	if _, err := os.Stat(outDir); err == nil {
		// backup it
		backUpDir := fmt.Sprintf("%s-%s", outDir, time.Now().Format("2006-01-02T15_04_05"))
		// clear it
		err = os.Rename(outDir, backUpDir)
		if err != nil {
			log.Fatal("Failed to do previous results backup: ", err)
		}
	}
	// create output dir
	err = os.MkdirAll(outDir, os.ModePerm)
	if err != nil {
		log.Fatal("Failed to create output directory: ", err)
	}

	// Override neatOptions configuration parameters with ones set from command line
	if *trialsCount > 0 {
		neatOptions.NumRuns = *trialsCount
	}
	if len(*logLevel) > 0 {
		neat.LogLevel = neat.LoggerLevel(*logLevel)
	}

	// The 100 generation XOR experiment
	expt := experiment.Experiment{
		Id:     0,
		Trials: make(experiment.Trials, neatOptions.NumRuns),
	}
	var generationEvaluator experiment.GenerationEvaluator
	switch *experimentName {
	case "XOR":
		expt.MaxFitnessScore = 16.0 // as given by fitness function definition
		generationEvaluator = xor.NewXORGenerationEvaluator(outDir)
	case "cart_pole":
		expt.MaxFitnessScore = 1.0 // as given by fitness function definition
		generationEvaluator = pole.NewCartPoleGenerationEvaluator(outDir, true, 500000)
	case "cart_2pole_markov":
		expt.MaxFitnessScore = 1.0 // as given by fitness function definition
		generationEvaluator = pole.NewCartDoublePoleGenerationEvaluator(outDir, true, pole.ContinuousAction)
	case "cart_2pole_non-markov":
		generationEvaluator = pole.NewCartDoublePoleGenerationEvaluator(outDir, false, pole.ContinuousAction)
	default:
		log.Fatalf("Unsupported experiment: %s", *experimentName)
	}

	if err = expt.Execute(neatOptions, startGenome, generationEvaluator); err != nil {
		log.Fatal("Failed to perform XOR experiment: ", err)
	}

	// Print statistics
	expt.PrintStatistics()

	fmt.Printf(">>> Start genome file:  %s\n", *genomePath)
	fmt.Printf(">>> Configuration file: %s\n", *contextPath)

	// Save experiment data
	expResPath := fmt.Sprintf("%s/%s.dat", outDir, *experimentName)
	if expResFile, err := os.Create(expResPath); err == nil {
		log.Fatal("Failed to create file for experiment results", err)
	} else if err = expt.Write(expResFile); err != nil {
		log.Fatal("Failed to save experiment results", err)
	}
}
