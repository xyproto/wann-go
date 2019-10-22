package wann

import (
	"fmt"
	"math/rand"
	"sort"
	"testing"
	"time"
)

var currentTime = time.Now().UTC().UnixNano()

func TestNetwork(t *testing.T) {
	rand.Seed(currentTime)
	net := NewNetwork(&Config{
		Inputs:          5,
		ConnectionRatio: 0.5,
		SharedWeight:    0.5,
	})
	fmt.Println(net)
}

// func NewNetwork(c *Config) *Network {
// func (net *Network) InsertNode(a, b, newNode *Neuron) error {
// func (net *Network) AddConnection(a, b *Neuron) error {
// func (net *Network) ChangeActivationFunction(n *Neuron, f func(float64) float64) {
// func (net *Network) String() string {

func TestUpDetection(t *testing.T) {
	fmt.Println("### Up detection ###")

	//  o
	// ooo
	up := []float64{0.0, 1.0, 0.0, 1.0, 1.0, 1.0}

	// ooo
	//  o
	down := []float64{1.0, 1.0, 1.0, 0.0, 1.0, 0.0}

	// ooo
	//   o
	left := []float64{1.0, 1.0, 1.0, 0.0, 0.0, 1.0}

	// ooo
	// o
	right := []float64{1.0, 1.0, 1.0, 0.1, 0.0, 0.0}

	_ = up
	_ = down
	_ = left
	_ = right

	// ---

	rand.Seed(currentTime)

	inputNumbers := up

	// 1. Create a neural network that is supposed to be able to detect "up"
	// 2. Use the inputs from up, down, left, right.
	// 3. The goal is that the output neuron should fire "1" for the up patterns and "0" for the rest. (or at least a high/low value)
	// 4. Train, according to the paper.

	config := &Config{
		Inputs:          len(inputNumbers),
		ConnectionRatio: 0.5,
		SharedWeight:    1.0,
	}

	// population of 10 networks
	N := 10
	population := make([]*Network, N)

	// Initialize the population
	for i := 0; i < N; i++ {
		population[i] = NewNetwork(config)
	}

	// For each network, for each weight, evaluate the result
	scoreMap := make(map[float64]int)
	for i := 0; i < N; i++ {
		net := population[i]
		bestResult := 0.0
		bestWeight := 0.0
		complexity := net.Complexity()
		for w := 0.0; w <= 1.0; w += 0.1 {
			net.SetWeight(w)
			result := net.Evaluate(inputNumbers)
			if result > bestResult {
				bestResult = result
				bestWeight = w
			}
		}
		score := bestResult / (complexity * 0.1)
		scoreMap[score] = i
		fmt.Println("Best weight for network", i, "is", bestWeight, "with score", score, "(best result", bestResult, ", complexity", complexity, ")")
	}

	// Prepare to sort the score -> population index map, by key
	keys := make(sort.Float64Slice, 0, len(scoreMap))
	for k := range scoreMap {
		keys = append(keys, k)
	}

	// --- RANK BY SCORE ---
	keys.Sort()

	for _, scoreIndex := range keys {
		fmt.Println("score", scoreIndex, "network index", scoreMap[scoreIndex])
	}
	if len(keys) > 0 {
		bestIndex := scoreMap[keys[len(keys)-1]]
		fmt.Println("Best network index: ", bestIndex)
		bestNetwork := population[bestIndex]
		fmt.Println(bestNetwork)
	}

	// Now take the best networks and make mutated offspring.
	// Delete the worst networks.

}
