package wann

import (
	"fmt"
	"math/rand"
	"testing"
)

// Use a pseudo-random seed
//var commonSeed = time.Now().UTC().UnixNano()

// Use a specific seed
var commonSeed int64 = 1571917826405889420

func TestNewNetwork(t *testing.T) {
	rand.Seed(commonSeed)
	net := NewNetwork(&Config{
		Inputs:          5,
		ConnectionRatio: 0.5,
		SharedWeight:    0.5,
	})
	fmt.Println(net)
	for i, n := range net.AllNodes {
		if NeuronIndex(i) != n.neuronIndex {
			t.Fail()
		}
	}
}

func TestGet(t *testing.T) {
	rand.Seed(commonSeed)
	net := NewNetwork()
	fmt.Println(net)
	fmt.Println(net.Get(0))
	if net.OutputNode != 0 {
		t.Fail()
	}
}

func TestIsInput(t *testing.T) {
	rand.Seed(commonSeed)
	net := NewNetwork(&Config{
		Inputs:          5,
		ConnectionRatio: 0.5,
		SharedWeight:    0.5,
	})
	if !net.IsInput(1) {
		t.Fail()
	}
	if net.IsInput(0) {
		t.Fail()
	}
}

func TestForEachConnected(t *testing.T) {
	rand.Seed(commonSeed)
	net := NewNetwork(&Config{
		Inputs:          5,
		ConnectionRatio: 0.5,
		SharedWeight:    0.5,
	})
	net.ForEachConnected(func(n *Neuron) {
		fmt.Printf("%d: %s, distance from output node: %d\n", n.neuronIndex, n, n.distanceFromOutputNode)
	})
}

func TestAll(t *testing.T) {
	rand.Seed(commonSeed)
	net := NewNetwork(&Config{
		Inputs:          5,
		ConnectionRatio: 0.7,
		SharedWeight:    0.5,
	})
	for _, node := range net.All() {
		fmt.Println(node)
	}
}

func TestEvaluate2(t *testing.T) {
	rand.Seed(commonSeed)
	net := NewNetwork(&Config{
		Inputs:          5,
		ConnectionRatio: 0.7,
		SharedWeight:    0.5,
	})
	_ = net.Evaluate([]float64{0.1, 0.2, 0.3, 0.4, 0.5})
}

func TestInsertNode(t *testing.T) {
	rand.Seed(commonSeed)
	net := NewNetwork(&Config{
		Inputs:          5,
		ConnectionRatio: 0.5,
		SharedWeight:    0.5,
	})
	_, newNeuronIndex := net.NewNeuron()
	if err := net.InsertNode(0, 2, newNeuronIndex); err != nil {
		t.Error(err)
	}
	//fmt.Println(net)
	_ = net.Evaluate([]float64{0.1, 0.2, 0.3, 0.4, 0.5})
}

func TestAddConnection(t *testing.T) {
	rand.Seed(commonSeed)
	net := NewNetwork(&Config{
		Inputs:          5,
		ConnectionRatio: 0.5,
	})
	_, newNeuronIndex := net.NewNeuron()
	if err := net.InsertNode(net.OutputNode, 2, newNeuronIndex); err != nil {
		t.Error(err)
	}
	// Add a connection from 1 to the new neuron.
	// This is the same as making the new neuron have an additional input neuron: index 1
	if err := net.AddConnection(1, newNeuronIndex); err != nil {
		t.Error(err)
	}
	// Add a connection from the output node to the output node. Should fail.
	if err := net.AddConnection(net.OutputNode, net.OutputNode); err == nil {
		t.Fail()
	}
	// Adding a made-up index should fail as well
	if err := net.AddConnection(net.OutputNode, 999); err == nil {
		t.Fail()
	}
}

func TestRandomizeActivationFunctionForRandomNeuron(t *testing.T) {
	rand.Seed(commonSeed)
	net := NewNetwork(&Config{
		Inputs:          5,
		ConnectionRatio: 0.5,
	})
	net.RandomizeActivationFunctionForRandomNeuron()
}

func TestNetworkString(t *testing.T) {
	rand.Seed(commonSeed)
	net := NewNetwork(&Config{
		Inputs:          5,
		ConnectionRatio: 0.5,
	})
	//fmt.Println(net.String())
	_ = net.String()
}

func TestSetWeight(t *testing.T) {
	net := NewNetwork()
	net.SetWeight(0.1234)
	if net.Weight != 0.1234 {
		t.Fail()
	}
}

func TestComplexity(t *testing.T) {
	rand.Seed(commonSeed)
	net := NewNetwork(&Config{
		Inputs:          5,
		ConnectionRatio: 0.0,
	})
	// The complexity will vary, because the performance varies when
	// estimating the complexity of each function.
	// But the complexity compared between networks should still hold.
	firstComplexity := net.Complexity()
	//fmt.Println("First network complexity:", firstComplexity)
	// Adding a connection increases the complexity
	net.AddConnection(0, 1)
	secondComplexity := net.Complexity()
	//fmt.Println("Second network complexity:", secondComplexity)
	if firstComplexity >= secondComplexity {
		t.Fail()
	}
}

