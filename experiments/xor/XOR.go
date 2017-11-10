// The XOR experiment serves to actually check that network topology actually evolves and everything works as expected.
// Because XOR is not linearly separable, a neural network requires hidden units to solve it. The two inputs must be
// combined at some hidden unit, as opposed to only at the out- put node, because there is no function over a linear
// combination of the inputs that can separate the inputs into the proper classes. These structural requirements make
// XOR suitable for testing NEAT’s ability to evolve structure.
package xor

import (
	"github.com/yaricom/goNEAT/neat"
	"os"
	"fmt"
	"github.com/yaricom/goNEAT/neat/genetics"
	"math"
	"github.com/yaricom/goNEAT/experiments"
	"time"
)

// The precision to use for XOR evaluation, i.e. one is x > 1 - precision and zero is x < precision
const precision = 0.5

// XOR is very simple and does not make a very interesting scientific experiment; however, it is a good way to
// check whether your system works.
// Make sure recurrency is disabled for the XOR test. If NEAT is able to add recurrent connections, it may solve XOR by
// memorizing the order of the training set. (Which is why you may even want to randomize order to be most safe) All
// documented experiments with XOR are without recurrent connections. Interestingly, XOR can be solved by a recurrent
// network with no hidden nodes.
//
// This method performs evolution on XOR for specified number of generations and output results into outDirPath
// It also returns number of nodes, genes, and evaluations performed per each run (context.NumRuns)
func XOR(context *neat.NeatContext, start_genome *genetics.Genome, out_dir_path string, experiment *experiments.Experiment) (err error) {

	if experiment.Trials == nil {
		experiment.Trials = make(experiments.Trials, context.NumRuns)
	}

	var pop *genetics.Population
	for run := 0; run < context.NumRuns; run++ {
		neat.InfoLog("\n>>>>> Spawning new population ")
		pop, err = genetics.NewPopulation(start_genome, context)
		if err != nil {
			neat.InfoLog("Failed to spawn new population from start genome")
			return err
		} else {
			neat.InfoLog("OK <<<<<")
		}
		neat.InfoLog(">>>>> Verifying spawned population ")
		_, err = pop.Verify()
		if err != nil {
			neat.ErrorLog("\n!!!!! Population verification failed !!!!!")
			return err
		} else {
			neat.InfoLog("OK <<<<<")
		}

		// start new trial
		trial := experiments.Trial{
			Id:run,
		}

		for gen := 0; gen < context.NumGenerations; gen++ {
			neat.InfoLog(fmt.Sprintf(">>>>> Epoch: %d\tRun: %d\n", gen, run))
			epoch := experiments.Epoch {
				Id:gen,
			}
			err = xor_epoch(pop, gen, out_dir_path, &epoch, context)
			if err != nil {
				neat.InfoLog(fmt.Sprintf("!!!!! Epoch %d evaluation failed !!!!!\n", gen))
				return err
			}
			epoch.Executed = time.Now()
			trial.Epochs = append(trial.Epochs, epoch)
			if epoch.Solved {
				neat.InfoLog(fmt.Sprintf(">>>>> The winner organism found in epoch %d! <<<<<\n", gen))
				break
			}
		}
		// store trial into experiment
		experiment.Trials[run] = trial
	}

	return nil
}

