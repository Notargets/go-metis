package main

import (
	"fmt"
	"log"

	"github.com/notargets/go-metis"
)

func main() {
	fmt.Printf("METIS Version: %s\n\n", metis.Version())

	// Example 1: Simple graph partitioning
	example1()

	// Example 2: Weighted graph partitioning
	example2()

	// Example 3: Mesh partitioning
	example3()

	// Example 4: Nested dissection ordering
	example4()
}

func example1() {
	fmt.Println("Example 1: Simple Graph Partitioning")
	fmt.Println("====================================")

	// Create a simple 6-node graph
	// Graph structure:
	//   0 --- 1 --- 2
	//   |     |     |
	//   3 --- 4 --- 5
	xadj := []int32{0, 2, 5, 7, 9, 12, 14}
	adjncy := []int32{
		1, 3, // neighbors of 0
		0, 2, 4, // neighbors of 1
		1, 5, // neighbors of 2
		0, 4, // neighbors of 3
		1, 3, 5, // neighbors of 4
		2, 4, // neighbors of 5
	}

	// Create options
	opts := make([]int32, metis.NoOptions)
	if err := metis.SetDefaultOptions(opts); err != nil {
		log.Fatal(err)
	}

	// Partition into 2 parts using recursive bisection
	fmt.Println("Recursive bisection (2 parts):")
	part, edgeCut, err := metis.PartGraphRecursive(xadj, adjncy, 2, opts)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Edge cut: %d\n", edgeCut)
	fmt.Printf("Partition assignment: %v\n", part)

	// Partition into 3 parts using k-way
	fmt.Println("\nK-way partitioning (3 parts):")
	part, edgeCut, err = metis.PartGraphKway(xadj, adjncy, 3, opts)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Edge cut: %d\n", edgeCut)
	fmt.Printf("Partition assignment: %v\n", part)
	fmt.Println()
}

func example2() {
	fmt.Println("Example 2: Weighted Graph Partitioning")
	fmt.Println("======================================")

	// Create a weighted graph
	xadj := []int32{0, 3, 6, 9, 12}
	adjncy := []int32{
		1, 2, 3, // neighbors of 0
		0, 2, 3, // neighbors of 1
		0, 1, 3, // neighbors of 2
		0, 1, 2, // neighbors of 3
	}

	// Vertex weights (importance of each vertex)
	vwgt := []int32{10, 20, 30, 40}

	// Edge weights (strength of connections)
	adjwgt := []int32{
		5, 3, 2, // weights for edges from 0
		5, 4, 1, // weights for edges from 1
		3, 4, 6, // weights for edges from 2
		2, 1, 6, // weights for edges from 3
	}

	opts := make([]int32, metis.NoOptions)
	metis.SetDefaultOptions(opts)

	// Partition with weights
	part, edgeCut, err := metis.PartGraphKwayWeighted(xadj, adjncy, vwgt, adjwgt, 2, nil, nil, opts)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Weighted edge cut: %d\n", edgeCut)
	fmt.Printf("Partition assignment: %v\n", part)

	// Calculate partition weights
	partWeights := make([]int32, 2)
	for i, p := range part {
		partWeights[p] += vwgt[i]
	}
	fmt.Printf("Partition weights: %v (total: %d)\n", partWeights, partWeights[0]+partWeights[1])

	// Try with target partition weights (30%/70% split)
	fmt.Println("\nWith target weights (30%/70%):")
	tpwgts := []float32{0.3, 0.7}
	part, edgeCut, err = metis.PartGraphKwayWeighted(xadj, adjncy, vwgt, nil, 2, tpwgts, nil, opts)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Edge cut: %d\n", edgeCut)
	fmt.Printf("Partition assignment: %v\n", part)

	partWeights = make([]int32, 2)
	for i, p := range part {
		partWeights[p] += vwgt[i]
	}
	fmt.Printf("Partition weights: %v\n", partWeights)
	fmt.Println()
}

func example3() {
	fmt.Println("Example 3: Mesh Partitioning")
	fmt.Println("============================")

	// Create a simple mesh with 3 triangular elements
	// Mesh structure:
	//   0 --- 1
	//   |\    |
	//   | \ 1 |
	//   |0 \  |
	//   |   \ |
	//   3 --- 2
	//       2

	ne := int32(3) // Number of elements
	nn := int32(4) // Number of nodes

	// Element connectivity (each element lists its nodes)
	eptr := []int32{0, 3, 6, 9} // Start index for each element
	eind := []int32{
		0, 3, 2, // Element 0: nodes 0, 3, 2
		0, 2, 1, // Element 1: nodes 0, 2, 1
		2, 3, 1, // Element 2: nodes 2, 3, 1
	}

	opts := make([]int32, metis.NoOptions)
	metis.SetDefaultOptions(opts)

	// Partition mesh using dual graph (element-based)
	fmt.Println("Dual graph partitioning (element-based):")
	objval, epart, npart, err := metis.PartMeshDual(ne, nn, eptr, eind, nil, nil, 2, 2, nil, opts)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Objective value: %d\n", objval)
	fmt.Printf("Element partition: %v\n", epart)
	fmt.Printf("Node partition: %v\n", npart)

	// Partition mesh using nodal graph (node-based)
	fmt.Println("\nNodal graph partitioning (node-based):")
	objval, epart, npart, err = metis.PartMeshNodal(ne, nn, eptr, eind, nil, nil, 2, nil, opts)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Objective value: %d\n", objval)
	fmt.Printf("Element partition: %v\n", epart)
	fmt.Printf("Node partition: %v\n", npart)
	fmt.Println()
}

func example4() {
	fmt.Println("Example 4: Nested Dissection Ordering")
	fmt.Println("=====================================")

	// Create a grid-like graph (good for testing nested dissection)
	// 3x3 grid:
	//   0 - 1 - 2
	//   |   |   |
	//   3 - 4 - 5
	//   |   |   |
	//   6 - 7 - 8

	xadj := []int32{0, 2, 5, 7, 10, 14, 17, 19, 22, 24}
	adjncy := []int32{
		1, 3, // 0
		0, 2, 4, // 1
		1, 5, // 2
		0, 4, 6, // 3
		1, 3, 5, 7, // 4
		2, 4, 8, // 5
		3, 7, // 6
		4, 6, 8, // 7
		5, 7, // 8
	}

	opts := make([]int32, metis.NoOptions)
	metis.SetDefaultOptions(opts)

	// Compute nested dissection ordering
	perm, iperm, err := metis.NodeND(xadj, adjncy, nil, opts)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Nested dissection ordering:")
	fmt.Printf("Permutation: %v\n", perm)
	fmt.Printf("Inverse permutation: %v\n", iperm)

	// Show the reordered vertices
	fmt.Print("Original order: ")
	for i := 0; i < 9; i++ {
		fmt.Printf("%d ", i)
	}
	fmt.Println()

	fmt.Print("New order:      ")
	for i := 0; i < 9; i++ {
		fmt.Printf("%d ", perm[i])
	}
	fmt.Println()

	// Compute vertex separator
	fmt.Println("\nVertex separator:")
	sepsize, part, err := metis.ComputeVertexSeparator(xadj, adjncy, nil, opts)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Separator size: %d\n", sepsize)
	fmt.Printf("Partition assignment (0/1/2=separator): %v\n", part)

	// Show which vertices are in the separator
	fmt.Print("Separator vertices: ")
	for i, p := range part {
		if p == 2 {
			fmt.Printf("%d ", i)
		}
	}
	fmt.Println()
}
