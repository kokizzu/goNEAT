package xor

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yaricom/goNEAT/v2/experiments"
	"github.com/yaricom/goNEAT/v2/experiments/utils"
	"github.com/yaricom/goNEAT/v2/neat"
	"math/rand"
	"testing"
	"time"
)

// The integration test running over multiple iterations in order to detect if any random errors occur.
func TestXOR(t *testing.T) {
	// the numbers will be different every time we run.
	rand.Seed(time.Now().Unix())

	outDirPath, contextPath, genomePath := "../../out/XOR_test", "../../data/xor.neat", "../../data/xorstartgenes"

	// Load Genome
	fmt.Println("Loading start genome for XOR experiment")
	context, startGenome, err := utils.LoadContextAndGenome(contextPath, genomePath)
	neat.LogLevel = neat.LogLevelInfo
	require.NoError(t, err)

	// Check if output dir exists
	err = utils.CreateOutputDir(outDirPath)
	require.NoError(t, err, "Failed to create output directory")

	// The 100 runs XOR experiment
	context.NumRuns = 100
	experiment := experiments.Experiment{
		Id:     0,
		Trials: make(experiments.Trials, context.NumRuns),
	}
	err = experiment.Execute(context, startGenome, NewXORGenerationEvaluator(outDirPath))
	require.NoError(t, err, "Failed to perform XOR experiment")

	// Find winner statistics
	avgNodes, avgGenes, avgEvals, _ := experiment.AvgWinner()

	// check results
	if avgNodes < 5 {
		t.Error("avg_nodes < 5", avgNodes)
	} else if avgNodes > 15 {
		t.Error("avg_nodes > 15", avgNodes)
	}

	if avgGenes < 7 {
		t.Error("avg_genes < 7", avgGenes)
	} else if avgGenes > 20 {
		t.Error("avg_genes > 20", avgGenes)
	}

	maxEvals := float64(context.PopSize * context.NumGenerations)
	assert.True(t, avgEvals < maxEvals)

	t.Logf("avg_nodes: %.1f, avg_genes: %.1f, avg_evals: %.1f\n", avgNodes, avgGenes, avgEvals)
	meanComplexity, meanDiversity, meanAge := 0.0, 0.0, 0.0
	for _, t := range experiment.Trials {
		meanComplexity += t.BestComplexity().Mean()
		meanDiversity += t.Diversity().Mean()
		meanAge += t.BestAge().Mean()
	}
	count := float64(len(experiment.Trials))
	meanComplexity /= count
	meanDiversity /= count
	meanAge /= count
	t.Logf("Mean best organisms: complexity=%.1f, diversity=%.1f, age=%.1f", meanComplexity, meanDiversity, meanAge)
}

// The XOR integration test for disconnected inputs running over multiple iterations in order to detect if any random errors occur.
func TestXOR_disconnected(t *testing.T) {
	// the numbers will be different every time we run.
	rand.Seed(time.Now().Unix())

	outDirPath, contextPath, genomePath := "../../out/XOR_disconnected_test", "../../data/xor.neat", "../../data/xordisconnectedstartgenes"

	fmt.Println("Loading start genome for XOR disconnected experiment")
	context, startGenome, err := utils.LoadContextAndGenome(contextPath, genomePath)
	neat.LogLevel = neat.LogLevelInfo
	require.NoError(t, err)

	// Check if output dir exists
	err = utils.CreateOutputDir(outDirPath)
	require.NoError(t, err, "Failed to create output directory")

	// The 100 runs XOR experiment
	context.NumRuns = 40 //100 reduce to shorten test time
	experiment := experiments.Experiment{
		Id:     0,
		Trials: make(experiments.Trials, context.NumRuns),
	}
	err = experiment.Execute(context, startGenome, NewXORGenerationEvaluator(outDirPath))
	require.NoError(t, err, "Failed to perform XOR disconnected experiment")

	// Find winner statistics
	avgNodes, avgGenes, avgEvals, _ := experiment.AvgWinner()

	// check results
	if avgNodes < 5 {
		t.Error("avg_nodes < 5", avgNodes)
	} else if avgNodes > 15 {
		t.Error("avg_nodes > 15", avgNodes)
	}

	if avgGenes < 7 {
		t.Error("avg_genes < 7", avgGenes)
	} else if avgGenes > 20 {
		t.Error("avg_genes > 20", avgGenes)
	}

	maxEvals := float64(context.PopSize * context.NumGenerations)
	assert.True(t, avgEvals < maxEvals)

	t.Logf("avg_nodes: %.1f, avg_genes: %.1f, avg_evals: %.1f\n", avgNodes, avgGenes, avgEvals)
	meanComplexity, meanDiversity, meanAge := 0.0, 0.0, 0.0
	for _, t := range experiment.Trials {
		meanComplexity += t.BestComplexity().Mean()
		meanDiversity += t.Diversity().Mean()
		meanAge += t.BestAge().Mean()
	}
	count := float64(len(experiment.Trials))
	meanComplexity /= count
	meanDiversity /= count
	meanAge /= count
	t.Logf("Mean best organisms: complexity=%.1f, diversity=%.1f, age=%.1f", meanComplexity, meanDiversity, meanAge)
}
