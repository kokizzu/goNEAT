package experiment

import (
	"bytes"
	"encoding/gob"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yaricom/goNEAT/v3/neat"
	"github.com/yaricom/goNEAT/v3/neat/genetics"
	"github.com/yaricom/goNEAT/v3/neat/math"
	"github.com/yaricom/goNEAT/v3/neat/network"
	"math/rand"
	"testing"
	"time"
)

func TestGeneration_Average(t *testing.T) {
	fitness := Floats{1.0, 4.0, 4.0}
	ages := Floats{10.0, 40.0, 40.0}
	complexities := Floats{2.0, 3.0, 4.0}
	g := createGenerationWith(fitness, ages, complexities)
	f, a, c := g.Average()
	assert.Equal(t, 3.0, f)
	assert.Equal(t, 30.0, a)
	assert.Equal(t, 3.0, c)
}

func TestGeneration_FillPopulationStatistics(t *testing.T) {
	rand.Seed(42)
	in, out, nmax := 3, 2, 5
	linkProb := 0.9
	conf := neat.Options{
		DisjointCoeff:   0.5,
		ExcessCoeff:     0.5,
		MutdiffCoeff:    0.5,
		CompatThreshold: 5.0,
		PopSize:         100,
	}
	pop, err := genetics.NewPopulationRandom(in, out, nmax, false, linkProb, &conf)
	require.NoError(t, err, "failed to create population")
	require.NotNil(t, pop, "population expected")
	require.Len(t, pop.Organisms, conf.PopSize, "wrong population size")
	assert.True(t, len(pop.Species) > 0, "population has no species")

	maxFitness := -1.0
	for i := range pop.Species {
		for j := range pop.Species[i].Organisms {
			fitness := rand.Float64() * 100
			pop.Species[i].Organisms[j].Fitness = fitness
			if fitness > maxFitness {
				maxFitness = fitness
			}
		}
	}

	gen := Generation{
		Id:      1,
		TrialId: 1,
	}
	gen.FillPopulationStatistics(pop)
	expectedSpecies := 5
	assert.Equal(t, expectedSpecies, gen.Diversity, "wrong diversity")
	assert.Equal(t, expectedSpecies, len(gen.Fitness))
	assert.Equal(t, expectedSpecies, len(gen.Age))
	assert.EqualValues(t, Floats{1, 1, 1, 1, 1}, gen.Age)
	assert.Equal(t, expectedSpecies, len(gen.Complexity))
	assert.Equal(t, Floats{10, 23, 36, 31, 38}, gen.Complexity)
	assert.NotNil(t, gen.Best)
	assert.Equal(t, maxFitness, gen.Best.Fitness)
}

func createGenerationWith(fitness Floats, ages Floats, complexities Floats) *Generation {
	return &Generation{
		Fitness:    fitness,
		Age:        ages,
		Complexity: complexities,
	}
}

// Tests encoding/decoding of generation
func TestGeneration_Encode_Decode(t *testing.T) {
	genomeId, fitness := 10, 23.0
	gen := buildTestGeneration(genomeId, fitness)

	var buff bytes.Buffer
	enc := gob.NewEncoder(&buff)

	// encode generation
	err := gen.Encode(enc)
	require.NoError(t, err, "failed to encode generation")

	// decode generation
	data := buff.Bytes()
	dec := gob.NewDecoder(bytes.NewBuffer(data))
	dgen := &Generation{}
	err = dgen.Decode(dec)
	require.NoError(t, err, "failed to decode generation")

	//  and test fields
	assert.EqualValues(t, gen, dgen)
}

func buildTestGeneration(genId int, fitness float64) *Generation {
	epoch := Generation{}
	epoch.Id = genId
	epoch.Executed = time.Now().Round(time.Second)
	epoch.Solved = true
	epoch.Fitness = Floats{10.0, 30.0, 40.0, fitness}
	epoch.Age = Floats{1.0, 3.0, 4.0, 10.0}
	epoch.Complexity = Floats{34.0, 21.0, 56.0, 15.0}
	epoch.Diversity = 32
	epoch.WinnerEvals = 12423
	epoch.WinnerNodes = 7
	epoch.WinnerGenes = 5

	genome := buildTestGenome(genId)
	org := genetics.Organism{Fitness: fitness, Genotype: genome, Generation: genId}
	epoch.Best = &org

	return &epoch
}

func buildTestGenome(id int) *genetics.Genome {
	traits := []*neat.Trait{
		{Id: 1, Params: []float64{0.1, 0, 0, 0, 0, 0, 0, 0}},
		{Id: 3, Params: []float64{0.3, 0, 0, 0, 0, 0, 0, 0}},
		{Id: 2, Params: []float64{0.2, 0, 0, 0, 0, 0, 0, 0}},
	}

	nodes := []*network.NNode{
		{Id: 1, NeuronType: network.InputNeuron, ActivationType: math.NullActivation, Incoming: make([]*network.Link, 0), Outgoing: make([]*network.Link, 0)},
		{Id: 2, NeuronType: network.InputNeuron, ActivationType: math.NullActivation, Incoming: make([]*network.Link, 0), Outgoing: make([]*network.Link, 0)},
		{Id: 3, NeuronType: network.BiasNeuron, ActivationType: math.SigmoidSteepenedActivation, Incoming: make([]*network.Link, 0), Outgoing: make([]*network.Link, 0)},
		{Id: 4, NeuronType: network.OutputNeuron, ActivationType: math.SigmoidSteepenedActivation, Incoming: make([]*network.Link, 0), Outgoing: make([]*network.Link, 0)},
	}

	genes := []*genetics.Gene{
		genetics.NewGeneWithTrait(traits[0], 1.5, nodes[0], nodes[3], false, 1, 0),
		genetics.NewGeneWithTrait(traits[2], 2.5, nodes[1], nodes[3], false, 2, 0),
		genetics.NewGeneWithTrait(traits[1], 3.5, nodes[2], nodes[3], false, 3, 0),
	}

	return genetics.NewGenome(id, traits, nodes, genes)
}
