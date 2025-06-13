package main

import (
	"fmt"
	"github.com/notargets/go-metis"
	"log"
	"os"
)

func main() {
	fmt.Println("METIS Mesh Partitioning Example")
	fmt.Println("==============================")

	// Example: Partition a tetrahedral mesh
	tetrahedralMeshExample()

	// Example: Partition a hexahedral mesh
	hexahedralMeshExample()

	// Example: Convert mesh from Gambit format
	if len(os.Args) > 1 {
		gambitMeshExample(os.Args[1])
	}
}

func tetrahedralMeshExample() {
	fmt.Println("\nTetrahedral Mesh Partitioning")
	fmt.Println("-----------------------------")

	// Create a simple tetrahedral mesh
	// This represents a cube divided into 5 tetrahedra
	ne := int32(5) // Number of elements
	nn := int32(8) // Number of nodes (cube vertices)

	// Node coordinates (for visualization, not used in partitioning)
	// 0: (0,0,0), 1: (1,0,0), 2: (1,1,0), 3: (0,1,0)
	// 4: (0,0,1), 5: (1,0,1), 6: (1,1,1), 7: (0,1,1)

	// Element connectivity
	eptr := []int32{0, 4, 8, 12, 16, 20}
	eind := []int32{
		0, 1, 3, 7, // Tetrahedron 0
		1, 2, 3, 7, // Tetrahedron 1
		1, 5, 7, 2, // Tetrahedron 2
		5, 6, 7, 2, // Tetrahedron 3
		1, 5, 7, 4, // Tetrahedron 4
	}

	opts := make([]int32, metis.NoOptions)
	metis.SetDefaultOptions(opts)

	// Try different numbers of partitions
	for _, nparts := range []int32{2, 3, 4} {
		fmt.Printf("\nPartitioning into %d parts:\n", nparts)

		// Dual graph partitioning
		objval, epart, npart, err := metis.PartMeshDual(ne, nn, eptr, eind, nil, nil, 3, nparts, nil, opts)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("  Dual graph - Objective: %d\n", objval)
		fmt.Printf("  Element assignment: %v\n", epart)

		// Count elements per partition
		counts := make([]int, nparts)
		for _, p := range epart {
			counts[p]++
		}
		fmt.Printf("  Elements per partition: %v\n", counts)

		// Nodal graph partitioning
		objval, epart, npart, err = metis.PartMeshNodal(ne, nn, eptr, eind, nil, nil, nparts, nil, opts)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("  Nodal graph - Objective: %d\n", objval)
		fmt.Printf("  Element assignment: %v\n", epart)
	}
}

func hexahedralMeshExample() {
	fmt.Println("\n\nHexahedral Mesh Partitioning")
	fmt.Println("----------------------------")

	// Create a 2x2x2 hexahedral mesh
	ne := int32(8)  // 8 hex elements
	nn := int32(27) // 27 nodes (3x3x3 grid)

	// Element connectivity (8 nodes per hex)
	eptr := make([]int32, ne+1)
	for i := range eptr {
		eptr[i] = int32(i * 8)
	}

	// Define connectivity for 8 hexahedra
	// Numbering nodes in a 3x3x3 grid from 0 to 26
	eind := []int32{
		// Layer 0 (z=0)
		0, 1, 4, 3, 9, 10, 13, 12, // Hex 0
		1, 2, 5, 4, 10, 11, 14, 13, // Hex 1
		3, 4, 7, 6, 12, 13, 16, 15, // Hex 2
		4, 5, 8, 7, 13, 14, 17, 16, // Hex 3
		// Layer 1 (z=1)
		9, 10, 13, 12, 18, 19, 22, 21, // Hex 4
		10, 11, 14, 13, 19, 20, 23, 22, // Hex 5
		12, 13, 16, 15, 21, 22, 25, 24, // Hex 6
		13, 14, 17, 16, 22, 23, 26, 25, // Hex 7
	}

	// Add element weights (volume-based)
	vwgt := []int32{1, 1, 1, 1, 1, 1, 1, 1} // All elements have equal weight

	opts := make([]int32, metis.NoOptions)
	metis.SetDefaultOptions(opts)

	// Partition into 4 parts
	fmt.Println("\nPartitioning hexahedral mesh into 4 parts:")

	// With equal weights
	objval, epart, npart, err := metis.PartMeshDual(ne, nn, eptr, eind, vwgt, nil, 4, 4, nil, opts)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Objective value: %d\n", objval)
	fmt.Printf("Element partition: %v\n", epart)

	// Visualize the partitioning
	fmt.Println("\nPartition visualization (2x2x2 mesh):")
	fmt.Println("Layer 0 (bottom):")
	fmt.Printf("  [%d] [%d]\n", epart[0], epart[1])
	fmt.Printf("  [%d] [%d]\n", epart[2], epart[3])
	fmt.Println("Layer 1 (top):")
	fmt.Printf("  [%d] [%d]\n", epart[4], epart[5])
	fmt.Printf("  [%d] [%d]\n", epart[6], epart[7])
}