// This method evaluates one epoch for given population and prints results into specified directory if any.
func xor_epoch(pop *genetics.Population, generation int, out_dir_path string, epoch *experiments.Epoch, context *neat.NeatContext) (err error) {
	// Evaluate each organism on a test
	for _, org := range pop.Organisms {
		res, err := xor_evaluate(org, context)
		if err != nil {
			return  err
		}

		if res {
			epoch.Solved = true
			epoch.WinnerNodes = len(org.Genotype.Nodes)
			epoch.WinnerGenes = org.Genotype.Extrons()
			epoch.WinnerEvals = context.PopSize * epoch.Id + org.Genotype.Id
			epoch.Best = org
			if (epoch.WinnerNodes == 5) {
				// You could dump out optimal genomes here if desired
				opt_path := fmt.Sprintf("%s/%s", out_dir_path, "xor_optimal")
				file, err := os.Create(opt_path)
				if err != nil {
					neat.ErrorLog(fmt.Sprintf("Failed to dump optimal genome, reason: %s\n", err))
				} else {
					org.Genotype.Write(file)
					neat.InfoLog(fmt.Sprintf("Dumped optimal genome to: %s\n", opt_path))
				}
			}
			break // we have winner
		}
	}

	// Fill statistics about current epoch
	max_fitness := 0.0
	epoch.Diversity = len(pop.Species)
	epoch.Age = make(experiments.Floats, epoch.Diversity)
	epoch.Compexity = make(experiments.Floats, epoch.Diversity)
	epoch.Fitness = make(experiments.Floats, epoch.Diversity)
	for i, curr_species := range pop.Species {
		epoch.Age[i] = float64(curr_species.Age)
		epoch.Compexity[i] = float64(curr_species.Organisms[0].Phenotype.Complexity())
		epoch.Fitness[i] = curr_species.Organisms[0].Fitness

		// find best organism in epoch if not solved
		if !epoch.Solved && curr_species.Organisms[0].Fitness > max_fitness {
			max_fitness = curr_species.Organisms[0].Fitness
			epoch.Best = curr_species.Organisms[0]
		}
	}

	// Only print to file every print_every generations
	if epoch.Solved || generation % context.PrintEvery == 0 {
		pop_path := fmt.Sprintf("%s/gen_%d", out_dir_path, generation)
		file, err := os.Create(pop_path)
		if err != nil {
			neat.ErrorLog(fmt.Sprintf("Failed to dump population, reason: %s\n", err))
		} else {
			pop.WriteBySpecies(file)
		}
	}

	if epoch.Solved {
		// print winner organism
		for _, org := range pop.Organisms {
			if org.IsWinner {
				// Prints the winner organism to file!
				org_path := fmt.Sprintf("%s/%s", out_dir_path, "xor_winner")
				file, err := os.Create(org_path)
				if err != nil {
					neat.ErrorLog(fmt.Sprintf("Failed to dump winner organism genome, reason: %s\n", err))
				} else {
					org.Genotype.Write(file)
					neat.InfoLog(fmt.Sprintf("Generation #%d winner dumped to: %s\n", generation, org_path))
				}
				break
			}
		}
	} else {
		// Move to the next epoch if failed to find winner
		neat.DebugLog(">>>>> start next generation")
		_, err = pop.Epoch(generation + 1, context)
	}

	return err
}

// This methods evalueates provided organism
func xor_evaluate(organism *genetics.Organism, context *neat.NeatContext) (bool, error) {
	// The four possible input combinations to xor
	// The first number is for biasing
	in := [][]float64{
		{1.0, 0.0, 0.0},
		{1.0, 0.0, 1.0},
		{1.0, 1.0, 0.0},
		{1.0, 1.0, 1.0}}

	net_depth, err := organism.Phenotype.MaxDepth() // The max depth of the network to be activated
	if err != nil {
		neat.ErrorLog(fmt.Sprintf("Failed to estimate maximal depth of the network with genome:\n%s", organism.Genotype))
		return false, err
	}
	neat.DebugLog(fmt.Sprintf("Network depth: %d for organism: %d\n", net_depth, organism.Genotype.Id))
	if net_depth == 0 {
		neat.DebugLog(fmt.Sprintf("ALERT: Network depth is ZERO for Genome: %s", organism.Genotype))
	}

	success := false  // Check for successful activation
	out := make([]float64, 4) // The four outputs

	// Load and activate the network on each input
	for count := 0; count < 4; count++ {
		organism.Phenotype.LoadSensors(in[count])

		// Relax net and get output
		success, err = organism.Phenotype.Activate()
		if err != nil {
			neat.ErrorLog("Failed to activate network")
			return false, err
		}

		// use depth to ensure relaxation
		for relax := 0; relax <= net_depth; relax++ {
			success, err = organism.Phenotype.Activate()
			if err != nil {
				neat.ErrorLog("Failed to activate network")
				return false, err
			}
		}
		out[count] = organism.Phenotype.Outputs[0].Activation

		organism.Phenotype.Flush()
	}

	error_sum := 0.0
	if (success) {
		// Mean Squared Error
		error_sum = math.Abs(out[0]) + math.Abs(1.0 - out[1]) + math.Abs(1.0 - out[2]) + math.Abs(out[3])
		organism.Fitness = math.Pow(4.0 - error_sum, 2.0)
		organism.Error = error_sum
	} else {
		// The network is flawed (shouldn't happen)
		error_sum = 999.0
		organism.Fitness = 0.001
	}

	if out[0] < precision && out[1] >= 1 - precision && out[2] >= 1 - precision && out[3] < precision {
		organism.IsWinner = true
		neat.InfoLog(fmt.Sprintf(">>>> Output activations: %e\n", out))

	} else {
		organism.IsWinner = false
	}
	return organism.IsWinner, nil
}