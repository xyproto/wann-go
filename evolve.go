package wann

import (
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

// ScorePopulation evaluates a population, given a slice of input numbers.
// It returns a map with scores, together with the sum of scores.
func ScorePopulation(population []*Network, weight float64, inputData [][]float64, correctOutputMultipliers []float64) (map[int]float64, float64) {

	scoreMap := make(map[int]float64)
	scoreSum := 0.0

	for i := 0; i < len(population); i++ {
		net := population[i]

		net.SetWeight(weight)
		result := 0.0
		for i := 0; i < len(inputData); i++ {
			result += net.Evaluate(inputData[i]) * correctOutputMultipliers[i]
		}
		score := result / net.Complexity()
		scoreSum += score
		scoreMap[i] = score
	}
	return scoreMap, scoreSum
}

// Modify the network using one of the three methods outlined in the paper:
// * Insert node
// * Add connection
// * Change activation function
func (net *Network) Modify(maxIterations int) {

	// Use method 0, 1 or 2
	// Method 1 and 2 are always fine, method 0 has had issues
	method := rand.Intn(3) // up to and not including 3

	// Perform a modfification, using one of the three methods outlined in the paper
	switch method {
	case 0:
		// Insert a node, replacing a randomly chosen existing connection
		counter := 0
		for net.InsertRandomNode() == false {
			counter++
			if counter > maxIterations {
				break
			}
		}
	case 1:
		//net.checkInputNeurons()
		nodeA, nodeB := net.GetRandomNode(), net.GetRandomNode()
		// A bit risky, time-wise, but continue finding random neurons until they work out
		// Create a new connection
		counter := 0
		for net.AddConnection(nodeA, nodeB) != nil {
			nodeA, nodeB = net.GetRandomNode(), net.GetRandomNode()
			counter++
			if maxIterations > 0 && counter > maxIterations {
				// Could not add a connection. The possibilities for connections might be saturated.
				return
			}
		}
		//net.checkInputNeurons()
	case 2:
		// Change the activation function to a randomly selected one
		net.RandomizeActivationFunctionForRandomNeuron()
	default:
		panic("implementation error: invalid method number: " + strconv.Itoa(method))
	}
}

// initialize the pseaudo-random number generator, either using the config.RandomSeed or the time
func (config *Config) initRandom() {
	randomSeed := config.RandomSeed
	if config.RandomSeed == 0 {
		randomSeed = time.Now().UTC().UnixNano()
	}
	if config.Verbose {
		fmt.Println("Using random seed:", randomSeed)
	}
	// Initialize the pseudo-random number generator
	rand.Seed(randomSeed)
}

// Init will initialize the pseudo-random number generator and estimate the complexity of the available activation functions
func (config *Config) Init() {
	config.initRandom()
	config.estimateComplexity()
	config.initialized = true
}

// Evolve evolves a neural network, given a slice of training data and a slice of correct output values.
// Will overwrite config.Inputs.
func (config *Config) Evolve(inputData [][]float64, correctOutputMultipliers []float64) (*Network, error) {

	const maxModificationInterationsWhenMutating = 10

	// Initialize, if needed
	if !config.initialized {
		config.Init()
	}

	inputLength := len(inputData)
	if inputLength == 0 {
		return nil, errors.New("no input data")
	}

	if len(correctOutputMultipliers) == 1 && inputLength != 1 {
		// Assume the first slice of floats in the input data is the correct and that the rest are examples of being wrong
		for i := 1; i < inputLength; i++ {
			correctOutputMultipliers = append(correctOutputMultipliers, -1.0)
		}
	} else if inputLength != len(correctOutputMultipliers) {
		// Assume that the list of correct output multipliers should match the length of the float64 slices in inputData
		return nil, errors.New("the length of the input data and the slice of output multipliers differs")
	}

	config.inputs = len(inputData[0])

	population := make([]*Network, config.PopulationSize)

	// Initialize the population
	for i := 0; i < config.PopulationSize; i++ {
		n := NewNetwork(config)
		population[i] = &n
	}

	var bestNetwork *Network
	var bestWeight float64

	// Keep track of the best scores
	bestScore := 0.0
	lastBestScore := 0.0
	noImprovementOfBestScoreCounter := 0

	// Keep track of the average scores
	averageScore := 0.0
	//lastAverageScore := 0.0
	//noImprovementOfAverageScoreCounter := 0

	// Keep track of the worst scores
	worstScore := 0.0
	//lastWorstScore := 0.0
	//noImprovementOfWorstScoreCounter := 0

	if config.Verbose {
		fmt.Printf("Starting evolution with population size %d, for %d generations.\n", config.PopulationSize, config.Generations)
	}

	// For each generation, evaluate and modify the networks
	for j := 0; j < config.Generations; j++ {

		bestNetwork = nil

		// Initialize the scores with unlikely values
		// TODO: Use the first network in the population for initializing these instead
		bestScore = -9999.0
		averageScore = 0.0
		worstScore = 9999.0

		// Random weight from -2.0 to 2.0
		w := rand.Float64()

		// The scores for this generation (using a random shared weight within ScorePopulation).
		// CorrectOutputMultipliers gives weight to the "correct" or "wrong" results, with the same index as the inputData
		// Score each network in the population.
		scoreMap, scoreSum := ScorePopulation(population, w, inputData, correctOutputMultipliers)

		// Sort by score
		scoreList := SortByValue(scoreMap)

		// Handle the best score stats
		lastBestScore = bestScore
		if scoreList[0].Value > bestScore {
			bestScore = scoreList[0].Value
		}
		if bestScore >= lastBestScore {
			bestNetwork = population[scoreList[0].Key]
			bestWeight = w
			noImprovementOfBestScoreCounter = 0
		} else {
			noImprovementOfBestScoreCounter++
		}

		// Handle the average score stats
		//lastAverageScore = averageScore
		averageScore = scoreSum / float64(config.PopulationSize)
		// if averageScore >= lastAverageScore {
		// 	noImprovementOfAverageScoreCounter = 0
		// } else {
		// 	noImprovementOfAverageScoreCounter++
		// }

		// Handle the worst score stats
		//lastWorstScore = worstScore
		if scoreList[len(scoreList)-1].Value < worstScore {
			worstScore = scoreList[len(scoreList)-1].Value
		}
		// if worstScore >= lastWorstScore {
		// 	noImprovementOfWorstScoreCounter = 0
		// } else {
		// 	noImprovementOfWorstScoreCounter++
		// }

		if bestNetwork == nil {
			panic("implementation error: no best network")
		}

		if config.Verbose {
			fmt.Printf("[generation %d] worst score = %f, average score = %f, best score = %f\n", j, worstScore, averageScore, bestScore)
			if noImprovementOfBestScoreCounter > 0 {
				fmt.Printf("No improvement in the best score for the last %d generations\n", noImprovementOfBestScoreCounter)
			}
		}

		bestThirdCountdown := len(population) / 3

		goodNetworks := make([]*Network, 0)

		// Now loop over all networks, sorted by score (descending order)
		for _, p := range scoreList {
			networkIndex := p.Key
			//networkScore := p.Value
			if bestThirdCountdown > 0 {
				bestThirdCountdown--
				// In the best third of the networks
				//fmt.Println("BEST THIRD:", networkIndex, "score", networkScore)
				goodNetworks = append(goodNetworks, population[networkIndex])
			} else {
				//fmt.Println("WORST TWO THIRDS:", networkIndex, "score", networkScore)
				randomGoodNetwork := goodNetworks[rand.Intn(len(goodNetworks))]
				//randomGoodNetwork.UpdateNetworkPointers()
				//randomGoodNetwork.checkInputNeurons()
				randomGoodNetworkCopy := randomGoodNetwork.Copy()
				//randomGoodNetworkCopy.UpdateNetworkPointers()
				//randomGoodNetworkCopy.checkInputNeurons()
				//randomGoodNetworkCopy.UpdateNetworkPointers()
				randomGoodNetworkCopy.Modify(maxModificationInterationsWhenMutating)
				//randomGoodNetworkCopy.UpdateNetworkPointers()
				//randomGoodNetworkCopy.checkInputNeurons()
				// Replace the "bad" network with the modified copy of a "good" one
				// It's important that this is a pointer to a Network and not
				// a bare Network, so that the node .Net pointers are correct.
				population[networkIndex] = randomGoodNetworkCopy
			}
		}
	}
	if config.Verbose {
		fmt.Printf("[all time best network, random weight ] weight=%f score=%f\n", bestWeight, bestScore)
	}

	// Now find the best weight for the best network, using a population of 1
	// and a step size of 0.0001 for the weight
	population = []*Network{bestNetwork}
	for w := 0.0; w <= 1.0; w += 0.0001 {
		scoreMap, _ := ScorePopulation(population, w, inputData, correctOutputMultipliers)

		// Sort by score
		scoreList := SortByValue(scoreMap)

		// Handle the best score stats
		if scoreList[0].Value > bestScore {
			bestScore = scoreList[0].Value
			bestWeight = w
		}
	}

	// Check if the best network is nil, just in case
	if bestNetwork == nil {
		return nil, errors.New("the total best network is nil")
	}

	// Save the best weight for the network
	bestNetwork.SetWeight(bestWeight)

	if config.Verbose {
		fmt.Printf("[all time best network, optimal weight] weight=%f score=%f\n", bestWeight, bestScore)
	}

	return bestNetwork, nil
}