func TestLeftRight(t *testing.T) {
	rand.Seed(commonSeed)
	net := NewNetwork(&Config{
		Inputs:          3,
		ConnectionRatio: 1.0,
	})
	net.AllNodes[1].ActivationFunctionIndex = Swish
	a, b, _ := net.LeftRight(0, 1)
	// output node to the right
	if a != 1 || b != 0 {
		t.Fail()
	}
	// output node to the right
	a, b, _ = net.LeftRight(1, 0)
	if a != 1 || b != 0 {
		t.Fail()
	}
	net.WriteSVG("before.svg")
	_, nodeIndex := net.NewNeuron()
	err := net.InsertNode(0, 1, nodeIndex)
	if err != nil {
		t.Error(err)
	}
	net.WriteSVG("after.svg")
	a, b, _ = net.LeftRight(0, nodeIndex)
	// output node to the right
	if a != nodeIndex || b != 0 {
		t.Fail()
	}
	a, b, _ = net.LeftRight(nodeIndex, 0)
	// output node to the right
	if a != nodeIndex || b != 0 {
		t.Fail()
	}
	a, b, _ = net.LeftRight(1, nodeIndex)
	// Here, the new node should be to the right, since it's between node 1 and the output node
	if a != 1 || b != nodeIndex {
		t.Fail()
	}
	//net.WriteSVG("c.svg")
	fmt.Println(net)
	a, b, _ = net.LeftRight(nodeIndex, 1)
	fmt.Println("nodeIndex:", nodeIndex)
	fmt.Println("1:", 1)
	fmt.Println("a:", a)
	fmt.Println("b:", b)

	if a != 1 || b != nodeIndex {
		t.Fail()
	}
}

func TestDepth(t *testing.T) {
	rand.Seed(commonSeed)
	net := NewNetwork(&Config{
		Inputs:          3,
		ConnectionRatio: 1.0,
	})
	fmt.Println(net.Depth())
	_, nodeIndex := net.NewBlankNeuron()
	_ = net.InsertNode(0, 1, nodeIndex)
	fmt.Println(net.Depth())
}

func ExampleCombine() {
	ac := []NeuronIndex{0, 1, 2, 3, 4}
	bc := []NeuronIndex{5, 6, 7, 8, 9}
	fmt.Println(Combine(ac, bc))
	// Output:
	// [0 1 2 3 4 5 6 7 8 9]
}

func TestGetRandomNeuron(t *testing.T) {
	rand.Seed(commonSeed)
	net := NewNetwork(&Config{
		Inputs:          5,
		ConnectionRatio: 1.0,
	})
	stats := make(map[NeuronIndex]uint)
	for i := 0; i < 1000; i++ {
		ni := net.GetRandomNode()
		if _, ok := stats[ni]; !ok {
			stats[ni] = 0
		} else {
			stats[ni]++
		}
	}
	fmt.Println(stats)
	// Check that the output node exists in the stats
	if _, ok := stats[0]; !ok {
		t.Fail()
	}

	// This is more a test of the random number generator than anything. Disable:
	// // This isn't 00% watertight, but each element should have been chosen around 160 times, +- 30
	// center := uint(160)
	// margin := uint(30)
	// for _, chosenCount := range stats {
	// 	if chosenCount < (center-margin) || chosenCount > (center+margin) {
	// 		t.Fail()
	// 	}
	// }
}

func TestGetRandomInputNode(t *testing.T) {
	rand.Seed(commonSeed)
	net := NewNetwork(&Config{
		Inputs:          5,
		ConnectionRatio: 1.0,
	})
	stats := make(map[NeuronIndex]uint)
	for i := 0; i < 1000; i++ {
		ni := net.GetRandomInputNode()
		if _, ok := stats[ni]; !ok {
			stats[ni] = 0
		} else {
			stats[ni]++
		}
	}
	fmt.Println(stats)
	// Check that the output node does not exist in the stats
	if _, ok := stats[0]; ok {
		t.Fail()
	}
}

func TestConnected(t *testing.T) {
	rand.Seed(commonSeed)
	net := NewNetwork(&Config{
		Inputs:          5,
		ConnectionRatio: 0.1,
	})
	connected := net.Connected()
	if connected[0] != 0 || connected[1] != 2 {
		t.Fail()
	}
}

func TestUnconnected(t *testing.T) {
	rand.Seed(commonSeed)
	net := NewNetwork(&Config{
		Inputs:          5,
		ConnectionRatio: 0.5,
	})
	unconnected := net.Unconnected()
	correct := []NeuronIndex{1, 3, 4}
	for i := 0; i < len(unconnected); i++ {
		if unconnected[i] != correct[i] {
			t.Fail()
		}
	}
}

func TestCopy(t *testing.T) {
	rand.Seed(commonSeed)
	net := NewNetwork(&Config{
		Inputs:          5,
		ConnectionRatio: 0.5,
	})
	net2 := net.Copy()
	n := NewUnconnectedNeuron()
	net2.AllNodes[1] = *n
	if net.String() != net2.String() {
		t.Fail()
	}
	net3 := net
	if net.String() != net3.String() {
		t.Fail()
	}

}

func TestForEachConnectedNodeIndex(t *testing.T) {
	//ForEachConnectedNodeIndex(f func(ni NeuronIndex, distanceFromOutputNode int))
}
