/*
Package metis provides Go bindings for METIS, a software package for partitioning
unstructured graphs, partitioning meshes, and computing fill-reducing orderings
of sparse matrices.

METIS is developed by the Karypis Lab at the University of Minnesota and is
widely used in scientific computing, particularly for parallel computing applications
where load balancing and minimizing communication overhead are critical.

# Overview

The package provides access to the main functionality of METIS:
  - Graph partitioning (recursive bisection and k-way)
  - Mesh partitioning (dual and nodal)
  - Nested dissection ordering for sparse matrices
  - Graph coarsening and refinement
  - Vertex separator computation

# Installation

Before using this package, you must have METIS installed on your system:

	# Ubuntu/Debian
	sudo apt-get install libmetis-dev

	# macOS
	brew install metis

	# From source
	git clone https://github.com/KarypisLab/GKlib.git
	git clone https://github.com/KarypisLab/METIS.git
	cd GKlib && make config prefix=/usr/local && sudo make install
	cd ../METIS && make config shared=1 prefix=/usr/local && sudo make install

Then install the Go package:

	go get github.com/Notargets/go-metis

# Basic Usage

Graph Partitioning

The most common use case is partitioning a graph for parallel processing:

	// Define a simple graph in CSR format
	// Graph: 0-1-2
	//        |X|X|
	//        3-4-5
	xadj := []int32{0, 2, 5, 7, 9, 12, 14}
	adjncy := []int32{1, 3, 0, 2, 4, 1, 5, 0, 4, 1, 3, 5, 2, 4}

	// Create options and use defaults
	opts := make([]int32, metis.NoOptions)
	metis.SetDefaultOptions(opts)

	// Partition into 2 parts
	part, edgeCut, err := metis.PartGraphKway(xadj, adjncy, 2, opts)
	if err != nil {
		log.Fatal(err)
	}

	// part[i] contains the partition assignment for vertex i
	fmt.Printf("Partition: %v, Edge cut: %d\n", part, edgeCut)

# Graph Format

Graphs are represented in Compressed Sparse Row (CSR) format:
  - xadj: Index array of size n+1, where n is the number of vertices
  - adjncy: Concatenated adjacency lists
  - xadj[i] points to the start of adjacency list for vertex i
  - Vertices are numbered from 0 (C-style)

Example for a triangle graph (0-1-2-0):

	xadj   = [0, 2, 4, 6]    // Vertex 0 starts at 0, vertex 1 at 2, etc.
	adjncy = [1, 2, 0, 2, 0, 1]  // Neighbors: 0->[1,2], 1->[0,2], 2->[0,1]

# Weighted Graphs

Both vertices and edges can have weights:

	vwgt := []int32{10, 20, 30, 40}      // Vertex weights
	adjwgt := []int32{1, 2, 1, 3, 3, 2}  // Edge weights

	part, edgeCut, err := metis.PartGraphKwayWeighted(
		xadj, adjncy, vwgt, adjwgt, nparts, nil, nil, opts)

# Mesh Partitioning

For finite element meshes, METIS provides specialized partitioning:

	// Define mesh connectivity
	ne := int32(4)  // Number of elements
	nn := int32(6)  // Number of nodes
	eptr := []int32{0, 3, 6, 9, 12}  // Element pointer
	eind := []int32{0, 1, 2, 1, 3, 2, 2, 3, 4, 3, 5, 4}  // Element-node list

	// Partition mesh using dual graph (element-based)
	objval, epart, npart, err := metis.PartMeshDual(
		ne, nn, eptr, eind, nil, nil, 3, 2, nil, opts)

# Options

METIS behavior can be controlled through the options array:

	opts := make([]int32, metis.NoOptions)
	metis.SetDefaultOptions(opts)

	// Set specific options
	opts[metis.OptionPType] = metis.PTypeKway      // Partitioning method
	opts[metis.OptionObjType] = metis.ObjTypeCut   // Minimize edge cut
	opts[metis.OptionNumBering] = 0                // C-style numbering
	opts[metis.OptionSeed] = 42                    // Random seed
	opts[metis.OptionDBGLvl] = 0                   // Debug level

# Partitioning Methods

Two main partitioning approaches are available:

1. Recursive Bisection: Recursively splits the graph in half
   - Better for small number of partitions (2-8)
   - Often produces better quality partitions
   - Use: PartGraphRecursive()

2. K-way Partitioning: Directly partitions into k parts
   - Better for large number of partitions (>8)
   - Generally faster
   - Use: PartGraphKway()

# Applications

Common use cases include:

1. Parallel Computing: Distribute computation across processors
   - Minimize communication (edge cut)
   - Balance computational load (vertex weights)

2. Finite Element Analysis: Partition meshes for parallel solvers
   - Element-based (dual) or node-based (nodal) partitioning
   - Minimize interface nodes/elements

3. Sparse Matrix Ordering: Reduce fill-in for direct solvers
   - Nested dissection ordering
   - Bandwidth/profile reduction

4. Graph Analytics: Process large graphs in parallel
   - Community detection preprocessing
   - Distributed graph algorithms

# Performance Considerations

1. Graph Size: METIS handles graphs with millions of vertices efficiently

2. Memory Usage: Approximately O(n + m) where n = vertices, m = edges

3. Time Complexity: O(m) for most algorithms

4. Quality vs Speed: Options allow trading partition quality for speed

# Error Handling

All functions return errors for invalid inputs or internal failures:

	part, edgeCut, err := metis.PartGraphKway(xadj, adjncy, nparts, opts)
	if err != nil {
		switch err {
		case metis.ErrorInput:
			// Invalid input parameters
		case metis.ErrorMemory:
			// Insufficient memory
		default:
			// Other errors
		}
	}

# Thread Safety

METIS functions are not thread-safe. Concurrent calls must be synchronized
externally. For parallel partitioning, create separate METIS instances or
use locking.

# References

For more information about METIS algorithms and options:
  - METIS Manual: http://glaros.dtc.umn.edu/gkhome/metis/metis/overview
  - Karypis Lab: http://glaros.dtc.umn.edu/gkhome/

Based on METIS version 5.1.0 by George Karypis and Vipin Kumar.
*/
package metis
