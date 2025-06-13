package metis

import (
	"math"
	"math/rand"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestVersion(t *testing.T) {
	v := Version()
	assert.NotEmpty(t, v)
	assert.Equal(t, "5.2.1", v) // Based on the header file
}

func TestSetDefaultOptions(t *testing.T) {
	opts := make([]int32, NoOptions)
	err := SetDefaultOptions(opts)
	require.NoError(t, err)

	// Check some default values
	assert.Equal(t, int32(0), opts[OptionNumbering]) // C-style numbering by default

	// Test with wrong size
	wrongOpts := make([]int32, 10)
	err = SetDefaultOptions(wrongOpts)
	assert.Error(t, err)
}

func TestPartGraphRecursive(t *testing.T) {
	// Simple 4-node square graph
	//  0 -- 1
	//  |    |
	//  2 -- 3
	xadj := []int32{0, 2, 4, 6, 8}
	adjncy := []int32{1, 2, 0, 3, 0, 3, 1, 2}

	opts := make([]int32, NoOptions)
	err := SetDefaultOptions(opts)
	require.NoError(t, err)

	// Test 2-way partition
	part, objval, err := PartGraphRecursive(xadj, adjncy, 2, opts)
	require.NoError(t, err)
	assert.Len(t, part, 4)
	assert.GreaterOrEqual(t, objval, int32(0))

	// Verify partition is valid
	counts := make(map[int32]int)
	for _, p := range part {
		assert.GreaterOrEqual(t, p, int32(0))
		assert.Less(t, p, int32(2))
		counts[p]++
	}
	assert.Len(t, counts, 2, "Should have exactly 2 partitions")

	// Test 3-way partition
	part, objval, err = PartGraphRecursive(xadj, adjncy, 3, opts)
	require.NoError(t, err)
	assert.Len(t, part, 4)

	// Test with nil options
	part, objval, err = PartGraphRecursive(xadj, adjncy, 2, nil)
	require.NoError(t, err)
	assert.Len(t, part, 4)
}

func TestPartGraphKway(t *testing.T) {
	// Simple 6-node graph (two triangles connected)
	//  0 -- 1     3 -- 4
	//   \  /       \  /
	//    2 -------- 5
	xadj := []int32{0, 2, 4, 7, 9, 11, 14}
	adjncy := []int32{1, 2, 0, 2, 0, 1, 5, 4, 5, 3, 5, 2, 3, 4}

	opts := make([]int32, NoOptions)
	err := SetDefaultOptions(opts)
	require.NoError(t, err)

	// Test 2-way partition
	part, objval, err := PartGraphKway(xadj, adjncy, 2, opts)
	require.NoError(t, err)
	assert.Len(t, part, 6)
	assert.GreaterOrEqual(t, objval, int32(0))

	// Verify partition
	for _, p := range part {
		assert.GreaterOrEqual(t, p, int32(0))
		assert.Less(t, p, int32(2))
	}

	// Test 3-way partition
	part, objval, err = PartGraphKway(xadj, adjncy, 3, opts)
	require.NoError(t, err)
	assert.Len(t, part, 6)
}

func TestPartGraphWeighted(t *testing.T) {
	// Weighted graph test
	xadj := []int32{0, 2, 5, 7, 9}
	adjncy := []int32{1, 3, 0, 2, 3, 1, 3, 0, 1, 2}
	vwgt := []int32{10, 20, 30, 40}                 // Vertex weights
	adjwgt := []int32{1, 2, 1, 3, 4, 3, 5, 2, 4, 5} // Edge weights

	opts := make([]int32, NoOptions)
	err := SetDefaultOptions(opts)
	require.NoError(t, err)

	// Test recursive with weights
	part, _, err := PartGraphRecursiveWeighted(xadj, adjncy, vwgt, adjwgt, 2, nil, nil, opts)
	require.NoError(t, err)
	assert.Len(t, part, 4)

	// Test k-way with weights
	part, _, err = PartGraphKwayWeighted(xadj, adjncy, vwgt, adjwgt, 2, nil, nil, opts)
	require.NoError(t, err)
	assert.Len(t, part, 4)

	// Test with target partition weights
	tpwgts := []float32{0.3, 0.7} // 30%/70% split
	part, _, err = PartGraphKwayWeighted(xadj, adjncy, vwgt, nil, 2, tpwgts, nil, opts)
	require.NoError(t, err)
	assert.Len(t, part, 4)
}

func TestMeshToDualAndNodal(t *testing.T) {
	// Simple mesh with 2 triangular elements sharing an edge
	// Elements: [0,1,2] and [1,3,2]
	ne := int32(2) // Number of elements
	nn := int32(4) // Number of nodes
	eptr := []int32{0, 3, 6}
	eind := []int32{0, 1, 2, 1, 3, 2}

	// Test MeshToDual
	xadj, _, err := MeshToDual(ne, nn, eptr, eind, 2)
	require.NoError(t, err)
	assert.Len(t, xadj, int(ne+1))
	assert.Equal(t, int32(0), xadj[0])

	// Test MeshToNodal
	xadj, _, err = MeshToNodal(ne, nn, eptr, eind)
	require.NoError(t, err)
	assert.Len(t, xadj, int(nn+1))
	assert.Equal(t, int32(0), xadj[0])
}

func TestPartMesh(t *testing.T) {
	// Simple tetrahedral mesh
	ne := int32(2) // 2 tetrahedra
	nn := int32(5) // 5 nodes
	eptr := []int32{0, 4, 8}
	eind := []int32{0, 1, 2, 3, 1, 2, 3, 4}

	opts := make([]int32, NoOptions)
	err := SetDefaultOptions(opts)
	require.NoError(t, err)

	// Test nodal partitioning
	objval, epart, npart, err := PartMeshNodal(ne, nn, eptr, eind, nil, nil, 2, nil, opts)
	require.NoError(t, err)
	assert.Len(t, epart, int(ne))
	assert.Len(t, npart, int(nn))
	assert.GreaterOrEqual(t, objval, int32(0))

	// Test dual partitioning
	objval, epart, npart, err = PartMeshDual(ne, nn, eptr, eind, nil, nil, 3, 2, nil, opts)
	require.NoError(t, err)
	assert.Len(t, epart, int(ne))
	assert.Len(t, npart, int(nn))
}

func TestNodeND(t *testing.T) {
	// Grid graph for testing nested dissection
	xadj := []int32{0, 2, 5, 7, 9, 12, 15, 18, 20}
	adjncy := []int32{1, 3, 0, 2, 4, 1, 5, 0, 4, 1, 3, 5, 2, 4, 6, 3, 5, 7, 4, 6}

	opts := make([]int32, NoOptions)
	err := SetDefaultOptions(opts)
	require.NoError(t, err)

	// Test without weights
	perm, iperm, err := NodeND(xadj, adjncy, nil, opts)
	require.NoError(t, err)
	assert.Len(t, perm, 8)
	assert.Len(t, iperm, 8)

	// Verify permutation validity
	for i := 0; i < 8; i++ {
		assert.Equal(t, int32(i), perm[iperm[i]], "perm[iperm[i]] should equal i")
		assert.Equal(t, int32(i), iperm[perm[i]], "iperm[perm[i]] should equal i")
	}

	// Test with vertex weights
	vwgt := []int32{1, 2, 3, 4, 5, 6, 7, 8}
	perm, iperm, err = NodeND(xadj, adjncy, vwgt, opts)
	require.NoError(t, err)
	assert.Len(t, perm, 8)
}

func TestComputeVertexSeparator(t *testing.T) {
	// Simple graph
	xadj := []int32{0, 2, 4, 6, 8}
	adjncy := []int32{1, 2, 0, 3, 0, 3, 1, 2}

	opts := make([]int32, NoOptions)
	err := SetDefaultOptions(opts)
	require.NoError(t, err)

	sepsize, part, err := ComputeVertexSeparator(xadj, adjncy, nil, opts)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, sepsize, int32(0))
	assert.Len(t, part, 4)

	// Count partitions (0, 1, or 2 for separator)
	counts := make(map[int32]int)
	for _, p := range part {
		assert.GreaterOrEqual(t, p, int32(0))
		assert.LessOrEqual(t, p, int32(2))
		counts[p]++
	}
}

