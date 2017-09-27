package genetics

import (
	"testing"
	"github.com/yaricom/goNEAT/neat"
	"math/rand"
	"strings"
)

func TestNewPopulationRandom(t *testing.T) {
	rand.Seed(42)
	size, in, out, nmax := 10, 3, 2, 5
	recurrent := false
	link_prob := 0.5
	conf := neat.Neat{
		CompatThreshold:0.5,
	}
	pop, err := NewPopulationRandom(size, in, out, nmax, recurrent, link_prob, &conf)
	if err != nil {
		t.Error(err)
	}
	if pop == nil {
		t.Error("pop == nil")
	}
	if len(pop.Organisms) != size {
		t.Error("len(pop.Organisms) != size")
	}
	if pop.currNodeId != 11 {
		t.Error("pop.currNodeId != 11")
	}
	if pop.currInnovNum != int64(101) {
		t.Error("pop.currInnovNum != 101")
	}
	if len(pop.Species) == 0 {
		t.Error("len(pop.Species) == 0")
	}

}

func TestNewPopulation(t *testing.T) {
	rand.Seed(42)
	size, in, out, nmax, n := 10, 3, 2, 5, 3
	recurrent := false
	link_prob := 0.5
	conf := neat.Neat{
		CompatThreshold:0.5,
	}
	gen := NewGenomeRand(1, in, out, n, nmax, recurrent, link_prob)

	pop, err := NewPopulation(gen, size, &conf)
	if err != nil {
		t.Error(err)
	}
	if pop == nil {
		t.Error("pop == nil")
	}
	if len(pop.Organisms) != size {
		t.Error("len(pop.Organisms) != size")
	}
	last_node_id, _ := gen.getLastNodeId()
	if pop.currNodeId != last_node_id {
		t.Error("pop.currNodeId != last_node_id")
	}
	last_gene_innov_num, _ := gen.getLastGeneInnovNum()
	if pop.currInnovNum != last_gene_innov_num {
		t.Error("pop.currInnovNum != last_gene_innov_num")
	}
}

func TestReadPopulation(t *testing.T) {
	rand.Seed(42)
	pop_str := "genomestart 1\n" +
		"trait 1 0.1 0 0 0 0 0 0 0\n" +
		"trait 2 0.2 0 0 0 0 0 0 0\n" +
		"trait 3 0.3 0 0 0 0 0 0 0\n" +
		"node 1 0 1 1\n" +
		"node 2 0 1 1\n" +
		"node 3 0 1 3\n" +
		"node 4 0 0 2\n" +
		"gene 1 1 4 1.5 false 1 0 true\n" +
		"gene 2 2 4 2.5 false 2 0 true\n" +
		"gene 3 3 4 3.5 false 3 0 true\n" +
		"genomeend 1\n" +
		"genomestart 2\n" +
		"trait 1 0.1 0 0 0 0 0 0 0\n" +
		"trait 2 0.2 0 0 0 0 0 0 0\n" +
		"trait 3 0.3 0 0 0 0 0 0 0\n" +
		"node 1 0 1 1\n" +
		"node 2 0 1 1\n" +
		"node 3 0 1 3\n" +
		"node 4 0 0 2\n" +
		"gene 1 1 4 1.5 false 1 0 true\n" +
		"gene 2 2 4 2.5 false 2 0 true\n" +
		"gene 3 3 4 3.5 false 3 0 true\n" +
		"genomeend 2\n"
	conf := neat.Neat{
		CompatThreshold:0.5,
	}
	pop, err := ReadPopulation(strings.NewReader(pop_str), &conf)
	if err != nil {
		t.Error(err)
	}
	if pop == nil {
		t.Error("pop == nil")
	}
	if len(pop.Organisms) != 2 {
		t.Error("len(pop.Organisms) != size")
	}
	if len(pop.Species) != 1 {
		// because genomes are identical
		t.Error("len(pop.Species) != 1", len(pop.Species))
	}
}

func TestPopulation_verify(t *testing.T) {
	// first create population
	rand.Seed(42)
	pop_str := "genomestart 1\n" +
		"trait 1 0.1 0 0 0 0 0 0 0\n" +
		"trait 2 0.2 0 0 0 0 0 0 0\n" +
		"trait 3 0.3 0 0 0 0 0 0 0\n" +
		"node 1 0 1 1\n" +
		"node 2 0 1 1\n" +
		"node 3 0 1 3\n" +
		"node 4 0 0 2\n" +
		"gene 1 1 4 1.5 false 1 0 true\n" +
		"gene 2 2 4 2.5 false 2 0 true\n" +
		"gene 3 3 4 3.5 false 3 0 true\n" +
		"genomeend 1\n" +
		"genomestart 2\n" +
		"trait 1 0.1 0 0 0 0 0 0 0\n" +
		"trait 2 0.2 0 0 0 0 0 0 0\n" +
		"trait 3 0.3 0 0 0 0 0 0 0\n" +
		"node 1 0 1 1\n" +
		"node 2 0 1 1\n" +
		"node 3 0 1 3\n" +
		"node 4 0 0 2\n" +
		"gene 1 1 4 1.5 false 1 0 true\n" +
		"gene 2 2 4 2.5 false 2 0 true\n" +
		"gene 3 3 4 3.5 false 3 0 true\n" +
		"genomeend 2\n"
	conf := neat.Neat{
		CompatThreshold:0.5,
	}
	pop, err := ReadPopulation(strings.NewReader(pop_str), &conf)
	if err != nil {
		t.Error(err)
	}
	if pop == nil {
		t.Error("pop == nil")
	}

	// then verify created
	res, err := pop.verify()
	if err != nil {
		t.Error(err)
	}
	if !res {
		t.Error("Population verification failed, but must not")
	}
}