func gambitMeshExample(filename string) {
	fmt.Println("\n\nGambit Mesh File Partitioning")
	fmt.Println("-----------------------------")
	fmt.Printf("File: %s\n", filename)

	// This is a placeholder for reading Gambit mesh files
	// In practice, you would parse the .neu file format

	fmt.Println("\nTo partition a Gambit mesh:")
	fmt.Println("1. Parse the .neu file to extract elements and nodes")
	fmt.Println("2. Convert to METIS format (eptr, eind arrays)")
	fmt.Println("3. Call PartMeshDual or PartMeshNodal")
	fmt.Println("4. Write partitioning back to file or use in solver")

	// Example of how to use the mesh partitioning
	exampleCode := `
	// After parsing Gambit file
	mesh := ReadGambitFile(filename)
	
	// Convert to METIS format
	ne := int32(mesh.NumElements)
	nn := int32(mesh.NumVertices)
	eptr, eind := mesh.GetConnectivity()
	
	// Partition
	opts := make([]int32, metis.NoOptions)
	metis.SetDefaultOptions(opts)
	
	nparts := int32(4) // Number of partitions
	objval, epart, npart, err := metis.PartMeshDual(
		ne, nn, eptr, eind, nil, nil, 3, nparts, nil, opts)
	
	// Use partitioning
	for i, part := range epart {
		mesh.Elements[i].Partition = part
	}
	`

	fmt.Println(exampleCode)
}

// Helper function to analyze mesh partitioning quality
func analyzeMeshPartitioning(ne int32, eptr, eind, epart []int32, nparts int32) {
	fmt.Println("\nPartitioning Analysis:")

	// Count elements per partition
	counts := make([]int, nparts)
	for _, p := range epart {
		counts[p]++
	}

	fmt.Printf("Elements per partition: %v\n", counts)

	// Calculate load imbalance
	min, max := counts[0], counts[0]
	for _, c := range counts {
		if c < min {
			min = c
		}
		if c > max {
			max = c
		}
	}

	imbalance := float64(max) / (float64(ne) / float64(nparts))
	fmt.Printf("Load imbalance factor: %.2f\n", imbalance)

	// Count interface elements
	interfaceElements := 0
	for e := int32(0); e < ne; e++ {
		start := eptr[e]
		end := eptr[e+1]

		// Check if element has neighbors in different partitions
		hasInterface := false
		myPart := epart[e]

		// This would require the dual graph to properly check
		// For now, just count elements on partition boundaries
		for i := start; i < end; i++ {
			// Check adjacent elements (would need dual graph)
			_ = eind[i] // Node index
		}

		if hasInterface {
			interfaceElements++
		}
	}

	fmt.Printf("Interface elements: ~%d (estimate)\n", interfaceElements)

	// Calculate communication volume (simplified)
	// In practice, you'd use the dual graph to count shared faces
	fmt.Println("\nNote: For accurate interface/communication metrics,")
	fmt.Println("use METIS dual graph construction functions.")
}

// Utility function to create a simple test mesh
func createTestMesh() (ne, nn int32, eptr, eind []int32) {
	// Create a 2D quad mesh (4x4 grid = 16 quads, 25 nodes)
	// Convert to triangles (32 triangles)
	ne = 32
	nn = 25

	eptr = make([]int32, ne+1)
	eind = make([]int32, ne*3) // 3 nodes per triangle

	// Fill in the connectivity
	idx := int32(0)
	elemIdx := int32(0)

	// Convert each quad to 2 triangles
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			// Node indices for this quad
			n0 := int32(i*5 + j)
			n1 := n0 + 1
			n2 := n0 + 6
			n3 := n0 + 5

			// First triangle: n0, n1, n2
			eptr[elemIdx] = idx
			eind[idx] = n0
			eind[idx+1] = n1
			eind[idx+2] = n2
			idx += 3
			elemIdx++

			// Second triangle: n0, n2, n3
			eptr[elemIdx] = idx
			eind[idx] = n0
			eind[idx+1] = n2
			eind[idx+2] = n3
			idx += 3
			elemIdx++
		}
	}
	eptr[ne] = idx

	return
}

// Example of partitioning a mesh read from file
func partitionMeshFromFile(graph *metis.Graph, nparts int32) {
	opts := make([]int32, metis.NoOptions)
	metis.SetDefaultOptions(opts)

	// Partition the graph
	part, edgeCut, err := metis.PartGraphKway(graph.Xadj, graph.Adjncy, nparts, opts)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\nGraph partitioning results:\n")
	fmt.Printf("Number of vertices: %d\n", graph.NumVertices())
	fmt.Printf("Number of edges: %d\n", graph.NumEdges())
	fmt.Printf("Number of partitions: %d\n", nparts)
	fmt.Printf("Edge cut: %d\n", edgeCut)

	// Analyze partition quality
	min, max, avg := metis.CalculatePartitionBalance(part, graph.Vwgt, nparts)
	fmt.Printf("\nPartition sizes:\n")
	fmt.Printf("  Min: %.0f\n", min)
	fmt.Printf("  Max: %.0f\n", max)
	fmt.Printf("  Avg: %.1f\n", avg)
	fmt.Printf("  Balance: %.2f\n", max/avg)

	// Write partition file
	outFile := "partition.txt"
	f, err := os.Create(outFile)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	if err := metis.WritePartitioning(f, part); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\nPartitioning written to: %s\n", outFile)
}