// Test helper functions

func createRandomGraph(nvtxs int) ([]int32, []int32) {
	rand.Seed(42)
	edges := make(map[[2]int]bool)

	// Create a connected graph with random edges
	for i := 0; i < nvtxs-1; i++ {
		// Ensure connectivity
		edges[[2]int{i, i + 1}] = true
		edges[[2]int{i + 1, i}] = true
	}

	// Add random edges
	numExtraEdges := nvtxs + rand.Intn(nvtxs*2)
	for i := 0; i < numExtraEdges; i++ {
		u := rand.Intn(nvtxs)
		v := rand.Intn(nvtxs)
		if u != v {
			edges[[2]int{u, v}] = true
			edges[[2]int{v, u}] = true
		}
	}

	// Build adjacency lists
	adjList := make([][]int, nvtxs)
	for edge := range edges {
		adjList[edge[0]] = append(adjList[edge[0]], edge[1])
	}

	// Sort adjacency lists and build CSR format
	xadj := make([]int32, nvtxs+1)
	adjncy := []int32{}

	for i := 0; i < nvtxs; i++ {
		sort.Ints(adjList[i])
		xadj[i+1] = xadj[i] + int32(len(adjList[i]))
		for _, v := range adjList[i] {
			adjncy = append(adjncy, int32(v))
		}
	}

	return xadj, adjncy
}

func TestLargeGraph(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping large graph test in short mode")
	}

	// Create a larger random graph
	xadj, adjncy := createRandomGraph(100)

	opts := make([]int32, NoOptions)
	err := SetDefaultOptions(opts)
	require.NoError(t, err)

	// Test various partitioning methods
	for nparts := int32(2); nparts <= 8; nparts *= 2 {
		// Recursive
		part, objval, err := PartGraphRecursive(xadj, adjncy, nparts, opts)
		require.NoError(t, err)
		assert.Len(t, part, 100)
		verifyPartitioning(t, xadj, adjncy, part, nparts, objval)

		// K-way
		part, objval, err = PartGraphKway(xadj, adjncy, nparts, opts)
		require.NoError(t, err)
		assert.Len(t, part, 100)
		verifyPartitioning(t, xadj, adjncy, part, nparts, objval)
	}
}

func verifyPartitioning(t *testing.T, xadj, adjncy, part []int32, nparts, reportedObjval int32) {
	nvtxs := len(xadj) - 1

	// Check all vertices are assigned
	assert.Len(t, part, nvtxs)

	// Check partition numbers are valid
	partCounts := make([]int, nparts)
	for i, p := range part {
		assert.GreaterOrEqual(t, p, int32(0), "Vertex %d has negative partition", i)
		assert.Less(t, p, nparts, "Vertex %d has partition %d >= nparts %d", i, p, nparts)
		partCounts[p]++
	}

	// Check all partitions are used (for small graphs this might not be true)
	if nvtxs >= int(nparts)*2 {
		for i := int32(0); i < nparts; i++ {
			assert.Greater(t, partCounts[i], 0, "Partition %d is empty", i)
		}
	}

	// Calculate edge cut
	edgeCut := int32(0)
	for i := 0; i < nvtxs; i++ {
		for j := xadj[i]; j < xadj[i+1]; j++ {
			if part[i] != part[adjncy[j]] {
				edgeCut++
			}
		}
	}
	edgeCut /= 2 // Each edge counted twice

	// The reported objval should match our calculation
	assert.Equal(t, reportedObjval, edgeCut, "Reported edge cut doesn't match calculated")

	// Check load balance (should be reasonably balanced)
	minSize := math.MaxInt32
	maxSize := 0
	for _, count := range partCounts {
		if count < minSize {
			minSize = count
		}
		if count > maxSize {
			maxSize = count
		}
	}

	// Allow up to 10% imbalance
	avgSize := float64(nvtxs) / float64(nparts)
	maxAllowed := int(math.Ceil(avgSize * 1.1))
	assert.LessOrEqual(t, maxSize, maxAllowed, "Partition imbalance too high")
}
